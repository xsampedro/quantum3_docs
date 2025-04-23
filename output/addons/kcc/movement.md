# movement

_Source: https://doc.photonengine.com/quantum/current/addons/kcc/movement_

# Movement

## Move step

Following equation is used to calculate desired position delta for single move step:

Step Position Delta = (Dynamic Velocity + Kinematic Velocity) \* Delta Time + External Delta

_Dynamic Velocity_

- Velocity that accumulates external forces
- Gravity, impulse from explosion, force field, jump pad
- **Can push the KCC upwards on slopes**
- Related API when using default `EnvironmentProcessor`:

  - `KCCData.DynamicVelocity`
  - `KCC.SetDynamicVelocity()`
  - `KCC.AddExternalForce()`
  - `KCC.AddExternalImpulse()`
  - `KCC.Jump()`

_Kinematic Velocity_

- Unconstrained velocity calculated from user input actions
- Usually based on `KCCData.InputDirection` (entry point)
- **Only walkable surfaces push the KCC upwards to prevent artifacts on steep slopes** (controlled by `KCCData.MaxGroundAngle`)
- Related API when using default `EnvironmentProcessor`:

  - `KCCData.KinematicVelocity`
  - `KCC.SetInputDirection()`
  - `KCC.SetKinematicDirection()`
  - `KCC.SetKinematicSpeed()`
  - `KCC.SetKinematicVelocity()`

_External Delta_

- Absolute position offset
- Useful for corrections on the end of move step (after depenetration)
- Related API:
  - `KCCData.ExternalDelta`
  - `KCC.SetExternalDelta()`

## Movement algorithm

The KCC virtual capsule always moves by calculated position delta and then depenetrates from overlapping colliders.

- ✅ This approach gives natural sliding against geometry.
- ✅ Most of the time single capsule overlap is needed which makes KCC very performant.
- ❌ Sometimes it results in jitter when moving against multiple colliders (usually some corners), but this can be usually resolved by running more depenetration steps.

## Continuous Collision Detection (CCD)

If the character moves too fast (position delta for single step is bigger than 75% or radius), the movement is divided into sub-steps to not pass through geometry.

- Movement is divided to smaller steps based on desired velocity.
- The maximum distance traveled in single step is 25-75% of radius and is controlled by `CCD Radius Multiplier` property in `KCC Settings` asset.
- 75% is good in most cases, lower this value only if you have problems with KCC running through geometry.
- Before lowering `CCD Radius Multiplier` try to increase `Max Penetration Steps` (improves quality of depenetration) in `KCC Settings` asset.

![Continuous Collision Detection](https://doc.photonengine.com/docs/img/quantum/v3/addons/kcc/ccd-feature.jpg)
Continuous collision detection.
## Collision filtering

- Primary filter is collision layer mask - controlled by `Collision Layer Mask` property in `KCC Settings` asset. The _Layer Collision Matrix_ from _Physics Settings_ is not used by KCC.
- It is possible to explicitly ignore a static collider or an entity with `KCC.SetIgnoreCollider()`.
- Use `KCCContext.PrepareUserContext()` to set `ResolveCollision` delegate for application of additional filtering rules.

Back to top

- [Move step](#move-step)
- [Movement algorithm](#movement-algorithm)
- [Continuous Collision Detection (CCD)](#continuous-collision-detection-ccd)
- [Collision filtering](#collision-filtering)