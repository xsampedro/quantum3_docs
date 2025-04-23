# brick-builder

_Source: https://doc.photonengine.com/quantum/current/technical-samples/brick-builder_

# Brick Builder Sample

![Level 4](/v2/img/docs/levels/level03-intermediate_1.5x.png)

Available in the [Gaming Circle](https://www.photonengine.com/gaming)

![Circle](/v2/img/docs/circles/icon-gaming_1x.png)

## Download

| Version | Release Date | Download |
| --- | --- | --- |
| 3.0.2 | Mar 30, 2025 | [Quantum BrickBuilder Sample 3.0.2 Build 611](https://dashboard.photonengine.com/download/quantum/quantum-brickbuilder-sample-3.0.2.zip) |

## Overview

By default, a [Quantum Map](/quantum/current/manual/maps/overview) is created during edit or build time and, because it is a static asset, does not allow changes during run-time. Having the majority of the collision geometry as static colliders is important to run performant collision detection.

This technical sample shows options on how to expand the common uses of a Quantum Map and is presented as a level editor allows players to create and modify maps in three ways.

1. **Download Map** \- downloading a Quantum map from a custom backend at runtime, **before** the Quantum game is stated.
2. **Procedural Generation** \- generating the map procedurally at runtime **before** the game is started.
3. **Level Editor** \- collaboratively changing the map at runtime using the **`DynamicMap`** API.

## Scenes

### Download From Backend

This scene shows how to save an asset (in this case, a Map) that was obtained from an external source, into the simulation before it is started.

This technique is typically used when a game has a map sharing system, where players can save their own creations to the game's backend server for example.

#### Simulating a Backend

In order to keep the sample generic, the `MockBackend` class is used to fake downloading a map from the Resources folder. The class implements the `IBackend` interface, allowing the developer to plug in their own backend.

The `IBackend` interface has two methods:

1. **DownloadMapAsync** \- Downloads a map from the backend asynchronously.

   - This version is used during real gameplay in order to not freeze the game while downloading the map.
2. **DownloadMap** \- Downloads a map from the backend synchronously.

   - This version is only used when debugging in the editor.
   - This is because the Debug Runner start callback is synchronous.

These methods return a struct of type `DownloadedMap`, which contains the Map asset along with the binary data asset containing the map's collision geometry.

This data is consumed before Quantum is started in the overriden connection behaviour. See: [Overriding Default Menu Behaviour](#overriding-default-menu-behaviour)

To make this sample work with a real backend the `IBackend` interface needs to be implemented and the `MockBackend` has to be replaced with code that calls a custom backend API.

#### ServerMapDownloader

The `ServerMapDownloader` class downloads the map from the `MockBackend` and adds it to the asset resources using the `QuantumUnityDB.Global.AddAsset` method.

This class has three methods:

1. **Download** \- Downloads the map from the backend synchronously and adds it to the asset database.
2. **DownloadAsync** \- Downloads the map from the backend asynchronously and adds it to the asset database.
3. **AddAsset** \- Adds the asset to the asset database. This is used by both of the previous methods.

   - Without registering the asset in the asset database, the asset would not be accessible via the `FindAsset` API in the simulation.

`AddAsset` works by using the `QuantumUnityDB.Global.AddAsset` method to add the provided asset to the asset database.

Example:

C#

```csharp
var asset = AssetObject.Create<MyAssetType>();
asset.Guid = SomeDeterministicGuid;
QuantumUnityDB.Global.AddAsset(asset);

```

`SomeDeterministicGuid` in this case is a deterministic guid that is generated based on the asset object. See: [Generating Deterministic Guids](#generating-deterministic-guids)

After the asset is added to the asset database, it is effectively the same as if the asset was created and added in the editor.

However, this is only valid if the game has not started yet. Once the game has started, the static asset database should not be mutated.

If you need to add or mutate assets during gameplay, you must use the `DynamicDB` API. This is covered in the [Level Editor](#level-editor) section.

#### Deterministic Guids

When adding assets at runtime, it is important for the asset to have a deterministic guid, otherwise, it can lead to a game desync.

There are two ways to generate deterministic guids:

1. **Using a constant**:

   - This is the simplest way to have a deterministic guid.
   - You must ensure that another asset is not assigned the same guid.
   - This is fine as long as you have a fixed amount of assets to add, and you are sure that the same asset will always be assigned the same guid.
   - The [Level Editor](#level-editor) scene uses this approach. See also: [LevelEditorConsts](#leveleditorconsts)
2. **Generation**:

   - The `QuantumUnityDB.CreateRuntimeDeterministicGuid` method provides an API to generate guids.
   - It uses the asset object's name as a seed to generate a deterministic guid.
   - You should use this approach if the amount of assets you need to add is not fixed.

Usage:

C#

```csharp
// create any asset
var assetObject = AssetObject.Create<MyAssetObjectType>();
// set it's name
assetObject.name = "Generated Map";
// get a deterministic guid
var guid = QuantumUnityDB.CreateRuntimeDeterministicGuid(assetObject);
// add the asset to the asset database
QuantumUnityDB.Global.AddAsset(assetObject);
// set the guid
assetObject.Guid = guid;

```

### Procedural Generation

Assets can also be procedurally generated (provided that the approach is deterministic) on each client in the application during runtime and added into the database before the game starts.

The `DynamicDB` API could also be used in this case, but it is not optimal because late joining clients would have to download the map before starting. This can be an issue if the map is very large, leading to timeouts while waiting for the download to complete.

In general, it is a good optimization step to avoid having to download anything at all, if it can be locally generated instead.

#### DeterministicBrickMapGenerator

This class uses two assets for procedural generation:

1. **Base Map** \- The base map asset that is used to generate the map.

   - This is just an empty map that is used as a base for the generated map.
2. **TextAsset Object** \- This is a text asset containing two serialized assets.

1. **Map** \- A map containing a single tree made out of bricks.
2. **Binary Data** \- The binary data asset containing the collision geometry of the tree.

This class has several methods:

1. **GenerateAndSetMapGuid** \- Generates a map, adds it to the asset database and sets the provided `RuntimeConfig` map field.

   - This is used during debugging in the editor. It is called when the debug runner is started.
2. **Generate** \- Generates a map and the binary data asset, then adds it to the asset database.

   - This is the main method used during runtime. Example: [Generation Example](#generation-example)
3. **MergeMaps** \- This is used during the generation process to merge the base map with the randomly positioned tree map.
4. **MergeMapEntities** \- This is used during the generation process to merge the entities of the base map with the entities of the tree map.

   - This is needed because the level editor utilizes entities to show the visual representation of the bricks. See: [Brick Visual](#brick-visual)
5. **MergeStaticColliders3D** \- This is used during the generation process to merge the static colliders of the base map with the static colliders of the tree map.

   - When the static collider triangles are read into memory, they are stored in an unmanaged array.


     Because of this, the triangles need to be converted to a managed array before they can be re-saved to a new binary data asset.


     This is done by using the provided utility method `Utils.FromUnManagedTriangleBuffer`.
6. **WriteBinaryData** \- Writes the combined triangle information to a `ByteStream` which is then saved to the created Binary Data.
7. **DuplicateMap** \- This duplicates a provided map asset.

   - This is needed in order to avoid modifying the base map and the tree map when combining them.

##### Generation Example:

To use the `DeterministicBrickMapGenerator` class, call the `Generate` method before the game starts, like so:

C#

```csharp
var generator = GetComponent<DeterministicBrickMapGenerator>();
generator.Generate(seed);

```

This will generate a map and add it to the asset database deterministically.

### Level Editor

This scene demonstrates adding and mutating assets at runtime by using `DynamicDB` API and the `DynamicMap` API.

The `DynamicDB` is a runtime database that allows for the addition and mutation of assets at runtime. See: [DynamicDB](/quantum/current/manual/assets/assets-simulation#dynamic-assets)

The `DynamicMap` API is an extension of the `Map` class that overrides the map's default behavior to allow for dynamic

changes to the map at runtime. Specifically, static colliders.

This approach is typically used when the map needs to be mutated in some way during gameplay.

NOTE: `DynamicMap` is not a replacement for the `Map`. It is limited by design and is not usable in the editor, only at runtime.

#### Why not use Dynamic Colliders?

Static colliders don't move, so the physics engine doesn't need to constantly recalculate their position, velocity, or acceleration. This significantly reduces the computational load on the physics engine.

#### Input

Instead of normal input, this sample scene utilizes commands. See: [Commands](/quantum/current/manual/commands)

This is because map editing does not occur every frame.

When these commands are executed, it calls the appropriate signal which is consumed by the `BrickBuildingSystem`. See: [BrickBuildingSystem](#brickbuildingsystem)

##### PlaceBrickCommand

This command is used to place a brick on the map. It contains the position, rotation, and the brick asset to place.

#### DeleteBrickCommand

This command is used to delete a brick from the map. It contains the static collider index of the brick that should be removed.

#### ClearBricksCommand

This command is used to clear all bricks from the map. It is effectively the same as the [DeleteBrickCommand](#deletebrickcommand), but it deletes all bricks at once.

#### BrickBuildingSystem

This is the main system that handles the input and applies the changes to the map.

First, in `OnInit`, the system uses `DynamicMap.FromStaticMap` to create a clone that can be mutated. Then, it is added to the `DynamicDB`.

The system listens for the `PlaceBrickCommand`, `DeleteBrickCommand`, and `ClearBricksCommand` signals.

When a signal is received, the system applies the changes with the methods below and updates the map using the `DynamicMap` API.

1. **BrickPlaced** \- Places a brick on the map.

1. This method first looks up the brick in the asset database, then checks if it can be placed at the desired position.
2. If the brick can be placed, it is added to the map using the `DynamicMap` API.
3. Depending on the collider type, the brick will be added to the map using either `AddMeshCollider` (Indicating a mesh collider) or `AddCollider3D` (Indicating a normal collider).
4. Finally, it creates a visual representation of the brick using the `BrickVisual` class. See: [Brick Visual](#brick-visual)
2. **BrickDeleted** \- Removes a brick from the map.

1. This method uses the static collider index of the brick to remove it from the map.
2. If the brick is found, it is removed from the map using either `RemoveMeshCollider` or `RemoveCollider3D` depending on the collider type.

      - One thing to note is that when a collider is removed from the map, the collider index is re-assigned to the last collider in the map.
      - This is done to keep the collider indices contiguous for the physics engine.
3. Finally, if all previous steps are successful, the visual representation of the brick is also removed.

#### Debugging the Map

The `BrickBuildingSystem` also has a debug mode which can be toggled by enabling the `DebugDrawBrickGrid` field in the RuntimeConfig.

![Debug Draw Level Editor](/docs/img/quantum/v3/technical-samples/brick-builder/level-editor-debug.png)

When this is enabled, the system will draw the grid of the map in the scene view.

#### Bricks

Bricks are assets which are the building blocks of the map. They contain the information about grid size, mesh collider type and what visual prefab to use when spawned.

#### Brick Visual

This sample uses one prefab to represent all bricks in the map.

This is done by reading the mesh from the associated brick asset and setting it on the prefab at runtime.

#### BrickRotation

The `BrickRotation` enum is used to represent the rotation of the brick. This is because the brick only has four possible rotations.

#### Brick Grid

The `BrickGrid` class is used to represent the grid that contains all the bricks.

It stores the bricks in a dictionary where the key is the position of the brick in the grid.

Because of the BrickGrid's potential size, it is not feasible to serialize the grid in the simulation.

Because of this, the grid is serialized and stored in the `FrameContext`.

#### Creating a Brick

To create a brick you need to:

1. Create a new ScriptableObject asset in the `Bricks` folder.
2. Populate the asset with the necessary data, such as the mesh to use, the grid size and any variants.
3. Add the brick to the `AllBricks` asset.

   - This is simply an asset that contains a list of all bricks that are available to the player.

The brick will now show up in the editor scene and can be placed in the map.

#### File I/O

The `LevelEditorFileController` class is used to save and load the map to and from a file. It simply serializes the

brick information to a `BitStream` and writes it to a file in the `Application.persistentDatapath`.

## Additional Info

### Overriding Default Menu Behaviour

The `LevelEditorConnectionBehaviour` class overrides the `ConnectAsync` method in the demo menu's `ConnectionBehaviour` class to allow the level editor to perform actions before

connecting to the server.

It can do two things:

1. **Download Map** \- Downloads the map from the backend before connecting. This is used in the [Download From Backend](#download-from-backend) scene.
2. **Generate Map** \- Generates the map procedurally before connecting. This is used in the [Procedural Generation](#procedural-generation) scene.

### Previous Version

This sample is the successor to the Quantum 2.1 Level Editor. It has been updated to Quantum 3.0 and has been improved

in many areas, such as:

1. Box Collider Support
   - In the previous version, the sample only supported mesh colliders.
   - Due to this limitation, the old sample was not able to support a lot of objects.
2. Irregular grid
   - In the previous version, you could only have grid objects be 1x1x1 in size.
   - In this version, you can have grid objects of any size.
3. DynamicMap
   - Instead of writing the static collider editing from scratch, we now use the new `DynamicMap` API.
4. Downloading Static Assets
   - In Quantum 2.1, downloading new static assets at runtime was very complex.
   - In Quantum 3.0, this has been simplified by using the new `QuantumUnityDB.Global.AddSource` API.

## Third Party Assets

This sample includes third-party free and CC0 assets. The full packages can be acquired for your own projects at their respective site:

- [Brick Kit](https://kenney.nl/assets/brick-kit) by Kenney

Back to top

- [Download](#download)
- [Overview](#overview)
- [Scenes](#scenes)

  - [Download From Backend](#download-from-backend)
  - [Procedural Generation](#procedural-generation)
  - [Level Editor](#level-editor)

- [Additional Info](#additional-info)

  - [Overriding Default Menu Behaviour](#overriding-default-menu-behaviour)
  - [Previous Version](#previous-version)

- [Third Party Assets](#third-party-assets)