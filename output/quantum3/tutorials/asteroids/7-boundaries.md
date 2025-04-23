# 7-boundaries

_Source: https://doc.photonengine.com/quantum/current/tutorials/asteroids/7-boundaries_

# 7 - Boundaries

## Overview

At the moment, the ship and asteroids can move without bounds and leave the camera view. In the asteroids game when a ship or asteroid leaves the screen it appears again on the opposite side of the screen. A simple way to replicate this behavior is with a system that checks whether entities are inside bounds and teleports them accordingly.

## Boundaries

Create a ```
BoundarySystem
```

 script and add the following code to it:

C#

```csharp
using UnityEngine.Scripting;
using Photon.Deterministic;

namespace Quantum.Asteroids
{
\[Preserve\]
public unsafe class BoundarySystem : SystemMainThreadFilter<BoundarySystem.Filter>
{
public struct Filter
{
public EntityRef Entity;
public Transform2D\* Transform;
}

public override void Update(Frame frame, ref Filter filter)
{
if (IsOutOfBounds(filter.Transform->Position, new FPVector2(10, 10), out FPVector2 newPosition))
{
filter.Transform->Position = newPosition;
filter.Transform->Teleport(frame, newPosition);
}
}

/// <summary>
/// Test if a position is out of bounds and provide a warped position.
/// When the entity leaves the bounds it will emerge on the other side.
/// </summary>
public bool IsOutOfBounds(FPVector2 position, FPVector2 mapExtends, out FPVector2 newPosition)
{
newPosition = position;

if (position.X >= -mapExtends.X && position.X <= mapExtends.X &&
position.Y >= -mapExtends.Y && position.Y <= mapExtends.Y)
{
// position is inside map bounds
return false;
}

// warp x position
if (position.X < -mapExtends.X)
{
newPosition.X = mapExtends.X;
}
else if (position.X > mapExtends.X)
{
newPosition.X = -mapExtends.X;
}

// warp y position
if (position.Y < -mapExtends.Y)
{
newPosition.Y = mapExtends.Y;
}
else if (position.Y > mapExtends.Y)
{
newPosition.Y = -mapExtends.Y;
}

return true;
}
}
}

```

This system loops over all entities with a transform on it and teleports then to the opposite side of the map when out of bounds. For a more selective approach a tag component could be used in the filter.

The ```
Teleport
```

function on the transform is used to signal to the ```
EntityViewInterpolator
```

that interpolation for this movement should be skipped. Without the teleport command the entity would be interpolated from one end of the screen to the other when teleporting from one edge of the screen to the other.

Add the ```
BoundarySystem
```

to the ```
AsteroidsSystemConfig
```

ScriptableObject and enter play mode. The spaceship and the asteroids are correctly teleported when reaching the boundary.

## Configurable Boundaries

Currently, the boundary is hard coded to a map size of 20x20 units. To make this more flexible the map size can be added to the ```
AsteroidsGameConfig
```

.

Open the ```
AsteroidsGameConfig
```

 script and add the following:

C#

```csharp
\[Header("Map configuration")\]
\[Tooltip("Total size of the map. This is used to calculate when an entity is outside de gameplay area and then wrap it to the other side")\]
public FPVector2 GameMapSize = new FPVector2(25, 25);

public FPVector2 MapExtends => \_mapExtends;

private FPVector2 \_mapExtends;

```

A public member field called ```
GameMapSize
```

is used to adjust the size of the map in the config as it is more intuitive for a designer to work with total map width and high instead of the extends. However, for the gameplay code using the extends is easier and more performant. A common pattern in config files is to calculate additional data once when the config is loaded in to the game. This can be done using the ```
Loaded
```

function in the following way:

C#

```csharp
public override void Loaded(IResourceManager resourceManager, Native.Allocator allocator)
{
base.Loaded(resourceManager, allocator);

\_mapExtends = GameMapSize / 2;
}

```

Return to the ```
BoundarySystem
```

script and adjust the code to load the config and use the extends from it:

C#

```csharp
public override void Update(Frame frame, ref Filter filter)
{
 AsteroidsGameConfig config = frame.FindAsset(frame.RuntimeConfig.GameConfig);

 if (IsOutOfBounds(filter.Transform->Position, config.MapExtends, out FPVector2 newPosition))
 {
 filter.Transform->Position = newPosition;
 filter.Transform->Teleport(frame, newPosition);
 }
}

```

Return to Unity and adjust the ```
GameMapSize
```

 in the config. A value of ```
(70, 40)
```

matches the visuals when using a ```
16:9
```

 aspect ratio.

Enter play mode. The spaceship and asteroids now use the bounds as set in the config file.

![GIF of boundaries during play mode](/docs/img/quantum/v3/tutorials/asteroids/7-boundaries.gif)Back to top

- [Overview](#overview)
- [Boundaries](#boundaries)
- [Configurable Boundaries](#configurable-boundaries)