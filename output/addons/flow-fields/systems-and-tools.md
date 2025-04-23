# systems-and-tools

_Source: https://doc.photonengine.com/quantum/current/addons/flow-fields/systems-and-tools_

# Systems and Tools

## Systems

The following systems are required to run FlowFields (listed in the same order as they are in `SystemSetup`)

- `FlowFieldPathfinderSystem\_Requests`: single thread system. Adds new destinations to `FlowFieldControllers` when requested by `FlowFieldPathfinder`.
- `FlowFieldMapSystem\_UpdateMap`: multi threaded system. Recalculates dirty controllers. A controller is dirty when the cost or portal has changed and when new destinations were added by `FlowFieldPathfinderSystem\_Requests`.
- `FlowFieldPathfinderSystem\_RequestPath`: multi threaded system. Calculates the identifier (used for path caching) for a requested path (combination of start position and destination) and registers the request to `FlowFieldMap`. When two pathfinders requests paths with the same identifier, in the same frame, the position and destination of pathfinder with the lowest entity index is used for the first step of the navigation (A\*).
- `FlowFieldMapSystem\_FindPaths`: multi threaded system. Calculates paths requested by `FlowFieldPathfinderSystem\_RequestPath` and adds them to the path cache.
- `FlowFieldPathfinderSystem\_CopyPaths`: single threaded system. Copies paths calculated by `FlowFieldMapSystem\_FindPaths` to `FlowFieldPathfinder`.
- `FlowFieldPathfinderSystem\_SmoothPath`: multi threaded system. Calculates the next smooth corner for each pathfinder who has new a destination or already arrived to its current smooth corner.
- `FlowFieldPathfinderSystem\_Removes`: single threaded system. Removes the destination requests from `FlowFieldControllers` added by `FlowFieldPathfinderSystem\_Requests`.
- `FlowFieldMapSystem\_ClearCache`: single threaded system. Clears cached paths when their count gets too high. The cache can't get too big because late join/reconnect snapshots might get too big.
- `FlowFieldPathfinderSystem`: signals only system. Initializes and deinitializes \`FlowFieldPathfinder\*.

## Tools

### FlowFieldMapDebug

Tool which shows the A\* graph used for navigation between controllers.

![Flow Field Map Debug](https://doc.photonengine.com/docs/img/quantum/v2/addons/flow-fields/ff-map-debug-1.png)
Green spheres represents portal positions and blue lines represents edges between portals.
### FlowFieldControllerDebug

Debug view of the `FlowFieldController`.

Can show the following data:

- **Cost Original** \- cost of the tiles **without** dynamic modifiers;
- **Cost** \- cost of the tiles **with** dynamic modifiers applied;
- **Integration** \- tiles integration (only contains data of the last calculated flow, one array is used for the entire Controller);
- **Flow** \- tiles directions towards the specified Destination or Portal;
- **Controller Location** \- tile location within the Controller;
- **Map Location** \- tile location within the Map;

![Flow Field Map Debug](https://doc.photonengine.com/docs/img/quantum/v2/addons/flow-fields/ff-controller-debug-1.png)
Cost field.
![Flow Field Controller Debug](https://doc.photonengine.com/docs/img/quantum/v2/addons/flow-fields/ff-controller-debug-2.png)
Flow towards portal.
![Flow Field Controller Debug](https://doc.photonengine.com/docs/img/quantum/v2/addons/flow-fields/ff-controller-debug-3.png)
Tool inspector.
Back to top

- [Systems](#systems)
- [Tools](#tools)
  - [FlowFieldMapDebug](#flowfieldmapdebug)
  - [FlowFieldControllerDebug](#flowfieldcontrollerdebug)