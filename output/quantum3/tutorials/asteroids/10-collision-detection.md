# 10-collision-detection

_Source: https://doc.photonengine.com/quantum/current/tutorials/asteroids/10-collision-detection_

# 10 - Collision Detection

## Collision System

Quantum is ECS based, so collision events are not per entity based but instead global signals that systems can listen to. A common pattern for performance and convenience is to have a single system that receives all collision and then filters them by type and invokes other signals accordingly.

Before implementing the collision system, create a new ```
AsteroidsAsteroid.qtn
```

 and add the following code to it:

Qtn

```cs
component AsteroidsAsteroid
{
}

```

Add this empty tag component to the ```
AsteroidsLarge
```

prefab.

Next, create a new c# script and name it ```
AsteroidsCollisionSystem
```

.

Add the following code to it:

C#

```csharp
using UnityEngine.Scripting;

namespace Quantum.Asteroids
{
\[Preserve\]
public unsafe class AsteroidsCollisionsSystem : SystemSignalsOnly, ISignalOnCollisionEnter2D
{
public void OnCollisionEnter2D(Frame frame, CollisionInfo2D info)
{
// Projectile is colliding with something
if (frame.Unsafe.TryGetPointer<AsteroidsProjectile>(info.Entity, out var projectile))
{
if (frame.Unsafe.TryGetPointer<AsteroidsShip>(info.Other, out var ship))
{
// Projectile Hit Ship
}
else if (frame.Unsafe.TryGetPointer<AsteroidsAsteroid>(info.Other, out var asteroid))
{
// projectile Hit Asteroid
}
}

// Ship is colliding with something
else if (frame.Unsafe.TryGetPointer<AsteroidsShip>(info.Entity, out var ship))
{
if (frame.Unsafe.TryGetPointer<AsteroidsAsteroid>(info.Other, out var asteroid))
{
// Asteroid Hit Ship
}
}
}
}
}

```

This code listens to the global collision signal and then filters it by entity type.

Next create a ```
AsteroidsCollisionSignals.qtn
```

and add the following signals to it:

C#

```csharp
signal OnCollisionProjectileHitShip(CollisionInfo2D info, AsteroidsProjectile\* projectile, AsteroidsShip\* ship);

signal OnCollisionProjectileHitAsteroid(CollisionInfo2D info, AsteroidsProjectile\* projectile, AsteroidsAsteroid\* asteroid);

signal OnCollisionAsteroidHitShip(CollisionInfo2D info, AsteroidsShip\* ship, AsteroidsAsteroid\* asteroid);

```

## Projectiles colliding with Ships

The projectile based interactions will be implemented in the ```
AsteroidsProjectileSystem
```

. Open the system and add the following code:

C#

```csharp
public void OnCollisionProjectileHitShip(Frame frame, CollisionInfo2D info, AsteroidsProjectile\* projectile, AsteroidsShip\* ship)
{
if (projectile->Owner == info.Other)
{
info.IgnoreCollision = true;
return;
}

frame.Destroy(info.Entity);
}

```

The first part of the code ignores collisions when a projectile hits the own ship. The second part otherwise destroys the projectile. Additional logic could be added to destroy enemy ships on hit.

## Projectiles colliding with Asteroids

When Projectiles collide with asteroids, the asteroid should be split into multiple smaller asteroids.

There is already a function for spawning asteroids in the ```
AsteroidsWaveSpawnerSystem
```

that can be repurposed to also spawn child asteroids by turning it into a signal.

Add the following signal to the ```
AsteroidsAsteroid.qtn
```

. Also add a ```
ChildAsteroid
```

file to the component.

Qtn

```cs
component AsteroidsAsteroid
{
 asset\_ref<EntityPrototype> ChildAsteroid;
}

signal SpawnAsteroid(AssetRef<EntityPrototype> childPrototype, EntityRef parent);

```

The ```
ChildAsteroid
```

 field won't be modified at runtime, so technically it would be better fit in a config file. However, creating and linking a config file to a component that only contains a single field is overkill and only results in worse performance and increased code complexity.

Open the ```
AsteroidsWaveSpawnerSystem
```

and have it inherit from ```
ISignalSpawnAsteroid
```

.

Add ```
EntityRef parent
```

parameter to the ```
SpawnAsteroid
```

function so that implements the signal correctly. In the line that calls ```
SpawnAsteroid()
```

in the ```
SpawnAsteroidWave
```

function pass ```
EntityRef.None
```

as the parent entity into the function.

Finally, replace the following line

C#

```csharp
asteroidTransform->Position = GetRandomEdgePointOnCircle(f, config.AsteroidSpawnDistanceToCenter);

```

with:

C#

```csharp
if (parent == EntityRef.None)
{
 asteroidTransform->Position = GetRandomEdgePointOnCircle(f, config.AsteroidSpawnDistanceToCenter);
}
else
{
 asteroidTransform->Position = f.Get<Transform2D>(parent).Position;
}

```

Now, when a parent entity is provided the asteroid will be spawned at the position of the parent instead of a random position.

Return to the ```
AsteroidsProjectileSystem
```

 and add the following function to it:

C#

```csharp
public void OnCollisionProjectileHitAsteroid(Frame frame, CollisionInfo2D info, AsteroidsProjectile\* projectile, AsteroidsAsteroid\* asteroid)
{
if (asteroid->ChildAsteroid != null)
{
frame.Signals.SpawnAsteroid(asteroid->ChildAsteroid, info.Other);
frame.Signals.SpawnAsteroid(asteroid->ChildAsteroid, info.Other);
}

frame.Destroy(info.Entity);
frame.Destroy(info.Other);
}

```

Finally, for the two signals in the ```
AsteroidsProjectileSystem
```

to work it has to implement ```
ISignalOnCollisionProjectileHitShip
```

and ```
ISignalOnCollisionProjectileHitAsteroid
```

so add the two interfaces to it.

C#

```csharp
public unsafe class AsteroidsProjectileSystem : SystemMainThreadFilter<AsteroidsProjectileSystem.Filter>, ISignalAsteroidsShipShoot, ISignalOnCollisionProjectileHitShip, ISignalOnCollisionProjectileHitAsteroid

```

## Ships colliding with Asteroids

Open the ```
AsteroidsShipSystem
```

 and have it implement the ```
ISignalOnCollisionAsteroidHitShip
```

interface.

Add the corresponding function with the following code to it:

C#

```csharp
public void OnCollisionAsteroidHitShip(Frame frame, CollisionInfo2D info, AsteroidsShip\* ship, AsteroidsAsteroid\* asteroid)
{
 frame.Destroy(info.Entity);
}

```

For now, when a ship gets hit by an asteroid it is simply destroyed. In a full game loop there would be additional logic here to respawn the ship, adjust the score and end the game once all ships are destroyed.

Return to the ```
AsteroidsCollisionsSystem
```

 and hook the signals up by replacing

C#

```csharp
// Projectile Hit Ship

```

with

C#

```csharp
f.Signals.OnCollisionProjectileHitShip(info, projectile, ship);

```

and,

C#

```csharp
// Projectile Hit Asteroid

```

with

C#

```csharp
f.Signals.OnCollisionProjectileHitAsteroid(info, projectile, asteroid);

```

and finally,

C#

```csharp
// Asteroid Hit Ship

```

with

C#

```csharp
f.Signals.OnCollisionAsteroidHitShip(info, ship, asteroid);

```

Return to Unity and add the ```
AsteroidsCollisionSystem
```

to the ```
AsteroidsSystemConfig
```

.

## Child Asteroid Prefab

Drag the ```
AsteroidLarge
```

prefab into the scene. Adjust the ```
Circle Radius
```

to ```
0.5
```

in the ```
EntityPrototype
```

. Also scale the child model down to ```
(0.6, 0.6, 0.6)
```

. Rename it to ```
AsteroidSmall
```

. Drag the ```
AsteroidSmall
```

object back into the ```
Resources
```

folder and select ```
Prefab Variant
```

in the popup to create a smaller prefab variant for the asteroid. Then delete it from the scene.

Drag the ```
AsteroidSmall VariantEntityPrototype
```

 into the ```
Child Asteroid
```

field of the ```
AsteroidLarge
```

 prefab's ```
EntityPrototype
```

.

## Enabling Collision Callbacks

Code-wise collisions are now fully setup. However, by default Quantum does not invoke collision callbacks between entities for performance reasons. They can be enabled by modifying the ```
Callback Flags
```

 on the ```
PhysicsCollider2D
```

part of the ```
EntityPrototype
```

.

Enable the ```
OnDynamicCollisionEnter
```

callback flag on the ```
AsteroidLarge
```

, ```
AsteroidsProjectile
```

and ```
AsteroidsShip
```

prefab. (there is no need to adjust the ```
AsteroidSmall
```

as it is a prefab variant of ```
AsteroidLarge
```

and thus is adjusted automatically).

Enter play mode and test the multiple collision scenarios. The following should occur:

- When a projectile hits an asteroid the asteroid is destroyed and split into two asteroids. The projectile is destroyed as well.
- When the ship collides with an asteroid the ship is destroyed.
- When shooting another ship the projectile is destroyed.

## Fixing Wave Spawning

Currently, only the first wave of asteroids is being spawned. To fix this add the following function to the ```
AsteroidWaveSpawnerSystem
```

:

C#

```csharp
public void OnRemoved(Frame frame, EntityRef entity, AsteroidsAsteroid\* component)
{
 if (frame.ComponentCount<AsteroidsAsteroid>() <= 1)
 {
 SpawnAsteroidWave(frame);
 }
}

```

And have it implement the ```
ISignalOnComponentRemoved<AsteroidsAsteroid>
```

 interface. This signal gets called whenever a component is removed from an entity (when an entity is destroyed all components are removed from it as well).

With that the core gameplay mechanics of asteroids are implemented.

![Gameplay](/docs/img/quantum/v3/tutorials/asteroids/10-multiplayer-gameplay.gif)Back to top

- [Collision System](#collision-system)
- [Projectiles colliding with Ships](#projectiles-colliding-with-ships)
- [Projectiles colliding with Asteroids](#projectiles-colliding-with-asteroids)
- [Ships colliding with Asteroids](#ships-colliding-with-asteroids)
- [Child Asteroid Prefab](#child-asteroid-prefab)
- [Enabling Collision Callbacks](#enabling-collision-callbacks)
- [Fixing Wave Spawning](#fixing-wave-spawning)