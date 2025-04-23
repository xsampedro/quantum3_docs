# 6-asteroids

_Source: https://doc.photonengine.com/quantum/current/tutorials/asteroids/6-asteroids_

# 6 - Asteroids

## Overview

With the ship player entity ready the next step is to spawn in asteroids.

## Asteroid Spawner Config

Asteroids are spawned in waves each wave containing more asteroids. An asset object used as a config file is a better way to store the parameters that declare how asteroids are spawned instead of hard coding the variables related to asteroid spawning into a system. Create a new ```
AsteroidGameConfig
```

 script in the ```
QuantumUser > Simulation
```

folder. Add the following code to it:

C#

```
```csharp
using Photon.Deterministic;
using UnityEngine;

namespace Quantum.Asteroids
{
public class AsteroidsGameConfig: AssetObject
{
\[Header("Asteroids configuration")\]
\[Tooltip("Prototype reference to spawn asteroids")\]
public AssetRef<EntityPrototype> AsteroidPrototype;
\[Tooltip("Speed applied to the asteroid when spawned")\]
public FP AsteroidInitialSpeed = 8;
\[Tooltip("Minimum torque applied to the asteroid when spawned")\]
public FP AsteroidInitialTorqueMin = 7;
\[Tooltip("Maximum torque applied to the asteroid when spawned")\]
public FP AsteroidInitialTorqueMax = 20;
\[Tooltip("Distance to the center of the map. This value is the radius in a random circular location where the asteroid is spawned")\]
public FP AsteroidSpawnDistanceToCenter = 20;
\[Tooltip("Amount of asteroids spawned in level 1. In each level, the number os asteroids spawned is increased by one")\]
public int InitialAsteroidsCount = 5;
}
}

```

```

Note that the class inherits from ```
AssetObject
```

 this turns it into a ```
ScriptableObject
```

that can be injected into the Quantum simulation just like the ```
MapData
```

 and ```
SystemConfig
```

assets earlier in the tutorial.

```
AssetRef<EntityPrototype>
```

is a special type. AssetRef is Quantum's way of referencing ```
AssetObjects
```

 in the asset database. The other fields are simple config values that will be used by the system responsible for spawning asteroids.

To create an instance of the asset right-click on the ```
Resources
```

folder and select ```
Create > Quantum > Asset.. > AsteroidsGameConfig
```

. Name the new config asset ```
AsteroidsGameConfig
```

. Keep all the values as the default values and drop the ```
AsteroidLargeEntityPrototype
```

into the ```
AsteroidPrototype
```

field of the config.

![The AsteroidGameConfig](/docs/img/quantum/v3/tutorials/asteroids/6-game-config.png)## Injecting the Config

There are two common ways to get access to a ```
AssetObject
```

 inside the Quantum simulation

1. Have a ```
   AssetRef
   ```

    field on an entity component on a scene object or prefab.
2. Link the ```
   AssetObject
   ```

    as a global object via ```
   RuntimeConfig
   ```

   .

In this case the later approach is used since the game config is a global config that can be used by many systems.

Create a new script in the ```
QuantumUser > Simulation
```

folder and name it ```
RuntimeConfig.Asteroids
```

. Add the following code to it:

C#

```
```csharp
namespace Quantum
{
 public partial class RuntimeConfig
 {
 public AssetRef<Asteroids.AsteroidsGameConfig> GameConfig;
 }
}

```

```

Note that the ```
Quantum
```

instead of the ```
Quantum.Asteroids
```

namespace has to be used because this class is a partial class extending an existing ```
RuntimeConfig
```

class.

Return to Unity and head to the ```
QuantumDebugRunner
```

 GameObject in the scene. There is a new ```
GameConfig
```

field on the ```
RuntimeConfig
```

 entry now. Drop in the ```
AsteroidGameConfig
```

asset created earlier.

![Linking the Game Config](/docs/img/quantum/v3/tutorials/asteroids/6-link-game-config.png)## Asteroid Spawner System

Create a new c# script in the ```
QuantumUser > Simulation
```

 folder and name it ```
AsteroidsWaveSpawnerSystem
```

. Add the following code to it:

C#

```
```csharp
using Photon.Deterministic;
using UnityEngine.Scripting;

namespace Quantum.Asteroids
{
\[Preserve\]
public unsafe class AsteroidsWaveSpawnerSystem : SystemSignalsOnly
{
public void SpawnAsteroid(Frame frame, AssetRef<EntityPrototype> childPrototype)
{
AsteroidsGameConfig config = frame.FindAsset(frame.RuntimeConfig.GameConfig);
EntityRef asteroid = frame.Create(childPrototype);
Transform2D\* asteroidTransform = frame.Unsafe.GetPointer<Transform2D>(asteroid);

asteroidTransform->Position = GetRandomEdgePointOnCircle(frame, config.AsteroidSpawnDistanceToCenter);
asteroidTransform->Rotation = GetRandomRotation(frame);

if (frame.Unsafe.TryGetPointer<PhysicsBody2D>(asteroid, out var body))
{
body->Velocity = asteroidTransform->Up \* config.AsteroidInitialSpeed;
body->AddTorque(frame.RNG->Next(config.AsteroidInitialTorqueMin, config.AsteroidInitialTorqueMax));
}
}

public static FP GetRandomRotation(Frame frame)
{
return frame.RNG->Next(0, 360);
}

public static FPVector2 GetRandomEdgePointOnCircle(Frame frame, FP radius)
{
return FPVector2.Rotate(FPVector2.Up \* radius , frame.RNG->Next() \* FP.PiTimes2);
}
}
}

```

```

This system currently only contains a function that can be called to spawn an asteroid. The function gets the prototype of the asteroid entity to spawn by finding the ```
GameConfig
```

 asset that was linked to the ```
RuntimeConfig
```

earlier.

```
FindAsset
```

is very efficient and there is no problem with calling it repeatedly whenever an asteroid is spawned.

Besides creating the asteroid entity the function initializes the asteroid with a random position and rotation and applies velocity and torque to it based on the values in the ```
GameConfig
```

.

## Adding State to a System

With the ```
SpawnAsteroid
```

function ready it is time to implement waves. Asteroids should spawn in waves and each wave should spawn one more asteroid than the previous wave. For that it is necessary to keep track of a wave counter.

Systems in ECS are stateless so adding a regular variable to the ```
AsteroidsWaveSpawnerSystem
```

 is not allowed. This is not only done as a good practice in Quantum, it is necessary for the predict-rollback simulation to run correctly. **Never put state on a system in Quantum**.

Instead of putting state on a system the following two approaches can be used.

- Put state on the global frame. The global frame contains singleton variables that are accessible to all systems.
- Put state on a singleton component. Singleton components are like regular entity components but in addition they can be easily fetched via convenient API. You can learn more about singleton components [here](/quantum/current/manual/quantum-ecs/components#singleton-component).

For simplicity in this tutorial the global frame is used. Create a new ```
Global.qtn
```

file in the ```
QuantumUser > Simulation
```

folder. Add the following code to it:

```
```
global
{
 Int32 AsteroidsWaveCount;
}

```

```

This adds the int to the global frame which makes it accessible from any system by calling ```
frame.Global->AsteroidsWaveCount
```

.

## Adding Waves

With the wave counter implemented return to the ```
AsteroidWaveSpawnerSystem
```

 and add the following function:

C#

```
```csharp
private void SpawnAsteroidWave(Frame frame)
{
 AsteroidsGameConfig config = frame.FindAsset(frame.RuntimeConfig.GameConfig);
 for (int i = 0; i < frame.Global->AsteroidsWaveCount + config.InitialAsteroidsCount; i++)
 {
 SpawnAsteroid(frame, config.AsteroidPrototype);
 }

 frame.Global->AsteroidsWaveCount++;
}

```

```

This function spawns a number of asteroids based on how many waves have spawned and then increments the counter.

Finally, to get the first wave to spawn the ```
OnInit
```

function can be used which is called on every system once the simulation starts:

C#

```
```csharp
public override void OnInit(Frame frame)
{
SpawnAsteroidWave(frame);
}

```

```

With the implementation of the wave spawning system complete return to Unity. There is one last step needed to make the wave system run. Add it to the system list in the ```
AsteroidSystemConfig
```

.

## Identifying the Ship

When entering play mode now there are many null reference errors in the console from the ```
AsteroidsShipSystem
```

even though that system has not been changed at all. The reason for that is the systems filter which is used:

C#

```
```csharp
public struct Filter
{
public EntityRef Entity;
public Transform2D\* Transform;
public PhysicsBody2D\* Body;
}

```

```

Since both the ship and the asteroids have a transform and a physics body the filter finds both objects and runs the ship update function over every asteroid. A simple way to fix this is to add a unique component to the ship.

Create a new ```
AsteroidsShip
```

 .qtn file in the ```
QuantumUser > Simulation
```

folder. Add an empty component to it like this:

Qtn

```
```cs
component AsteroidsShip
{
}

```

```

**NOTE:** Empty components in ECS are referred to as tag components. As their use case is to identify entities just like with a tag system.

Return to Unity and open the ```
AsteroidsShip
```

 Prefab. On the ```
QuantumEntityPrototype
```

component press on the ```
+
```

 button on the ```
Entity Components
```

list and select the ```
AsteroidsShip
```

 component.

![Add AsteroidsShip Component](/docs/img/quantum/v3/tutorials/asteroids/6-add-asteroids-ship.png)

Finally, adjust the filter in the ```
AsteroidsShipSystem
```

to include the component:

C#

```
```csharp
public struct Filter
{
public EntityRef Entity;
public Transform2D\* Transform;
public PhysicsBody2D\* Body;
public AsteroidsShip\* AsteroidsShip;
}

```

```

Return to Unity and enter play mode. 5 asteroids are spawned in addition to the player object as part of the first wave.

![Asteroid Wave in Play Mode GIF](/docs/img/quantum/v3/tutorials/asteroids/6-asteroid-wave.gif)Back to top

- [Overview](#overview)
- [Asteroid Spawner Config](#asteroid-spawner-config)
- [Injecting the Config](#injecting-the-config)
- [Asteroid Spawner System](#asteroid-spawner-system)
- [Adding State to a System](#adding-state-to-a-system)
- [Adding Waves](#adding-waves)
- [Identifying the Ship](#identifying-the-ship)