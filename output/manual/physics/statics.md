# statics

_Source: https://doc.photonengine.com/quantum/current/manual/physics/statics_

# Static Colliders

## Introduction

Adding static colliders to a Scene takes three simple steps:

1. Attach a Quantum Static Collider Script to a Unity GameObject;
2. Edit the properties to resemble the geometry you want for the static obstacle in the scene; and,
3. Bake the Scene via the MapData Script.

![Step 1 & 2 - Add Static Colliders to GameObject in Unity Scene and adjust Settings](https://doc.photonengine.com/docs/img/quantum/v2/manual/physics/physics-static-colliders-2d.png)
Step 1 & 2 - Add Static Colliders to GameObject in Unity Scene and adjust Settings.
![Step 3 - Baking the Map Saves the Scene Colliders as a Quantum Asset (Map)](https://doc.photonengine.com/docs/img/quantum/v2/manual/physics/physics-static-colliders-map-baking.png)
Step 3 - Baking the Map Saves the Scene Colliders as a Quantum Asset (Map).
## Unity Collider as a Source

The Quantum static collider can also mirror the properties from a Unity collider.

To do that, simply drag and drop the desired collider into the `Source Collider` field on the Quantum Static Collider component:

![unity-collider-source](https://doc.photonengine.com/docs/img/quantum/v2/manual/physics/physics-static-collider-source.png)## Shapes

The 2D Physics shapes are:

- Circle
- Box
- Polygon

**N.B.:** All of them have a `Height` field, which allows the creation of 2.5D shapes.

The 3D Physics shapes are:

- Sphere
- Box
- Mesh

## Configuration

Static colliders can be fitted with a _PhysicsMaterial_ and a _User Asset_. The latter is available in the simulation via the collision callbacks.

### Smooth Sphere Mesh Collider

A `Static Mesh Collider 3D` has an option called `Smooth Sphere Mesh Collision`. When toggling this option on, the physics solver will resolve sphere-mesh collisions as if the mesh was a regular flat and smooth plane. This prevents adding spin to a sphere colliding with triangle edges.

![Static Smooth Mesh Collider](https://doc.photonengine.com/docs/img/quantum/v2/manual/physics/physics-static-colliders-smooth-mesh-collision.png)
Static Mesh Collider.

If the \`Static Mesh Collider 3D\` is marked with the \`Smooth Sphere Mesh Collision\` option but the mesh is not completely flat, it might result in undesirable collision responses.

## Enable / Disable at Runtime

This section will present several approaches to enable and disable static colliders at runtime in simulation.

### Physics Engine

It is possible to toggle static colliders on and off at runtime directly in the Physics Engine.

For a static collider to be toggle-able, its `Mode` needs to be set at _edit-time_ (Unity) and baked into the `Map` asset. The mode can be set to:

- `Immutable` (default): the collider cannot be enabled or disabled at runtime.
- `Toggleable Start On`: the collider can be toggled at runtime and starts **enabled**.
- `Toggleable Start Off`: the collider can be toggled at runtime and starts **disabled**.

![Enable toggle on 3D Static Mesh Colliders](https://doc.photonengine.com/docs/img/quantum/v2/manual/physics/physics-static-colliders-mesh-toggle.png)
Enable toggle on a 3D Static Mesh Collider component.

Once a static collider has been marked as toggle-able and baked, it becomes possible to enable and disable the collider at runtime from the simulation (Quantum) using `SetStaticColliderEnabled()` in `Frame.Physics3D` and `Frame.Physics2D` for 3D and 2D static colliders respectively.

The `index` to be passed as a parameter is the collider's index in the `frame.Map.StaticColliders` array. Collision callbacks return the index (`ColliderIndex`) of a static collider as part of the `StaticData` in their `TriggerInfo` or `CollisionInfo`.

**IMPORTANT:** A disabled static mesh collider is ignored by physics queries and will not trigger collision signals.

### Manual Tracking

Although static colliders can be enabled / disabled at the physics engine level, there are various approaches to do the same manually.

#### Keep a global bitset for the state

If the only purpose is to keep track of which static colliders to ignore or take into account in a collision callback, the most convenient approach is to define a global `BitSet` which is of the same length or bigger than the `frame.Map.StaticColliders` array. This can be done as part of the `Frame` object or as a singleton component.

C#

```csharp
singleton component StaticColliderState {
    bitset[256] colliders;
}

```

This allows to use the bitset instance with the collider indices to set its bits.

C#

```csharp
// loops through the bitset to initialize all bits as &#34;On&#34; to mark all colliders as active
public override void OnInit(Frame frame)
{
    var collidersState = frame.Unsafe.GetPointerSingleton<StaticColliderState>();
    for (int i = 0; i < collidersState->colliders.Length; i++) {
        collidersState->colliders.Set(i);
    }
}

public void OnTrigger3D(Frame frame, TriggerInfo3D info)
{
    if (info.IsStatic == false) return;

    // Use a custom asset slotted in the UserAsset field to identify toggleable colliders
    var colliderAsset = frame.FindAsset<MyColliderAsset>(info.StaticData.Asset);
    if (colliderAsset == null) return;

    var collidersState = frame.Unsafe.GetPointerSingleton<StaticColliderState>();
    collidersState->colliders.Clear(info.StaticData.ColliderIndex);
}

```

The values can then be read using `IsSet()` and used to check whether a collision signal should be handled or ignored. This is particularly useful when dealing with static interactable objects, environmental barriers or implementing `IKCCCallbacks3D` for movement.

#### Toggle with Behaviour

Static colliders are assets, i.e. they are stateless and immutable at runtime. However, there are instance where static objects should be enabled / disabled based on dynamic conditions.

For example, pick-ups can usually be represented with a static position and a trigger collider; turning those into static colliders will avoid the over associated with dynamic entities. Unfortunately, the timer commonly to re-spawn a power-up after its pick-up cooldown requires a state. It is possible to solve this conundrum by extending the concept presented in the previous section.

First, the state of the static colliders representing power-ups needs to be held somewhere.

C#

```csharp
singleton component PowerUps {
    [ExcludeFromPrototype] bitset[256] IsPowerUp;
    [ExcludeFromPrototype] bitset[256] State;
    [ExcludeFromPrototype] array<FP>[256] Timers;
    FP SpawnCooldown;
}

```

Then a system can be created to handled the enabling and disabling of the power-ups.

C#

```csharp
public unsafe class MyPowerUpSystem : SystemMainThread {
public override void OnInit(Frame frame)
{
    var powerUps = frame.Unsafe.GetPointerSingleton<PowerUps>();
    for (int i = 0; i < powerUps->IsPowerUp.Length; i++)
    {
        var powerUp = frame.FindAsset<MyPowerUpAsset>(frame.Map.StaticColliders3D[i].StaticData.Asset);
        if (powerUp == null) {
            powerUps->IsPowerUp.Clear(i);
            continue;
        }

        powerUps->IsPowerUp.Set(i);
        powerUps->State.Set(i);
        powerUps->Timers[i] = FP._0;
    }
}

public override void Update(Frame frame)
{
    var powerUps = frame.Unsafe.GetPointerSingleton<PowerUps>();
    for (int i = 0; i < powerUps->IsPowerUp.Length; i++)
    {
        if (powerUps->IsPowerUp.IsSet(i) == false) continue;
        if (powerUps->State.IsSet(i)) continue;

        powerUps->Timers[i] -= frame.DeltaTime;
        if(powerUps->Timers[i] > 0) continue;

        powerUps->State.Set(i);
        // Other code visualizing the spawned / re-enabled power-up
        // can use frame event to trigger VFX, SFX, re-enable visual / GameObject
    }

}

public void OnTrigger3D(Frame frame, TriggerInfo3D info)
{
    if(info.IsStatic == false) return;

    var powerUps = f.Unsafe.GetPointerSingleton<PowerUps>();

    if(powerUps->IsPowerUp.IsSet(info.StaticData.ColliderIndex) == false) return;
    if(powerUps->State.IsSet(info.StaticData.ColliderIndex) == false) return;

    powerUps->State.Clear(info.StaticData.ColliderIndex);
    powerUps->Timers[info.StaticData.ColliderIndex] = powerUps->SpawnCooldown;

    // Remember to communicate the disabled state visually, e.g. trigger a frame event to disable the GameObject in Unity
}

```

Back to top

- [Introduction](#introduction)
- [Unity Collider as a Source](#unity-collider-as-a-source)
- [Shapes](#shapes)
- [Configuration](#configuration)

  - [Smooth Sphere Mesh Collider](#smooth-sphere-mesh-collider)

- [Enable / Disable at Runtime](#enable-disable-at-runtime)
  - [Physics Engine](#physics-engine)
  - [Manual Tracking](#manual-tracking)