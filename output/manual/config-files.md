# config-files

_Source: https://doc.photonengine.com/quantum/current/manual/config-files_

# Configuration Files

## Introduction

There are seven different configuration files necessary to create and run a Quantum game.

| Parameter Name | Description |
| --- | --- |
| `PhotonServerSettings` | Stores details of the Photon cloud connection. |
| `SessionConfig` | Stores configurations of the deterministic simulation and server. |
| `SimulationConfig` | Stores configuration of the Quantum ECS layer and core systems like Physics. |
| `RuntimeConfig` | Stores data about the actual game/application. |
| `RuntimePlayer` | Stores data about the individual player. |
| `QuantumEditorSettings` | Stores configurations of the Unity project. |
| `QuantumGameGizmoSettings` | Stores configurations of Unity gizmos shown during Quantum simulations. |

`PhotonServerSettings`, `SessionConfig`, `SimulationConfig`, `QuantumEditorSettings` and `QuantumGameGizmoSettings` are stored inside `ScriptableObjects` in the Unity project while `RuntimeConfig` and `RuntimePlayer` are usually assembled during runtime.

The Quantum menu `Quantum > Find Config > ..` will help locating the global instances of the config assets in the Unity project.

![Find Quantum Configs Menu](https://doc.photonengine.com/docs/img/quantum/v3/manual/config-files/find-configs-menu.png)## PhotonServerSettings

Quantum 3.0 uses Photon Realtime 5 to connect and communicate with the Photon Cloud. This config stores all required information to establish the connection.

See the [Photon Realtime Introduction](/realtime/current/getting-started/realtime-intro "Realtime Introduction") for more information about Photon Realtime.

The most important setting is the AppId. Read the [Quantum Asteroids Tutorial - Project Setup](/quantum/current/tutorials/asteroids/2-project-setup) page, to learn how to set up a Photon AppId.

![Photon Server Settings](https://doc.photonengine.com/docs/img/quantum/v3/manual/config-files/photon-server-settings.png)

| Parameter Name | Description |
| --- | --- |
| `App Settings` | See the inline code summary or Photon Realtime 5 API documentation for details. |
| `Player TTL In Seconds` | The default Time-To-Live for players, used when creating Photon rooms. See the Realtime docs for more information. |
| `Empty Room TTL In Seconds` | The default room Time-To-Live, set when creating Photon rooms. See the Realtime docs for more information. |
| `Best Region Summary Key` | When connecting to the best region (`FixedRegion:null`) Photon does a ping to each available region for the `AppId`. The results are stored under Unity PlayerPrefs using this key to reuse during the next application start. |
| `Manage App Id` | Opens the Photon AppId dashboard URL. |
| `Best Region Cache` | The content of the best region cache read from PlayerPrefs. |
| `Reset Best Region Cache` | Delete the content of the PlayerPrefs. |
| `Open Region Dashboard` | Open the Photon region dashboard URL. |
| `Load App Settings` | The buttons configure the `AppSettings` to connect to cloud or local servers. |

## SessionConfig

Other names for the SessionConfig are DeterministicConfig or DeterministicSessionConfig.

Through the SessionConfig developers can parametrize internals of the deterministic simulation and plugin (the Quantum server component).

Each client sends its SessionConfig to the server as part of the `SessionRunner.Arguments`. The server will select the first instance it received or overwrites it by calling a webhook or running a custom plugin. The SessionConfig will be synchronized between all clients of a session before starting the simulation. See the [Online Config Sequence Diagram](#online_config_sequence_diagram) section for protocol details.

The content on this config is included in the checksum generation.

![Session Config](https://doc.photonengine.com/docs/img/quantum/v3/manual/config-files/session-config.png)

| Parameter Name | Alternative Name | Unit | Description |
| --- | --- | --- | --- |
| `Simulation Rate` | `Update FPS` | Hz | How many ticks per second Quantum should execute. |
| `Lockstep` | `Lockstep Simulation` | Boolean | Runs the quantum simulation in lockstep mode, where no rollbacks are performed. It is recommended to set input `InputDelayMin` to at least `10`. |
| `Rollback Window` |  | Tick Count | How many frames are kept in the local ring buffer on each client. Controls how much Quantum can predict into the future. Not used in lockstep mode. |
| `Checksum Interval` |  | Tick Count | How often we should send checksums of the frame state to the server for verification (useful during development, set to zero for release). Defined in frames. |
| `Checksum Cross Platform Determinism` |  | Boolean | This allows Quantum frame checksums to be deterministic across different platforms, however it comes with quite a cost and should only be used during debugging. |
| `Input Delta Compression` |  | Boolean | If the server should delta-compress inputs against previous tick-input-set. Reduces overall bandwidth. |
| `Offset Min` | `Input Delay Min` | ms | The minimum input delay a player can have. |
| `Offset Max` | `Input Delay Max` | ms | The maximum input delay a player can have. |
| `Offset Ping Start` | `Input Delay Ping Start` | ms | At which ping value Quantum starts applying input delay. |
| `Send Redundancy` | `Input Redundancy` | Tick Count | How much staggering the Quantum client should apply to redundant input resend. `1` = Wait one frame, `2` = Wait two frames, etc. |
| `Input Repeat Max Distance` |  | Tick Count | How many frames Quantum will scan for repeatable inputs. `5` = Scan five frames forward and backwards, `10` = Scan ten frames, etc. |
| `Hard Tolerance` |  | Tick Count | How many frames the server will wait until it expires a frame and replaces all non-received inputs with repeated inputs or null's and sends it out to all players. |
| `Offset Correction Limit` | `Min Time Correction Frames` | Tick Count | How much the local client time must differ with the server time when a time correction package is received for the client to adjust it's local clock. |
| `Correction Send Rate` | `Time Correction Rate` | Hz | How many times per second the server will send out time correction packages to make sure every clients time is synchronized. |
| `Correction Frames Limit` | `Min Offset Correction Diff` | Tick Count | How many frames the current local input delay must diff to the current requested offset for Quantum to update the local input offset. |
| `Room Wait Time` | `Session Start Timeout` | s | An artificial wait time to control how long the server waits for other players after the online simulation has been requested to start. |
| `Time Scale Minimum` | `Time Scale Min` | % | The smallest timescale that can be applied by the server. |
| `Time Scale Ping Start` | `Time Scale Ping Min` | ms | The ping value that the server will start lowering the time scale towards `Time Scale Minimum`. |
| `Time Scale Ping End` | `Time Scale Ping Max` | ms | The ping value that the server will reach the `Time Scale Minimum` value at, i.e. be at its slowest setting. |
| `Player Count` |  | int | Player count the simulation is initialized for. Can be left 0 because this parameter is overwritten by `SessionRunner.Arguments.PlayerCount` (when > 0) when starting the SessionRunner. |
| `Input Fixed Size` |  | int | The size of the input struct. This will be set internally after starting the session. The method to compute this is `QuantumGame.GetInputSerializedFixedSize()` |

## SimulationConfig

The SimulationConfig holds parameters used by the Quantum ECS layer and inside core systems like physics and navigation. See the related system sections in the manual for more details of each value.

The SimulationConfig is part of the Quantum DB and multiple instances of this config are supported. Reference the SimulationConfig with the RuntimeConfig to select which one to use for the simulation.

During the Quantum Unity project initialization by the Quantum Hub the following files are created.

- `QuantumUser/Resources/QuantumDefaultConfigs.asset` which includes a SimulationConfig asset and it's referenced default config assets such as PhysicsMaterial, CharacterController2DConfig, NavMeshAgentConfig, etc. as sub assets.

![Default Configs](https://doc.photonengine.com/docs/img/quantum/v3/manual/config-files/default-configs.png)

- The file `QuantumUser/Game/SimulationConfig.User.cs` contains a partial class definition that can be used to extend the content of the SimulationConfig.

C#

```csharp
namespace Quantum {
  public partial class SimulationConfig : AssetObject {
    public int Foo;
  }
}

```

![Simulation Config](https://doc.photonengine.com/docs/img/quantum/v3/manual/config-files/simulation-config.png)

| Parameter Name | Description |
| --- | --- |
| `Entities` | See API docs. |
| `Physics` | See physics docs. |
| `Navigation` | See navigation docs. |
| `Auto Load Scene From Map` | This option will trigger a Unity scene load during the Quantum start sequence. Which might be convenient to start with but once the starting sequence is customized it should be disabled and replaced by custom scene loading. "Previous Scene" refers to a scene name in Quantum Map.<br> <br> The demo menu for example has a step that can load the scene before starting the Quantum simulation when AutoLoadSceneFromMap is disabled. |
| `Thread Count` | Override the number of threads used internally. Default is 2. |
| `Checksum Snapshot History Length` | How long to store checksums of verified frames. They are used to generate a frame dump in case of a checksum error happening. Not used in Replay and Local mode. Default is 3. |
| `Checksum Error Dump Options` | Additional options for checksum dumps, if the default settings don't provide a clear picture. |
| `Heap Tracking Mode` | If and to which extent allocations in the Frame Heap should be tracked when in Debug mode. <br> Recommended modes for development is `DetectLeaks`. While actively debugging a memory leak, `TraceAllocations` mode can be enabled (warning: tracing is very slow). |
| `Heap Page Shift` | Define the max heap size for one page of memory the frame class uses for custom allocations like QList, for example. The default is 15. <br> Example: 2^15 = 32.768 bytes<br>`TotalHeapSizeInBytes = (1 << HeapPageShift) \* HeapPageCount` |
| `Heap Page Count` | Define the max heap page count for memory the frame class uses for custom allocations like QList for example. Default is 256.<br>`TotalHeapSizeInBytes = (1 << HeapPageShift) \* HeapPageCount` |
| `Heap Extra Count` | Sets extra heaps to allocate for a session in case you need to create more (auxiliary) frames than actually required for the simulation itself. Default is 0. |

## RuntimeConfig

In contrast to the SimulationConfig the RuntimeConfig holds information that can be **different from game to game**. By default it defines for example what map to load and the random seed. It does not have an asset to store the configs but it is assembled during runtime most likely based on the selection player do (e.g. game mode).

To use a certain RuntimeConfig, assign it to the `SessionRunner.Arguments` when starting the session.

C#

```csharp
var map              = new AssetRef<Map>(QuantumUnityDB.GetGlobalAssetGuid(&#34;Photon/Quantum/Samples/SampleScenes/Resources/SampleMap&#34;));
var simulationConfig = new AssetRef<SimulationConfig>(QuantumUnityDB.GetGlobalAssetGuid(&#34;QuantumUser/Resources/QuantumDefaultConfigs|DefaultConfigSimulation&#34;));
var systemsConfig    = new AssetRef<SystemsConfig>(QuantumUnityDB.GetGlobalAssetGuid(&#34;Photon/QuantumUser/Resources/DefaultSystemsConfig&#34;));

var sessionRunnerArguments = new SessionRunner.Arguments {
    RuntimeConfig = new RuntimeConfig() {
        Map              = map,
        Seed             = DateTime.Now.Millisecond
        SimulationConfig = simulationConfig,
        SystemsConfig    = systemsConfig },
    // ..
};

```

Similar to the SessionConfig the RuntimeConfig is sent to the server by each client, it can be validated by webhooks or a custom server then one version is distributed to every client during the Quantum start sequence.

Unlike SessionConfig, which is serialized in a binary form when send to the server, RuntimeConfig uses by default zip-compressed **Json serialization** to upload and download the config from the Quantum server.

During the Quantum Unity project installation a `QuantumUser/Game/RuntimeConfig.User.cs` script is created that allows the RuntimeConfig to be extended.

C#

```csharp
namespace Quantum {
  public partial class RuntimeConfig {
    // Add your own fields (don&#39;t use properties).
    public int Foo;

    // Implement DumpUserData() to add information to a debug string that is returned when using Dump().
    partial void DumpUserData(ref String dump) {
    }
  }
}

```

Create a copy of a RuntimeConfig, for example, to test the serialization.

C#

```csharp
var copy = RuntimeConfig.Copy(runtimeConfig, new QuantumUnityJsonSerializer());

```

If desired the config can also be stored on a game object like the `QuantumRunnerLocalDebug` script does.

C#

```csharp
public class QuantumRunnerLocalDebug : QuantumMonoBehaviour {
    public RuntimeConfig RuntimeConfig;
    // ..
}

```

| Parameter Name | Description |
| --- | --- |
| `Seed` | The seed to initialize the randomization session under `Frame.RNG`. |
| `Map` | Asset reference of the Quantum Map used with the upcoming game session. |
| `SimulationConfig` | Asset reference to the SimulationConfig used with the upcoming game session. |
| `SystemsConfig` | Asset reference to the Quantum systems configuration.<br> If no config is assigned then a default selection of built-in systems is used `DeterministicSystemSetup.CreateSystems(RuntimeConfig, SimulationConfig, SystemsConfig)`.<br> The systems to be used can always be post processed by code using the partial method `DeterministicSystemSetup.AddSystemsUser(System.Collections.Generic.ICollection{SystemBase}, RuntimeConfig, Quantum.SimulationConfig, Quantum.SystemsConfig).` |

## RuntimePlayer

Similar to the RuntimeConfig the RuntimePlayer describes run-time properties for one player.

The data for a player behaves differently to the other configs, because it is sent by each player individually after the actual game has been started (AddPlayer). See the [Player](/quantum/current/manual/player/player) and [Online Session](/quantum/current/manual/game-session/starting-session) documentation for more information.

During the Quantum Unity project installation a `QuantumUser/Game/RuntimePlayer.User.cs` script is created that allows the RuntimePlayer to be extended.

C#

```csharp
namespace Quantum {
  public partial class RuntimePlayer {
    // Add your own fields (don&#39;t use properties).
    public int Foo;

    // Implement DumpUserData() to add information to a debug string that is returned when using Dump().
    partial void DumpUserData(ref String dump) {
    }
  }
}

```

Similar to the RuntimeConfig the RuntimePlayer data is serialized with Json when sending and receiving from the Quantum server.

The default RuntimePlayer already contains two fields that can be used as a starting point for player visualization and nickname. The demo menu uses them for example.

C#

```csharp
public partial class RuntimePlayer {
    // This is a proposal how to let players select an avatar prototype using RuntimePlayer. Can be ignored.
    public AssetRefEntityPrototype PlayerAvatar;
    // This is a proposal how to assign a nickname to players using RuntimePlayer. Can be ignored.
    public string PlayerNickname;
}

```

## QuantumEditorSettings

The Quantum editor settings hold information vital for the Unity integration to function.

During Quantum Unity project installation by the Quantum Hub a default editor settings asset is created at `QuantumUser/Editor/QuantumEditorSettings.asset`.

![Quantum Editor Settings](https://doc.photonengine.com/docs/img/quantum/v3/manual/config-files/editor-settings.png)

| Parameter Name | Description |
| --- | --- |
| `Asset Search Paths` | Locations that the QuantumUnityDB uses to find Quantum assets. Changing this requires reimporting all Unity (Quantum) assets manually. |
| `Default New Assets Location` | Default folder where new Quantum assets are created. |
| `Use Quantum Unit DB Asset Postprocessor` | The post processor enables duplicating Quantum assets and prefabs and make sure a new GUID and correct path are set. This can make especially batched processes slow and can be toggled off here. |
| `Use Quantum Toolbar Utilities` | If enabled a scene loading dropdown is displayed next to the play button. |
| `Quantum Toolbar Zone` | Where to display the toolbar. Requires a domain reload after change. |
| `Use Photon App Versions Postprocessor` | If enabled a local PhotonPrivateAppVersion scriptable object is created to support the demo menu scene. |
| `Entity Component Inspector Mode` | If enabled entity components are displayed inside of EntityPrototype inspector |
| `FP Display Precision` | How many decimal places to round to when displaying FPs. Default is 5. |
| `Auto Build On Scene Save` | Automatically trigger bake when saving a scene. |
| `Auto Build On Playmode Changed` | If set, MapData will be automatically baked when entering play mode, when saving a scene and when building a player. |
| `Auto Build On Build` | If set MapData will be automatically baked when building, when saving a scene and when building a player. |
| `Auto Run Qtn CodeGen` | If enabled any changes in .qtn files in quantum.code will run the codegen immediately. |
| `Asset Guid Overrides` | A list of Quantum assets that enabled GUID Override. This list is tracked automatically. |
| `Quantum Debug (All Platforms)` | Toogles QUANTUM\_DEBUG scripting define for all platforms to make use Quantum debug dlls. |
| `Quantum Debug (Current Platform)` | Toogles QUANTUM\_DEBUG scripting define for the current platform to make use Quantum debug dlls. |

## QuantumGameGizmoSettings

The gizmo settings contain information about Quantum debug and gizmo rendering, colors and toggles.

During Quantum Unity project installation by the Quantum Hub a default editor settings asset is created at `QuantumUser/Editor/QuantumGameGizmosSettings.asset`.

Usually the global config retrieved by `QuantumGameGizmosSettingsScriptableObject.Global.Settings` is used. To use a different one, assign it to `QuantumRunner.GizmoSettings` during runtime.

![Quantum Editor Settings](https://doc.photonengine.com/docs/img/quantum/v3/manual/config-files/gizmo-settings.png)## Online Config Sequence Diagram

The flow of client controlled config files through the Quantum connection protocol.

![Online Config Sequence Diagram](https://doc.photonengine.com/docs/img/quantum/v3/manual/config-files/connection-flow.jpg)
Online Config Sequence Diagram
Back to top

- [Introduction](#introduction)
- [PhotonServerSettings](#photonserversettings)
- [SessionConfig](#sessionconfig)
- [SimulationConfig](#simulationconfig)
- [RuntimeConfig](#runtimeconfig)
- [RuntimePlayer](#runtimeplayer)
- [QuantumEditorSettings](#quantumeditorsettings)
- [QuantumGameGizmoSettings](#quantumgamegizmosettings)
- [Online Config Sequence Diagram](#online-config-sequence-diagram)