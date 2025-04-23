# 5-player-spawning

_Source: https://doc.photonengine.com/quantum/current/tutorials/asteroids/5-player-spawning_

# 5 - Player Spawning

## Overview

With the player character entity created the next steps are spawning a player character for each player that joins the game and linking up the input of the player to the entity.

## Player Prefab

Currently, the ship is a scene object. Quantum can spawn entities at runtime similar to how prefabs can be spawned at runtime in a single player Unity game. Turn the ship into a prefab by dragging the ```
AsteroidsShip
```

 GameObject from the scene into the ```
Resources
```

folder. After doing so, delete the ```
AsteroidsShip
```

 from the scene.

![The AsteroidsShip Prefab](/docs/img/quantum/v3/tutorials/asteroids/5-player-prefab.png)

Note how a ```
EntityPrototype
```

and a ```
EntityView
```

file have been created nested into the prefab. These are used by Quantum to spawn the entity and link it up to the Unity view.

## Player Link Component

Quantum has the concept of players. Each client can own one or multiple players. However, Quantum does not have a built-in concept of a player object/avatar. Each player that is connected to the game is given a unique ID. This ID is called a ```
PlayerRef
```

. To link an entity to a specific player we will create a PlayerLink component that contains the ```
PlayerRef
```

of the owner.

Create a ```
PlayerLink.qtn
```

file in the ```
Assets/QuantumUser/Simulation
```

folder and add the following code to it:

Qtn

```cs
component PlayerLink
{
player\_ref PlayerRef;
}

```

## Player Data

To dynamically spawn a character we need to let the gameplay code know which entity it should create. Quantum has the concept of player data. Player data allows each player to pass information into the simulation upon connection. It can be information such as which character a player has selected or which skin they are using.

By default the player data contains an ```
Avatar
```

entity and a ```
Player Nickname
```

, however the data can also be extended in the ```
Photon/QuantumUser/Simulation/RuntimePlayer.User.cs
```

file.

To spawn the player entity we will use the predefined avatar field. An ```
AssetRefEntityPrototype
```

 is the Quantum equivalent to a prefab.

When entering play mode the Quantum simulation automatically runs. This is driven by the ```
QuantumRunnerLocalDebug
```

component on the ```
QuantumDebugRunner
```

GameObject in the ```
QuantumGameScene
```

scene in Unity. This component is used to run a single player version of the ```
Game
```

locally for development purposes.

The ```
QuantumLocalRunnerDebug
```

allows to simulate any numbers of local players. In the ```
QuantumLocalRunnerDebug
```

under ```
Local Players
```

drag and drop the ```
AsteroidsShipEntityPrototype
```

file that can be found under the ```
AsteroidsShip
```

prefab into the ```
Player Avatar
```

field of the player.

![Inspector view of the QuantumDebugRunner](/docs/img/quantum/v3/tutorials/asteroids/5-runner-link-playerobject.png)

Now that the prefab is linked up to the player data all that is left is to write code to spawn the entity when a player joins.

## Spawning Player Objects

Create a new ```
ShipSpawnSystem.cs
```

class. Add the following code:

C#

```csharp
using UnityEngine.Scripting;

namespace Quantum.Asteroids
{
 \[Preserve\]
 public unsafe class ShipSpawnSystem : SystemSignalsOnly, ISignalOnPlayerAdded
 {
 public void OnPlayerAdded(Frame frame, PlayerRef player, bool firstTime)
 {
 {
 RuntimePlayer data = frame.GetPlayerData(player);

 // resolve the reference to the avatar prototype.
 var entityPrototypAsset = frame.FindAsset<EntityPrototype>(data.PlayerAvatar);

 // Create a new entity for the player based on the prototype.
 var shipEntity = frame.Create(entityPrototypAsset);

 // Create a PlayerLink component. Initialize it with the player. Add the component to the player entity.
 frame.Add(shipEntity, new PlayerLink { PlayerRef = player });
 }
 }
 }
}

```

This code creates the ship entity when a player joins and links it up to the player by adding a PlayerLink component to it.

Signals are similar to events in C#. They are used by Quantum systems to communicate with each other. Quantum comes with a lot of existing signals such as the ```
ISignalOnPlayerAdded
```

 which gets called after a player has joined the session and shared their player data.

SystemSignalsOnly is a special type of system that doesn't have an Update routine, which makes it leaner and can be used solely for reacting to signals.

Add the ```
ShipSpawnSystem
```

to the list of systems in the ```
AsteroidsShipConfig
```

asset after the ```
AsteroidsShipSystem
```

.

## Update the Movement

Until now the ship movement in the ```
AsteroidsShipSystem
```

 always moved using inputs from player 0:

C#

```csharp
var input = f.GetPlayerInput(0);

```

Replace it with the following code to get the input from the linked player:

C#

```csharp
Input\* input = default;
if(f.Unsafe.TryGetPointer(filter.Entity, out PlayerLink\* playerLink))
{
input = f.GetPlayerInput(playerLink->PlayerRef);
}

```

Note that the filter has not been adjusted so the system will still filter for entities with a PhysicsBody2D but no PlayerLink component. In this case it will use the ```
default
```

value for the input. This results in no movement besides gravity being applied. This pattern makes it easy to add AI controlled ships by passing in input from a different source.

Getting a component using ```
TryGet
```

 in Quantum is very fast ```
O(1)
```

because Quantum uses a sparse set ECS.

```
f.Unsafe
```

provides access to Quantum's unsafe API which is generally faster and more convenient to use than the no-pointer counterparts. However, if you want to avoid using pointers the safe API such as ```
f.Get
```

 can be used.

Switch to Unity and enter play mode. The ship will be spawned and it reacts to keyboard inputs.

Back to top

- [Overview](#overview)
- [Player Prefab](#player-prefab)
- [Player Link Component](#player-link-component)
- [Player Data](#player-data)
- [Spawning Player Objects](#spawning-player-objects)
- [Update the Movement](#update-the-movement)