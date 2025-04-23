# avoidance

_Source: https://doc.photonengine.com/quantum/current/manual/navigation/avoidance_

# Agent Avoidance

Quantum implements a variation of the collision avoidance technique called _Hybrid Reciprocal Velocity Obstacles (HRVO)_.

![Navmesh Agent Prototype](/docs/img/quantum/v3/manual/navigation/avoidance-banner.gif)## Setting up Avoidance Agents

The ```
NavMeshAvoidanceAgent
```

 component requires the entity to already have a ```
NavMeshPathfinder
```

**and** a ```
NavMeshSteeringAgent
```

 component.

### Simulation Config

The avoidance requires the following global parameters located on the ```
SimulationConfig
```

in the ```
Navigation
```

section.

![Simulation Config](/docs/img/quantum/v3/manual/navigation/avoidance-simulationconfig.png)

Disable ```
EnableAvoidance
```

to completely remove any overhead from the avoidance system.

```
AvoidanceRange
```

is crucial to the quality and performance cost of the avoidance system. It defines the range in which agents start to influence each other. The range is measured between the radii of two agents. (A) is the avoidance radius of an individual agent, and (B) is the avoidance radius of the agent plus the global ```
AvoidanceRange
```

.

```
inRange = Distance(positionAgentA, positionAgentB) - radiusAgentA - radiusAgentB < AvoidanceRange

```

![Simulation Config](/docs/img/quantum/v3/manual/navigation/avoidance-ranges.png)

```
MaxAvoidanceCandidates
```

defines the maximum number of avoidance candidates used by each agent. More candidates requires more memory and CPU but also increase the quality. The higher the ```
AvoidanceQuality
```

and the more agents that are influencing each other, the higher this number has to be.

```
VelocityObstacleTruncationFactor
```

defines how much a VO is truncated for non-moving obstacles.

### NavMeshAgentConfig - Avoidance

|     |     |
| --- | --- |
| **AvoidanceType** | Sets the active avoidance mode of the agent.<br> <br>```<br>None<br>```<br> = the agent will not avoid others but others will avoid it<br>```<br>Internal<br>```<br> = the agent will actively avoid other agents using the internal avoidance systems |
| **AvoidanceQuality** | Sets the active avoidance mode of the agent.<br> <br>```<br>None<br>```<br> = the agent will not avoid others but others will avoid it<br>```<br>Internal<br>```<br> = the agent will actively avoid other agents using the internal avoidance systems |
| **Priority** | The agent ```<br>Priority<br>```<br> works as in Unity. <br> <br> Default = ```<br>50<br>```<br> Most important = ```<br>0<br>```<br> Least important = ```<br>99<br>```<br> Because the avoidance system relies on reciprocity, the avoiding-work (who will avoid whom, and by how much) is always split between the agents. Higher priority agents do only ```<br>25%<br>```<br> of the work while agents of the the same priority split the work ```<br>50<br>```<br>/```<br>50<br>```<br>. |
| **AvoidanceRadius** | The avoidance radius of this agent should roughly match the visual character size.<br> <br> Together with the avoidance radius from the ```<br>SimulationConfig<br>```<br> it is used during the avoidance broadphase that creates agent pairs that influence each other. |
| **AvoidanceLayer** | Sets the avoidance layer of this agent. <br> <br> The Unity layers are used for ```<br>AvoidanceLayer<br>```<br> and ```<br>AvoidanceMask<br>```<br> to filter agents. |
| **AvoidanceMask** | Sets the avoidance mask of this agent. Agents that have an ```<br>AvoidanceLayer<br>```<br> that is not contained in the mask will be ignored. |
| **MaxAvoidanceCandidates** | Sets the maximum number of avoidance candidates for this agent type. The global max is set in the ```<br>SimulationConfig<br>```<br>. |
| **ReduceAvoidanceAtWaypoints** | Solving avoidance while also trying to follow waypoints to steer around corners or through narrow passages is hard. To mitigate the problem and to accept visual overlapping in favor of agents blocking each other toggle ```<br>ReduceAvoidanceAtWaypoints<br>```<br>. |
| **ReduceAvoidanceFactor** | The avoidance applied when an agent is getting close to a waypoint is reduced. The ```<br>ReduceAvoidanceFactor<br>```<br> value is multiplied with the agent radius and then represents the distance in which the avoidance influence is reduced quadratically. |
| **AvoidanceCanReduceSpeed** | This option allows velocity candidates to decrease the agent's speed, making avoidance maneuvers look more natural. |
| **ShowDebugAvoidance** | Defines if velocity obstacles and candidates are drawn as gizmos during run-time. |

## Setting up Avoidance Obstacles

```
Avoidance Obstacles
```

are static or moving entities that influence the avoidance behavior of Navmesh agents but are not agents themselves. They do not influence the path finder and should **not** be used to block parts of the game level.

A ```
NavMeshAvoidanceObstacle
```

 component requires a ```
Transform2D
```

or ```
Transform3D
```

 component to work properly.

If the entity that has a ```
NavMeshAvoidanceObstacle
```

component is moving, other agents require its velocity information to predict its future position and the ```
NavMeshAvoidanceObstacle.Velocity
```

has to be set manually.

Add a ```
NavMeshAvoidanceObstacle
```

to a Quantum Entity Prototype in Unity.

![Avoidance Obstacle Prototype](/docs/img/quantum/v3/manual/navigation/avoidance-obstacle.png)

Or add it in code.

C#

```csharp
var c = f.Create();
f.Set(c, new Transform2D { Position = new FPVector2(8,-2) });
var obstacle = new NavMeshAvoidanceObstacle();
obstacle.AvoidanceLayer = 0;
obstacle.Radius = FP.\_0\_50;
obstacle.Velocity = FPVector2.Zero;
f.Set(c, obstacle);

```

## Jittering Agents

The agent movement, particularly their movement direction, can be unstable due to the nature of the avoidance calculations in conjunction with fixed-point math. This causes the agents to jitter visibly.

To mitigate this, the ```
Angular Speed
```

 of the agents can be tuned down, or view smoothing can be added. This can be done by overriding the ```
QuantumEntityView
```

class and adding blending to the rotation when applying the transforms.

C#

```csharp
namespace Quantum {
 using UnityEngine;

 public class SmoothRotationEntityView : QuantumEntityView {
 public float Blending = 15;

 private Quaternion rotation;

 protected override void ApplyTransform(ref UpdatePositionParameter param) {
 // Override this in subclass to change how the new position is applied to the transform.
 transform.position = param.NewPosition + param.ErrorVisualVector;

 // Unity's quaternion multiplication is equivalent to applying rhs then lhs (despite their doc saying the opposite)
 rotation = param.ErrorVisualQuaternion \* param.NewRotation;
 transform.rotation = Quaternion.Lerp(transform.rotation, rotation, Time.deltaTime \* Blending);
 }
 }
}

```

Back to top

- [Setting up Avoidance Agents](#setting-up-avoidance-agents)

  - [Simulation Config](#simulation-config)
  - [NavMeshAgentConfig - Avoidance](#navmeshagentconfig-avoidance)

- [Setting up Avoidance Obstacles](#setting-up-avoidance-obstacles)
- [Jittering Agents](#jittering-agents)