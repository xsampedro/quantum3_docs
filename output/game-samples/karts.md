# karts

_Source: https://doc.photonengine.com/quantum/current/game-samples/karts_

# Karts

![Level 4](/v2/img/docs/levels/level03-intermediate_1.5x.png)

## Overview

Quantum Karts demonstrates how to build an arcade racing game with a full race loop, AI opponents, powerups and more.

The sample supports up to 12 karts, has 2 tracks, 4 different powerups, 3 kart types and 3 kart looks.

![](https://doc.photonengine.com/docs/img/quantum/v3/game-samples/karts/karts-1.png)

![](https://doc.photonengine.com/docs/img/quantum/v3/game-samples/karts/karts-2.png)

![](https://doc.photonengine.com/docs/img/quantum/v3/game-samples/karts/karts-3.png)

![](https://doc.photonengine.com/docs/img/quantum/v3/game-samples/karts/karts-4.png)

## Download

| Version | Release Date | Download |
| --- | --- | --- |
| 3.0.2 | Mar 24, 2025 | [Quantum Karts 3.0.2 Build 607](https://dashboard.photonengine.com/download/quantum/quantum-karts-3.0.2.zip) |

## Technical Info

- Unity: 2021.3.30f1.
- Platforms: PC (Windows / Mac), WebGL

## Before You Start

To run the sample in online multiplayer mode, first create a Quantum AppId in the [PhotonEngine Dashboard](https://dashboard.photonengine.com) and paste it into the `AppId` field in `PhotonServerSettings` asset.

Then load the `Menu` scene in the Scenes menu and press `Play`.

## Highlights

###

#### Technical

- Arcade racing physics utilizing Broadphase Queries for wheels and collisions.
- Custom friction and drifting physics utilizing various different FP Math features.
- Different surfaces to offer varying driving conditions.
- Input encoding (Vector2 <-> byte).
- Entity pool for re-usable powerup entities.
- AI drivers with adjustable difficulty that can drive through the tracks, drift and use powerups.
- An extendable powerup system that supports different AI behaviours for each powerup.

#### Gameplay

- Satisfying arcade driving with a rewarding skill-based drifting mechanic.
- Full race loop: Ready check, Countdown, Race, Scoreboard.
- Different Kart stats and looks to choose from.
- Get speed boosts by executing long drifts.
- 4 weapons plus variations of them (Bomb, Mine, Shield, Boost).
- 2 tracks with different driving surfaces.

## Controls

- A/D or Left/Right Arrow Keys to Steer
- W or Up Arrow Key to Accelerate
- S or Down Arrow Key to Reverse
- Space to Hop/Drift
- Shift to use Items

## 3rd Party Assets

The sample game was developed by Zaibatsu Studio for Photon Engine.

Back to top

- [Overview](#overview)
- [Download](#download)
- [Technical Info](#technical-info)
- [Before You Start](#before-you-start)
- [Highlights](#highlights)

- [Controls](#controls)
- [3rd Party Assets](#rd-party-assets)