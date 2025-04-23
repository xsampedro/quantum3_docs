# download

_Source: https://doc.photonengine.com/quantum/current/addons/kcc/download_

# Download

Here you can download sample project which includes step-by-step explanations and several gameplay oriented examples. Download addon if you are updating to new version or want to start from scratch.

## Download - Sample

| Version | Release Date | Download |
| --- | --- | --- |
| 3.0.5 | Mar 06, 2025 | [Quantum KCC Sample 3.0.5 Build 596](https://dashboard.photonengine.com/download/quantum/quantum-kcc-sample-3.0.5.zip) |

## Download - Addon

| Version | Release Date | Download |
| --- | --- | --- |
| 3.0.5 | Mar 06, 2025 | [Quantum KCC 3.0.5 Build 595](https://dashboard.photonengine.com/download/quantum/quantum-kcc-3.0.5.unitypackage) | [Release Notes](#3.0.5) |

## Requirements

- Unity 2021.3
- Quantum V3 AppId: To run the sample, first create a Quantum V3 AppId in the [PhotonEngine Dashboard](https://dashboard.photonengine.com) and paste it into the `App Id Quantum` field in Photon Server Settings (reachable from the `Tools/Quantum/Find Config/Photon Server Settings` menu in Unity editor).

## Release Notes

### Photon Quantum KCC Addon

Last tested with `Quantum SDK 3.0.2 Stable 1660`

#### 3.0.5

- Fixed position in CapsuleCast.
- Improved mesh collider depenetration algorithm.
- Added support for processor priority and sorting - virtual FP KCCProcessor.GetPriority(). Processors with higher priority are executed first.
- Added full support for processor suppressing - KCCContext.StageInfo.SuppressProcessors<T>(). This skips execution of all pending processors of type T in current stage - e.g. IBeforeMove.
- Environment Processor - added IPrepareData interface to allow modifying KCCData properties (like Gravity) before calculating velocities.

#### 3.0.4

- Added Entity reference to KCC.
- KCC.Teleport(FPVector3) marked as obsolete and replaced by KCC.Teleport(Frame, FPVector3). This propagates the teleport also to Transform3D component.
- Drawing KCC collider gizmo based on linked KCC settings when the game object is selected.
- Quantum.Unity.asmref file moved from KCC/View/Generated to KCC/View.

#### 3.0.3 (Breaking Changes)

- Code generated prototype scripts are now included in the addon package instead of being generated in the user project.


After upgrading the addon old prototype scripts are deleted and `QuantumEntityPrototypes` that used them will break.


To migrate choose one of the options before upgrading:

  - A) Search And Replace Guids
    - Close Unity Editor
    - Open the `QPrototypeKCC.cs.meta` file, copy the Unity script guid and search and replace all occurrences of it with `ef7706d4b9fc4dc468d3a1cf0c2dde40` in files inside the Assets folder.
    - Open the `QPrototypeKCCProcessorLink.cs.meta` file, copy the Unity script guid and search and replace all occurrences of it with `6a24c6b0be5af364298f6a16f1d81325` in files inside the Assets folder.
  - B) Rename And Inherit Class


    - Open the `QPrototypeKCC.cs` file in Rider or VS, rename the class to `LegacyQPrototypeKCC`, inherit it from `QPrototypeKCC` and delete its content.

C#
```csharp
public unsafe partial class LegacyQPrototypeKCC : QPrototypeKCC {
}

```

    - Open the `QPrototypeKCCProcessorLink.cs` file in Rider or VS, rename the class to `LegacyQPrototypeKCCProcessorLink`, inherit it from `QPrototypeKCCProcessorLink` and delete its content.

C#
```csharp
public unsafe partial class LegacyQPrototypeKCCProcessorLink : QPrototypeKCCProcessorLink {
}

```

    - Rename the `QPrototypeKCC` script inside Unity Editor to `LegacyQPrototypeKCC` and move it into a subfolder called `Legacy.`
    - Rename the `QPrototypeKCCProcessorLink` script inside Unity Editor to `LegacyQPrototypeKCCProcessorLink` and move it into a subfolder called `Legacy`.

#### 3.0.2

- Improved penetration solver.
- Range of CCD Radius Multiplier increased to 10-90%.
- Ground snap processor now does all checks on separate `KCCData` instance to avoid side effects.
- `KCC.ResolvePenetration()` now takes data from `KCCData` passed as argument and stores results there.
- Added `\[Preserve\]` attribute to `KCCSystem` prevent stripping.

#### 3.0.1

- Improved multi-collider penetration correction.
- Fixes for `StepUpProcessor`. Now it requires horizontal movement push-back to activate.
- Fixed initialization of `KCCData.LookYaw` after spawning an entity based on its transform component.

#### 3.0.0

- Initial release.

Back to top

- [Download - Sample](#download-sample)
- [Download - Addon](#download-addon)
- [Requirements](#requirements)
- [Release Notes](#release-notes)
  - [Photon Quantum KCC Addon](#photon-quantum-kcc-addon)