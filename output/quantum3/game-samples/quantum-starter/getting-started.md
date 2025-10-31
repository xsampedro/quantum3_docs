# getting-started

_Source: https://doc.photonengine.com/quantum/current/game-samples/quantum-starter/getting-started_

# Getting Started

## Requirements

- Unity: 6000.0.44f1 or later
- Quantum 3 AppId, which can be created in the [Photon Engine Dashboard](https://dashboard.photonengine.com). For detailed instructions, see [How to create a Quantum 3 AppId](/quantum/current/reference/create-quantum-appid).

Before diving into this sample, we recommend having a basic understanding of Quantum (see [Quantum 3 Intro](/quantum/current/quantum-intro)) and completing the [Asteroids Tutorial](/quantum/current/tutorials/asteroids/1-overview) first.

## Download

| Version | Release Date | Download |
| --- | --- | --- |
| 3.0.5 | Aug 07, 2025 | [Quantum Starter 3.0.5](https://downloads.photonengine.com/download/quantum/quantum-starter-3.0.5.zip?pre=sp) |
| 3.0.2 | May 05, 2025 | [Quantum Starter 3.0.2](https://downloads.photonengine.com/download/quantum/quantum-starter-3.0.3.zip?pre=sp) |

## Running in Editor

To run the sample, set your _AppId_ in the _Quantum Hub_ window (menu _Window/Quantum/Quantum Hub_).

Each example within **Quantum Starter** is accessible via its own scene file, which can be opened and played directly. Alternatively, all examples can be launched from the **MainMenu** scene located at `/00\_MainMenu/00\_MainMenu`.

![Main Menu Scene](/docs/img/quantum/v3/game-samples/starter/MainMenu.jpg)

Upon starting the game, a small game menu appears where players can enter a nickname or a session name (also known as a room name). Clicking the _Start Game_ button will initiate a new game session, or connect the player to an already existing game.

![Game Menu](/docs/img/quantum/v3/game-samples/starter/ThirdPersonCharacterMenu.jpg)

To play the game offline in _Local Mode_, enable the _ForceLocalMode_ flag on the _UIGameMenu_ GameObject in your chosen scene (for example, on the _UI/UIGameMenu_ object in the _03\_Shooter_ scene).

![Local Mode](/docs/img/quantum/v3/game-samples/starter/LocalMode.jpg)

## Gameplay Controls

- Use the `W`, `S`, `A`, `D` keys for movement.
- Use the `Shift` key for sprint and the `Space` key for jump.
- Use the mouse for looking around.
- Use `Left mouse button` for weapon fire _(only in the Shooter example)_
- Press the `Esc` key to open the game menu during play.

## Making a Build

The sample can be built for all desktop platforms (Windows, Mac, Linux) and web platforms. While it is possible to create a mobile version of the sample, please note that touch controls are not currently implemented.

Back to top

- [Requirements](#requirements)
- [Download](#download)
- [Running in Editor](#running-in-editor)
- [Gameplay Controls](#gameplay-controls)
- [Making a Build](#making-a-build)