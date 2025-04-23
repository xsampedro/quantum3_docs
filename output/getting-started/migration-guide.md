# migration-guide

_Source: https://doc.photonengine.com/quantum/current/getting-started/migration-guide_

# Migration Guide

The complete guide to migrate Quantum 2.1 projects to Quantum 3.0.

Make sure to upgrade the projects to the latest Quantum 2.1 release and Unity 2021.3 before starting the Quantum 3 migration.

The automated process is recommended to use. When facing problems use the instructions for debugging and contact us.

Scripts GUIDs have been preserved, so prefabs and scenes are not going to suffer from missing components.

To reduce the number of possible compile errors during migration, obsolete scripts with legacy names are provided.

The following addons are **not** yet compatible with Quantum 3:

- Asset Injection

## Quantum 3.0 Migration

- Ensure that the project has a backup and/or is under source control
- Upgrade the Quantum project to the **latest 2.1 Quantum SDK**
- Upgrade the Quantum project to **Unity Editor 2021.3** or later
- Download the latest `Photon-Quantum-3.0.0-Stable-Migration-XXXX.zip` package
- Extract the content into the root folder of the Quantum 2.1 project
- **Bot SDK specifics**: find on the Bot SDK v3 package the migration instructions which are specific to the addon

### Before Starting

- Because of changes to the Quantum Asset:
  - Each non-abstract and non-generic Asset has to live in a separate file that is named as this: `<AssetName>.cs`
  - Methods on assets which are common for `ScriptableObjects` have to be renamed: `Update()`, `Start()`, etc.
- If the `GraphProfiler` addon has been added, and it has been moved from its default location `Assets/Photon/Profiling`, it has to be deleted before migration.
- Fields defined in a qtn-file called `EntityRef` have to be renamed.
- It's recommended to save intermediate steps as a git commit when doing either an automated or manual migration.

### A) Automated Process

- Run the `Quantum3Migration.ps1` Powershell script

  - To troubleshoot inspect the migration logs and enable `-PauseAfterEachStep` then wrap the changes into a git commit after each step.
- Open Unity Editor
- Continue with `Final Steps`

```
Quantum3Migration.ps1 -UnityEditorPath <path> -Quantum2MigrationPreparationPackagePath <path> -Quantum3PackagePath <path> -Quantum3MigrationPackagePath <path> [..]

Usage:
  -UnityEditorPath <path>                 The path to the UnityEditor exe.
  -QuantumUnityPath <path>                The path to the Quantum Unity project folder. By default &#34;quantum_unity&#34;.
  -QuantumCodePath <path>                 The path to the Quantum code project. By default &#34;quantum_code/quantum.code&#34;.
  -AssetDBPath <path>                     The path to the temporary exported assets. By default &#34;Quantum3MigrationAssets&#34;.
  -Quantum2MigrationPreparationPackagePath <path>
                                         The path to the Quantum 3 Migration Preparation unitypackage.
  -Quantum3PackagePath <path>            The path to the Quantum 3 SDK unitypackage.
  -Quantum3MigrationPackagePath <path>   The path to the Quantum 3 Migration unitypackage.
  -LogBasePath <path>                    The folder where the migration logs are stored. If not set logs are created inside the executing directory.
  -AssemblyDefinitionsDecision           Define the answer for the &#34;remove assembly definition&#34; prompt. &#34;yes&#34; or &#34;no&#34;.
  -PauseAfterEachStep                    Pause and wait for user input after each step. Default is false.
  -SkipQuantum2Preparation               Disable Quantum 2 preparation steps.
  -SkipQuantum3PackageImports            Disable the package import steps.
  -SkipQuantumCodeCopy                   Disable importing the quantum code.
  -SkipInitialCodeGen                    Disable initial CodeGen steps.
  -SkipCompileErrorDetection             Disable waiting for compiler error fixes.
  -SkipAssetsUpgrade                     Disable upgrading the Quantum assets.

```

### B) Manual Process

- Import `Quantum3MigrationPreparation.unitypackage`
- Run `Quantum` \> `Migration Preparation` \> `Add Migration Defines`
- Run `Quantum` \> `Migration Preparation` \> `Delete Prefab Standalone Assets`
- Run `Quantum` \> `Migration Preparation` \> `Export Assets`
- Run `Quantum` \> `Migration Preparation` \> `Delete Photon`
- Import `Photon-Quantum-3.0.0-XXXX.unitypackage` (if the Unity Editor crashes, restart Unity)
- Import `Photon-Quantum-3.0.0-Stable-Migration-XXXX.unitypackage`
- Restart Unity Editor (click `Ignore` on the `Enter Safemode` dialog)
- Run `Tools` \> `Quantum` \> `Migration` \> `Import Simulation Project`
- Run `Tools` \> `Quantum` \> `Migration` \> `Run Initial CodeGen`
- Optionally run `Tools` \> `Quantum` \> `Migration` \> `Run Delete Assembly Definitions`
- Restart Unity Editor (click `Ignore` on the `Enter Safemode` dialog)
- Fix compilation errors
  - If Quantum AssetObjects have existing `abstract void Update(Frame f)` methods Unity will complain that ScriptableObject have to have parameter-less Update() methods, add `void Update() {}` to fix the compilation error
- Run `Tools` \> `Quantum` \> `Migration` \> `Check Asset Object Scripts`
- Run `Tools` \> `Quantum` \> `Migration` \> `Transfer AssetBase Guids To AssetObjects`
- Run `Tools` \> `Quantum` \> `Migration` \> `Upgrade AssetsObjects`
- Run `Tools` \> `Quantum` \> `Migration` \> `Enable AssetObject Postprocessor`
- Run `Tools` \> `Quantum` \> `Migration` \> `Reimport All AssetObjects`
- Run the Quantum user file installation from the `QuantumHub` (Ctrl + H)

### Final Steps

- While the migration define `QUANTUM\_ENABLE\_MIGRATION` is enabled, the detection if CodeGen has to run (qtn files changed) is disabled and needs to be run manually using the Unity menu: `Tools` \> `Quantum` \> `CodeGen` \> `Run Qtn CodeGen`
- Scripts have gone through an extensive refactor:
  - Component prototype suffix has been changed from `\_Prototype` to `Prototype`.
  - Component prototype wrapper prefix changed from `EntityComponent` to `QPrototype` (e.g. `EntityComponentTransform2D` -\> `QPrototypeTransform2D`).
  - All Unity-only scripts have been put in `Quantum` namespace and have been given `Quantum` prefix (e.g. `MapData` -\> `QuantumMapData`, `EntityPrototype` -\> `QuantumEntityPrototype`).
- The separate Quantum 2 code solution `quantum.code.sln` is discontinued and can be deleted with all its projects, the simulation source files have been migrated to the Unity project
- The config assets may have different names and different locations, browse them inside the `Quantum Setup` section.
- `QuantumEditorSettings` was split into `QuantumEditorSettings` and `QuantumGameGizmosSettings`.
- A new set of default configs (previously located in Resources/DB/Configs) has been generated and does not support static Guids anymore.
- A Quantum3 AppId has to be created on the Photon dashboard.
- `QuantumEditorSettings.AssetSearchPaths` can be set to `Assets` to search the complete project. It's fast enough now.
- Let Unity reimport the Quantum assets (right-click > Reimport) located in the Resources folder for example and reimport the `QuantumUnityDB.qunitydb` asset (which will fix broken Quantum Asset Guids).
- Re-bake all maps and navmeshes.
- Quantum 3 API changes have all been wrapped into `\[Obsolete\]` attributes to make a successfully migrated project compile. Make sure to check and fix all related warnings after that especially Quantum asset and prototype scripts and prefer their new names to make `GameObject.AddComponent<QAssetEntityView>` work.

## Breaking Changes Realtime 5

- The namespaces `ExitGames.\*` are obsolete. Use `Photon.Client` and `Photon.Realtime` instead.
- The `LoadBalancingClient` was renamed to `RealtimeClient`.
- The `LoadBalancingClient.LoadBalancingPeer` was renamed to `RealtimeClient.RealtimePeer`.
- Removed `LoadBalancingPeer` class. Operations are now in the `RealtimeClient`.
- The `ExitGames.Client.Photon.Hashtable` class was renamed to `Photon.Client.PhotonHashtable`. The client can only send `PhotonHashtable` now.
- `ConnectToRegionMaster()`, `ConnectToMasterServer()`, `ConnectToNameServer()` have been removed. Use only `ConnectUsingSettings()` instead.
- `LoadBalancingClient.CloudRegion` is obsolete. Use `RealtimeClient.CurrentRegion` instead.
- `LoadBalancingClient.AppVersion` is obsolete. Set the version via AppSettings and `ConnectUsingSetting()` instead.
- `AppId` setter now throws a `NotImplementedException`. Set the AppId in the `AppSettings` and use `ConnectUsingSettings()` instead.
- The `RaiseEventOptions` class was changed to a struct called `RaiseEventArgs`.
- `EnterRoomParams` was renamed to `EnterRoomArgs`.
- `OpJoinRandomRoomParams` was renamed to `JoinRandomRoomArgs`.
- `RealtimeClient.ConnectionCallbackTargets` is now internal. Use `AddCallbackTarget()` and `RemoveCallbackTarget()` instead.
- Enum `DebugLevel` was renamed to `LogLevel`. The enum names are now Pascal case.

## FAQ

### "Couldn't set project path"

You renamed the Quantum Unity project folder. Set the `QuantumUnityPath` parameter to the correct path.

### Quantum Hub

During migration, it is possible due to user error or other bugs that Quantum Hub does not allow the user to complete the installation process. To remedy this you can try to delete the `PhotonServerSettings` scriptable object in your project in order to let the Hub regenerate the files.

### Failed Importing Package

Your Quantum 3 package file path is wrong. Set the `Quantum3PackagePath` parameter to the correct path.

OR

Your Quantum 3 Migration package file path is wrong. Set the `Quantum3MigrationPackagePath` parameter to the correct path.

### Unity Editor Path

The Unity Editor path is the path to the editor executable.

On Windows, it is usually located in `C:\\Program Files\\Unity\\Hub\\Editor\\UNITY\_VERSION\\Editor\\Unity.exe`.

The powershell script only works on Windows.

### Quantum.Code Project

If you have a different quantum code csproj name other than the default, the migration may fail. You will need to set the `QuantumCodePath` parameter to the correct path.

### Bot SDK

The Bot SDK must also be migrated in parallel if you are using it in your project.

Bot SDK data is now a singleton, rather than stored globally.

### Powershell Arguments

Powershell arguments must be passed with a single '-' and not '--'.

Example: `-UnityEditorPath "C:\\Program Files\\Unity\\Hub\\Editor\\2021.3.0f1\\Editor\\Unity.exe"`

Wrong: `--UnityEditorPath "C:\\Program Files\\Unity\\Hub\\Editor\\2021.3.0f1\\Editor\\Unity.exe"`

### .Qtn Fields

It is possible to name a field defined in a .qtn file `EntityRef`. This is supported by codegen, but not by migration. If you are having errors it is recommended to rename these fields before migrating.

### Quantum Version

If you are getting undefined errors, ensure you are updated to latest Quantum 2.1 release before attempting a migration.

### Max Component Count 512

Projects that use the 512 component count libraries need to add a new pragma to a qtn-file.

```
#pragma max_components 512

```

Back to top

- [Quantum 3.0 Migration](#quantum-3.0-migration)

  - [Before Starting](#before-starting)
  - [A) Automated Process](#a-automated-process)
  - [B) Manual Process](#b-manual-process)
  - [Final Steps](#final-steps)

- [Breaking Changes Realtime 5](#breaking-changes-realtime-5)
- [FAQ](#faq)
  - ["Couldn't set project path"](#couldnt-set-project-path)
  - [Quantum Hub](#quantum-hub)
  - [Failed Importing Package](#failed-importing-package)
  - [Unity Editor Path](#unity-editor-path)
  - [Quantum.Code Project](#quantum.code-project)
  - [Bot SDK](#bot-sdk)
  - [Powershell Arguments](#powershell-arguments)
  - [.Qtn Fields](#qtn-fields)
  - [Quantum Version](#quantum-version)
  - [Max Component Count 512](#max-component-count-512)