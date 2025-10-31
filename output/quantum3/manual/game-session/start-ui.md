# start-ui

_Source: https://doc.photonengine.com/quantum/current/manual/game-session/start-ui_

# Start UI

The Quantum `StartUI` consist of a `Prefab` and a few Unity scripts that make a Quantum scene run-able in local and online mode.

It can be used as a simple connection menu, it's easy to add small variations to and it should be used as reference to start a custom game menu in each developers own style.

![](/docs/img/quantum/v3/manual/game-session/start-ui.png)## Use the StartUI

- open a Unity scene that includes a Quantum `Map`
- run the Unity Editor menu command `Tools > Quantum > Setup > Add Start UI To Current Scene`

A new GameObject is created with a pre-configured Unity `Canvas`. The `QuantumStartUI` prefab is instantiated as a child game object.

The prefab receives one variation which is saved on the scene: the `QuantumStartUIConnection` script. It contains the actual Photon cloud connection logic and can start the Quantum simulation using Quantum configs like `PhotonServerSettings`, `RuntimeConfig` and `RuntimePlayer`. When adding the UI via the Unity Editor menu script the current scene is tested for a `QuantumRunnerLocalDebug` script to pre-configure `QuantumStartUIConnection` members.

![](/docs/img/quantum/v3/manual/game-session/start-ui-connection.png)## Hierarchy and Classes

The `StartUI` is built with two main classes:

- `QuantumStartUI` controls the UI
- `QuantumStartUIConnection` starts and shuts down the connection and Quantum simulation

### QuantumStartUI

This script contains the UI bindings and UI state machine. The bindings and control registrations are done verbosely in code and there are no hidden settings in the Unity inspector.

The class has a subclass `UIBindings` to make the inspector have a foldout. _Elements_ are the actual UI controls and _Toggles_ are mostly on/off switches. Some of the elements can be set to `null` to disable their feature which is mentioned in the fields code comment.

![](/docs/img/quantum/v3/manual/game-session/start-ui-bindings.png)

Possible UI states are `Idle`, `Starting`, `Running` and `ShuttingDown`. By default the complete scene is reloaded after shutting down, this can be disabled by toggling off `ReloadSceneAfterShutdown`.

It has a rudimentary tab system. Some tab names are unused like `Settings` and `Custom`. Getting the tabs and a few other toggles to work without additional wrappers uses a small trick to toggle two states:

C#

```cs
UI.InputTabs[(int)Tab.Online].transform.Find("On").gameObject.SetActive(tab == Tab.Online);

```

All of the members are `protected` and all of the methods are `protected virtual` to be easily extendable to include small variations quickly.

When specifying a room name a toggle called "private" appears. This toggle will set the room to `IsVisible` false preventing the random matchmaking to add players. Instead other player have to enter the same room name and region to join.

### QuantumStartUIRegionDropdown

A wrapper around the Unity `Dropdown` script that downloads available Photon regions when being clicked.

The class has a partial implementation with the actual connection logic because of internal code organization.

### QuantumStartUIConnection

The class performs the actual connection, session starting and shutdown logic. It implements the abstract class `QuantumStartUIConnectionBase`, and its main functionality is inside

- `Task ConnectAsync(StartParameter startParameter)`and
- `Task DisconnectAsync()` which are both called by the UI class `QuantumStartUI`.

The menu uses `Task` and `async` keywords. Read more about them on this page [Starting Session > Async Extensions](/quantum/current/game-session/async-extensions). The methods are not strictly pure async because they are not enqueued and perform actions before synchronizing using the `await` keyword, but for simplicity this is omitted.

Online and local game modes are supported.

`ConnectAsync()` performs four steps:

1. Connect to the Photon cloud and enter a Photon room: `MatchmakingExtensions.ConnectToRoomAsync()` (online mode only)
2. Subscribe to server error channel: `QuantumCallback.SubscribeManual<CallbackPluginDisconnect>` (online mode only)
3. Configure the game configs and start the Quantum simulation: `SessionRunner.StartAsync()`
4. Add the players to the game: `Runner.Game.AddPlayer()`

`DisconnectAsync()` is intended to be used at any time - when still connecting or when the online session has been started and the game is running.

It uses a `CancellationTokenSource` to signal the async methods `MatchmakingExtensions.ConnectToRoomAsync()` and `SessionRunner.StartAsync()` to cancel.

Afterwards the `CallbackPluginDisconnect` is disposed and runner and client are shutdown and disconnected.

Read more about connecting and starting the Quantum session:

- [Starting Session > Game Server Connection](/quantum/current/game-session/starting-session#game-server-connection)
- [Starting Session > Quantum Simulation Start Sequence](/quantum/current/game-session/starting-session#quantum-simulation-start-sequence)

To catch errors that occur between requesting the online Quantum simulation to start and finally starting locally and to let the UI know about them and also runtime errors that the Quantum server may invoke on the client the process subscribes to the `CallbackPluginDisconnect` callback before calling `SessionRunner.StartAsync()`.

C#

```cs
_pluginDisconnectSubscription = QuantumCallback.SubscribeManual<CallbackPluginDisconnect>(m =>
    startParameter.OnConnectionError(m.Reason)
);

```

Errors are reported back to the UI using the `OnConnectionError` action on the start parameter. The subscription is disposed when the game and connection are shut down.

#### App Version

The `AppVersion` used when connecting to the cloud controls which players can be matched against each other. Quantum is very sensitive to any asset or code differences between builds or Editors. By default the `AppVersionMachineId` is selected which generates a machine and checkout specific unique id. Only players using builds from this machine will be matched together.

Use the `AppVersionMachineIdPostfix` to additionally separate the players by the map name.

Use `AppVersionOverride` to set an explicit app version (requires AppVersionMachineId to be `null`).

### QuantumStartUIMppmCommand

The Quantum `MppmCommands` are used to communicate between Editor instances using Unity 6 Multiplayer Play Mode. It will synchronize clicking start in one of the instances and makes sure that subsequent Editor instances automatically start and get into the same game.

The `QuantumStartUIMppmConnectCommand` is used inside the `QuantumStartUI` class. Invoked inside `ConnectAsync()` and handled by `TryExecuteMppmCommand()`.

## Features that are Not Implemented

- (Custom) authentication, usually requires an API call to a third-party authentication provider and then configuring `MatchmakingArguments.AuthValues` accordingly.
- Reconnection logic, usually requires a switch to call `MatchmakingExtensions.ReconnectToRoomAsync()` instead and a proper usage of `QuantumReconnectInformation`.
- A main menu that can start different scenes.
  - Instead create a simple main menu that loads individual scenes without connection logic, each with their own instance of `StartUI`.
  - Or use the `QuantumMenu` which provides a more complex implementation of a start menu.

Back to top

- [Use the StartUI](#use-the-startui)
- [Hierarchy and Classes](#hierarchy-and-classes)

  - [QuantumStartUI](#quantumstartui)
  - [QuantumStartUIRegionDropdown](#quantumstartuiregiondropdown)
  - [QuantumStartUIConnection](#quantumstartuiconnection)
  - [QuantumStartUIMppmCommand](#quantumstartuimppmcommand)

- [Features that are Not Implemented](#features-that-are-not-implemented)