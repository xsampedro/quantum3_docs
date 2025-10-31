# extending-assets

_Source: https://doc.photonengine.com/quantum/current/manual/assets/extending-assets_

# Extending Assets for Unity

## Overview

Quantum assets can be extended with Unity types that are not relevant for the Simulation like data for the UI (colors, texts, icons...).

Let's take this `CharacterSpec` asset as an example.

C#

```csharp
public class CharacterSpec : AssetObject
{
  public FP MaxEnergy;
  public FP MaxHealth;
#if QUANTUM_UNITY
  [Header("Unity")]
  public Sprite Icon;
  public Color Color;
  public string DisplayName;
#endif
}

```

Using #if QUANTUM\_UNITY

Fields that use Unity types should only be accessed in the View (Unity) and should NEVER be accessed or used in the Simulation (Quantum). To ensure that doesn't happen, it's good practice to wrap any Unity only references (sound, icons, etc) in an `#if QUANTUM\_UNITY` block.

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

## Extending Asset Scope

Out of the box, `AssetObjects` can only use Unity types for fields used by the View. If expanding this to classes and types that are not a part of Unity is needed, a few extra steps are required.

### Creating an Assembly Definition

First, to reference types that are not a part of Unity, an `AssemblyDefinitionAsset` (or `.asmdef`) must be created. These can be created by `Assets -> Create -> Scripting -> Assembly Definition`. Any scripts within the same folder and subfolders as the `AssemblyDefinitionAsset` will now be a part of it. More about `AssemblyDefinitionAssets` can be found in _[Unity's official documentation](https://docs.unity3d.com/6000.0/Documentation/Manual/assembly-definition-files.html)_.

The following is how a user-created `AssemblyDefinitionAsset` looks.

![User-Defined Assembly Definition.](/docs/img/quantum/v3/manual/extending-assets-asmdef.png)

The two user-created scripts, `VisualsScriptableObject` and `UIScriptableObject` are now part of the `MyAssemblyDefinition` asset.

### Adding the Assembly Definition to Quantum.simulation

Within `Assets/Photon/Quantum/Simulation` is an `AssemblyDefinitionAsset` named `Quantum.Simulation`. To use the previously created `AssemblyDefinitionAsset` it must be included as a part of its `Assembly Definition References` list. The following demonstrates the inclusion of user-created `AssemblyDefinitionAsset` as well as `Unity.TextMeshPro`'s `AssemblyDefinitionAsset`; this will allow an `AssetObject` to reference user-created scripts as well as `TextMeshPro`.

![User-Defined Assembly Definition.](/docs/img/quantum/v3/manual/extending-assets-quantum-sim-asmdef.png)### Wrapping with QUANTUM\_UNITY

As mentioned previously, any fields used in the View, regardless if they are a part of Unity or not, must be wrapped in an `#if QUANTUM\_UNITY` block. The following is an example of doing so with a user-created class and `TextMeshPro` after their `AssemblyDefinitionAssets` were added to `Quantum.Simulation`.

C#

```csharp
public class CharacterSpec : AssetObject
{
  public FP MaxEnergy;
  public FP MaxHealth;
#if QUANTUM_UNITY
  [Header("Unity")]
  public Sprite Icon;
  public Color Color;
  public string DisplayName;
  public TMPro.TMP_StyleSheet StyleSheet;
  public UIScriptableObject UIVisualData;
#endif
}

```

In the above, if the user-created, `UIScriptableObject` and `TextMeshPro``AssemblyDefinitionAssets` were not setup properly, this would not compile, citing that the type or namespace could be not found. Again, the field wrapped in `#if UNITY\_QUANTUM` cannot be used during the Simulation, only the View.

Updating Photon Quantum

It's important to note that when changing to a different version of Photon Quantum, `Quantum.Simulation` may be updated and any `AssemblyDefiontionAssets` missing will needed to be added back to the list of `Assembly Definition References`.

Back to top

- [Overview](#overview)

  - [Access at Runtime](#access-at-runtime)
  - [Access at Edit-time](#access-at-edit-time)

- [Extending Asset Scope](#extending-asset-scope)
  - [Creating an Assembly Definition](#creating-an-assembly-definition)
  - [Adding the Assembly Definition to Quantum.simulation](#adding-the-assembly-definition-to-quantum.simulation)
  - [Wrapping with QUANTUM\_UNITY](#wrapping-with-quantum_unity)