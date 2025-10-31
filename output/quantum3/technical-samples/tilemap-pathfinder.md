# tilemap-pathfinder

_Source: https://doc.photonengine.com/quantum/current/technical-samples/tilemap-pathfinder_

# Tilemap Pathfinder

Level

INTERMEDIATE

Topology

**DETERMINISTIC**

## Overview

This sample demonstrates an implementation of a generic pathfinder system based on tiles which is simpler and cheaper than navigating on navmeshes and can fit a lot of tile-, grid- and hex grid-based games.

This sample shows possibilities of steering using the transform, a kinematic character controller and custom pathing. The systems in the project are open to customization and can be adapted to different gameplay styles but are limited to flat (2D) levels using the X and Z axis. The Y axis will be ignored.

**NOTE:**: This system uses a 3D axis from Transform3D internally and does NOT work using entities with Transform2D components.

This pathfinder system contains three parts:

- **Tile Map Data**: The number of tiles, their scale and neighboring tiles. All this data is baked during edit mode using a tool.

- **Pathfinder Function**: This sample uses the A\* algorithm to find a list of points in the tilemap. The points are converted to indexes and stored in an array.

- **Agent**: The entity in the game world that can move through the map using the list of waypoints.


**NOTE:** This is _NOT_ a game.

## Technical Info

The project has been developed with:

- Unity `6.0 (6000.0.47f1)`

- Quantum `3.0.3 Stable 1685`


## Download

| Version | Release Date | Download |
| --- | --- | --- |
| 3.0.3 | May 16, 2025 | [Quantum TilemapPathfinder 3.0.3](https://downloads.photonengine.com/download/quantum/quantum-tilemappathfinder-3.0.3.zip?pre=sp) |

## Sample Scenes

The project has four samples which exemplifies how to use the `TilePathfinder` component. The features are shared between the systems.

### 1) Simple Click Movement

This scene demonstrates a simple entity with a `TilePathfinder` and a `PlayerLink` component. Click anywhere to move it. Open the Quantum Code project to view the `ClickMoveSystem`.

![Alternative Text](/docs/img/quantum/tech-samples/click-to-move-system-sample.png)
The gizmos shows the waypoints and the agent movement.


The scene view in the Unity Editor will draw gizmos of the target and the path.

### 2) Move With Callback

The `ClickMoveSystem` reads the player input and set the target for the entity. If the `TileAgentConfig` has the movement type set to `CUSTOM`, the entity will not move automatically, instead, it needs to listen to the `OnTileMapMoveAgent` callback and use the desiredDirection vector to move the position of the transform of the entity.

Check the `CallbackMoveSystem` to see other callback called when the system doesn’t find the path or when the agent reaches the path.

### 3) Edit Tile Map

Every grid cell is represented in the frame, as part of a bitset. This bitset informs whether the tile is traversable or not. The initial

values on the bitset are constructed during the game start based on the data baked in the scene, in edit time.

It is also possible to change the tile map during runtime, which is the objective of this sample scene. Click with the right mouse button to send a `Quantum Command` with the `TileType` (None or Wall) and the position. The `EditTileMapSystem` converts that position to an index and uses a bitset to store if the tile is traversable.

![Alternative Text](/docs/img/quantum/tech-samples/hexagon-grid-base-sample.png)

This entity uses a `CUSTOM` movement and evaluates each step if there is a wall in direction of the movement. Otherwise, it will repath the movement.

### 4) Decide the best target

This sample shows an entity chasing points on the map. It has a **Dijkstra's** algorithm that is used to decide which is the best item to chase. Different to the A\* algorithm, it can be used to measure the cost of movement of multiple entities in the map with only one search. The objective of the `AIDecisionSystem` is to provide NPC or Bots a way to evaluate which is the best decision based on the distance of interesting points in the map.

![Alternative Text](/docs/img/quantum/tech-samples/ai-decision-system-sample.png)
Entity moving into direction of the items.


The agent maps all available items and chooses the nearest to chase. In the Image, the blue entity chases alternately the yellow items in the map. The red point is the current target and when it finishes visiting all items, it repeats the search.

To identify which entities are valid to search, this sample uses the `ItemTileSystem` that maps the indexes of all entities with the `ItemTile` component. If an entity exists on the map, its position index can be found in the global `TileMapItems` dictionary.

## The Asset, Components and System

### Tilemap Asset

The `TileMapData` is an asset with a one-dimensional array of indexes. The world position can be converted to an index and an index can be converted to a world position. Different functions can be used to do this conversion but a grid-base tilemap normally uses the function:

![Alternative Text](/docs/img/quantum/tech-samples/convert-function.png)
The function to convert the X-Y axis into indexes.


Using indexes is cheaper and more flexible to use in the pathfinder algorithm. Look in the file TileMapData.cs the implementation of the methods IndexToPosition and PositionToIndex. These methods will be used every time in the project to handle inputs, that normally are given in X-Y coordinates.

C#

```csharp
// gets the index position of a world point position

var map = f.GetSingleton<RuntimeTileMap>();
var asset = map.Asset(f);

var worldPosition = frame.Get<Transform3D>(entity).Position;
var positionIndex = asset.PositionToIndex(worldPosition);

```

The `positionIndex` can be used to check if a specific tile is traversable or its neighboring tile is. The return indexes of GetNeighbors method can be converted to a world position.

C#

```csharp
var isPassable =  map.IsWall(positionIndex);
var neighbors = asset .GetNeighbors(positionIndex);

```

There are other fields in `TileMapData` and the asset can be extended if necessary. Remember this asset must contain only read-only data and any data that needs to be modified in runtime must be in a component. This project uses a singleton called `RuntimeTileMap` to store the information on which tiles are wall.

| Field | Description |
| --- | --- |
| Tiles | List of indexes of each tile in the map. |
| HEIGHT | Number of lines of Tilemap |
| WIDTH | Number of columns of Tilemap |
| TileWidth | Width of a tile based on world size. |
| TileHeight | Height of the tile based on world size. |
| OffsetX | Horizontal distance from the center of the level to the center of the first tile of the tilemap. |
| OffsetY | Vertical distance from the center of the level to the center of the first tile of the tilemap. |
| Neighbors | List of arrays of indexes for each tile in the map. Each index represents one neighbor and negative values represent invalid tiles. |

### TilePathfinder Component

`TilePathfinder` component has the essential data to perform the movement:

| Field | Description |
| --- | --- |
| Agent | The reference to TileAgentConfig with some parameters to limit and steer the movement. |
| Waypoints | List of tile indexes of tiles that sequentially represents the path to the target. |
| TargetPosition | The index of the target tile. |
| CurrentWaypoint | The index of the next tile to move towards. |

The first value of `Waypoints` normally is the target while the last value is the next point to move towards. To get the world position of the next waypoint use the `CurrentWaypoint` as an index of the waypoint list.

C#

```csharp
// gets the asset reference
var asset = f.GetSingleton<RuntimeTileMap>().Asset(f);
// gets the next point
var pathfinder = f.Get<TilePathfinder>(entity);
var index = pathfinder.Waypoints[pathfinder.CurrentWaypoint];
var nextPoint = asset.IndexToPosition(index);

```

The default value of `CurrentWaypoint` and `TargetPosition` is -1. This value means there is no path or target, respectively. Setting these values to -1 will cancel the movement of the agents.

### Agents

Agents are entities that can move on the map using `TilePathfinder` components and waypoints as references. The component needs a `TileAgentConfig` with the parameters for the movement. The table below shows the fields of the `TileAgentConfig`.

| Field | Description |
| --- | --- |
| Movement Type | The type of movement to adapt to gameplay style. |
| Velocity | Velocity of movement used in the TRANSFORM movement type. The character controller uses its own movement parameter. |
| Distance to Reach | The minimum distance from the agent to the waypoint to change to the next waypoint. |
| Max Number Of Waypoints | Limits the max number of waypoints in the component. |
| Max Number Of Waypoints | Limits the max number of waypoints in the component. |
| Max Cost Of Search | Limits the distance of search based on the cost of movement. |
| Draw Gizmos | Draw the waypoints and the target of that agent movement. |

The type of movement can be set in the `TileAgentConfig` asset file as an enum type. Change this option to select the movement type to be used.

![Alternative Text](/docs/img/quantum/tech-samples/agent-config-editor.png)
Set the movement type in the agent config asset.


The sample shows three possible types of movements:

- **Transform**: Set the position of the entity’s transform using the movement direction multiplied by the agent velocity.

- **Kinematic Character Controller(KCC)**: Only sets the movement direction in the kcc based on the next waypoint to chase.

- **Customized**: Sends signals with the entity reference and normalized direction of movement. The movement can be customized and applied to different systems according to user preferences.


### TilePathFinderSystem

The `TilePathFinderSystem` calculates the path and sets the waypoints in the `TilePathfinder` component. It drives the movement of the agent to the target using waypoints or sends signals to customized movement types. To start moving an agent, set a target position in the component. The SetTarget function will clear the `Waypoints` and try to find a new path to the target.

C#

```csharp
var pathfinder = frame.Unsafe.GetPointer<TilePathfinder>(entity);
pathfinder->SetTarget(frame, pathfinder, target);

```

Then the system performs the A\* algorithm calling the `AStar` method. Based on the tilemap and the limits of the agent’s movement, it may or may not find a path. If it is not found, a signal is sent and the indexes are invalidated by setting a negative value. Otherwise, a target is fixed and the `CurrentWaypoint` points to the last value in the `Waypoints` list.

The snippet below shows how `AStar` is used and how the parameters are applied. The method is static and can be used by other custom systems if necessary. As this function generates a lot of temporary data, the first parameter is a `PathfinderData` object that caches the reference of some lists to avoid GC allocation.

C#

```csharp
var result = AStar(pathfinderData, asset, frame, entityPosition, targetPosition, waypoints, maxCostOfSearch, maxNumberOfWaypoints);
if (result == PathFindStatus.NOT_FOUND) {
    waypoints.Clear();
    pathfinder->CurrentWaypoint = -1;
    pathfinder->TargetPosition = -1;
    f.Signals.OnTileMapSearchFailed(entity);
    return;
} else {
    pathfinder->TargetPosition = targetPosition;
    pathfinder->CurrentWaypoint = waypoints.Count - 1;
}

```

If no path is found or the max cost of the search is reached, the function will return the `NOT\_FOUND` enum. When it finds the target, the SetPath function is called to fill the waypoints with the index of the path and the `SUCCESS` enum is returned.

### Custom Movement Steering

The movement of the agent can be customized using callbacks when the movement type of `TileAgentConfig` is set to `CUSTOM`. To receive the callback, the system needs to implement the corresponding signal interface:

| Interface | Description |
| --- | --- |
| `ISignalOnTileMapMoveAgent` | Called in TilePathfinderSystem each frame while the agent has a valid target. It is called only when the agent is a type of CUSTOM. |
| `ISignalOnTileMapWaypointReached` | Called when the agent reaches a waypoint from the list. |
| `ISignalOnTileMapSearchFailed` | Called in the AStar function when the target is invalid or wasn’t found. |

The example below shows how to implement `ISignalOnTileMapMoveAgent` interface and steer the movement of the agent.

C#

```csharp
public class MoveSystem : SystemMainThread, ISignalOnTileMapMoveAgent{
    public void OnTileMapMoveAgent(Frame frame, EntityRef entity, FPVector3 direction) {
        if (f.Unsafe.TryGetPointer<Transform3D>(entity, out var transform)){
            transform->Position += direction.Normalized * frame.DeltaTime;
        }
    }
}

```

## Setting up a Level

To prototype, the `TileMapBaker` tool can be used to easily create a tilemap from a 3D level. This tool will check each tile position if there is some Unity’s Collider component. Then it will set in `TileMapData` the index of the tile and its type: `NONE` or `WALL`.

![Alternative Text](/docs/img/quantum/tech-samples/tile-map-baker-editor.png)### Creating a tilemap asset in Baker Tool

1. Create a new tilemap asset: Right button in DB folder > Create > Quantum > TileMapData > TileMap
2. In an empty game object, add the `TileMapBaker` component.

| Field | Description |
| --- | --- |
| Level | The root of the game object level. |
| Width | Number columns in the grid of tilemap. This is the value used in the conversion from the world position to the index position. |
| Height | Called in the AStar function when the target is invalid or wasn’t found.Number rows in the grid of the tilemap. |
| Bake Height Offset | As the level is a 3D scene, this offset is used to fix the height of colliders checking. |
| Tile Tolerance | The tolerance of tile size. |
| Layer | The specific layer that represents the walls or non-passable objects. |

3. Drag the game object root to `TileMapBaker`.
4. Click on the “Bake” button on the `TileMapBaker` component.

After baking, the traversable tiles will be visible in green and non-traversable tiles in red.

![Alternative Text](/docs/img/quantum/tech-samples/scene-tile-map-gizmos.png)### Initialize the map

In a default quantum scene, create an entity prototype and add the component `RuntimeTileMap`. This component is a singleton component that reads the asset and loads into a bitset the data about the type of tiles during initialization. This bitset structure is used to read and write the state of the tiles during runtime.

1. Drag the tilemap asset to this component.

![Alternative Text](/docs/img/quantum/tech-samples/runtime-map-editor.png)

2. When this scene starts, the `RuntimeTileMap` will start as a singleton component that loads the data baked in the asset to a bitset structure where each bit says if the tile is a wall or not.

## Extending and Adapting

The tile-based maps can have different types of layouts. To adapt the algorithm to different layouts, the `TileMapData` asset can be inherited. The `AStar` function tries to find the path by checking the neighbors of the tile, but there aren’t limits on how many neighbors a tile can have. Each tile is a node linked to other nodes, then is possible to create a customized `TileMapdata` with different properties changing the relationship between these nodes.

This sample exemplifies two different layouts: **Hexagonal** and **8-directions** map.

### Tilemap with 8 possible directions of movement

A tilemap with 8 directions can be easily created by overriding the neighbor's size and offset. The `GetNeighborsOffset` function is used to get the number of neighbors and their directions based on the index.

C#

```csharp
public class TileMap8 : TileMapData{
    public override int[] GetNeighborsOffset(int index) {
        return new int[8] {
            WIDTH, -WIDTH, 1, -1, WIDTH+1, WIDTH-1, -WIDTH-1, -WIDTH+1
        };
    }
}

```

This function is used by the baker tool to register the neighbors of each tile and check if they are valid. The value in the array is only the direction then the neighbor index can be found using the index + direction.

C#

```csharp
List<int> neighbors = new List<int>();
var neighborDirections = tileMap.GetNeighborsOffset(tile.Index);
// check if each neighbor is valid
foreach (var direction in neighborDirections ) {
    if (HasTileInDirection(tileMapAsset,tile.Index,direction)) {
        neighbors.Add(tile.Index + direction );
    } else {
        neighbors.Add(-1);
    }
}
// add the neighbors of that index to tilemap
tileMap.Neighbors[tile.Index] = new NeighborList() {
    Values = neighbors.ToArray()
};

```

### Hexagonal map

There are multiple approaches to deal with hexagonal tilemaps, but the simplest is to apply an offset in the rows or columns of a grid-based tilemap. Then the logic of neighbors must be changed.

In the example below, the HexagonOffset can be +1 or -1, depending on the desired direction of the offset. The column or row that will offset also changes if the top of the tile is pointy or flat.

C#

```csharp
public override FPVector3 IndexToPosition(int index) {
    int x = index % WIDTH;
    int y = index / HEIGHT;
    FP offsetX = 0;
    FP offsetY = 0;
    if (HexagonTop == HexagonTop.Pointy) {
        offsetX += (y % 2) * (TileWidth / 2) * (int)HexagonOffset;
    } else {
        offsetY += (x % 2) * (TileWidth / 2) * (int)HexagonOffset;
    }
    return new FPVector3((x * TileWidth) + offsetX, 0, (y * TileHeight) + offsetY) + Offset;
}

```

Look at the `HexagonalMap.cs` file to view other details of the hexagonal map implementation.

### Heuristics

The `AStar` function can use a heuristic to try to find the shortest path. The heuristic has a great impact on the performance of the algorithm and not using it is guaranteed to always find the shortest path, but with the worst performance. The distance value can be changed by inheriting from the `TileMapdata` and overriding `Heuristic` method.

C#

```csharp
// The Manhattan distance is better for grid base
public virtual ushort Heuristic(int from, int to) {
    FPVector3 a = IndexToPositionRaw(from);
    FPVector3 b = IndexToPositionRaw(to);
    return (ushort)(FPMath.Abs(a.X - b.X) + FPMath.Abs(a.Z - b.Z));
 }
// The linear distance is better fo hexagonal maps
public override ushort Heuristic(int from, int to) {
    FPVector3 a = IndexToPosition(from);
    FPVector3 b = IndexToPosition(to);
    return (ushort)(FPVector3.Distance(a, b));
}

```

## Highlights

### Heuristic Tradeoff

There is a tradeoff between speed or accuracy and once the project has only two samples of functions used on A\* algorithm, the final result can not be the enough to other projects. Without a heuristic function or if it returns 0, the A\* algorithm turns into a Dijkstra’s algorithm and the shortest path is guaranteed. But most of the games don’t really need the best path and if the value of the heuristic is the same as the cost of movement the algorithm will execute very fast, but finding the shortest path is not guaranteed.

### Priority Queue

The A\* algorithm uses a Priority Queue that sorts the tiles by the cost of movement. The project has multiple implementations of Priority Queue and a performance test shows the Binary Heap is the fastest data structure to handle this problem. Insertion Sort and other algorithms take at least the double time of the Binary Heap.

Check the `PriorityQueue` file to see the other's implementation. The implementation of the Priority Queue has an impact on the selected tiles and the final waypoint list. For example, if you change the comparison signal in the sort function, the movement will look more natural but more waypoints will be added.

![Alternative Text](/docs/img/quantum/tech-samples/priority-queue-using-equality.png)
Using `\_heap\[i\].Priority <= \_heap\[Parent(i)\].Priority` in the sort function.
![Alternative Text](/docs/img/quantum/tech-samples/priority-queue-whithout-equality.png)
Using `\_heap\[i\].Priority < \_heap\[Parent(i)\].Priority` in the sort function.
Back to top

- [Overview](#overview)
- [Technical Info](#technical-info)
- [Download](#download)
- [Sample Scenes](#sample-scenes)

  - [1) Simple Click Movement](#simple-click-movement)
  - [2) Move With Callback](#move-with-callback)
  - [3) Edit Tile Map](#edit-tile-map)
  - [4) Decide the best target](#decide-the-best-target)

- [The Asset, Components and System](#the-asset-components-and-system)

  - [Tilemap Asset](#tilemap-asset)
  - [TilePathfinder Component](#tilepathfinder-component)
  - [Agents](#agents)
  - [TilePathFinderSystem](#tilepathfindersystem)
  - [Custom Movement Steering](#custom-movement-steering)

- [Setting up a Level](#setting-up-a-level)

  - [Creating a tilemap asset in Baker Tool](#creating-a-tilemap-asset-in-baker-tool)
  - [Initialize the map](#initialize-the-map)

- [Extending and Adapting](#extending-and-adapting)

  - [Tilemap with 8 possible directions of movement](#tilemap-with-8-possible-directions-of-movement)
  - [Hexagonal map](#hexagonal-map)
  - [Heuristics](#heuristics)

- [Highlights](#highlights)
  - [Heuristic Tradeoff](#heuristic-tradeoff)
  - [Priority Queue](#priority-queue)