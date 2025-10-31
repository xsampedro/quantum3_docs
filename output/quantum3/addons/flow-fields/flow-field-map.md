# flow-field-map

_Source: https://doc.photonengine.com/quantum/current/addons/flow-fields/flow-field-map_

# Flow Fields Map

### Map Hierarchy

The `FlowFieldMap` is subdivided into smaller chunks named `FlowFieldController`. Each controller is subdivided into tiles. It is possible to setup more than one map to work at same runtime, but they are not connected.

![Map Hierarchy](/docs/img/quantum/v2/addons/flow-fields/map-hierarchy-1.png)
Example of a map with dimensions 32x32 and controller size 8.
### FlowFieldMapUtility

`FlowFieldMapUtility` is a static helper class with useful methods like LineOfSight, converting world positions to map location, etc.

### Modifying a Flow Field map

There are two ways to modify an existing map:

1. Setting original costs - used for permanent cost changes in single tiles on the map:

C#

```csharp
public void SetOriginalTileCost(Frame frame, Vector2Byte location, byte cost)

```

2. Area cost modifiers - used for temporary changes, for example like in RTS buildings (which can be created and destroyed). It is applied to a group of tiles based on the parameterized area. Adding a modifier retrieves it's integer ID, which can further be used to easily remove the same modifier:

C#

```csharp
public int AddCostModifier(Frame frame, FPVector2 minPosition, FPVector2 maxPosition, byte cost)
public bool RemoveCostModifier(Frame frame, int modifierID)

```

## Map Creation

The `FlowFieldMap` is created in runtime and stored in the FrameContext.

### Parameters

- **Dimensions** \- (X, Y) size of the map. Maximum dimension supported is 256x256;
- **Offset** \- (X, Y) the offset in world position of the center of the map.
- **Tile Size** \- size of each tile;
- **Controller Size** \- size of the map's subdivisions (the map _Dimensions_ must be a multiple of the _Controller Size_). Recommended _Controller Size_ is between 8 and 20;
- **Max Portal Length** \- maximum length of a portal connecting two neighboring controllers. If _Max Portal Length_ is greater than _Controller size_, there can be multiple portals between two controllers. More portals means greater precision but it is slower CPU-wise;
- **Cost Field** \- costs of individual tiles.

C#

```csharp
var ffMap = new FlowFieldMap(new Vector2Int(16, 16), FP._2, 8, 4, COSTS);
ffMap.Initialize(frame.SimulationConfig.ThreadCount, false);
frame.Context.FlowFieldMap = ffMap;

```

## Map Data

### Frame Context

The `FlowFieldMap` can hold lots of data in case of large maps, so it is not suitable to keep such data in the Frame.

FlowFieldMap is stored in the FrameContext which means changes to it has to be done very carefully. To prevent any desyncs between clients it is crucial to modify it only in Verified frames.

Having too big frame size can lead to lower performance and serialized data for rejoin/late join could be too big to transfer (FlowFieldMap has custom serialization which mitigates this issue).

### Late join and Reconnect

Due to the fact that `FlowFieldMap` is not part of the frame, it implements custom serialization - see `FlowFieldMap.Serialize()`.

By default, the addon comes with a call for such serialization method in `Frame.User.cs`, on the partial `SerializeUser` implementation. When using the addon, please make sure this code is not removed. Either use it as provided on the addon, or add a call to `FlowFieldMap.Serialize()` in your `SerializeUser` method if you already have custom code in it.

## Multiple Maps

This addon supports simulating multiple maps at same time. But the maps are not connected which means an entity will not change or move between the maps automatically. It is possible to include more maps directly in the Runtime Config being used. For initializing multiple maps, create a partial implementation of `FrameContextUser` and include a function for initialization:

C#

```csharp
public partial class FrameContextUser
{
    public void InitFlowFieldMaps(Frame frame, AssetRef<TileMapSetup>[] maps, int threadCount)
    {
        _flowFieldMaps = new FlowFieldMap[maps.Length];
        for (int i = 0; i < maps.Length; i++)
        {
            var tileMapSetup = frame.FindAsset<TileMapSetup>(maps[i].Id);
            _flowFieldMaps[i] = new FlowFieldMap(tileMapSetup.Offset, tileMapSetup.Dimensions, tileMapSetup.TileSize, tileMapSetup.ControllerDimension, 6, tileMapSetup.CostField, threadCount);
            _flowFieldMaps[i].Initialize(frame.SimulationConfig.ThreadCount, false);
        }
        CheckMap();
    }
}

```

Then it is possible to initialize it from any system:

C#

```csharp
// Initialize it in any system
public override void OnInit(Frame frame)
{
    var threadCount = frame.SimulationConfig.ThreadCount;
    frame.Context.InitFlowFieldMaps(frame, frame.RuntimeConfig.TileMapSetup, threadCount);
}

```

## Map Baking

The maps are baked automatically during the `OnBake` callback. It is possible to customize the baking process for different maps by using the `FlowFieldMapBaker` component on the root of the GameObject that contains all the relevant colliders for that map in its hierarchy. In this component it is possible to setup:

- Bake Mode:
  - Default: it will include all the colliders of the scene inside the map's range and offset.
  - Hierarchy: it will include only the colliders in the hierarchy of the `FlowFieldMapBaker` game object.
  - Layer: It will include all the colliders of the setup layer.
- Tile Map Setup: the reference of the tile map asset that must be baked.
- Layer Mask: the mask that the bake mode will use to select the colliders.

It is possible to use two maps in overlap to create, for example, ground units that are blocked by simple obstacles and flying units that can pass through simple obstacles. To do so, setup each map with the same offset but different layers.

Back to top

- [Map Hierarchy](#map-hierarchy)
- [FlowFieldMapUtility](#flowfieldmaputility)
- [Modifying a Flow Field map](#modifying-a-flow-field-map)

- [Map Creation](#map-creation)

- [Parameters](#parameters)

- [Map Data](#map-data)

- [Frame Context](#frame-context)
- [Late join and Reconnect](#late-join-and-reconnect)

- [Multiple Maps](#multiple-maps)
- [Map Baking](#map-baking)