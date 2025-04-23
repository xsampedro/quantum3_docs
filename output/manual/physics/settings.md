# settings

_Source: https://doc.photonengine.com/quantum/current/manual/physics/settings_

# Settings

## Overview

The Physics settings can be edited in the `Map` Asset associated with a Scene in its `Quantum Map Data` script, and in the `Simulation Config` Asset linked in the `Quantum Runner` script.

The settings in the `Map` are specific to a given scene, while the `Simulation Config` can be shared among multiple scenes.

### Map Data

The scene's playable area related settings can be found in the scene's `QuantumMapData` script or the `Q Asset Map` slotted in the Asset field.

![Adjust the World Size accordingly to the Playable Area](https://doc.photonengine.com/docs/img/quantum/v3/manual/physics/physics-settings-mapdata.png)
The Physics settings as seen in the MapData script in the Scene.

Aside from `World Size`, tweaking these settings should only be a concern if the physics simulation is a bottleneck in the game.

| Setting | Description |
| --- | --- |
| **World Size** | The physics scene size in the bucketing axis. The broad phase is clamped by a bounding box of all physics entries between `-WorldSize/2` to `WorldSize/2`. It is therefore crucial to ensure the world is big enough to encompass all entities. If an entity is outside the world, it will cost performance as it is added to either the first or last bucket. Everything outside the bounding box is considered to be at the world's edge, from the physics engine perspective, which will result in false collision candidates. <br> <br>In the non bucketing axis, the physics world is only limited by the value range of `FP.UsableMin` to `FP.UsableMax`. |
| **Buckets Count** | The amount of buckets used in the broad phase, which are resolved in parallel. Use a reasonable amount according to the amount of physics entries (colliders). Too many buckets and the handling overhead increases without any performance gain because there are only few entries in each one; too few buckets and there will be an excessive amount of entries in each, slowing down the broad phase performance. |
| **Buckets Subdivisions** | Regular queries (overlaps and raycasts) use a stabbing approach for checking as few entries as possible in the buckets subdivisions. Tweak the number in accordance with the expected amount of entries and regular queries performed. Too many subdivisions will add overhead without performance, while too few will result in queries taking longer to resolve, because they will have to check too many entries. |
| **NOTE on Buckets Count & Buckets Subdivisions**<br>The default buckets count and bucket subdivisions values (16 buckets with 8 subdivisions) are usually good for up to 1~2K entries. You should thus not have to worry about tweaking them unless the physics are the bottleneck of the game. In that case, use the Task Profiler for evaluating the performance and tweak the values based on the findings (broad phase and regular queries resolution respectively). |
| **Bucketing Axis** | Physics entries are put into buckets according to their position in the bucketing axis. |
| **Sorting Axis** | The queries in a bucket are sorted according to their position in the sorting axis. |
| **NOTE on Bucketing Axis & Sorting Axis**<br>The Y-Axis represents the vertical axis of the physics simulation. In 2D this is equal to the Y-Axis, whereas in 3D the Y-Axis is mapped to the Z-Axis as the 3D space partitioning is performed on the XZ-Plane.<br> <br>Choose these a bucketing and sorting axis based on how the entries are spread out in the world. Selecting a different axis for bucketing and sorting (e.g. X-Y or Y-X) is good for uniformly spread entries across that plane. If entries are concentrated in one axis, consider using the same axis for both bucketing and sorting. |
| **Triangle Mesh Cell Size** | Defines the size of the cells into which the 3D triangle soup is divided. This number should be adapted based on how dense the meshes' triangles density to get a reasonable amount of triangles per cell. <br> <br>For better visualization enable related fields in the `QuantumEditorSettings` asset's Collider gizmos section.<br> <br>This will affect the performance of both the broad phase and regular queries. Use the Task Profiler to analyse the performance and find the most suitable number for the game. |

### Simulation Config

The `SimulationConfig` data asset contains an extensive set of settings for the physics engines:

![Physics Settings on SimulationConfig](https://doc.photonengine.com/docs/img/quantum/v3/manual/physics/physics-settings-simulationconfig.png)
SimulationConfig Asset.

**Layers** and the corresponding **Layer Collision Matrix** can be imported from Unity ones. Once imported, the _Collision Matrix_ can be edited directly in the settings.

## Optimization Tips

In this section we are covering some general considerations to take when optimizing the physics settings for enhanced performance:

- _Collision Matrix_, make sure to only **enable** collisions between layers that actually require to be checked against each other;
- _Angular Velocity_ (physics controlled rotation), **disabling** this option leads to faster and more stable physics simulation;
- _Kinematic Entities_, use kinematic entities rather than dynamic entities whenever possible. Kinematics do not check for collisions against each other, unless one of them is a trigger kinematic.
- _Raycasts_, use reasonable **distances** for rays to prevent them from becoming a bottleneck.
- _PhysicsBody_, **enabling** _resting bodies_ in the settings allows resting entities to be excluded from collisions checks and reduce the load on the collision detection system. Resting bodies can be awaken again by either another moving body, or by code.
- _Thread Count_, **tweaking** this options allows to raise the amount of threads available to the Quantum Simulation at runtime.
- _Profiler_, run it on the code systems before and during performance tweaks. Bottlenecks are often tied to custom code rather than the physics engine. Furthermore, the profiler helps to identify which settings work best under the game specific load. Find more information in the [Profiling](/quantum/current/manual/profiling) documentation.

Back to top

- [Overview](#overview)

  - [Map Data](#map-data)
  - [Simulation Config](#simulation-config)

- [Optimization Tips](#optimization-tips)