# quantum-project

_Source: https://doc.photonengine.com/quantum/current/manual/quantum-project_

# Quantum Project

## Unitypackage Content

The Quantum 3 SDK is distributed as a `unitypackage` file. It's roughly separated into two sections: `Assets/Photon` and `Assets/QuantumUser`. The former represents the Quantum SDK and the latter mutable user code.

After installing Quantum, the user is presented with following folder structure:

```
Assets
├───Photon
│   ├───PhotonLibs
│   ├───PhotonRealtime
│   ├───Quantum
│   ├───QuantumAsteroids
│   └───QuantumMenu
└───QuantumUser
    ├───Editor
    │   ├───CodeGen
    |   └───Generated
    ├───Resources
    ├───Scenes
    ├───Simulation
    │   └───Generated
    └───View
        └───Generated

```

Upgrading will replace all files in Photon subfolders but not in QuantumUser.

All deterministic simulation code MUST be in `QuantumUser/Simulation` or included in the Quantum.Simulation assembly reference.

Code extending Quantum view scripts CAN be in `Quantum.Unity` or `Quantum.Unity.Editor` (e.g. partial methods) using their respective assembly references.

It is possible to opt out of the Quantum assembly definitions, for example for users migrating from such setup in Quantum 2.1 projects.

**Assets/Photon**

- `Assets/Photon`: This folder contains Photon Quantum and other packages. Changes to files in here will be overwritten during upgrades
- `Assets/Photon/PhotonLibs (and PhotonRealtime)`: Contains Photon dependencies that Quantum uses to connect and communicate with the Quantum cloud.
- `Assets/Photon/Quantum/Assemblies`: Contains Quantum libraries for Unity and their debug versions.
- `Assets/Photon/Quantum/Editor`: Quantum Unity editor scripts. Compiles to `Quantum.Unity.Editor.dll`.
- `Assets/Photon/Quantum/Editor/Assemblies`: Contains the Quantum CodeGen libraries.
- `Assets/Photon/Quantum/Editor/CodeGen`: Includes the Quantum CodeGen scripts. Compiles to `Quantum.CodeGen.Qtn.dll`.
- `Assets/Photon/Quantum/Runtime`: Quantum Unity runtime scripts. Compiles to `Quantum.Unity.dll`.
- `Assets/Photon/Quantum/Simulation`: Quantum simulation core scripts. Compiles to `Quantum.Simulation.dll`.
- `Assets/Photon/Quantum/Samples`: Contains the demo menu scene and currently the GraphProfiler.

**Assets/QuantumUser**

- `Assets/QuantumUser`: Files inside this folder will never be overwritten by an upgrade and are under the developers control. It uses assembly definition references to add code to the Quantum and simulation libraries.
- `Assets/QuantumUser/Editor/CodeGen`: Files that extend the Qtn CodeGen.
- `Assets/QuantumUser/Editor/Generated`: Generated Unity Editor scripts.
- `Assets/QuantumUser/Simulation`: The actual simulation code.
- `Assets/QuantumUser/Simulation/Generated`: Generated Quantum C# code.
- `Assets/QuantumUser/Resources`: Runtime config files.
- `Assets/QuantumUser/Scenes`: The default folder for new Quantum scenes.
- `Assets/QuantumUser/View`: The default location for view scripts.
- `Assets/QuantumUser/View/Generated`: Generated Quantum prototype scripts, scripts inside this folder extend the Quantum.Unity.dll.

**Quantum is split into four assemblies:**

- `Quantum.Simulation`: contains simulation code. Any user simulation code should be added to this assembly with `AssemblyDefinitionReferences`. Unity/Odin property attributes can be used at will, but any use of non-deterministic Unity API is heavy discouraged. Code from this assembly can be easily worked on as a standalone `.csproj`, similar to `quantum.code.csproj` in Quantum 2.
- `Quantum.Unity`: contains code specific to Quantum's integration with Unity. Additionally, CodeGen emits `MonoBehaviours` that wraps component prototypes.
- `Quantum.Unity.Editor`: contains editor code for `Quantum.Simulation` and `Quantum.Unity`
- `Quantum.Unity.Editor.CodeGen`: contains CodeGen integration code. It is fully independent of other Quantum assemblies, so can always be run, even if there are compile errors - this may require exiting Safe Mode.

## Quantum Dependencies

![Quantum Dependencies](/docs/img/quantum/v3/manual/quantum-project/quantum-dependencies.jpg)## Quantum Hub

Open by pressing `Ctrl+H` or via the Quantum menu.

The Hub window will also pop up when a vital Quantum configuration files is missing (e.g. PhotonServerSettings) and recommend pressing the installation button.

The installation process takes care of installing files locally that cannot originate from the unitypackage because they would be overwritten by the next Quantum version upgrade.

![Quantum Hub](/docs/img/quantum/v3/manual/quantum-project/quantum-hub.png)## Installing The Quantum Menu

The Quantum menu is an addon that can be installed with the SDK using the Hub. The unitypackage can be found under `Asset/Photon/QuantumMenu/Quantum-Menu.unitypackage`.

The menu is a functional and graphical in-game menu to start random online matches or create parties.

![Installing The Quantum Menu](/docs/img/quantum/v3/manual/quantum-project/quantum-hub-menu.png)

Find more information about customization possibilities here: [Sample Menu Customization](/quantum/current/manual/sample-menu/sample-menu-customization).

## Release And Debug Builds

The Quantum libraries (Quantum.Deterministic.dll, Quantum.Engine.dll and Quantum.Corium.dll) come both in release and debug configurations. To make Unity recognize the correct library the `QUANTUM\_DEBUG` global scripting define is used.

To toggle between debug and release use the menu:

![Debug Toggle](/docs/img/quantum/v3/manual/quantum-project/menu-debug-toggle.png)

Caveat: because the define needs to be set per platform, there is a risk that the debug version is not disabled for every platform when making release builds for example.

The debug build has **significant performance penalties** compared to a release build. For performance tests always use a Quantum release build (and Unity IL2CPP). Read more about this in the [Profiling](/quantum/current/manual/profiling "Profiling") section.

The development build contains assertions, exceptions, checks and debug outputs that help during development and which are disabled in **release configuration**. For example:

- `Log.Debug()` and `Log.Trace()`, when called from the quantum code project, will not be outputting log anymore.
- As well as all `Draw.Shape()` methods.
- `NavMeshAgentConfig.ShowDebugAvoidance` and `ShowDebugSteering` will not draw gizmos anymore.
- Assertions and exceptions inside low level systems like physics are disabled.

## Logging

Quantum provides the static `Quantum.Log` class to log from the simulation code and also to produce all its owns log outputs.

C#

```csharp
namespace Quantum {
  public unsafe class MyQuantumSystem : SystemMainThread
    public override void Update(Frame frame) {
      Log.Debug($"Updating MyQuantumSystem tick {frame.Number}");
    }
  }
}

```

The Unity SDK has wrapper around that called `QuantumUnityLogger`. It will statically initialize itself during `RuntimeInitializeOnLoadMethod` and/or `InitializeOnLoadMethod` (see `QuantumUnityLogger.Initialize()`)

The partial method `InitializePartial()` can be used to customize the initialization.

The `QuantumUnityLogger` class offers various customizations fields to for example define the color schemes.

### Log Level

The global log level is controlled by `Quantum.Log.LogLevel`. To set the **initial log level** use the following scripting defines:

- `QUANTUM\_LOGLEVEL\_TRACE`
- `QUANTUM\_LOGLEVEL\_DEBUG`
- `QUANTUM\_LOGLEVEL\_INFO`
- `QUANTUM\_LOGLEVEL\_WARN`
- `QUANTUM\_LOGLEVEL\_ERROR`

They are toggled with the `QuantumEditorSettings` inspector.

![Editor Settings - LogLevel](/docs/img/quantum/v3/manual/quantum-project/editor-settings-loglevel.png)

As mentioned in the last section about debug builds `Log.Trace()` and `Log.Debug()` messages will only be logged with `TRACE` and `DEBUG` defines respectively. Notice that running the UnityEditor always defines `DEBUG`.

### Photon Realtime Log

The Photon Realtime library used by Quantum has its own log levels that can be controlled with the `Photon.Realtime.AppSettings` found on the `PhotonServerSettings` ScriptableObject.

- `AppSettings.NetworkLogging`: Log level for the PhotonPeer and connection. Useful to debug connection related issues.
- `AppSettings.ClientLogging`: Log level for the RealtimeClient and callbacks. Useful to get info about the client state, servers it uses and operations called.

In non-development builds Realtime logs less than the `ERROR` severity will not log unless defining `LOG\_WARNING`, `LOG\_INFO` and `LOG\_DEBUG`.

## Exporting Simulation

The Quantum simulation code can be exported to a standalone C# project that is free of Unity dependencies. The project is generated into the Unity project folder, it will link to all simulation source files found and uses non-Unity Quantum dependencies extracted from `Assets/Photon/Quantum/Editor/Dotnet/Quantum.Dotnet.Debug.zip` for example.

1. Select the `QuantumDotnetProjectSettings` asset in the Unity project to define what simulation sources are searched for.

   - Add additional `IncludePaths` manually or mark the folders with the `QuantumDotNetInclude` Unity asset label.
2. Select the `QuantumDotnetBuildSettings` asset in the Unity project control how the DotNet project is generated and build.


![QuantumDotnetBuildSettings](/docs/img/quantum/v3/manual/quantum-project/quantumdotnetbuildsettings.png)

3. Enable `Show Folder After Generation` and press `Generate Dotnet Project`
   - The project will be generated at `Project Base Path`, relative to the Unity project folder.

![Exported Project](/docs/img/quantum/v3/manual/quantum-project/exported-dotnet.png)### Exported Solution Structure

The generated solution comes with two projects:

- `Lib` \- The unzipped non-Unity Quantum dependencies in Debug and Release configuration
- `Quantum.Runner.Dotnet` \- Contains a lightweight console runner that can run Quantum replays
- `Quantum.Simulation.Dotnet` \- Contains the non-Unity Quantum simulation project and code

#### Using the Console Runner

1. Build the solution in your IDE.

2. Open a terminal and navigate to the folder with the built exe (`Quantum.Runner.exe`).

3. Run the exe with the following arguments:


bash

```bash
Quantum.Runner.exe --replay-path path/to/replay --lut-path path/to/lut --db-path path/to/db --checksum-path path/to/checksum

```

The `--replay-path` and `--lut-path` arguments are required.

The `--db-path` argument is optional if the replay contains the db.

The `--checksum-path` argument is completely optional.

To acquire a standalone assets file, you can press the `Asset Database` menu button under the `Tools/Quantum/Export/AssetDatabase` menu.

![export assets](/docs/img/quantum/v3/manual/quantum-project/export-assets.png)

To acquire the LUT files, you can access them in the `Assets/Photon/Quantum/Resources/LUT` folder in the Unity project.

## Syntax Highlighting in .qtn files

To enable syntax highlighting in the DSL (files with `.qtn` extension), follow the the IDE specific guides below.

### Visual Studio

In Visual Studio, it is possible to add syntax highlighting for QTN files by associating it with another type (e.g. C# or Microsoft Visual C++). To do this go to `Tools -> Options -> Text Editor -> File Extension`.

![File Types](/docs/img/quantum/v3/manual/visual-studio-dsl-highlight.png)
DSL Syntax Highlighting in .qtn files (Visual Studio).
### Visual Studio Code

In Visual Studio, it is possible to add syntax highlighting for QTN files by associating it with another type (e.g. C#).

1. Open Settings (Ctrl+, or Cmd+, on macOS).
2. Search for "Files: Associations".
3. Add a new file association of "\*.qtn" to "csharp".

![Visual Studio Code File Association settings](/docs/img/quantum/v3/manual/visual-studio-code-dsl-highlight.png)
DSL Syntax Highlighting in .qtn files (Visual Studio).
### JetBrains Rider

In JetBrains Rider, it is possible to add syntax highlighting to the DSL files (the ones with `.qtn` extension) by defining a new file type.

- **Step 1:** Navigate to `File -> Settings -> Editor -> File Types`.

![File Types](/docs/img/quantum/v3/manual/dsl-syntax-highlighting-overview.png)
The \`File Types\` settings in JetBrains Rider.


- **Step 2:** In the `Recognized File Types` category, press the `+` sign at the right of the to add a new file type.

![New File Type](/docs/img/quantum/v3/manual/dsl-syntax-highlighting-definition.png)
The \`New File Type\` window in JetBrains Rider.


- **Step 3:** Check the settings for line comments, block comments, etc...
- **Step 4:** Paste the list below into the keywords level 1.

C#

```csharp
#define
#pragma
abstract
any
array
asset
asset_ref
bitset
button
byte
component
dictionary
entity_ref
enum
event
fields
filter
flags
global
has
import
input
int
list
local
long
not
player_ref
remote
sbyte
set
short
signal
struct
synced
uint
ulong
union
use
ushort
using

```

- **Step 5:** Paste the list below into the keywords level 2 then press `Ok`.

C#

```csharp
(
)
*
:
;
<
=
>
?
[
]
{
}

```

- **Step 6:** In the `File Name Patterns` category, press the `+` sign at the right side.
- **Step 7:** Enter `\*.qtn` as the wildcard for the type.

![DSL Syntax Highlighting](/docs/img/quantum/v3/manual/rider-dsl-highlight.png)
DSL Syntax Highlighting in .qtn files (JetBrains Rider).
Back to top

- [Unitypackage Content](#unitypackage-content)
- [Quantum Dependencies](#quantum-dependencies)
- [Quantum Hub](#quantum-hub)
- [Installing The Quantum Menu](#installing-the-quantum-menu)
- [Release And Debug Builds](#release-and-debug-builds)
- [Logging](#logging)

  - [Log Level](#log-level)
  - [Photon Realtime Log](#photon-realtime-log)

- [Exporting Simulation](#exporting-simulation)

  - [Exported Solution Structure](#exported-solution-structure)

- [Syntax Highlighting in .qtn files](#syntax-highlighting-in.qtn-files)
  - [Visual Studio](#visual-studio)
  - [Visual Studio Code](#visual-studio-code)
  - [JetBrains Rider](#jetbrains-rider)