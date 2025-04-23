# material

_Source: https://doc.photonengine.com/quantum/current/manual/physics/material_

# Materials

## Overview

Every _PhysicsBody_ requires a _PhysicsMaterial_ (a quantum data-asset). The PhysicsMaterial holds properties necessary for the physics engine to resolve collisions, integration of forces and velocities.

## PhysicsMaterial Data-Asset

The PhysicsMaterial holds the parameters for:

- Restitution (sometimes referred to as "bounciness", or "bounce")
- Restiution Combine Function
- Friction Static
- Friction Dynamic
- Friction Combine Function

If no _PhysicsMaterial_ asset is slotted, the default physics material will be assigned; the default physics material is the one linked in the _SimulationConfig_ physics settings.

![Adjusting Properties to Physics Materials](/docs/img/quantum/v3/manual/physics/physics-material-asset.png)
Adjusting Properties to Physics Materials.


A _PhysicsMaterial_ asset can be assigned to a _PhysicsCollider_ directly:

C#

```
```csharp
var material = f.FindAsset<PhysicsMaterial>("steel");
collider.Material = material;

f.Set(entity, collider);

```

```

### Important Note

A _PhysicsMaterial_ is a data asset and lives in the Quantum Asset Database. As assets are not part of the rollback-able game state, every _PhysicsMaterial_ is therefore to be considered immutable at runtime. Changing its properties while the game running leads to non-deterministic behaviour.

_PhysicsMaterial_ s follow the same rules as other data-assets.

C#

```
```csharp
// this is NOT safe and cannot be rolled-back:
collider->Material.Restitution = FP.\_0;

// switching a reference is safe and can be rolled back:
var newMaterial = f.FindAsset<PhysicsMaterial>("ice");
collider->Material = newMaterial;

```

```

## Combine Functions

The Combine Function used to resolve the restitution and friction for each collision manifold (a collision pair) is based on the combine functions' precedence order. The Physics system will chose the function with the highest precedent from the two colliders.

The precedence order is:

1. Max
2. Min
3. Average
4. Multiply

For instance: take a collision manifold with a Collider A and Collider B. Collider A's physics material has a _Restitution Combine Function_ set to **Max**, while Collider B's physics material has its set to **Average**. Since _Max_ has a higher priority than _Average_, the restitution for this collision will be solved using the _Max_ function.

The same logic applies to the _Friction Combine Function_.

**N.B.:** The _Friction Combine Function_ and _Restitution Combine Function_ are resolved separately and thus carry different settings.

Back to top

- [Overview](#overview)
- [PhysicsMaterial Data-Asset](#physicsmaterial-data-asset)

  - [Important Note](#important-note)

- [Combine Functions](#combine-functions)