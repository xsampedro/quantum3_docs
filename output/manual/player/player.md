# player

_Source: https://doc.photonengine.com/quantum/current/manual/player/player_

# Overview

## Introduction

Quantum is agnostic to the concept of player entities. All entities are the same in the eyes of the simulation. Therefore, when we refer to "the player" in this document, we mean _the player controlled entity_.

## Player Identification

A player can be identified in two ways:

- their player index
- their PlayerRef

### Player Index Assignment

The Quantum `player index` is assigned by the server based on the order in which the `QuantumGame.AddPlayer()` operations arrive. This is not to be confused with the _Photon/Actor Id_ which is based on order in players joined the _Photon room_.

It is not possible to set the "Desired Quantum Id".

**N.B.:** In the event of a disconnect, we guarantee the client gets the same player index **IF** it reconnects with the same _ClientId_; regardless of their Photon Id - `public static QuantumRunner StartGame(String clientId, Int32 playerCount, StartParameters param)`.

### Player Index vs PlayerRef

The `PlayerRef` is a wrapper for the `player index` in the Quantum ECS. The `PlayerRef` is 1-based, while `player index` starts at 0. The reason is that `default(PlayerRef)` will return a "null/invalid" player ref struct for convenience.

- default(PlayerRef), internally a 0, means NOBODY
- PlayerRef, internally 1, is the same as player index 0
- PlayerRef, internally 2, is the same as player index 1

Automatic cast operators convert an `Integer` to a `PlayerRef` and vice-versa.

C#

```csharp
// DeterministicCommand GetPlayerCommand(PlayerRef player);
for (int p = 0; p < f.PlayerCount; p++) {
  var command = f.GetPlayerCommand(p);
}

```

### PlayerSlot

Quantum supports running multiple local players within one game client. The 0-based `player slot ` refers to all local players.

Some API explicitly require choosing a player slot which would be `0` when only one local player is used.

Read more in the following chapter `Multiple Local Players`.

**Is Local Player**

Quantum offers to APIs in the View to check if a player is local:

- `QuantumRunner.Default.Game.Session.IsLocalPlayer(int player)`
- `QuantumRunner.Default.Game.PlayerIsLocal(PlayerRef playerRef)`

**Get Local Players**

`QuantumRunner.Default.Game.GetLocalPlayers()` returns an array that is unique on every client and represents the global player indices that the local client controls in the Quantum simulation.

`QuantumRunner.Default.Game.GetLocalPlayerSlots()` returns an array that is unique on every client and represents the player slots that the local client controls in the Quantum simulation.

- The methods returns one index if there is only one local player. Should several players be on the same local machine, then the arrays will have the length of the local player count.
- When rejoining the game the same player index will be assigned to the client when starting the session with the same Quantum `SessionRunner.Arguments.ClientId`.

### Photon Id

You can identify a player's corresponding Photon Id via the `Frame` API:

- `Frame.PlayerToActorId(PlayerRef player)` converts a Quantum PlayerRef to an ActorId (Photon client id); or,
- `Frame.ActorIdToAllPlayers(Int32 actorId)` the reverse process of the previous method.

**IMPORTANT:** The Photon Id is irrelevant to the Quantum simulation.

## Starting the Game

After calling `SessionRunner.Start()` the followings start protocol sequence will happen:

1. A StartRequest is send to server attaching the SessionConfig and RuntimeConfig
2. The request is received by the server, validated and a `SimulationStart` confirmation is send back to the client containing the chose configs. If the information attached to the request is not valid, the request will be refused by a `PluginDisconnect` message.
3. The `SimulationStart` game message is received by the client. Clients do not have a player, yet. They will all start as spectator only and cannot send inputs or commands.
4. _(OPTIONAL)_ \- In case of a late join, the client may receive a snapshot-resync.
5. `QuantumGame.AddPlayer()` is used to register a Quantum player and sending individual RuntimePlayer configs.
6. The server encodes successfully joined players into the input which will be processed by all clients at the same tick and call the `ISignalOnPlayerAdded` signals. Unsuccessful attempts are answered by `OnLocalPlayerAddFailed` callbacks.

![Config Sequence Diagram](https://doc.photonengine.com/docs/img/quantum/v3/manual/player/connection_flow.jpg)
Sequence Diagram

For more information on the configuration files involved, please refer to the [Configuration Files](/quantum/current/manual/config-files) manual.

## Add And Remove Players

Before a client can add players it has to wait until the game start has commenced. This can be done in two ways:

A) Using the async version of `SessionRunner.StartAsync` and await the return.

C#

```csharp
// this will return once the connection logic is complete (e.g. received snapshot if needed)
var runner = (QuantumRunner)await SessionRunner.StartAsync(sessionRunnerArguments);

// adding player to the online simulation
var runtimePlayer = new RuntimePlayer { PlayerNickname = &#34;whiskeyjack29&#34; };
runner.Game.AddPlayer(runtimePlayer);

```

B) Registering to the `CallbackGameStarted` like shown in `QuantumAddRuntimePlayers.cs` or `QuantumRunnerLocalDebug.cs` script.

C#

```csharp
  public class QuantumAddRuntimePlayers : QuantumMonoBehaviour {
    public RuntimePlayer[] Players;

    public void Awake() {
      QuantumCallback.Subscribe(this, (CallbackGameStarted c) => OnGameStarted(c.Game, c.IsResync), game => game == QuantumRunner.Default.Game);
    }

    public void OnGameStarted(QuantumGame game, bool isResync) {
      for (int i = 0; i < Players.Length; i++) {
        game.AddPlayer(i, Players[i]);
      }
    }
  }

```

AddPlayer() can only be called once per slot. The server will send a `AddPlayerFailed` protocol event (see callbacks below) on any errors.

Because the operation can invoke a webrequest there is a rate limitation to protect third-party backends.

The total number of players that a client can request can be limited by webhooks or the Photon dashboard parameter `MaxPlayerSlots`.

The following API sends player related operations:

C#

```csharp
class QuantumGame {
  // Sends player data to the server and request adding a player (using player slot 0)
  void AddPlayer(RuntimePlayer data);
  // Sends player data a certain for player slot to the server
  void AddPlayer(int playerSlot, RuntimePlayer data);
  // Request removing player slot 0
  void RemovePlayer();
  // Removing certain player float
  void RemovePlayer(int playerSlot);
  // Remove all players that belong to this client
  void RemoveAllPlayers();
}

```

The following Quantum Callbacks can be listed to and used in the view:

C#

```csharp
CallbackLocalPlayerAddConfirmed {
  public Frame Frame;
  public int PlayerSlot;
  public PlayerRef Player;
}

CallbackLocalPlayerAddFailed {
  public int PlayerSlot;
  public string Message;
}

CallbackLocalPlayerRemoveConfirmed {
  public Frame Frame;
  public int PlayerSlot;
  public PlayerRef Player;
}

CallbackLocalPlayerRemoveFailed {
  public int PlayerSlot;
  public string Message;
}

// for example
QuantumCallback.Subscribe(this, (CallbackLocalPlayerAddConfirmed c) => OnLocalPlayerAddConfirmed(c));
private void OnLocalPlayerAddConfirmed(CallbackLocalPlayerAddConfirmed c) { }

```

The following Quantum signals can be listed to and used in the simulation. The information is missing the player slot, as this is information only available for the local player.

C#

```csharp
// The first time that this player ref was assigned to a player at all.
// When firstTime is false the player ref is being reused by a different player.
ISignalOnPlayerAdded(Frame frame, PlayerRef player, bool firstTime)
ISignalOnPlayerRemoved
(Frame frame, PlayerRef player)

```

## PlayerConnectedSystem

The `PlayerConnectedSystem` is another way of tracking the online connection status of players using [Input & Connection Flags](/quantum/current/manual/player/input-flags). Add it to the `SystemsConfig` to make use of it.

To receive connection and disconnection Quantum system callbacks the `ISignalOnPlayerConnected` and `ISignalOnPlayerDisconnected` have to be implemented.

## RuntimePlayer

The class `RuntimePlayer` holds player specific information for example the selected character.

The class is Json serialized when sent to and received from the server.

`RuntimePlayer` is a partial class in the `quantum.code` project and custom implementations are to be done in `RuntimePlayer.User.cs`. Keep Json serialization compatibility in mind when adding data to it.

**Accessing At Runtime**

The `RuntimePlayer` asset associated with any player can be retrieved by querying `Frame.GetPlayerData()` with their `PlayerRef`.

C#

```csharp
public void OnPlayerDataSet(Frame frame, PlayerRef player){
  var data = frame.GetPlayerData(player);
}

```

## Initializing A Player Entity

The entity controlled by a player can be initialized at any point during the simulation. A common approach is to initialize it when the player connects(`ISignalOnPlayerConnected`) and / or the player data is received (`ISignalOnPlayerAdded`).

- `ISignalOnPlayerConnected`: The player entity can be initialized with whatever information is already available in the simulation or the asset database.
- `ISignalOnPlayerAdded`: The player entity can be initialized with the information associated with the player's `RuntimePlayer` specific information. This is convenient for things such as selected character model or inventory load-out.

## Simulation Vs View

First, a few clarifications:

1. From the simulation's perspective (Quantum), player controlled entity are entities with player input. It does not know of local or remote players.
2. From the view's perspective (Unity), we poll input from the players on the local client.

To recap, in the simulation there is no such thing as "local" or "remote" players; however, in the view a player is either "local" or it is not.

C#

```csharp
Photon.Deterministic.DeterministicGameMode.Local
Photon.Deterministic.DeterministicGameMode.Multiplayer
Photon.Deterministic.DeterministicGameMode.Replay

```

### Max Player Count

The max player count is essential to know in advance for it defines how much space needs to be allocated inside the memory chunk for each frame. By default the maximum amount of players is `6`.

To change it add the following lines to any of your `qtn`-files:

```
#define PLAYER_COUNT 8

#pragma max_players PLAYER_COUNT

```

- The `define` acts like a define and can be used inside the DSL (e.g. for allocating arrays with the player count).
- The `pragma` actually defines how many player the simulation can handle.

Back to top

- [Introduction](#introduction)
- [Player Identification](#player-identification)

  - [Player Index Assignment](#player-index-assignment)
  - [Player Index vs PlayerRef](#player-index-vs-playerref)
  - [PlayerSlot](#playerslot)
  - [Photon Id](#photon-id)

- [Starting the Game](#starting-the-game)
- [Add And Remove Players](#add-and-remove-players)
- [PlayerConnectedSystem](#playerconnectedsystem)
- [RuntimePlayer](#runtimeplayer)
- [Initializing A Player Entity](#initializing-a-player-entity)
- [Simulation Vs View](#simulation-vs-view)
  - [Max Player Count](#max-player-count)