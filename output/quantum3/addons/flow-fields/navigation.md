# navigation

_Source: https://doc.photonengine.com/quantum/current/addons/flow-fields/navigation_

# Flow Fields Navigation

Path calculations are split in two main steps:

1. Calculate the path using portals data and A\*;
2. Calculate the smoothed path

![Navigation](/docs/img/quantum/v2/addons/flow-fields/navigation-1.png)
Red path shows the raw path, green path shows the smoothed path
## A\* Path

Every portal on the map is a node on the A\* graph. Possible paths between two portals are represented as edges in the A\* graph. When a path is requested, two extra nodes (start position and destination) are added to the graph with edges connected to portals which are reachable from the new nodes within the flow field controller. The A\* result is a sequence of portals which the pathfinder have to cross to reach its destination.

## Smooth Path

The smooth path is calculated per agent. The goal is to remove unnecessary corners on the path created by the A\* navigation and directions in the flow fields (only 8 possible directions).

Only the first corner of the smooth path is calculated in order improve the performance of the pathfinders which might change their destination before reaching it, meaning that the general path will be gradually smoothed, and it progresses as soon as a smooth corner is reached.

![Navigation](/docs/img/quantum/v2/addons/flow-fields/navigation-2.png)
Green dotted line represents the A\\\* result. 1. Start Position 2. Current smooth segment 3. Current smooth corner 4. Destination
## Path Caching

Similar paths (paths with the same start/end positions, controllers and same closest portal within the controller) are cached and reused by the pathfinders.

## Portal Fields

Flows towards each portal are precalculated and reused when needed. Recalculation happens every time a tile's cost changes in the controller or a portal's size/position changes (when cost changes in the neighboring controller on the border).

## Movement

Units movement is not part of this addon. The `FlowFieldPathfinder` provides the direction in which the unit should move to follow calculated path.

There are two movement implementations included in this sample:

1. Simple movement setting the position directly in the Transform2D component - see `MovementBasic`;
2. More advanced movement using the PhysicsBody2D component - see `MovementAdvanced`.

### MovementBasic Example

C#

```csharp
var pathfinder  = frame.GetPointer<FlowFieldPathfinder>(entity);
if (pathfinder->HasDestination == false || pathfinder->AtDestination == true)
    return;
var direction = pathfinder->GetDirection(frame, entity);
if (direction.Valid == false)
    return;
var transform        = frame.GetPointer<Transform2D>(entity);
transform->Position += direction.Direction * Speed * frame.DeltaTime;

```

### MovementAdvanced Example

C#

```csharp
var pathfinder  = frame.GetPointer<FlowFieldPathfinder>(entity);
var physicsBody = frame.GetPointer<PhysicsBody2D>(entity);
if (pathfinder->HasDestination == false || pathfinder->AtDestination == true)
{
    physicsBody->Velocity = default;
    return;
}
var direction = pathfinder->GetRotationDirection(frame, entity);
if (direction.Valid == false)
{
    physicsBody->Velocity = default;
    return;
}
physicsBody->Velocity = FPVector2.Rotate(FPVector2.Up, direction.Rotation) * Speed;
physicsBody->WakeUp();

```

## Change Units Between Maps

Is it possible to move entities between maps by manually change the `MapIndex` of the `FlowFieldPathfinder` component. This index is 0 by default but you can set it directly on inspector or manullay:

C#

```csharp
var pathfinder  = frame.GetPointer<FlowFieldPathfinder>(entity);
pathfinder->ChangeMap((Frame)frame, flowFieldIndex);

```

The maps doesn't have an index itself, this index is used to get the map from the array of maps in the `FrameContext`. So the order of the maps is the same as the order in the array that was set in the `Quantum Runner Local Debug` on Unit's inspector.

Back to top

- [A\* Path](#a-path)
- [Smooth Path](#smooth-path)
- [Path Caching](#path-caching)
- [Portal Fields](#portal-fields)
- [Movement](#movement)

  - [MovementBasic Example](#movementbasic-example)
  - [MovementAdvanced Example](#movementadvanced-example)

- [Change Units Between Maps](#change-units-between-maps)