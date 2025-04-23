# prediction-culling

_Source: https://doc.photonengine.com/quantum/current/manual/prediction-culling_

# Prediction Culling

## Introduction

_Prediction Culling_ is used in games where players only have a partial view of the game world at any given time. It is safe to use and simple to activate.

_Prediction Culling_ allows to save CPU time during the Quantum prediction and rollback phases. Having it enabled will allow the predictor to run exclusively for important and visible entities to the local player(s), while leaving anything outside the view to be simulated only once per tick once the inputs are confirmed by the server; thus avoiding rollbacks wherever possible.

Although the performance benefits vary from game to game, they can be quite large; this is particularly important the more players the game supports, as the predictor will eventually miss at least for one of them.

Take for instance a game running at a 30Hz simulate rate. If the game requires an average of ten ticks rollback per confirmed input, this means the game simulation will have to be lightweight enough to run close to 300Hz (including rollbacks). Using _Prediction Culling_, full frames will be simulated at the expected 30/60Hz at all times, and the culling will be applied to the prediction area that is running within the prediction buffer.

## Setting Up Prediction Culling

As a consequence, using _Prediction Culling_ means the predicted simulation can never be accepted as the final result of a frame since part of it was culled, thus it never advanced the simulation of the whole game state.

To set up prediction culling, there are two steps; one in Quantum and one in Unity.

### In Quantum

By default, the Prediction Culling systems are already enabled in the sample Systems Config asset, located in ```
Assets/Photon/Quantum/Samples/SampleScenes/Resources
```

:

![Prediction Culling Systems](/docs/img/quantum/v3/manual/prediction-culling-systems.png)### In Unity

In Unity, it is necessary to set the prediction area. This will be used to decide which entities to cull from prediction.

Update the prediction area by calling ```
SetPredictionArea()
```

on every Unity update:

C#

```csharp
// center is either FPVector2 or FPVector3
// radius is an FP
QuantumRunner.Default.Game.SetPredictionArea(center, radius);

```

## What To Expect

### Physics And Navmesh Agents

The physics engines and the NavMesh related systems are affected by _Prediction Culling_.

When _Prediction Culling_ is enabled, they will only consider and update entities within the visible area on non-verified, i.e. predicted, frames.

CPU cycles are saved on account of physics and navmesh related agents skipping updates for any entity with the relevant component ( _PhysicsCollider_, _PhysicsBody_, _NavMeshPathFinder_, _NavMeshSteeringAgent_, _NavMeshAvoidanceAgent_) and outside the area of interest as defined by the Prediction Area center point and radius on the local machine.

### Iterators

The game codes can also benefit from _Prediction Culling_. Any filter that includes a ```
Transform2D
```

 or ```
Transform3D
```

will be subject to culling based on their positions.

Essentially whenever a prediction frame is running, calling any of the methods below will only return entities within the limits of the prediction radius, while the same call will return all active instances when simulating the verified frames after input confirmations arrive).

- ```
f.Filter()
```

- ```
f.Unsafe.FilterStruct()
```


**N.B.:** While **filters** benefit from _Prediction Culling_, **component iterators** _do NOT_.

- ```
f.GetComponentIterator()
```

- ```
f.Unsafe.GetComponentBlockIterator()
```


### Manual Culling Control Flags

It is also possible to manually **flag** entities for culling on predicted frames via the API provided via the Frame.

| Method | Description |
| --- | --- |
| SetCullable(EntityRef entityRef, bool cullable) | Sets if an entity can be culled or not. Does nothing if the entity does not exist (including invalid entity refs). |
| IsCulled(EntityRef entityRef) | If an entity is currently culled from the simulation, regardless of the frame state (Predicted or Verified).<br> <br>True if the entity is culled (for instance, not inside the prediction area) or does not exist.<br> <br>False otherwise (if the entity exists and is not being culled). |
| Culled(EntiyRef entityRef) | If an entity is prediction-culled.<br> <br>True if the frame is Predicted AND the entity IsCulled.<br> <br>False otherwise (if the frame is Verified or the entity is not culled). |
| Cull(EntiyRef entityRef) | Manually marks a cullable and existing entity as culled for this tick. | Does nothing if the entity does not exist or is not cullable. |
| ClearCulledState() | Resets the culling state of all entities on that frame. | Called automatically at the beginning of every frame simulation. |

To keep a consistent state and avoid desync, **de-flag** the culled entities on verified frames in the same systems you originally flag them. from the same system, so you keep a consistent state and do not desync.

## Avoiding RNG Issues

Using _RNGSession_ instances with _Prediction Culling_ is perfectly safe and determinism is guaranteed. However, their combined use can result in some visual jitter when two entities share a _RNGSession_, such as the default one stored in Quantum's ```
\_globals\_
```

. This is due to new a RNG value being generated for a _predicted_ entity after a _verifited_ frame was simulated, thus changing/modifying the entity's final position.

The solution is to store an isolated RNGSession struct in each entity subject to culling. The isolation guarantees culling will not affect the final positions of predicted entities unless the rollback actually required it.

chsarp

```chsarp
struct SomeStruct {
RNGSession MyRNG;
}

```

You can inject each RNGSession with their seeds in any way you desire.

Back to top

- [Introduction](#introduction)
- [Setting Up Prediction Culling](#setting-up-prediction-culling)

  - [In Quantum](#in-quantum)
  - [In Unity](#in-unity)

- [What To Expect](#what-to-expect)

  - [Physics And Navmesh Agents](#physics-and-navmesh-agents)
  - [Iterators](#iterators)
  - [Manual Culling Control Flags](#manual-culling-control-flags)

- [Avoiding RNG Issues](#avoiding-rng-issues)