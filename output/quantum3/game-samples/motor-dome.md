# motor-dome

_Source: https://doc.photonengine.com/quantum/current/game-samples/motor-dome_

# Motor Dome

Level

INTERMEDIATE

Topology

**DETERMINISTIC**

## Overview

The **Quantum Motor Dome** sample demonstrates an approach on how to build a free-for-all Snake-like sports game for up to 6 players, featuring hands-off host-or-join lobbying, automatic game state progression with a custom-built FSM system, broadphase hit detection, and player customization.

![](/docs/img/quantum/v3/game-samples/motor-dome/MD1.gif)

![](/docs/img/quantum/v3/game-samples/motor-dome/MD3.gif)

![](/docs/img/quantum/v3/game-samples/motor-dome/MD2.gif)

![](/docs/img/quantum/v3/game-samples/motor-dome/MD4.gif)

## Before You Start

Create a Quantum AppId from the [Photon Dashboard](https://dashboard.photonengine.com/) and paste it into the `AppId` in the `PhotonServerSettings` asset, located in the Unity project under `Assets/Resources`.

Select `StartScene` from the scenes dropdown and press the play button.

## Download

| Version | Release Date | Download |
| --- | --- | --- |
| 3.0.8 | Oct 21, 2025 | [Quantum Motor Dome 3.0.8](https://downloads.photonengine.com/download/quantum/quantum-motor-dome-3.0.8.zip?pre=sp) |
| 3.0.4 | Jun 25, 2025 | [Quantum Motor Dome 3.0.4](https://downloads.photonengine.com/download/quantum/quantum-motor-dome-3.0.4.zip?pre=sp) |
| 3.0.2 | Mar 18, 2025 | [Quantum Motor Dome 3.0.2](https://downloads.photonengine.com/download/quantum/quantum-motor-dome-3.0.2.zip?pre=sp) |

## Highlights

### Technical

- Broadphase Queries for collision handling
- Custom state machine for game state
- Custom baked map data using `MapDataBakerCallback`
- Player nicknames and player character customization

### Gameplay

The core gameplay features include:

- Full game loop, including pre-game lobby and post-game screen.
- Spawn protection
- Randomly Spawned Pickups
- News feed for gameplay events

## Project

### Input

The project uses the legacy Unity input system.

(gamepad inputs are listed as Xbox buttons)

local input:

- `P` key to open the pause menu

quantum input:

- steer with `A` and `D` keys / left and right arrows / left thumbstick X axis
- boost with `W` key / up arrow / A button / RB
- brake with `S` key / down arrow / B button / LB

### Getting Into A Game

Matchmaking:

To play an online match, select `Matchmaking` followed by `Start Queue`. At this point, the Matchmaker takes over and the Quantum session begins. After a brief waiting period, the `LoadSyncSystem` will progress to the intro sequence, after which the game will begin. The game state is entirely autonomous, not relying on player input to progress.

Practice Mode:

To play a solo match, select practice. The `QuantumRunnerLocalDebug` in the game scene handles starting up the Quantum session. The intro phase will be bypassed by the `LoadSyncSystem`.

### Matchmaker

Deriving from `QuantumCallbacks` and implementing `IConnectionCallbacks`, `IMatchmakingCallbacks`, `IInRoomCallbacks`, `IOnEventCallback`, the Matchmaker class takes care of the flow of connecting to a game and delegating the various callbacks. The main entry point is the static `Connect` method, which signals back to the caller when the connection state updates.

### Game State

The Game State System enables and disables systems which implement any of the game state interfaces. Systems may implement one or many interfaces, and will be enabled while in any of the corresponding game states. Systems which do not implement a game state interface remain untouched by the game state system.

When the game state changes, a `GameStateChanged` event is sent, containing both the old and new game state.

The system is designed to be reusable, and can have states changed or added with little effort.

The files relevant to the game state system are:

pre

```pre
quantum.code
└ Game
  ├ Game State
  │ ├ gameState.qtn
  │ └ IGameStates.cs
  └ Systems
    └ GameStateSystem.cs

```

### Events

The `EventSubscriptions` class listens for the majority of events. All "player" events are delegated to the `InterfaceManager` for handling.

State System:

- `GameStateChanged`

Gameplay:

- `PickupCollected`
- `PlayerDied`
- `PlayerLeadGained`
- `PlayerReconnected`
- `PlayerScoreChanged`

Other:

- `PlayerLeft`
- `PlayerDataChanged`
- `Shutdown`

Events listened to by the `ShipView` class:

- `PlayerVulnerable`
- `PlayerDataChanged`

### Collisions

Collisions between players leverage Quantum's [broadphase queries](https://doc.photonengine.com/quantum/current/manual/physics/queries#broadphase_queries). Pickups use `PhysicsCollider3D` and `OnTriggerEnter3D` for collision detection.

Broadphase:

`ShipCollisionInjectionSystem.cs` \- For each ship, a linecast query is added for each segment of the ship's trail.

`ShipCollisionRetrievalSystem.cs` \- Evaluates collisions between ships and trails, and differentiates a ship colliding with the end of its own trail for scoring.

Pickups:

`PickupSystem.cs` \- Collisions with pickups rely on ships having a `PhysicsCollider3D` with the callback flag `OnDynamicTriggerEnter`, and pickups having a `PhysicsCollider3D` with `IsTrigger` ticked. These are configured on the entity prototypes, and can be found in `Assets/Resources/DB/Quantum Prefabs` of the Unity project.

### Custom Baked Map Data

This sample utilizes the `UserAsset` field of the `MapData` asset to associate additional information with the map.

See [Map Baking](https://doc.photonengine.com/quantum/current/manual/map-baking) for more information.

Quantum:

- `MapMeta.cs` \- Stores the position of spawnpoints, and the size and origin of the play space.

Unity:

- `MapDataMeta.cs` \- A MonoBehaviour which holds information to be baked to the `MapMeta` asset.
- `CustomMapBaker.cs` \- Derives from `MapDataBakerCallback` in order to bake the view-side map data into the Quantum-side `MapMeta` asset.

### Customization

Unity:

- `LocalData.cs` \- Holds customization data before it has been sent to Quantum.

Quantum:

- `RuntimePlayer.User.cs` \- Serializes the player's nickname, ship prototype, and selected colors. The view can then access these by retrieving the `RuntimePlayer` via the frame's `GetPlayerData` method.

The below snippet shows how a `ColorRGBA` may be serialized. Note that the alpha is not used, and so does not need to be serialized. Each color then uses only 3 bytes.

C#

```cs
partial void SerializeUserData(BitStream stream)
{
    /* ... */
    stream.Serialize(ref primaryColor.R);
    stream.Serialize(ref primaryColor.G);
    stream.Serialize(ref primaryColor.B);
    /* ... */
}

```

## 3rd Party Assets

The sample game was developed by Nthusia Studio for Photon Engine.

Back to top

- [Overview](#overview)
- [Before You Start](#before-you-start)
- [Download](#download)
- [Highlights](#highlights)

  - [Technical](#technical)
  - [Gameplay](#gameplay)

- [Project](#project)

  - [Input](#input)
  - [Getting Into A Game](#getting-into-a-game)
  - [Matchmaker](#matchmaker)
  - [Game State](#game-state)
  - [Events](#events)
  - [Collisions](#collisions)
  - [Custom Baked Map Data](#custom-baked-map-data)
  - [Customization](#customization)

- [3rd Party Assets](#rd-party-assets)