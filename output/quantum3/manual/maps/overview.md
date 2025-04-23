# overview

_Source: https://doc.photonengine.com/quantum/current/manual/maps/overview_

# Overview

Maps are Quantum assets representing pre-made collections of entity prototypes, static collider geometries and referenced navmesh assets.

Quantum maps are the equivalent of scenes in Unity and contain gameplay content for the deterministic simulation.

They are coupled 1-to-1 with a Unity scene which they are usually baked from during edit time. The scenes may contain view game objects (EntityViews) that can be associated with Quantum simulation entities at run-time.

## Map API

The initial map that a Quantum session starts with must set by ```
RuntimeConfig.Map
```

.

Set a new map during run-time from a system update method for example:

C#

```cs
public override void Update(Frame frame) {
// only change the map during verified frames
if (frame.IsVerified) {
// nextMap is of type AssetRef<Quantum.Map>
frame.Map = frame.FindAsset(nextMap);
}
}

```

Loading the initial map and changing a map will invoke the ```
ISignalOnMapChanged
```

signal on all systems.

C#

```cs
public void OnMapChanged(Frame frame, AssetRef<Map> previousMap) {
 // new map is frame.Map
}

```

## Map Entities

```
Map Entities
```

are a list of Quantum entity prototypes that are created and configured when the map is loaded by the ```
EntityPrototypeSystem
```

.

During map baking for example all entity prototype scripts in a Unity scene are converted into the referenced map ```
Map Entities
```

.

The entity prototypes can optionally have an entity view scripts which will be linked to entities during run-time by the ```
QuantumEntityViewUpdater
```

 using the ```
MapEntityLink
```

component.

## Entity Prototype System

This build-in system handles the lifecycle of ```
Map Entities
```

. When a map is loaded or changed the system reacts to ```
SignalOnMapChanged
```

and perform the following tasks:

- Delete all entities that have a ```
MapEntityLink
```

component, which indicates that they were created by a previous map.
- Creates new Quantum entities from the entity prototypes (Map Entities) of the new map that is loaded.
- Add a ```
MapEntityLink
```

component to each entity created this way which saves the index of the entity prototype it was created from.

## Map Entity Link

This built-in component is added to each entity that was created from ```
Map Entity
```

 prototypes. The ```
Index
```

field on the component references the array index of the prototype on a map that it was created from.

The ```
QuantumEntityViewUpdater
```

 uses this information to identify the associated entity view using the game object list ```
MapEntityReferences
```

, which order is identical to the ```
Map Entities
```

.

## Static Colliders

Static colliders are 2D and 3D physics colliders that do not move during the lifetime of a map and are saved on the ```
StaticColliders2D
```

and ```
StaticColliders3D
```

collections on the map.

Static colliders have a considerable performance advantage compared to dynamic colliders. Read

more about [Static Physics Colliders](/quantum/current/manual/physics/statics).

## User Asset

Optionally associate any Quantum asset with a map and load it like this:

C#

```cs
public override void OnInit(Frame frame) {
if (frame.TryFindAsset(frame.Map.UserAsset.Id, out FooAsset mapAsset)) { }
}

```

## Physics Settings

Read more about the [Physics Map Settings](/quantum/current/manual/physics/settings).

## NavMesh Settings

These settings are required when navmeshes are generated for the map. The grid size affects how performant searching for neighbor triangles is for the navigation system. It's a trade-of between size and performance. Toggle the Quantum navmesh gizmos NavMesh Area and NavMesh Grid to make the cells visible.

```
Grid Size X
```

\- the number of grid cells in x dimension.

```
Grid Size Y
```

\- the number of grid cells in the z dimension.

```
Grid Node Size
```

\- The number of Unity units per grid cell, must be a multiple or 2.

Read more about the [Navigation Workflow](/quantum/current/manual/navigation/workflow-navmesh).

The ```
NavMeshLinks
```

are the final collection of navmeshes assets created for this map, which will be loaded when the map is loaded.

The ```
Regions
```

 are a part of the navmeshes, where pre-defined regions can be toggled on and off during run-time. Read more about [Navigation Workflow](/quantum/current/manual/navigation/regions).

## Plans For Quantum 3.1

For future Quantum versions additional features are planned to support more map features:

- Additive map loading
- Currently, the dynamic map system needs to re-initialize the physics system if a collider was added or removed. This is not ideal for performance. The plan is to optimize this process so that the physics system does not need to be re-initialized.

Back to top

- [Map API](#map-api)
- [Map Entities](#map-entities)
- [Entity Prototype System](#entity-prototype-system)
- [Map Entity Link](#map-entity-link)
- [Static Colliders](#static-colliders)
- [User Asset](#user-asset)
- [Physics Settings](#physics-settings)
- [NavMesh Settings](#navmesh-settings)
- [Plans For Quantum 3.1](#plans-for-quantum-3.1)