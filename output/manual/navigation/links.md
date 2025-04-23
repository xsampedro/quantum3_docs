# links

_Source: https://doc.photonengine.com/quantum/current/manual/navigation/links_

# Using Navmesh Off Mesh Links

Unity `Nav Mesh Links` scripts are baked into the Quantum navmesh (`Quantum.NavMesh.Links`).

The following fields are used during baking:

- `Bidirectional` \- creates two links in each direction
- `CostModifier` \- used as input for `Quantum.NavMeshLink.CostOverride`
- `Area Type` \- set the region of the link (see Toggling a Navmesh Link at Run-Time section)
- `Activated` \- if not activated the link will be skipped during baking
- `GameObject.name` \- used as input for Quantum `Quantum.NavMeshLink.Name`

## Creating a Navmesh Link

- Create and configure a Unity `NavMesh Link`![Off Mesh Link Setup](https://doc.photonengine.com/docs/img/quantum/v3/manual/navigation/links-setup.png)
- Bake the Quantum navmesh and verify the link by checking the gizmos for blue arrows (Quantum gizmo overlay menu, enable `NavMeshLinks` and select the `QuantumMapData` GameObject) ![Nav Mesh Link Gizmo](https://doc.photonengine.com/docs/img/quantum/v3/manual/navigation/links-gizmo.png)
- Agents now considers the link for path-finding. ![Nav Mesh Link Path](https://doc.photonengine.com/docs/img/quantum/v3/manual/navigation/links-path.png)

## Toggling a Navmesh Link at Run-Time

Links can be toggled on and off during runtime using Quantum navmesh regions. They can be baked with a region id and toggled globally or per agent. Read the [Regions Manual](regions) for more information.

### Simple Region Mode

Choose the `Area Type` of the Unity `NavMesh Link` to set the region.

### Advanced Region Mode

- Set the `Area Type` of the Unity `NavMesh Link` to a region detection area.
- Attach a `MapNavMeshRegion` script to the Unity `NavMesh Link`, set a region `Id` and set `Cast Region` to `No Region`.

![Off Mesh Link Regions](https://doc.photonengine.com/docs/img/quantum/v3/manual/navigation/links-region.png)## Customization and Signals

The link data structure can be queried from the navmesh asset using `NavMesh.Links`. If a mapping from link `index` to `name` or vice-versa is required, a dictionary can be constructed.

The `NavMeshPathfinder` component has the following Link API:

- `bool IsOnLink(FrameBase)` \- returns `true` when the agent is currently on a link
- `int CurrentLink(FrameBase)` \- returns the current link index that points into `NavMesh.Links` or `-1` when currently not on a link

Waypoint have link related `WaypointFlags`:

- `LinkStart` \- this waypoint is the start of a link
- `LinkStop` \- this waypoint is the end of a link
- `RepathWhenReached` \- after reaching this waypoint the agent re-runs path finding

When setting a new target while the agent is on a link the agent will finish the current link before executing the path-finding (see `WaypointFlag.RepathWhenReached`).

No automatic re-pathing (for example for `NavMeshAgentConfig.MaxRepathTimeout`) will be executed as long as the agent is traversing a link.

By default the agent will traverse the link with its normal speed.

To take control of the agent movement when the link has been reached, the `ISignalOnNavMeshWaypointReached` signal can be used. Afterwards the agent can be disabled until an animation is complete or the movement can be overwritten in following `ISignalOnNavMeshMoveAgent` callbacks.

- Receiving any navigation callbacks requires enabled `SimulationConfig.Navigation.EnabledNavigationCallbacks`
- Receiving an `ISignalOnNavMeshMoveAgent` callback requires the `NavMeshAgentConfig.MovementType` to be set to `Callback`, it is possible to change the config of an agent at run-time.

The sample code performs a teleport when an agent steps on the link start waypoint.

C#

```csharp
namespace Quantum
{
  using Photon.Deterministic;
  using UnityEngine.Scripting;

  [Preserve]
  public unsafe class NewQuantumSystem : SystemMainThread, ISignalOnNavMeshWaypointReached
  {
    public override void Update(Frame frame)
    {
    }

    public void OnNavMeshWaypointReached(Frame frame, EntityRef entity, FPVector3 waypoint, Navigation.WaypointFlag waypointFlags, ref bool resetAgent)
    {
      var agent = frame.Get<NavMeshPathfinder>(entity);
      var waypointIndex = agent.WaypointIndex;
      if ((waypointFlags & Navigation.WaypointFlag.LinkStart) == Navigation.WaypointFlag.LinkStart)
      {
        // There always is another waypoint after the LinkStart
        var linkDestination = agent.GetWaypoint(frame, waypointIndex + 1);
        f.Unsafe.GetPointer<Transform2D>(entity)->Position = linkDestination.XZ;
      }
    }
  }
}

```

Back to top

- [Creating a Navmesh Link](#creating-a-navmesh-link)
- [Toggling a Navmesh Link at Run-Time](#toggling-a-navmesh-link-at-run-time)

  - [Simple Region Mode](#simple-region-mode)
  - [Advanced Region Mode](#advanced-region-mode)

- [Customization and Signals](#customization-and-signals)