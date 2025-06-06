# async-extensions

_Source: https://doc.photonengine.com/realtime/current/reference/async-extensions_

# Async Extensions

The Async Extensions for the Realtime API are built as an optional model on top of the well-known callback approach and have been added to the .Net Realtime SDKs v5 and up.

These async methods are based on C# async/await and Tasks (TAP) which are supported on Unity for a while now (for more information, read the TAP documentation from Microsoft).

The implementation of these async APIs is not optimized for performance but for ease of use. As they are not involved in the gameplay loop (e.g. Quantum simulation), this should be fine.

Establishing a connection and sending operations is, by itself, an asynchronous process and Task-based asynchronous patterns produce much more readable and maintainable code.

C#

```csharp
var appSettings = new new AppSettings();
var client = new RealtimeClient();
await client.ConnectUsingSettingsAsync(appSettings);
var joinRandomRoomParams = new JoinRandomRoomArgs();
var enterRoomArgs = new EnterRoomArgs();
var result = await client.JoinRandomOrCreateRoomAsync(joinRandomRoomParams, enterRoomArgs);

```

An Async-version of most Realtime API is already implemented. When something is missing it should be easy to add locally or talk to the Photon team with a request. The code can be found in the `AsyncExtensions.cs` file.

None of the async methods will start processing or sending before being awaited.

The `RealtimeClient` does not need to be updated (`RealtimeClient.Service()`) while waiting for an async operation method to complete.

### Error Handling

**All** Async-methods will throw exceptions when encountering errors. Different types of exceptions are thrown to improve further error handling (see the related method summary).

It's a bit unfortunate to pollute the code with try/catch blocks but it makes the API much simpler.

C#

```csharp
try {
  await client.ConnectUsingSettingsAsync(appSettings);
} catch (Exception e) {
  Debug.LogException(e);
}

```

C#

```csharp
try {
  // Disconnecting can also fail
  await client.DisconnectAsync();
} catch (Exception e) {
  Debug.LogException(e);
}

```

### Unity And Async

When working with `async`/`await` in Unity there are a few specialties to consider:

_Contrary to .Net, using `await` from the Unity thread will always resume execution on the Unity thread._

- This makes using await for our purpose quite harmless in Unity. But it will cause multi-threading issues when used outside of Unity.

_New `Tasks` are ran on the thread pool (which in most cases is very undesirable in Unity) if not created with a custom TaskFactory based on the Unity SynchronizationContext._

- The global default `AsyncConfig` creates a `TaskFactory` that is internally used to create and continue the task.
- Follow the trail of `AsyncConfig.InitForUnity()`.

C#

```csharp
var taskFactory = new TaskFactory(
  CancellationToken.None,
  TaskCreationOptions.DenyChildAttach,
  TaskContinuationOptions.DenyChildAttach | TaskContinuationOptions.ExecuteSynchronously,
  TaskScheduler.FromCurrentSynchronizationContext());

```

_Unity does not stop running `Tasks` when switching play mode._

- This is a headache. We solve this by using a global `CancellationTokenSource` inside the `AsyncSetup` class that gets triggered on play mode change callbacks (see `AsyncSetup.Startup()`).
- All task and continuations created internally use either an explicit `AsyncConfig` passed as an argument or `AsyncConfig.Global`.

_Exceptions in `Tasks` can be suppressed by Unity under some circumstances._

- In short, use this pattern:

**`public async void Update() {}`**

C#

```csharp
// Does NOT log exception.
// Why? Because Unity does not handle exception inside tasks by design.
public Task Update1() {
  return Task.Run(() => throw new Exception("peng"));
}
​
// Does NOT log exception.
// Why? Because we return at await and continue as a task object and Unity swallows the exception.
public async Task Update3() {
  await Task.Delay(100);
  throw new Exception("peng");
}
​
// Logs exception.
// Why? because we unwrap the task and run it synchronously with .Wait().
public void Update2() {
  Task.Run(() => throw new Exception("peng")).Wait();
}
// Logs exception.
// Why? Because we resume the execution in this method and not return a task.
public async void Update4() {
  await Task.Delay(100);
  throw new Exception("peng");
}
​
// Logs exception.
// Why? We add a continuation task that logs (in any thread) when the task faulted.
public Task Update5() {
  var task = Task.Run(() => throw new Exception("peng")).ContinueWith(t => {
    if (t.IsFaulted) {
      Debug.LogException(t.Exception.Flatten().InnerException);
    };
  });
​
  return task;
}

```

### WebGL Requirements

The Realtime Async extensions are supported for WebGL. Because of threading restrictions in the browsers, multi-threading code can not be used:

e.g. `Task.Delay()`

## Matchmaking Async Extensions

The Realtime Matchmaking extensions are combining the most common connection and reconnection logic into two comfortable extension methods for the `RealtimeClient` class.

### ConnectToRoomAsync

C#

```csharp
public Task<RealtimeClient> ConnectToRoomAsync(MatchmakingArguments arguments)

```

`ConnectToRoomAsync` will perform a couple things:

- Connect to the Photon Cloud using provided the `PhotonSettings` and `AuthValues`
- Perform simple matchmaking based on the configuration
  - Random matchmaking: `RoomName:null`, `CanOnlyJoin:false`
  - Join an existing room: `RoomName:"room-name"`, `CanOnlyJoin:false`
  - Join or create a room: `RoomName:"room-name"`, `CanOnlyJoin:true`
  - Use a typed lobby: `Lobby:MyLobby`
  - Use lobby properties: `CustomLobbyProperties:MyLobbyProperties`

These values have to be set: `PhotonSettings`, `MaxPlayer`, `PluginName`, `AuthValues` / `UserId`

### MatchmakingArguments

| Property | Type | Description |
| --- | --- | --- |
| `PhotonSettings` | `AppSettings` | The Photon AppSetting class containing information about the AppId and Photon server addresses. |
| `PlayerTtlInSeconds` | `int` | Player TTL, in seconds. |
| `EmptyRoomTtlInSeconds` | `int` | Empty room TTL, in seconds. |
| `RoomName` | `string` | Set a desired room name to create or join. If the RoomName is null, random matchmaking is used instead. |
| `MaxPlayers` | `int` | Max clients for the Photon room. 0 = unlimited. |
| `CanOnlyJoin` | `bool` | Configure if the connect request can also create rooms or if it only tries to join. |
| `CustomProperties` | `Hashtable` | Custom room properties that are configured as `EnterRoomArgs.RoomOptions.CustomRoomProperties`. |
| `CustomLobbyProperties` | `string\[\]` | List of room properties that are used for lobby matchmaking. Will be configured as `EnterRoomArgs.RoomOptions.CustomRoomPropertiesForLobby`. |
| `AsyncConfig` | `AsyncConfig` | Async configuration that include TaskFactory and global cancellation support. If null, then `AsyncConfig.Global` is used. |
| `NetworkClient` | `RealtimeClient` | Optionally provide a client object. If null, a new client object is created during the matchmaking process. |
| `AuthValues` | `AuthenticationValues` | Provide authentication values for the Photon server connection. Use this in conjunction with custom authentication. This field is created when `UserId` is set. |
| `PluginName` | `string` | Photon server plugin to connect to. |
| `ReconnectInformation` | `MatchmakingReconnectInformation` | Optional object to save and load reconnect information. |
| `Lobby` | `TypedLobby` | Optional Realtime lobby to use for matchmaking. |

### ReconnectToRoomAsync

C#

```csharp
public Task<RealtimeClient> ReconnectToRoomAsync(MatchmakingArguments arguments)

```

`ReconnectToRoomAsync` will try to return to the previous room.

It will attempt to fast reconnect to a room (skipping master server) when the client object is in a reusable state (e.g. after a timeout). Otherwise, a complete connection sequence runs and attempts to rejoin the room.

If a previous connection for the user is not yet discarded by the server (e.g. during the 10 second timeout), a re-join fails. `ReconnectToRoomAsync` will automatically cover this specific case and make several attempts to rejoin a room, before finally failing.

`ReconnectToRoomAsync` can also be used after restarting the application or with a new client object. In those cases, the `arguments.ReconnectInformation` of type `MatchmakingReconnectInformation` must be set to provide the information to re-join a room.

Any `MatchmakingReconnectInformation` instance has a lifetime to prevent rejoin attempts for outdated sessions. Call `Set(client)` repeatedly during the online game to refresh this timeout. It is automatically called from the matchmaking extension methods after a successful connection or re-connection.

The virtual method `MatchmakingReconnectInformation.Set(client)` can be overwritten if needed.

C#

```csharp
virtual void Set(RealtimeClient client)

```

### MatchmakingReconnectInformation

| Property | Type | Description |
| --- | --- | --- |
| `Room` | `string` | The room name that the client was connected to. |
| `Region` | `string` | The region the client was connected to. |
| `AppVersion` | `string` | The app version used in the former connection. |
| `UserId` | `string` | The user id the client used to connect to the server. |
| `TimeoutInTicks` | `long` | The timeout after this information is considered to be unusable. Use the `Timeout` property to set and get this value. |

Back to top

- [Error Handling](#error-handling)
- [Unity And Async](#unity-and-async)
- [WebGL Requirements](#webgl-requirements)

- [Matchmaking Async Extensions](#matchmaking-async-extensions)
- [ConnectToRoomAsync](#connecttoroomasync)
- [MatchmakingArguments](#matchmakingarguments)
- [ReconnectToRoomAsync](#reconnecttoroomasync)
- [MatchmakingReconnectInformation](#matchmakingreconnectinformation)