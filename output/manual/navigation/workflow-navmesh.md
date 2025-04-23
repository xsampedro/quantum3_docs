# workflow-navmesh

_Source: https://doc.photonengine.com/quantum/current/manual/navigation/workflow-navmesh_

# Importing A Unity Navmesh

The Quantum navmesh generation is build on top of the Unity navmesh pipeline. A Unity navmesh triangulation is imported and baked into Quantum navmesh binary format.

The baking writes world positions into the Quantum navmesh because only navmeshes located in the origin are supported. Similar to Quantum physics it is encouraged that gameplay takes place near the origin to prevent fixed point precision issues.

## Importing A Unity Navmesh

- Add the AI Navigation package ```
com.unity.ai.navigation
```

to the Unity project using the ```
Unity Package Manager > Unity Registry
```

.
![AI Navigation package](/docs/img/quantum/v3/manual/navigation/workflow-package.png)
- Create and set up a Unity navmesh surface.
![Create a surface](/docs/img/quantum/v3/manual/navigation/workflow-create-surface.png)
- Create a new GameObject under the map ```
QuantumMapData
```

and add a ```
QuantumMapNavMeshUnity
```

script.
![Create script](/docs/img/quantum/v3/manual/navigation/workflow-createscript.png)
- Add the surface GameObject to the ```
Nav Mesh Surfaces
```

list.
![Add a Surface](/docs/img/quantum/v3/manual/navigation/workflow-add-surface.png)
- Select the map object and enable ```
Everything
```

on the ```
Bake All Mode
```

options. Then press **Bake All**. A navmesh imported message will be logged on the console.


```
```
Imported Unity NavMesh 'Navmesh', cleaned up 7 vertices, found 1 region(s), found 0 link(s)

```


```

- The Quantum navmesh asset is baked into two asset files. The navmesh asset and a second binary data asset (```
\_data
```

suffix), which are placed next to the map asset.
![Navmesh assets](/docs/img/quantum/v3/manual/navigation/workflow-navmesh-asset.png)
- After baking the navmesh asset is referenced in the ```
QuantumMapData
```

script under ```
Nav Mesh Links
```

.
![Map navmesh references](/docs/img/quantum/v3/manual/navigation/workflow-map-navmesh-ref.png)
- Select the map object to show the Quantum navmesh gizmos. Open the Quantum Gizmo Overlay Menu to toggle different types of gizmo drawing (right-click on the scene tab > Overlay Menu > Quantum Gizmos).
![Show Quantum navmesh gizmos](/docs/img/quantum/v3/manual/navigation/workflow-gizmos.png)
- Optionally toggle navmesh building on different auto build mode triggers in ```
QuantumEditorSettings
```

.
![Auto build modes](/docs/img/quantum/v3/manual/navigation/workflow-auto-build-mode.png)
- Add the ```
Quantum.Core.NavigationSystem
```

to the ```
SystemsConfig
```

that is set on the ```
RuntimeConfig
```

when starting the game. The default config has the system correctly added and enabled.
![Auto build modes](/docs/img/quantum/v3/manual/navigation/workflow-navigation-system.png)

## Import Settings

The conversion from a Unity navmesh to a Quantum navmesh uses import settings that can be customized using the ```
QuantumMapNavMeshUnity
```

 script.

|     |     |
| --- | --- |
| **WeldIdenticalVertices** | The Unity NavMesh is a collection of non-connected triangles. This option is very important and combines shared vertices. |
| **WeldVertexEpsilon** | Don't make the epsilon too small, vertices required to fuse can be missed, also don't make the value too big as it will deform your navmesh. |
| **DelaunayTriangulation** | This option will post processes the imported Unity navmesh with a Delaunay triangulation to produce more evenly distributed triangles (it reorders long triangles). |
| **DelaunayTriangulationRestrictToPlanes** | On 3D navmeshes the Delaunay triangulation can deform the navmesh on slopes while rearranging the triangles. This behavior is also noticeable on Unitys navmesh and can affect a game when the navmesh height is used for gameplay (e.g. walking on the Navmesh).<br> Check this option to restrict the triangulation to triangles that lie in the same plane. |
| **FixTrianglesOnEdges** | Imported vertices are sometimes lying on other triangle edges, which leads to unwanted border detection. With this option such triangles are split. |
| **FixTrianglesOnEdgesEpsilon** | Large navmeshes may require to increase this value (e.g. to 0.001) when false-positive borders are detected. Min = float.Epsilon. |
| **FixTrianglesOnEdgesHeightEpsilon** | Make the height offset considerably larger than FixTrianglesOnEdgesEpsilon to better detect degenerate triangles. If the navmesh becomes deformed chose a smaller epsilon. Min = float.Epsilon. Default is 0.05. |
| **LinkErrorCorrection** | Automatically correct navmesh link position to the closest triangle by searching this distance (default is 0). |
| **ClosestTriangleCalculation** | Areas in the map grid without a navmesh will need to detect nearest neighbors. This computation is very slow. The SpiralOut option will be much faster but fallback triangles can be null. |
| **ClosestTriangleCalculation Depth** | Number of cells to search triangles into each direction when using SpiralOut. |
| **EnableQuantum\_XY** | Only visible when the QUANTUM\_XY define is set. Toggle this on and the navmesh baking will flip Y and Z to support navmeshes generated in the XY plane. |
| **MinAgentRadius** | The minimum agent radius supported by the navmesh. This value is the margin between the navmesh and a visual border. The value is overwritten by retrieving it from Unity navmesh bake settings (or the surface settings) when baking in the Editor. |
| **ImportRegionMode** | Disable or change the import region mode. See the section Using Navmesh Regions for more details. Default is Simple. |
| **RegionDetectionMargin** | The artificial margin is necessary because a navmesh triangulation does not fit the source size very well. The value is added to the navmesh area and checked against all Quantum Region scripts to select the correct region id. |
| **RegionAreaIds** | The Unity area ids that will be converted to Quantum regions. Either using the area name (Simple) or an explicit region name (Advanced) |

## Navmesh Settings on the Map

The ```
QuantumMapData
```

script has two relevant navmesh settings: the serialization type and the navmesh grid and world size. Both settings get applied and copied to the navmesh asset during baking.

### NavMesh Serialize Type

The setting controls how much of the navmesh asset is serialized with the asset and how much is generated at load time. It can be used to reduce the asset size, which could be interesting when the asset is created and send to clients at runtime (see [Custom Navmesh Generation](customized-navmesh)).

![Navmesh Area Gizmos](/docs/img/quantum/v3/manual/navigation/workflow-mapdata-serialization.png)

|     |     |     |
| --- | --- | --- |
| Full | Serializes the full data but computes metadata after loading.<br> <br> Default value. | Medium size |
| FullWithMetaData | Serializes the full data and meta data (e.g. normals). | Largest size |
| BakeDataOnly | Serializes only the ```<br>NavMeshBakeData<br>```<br> and bakes the navmesh NavMesh at runtime. | Smallest size |

### NavMesh Area and Grid

Gizmos for ```
NavMesh Area
```

 and ```
NavMesh Grid
```

are rendered to the scene when toggled on in the Quantum Navmesh Overlay Menu.

![Navmesh Area Gizmos](/docs/img/quantum/v3/manual/navigation/workflow-mapdata-navmesh-area-gizmos.png)

The grid is configured by selecting the ```
QuantumMapData
```

 script under the ```
NavMesh Settings
```

headline. The number of grid cells in ```
X
```

 and ```
Y/Z
```

direction and the grid node size (cell dimension in Unity units per cell) result in the overall grid size.

![Navmesh Area Map Data](/docs/img/quantum/v3/manual/navigation/workflow-mapdata-navmesh-area-mapdata.png)

The navmesh grid area needs to encompass the complete navmesh (blue box on the left screenshot).

The grid cells (yellow square on the right screenshots) require a reasonable size to not include too many individual navmesh triangles. It's a trade-of between navmesh data size and performance. The smaller the cells the larger the data structure but less triangles need to be touched during pathfinding.

|     |     |
| --- | --- |
| ![Navmesh Area](/docs/img/quantum/v3/manual/navigation/workflow-mapdata-navmesh-area.png) | ![Navmesh Area Grid](/docs/img/quantum/v3/manual/navigation/workflow-mapdata-navmesh-area-grid.png) |

Even if the navmesh world position is not near the origin it still has to be inside the grid area.

## Custom Navmesh Baking Callback

The navmesh baking process can be extended by callbacks extending the ```
MapDataBakerCallback
```

 class.

Add the following script into ```
QuantumUser/Editor
```

.

The ```
QuantumMapBakeAssembly
```

 attribute is required for the baking process to find the callbacks.

C#

```
```csharp
\[assembly: Quantum.QuantumMapBakeAssembly\]

namespace Quantum.Editor
{
 using System.Collections.Generic;
 using UnityEngine;

 public class NavmeshBakeCallback : MapDataBakerCallback
 {
 public override void OnBeforeBakeNavMesh(QuantumMapData data)
 {
 // Before any navmesh baking takes place.
 }

 public override void OnCollectNavMeshBakeData(QuantumMapData data, List<NavMeshBakeData> navMeshBakeData)
 {
 // Unity navmesh surfaces have been imported and bake data is already filled out.
 Debug.Log($"Found {navMeshBakeData.Count} navmesh bake data");
 }

 public override void OnCollectNavMeshes(QuantumMapData data, List<NavMesh> navmeshes)
 {
 // Quantum navmesh have been baked.
 }

 public override void OnBakeNavMesh(QuantumMapData data) {
 // Quantum navmeshes have been saved to assets.
 }

 // abstract methods have to be implemented but not needed here
 public override void OnBake(QuantumMapData data) { }
 public override void OnBeforeBake(QuantumMapData data) { }
 }
}

```

```

Back to top

- [Importing A Unity Navmesh](#importing-a-unity-navmesh)
- [Import Settings](#import-settings)
- [Navmesh Settings on the Map](#navmesh-settings-on-the-map)

  - [NavMesh Serialize Type](#navmesh-serialize-type)
  - [NavMesh Area and Grid](#navmesh-area-and-grid)

- [Custom Navmesh Baking Callback](#custom-navmesh-baking-callback)