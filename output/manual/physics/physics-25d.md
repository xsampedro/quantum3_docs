# physics-25d

_Source: https://doc.photonengine.com/quantum/current/manual/physics/physics-25d_

# 2.5D Physics

## Introduction

When using 2.5D Physics it is possible to add Height (or thickness) while still benefiting from most performance advantages available in the 2D physics engine.

**N.B.:** _Use Vertical Transform_ is a field (enabled by default) that can be toggled in the _SimulationConfig_ asset's Physics settings.

## 2.5D Physics with Vertical Data

_StaticCollider2D_ can have 'thickness' in the 3rd dimension using Quantum's 2.5D physics; simply set the _Height_:

![Adding Height to a Static Collider](https://doc.photonengine.com/docs/img/quantum/v3/manual/physics/physics-25d-static-collider-height.png)
Adding Height to a Static Collider.

For Entities, just add the _Transform2DVertical_ component and set its _Height_ and _Position Offset_. On a Quantum XZ-oriented game, this adds height on the Y axis, for example.

**N.B.: _Transform2DVertical_ requires the _Transform2D_ component.**

C#

```csharp
    var transform2dVertical = new Transform2DVertical();
    transform2dVertical.Height = FP._1;
    transform2dVertical.Position = FP._1;

    f.Set(entity, transform2dVertical);

```

![Adding Height to an Entity Prototype](https://doc.photonengine.com/docs/img/quantum/v3/manual/physics/physics-25d-entityprototype-transform2dvertical-height.png)
Adding Height to an Entity Prototype.

If entities or statics have a 3rd dimension, the physics engine will take into consideration when solving collisions. This allows for 'aerial' entities to fly over 'ground-based' ones, etc.

## Physics Engine Implications

### Entity Separation

**Important**: When a collision is detected, the collision solver does not use the extra dimension information. This can result in entity bounce when separation is performed on the basic 2D plane of the physics engine.

It is possible to simulate 3-dimensional gravity by manually applying speed and forces directly on _Transform2DVertical.Position_. The physics engine will use that information only for collision detection though.

### Raycast and Overlaps

These functions are by default all flat and only execute on the 2D plane. To take advantage of 2.5D, use the _overloaded_ version that takes the _height_ and vertical _offset_ as parameters.

Back to top

- [Introduction](#introduction)
- [2.5D Physics with Vertical Data](#d-physics-with-vertical-data)
- [Physics Engine Implications](#physics-engine-implications)
  - [Entity Separation](#entity-separation)
  - [Raycast and Overlaps](#raycast-and-overlaps)