# async-extensions

_Source: https://doc.photonengine.com/quantum/current/manual/game-session/async-extensions_

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

An Async-version of most Realtime API is already implemented. When something is missing it should be easy to add locally or talk to the Photon team with a request. The code can be found in the ```
AsyncExtensions.cs
```

 file.

None of the async methods will start processing or sending before being awaited.

The ```
RealtimeClient
```

does not need to be updated (```
RealtimeClient.Service()
```

) while waiting for an async operation method to complete.

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

When working with ```
async
```

/```
await
```

in Unity there are a few specialties to consider:

_Contrary to .Net, using ```_
_await_
_```_

_from the Unity thread will always resume execution on the Unity thread._

- This makes using await for our purpose quite harmless in Unity. But it will cause multi-threading issues when used outside of Unity.

_New ```_
_Tasks_
_```_

_are ran on the thread pool (which in most cases is very undesirable in Unity) if not created with a custom TaskFactory based on the Unity SynchronizationContext._

- The global default ```
  AsyncConfig
  ```

   creates a ```
  TaskFactory
  ```

   that is internally used to create and continue the task.
- Follow the trail of ```
  AsyncConfig.InitForUnity()
  ```

  .

C#

```csharp
var taskFactory = new TaskFactory(
CancellationToken.None,
TaskCreationOptions.DenyChildAttach,
TaskContinuationOptions.DenyChildAttach \| TaskContinuationOptions.ExecuteSynchronously,
TaskScheduler.FromCurrentSynchronizationContext());

```

_Unity does not stop running ```_
_Tasks_
_```_

_when switching play mode._

- This is a headache. We solve this by using a global ```
  CancellationTokenSource
  ```

   inside the ```
  AsyncSetup
  ```

   class that gets triggered on play mode change callbacks (see ```
  AsyncSetup.Startup()
  ```

  ).
- All task and continuations created internally use either an explicit ```
  AsyncConfig
  ```

   passed as an argument or ```
  AsyncConfig.Global
  ```

  .

_Exceptions in ```_
_Tasks_
_```_

_can be suppressed by Unity under some circumstances._

- In short, use this pattern:

**```**
**public async void Update() {}**
**```**

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

e.g. ```
Task.Delay()
```

## Matchmaking Async Extensions

The Realtime Matchmaking extensions are combining the most common connection and reconnection logic into two comfortable extension methods for the ```
RealtimeClient
```

 class.

### ConnectToRoomAsync

C#

```csharp
public Task<RealtimeClient> ConnectToRoomAsync(MatchmakingArguments arguments)

```

```
ConnectToRoomAsync
```

will perform a couple things:

- Connect to the Photon Cloud using provided the ```
  PhotonSettings
  ```

   and ```
  AuthValues
  ```

- Perform simple matchmaking based on the configuration
  - Random matchmaking: ```
    RoomName:null
    ```

    , ```
    CanOnlyJoin:false
    ```

  - Join an existing room: ```
    RoomName:"room-name"
    ```

    , ```
    CanOnlyJoin:false
    ```

  - Join or create a room: ```
    RoomName:"room-name"
    ```

    , ```
    CanOnlyJoin:true
    ```

  - Use a typed lobby: ```
    Lobby:MyLobby
    ```

  - Use lobby properties: ```
    CustomLobbyProperties:MyLobbyProperties
    ```

These values have to be set: ```
PhotonSettings
```

, ```
MaxPlayer
```

, ```
PluginName
```

, ```
AuthValues
```

/ ```
UserId
```

### MatchmakingArguments

| Property | Type | Description |
| --- | --- | --- |
| ```<br>PhotonSettings<br>``` | ```<br>AppSettings<br>``` | The Photon AppSetting class containing information about the AppId and Photon server addresses. |
| ```<br>PlayerTtlInSeconds<br>``` | ```<br>int<br>``` | Player TTL, in seconds. |
| ```<br>EmptyRoomTtlInSeconds<br>``` | ```<br>int<br>``` | Empty room TTL, in seconds. |
| ```<br>RoomName<br>``` | ```<br>string<br>``` | Set a desired room name to create or join. If the RoomName is null, random matchmaking is used instead. |
| ```<br>MaxPlayers<br>``` | ```<br>int<br>``` | Max clients for the Photon room. 0 = unlimited. |
| ```<br>CanOnlyJoin<br>``` | ```<br>bool<br>``` | Configure if the connect request can also create rooms or if it only tries to join. |
| ```<br>CustomProperties<br>``` | ```<br>Hashtable<br>``` | Custom room properties that are configured as ```<br>EnterRoomArgs.RoomOptions.CustomRoomProperties<br>```<br>. |
| ```<br>CustomLobbyProperties<br>``` | ```<br>string\[\]<br>``` | List of room properties that are used for lobby matchmaking. Will be configured as ```<br>EnterRoomArgs.RoomOptions.CustomRoomPropertiesForLobby<br>```<br>. |
| ```<br>AsyncConfig<br>``` | ```<br>AsyncConfig<br>``` | Async configuration that include TaskFactory and global cancellation support. If null, then ```<br>AsyncConfig.Global<br>```<br> is used. |
| ```<br>NetworkClient<br>``` | ```<br>RealtimeClient<br>``` | Optionally provide a client object. If null, a new client object is created during the matchmaking process. |
| ```<br>AuthValues<br>``` | ```<br>AuthenticationValues<br>``` | Provide authentication values for the Photon server connection. Use this in conjunction with custom authentication. This field is created when ```<br>UserId<br>```<br> is set. |
| ```<br>PluginName<br>``` | ```<br>string<br>``` | Photon server plugin to connect to. |
| ```<br>ReconnectInformation<br>``` | ```<br>MatchmakingReconnectInformation<br>``` | Optional object to save and load reconnect information. |
| ```<br>Lobby<br>``` | ```<br>TypedLobby<br>``` | Optional Realtime lobby to use for matchmaking. |

### ReconnectToRoomAsync

C#

```csharp
public Task<RealtimeClient> ReconnectToRoomAsync(MatchmakingArguments arguments)

```

```
ReconnectToRoomAsync
```

will try to return to the previous room.

It will attempt to fast reconnect to a room (skipping master server) when the client object is in a reusable state (e.g. after a timeout). Otherwise, a complete connection sequence runs and attempts to rejoin the room.

If a previous connection for the user is not yet discarded by the server (e.g. during the 10 second timeout), a re-join fails. ```
ReconnectToRoomAsync
```

 will automatically cover this specific case and make several attempts to rejoin a room, before finally failing.

```
ReconnectToRoomAsync
```

can also be used after restarting the application or with a new client object. In those cases, the ```
arguments.ReconnectInformation
```

of type ```
MatchmakingReconnectInformation
```

must be set to provide the information to re-join a room.

In Quantum, the ```
QuantumReconnectInformation
```

can be used, which automatically saves the rejoin-information on the ```
PlayerPrefs
```

.

**Caveat**: Saving the ```
UserId
```

to ```
PlayerPrefs
```

is likely a security risk and should always be replaced with Custom Authentication before going into the reconnection logic after restarting the app.

The Quantum demo menu can be configured to check saved ```
QuantumReconnectInformation
```

on start up. Use ```
QuantumMenuUIMain.IsReconnectionCheckEnabled
```

to enable this.

Any ```
MatchmakingReconnectInformation
```

instance has a lifetime to prevent rejoin attempts for outdated sessions. Call ```
Set(client)
```

repeatedly during the online game to refresh this timeout. It is automatically called from the matchmaking extension methods after a successful connection or re-connection.

The virtual method ```
MatchmakingReconnectInformation.Set(client)
```

can be overwritten if needed.

C#

```csharp
virtual void Set(RealtimeClient client)

```

### MatchmakingReconnectInformation

| Property | Type | Description |
| --- | --- | --- |
| ```<br>Room<br>``` | ```<br>string<br>``` | The room name that the client was connected to. |
| ```<br>Region<br>``` | ```<br>string<br>``` | The region the client was connected to. |
| ```<br>AppVersion<br>``` | ```<br>string<br>``` | The app version used in the former connection. |
| ```<br>UserId<br>``` | ```<br>string<br>``` | The user id the client used to connect to the server. |
| ```<br>TimeoutInTicks<br>``` | ```<br>long<br>``` | The timeout after this information is considered to be unusable. Use the ```<br>Timeout<br>```<br> property to set and get this value. |

Back to top

- [Error Handling](#error-handling)
- [Unity And Async](#unity-and-async)
- [WebGL Requirements](#webgl-requirements)

- [Matchmaking Async Extensions](#matchmaking-async-extensions)
- [ConnectToRoomAsync](#connecttoroomasync)
- [MatchmakingArguments](#matchmakingarguments)
- [ReconnectToRoomAsync](#reconnecttoroomasync)
- [MatchmakingReconnectInformation](#matchmakingreconnectinformation)