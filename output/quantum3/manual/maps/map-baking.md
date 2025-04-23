# map-baking

_Source: https://doc.photonengine.com/quantum/current/manual/maps/map-baking_

# Map Baking

## Overview

In Quantum, a map consists of two parts: the Unity Scene and the QuantumMapData.

The Unity Scene is loaded when the map is loaded and contains the visible elements of the map such as the EntityViews. These Unity GameObjects are responsible for rendering the entities in the game.

The QuantumMapData, on the other hand, contains information about the map that is used by the deterministic simulation to drive gameplay.

## Map Baking

The QuantumSampleGame scene already comes with a QuantumMapData component which references an asset for the Map called ```
SampleMap
```

, which is located in the scene's Resources folder.

When creating a new game scene, there are a few alternatives to create and populate the new QuantumMapData component:

1. Create a new Quantum scene from the top toolbar: ```
   Quantum/Setup/Create New Quantum Scene
   ```

   ;
2. Or transform an already created scene into a Quantum scene with: ```
   Quantum/Setup/Add Quantum To Current Scene
   ```

   ;
3. Or do the map setup manually by creating the Game Object with the ```
   QuantumMapData
   ```

    component and create, wherever preferred, a new ```
   Map
   ```

    asset from the context menu in ```
   Create/Quantum/Map
   ```

   .

In the ```
QuantumEditorSettings
```

asset in the project under the ```
Editor Features
```

section, automatic map baking for scene saving, playmode changes and building the app can be enabled or disabled.

While disabling automatic map baking may be useful in certain scenarios, it is important to note that manual map baking can be time-consuming and may introduce human error into the pipeline. As a result, it is generally recommended to keep automatic map baking enabled for most projects.

![QuantumEditorSettings MapBaking](/docs/img/quantum/v3/manual/manual-quantumeditorsettings-mapbaking.png)
 QuantumEditorSettings Map Baking


For a Scene to be baked in Quantum, it is necessary to have a ```
QuantumMapData
```

component present on a GameObject in the scene. Additionally, if navmesh is present, a ```
QuantumMapNavMeshUnity
```

component is needed as well. The ```
QuantumSampleGame
```

scene provided with the Quantum SDK comes with the necessary ```
QuantumMapData
```

setup already in place.

The ```
QuantumMapData
```

MonoBehaviour component also includes buttons that can be used to manually trigger the baking process if needed. The ```
Bake All Mode
```

can be adjusted to skip certain steps in the baking process.

![MapData Component](/docs/img/quantum/v3/manual/manual-mapdata-mapbaking.png)
 QuantumMapData Component
 ## QuantumMapData

When Quantum bakes a map, it generates a ```
Map
```

asset that can be found under ```
Resources/DB/Configs
```

. However, it is possible to move these assets to another location if desired.

It's important to note that the values of the ```
Map
```

fields should generally not be changed manually, as re-baking the map will override any manual changes.

The ```
User Asset
```

 field in the Map Asset can be used to inject any asset into the map, which can then be retrieved on the simulation side. This can be done either by manually linking an asset in the inspector or by assigning one from a custom map baking callback. An example for this is shown [below](#adding_custom_data_to_a_map).

### Quantum Editor Settings BakeMapData

Map baking callbacks in Quantum allows users to inject custom steps into the QuantumMapData baking process. To implement a map baking callback, it is possible to create a class that derives from ```
MapDataBakerCallback
```

. This also requires the assembly to be marked with a custom attribute, which can be done in a single file by adding ```
\[assembly: Quantum.QuantumMapBakeAssemblyAttribute\]
```

. This signs the entire assembly, so it is not necessary to add the same attribute in multiple files. Here's an example implementation:

C#

```csharp
\[assembly: Quantum.QuantumMapBakeAssemblyAttribute\]

namespace Quantum
{
public class ExampleMapDataBaker : MapDataBakerCallback
{
public override void OnBeforeBake(QuantumMapData data)
{
}

public override void OnBake(QuantumMapData data)
{
}
}
}

```

The ```
MapDataBakerCallback
```

attribute in Quantum can be used to specify the order in which multiple custom map baking callbacks are executed. To use this attribute, add it to the custom map baking callback class and specify a ```
invokeOrder
```

value. The lower the value, the earlier the callback will be executed during the map baking process.

Here's an example:

C#

```csharp
\[MapDataBakerCallback(invokeOrder:5)\]
public class ExampleMapDataBaker : MapDataBakerCallback

```

```
OnBeforeBake
```

is called before any other MapData baking is executed. It allows to adjust the Unity scene by adding or removing components, for example.

```
OnBake
```

is called after all built-in baking steps for the map have been executed, but before the ```
Map
```

asset is saved. This allows for adjusting the ```
Map
```

while accessing all the data of the baked map.

There are additional virtual callbacks that can be overridden if needed. ```
OnBeforeBakeNavmesh
```

, ```
OnCollectNavMeshBakeData
```

, ```
OnCollectNavMeshes
```

, ```
OnBakeNavMesh
```

and an overload of ```
OnBeforeBake
```

that provides the bake flags and what triggered the baking.

## Adding Custom Data to a Map

Aside from colliders and navmeshes, game maps may include other elements that are significant for gameplay. When immutable data needs to be included, it's not mandatory to add it in entities on the map. In such cases, the User Asset field in the MapData Unity component can be used to pass any form of data into the simulation. This makes it possible to include things like objects locations, game rules and more to the map. The data assigned to the User Asset field can be accessed and used within the simulation.

### Example: Spawn Points

An example of utilizing custom data when baking a map is to include spawn points in the baked map. Designers can then place and freely move these spawn points in the Unity Scene.

In a DSL file add an asset declaration:

C#

```csharp
asset MapCustomData;

```

Then, create a new class in the ```
Quantum
```

 project to store the spawnpoint data:

C#

```csharp
namespace Quantum
{
using System;
using Photon.Deterministic;

public unsafe partial class MapCustomData
{
\[Serializable\]
public struct SpawnPointData
{
public FPVector3 Position;

public FPQuaternion Rotation;
}

public SpawnPointData DefaultSpawnPoint;
public SpawnPointData\[\] SpawnPoints;

public void SetEntityToSpawnPoint(Frame frame, EntityRef entity, Int32? index)
{
var transform = frame.Unsafe.GetPointer<Transform3D>(entity);
var spawnPoint = index.HasValue && index.Value < SpawnPoints.Length ? SpawnPoints\[index.Value\] : DefaultSpawnPoint;
transform->Position = spawnPoint.Position;
transform->Rotation = spawnPoint.Rotation;
}
}
}

```

Then, create a new class that will handle the baking of spawn points:

C#

```csharp
namespace Quantum
{
using UnityEditor;
using UnityEngine;

public class SpawnPointBaker : MapDataBakerCallback
{
public override void OnBeforeBake(QuantumMapData data)
{
}

public override void OnBake(QuantumMapData data)
{
var customData = QuantumUnityDB.GetGlobalAssetEditorInstance<Map>(data.Asset.UserAsset.Id);
var spawnPoints = GameObject.FindGameObjectsWithTag("SpawnPoint");

if (customData == null \|\| spawnPoints.Length == 0)
{
return;
}

var defaultSpawnPoint = spawnPoints\[0\];
if (customData.DefaultSpawnPoint.Equals(default(MapCustomData.SpawnPointData)))
{
customData.DefaultSpawnPoint.Position = defaultSpawnPoint.transform.position.ToFPVector3();
customData.DefaultSpawnPoint.Rotation = defaultSpawnPoint.transform.rotation.ToFPQuaternion();
}

customData.SpawnPoints = new MapCustomData.SpawnPointData\[spawnPoints.Length\];
for (var i = 0; i < spawnPoints.Length; i++)
{
customData.SpawnPoints\[i\].Position = spawnPoints\[i\].transform.position.ToFPVector3();
customData.SpawnPoints\[i\].Rotation = spawnPoints\[i\].transform.rotation.ToFPQuaternion();
}

#if UNITY\_EDITOR
EditorUtility.SetDirty(customData);
#endif
}
}
}

```

This baker is relatively simple. It gathers all GameObjects with a ```
SpawnPoint
```

tag in the scene and extracts their position and rotation, which are then saved into the custom asset. Finally, the asset is marked dirty so that the changes are stored to the disk.

To use the custom data, add GameObjects with the SpawnPoint tag to the scene. Then create a ```
MapCustomDataAsset
```

 and assign it to the ```
UserAsset
```

field of the map. To make use of the spawn points in the simulation, employ the following code:

C#

```csharp
var data = frame.FindAsset<MapCustomData>(frame.Map.UserAsset);
data.SetEntityToSpawnPoint(f, entity, spawnPointIndex);

```

## Baking at Edit time

For information on baking at edit time, see the [Map Baking](/quantum/current/manual/map-baking) page.

## Baking at Runtime

It is possible to bake maps at runtime. This is useful for procedurally generated maps or for maps that are not known at edit time.

### Before Quantum Starting

One approach is to bake a new map before the Quantum session is started. With this method, you have two options:

1. Have all clients bake the map deterministically before starting the session.
   - This method is typically used when the map generation is deterministic.
   - This method is more complex, as you need to ensure that all clients bake the map in the same way. This can be achieved by sharing the seed used for the map generation and by making sure your code is deterministic.
   - This saves bandwidth, as you only need to send the seed to the clients, and they can generate the map themselves.

How to add static assets to the ```
QuantumUnityDB
```

 at runtime:

C#

```csharp
var generatedMap = new Map();

// generate...

// add the map to the QuantumUnityDB

QuantumUnityDB.Global.AddAsset(generatedMap);

```

2. Have one client bake the map and share with the other clients.
   - This way is more simple, because you can reuse the same map baking code that you use in the editor.
   - But, due to the editor map baking not being deterministic, you need to make the asset a dynamic one and add it to the initial dynamic assets before starting Quantum. This will make the map available to all clients.

### During Quantum Session

This method is more complicated and highly game-specific.

Some approaches:

1. Have one client bake the map and share with the other clients via command.

   - This allows the reusing of the same map baking code that you use in the editor.
   - Map size is limited to command size (or needs to be manually split up and managed).
2. Have all clients bake the map deterministically during the session.

   - If your map generation is deterministic, you can have all clients bake the map during the session.
   - This is effectively the same as baking the map before the session, but the asset will still need to be dynamic, as late joiners could miss the baking event.

### Caveats

Baking maps at runtime comes with some caveats:

- The map baking process used in the Unity Editor is not deterministic by default (due to the float to FP conversion).
- Map Entity entries cannot be added or removed after the asset is added to the AssetDB.

Back to top

- [Overview](#overview)
- [Map Baking](#map-baking)
- [QuantumMapData](#quantummapdata)

  - [Quantum Editor Settings BakeMapData](#quantum-editor-settings-bakemapdata)

- [Adding Custom Data to a Map](#adding-custom-data-to-a-map)

  - [Example: Spawn Points](#example-spawn-points)

- [Baking at Edit time](#baking-at-edit-time)
- [Baking at Runtime](#baking-at-runtime)
  - [Before Quantum Starting](#before-quantum-starting)
  - [During Quantum Session](#during-quantum-session)
  - [Caveats](#caveats)