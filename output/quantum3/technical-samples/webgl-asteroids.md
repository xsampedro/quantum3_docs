# webgl-asteroids

_Source: https://doc.photonengine.com/quantum/current/technical-samples/webgl-asteroids_

# Quantum Crazy Starter WebGL

![Level 4](/v2/img/docs/levels/level01-beginner_1.5x.png)

![Crazy Games](/docs/img/quantum/v3/technical-samples/webgl-asteroids/img-init.png)## Overview

Take aim and blast some asteroids in this WebGL sample, created aiming the [Crazy Games](https://www.crazygames.com/) platform, [Gameplay link](https://www.crazygames.com/game/quantum-crazy-starter-webgl)! This project demonstrates setting up Photon Quantum and Unity3D targeting WebGL as platform, featuring a streamlined Asteroids-inspired game integrated with the `CrazySDK` from Crazy Games. Using [Unity Addressables](https://docs.unity3d.com/Packages/com.unity.addressables@2.3/manual/index.html) to reduces build size while delivering smoother gameplay.

This WebGL sample is a great first step for those new to Photon products!

For a more in-depth, multiplayer-focused experience, check out our [Asteroids tutorial](https://doc.photonengine.com/quantum/current/tutorials/asteroids/1-overview), which guides developers through building a multiplayer game using Photon Quantum.

![Crazy Games](/docs/img/quantum/v3/technical-samples/webgl-asteroids/img-3.png)## Download

| Version | Release Date | Download |
| --- | --- | --- |
| 3.0.2 | Apr 17, 2025 | [Quantum WebGL Sample 3.0.2 Build 623](https://dashboard.photonengine.com/download/quantum/quantum-webgl-sample-3.0.2.zip) |

## Download Vanilla Project

| Version | Release Date | Download |
| --- | --- | --- |
| 3.0.2 | Apr 22, 2025 | [Quantum CrazyGames WebGL Vanilla Sample 3.0.2 Build 630](https://dashboard.photonengine.com/download/quantum/quantum-crazygames-webgl-vanilla-sample-3.0.2.zip) |

## Technical Info

- Unity: 6000.0.17f1.
- Platforms: WebGL

To run the sample in online multiplayer mode, first create a Quantum AppId in the [PhotonEngine Dashboard](https://dashboard.photonengine.com) and paste it into the `AppId` field in `PhotonServerSettings` asset.

Then load the `Menu` scene in the Scenes menu and press `Play`.

## How to Play

1. Pilot your ship through space, taking down asteroids and other players along the way!
2. Choose your style: working with other players to clear the field, or detroy other players to dominate the battlefield.
3. Destroy all asteroids to advance to next level and test your skills against a map with more asteroids!

###

#### Controls

- `A/D` keys to steer;
- `W` key to throttle;
- `Space/Left Click` to shoote;

## Gameplay Implementation

This game was built using the [Asteroids tutorial](https://doc.photonengine.com/quantum/current/tutorials/asteroids/1-overview) as a foundation. Check it out for a step-by-step guide on creating a similar game, from setup to final touches!

###

#### Addressables

The project makes use of [Unity Addressables](https://docs.unity3d.com/Packages/com.unity.addressables@2.3/manual/index.html) assets, so most of the assets used in gameplay are loaded only moments before the match starts. This way the project build size was reduced, which is a good

practice on WebGL games.

![Addressables List](/docs/img/quantum/v3/technical-samples/webgl-asteroids/img-1.png)

Asset loading starts after the play button is pressed, and progress is shown during the connection phase:

![Loading Scene](/docs/img/quantum/v3/technical-samples/webgl-asteroids/img-2.png)

In this sample, the Quantum `Look-up-tables files (LUT)` is treated as an Addressable. As a result, it was moved from the Resources folder to a new location, `Resources\_moved`. This optimization reduces the build size by approximately 2MB.

## CrazySDK

The CrazySDK is integrated into this project, demonstrating the basic features required to have a multiplayer game approved during the project review in the CrazyGames platform.

The [multiplayer requirements](https://docs.crazygames.com/requirements/multiplayer/) are handled by `CrazyManager.cs`, it is responsible for checking if the game needs to be treated as [instantMuliplayer](https://docs.crazygames.com/sdk/game/#instant-multiplayer), generates [invite link](https://docs.crazygames.com/sdk/game/#invite-link) and trigger the [invite button](https://docs.crazygames.com/sdk/game/#invite-button).

More info about these features can be found under CrazyGames documentation.

- Check for updates on the [CrazyGames SDK](https://docs.crazygames.com/).
- The CrazyGames SDK has other functionality available (user, game etc), documentation via [docs.crazygames.com](https://docs.crazygames.com/).

Currently, due to limitations of the CrazyGames platform, running applications in multithreaded mode is not supported. However, this feature can be utilized on platforms like itch.io. Please note that in order to run the sample on an itch.io page, the CrazySDK must be removed from the project.

The next image shows the Multithread disabled in the project:

![Multithread](/docs/img/quantum/v3/technical-samples/webgl-asteroids/img-4.png)## Third Party Assets

This sample includes third-party free and CC0 assets. The full packages can be acquired for your own projects at their respective site:

- [Ulukais Space Skyboxes](https://opengameart.org/content/ulukais-space-skyboxes) by Ulukais
- [Space Kit](https://www.kenney.nl/assets/space-kit) by Kenney
- [Sci-Fi Sounds](https://kenney.nl/assets/sci-fi-sounds) by Kenney
- [Retro Game Weapons Sound Effects](https://freesound.org/) by various
- [Sci-Fi Sound Effects Library](https://opengameart.org/content/sci-fi-sound-effects-library) by Little Robot Sound Factory

Back to top

- [Overview](#overview)
- [Download](#download)
- [Download Vanilla Project](#download-vanilla-project)
- [Technical Info](#technical-info)
- [How to Play](#how-to-play)

- [Gameplay Implementation](#gameplay-implementation)

- [CrazySDK](#crazysdk)
- [Third Party Assets](#third-party-assets)