# project-architecture

_Source: https://doc.photonengine.com/quantum/current/game-samples/quantum-starter/project-architecture_

# Project Architecture

## Project Organization

The project consists of three independent examples, organized in the following folder structure:

|     |     |
| --- | --- |
| /00\_MainMenu | Contains the main menu scene |
| /01\_ThirdPersonCharacter | Third Person Character example |
| /02\_Platformer | Platformer example |
| /03\_Shooter | Shooter example |
| /Common | Shared resources (prefabs, scripts, and graphics) used across all examples |

Each example folder and the _Common_ folder contains a _Simulation_ folder with Quantum simulation-related scripts. For convenience in the Starter sample, simulation code is placed next to its related assets. While the default simulation folder is _QuantumUser/Simulation_, you can choose different folders for simulation code as long as the `Quantum.Simulation.asmref` file is present.

The _Common/Simulation_ folder contains shared simulation code (like input and movement system) that all examples use. Make sure to check this folder when looking for common functionality.

Each example has its own _SystemConfig_ asset that enables only the Quantum systems needed for that specific example. For reference, you can check the _/02\_Platformer/Configs/SystemsConfig-Platformer.asset_ file to see how systems are configured for the Platformer example.

![System Config](/docs/img/quantum/v3/game-samples/starter/SystemConfig.jpg)

The game launches from the `UIGameMenu` component, which is located on the _UIGameMenu_ GameObject in each example scene (for example, the _UI/UIGameMenu_ GameObject in the _02\_Platformer_ scene). This component contains the specific setup required for each example (e.g. the system config).

## 1 - Third Person Character

![Third Person Character](/docs/img/quantum/v3/game-samples/starter/ThirdPersonCharacter.jpg)

The **Third Person Character** sample is a conversion of Unity's [Starter Assets - Third Person](https://assetstore.unity.com/packages/essentials/starter-assets-thirdperson-updates-in-new-charactercontroller-pa-196526) into a multiplayer environment. Players can spawn as third-person characters, walk, and run around in a prototype environment. Instead of using the default `CharacterController` component, this example uses the [KCC addon](/quantum/current/addons/kcc/overview) to drive player movement. The KCC addon is our kinematic character controller solution specifically tailored for Quantum. It provides smooth movement and rotation even under the most challenging networking conditions while being highly optimized to handle dozens or even hundreds of characters.

### Where to go next

To learn more about character movement, check out the [KCC Sample Project](/quantum/current/addons/kcc/overview). For more advanced character animation, explore the [Simple FPS](/quantum/current/game-samples/quantum-simple-fps/overview) project. If you're making a fighting or sports game where animations need precise synchronization, look into our [Quantum Animator](/quantum/current/addons/animator/overview) addon.

## 2 - Platformer

![Platformer](/docs/img/quantum/v3/game-samples/starter/PlatformerJumping.jpg)

The **Platformer** example builds on a similar foundation as the _Third Person Character_ example and adds interaction with game objects (e.g., coins, falling platforms), use of RuntimePlayer data (player nicknames), and a simple game loop where players race to collect 10 coins and reach the top flag. After each round, all players return to the starting point.

### Where to go next

The [Platform Shooter 2D](/quantum/current/game-samples/platform-shooter-2d/overview) sample provides a fast and action-packed platformer experience. Physics interactions are where Quantum truly shines. See how players interact with each other and with a ball in [Sports Arena Brawler](/quantum/current/game-samples/sports-arena-brawler).

## 3 - Shooter

![Shooter](/docs/img/quantum/v3/game-samples/starter/Shooter.jpg)

The **Shooter** example showcases a simple first-person shooter. Players compete to be the best hunter by shooting flying chickens. The chicken counter resets when a player dies, either by falling from a platform or being killed by another player. This example demonstrates simple raycast logic and the use of custom assets (`ChickenConfig`).

### Where to go next

Since shooters are the most popular multiplayer genre, we provide plenty of resources to jumpstart your shooter game development journey. [Simple FPS](/quantum/current/game-samples/quantum-simple-fps/overview) is a natural evolution of this example — it provides a complete game loop, ammo handling, various weapons, pickups, player statistics and leaderboards, and integration with the Quantum Menu. [Platform Shooter 2D](/quantum/current/game-samples/platform-shooter-2d/overview) provides a comprehensive platformer shooter experience in 2D. Finally, our top-down sample [Twin Stick Shooter](/quantum/current/game-samples/twin-stick-shooter) will teach you how to build a game similar to Brawl Stars and includes full implementation of player AI using the [BotSDK addon](/quantum/current/addons/bot-sdk/overview).

[Quantum Simple FPS](/quantum/current/game-samples/quantum-simple-fps/overview) also provides an example implementation of lag compensation and snapshot interpolation, but keep in mind that [Fusion](/fusion/current/fusion-choose) might be a better choice for shooters, especially when creating a competitive game or targeting PC and consoles.

## More Quantum Samples

Visit our [Samples](https://www.photonengine.com/samples) page to discover more example projects.

Back to top

- [Project Organization](#project-organization)
- [1 - Third Person Character](#third-person-character)

  - [Where to go next](#where-to-go-next)

- [2 - Platformer](#platformer)

  - [Where to go next](#where-to-go-next-1)

- [3 - Shooter](#shooter)

  - [Where to go next](#where-to-go-next-2)

- [More Quantum Samples](#more-quantum-samples)