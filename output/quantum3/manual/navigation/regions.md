# regions

_Source: https://doc.photonengine.com/quantum/current/manual/navigation/regions_

# Using Navmesh Regions

Quantum navmesh regions are pre-defined (edit time) areas of the navmesh that can be turned on/off or filtered per agent during runtime with very little performance overhead. Regions are a compromise to Unity dynamic navmesh carving based on performance considerations for deterministic roll-backs.

Two region generation modes are available.

|     |     |     |
| --- | --- | --- |
|  | **Simple** | **Advanced** |
| Maximum region count | 29 | 512 |
| Tooling | Directly uses Unity navigation area ids | Requires setting up extra region scripts inside the level geometry |

Select the `QuantumMapNavMeshUnity` script to change the mode.

The `Walkable` area is always converted to region `0` and is called `MainArea`.

All regions for a map are accessible in this list `Map.Regions`. A region id is an integer and refers to the index of a region inside that map which can be retrieved using the `Map.RegionMap` lookup.

![Region Import Modes](/docs/img/quantum/v3/manual/navigation/region-modes.png)## Importing Navmesh Regions - Simple Mode

Select a scene object, add the `NavMesh Modifier` script to it, enable `Override Area` and set the desired area type. Define more Unity navigation areas types in the Unity Navigation Window.

![NavMesh Modifier](/docs/img/quantum/v3/manual/navigation/region-simple-modifier.png)

Select the `QuantumMapNavMeshUnity` script and enable the `Simple` region import mode. Then toggle the Unity Areas that should be converted to Quantum regions.

![NavMesh Modifier](/docs/img/quantum/v3/manual/navigation/region-simple-import.png)

Bake the Quantum navmesh and verify the imported region using the NavMesh gizmos. regions are rendered in magenta.

![Baked simple region](/docs/img/quantum/v3/manual/navigation/region-baked.png)

Select the `QuantumMapEntity` to see the imported region names.

![NavMesh map](/docs/img/quantum/v3/manual/navigation/region-simple-map.png)## Importing Navmesh Regions - Advanced Mode

Instead of using the Unity navigation area ids directly (because their number is limited) parts of the Unity navmesh are "marked" with a navigation area during baking using the `QuantumNavMeshRegion` script. The marked navmesh triangles are expanded into islands and then matched back to their origin `QuantumNavMeshRegion` bounding box to finally receive the correct region id.

![Unity navmesh](/docs/img/quantum/v3/manual/navigation/region-unity-navmesh.png)

Step 1: Open the Unity Navigation windows (`Window > AI > Navigation`) and add a new navigation area named `QuantumRegion` for example.

![Unity Navigation window](/docs/img/quantum/v3/manual/navigation/region-advanced-unity-areas.png)

Step 2: Select the `MapNavMeshUnity` script, set the region import mode to `Advanced` and enable the Unity navigation area `QuantumRegion`.

Because navmesh triangles are not generated 100% exact, `RegionDetectionMargin` defines a margin during the fitting. The value can be increased when regions fail to be imported. But when it becomes too large there may be problems detecting neighboring regions.

![Add Region Area](/docs/img/quantum/v3/manual/navigation/region-advanced-add-areas.png)

Step 3: Create regions caster scripts.

Select or create a GameObject with a `MeshRenderer` and attach the `MapNavMeshRegion` script.

![Region Setup](/docs/img/quantum/v3/manual/navigation/region-advanced-setup.png)

Unity uses a `MeshRenderer` to project areas onto the navmesh. Create a GameObject with a `MeshRenderer` and attach the `MapNavMeshRegion`. The assumption is that the Navmesh surface is configured to collect Render Meshes.

The `MeshRenderer` and region script only have to be active during the baking.

![Region Setup](/docs/img/quantum/v3/manual/navigation/region-advanced-script.png)

The region `Id` marks the region with a name. Multiple `MapNavMeshRegion` can use the same name and contribute to that region. The `Id` will be registered in the map under `Map.RegionMap`. A region index refers to the position of the region inside that array.

Set `CastRegion` to `CastRegion`.

Add a `NavMesh Modifier` script to the GameObject, enable `Override Area` and set `QuantumRegion` as area type. In `Unity 2022` and before the `NavMeshHelper` on the region script can be used to toggle the area.

![Region Setup](/docs/img/quantum/v3/manual/navigation/region-advanced-modifier.png)

When regions are close to each other, choose a different area type for each region. Region `A` uses `QuantumRegion` and `B` uses `QuantumRegion2`.

![Region Setup](/docs/img/quantum/v3/manual/navigation/region-advanced-multiple.png)

Step 4: Bake the map and navmesh

![Baked simple region](/docs/img/quantum/v3/manual/navigation/region-baked.png)### Region API

Get the region `Foo` and disable it globally.

C#

```csharp
var regionId = frame.Map.RegionMap["Foo"];
frame.NavMeshRegionMask->ToggleRegion(regionId, false);

```

The `NavMeshPathfinder` component also has a RegionMask that can be used to toggle off regions for individual agents.

C#

```csharp
var regionId = frame.Map.RegionMap["Foo"];
var agent = f.Unsafe.GetPointer<NavMeshPathfinder>(entity);
agent->RegionMask.ToggleRegion(regionId, false)

```

The runtime gizmos show enabled and disabled regions.

|     |     |
| --- | --- |
| ![Inactive Regions](/docs/img/quantum/v3/manual/navigation/region-active.png) | ![Inactive Regions](/docs/img/quantum/v3/manual/navigation/region-inactive.png) |

When changing the map, the region state can be reset by running `FrameBase.ClearAllNavMeshRegions()` during the `ISignalOnMapChanged` signal.

C#

```csharp
public class ResetRegionsSystem : SystemSignalsOnly, ISignalOnMapChanged {
    public void OnMapChanged(Frame frame, AssetRefMap previousMap) {
        frame.ClearAllNavMeshRegions();
    }
}

```

The `NavMeshRegionMask` object controls what regions are enabled using an internal bitset.

|     |     |
| --- | --- |
| **Default** | Creates a region mask with all regions enabled |
| **MainArea** | Creates a region mask that includes only the MainArea |
| **Empty** | Creates a region mask that has no enabled regions |
| **HasValidRegions** | Returns `true` if the mask has at least one valid region set including the main area. |
| **HasValidNoneMainRegion** | Returns `true` if the mask includes a region other than the main area. |
| **IsMainArea** | Checks if the mask only contains the main area |
| **Create(int region)** | Create a NavMeshRegionMask with one region enabled |
| **Create(int\[\] regionList)** | Create a mask using a list of region ids |
| **ToggleRegion(int region, bool enabled)** | Toggle a region on or off using the region id |
| **IsRegionEnabled(int region)** | Test if a region is enabled |
| **Clear()** | Reset the mask and sets all regions to enabled |

Back to top

- [Importing Navmesh Regions - Simple Mode](#importing-navmesh-regions-simple-mode)
- [Importing Navmesh Regions - Advanced Mode](#importing-navmesh-regions-advanced-mode)
  - [Region API](#region-api)