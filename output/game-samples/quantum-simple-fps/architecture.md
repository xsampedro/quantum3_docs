# architecture

_Source: https://doc.photonengine.com/quantum/current/game-samples/quantum-simple-fps/architecture_

# Architecture

## Simulation Systems

#### Gameplay System

The main system responsible for controlling state of the game, player spawn, calculation of statistics and checking win conditions.

#### Player System

Player system processes player input and propagates actions to specific components - for example it sets input direction for character controller (Quantum KCC).

#### Health System

Health system controls player health, taking damage and short immortality after respawn.

#### Pickup System

Pickup system updates pickups in the world and is responsible for collecting when a player passes through a pickup.

#### Weapons System

Weapons system maintains player weapons, switching, reloading and shooting.

#### Lag Compensation Systems

Lag compensation systems provide a way to make physics queries (for example raycasts) against snapshot interpolated entities.

These systems dynamically create and destroy "proxy" entities within Quantum simulation. A proxy mimicks its reference entity at the moment of snapshot interpolation.

Player shooting is then evaluated against proxies instead of other player entities, which eliminates mispredictions and makes shooting accurate.

## Views

#### Player View

Main player script which handles visuals and camera.

#### Weapon/Weapons View

Scripts that synchronize weapons, their visuals and react on simulation events.

#### Health View

Reacts on health changes - damage, spawns hit effects and controls immortality indicator.

#### Pickup View

Controls visual of the pickup based on its simulation state.

## Others

#### Projectile Visual

Regular MonoBehaviour script which controls projectile flying through the environment with hit effect at the end.

#### UI

The user interface in the **Quantum Simple FPS** is handled in a straightforward manner, without relying on any specific UI framework. The main script responsible for UI management is the ```
GameUI
```

 script. It enables the appropriate UI game objects that should be visible during gameplay. Additionally, the ```
GameUI
```

script maintains references to other UI elements and updates the player UI when a player is spawned.

Back to top

- [Simulation Systems](#simulation-systems)
- [Views](#views)
- [Others](#others)