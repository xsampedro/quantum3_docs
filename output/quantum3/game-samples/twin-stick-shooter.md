# twin-stick-shooter

_Source: https://doc.photonengine.com/quantum/current/game-samples/twin-stick-shooter_

# Twin Stick Shooter

![Level 4](/v2/img/docs/levels/level03-intermediate_1.5x.png)

## Overview

This sample is provided with full source code and demonstrates how Quantum can be used to create a top-down twin stick shooter game.

It showcases AI systems present in Bot SDK in conjunction to other auxilliary implementations such as an AI Director to provide team strategy, data-driven AI Sensors and Habilities architecture and much more.

## Download

| Version | Release Date | Download |
| --- | --- | --- |
| 3.0.2 | Apr 08, 2025 | [Quantum TwinStickShooter 3.0.2 Build 620](https://dashboard.photonengine.com/download/quantum/quantum-twinstickshooter-3.0.2.zip) |

## Technical Info

- Unity: 2021.3.30f1;
- Platforms: PC (Windows / Mac), and Mobile (Android);

## Highlights

### AI

- Bot SDK sample usage;
- Hierarchical Finite State Machine (HFSM) used as the "Brain" for characters controlled by AI;
- Filling a game match with Bots and picking random names from a text file;
- Replacement of disconnected players by Bots;
- Data driven architecture for AI Sensors;
- Tactical sensors, used by the Bots to judge what tactics they want to perform;
- Re-usable Input struct for both Players and Bots: both uses the same data and systems;
- AI Director which polls team relevant data and defines the team Strategies;
- AIMemory: store data and only make it available after an interval, and "forget" the data after an interval;

### General

- Usage of the HFSM to create the Game Manager, which dictates the flow of the game mode;
- Data driven architecture for Habilities;
- Usage of a Top-Down KCC;
- Union-based strategy for character Attributes (such as Health, Speed, etc);
- Custom baking of Level Design markers, used by the Bots as data on their decision making;
- Context Steering: consider many "movement desire" vectors to result in a single one, used to movement the Bots;
- Usage of the Callbacks version of the navigation system;

### Game/Level Design

- Usage of entity prototypes (prefabs and scene prototypes);
- 3 unique characters with 2 habilities each;
- Coin Grab game mode: collect the coins in the map. The team who keeps 10+ coins for 15 seconds wins the match;

## Stream Videos

### Photon Insiders - Fireside Chat - Twin Stick Shooter (04 mar 2022)

- Gameplay session with 4 players and 2 Bots;
- Quick view on the Unity project;
- Overall analysis of the game architecture and some of it's main features;
- Q & A, mostly focused on Bot SDK and AI coding;

Game Start - Players and Bots joining

![Game Start](/docs/img/quantum/v2/game-samples/twin-stick-shooter/Game Start - Player and Bot joining.jpg)

Input Polling - Includes player replacement

![Input Polling](/docs/img/quantum/v2/game-samples/twin-stick-shooter/Polling Input.jpg)

Applying the Input to characters

![Applying the Input to characters](/docs/img/quantum/v2/game-samples/twin-stick-shooter/Applying the Input.jpg)

AI Building Blocks

![AI Building Blocks](/docs/img/quantum/v2/game-samples/twin-stick-shooter/The AI building blocks.jpg)

AI Strategy and Tactics

![AI Strategy and Tactics](/docs/img/quantum/v2/game-samples/twin-stick-shooter/Strategy and Tactics.jpg)

The Game Manager

![The Game Manager](/docs/img/quantum/v2/game-samples/twin-stick-shooter/Game Management.jpg)## Screenshots

![](/docs/img/quantum/v2/game-samples/twin-stick-shooter/gameplay.png)

![](/docs/img/quantum/v2/game-samples/twin-stick-shooter/gameplay.gif)

![](/docs/img/quantum/v2/game-samples/twin-stick-shooter/hfsm-debug.gif)

## Third Party Assets

This sample includes third-party free and CC0 assets. The full packages can be acquired for your own projects at their respective site:

- [Carton FX Remaster Free](https://assetstore.unity.com/packages/vfx/particles/cartoon-fx-remaster-free-109565) by Jean Moreno
- [Kay Loysberg - Dungeon assets](https://kaylousberg.com/game-assets) by Kay Lousberg
- [100-cc0-fx](https://opengameart.org/content/100-cc0-sfx) by rubberduck

Back to top

- [Overview](#overview)
- [Download](#download)
- [Technical Info](#technical-info)
- [Highlights](#highlights)

  - [AI](#ai)
  - [General](#general)
  - [Game/Level Design](#gamelevel-design)

- [Stream Videos](#stream-videos)

  - [Photon Insiders - Fireside Chat - Twin Stick Shooter (04 mar 2022)](#photon-insiders-fireside-chat-twin-stick-shooter-04-mar-2022)

- [Screenshots](#screenshots)
- [Third Party Assets](#third-party-assets)