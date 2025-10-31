# reconnecting

_Source: https://doc.photonengine.com/quantum/current/manual/game-session/reconnecting_

# Reconnecting

The following documentation will shed light into all aspects of reconnecting into a running Quantum session. The basic flow is implemented in the Quantum Menu that is shipped with the Quantum SDK.

The reconnection process consists of two main parts: How to get back into the Photon Realtime room and what to do the Quantum simulation?

## Detecting Disconnects

- The socket connection managed by Photon Realtime reports an error `RealtimeClient.CallbackMessage.ListenManual<OnDisconnectedMsg>(OnDisconnect)` with a different reason than `DisconnectCause.DisconnectByClientLogic`.
- The Quantum server detects an error and disconnects the client, register to `QuantumCallback.SubscribeManual<CallbackPluginDisconnect>(OnPluginDisconnect)` to catch these cases.
- Mobile app loses focus. In most cases the online session cannot be resumed and the client should perform a complete reconnection.

**Simulation Network Errors To Test Disconnects**

- In the Unity Editor just hit Play to stop and start the application.
- `RealtimeClient.SimulateConnectionLoss(true)` will stop sending and receiving and will result in a `DisconnectCause.ClientTimeout` disconnection after 10 seconds.
- Use an external network tool (for example clumsy) to block the game server ports.

```
Clumsy:
Filter  (udp.DstPort == 5056 or udp.SrcPort == 5056) or (tcp.DstPort == 4531 or tcp.SrcPort == 4531)
Drop    100%

```

## Photon Realtime Fast Reconnect

To return to the room the `MatchmakingExtensions` can be used to reconnect and rejoin:

C#

```csharp
RealtimeClient Client = await MatchmakingExtensions.ReconnectToRoomAsync(arguments);

```

The method will try to directly connect to the game server and rejoin the room using data from `MatchmakingReconnectInformation` on the `MatchmakingArguments` previously configured, such as region, appVersion, room name, and so on.

Rejoining a room will assign the same Photon `Actor` id to the client.

Rejoining a room can also be performed after reconnecting or connecting to the master server:

C#

```csharp
RealtimeClient.ReconnectToMaster()
// ..
public void IConnectionCallbacks.OnConnectedToMaster() {
    _client.OpReJoinRoom(roomName);
}

```

The rejoin operation only works if the client has not left the room, yet. (see next section PlayerTTL)

### Requirements: PlayerTTL

Clients inside a room are generally `active`. They become `inactive`..

- after a timeout of 10 seconds (by default) without answering the server
- after calling `RealtimeClient.Disconnect()`
- after calling `RealtimeClient.DisconnectAsync()`
- after calling `RealtimeClient.LeaveRoomAsync()`

Now, there are two options:

A) the room runs the player-left logic (`PlayerTTL` is 0, default)

B) the player is marked `inactive` and kept in that state for `PlayerTTL` milliseconds before running the leave-room routine. The `PlayerTTL` value needs to be explicitly set in the `RoomOptions` when creating the room. Usually 20 seconds (20000 ms) is a good value to start with.

Fast Reconnect will allow clients back into their room when they are still `active` (before the 10 second timeout which `Realtime.Client.OpJoinRoom()` does not allow) and `inactive` (during the PlayerTTL) timeout.

When the client rejoined successfully `IMatchmakingCallbacks.OnJoinedRoom()` is called.

The sample menu implements an options to check out (`QuantumMenuUIMain.RunReconnection()`).

### Requirements: RoomTTL (Waiting For Snapshots)

When the room detects that all clients are `inactive` it will close itself right away. To prevent that set `RoomOptions.EmptyRoomTTL`. This may be important when your room only has a small number of players and the probability that all of them have connection problems at the same time is given. Because there needs to be someone present to send a snapshot, this will **only** work reliably with custom server plugin and server-side snapshots.

Consider this case: In a two player online game one player is reconnecting or late-joining and waiting for a snapshot while the other player starts to have connections problems. The snapshot is never send and the player is stuck waiting.

This problem is handled by `SessionRunner.StartAsync` and `SessionRunner.WaitForStartAsync`taking into account the `SessionRunner.Arguments.StartGameTimeoutInSeconds`.

### Requirements: Photon UserId

[Photon Realtime: Lobby And Matchmaking \| UserIds And Friends](/realtime/current/lobby-and-matchmaking/userids-and-friends)

In Photon, a player is identified using a unique UserID. To return to the room using rejoin the **UserId has to be the same**. It does not matter if the UserId originally has been set explicitly or by Photon.

Once in the room, the UserId does not matter for Quantum as it uses a different id to identify players (see section Quantum CliendId).

The Photon UserId can be..

1. set by the client when connecting (`AuthenticationValues.UserId`)
2. if left empty it is set by Photon
3. or set by an external authentication service

To complete the background info about Photon ids:

- `Photon Actor Number` (also referred to as actor id) identifies the player in his current room and is assigned per room and only valid in that context. Clients leaving and joining back a room will get a new actor id. A successful OpRejoinRoom() or ReconnectAndRejoin() will retain the actor id. Quantum provides a way to backtrace the actor ids and match them to a player `Frame.PlayerToActorId(PlayerRef)`. But keep in mind, that they can change for player leaving and joining back (not rejoining).

- `Photon Nickname` is a Photon client property that gets propagated in the rooms to know a bit more about the other clients. Has nothing to do with Quantum.


### Possible Error: ReconnectAndRejoin Returns False

The current connection handler `RealtimeClient` is missing relevant data to perform a reconnect. Run your default connection sequence and try to join or rejoin the room in a regular way.

Also see section `Reconnecting After App Restart`.

### Possible Error: PlayerTTL Ran Out

When rejoining past the PlayerTTL timeout `ErrorCode.JoinFailedWithRejoinerNotFound` is thrown.

This also means that we are connected to the MasterServer and can join the room with `JoinRoomAsync()`.

### Possible Error: Authentication Token Timeout

The authentication ticket expires after 1 hour. It will be refreshed automatically before running out in the course of the Quantum game session ( [Photon Realtime: Encryption \| Token Refresh](/realtime/current/reference/encryption#token_refresh)). If your general game sessions are long and you want to support reconnecting players after around 20 minutes you need to handle this error. Resolution is to restart the default connection routine and try to join back into the room.

C#

```csharp
public void OnDisconnected(DisconnectCause cause) {
    switch (cause) {
        case DisconnectCause.AuthenticationTicketExpired:
        case DisconnectCause.InvalidAuthentication:
            // Restart with your default connection sequence
        break;

```

### Possible Error: Connection Still Unavailable

Of course the connection can still be obstructed or other errors can occur. In each case a `IConnectionCallbacks.OnDisconnected(DisconnectCause cause)` is called.

### Reconnecting After App Restart

The `MatchmakingReconnectInformation` object caches data relevant to the rejoining operation and that information can get lost when restarting the application.

In that case the connection has to be restarted from scratch while reusing the same `UserId`, `FixedRegion` and `AppVersion`. When arriving at the master server either `Rejoin()` or `Join()` back into the room.

Due to the lost connection cache, rejoining may fail with `ErrorCode.JoinFailedFoundActiveJoiner` because the server did not register the disconnect, yet (10 sec timeout). In this case retry until rejoining worked or another error occurs.

Saving the Photon `UserId` to PlayerPrefs can of course be replaced by custom authentication.

It is also possible to store and load a snapshot inside PlayerPrefs, which may be interesting for games with a very low player count. To store binary data in PlayerPrefs as a `string` use `base64` en- and decoding.

### Different Master Server

`ReconnectAndRejoin()` and `ReconnectToMaster()` both prevent a fringe case that would let the client end up on a different master server than before when connecting back via the cloud. Reasons are:

- There are multiple cluster for one app
- Master server has been replaced (rotated out)
- Best region ping has a new result

### Other Photon Realtime Topics

These features are not important for reconnection but are part of the demo menu sample so we might as well cover them here.

#### Best Region Summary

The region ping is forced from time to time but to be certain players are not stuck with a bad ping result it could be smart to implement invalidation so players are not stuck with a bad or wrong result forever (e.g. if ping is above threshold clear BestRegionSummary every other day). Also it could happen that players travel to other parts of the world where a new ping would be required to find the closest region.

#### AppVersion

On the demo menu sample players can choose the `AppVersion` via `QuantumMenuViewSettings`. The `AppVersion` we supply with the `AppSettings` will group the player bases for the same AppId into separate groups. Players connecting to the same AppId and different AppVersions will not find each other at all.

This is useful when running multiple game versions live as well as during development to prevent other clients (that have a different code base and would instantly desync the game) to join a game running by a developer.

### Further Readings

- [Photon Realtime: Analysing Disconnects \| Quick Rejoin (ReconnectAndRejoin)](/realtime/current/troubleshooting/analyzing-disconnects#quick-rejoin-reconnectandrejoin)
- [Photon Realtime: Known Issues \| Mobile Background Apps](/realtime/current/troubleshooting/known-issues#running-in-background)
- [Photon Realtime: .NET Client API \| LoadBalancingClient](https://doc-api.photonengine.com/en/dotnet/current/class_photon_1_1_realtime_1_1_load_balancing_client.html#a8d741259a7eceee18c72db9b469d2741)

## Reconnecting Into A Running Quantum Game

### Quantum ClientId

The `ClientId` is a secret between the client and the server. Other clients never know it. It is passed when starting the QuantumRunner.

C#

```csharp
var sessionRunnerArguments = new SessionRunner.Arguments {
        ClientId = Client.UserId,
        //Other arguments are needed
      };
var runner = (QuantumRunner)await SessionRunner.StartAsync(sessionRunnerArguments);

```

Independently of having joined as a new Photon room actor or having rejoined, reconnecting clients are identified by their `ClientId` and will be assigned to the same player index they previously had if the slot was not filled by another player in the meantime. In short: player must use the **same `ClientId`** when reconnecting.

Quantum will not let a client start the session while another `active` player with the same `ClientId` is inside the room and waits for the disconnect timeout (10 seconds):

```
DISCONNECTED: Error #5: Duplicate client id

```

This is why **`ReconnectAndRejoin()` is required** to recover from short term connection losses.

#### Further Readings

- [Quantum: Player Manual](/quantum/current/manual/player/player)

### Restarting The Quantum Session

After a disconnect the `QuantumRunner` and `DeterministicSession` are **not** usable any more and must be destroyed and recreated.

When the client either joined or re-joined back into the room that runs the Quantum game the `QuantumRunner` needs to be restarted. The simulation will be paused until the snapshot arrives from another client. Then will catch-up and sync to the most recent game time.

Rough outline:

- detect disconnect, destroy QuantumRunner
- reconnect and rejoin the room
- re-start Quantum by calling SessionRunner.StartAsync()

To stop and destroy the QuantumSession call:

C#

```csharp
QuantumRunner.ShutdownAll(true);

```

**Only** call this method with `immediate:true` when you are on the Unity main thread and **never** from inside a Quantum callback. Call with `immediate:false` or delay the call manually until it gets picked up from a Unity update call.

The demo menu sample demonstrates how starting a new game, late-joining a running game. In `QuantumMenuUIParty.ConnectAsync()` we detect that that game has already been started by evaluating the `ConnectResult`.

### EntityViews And UI

Late-joins and reconnection players put high demands on how flexible your game is constructed. It needs to support starting the game from any point in time and possibly reusing instantiated prefabs and UI as well as stopping and cleaning up the game at any possible moment. Side effects are high loading times, having unwanted VFX and animations in the new scene, being stuck in UI transitions, etc.

If you want to keep the `QuantumEntityViewUpdater` and the `QuantumEntityViews` alive to reuse them, they need to manually be stopped from being updated, re-match them with the new `QuantumGame` instance, subscribe to the new callbacks, etc.

On the other side the handling of Quantum is extremely simple: shutdown runner, start runner.

### Events

The client will not receive previous events that were raised before the player joined or rejoined. The game view should be able to fully initialize/reset itself by polling the current state of the simulation and use future events/polling to keep itself updated.

### SetPlayerData

Calling `QuantumGame.AddPlayer(RuntimePlayer data)` for reconnecting players is optional. It depends if your avatar setup logic in the simulation requires this.

### StartParameters.QuitBehaviour

When the Quantum shutdown sequence is being executed (QuantumRunner.ShutdownAll) the QuantumNetworkCommunicator class will optionally perform room leave operations or disconnect the LoadBalancing client. Set to `QuitBehaviour.None` on the QuantumRunner.StartParameters to handle it yourself.

### Late-Joining And Buddy-Snapshots

A Quantum game snapshot is a platform independent blob of data that contains the complete state of the game after a verified (all input has been received) tick. The Quantum simulation can be started from a snapshot and seamlessly continue from that state on.

A client can create its own snapshot when the simulation is still running (local snapshot), the snapshot can be requested from other clients (buddy snapshot) or it can be send down from a custom server plugin that runs the simulation.

Starting or restarting from snapshots is very handy and is provided turn-key by Quantum. Otherwise late-joining or reconnecting clients would have to start the game session from the very beginning and fast-forward through the input history send by the server which can render the client app useless until it caught up and also input history stored on the server is limited to ~10 minutes.

The buddy snapshot process is started automatically when any client is starting its `QuantumRunner` (no matter if the client is starting the session for the first time, late-joining or reconnecting). The session will be put into paused mode `DeterministicSession.IsPaused` and a snapshot will be requested. Successful late joins will log the following messages:

```
Waiting for snapshot. Clock paused.
Detected Resync. Verified tick: 6541

```

Buddy snapshots are requested for clients connecting 5 seconds after the initial start.

The server uses a load balancing mechanism to decide which client it will ask for a buddy snapshot to not overburden individual clients.

Errors during the snapshot process will be forwarded to the client using the `Disconnect` message (e.g. the snapshot waiting state will time out after 15 seconds):

| Name | Description |
| --- | --- |
| `Error #13: Snapshot request failed` | For a late-joining or rejoining client when requesting a snapshot there is no other client in the room/game that can send a buddy snapshot. |
| `Error #51: Snapshot download timeout` | The server was not able to send all required snapshot fragment within the default timeout of 20 seconds. |
| `Error #52: Snapshot upload timeout` | A requested buddy snapshot was not able to be uploaded within the default timeout of 10 seconds. |
| `Error #53: Snapshot upload error` | The uploaded buddy snapshot contained an error. |
| `Error #54: Snapshot upload disconnected` | Late-joining was interrupted due to the buddy snapshot uploading client disconnecting. |

There are a **few differences** when starting from a snapshot during the game starting routines:

- Instead of `CallbackGameStarted` the callback `CallbackGameResynced` is executed.
- `System.OnInit()` is called before the snapshot is received.

### Local Snapshots

As an optional reconnection strategy a local snapshot of the last verified tick can be saved and used when starting the new `QuantumRunner`. This works best when the anticipated time offline is small. Local snapshots are generally more bandwidth friendly and faster.

**Guidelines**

Quantum enforces tight limitations around the local snapshot acceptance timing, because starting from a snapshot that is too old can degrade the user experience.

By default local snapshots that are older than **10 seconds** are not accepted by the server and instead a buddy-snapshot is requested. The process works transparently and from the clients perspective the only difference is the received snapshot age.

For games that have a low user count (e.g. 1 vs 1) the chance that there is no other client online that can provide a buddy snapshot is high. These types of games usually require to work with the `EmptyRoomTTL` value and Quantum prolongs the local snapshot acceptance time to **`EmptyRoomTTL`** but to a maximum of **two minutes**.

**Workflow**

- Detect disconnect
- Take snapshot
- Shutdown QuantumRunner
- Fast Photon Reconnect
- restart Quantum with snapshot

#### Local Snapshot Snippet

To use this snippet create an empty scene, create a game object and place this script on it.

As a minimum set the `Map` and the `SimulationConfig` on the `RuntimeConfig`.

Add at least one entry in `RuntimePlayers`.

Add the `QuantumStats` prefab to see if a Quantum simulation runs.

Press `Connect` to start the online game. Press `Disconnect` to stop, wait a couple seconds and press `Reconnect` and see that it the session is already in progress and has a tick larger than 60.

C#

```csharp
namespace Quantum.Demo {
  using System;
  using System.Collections.Generic;
  using Photon.Deterministic;
  using Photon.Realtime;
  using UnityEngine;
  using UnityEngine.SceneManagement;
  /// <summary>
  /// A Unity script that demonstrates how to connect to a Quantum cloud and start a Quantum game session.
  /// </summary>
  public class QuantumSimpleReconnectionGUI : QuantumMonoBehaviour {
    /// <summary>
    /// The RuntimeConfig to use for the Quantum game session. The RuntimeConfig describes custom game properties.
    /// </summary>
    public RuntimeConfig RuntimeConfig;
    /// <summary>
    /// The RuntimePlayers to add to the Quantum game session. The RuntimePlayers describe individual custom player properties.
    /// </summary>
    public List<RuntimePlayer> RuntimePlayers;
    /// <summary>
    /// Room keep alive time.
    /// </summary>
    public int EmptyRoomTtlInSeconds = 20;
    RealtimeClient _client;
    string _loadedScene;
    QuantumReconnectInformation _reconnectInformation;
    int _disconnectedTick;
    byte[] _disconnectedFrame;
    bool CanReconnect => _reconnectInformation != null && _reconnectInformation.HasTimedOut == false;
    async void OnGUI() {
      if (_client != null && _client.IsConnectedAndReady) {
        if (GUI.Button(new Rect(10, 60, 160, 40), "Disconnect")) {
          await Stop();
        }
      } else {
        if (GUI.Button(new Rect(10, 60, 160, 40), CanReconnect ? "Reconnect" : "Connect")) {
          await Run();
        }
      }
    }
    async System.Threading.Tasks.Task Run() {
      var connectionArguments = new MatchmakingArguments {
        PhotonSettings = PhotonServerSettings.Global.AppSettings,
        PluginName = "QuantumPlugin",
        MaxPlayers = Quantum.Input.MAX_COUNT,
        // Keep the client connection object, it has cached authentication information
        NetworkClient = _client,
        // Keep an empty room open for a time
        EmptyRoomTtlInSeconds = EmptyRoomTtlInSeconds,
        // Set the stored reconnection information
        ReconnectInformation = _reconnectInformation,
        // Don't let random matchmaking get into this room
        IsRoomVisible = false
      };
      if (CanReconnect) {
        // Switch to reconnecting mode
        _client = await MatchmakingExtensions.ReconnectToRoomAsync(connectionArguments);
      } else {
        _client = await MatchmakingExtensions.ConnectToRoomAsync(connectionArguments);
        // Remove the disconnect information, it would break a new room
        _disconnectedTick = 0;
        _disconnectedFrame = null;
      }
      // Load the map if AutoLoadSceneFromMap is not set
      if (QuantumUnityDB.TryGetGlobalAsset(RuntimeConfig.SimulationConfig, out Quantum.SimulationConfig simulationConfigAsset)
        && simulationConfigAsset.AutoLoadSceneFromMap == SimulationConfig.AutoLoadSceneFromMapMode.Disabled) {
        if (QuantumUnityDB.TryGetGlobalAsset(RuntimeConfig.Map, out Quantum.Map map) == false) {
          throw new Exception("Map not found");
        }
        using (new ConnectionServiceScope(_client)) {
          await SceneManager.LoadSceneAsync(map.Scene, LoadSceneMode.Additive);
          SceneManager.SetActiveScene(SceneManager.GetSceneByName(map.Scene));
          _loadedScene = map.Scene;
        }
      }
      var sessionRunnerArguments = new SessionRunner.Arguments {
        RunnerFactory = QuantumRunnerUnityFactory.DefaultFactory,
        GameParameters = QuantumRunnerUnityFactory.CreateGameParameters,
        ClientId = _client.UserId,
        RuntimeConfig = new QuantumUnityJsonSerializer().CloneConfig(RuntimeConfig),
        SessionConfig = QuantumDeterministicSessionConfigAsset.DefaultConfig,
        GameMode = DeterministicGameMode.Multiplayer,
        PlayerCount = Quantum.Input.MAX_COUNT,
        Communicator = new QuantumNetworkCommunicator(_client),
        // Set the initial tick
        InitialTick = _disconnectedTick,
        // Set the serialized frame
        FrameData = _disconnectedFrame
      };
      // Add a player to the game
      var runner = (QuantumRunner)await SessionRunner.StartAsync(sessionRunnerArguments);
      for (int i = 0; i < RuntimePlayers.Count; i++) {
        runner.Game.AddPlayer(i, RuntimePlayers[i]);
      }
    }
    async System.Threading.Tasks.Task Stop() {
      // Save the serialized frame
      _disconnectedTick = QuantumRunner.DefaultGame.Frames.Verified.Number;
      _disconnectedFrame = QuantumRunner.DefaultGame.Frames.Verified.Serialize(DeterministicFrameSerializeMode.Serialize);
      // Save the reconnect information
      _reconnectInformation = new QuantumReconnectInformation();
      // Set the timeout to empty room ttl
      _reconnectInformation.Set(_client, TimeSpan.FromSeconds(EmptyRoomTtlInSeconds));
      if (string.IsNullOrEmpty(_loadedScene) == false) {
        // Unload a scene if we loaded one
        await SceneManager.UnloadSceneAsync(_loadedScene);
      }
      // Shutdown the runner
      if (QuantumRunner.Default != null) {
        await QuantumRunner.Default.ShutdownAsync();
      }
      // Make sure the client has disconnected
      await _client.DisconnectAsync();
    }
  }
}

```

Back to top

- [Detecting Disconnects](#detecting-disconnects)
- [Photon Realtime Fast Reconnect](#photon-realtime-fast-reconnect)

  - [Requirements: PlayerTTL](#requirements-playerttl)
  - [Requirements: RoomTTL (Waiting For Snapshots)](#requirements-roomttl-waiting-for-snapshots)
  - [Requirements: Photon UserId](#requirements-photon-userid)
  - [Possible Error: ReconnectAndRejoin Returns False](#possible-error-reconnectandrejoin-returns-false)
  - [Possible Error: PlayerTTL Ran Out](#possible-error-playerttl-ran-out)
  - [Possible Error: Authentication Token Timeout](#possible-error-authentication-token-timeout)
  - [Possible Error: Connection Still Unavailable](#possible-error-connection-still-unavailable)
  - [Reconnecting After App Restart](#reconnecting-after-app-restart)
  - [Different Master Server](#different-master-server)
  - [Other Photon Realtime Topics](#other-photon-realtime-topics)
  - [Further Readings](#further-readings)

- [Reconnecting Into A Running Quantum Game](#reconnecting-into-a-running-quantum-game)
  - [Quantum ClientId](#quantum-clientid)
  - [Restarting The Quantum Session](#restarting-the-quantum-session)
  - [EntityViews And UI](#entityviews-and-ui)
  - [Events](#events)
  - [SetPlayerData](#setplayerdata)
  - [StartParameters.QuitBehaviour](#startparameters.quitbehaviour)
  - [Late-Joining And Buddy-Snapshots](#late-joining-and-buddy-snapshots)
  - [Local Snapshots](#local-snapshots)