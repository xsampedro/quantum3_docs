# extending-assets

_Source: https://doc.photonengine.com/quantum/current/manual/assets/extending-assets_

# Extending Assets for Unity

## Overview

Quantum assets can be extended with Unity-specific data not relevant for the simulation like data for the UI (colors, texts, icons...).

## Changes in 3.0

In Quantum 3, you no longer need to define your asset as partial and import it in a DSL file.

Instead, you should derive directly from `AssetObject` and add your desired information.

## Example

Let's take the `CharacterSpec` asset as an example.

C#

```csharp
public class CharacterSpec : AssetObject {
#if QUANTUM_UNITY
  [Header("Unity")]
  public Sprite Icon;
  public Color Color;
  public string DisplayName;
#endif
}

```

These fields should only be accessed in the View (Unity) and should NEVER be accessed or used in the simulation (Quantum). To ensure that doesn't happen, it's good practice to wrap any unity only references (sound, icons, etc) in an #if QUANTUM\_UNITY block.

### Access at Runtime

To access the extra fields at runtime, use any of the overloads of `QuantumUnityDB.GetGlobalAsset` method.

C#

```csharp
CharacterSpec characterSpec = QuantumUnityDB.GetGlobalAsset(assetRef);
Debug.Log(characterSpec.DisplayName);

```

Alternatively, `QuantumUnityDB.TryGetGlobalAsset` can be used.

C#

```csharp
if (QuantumUnityDB.TryGetGlobalAsset(assetPath, out CharacterSpec characterSpec)) {
  Debug.Log(characterSpec.DisplayName);
}

```

Both of the approaches will result in the asset being loaded into Quantum's AssetDB using the appropriate method, as discussed here: _[Resources and Addressables](/quantum/current/manual/assets/assets-unity#resources__addressables_and_asset_bundles)_.

### Access at Edit-time

To load an asset using its path while in the Unity Editor, the `UnityEditor.AssetDataBase.LoadAssetAtPath<T>()` method can be used.

C#

```csharp
CharacterSpecAsset characterSpecAsset = UnityEditor.AssetDatabase.LoadAssetAtPath<CharacterSpecAsset>(path);
Debug.Log(characterSpecAsset.DisplayName);

```

Alternatively, the asset can be loaded using any of `QuantumUnityDB.GetGlobalAssetEditorInstance` of `QuantumUnityDB.TryGetGlobalAssetEditorInstance` method.

C#

```csharp
CharacterSpec characterSpec = QuantumUnityDB.GetGlobalAssetEditorInstance<CharacterSpec>(guid);
Debug.Log(characterSpec.DisplayName);

```

Back to top

- [Overview](#overview)
- [Changes in 3.0](#changes-in-3.0)
- [Example](#example)
  - [Access at Runtime](#access-at-runtime)
  - [Access at Edit-time](#access-at-edit-time)