# ccd

_Source: https://doc.photonengine.com/quantum/current/manual/physics/ccd_

# Continuous Collision Detection (CCD)

## Overview

Continous Collision Detection is used to fast moving physics entities collide with other physics colliders instead of tunnelling through them.

There are two common approaches for CCD algorithms, _speculative_ and _sweep based_. Quantum implements a _speculative Continuous Collision Detection_ due to the performance considerations tied to its stateless physics engine. The speculative CCD approach is better suited for parallelism while also handling angular motion well; the former is needed for performance and the latter is necessary in many gameplay scenarios.

The speculative CCD algorithm increases the minimum bounding box used during the broad-phase for an entity based on its `PhysicsBody` component linear `Velocity` and `AngularVelocity`. It is called _speculative_ because it "speculates" the entity may collide with any of the other physics objects in that area and feeds all these candidates into the solver. This speculation ensures all contact contrains are taken into account when solving the collision thus preventing tunnelling.

## Set-up

Two simple steps are required to set up the CCD; both of which can done at edit-time and / or runtime.

**N.B.:** Given the performance impact CCD has on the simulation, the CCD functionality is enabled on a _per-entity_ basis and _NOT globally_!

### Edit-Time

Step 1: Check the `Allow CCD` boolean in the `Physics` section of the `Simulation Config` asset.

![Allow CCD in the Simulation Config](https://doc.photonengine.com/docs/img/quantum/v3/manual/physics/ccd-simulationconfig-toggle.png)
Enable the CCD in the Simulation Config.

Step 2: Enable the `Use Continuous Collision Detected` flag in the `Config` found on the `PhysicsBody` component of the Entity Prototype$.

![CCD Flag in the PhysicsBody Config](https://doc.photonengine.com/docs/img/quantum/v3/manual/physics/ccd-physicsbody-flag.png)
Select the CCD Flag in the PhysicsBody Config.
### Runtime

Should the CCD only be necessary in particular situation or moments of the game, it is possible to dynamically toggle the CCD and entities using it on and off.

Step 1: Toggle the `AllowCCD` property in the current game state's `PhysicsSceneSettings` which are part of the frame and initialized with the Physics values found in the `SimulationConfig` asset. **IMPORTANT:** Do _NOT_ modify the `SimulationConfig` asset at runtime, this is undeterministic and will result in desynchronization!

C#

```csharp
frame.PhysicsSceneSettings->CCDSettings.AllowCCD = true;

```

Step 2: Toggle the `UseContinuousCollisionDetection` property on the `PhysicsBody` component for the entity which should be using CCD.

C#

```csharp
var physicsBody = f.Unsafe.GetPointer<PhysicsBody3D>(myEntityRef);
physicsBody->UseContinuousCollisionDetection = true;

```

## Config

The `SimulationConfig` assets holds the default values for initializing the physics engine; including the aspects related to the CCD. The default values found in the `Continuous Collision Detection (CCD)` section are optimal for most games and should only be tweaked with care if edge cases were to arise.

- `AllowCCD`: Allows CCD to be performed if the Physics Body has CCD enabled on its Config flags.
- `CCDLinearVelocityThreshold`: If CCD is allowed, it will be performed on all Physics Bodies that have it enabled and have a linear velocity magnitude above this threshold.
- `CCDAngularVelocityThreshold`: If CCD is allowed, it will be performed on all Physics Bodies that have it enabled and have a angular velocity magnitude above this threshold.
- `CCDDistanceTolerance`: The absolute distance value below which the Physics Bodies under CCD check can be considered as touching.
- `MaxTimeOfImpactIterations`: The maximum number of iterations performed by the CCD algorithm when computing the time of impact between two Physics Bodies.
- `MaxRootFindingIterations`: The maximum number of iterations performed when computing the point in time when the distance between two Physics Bodies in a given separation axis is below the tolerance.

## Known Limitations

Although the _speculative CCD_ is feature complete, one needs to be aware of the know limitations of the speculative approach.

The current algorithm runs a single CCD iteration alongside the regular physics collision resolution. In other words, after a CCD collision is detected and resolved, the remaining delta-time for that entity is integrated regardless of CCD. Thus there is a chance of tunnelling occurring in highly dense environments with extremely fast moving entities.

Back to top

- [Overview](#overview)
- [Set-up](#set-up)

  - [Edit-Time](#edit-time)
  - [Runtime](#runtime)

- [Config](#config)
- [Known Limitations](#known-limitations)