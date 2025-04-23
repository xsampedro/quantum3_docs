# starting-session

_Source: https://doc.photonengine.com/quantum/current/manual/game-session/starting-session_

# Online Session

## Overview

The Quantum online services are build on top of the common Photon online infrastructure (Photon Realtime). Connecting to an online session usually goes through three connection phases:

1. **Custom Authentication**: Photon does not offer player accounts and recommends securing the logins using a proprietary or third-party authentication provider and set up Photon Custom Authentication.
2. **Game Server Connection**: Before starting the online simulation the clients have to connect to the Photon Cloud and enter a Photon Room using the Photon Realtime API.
3. **Quantum Simulation Start Sequence**: In this phase the Quantum simulation is started and synchronized and client configuration and player data is send.

## Game Server Connection

This is the simplest way to establish a connection to the Photon cloud and match-make onto a game server.

C#

```csharp
var connectionArguments = new MatchmakingArguments {
    // The Photon application settings include information about the app.
    PhotonSettings = PhotonServerSettings.Global.AppSettings,
    // The plugin to request from the Photon cloud set to the Quantum plugin.
    PluginName = "QuantumPlugin"
    // Setting an explicit room name will try to create or join the room based on how CanOnlyJoin is set. The RoomName can be null to create a unique name on create.
    RoomName = "My Room Name",
    // The maximum number of clients that can connect to the room, it most cases this is equal to the max number of players in the Quantum simulation.
    MaxPlayers = Input.MAX_COUNT,
    // Configure if the connect request can also create rooms or if it only tries to join
    CanOnlyJoin = false,
    // This sets the AuthValues and should be replaced with custom authentication and setting AuthValues explicitly
    UserId = Guid.NewGuid().ToString(),
};
// This line connects to the Photon cloud and performs matchmaking based on the arguments to finally enter a room.
RealtimeClient Client = await MatchmakingExtensions.ConnectToRoomAsync(connectionArguments);

```

It' a good practice to encapsulate `ConnectToRoomAsync()` and `ReconnectToRoomAsync()` procedures with `Try Catch`. Or let the exception bubble into other Tasks.

C#

```csharp
var client = default(RealtimeClient);
try {
    client = await MatchmakingExtensions.ConnectToRoomAsync(connectionArguments);
} catch (Exception e) {
    // Something unexpected happened.
    // In nearly all cases it's not worth to create detailed error handling. Log out the error, show generic feedback to the user and let him/her retry.
    Debug.LogException(e);
}

```

Examples for exception thrown during **ConnectToRoomAsync()** are:

- `Failed to connect and join with error 'X'` \- connecting to the Photon cloud failed with an error
- `MaxPlayer must be greater or equal than 0`
- `MaxPlayer must be less than 256`
- `PhotonSettings must be set`
- `TaskCanceledException` \- a cancellation was requested on `AsyncConfig.CancellationToken`
- ...

Examples for exception thrown during **ReconnectToRoomAsync()** are:

- `AppVersion mismatch` \- `ReconnectInformation.AppVersion` is different from `PhotonSettings.AppVersion`
- `UserId not set` \- `ReconnectInformation.UserId` is `null` or empty
- `UserId mismatch` \- `ReconnectInformation` and `AuthValues` have different `UserId`
- `ReconnectInformation timed out` \- `ReconnectInformation` is outdated
- `ReconnectInformation missing` \- `ReconnectInformation` is `null`
- ...

### AppSettings

To connect to the Photon Cloud the configuration called `AppSettings` is required. In the Unity project it can be accessed via a global asset: `PhotonServerSettings.Instance.AppSettings`.

Always copy the settings object to not risk saving changes to the `PhotonServerSettings` asset in Unity Editor.

C#

```csharp
var appSettings = new AppSettings(PhotonServerSettings.Global.AppSettings);

```

To set up a Photon AppId follow the instruction on: [Quantum Asteroids Tutorial - Project Setup](/quantum/current/tutorials/asteroids/2-project-setup)

### Optional Matchmaking Arguments

The `MatchmakingArguments` can be enhanced by these optional arguments.

Read more about the Photon Realtime matchmaking in the [Matchmaking Guide](/realtime/current/lobby-and-matchmaking/matchmaking-and-lobby).

|     |     |
| --- | --- |
| AsyncConfig | Async configuration that include `TaskFactory` and global cancellation support. If `null` then `AsyncConfig.Global` is used. |
| NetworkClient | Provide a Realtime client object. If `null` a new client object is created during the matchmaking process. |
| ReconnectInformation | Use this to support reconnection after losing the `NetworkClient` instance. For example after an application restart.<br> <br>`QuantumReconnectInformation` stores information on `PlayerPrefs` for example. See the Quantum demo menu for an example of how it can be used. |
| EmptyRoomTtlInSeconds | The time in seconds that an empty Photon room remains open after the last client left.<br> <br> Internally sets `RoomOptions.EmptyRoomTtl`. |
| PlayerTtlInSeconds | The time in seconds that disconnected clients will be marked as inactive on the server and support quick rejoining with `ReconnectAndRejoin()` for a longer period.<br> <br> Internally sets `RoomOptions.PlayerTtl`. |
| AuthValues | Provide authentication values for the Photon server connection. Use this together with custom authentication. <br> <br> This field is created when `UserId` is set. |
| CustomProperties | Set the room initial or expected room properties for the room. List the names of the properties used for matchmaking as `CustomLobbyProperties`. <br> <br> Internally sets `RoomOptions.CustomRoomProperties`. |
| CustomLobbyProperties | Name the custom room properties that should be available to clients that are in a lobby. This will separate clients inside the matchmaking and only match identical `CustomProperties` key-value pairs.<br> <br> Internally sets `RoomOptions.CustomRoomPropertiesForLobby`. |
| Lobby | The lobby in which to match. The type affects how filters are applied.<br> <br> Internally sets `EnterRoomArgs.Lobby` and `JoinRandomRoomArgs.Lobby`. |
| SqlLobbyFilter | SQL query to filter room matches. For default-typed lobbies, use ExpectedCustomRoomProperties instead.<br> <br> Internally sets `JoinRandomRoomArgs.SqlLobbyFilter`. |
| Ticket | Custom server signed matchmaking ticket.<br> <br> Internally sets `EnterRoomArgs.Ticket` and `JoinRandomRoomArgs.Ticket`. |
| RandomMatchingType | The `MatchmakingMode` affects how rooms get filled.<br> <br>`FillRoom` \- (Default) Fills up rooms (oldest first) to get players together as fast as possible.<br>`SerialMatching` -Distributes players across available rooms sequentially but takes filter into account. Without filter, rooms get players evenly distributed.<br>`RandomMatching` \- Joins a (fully) random room. Expected properties must match but aside from this, any available room might be selected.<br> <br> Internally sets `JoinRandomRoomArgs.MatchingType`. |
| ExpectedUsers | A list of users who are expected to join the room along with this client. Reserves slots for rooms with `MaxPlayers` value.<br> <br> Internally sets `EnterRoomArgs.ExpectedUser` and `JoinRandomRoomArgs.ExpectedUsers`. |
| CustomRoomOptions | Completely replaces the `EnterRoomArgs.RoomOptions` that would be composed from other arguments. |
| IsRoomVisible | If not `null` the initial value for `RoomOptions.IsVisible`. |
| IsRoomOpen | If not `null` the initial value for `RoomOptions.IsOpen`. |
| EnableCrc | While turned on, the client and server will add a CRC checksum to every sent package. The checksum enables both sides to detect and ignore packages that were corrupted during transfer. Corrupted packages have the same impact as lost packages: They require a re-send, adding a delay and could lead to timeouts. Building the checksum has a low processing overhead but increases integrity of sent and received data. Packages discarded due to failed CRC checks are counted in `PhotonPeer.PacketLossByCrc`. |

## Quantum Simulation Start Sequence

The snippet below shows the basic operation start an online Quantum session.

C#

```csharp
`var sessionRunnerArguments = new SessionRunner.Arguments {
    // The runner factory is the glue between the Quantum.Runner and Unity
    RunnerFactory = QuantumRunnerUnityFactory.DefaultFactory,
    // Creates a default version of `QuantumGameStartParameters`
    GameParameters = QuantumRunnerUnityFactory.CreateGameParameters,
    // A secret user id that is for example used to reserved player slots to reconnect into a running session
    ClientId = Client.UserId,
    // The player data
    RuntimeConfig = runtimeConfig,
    // The session config loaded from the Unity asset tagged as `QuantumDefaultGlobal`
    SessionConfig = QuantumDeterministicSessionConfigAsset.DefaultConfig,
    // GameMode has to be multiplayer for online sessions
    GameMode = DeterministicGameMode.Multiplayer,
    // The number of player that the session is running for, in this case we use the code-generated max possible players for the Quantum simulation
    PlayerCount = Input.MAX_COUNT,
    // A timeout to fail the connection logic and Quantum protocol
    StartGameTimeoutInSeconds = 10,
    // The communicator will take over the network handling after the simulation has started
    Communicator = new QuantumNetworkCommunicator(Client),
};
// This method completes when the client has successfully joined the online session
QuantumRunner runner = (QuantumRunner)await SessionRunner.StartAsync(sessionRunnerArguments);
`
```

## Adding and Removing Players

Quantum has a concept of a player. Each client can have zero or multiple players. When initially starting the Quantum online session the client is in the state of a spectator until Players are explicitly added.

Each player that is connected to the game is given a unique ID. This ID is called a `PlayerRef` and is often referred to as `Player`.

Unlike in Quantum 2.1 player can be added and removed at any time.

Players will occupy so called `PlayerSlots` which will always refer to how one client manages his players. If only one local player slot is used it will be slot `0`. While a second player that is controlled by that client could have the the `PlayerSlot` 1\. A typical usage for example is the input callback `CallbackPollInput` that uses the `PlayerSlot` property to control for what local player input is polled for: `QuantumGame.AddPlayer(Int32 playerSlot, RuntimePlayer data)`.

Because `AddPlayer()` could spawn a HTTP request to the customer backend the operation has a rate limitation on the server and cannot be spammed.

### RuntimePlayer

Quantum does not have a built-in concept of a player object or player avatar.

Any player related game information (such as a character load-out, character levels, etc.) is passed into the simulation by each clients using the `RuntimePlayer` object.

C#

```csharp
// Will add the player to player slot 0
QuantumRunner.Default.Game.AddPlayer(runtimePlayer);

```

To specify the exact player slot use `QuantumGame.AddPlayer(int playerSlot, RuntimePlayer data)`, as the example below.

C#

```csharp
// Will add the player to player slot 1
QuantumRunner.Default.Game.AddPlayer(1, runtimePlayer);

```

### PlayerConnectedSystem

To keep track of players' connection to a quantum session [Input & Connection Flags](/quantum/current/manual/player/input-flags) are used. The `PlayerConnectedSystem` automates the procedure and notifies the simulation if a player has connected to the session or disconnected from it.

In order to received the connection and disconnection callbacks the `ISignalOnPlayerConnected` and `ISignalOnPlayerDisconnected` have to be implemented in a system.

### Useful Local Player Callbacks

There is also some useful callbacks to handle with local players: `CallbackLocalPlayerAddConfirmed`, `CallbackLocalPlayerRemoveConfirmed`, `CallbackLocalPlayerAddFailed` and `CallbackLocalPlayerRemoveFailed`.

C#

```csharp
QuantumCallback.Subscribe(this, (CallbackLocalPlayerAddConfirmed c)    => OnLocalPlayerAddConfirmed(c));
QuantumCallback.Subscribe(this, (CallbackLocalPlayerRemoveConfirmed c) => OnLocalPlayerRemoveConfirmed(c));
QuantumCallback.Subscribe(this, (CallbackLocalPlayerAddFailed c)       => OnLocalPlayerAddFailed(c));
QuantumCallback.Subscribe(this, (CallbackLocalPlayerRemoveFailed c)    => OnLocalPlayerRemoveFailed(c));

```

## Stopping The Session And Disconnecting

To stop the Quantum simulation execute `QuantumRunner.ShutdownAll(bool immediate)`. Only set `immediate:true` when it's **not** called from within a Quantum callback. In case the shutdown command is called within a Quantum callback, is vital to set `immediate:false` so the shutdown is postponed until the next Unity update.

`ShutdownAll` will destroy the QuantumRunner object which triggers the local Quantum simulation to be stopped. It will also result in either a connection `Disconnect()` or `LeaveRoom()` depending what is set as `StartParameters.QuitBehaviour`.

If the client should exit the game gracefully, for example to clean up the player avatar for remote clients, extra logic has to be implemented into the simulation. Either a client issued command or monitoring the player connected state (see `PlayerConnectedSystem`).

Considering that players also close their app or Alt+F4 their games there might not always be an opportunity to send a graceful disconnect.

C#

```csharp
async void Disconnect() {
    // Signal all runners to shutdown and wait until each one has disconnected
    await QuantumRunner.ShutdownAllAsync();
    // OR just signal their shutdown
    QuantumRunner.ShutdownAll();
}

```

## Plugin Disconnected Errors

When the Quantum plugin encounters errors with the client start protocol or input messages it will send an operation to terminate the connection gracefully. When this happens the `CallbackPluginDisconnect` callback is invoked on the client, which includes a `Reason` string with more details. An error like this is not recoverable and the client would need to reconnect and restart the simulation.

|     |     |     |
| --- | --- | --- |
| Error #3 | Must request start before any protocol message or input messages other than 'StartRequest' is accepted | The client tried to send a protocol or input message before sending the start request. |
| Error #5 | Duplicate client id | The client tried to start the game with a `ClientId` that is already used. |
| Error #7 | Client protocol version '2.2.0.0' is not matching server protocol version '3.0.0.0' | The client tried to start the game with an incompatible protocol version. |
| Error #8 | Invalid client id 'NULL' | The client tried to start the game without specifying a `ClientId`. |
| Error #9 | Server refused client | The custom plugin refused the client to join the session. |
| Error #12 | Operation not allowed when running in spectating mode | The client tried to send a command while in spectator mode (no player added). |
| Error #13 | Snapshot request failed to start | The client late-joined a game but there is no suitable connected player to provide a snapshot to join. |
| Error #16 | Player corrupted | The client caused an exception on the plugin during either protocol or input message deserialization. The error can be tackled by enabling packet CRC checksums on the client connection: `RealtimeClient.RealtimePeer.CrcEnabled = true`. |
| Error #17 | PlayerCount is not valid | The client started the online game with an invalid player count. |
| Error #19 | Player not found | The client send a command for a `PlayerSlot` that they don't own. |
| Error #20 | RPC data corrupted | The `RuntimePlayer` object when adding a player or the command data was too large (max. 24 KB). |
| Error #21 | Quantum SDK 2 not supported on Quantum 3 AppIds, check Photon dashboard to set correct version | The client using Quantum SDK 2.1 tried to connect with a Quantum 3 AppId. |
| Error #33 | Player data was rejected 'Webhook Error Message' | The add player webhook failed. |
| Error #34 | Game configs not loaded 'Webhook Error Message' | The game configs webhook failed, disconnecting all clients. |
| Error #40 | Caught exception receiving protocol messages | An exception was raised when the client processed a protocol message, this is not necessarily a server error, but the client state is unrecoverable and it will disconnect. |
| Error #41 | Input cache full | The local delta compression input buffer for the client is full, this could happen after long breakpoint pauses for example. This state is not recoverable for the client. |
| Error #42 | Communicator not connected | The connection was lost while the simulation was running. This can be prevented by detecting disconnects earlier using Photon Realtime callbacks for example. |
| Error #51 | Snapshot download timeout | The server was not able to send all required snapshot fragment within the default timeout of 20 seconds. |
| Error #52 | Snapshot upload timeout | A requested buddy snapshot was not able to be uploaded within the default timeout of 10 seconds. |
| Error #53 | Snapshot upload error | The uploaded buddy snapshot contained an error. |
| Error #54 | Snapshot upload disconnected | Late-joining was interrupted due to the buddy snapshot uploading client disconnecting. |

Back to top

- [Overview](#overview)
- [Game Server Connection](#game-server-connection)

  - [AppSettings](#appsettings)
  - [Optional Matchmaking Arguments](#optional-matchmaking-arguments)

- [Quantum Simulation Start Sequence](#quantum-simulation-start-sequence)
- [Adding and Removing Players](#adding-and-removing-players)

  - [RuntimePlayer](#runtimeplayer)
  - [PlayerConnectedSystem](#playerconnectedsystem)
  - [Useful Local Player Callbacks](#useful-local-player-callbacks)

- [Stopping The Session And Disconnecting](#stopping-the-session-and-disconnecting)
- [Plugin Disconnected Errors](#plugin-disconnected-errors)