# migration-guide

_Source: https://doc.photonengine.com/quantum/current/getting-started/preview-3-1/migration-guide_

# Migration Guide

## Migrating Quantum From 3.0 to 3.1

![Icon](/v2/img/docs/icons/alert-warning.svg)

Projects in Quantum SDK 2.1 should be upgraded to the latest 3.0 SDK before being upgraded to 3.1.

### Migration Steps

- Download the latest [Quantum SDK 3.1 unitypackage](https://downloads.photonengine.com/download/latest/photon-quantum-sdk-3-1)
- Check for modifications to any Quantum Unity scripts and backup the changes.
- Import the `.unitypackage` into the Unity project.
- Delete remnants of the Quantum.Deterministic library: `Quantum.Deterministic.dll`, `Quantum.Deterministic.pdb` and `Quantum.Deterministic.xml`.![](/docs/img/quantum/v3/getting-started/migration/3.1-deleting-quantum-deterministic.png)
- Restart the Unity Editor.
- Fix any remaining compilation errors.
- (Optional) If the project uses Bot SDK, close its AI documents and editor window, download Bot SDK 3.1.0 Preview package and import it into the project.
- Reimport the asset `Assets/QuantumUser/Resources/QuantumDefaultConfigs.asset` to make sure the `InputActionAsset` is set.
- The multi client scripts and prefabs have been moved into a additional `.unitypackage`, delete the these files from the SDK.


  - `Assets/Photon/Quantum/Runtime/RuntimeAssets/QuantumMultiClientRunner.prefab`
  - `Assets/Photon/Quantum/Runtime/QuantumMultiClientPlayer.cs`
  - `Assets/Photon/Quantum/Runtime/QuantumMultiClientPlayerView.cs`
  - `Assets/Photon/Quantum/Runtime/QuantumMultiClientRunner.cs`

![](/docs/img/quantum/v3/getting-started/migration/3.1-deleting-multirunner.png)
- Import optional used Quantum packages from `Assets/Photon/Quantum/PackageResources`: Quantum-DemoInput, Quantum-Menu or Quantum-MultiClient by importing the `.unitypackage` manually or by reinstalling it using the Quantum Hub Sample section.![](/docs/img/quantum/v3/getting-started/migration/3-1-import-unitypackages.png)

Back to top

- [Migrating Quantum From 3.0 to 3.1](#migrating-quantum-from-3.0-to-3.1)
  - [Migration Steps](#migration-steps)