# 9-shooting

_Source: https://doc.photonengine.com/quantum/current/tutorials/asteroids/9-shooting_

# 9 - Shooting

## ShipConfig

Currently, values like the acceleration and turn speed of the ship are hard coded into the `AsteroidsShipSystem`. A much better pattern is to put these variables plus additional variables that will be used to control the ship's shooting into a config file.

Start by creating a new c# script and name it `AsteroidsShipConfig`. Then add the following code:

C#

```csharp
using Photon.Deterministic;
using UnityEngine;

namespace Quantum
{
  public class AsteroidsShipConfig : AssetObject
  {
    [Tooltip(&#34;The speed that the ship turns with added as torque&#34;)]
    public FP ShipTurnSpeed = 8;

    [Tooltip(&#34;The speed that the ship accelerates using add force&#34;)]
    public FP ShipAceleration = 6;

    [Tooltip(&#34;Time interval between ship shots&#34;)]
    public FP FireInterval = FP._0_10;

    [Tooltip(&#34;Displacement of the projectile spawn position related to the ship position&#34;)]
    public FP ShotOffset = 1;

    [Tooltip(&#34;Prototype reference to spawn ship projectiles&#34;)]
    public AssetRef<EntityPrototype> ProjectilePrototype;
  }
}

```

Then add a reference to the ship config in the `AsteroidsShip.qtn`:

Qtn

```cs
component AsteroidsShip
{
    AssetRef<AsteroidsShipConfig> ShipConfig;
}

```

Create a new instance of the ship config in the Resources folder. Name it `DefaultShipConfig`. Then link it to the ship prototype by dropping it into the `ShipConfig` field on the `AsteroidsShip` GameObject.

Next adjust the code in the `UpdateShipMovement` function of the `AsteroidsShipSystem` to use the values from the config:

C#

```csharp
var config = f.FindAsset(filter.AsteroidsShip->ShipConfig);
FP shipAcceleration = config.ShipAceleration;
FP turnSpeed = config.ShipTurnSpeed;

```

## shooting

First, add a fire interval field to the `AsteroidsShip.qtn` to keep track of the time between firing intervals when shooting.

Qtn

```cs
component AsteroidsShip
{
    AssetRef<AsteroidsShipConfig> ShipConfig;
    FP FireInterval;
}

```

Next, open the `AsteroidsShipSystem` and add a function to handle firing:

C#

```csharp
private void UpdateShipFire(Frame frame, ref Filter filter, Input* input)
{
    var config = frame.FindAsset(filter.AsteroidsShip->ShipConfig);

    if (input->Fire && filter.AsteroidsShip->FireInterval <= 0)
    {
        filter.AsteroidsShip->FireInterval = config.FireInterval;
        // TODO create projectile
    }
    else
    {
        filter.AsteroidsShip->FireInterval -= frame.DeltaTime;
    }
}

```

This function allows the ship to fire at regular intervals when the fire button is pressed. The `TODO` line will be replaced with code that creates a bullet entity.

Call this function from the `Update` method of the system after calling `UpdateShipMovement`:

C#

```csharp
UpdateShipFire(f, ref filter, input);

```

Create a new c# script and name it `AsteroidsProjectileSystem`. Add the following code to it:

C#

```csharp
using Photon.Deterministic;
using UnityEngine.Scripting;

namespace Quantum.Asteroids
{
    [Preserve]
    public unsafe class AsteroidsProjectileSystem : SystemSignalsOnly
    {
        public void AsteroidsShipShoot(Frame frame, EntityRef owner, FPVector2 spawnPosition, AssetRef<EntityPrototype> projectilePrototype)
        {
            EntityRef projectileEntity = frame.Create(projectilePrototype);
            Transform2D* projectileTransform = frame.Unsafe.GetPointer<Transform2D>(projectileEntity);
            Transform2D* ownerTransform = frame.Unsafe.GetPointer<Transform2D>(owner);

            projectileTransform->Rotation = ownerTransform->Rotation;
            projectileTransform->Position = spawnPosition;
        }
    }
}

```

Add this system to the `AsteroidsSystemConfig` in Unity.

Next, create a `AsteroidsProjectileConfig` c# script to act as a config file for projectiles:

C#

```csharp
using Photon.Deterministic;
using UnityEngine;

namespace Quantum
{
    public class AsteroidsProjectileConfig : AssetObject
    {
        [Tooltip(&#34;Speed applied to the projectile when spawned&#34;)]
        public FP ProjectileInitialSpeed = 15;

        [Tooltip(&#34;Time until destroy the projectile&#34;)]
        public FP ProjectileTTL = 1;
    }
}

```

Finally, create a `AsteroidsProjectile.qtn` file to act as a data component for the projectile:

Qtn

```cs
component AsteroidsProjectile
{
    FP TTL;
    EntityRef Owner;
    AssetRef<AsteroidsProjectileConfig> ProjectileConfig;
}

```

Adjust the `AsteroidsShipShoot` function in the `AsteroidsProjectileSystem` by adding the following line to initialize the projectile to the end of it:

C#

```csharp
AsteroidsProjectile* projectile = f.Unsafe.GetPointer<AsteroidsProjectile>(projectileEntity);
var config = f.FindAsset(projectile->ProjectileConfig);
projectile->TTL = config.ProjectileTTL;
projectile->Owner = owner;

PhysicsBody2D* body = f.Unsafe.GetPointer<PhysicsBody2D>(projectileEntity);
body->Velocity = ownerTransform->Up * config.ProjectileInitialSpeed;

```

## Signals

`AsteroidsShipShoot` should be called whenever the ship is shooting. However, the function is in a different system. One solution to call it would be to make the function static and call it directly however this lead to tight coupling of systems which is not ideal especially for larger code bases.

A solution for that is to use signals. Signals are an event based way for systems to communicate to each other. Signals can be defined in any .qtn file. In this case open the `AsteroidsProjectile.qtn` and add a signal to it like this:

Qtn

```cs
component AsteroidsProjectile
{
    FP TTL;
    EntityRef Owner;
    AssetRef<AsteroidsProjectileConfig> ProjectileConfig;
}

signal AsteroidsShipShoot(EntityRef owner, FPVector2 spawnPosition, AssetRef<EntityPrototype> projectilePrototype);

```

Now, have the `AsteroidsProjectileSystem` implement the `ISignalAsteroidsShipShoot` interface.

C#

```csharp
public unsafe class AsteroidsProjectileSystem : SystemSignalsOnly, ISignalAsteroidsShipShoot

```

This interface causes the `AsteroidsShipShoot` function to be called whenever the signal is invoked.

To invoke the signal adjust the `AsteroidsShipSystem` by replacing the `// TODO create projectile` line with the following:

C#

```csharp
var relativeOffset = FPVector2.Up * config.ShotOffset;
var spawnPosition = filter.Transform->TransformPoint(relativeOffset);
f.Signals.AsteroidsShipShoot(filter.Entity, spawnPosition, config.ProjectilePrototype);

```

## Creating the Projectile Entity

Create a new `Quantum > 2D > Circle Entity` GameObject in the scene. Name it `AsteroidsProjectile`. Remove the `MeshRenderer` and `MeshFilter` from it. Adjust the radius of the circle collider to `0.16` and check the `PhysicsBody2D` checkbox.

Add a `AsteroidsProjectile` component in the `Entity Components` list of the `QuantumEntityPrototype`. The `TTL` and `Owner` fields are assigned at runtime via the config, so they can be kept empty. For the config field create a new `AsteroidsProjectileConfig` asset and name it `DefaultProjectileConfig`. Set the `Projectile Initial Speed` to `20` and the `Projectile TTL` to `2` then assign the config to the projectile in the config field of the component.

Create a `3D > Cube` GameObject as the child of the `AsteroidsProjectile` and scale it down to `(0.17, 0.17, 0.17)`. This will act as the visual of the projectile. Remove the box collider from the object.

Drag the `AsteroidsProjectile` into the `Resources` folder to make it a prefab then remove it from the scene. Next link the `prototype` of the prefab to the `Projectile Prototype` field of the `DefaultShipConfig`.

## Time to Live

Currently, projectiles do not get destroyed yet. Projectiles can be destroyed in two ways. Either once their time to live runs out or when colliding with an asteroid or spaceship.

To destroy projectiles after their TTL expires a system needs to check them each frame. Have the `AsteroidsProjectileSystem` inherit from `SystemMainThreadFilter<AsteroidsProjectileSystem.Filter>` instead of `SystemSignalsOnly` and add the following to it:

C#

```csharp
public struct Filter
{
    public EntityRef Entity;
    public AsteroidsProjectile* Projectile;
}

public override void Update(Frame frame, ref Filter filter)
{
    filter.Projectile->TTL -= frame.DeltaTime;
    if (filter.Projectile->TTL <= 0)
    {
        frame.Destroy(filter.Entity);
    }
}

```

Enter play mode. The spaceship can shoot now by pressing `space`. The projectiles automatically despawn after `2` seconds, however they currently collide with asteroids and the players as if they were regular physics bodies. In the next part of the tutorial collision detection will be added to the bullets.

Back to top

- [ShipConfig](#shipconfig)
- [shooting](#shooting)
- [Signals](#signals)
- [Creating the Projectile Entity](#creating-the-projectile-entity)
- [Time to Live](#time-to-live)