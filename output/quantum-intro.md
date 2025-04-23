# quantum-intro

_Source: https://doc.photonengine.com/quantum/current/quantum-intro_

# Quantum 3 Intro

## Overview

Quantum is a high-performance deterministic ECS (Entity Component System) framework for online multiplayer games made with Unity, for up to 128 players.

Using predict/rollback networking, it is both ideal for latency-sensitive online games like sports games, fighting games, & FPS and also robust when faced with larger latency network connections.

Quantum is free for development. Start building today, launch a game and pay only when its success grows beyond our free tier. [Quantum's pricing](https://www.photonengine.com/quantum/pricing) enables small titles to grow big.

![Quantum Decoupled Architecture](https://doc.photonengine.com/docs/img/quantum/v2/getting-started/quantum-intro/quantum-sdk-layers.jpg)
Quantum helps developers write clean code. It decouples simulation logic (Quantum ECS) from the view/presentation (Unity), and takes care of the network implementations specifics (predict/rollback + transport layer + game agnostic server logic).

The state-of-the-art tech stack is composed of the following pieces:

- Server-managed predict/rollback simulation core.
- Sparse-set ECS memory model and API.
- Complete set of deterministic libraries (math, 2D and 3D physics, navigation, animation, bots, etc.).
- Rich Unity editor integration and tooling.

All built on top of mature and industry-proven existing Photon products and infrastructure (photon realtime transport layer, photon server plugin to host server logic, etc.);

### Predict / Rollback Networking

In deterministic systems, game clients only exchange player input with the simulation running locally on all clients.

Game clients are free to advance the simulation locally using input prediction, and an advanced rollback system takes care of restoring game state and re-simulating any mispredictions. The developer defines the player inputs which will be distributed to all clients and doesn’t write any other networking code.

Quantum uses a game-agnostic authoritative server component, running on our photon servers, to synchronize clocks and manage input latency. Clients never need to wait for a slower client or a client with a worse network connection as they would in lockstep networking.

This server component can be extended to integrate with customer hosted back-end systems, e.g. matchmaking, player services, authoritative ‘referee’ simulation.

![Quantum Server-Managed predict/Rollback](https://doc.photonengine.com/docs/img/quantum/v2/getting-started/quantum-intro/quantum-client-server.jpg)
In Quantum, input exchange is managed via game-agnostic server logic.

The Quantum client library communicates with Quantum server and runs local simulation, performing all input prediction and rollbacks.

Custom game code is written by the developer as an isolated, pure C# simulation, using the Quantum ECS, which is decoupled from Unity. Quantum's interface offers a great range of deterministic APIs that can be reused in any game. For example: vector math, 2D and 3D physics engines, navigation pathfinder, animation, bots.

## Entity Component System

To enable all the simulation code to be high-performance out of the box, Quantum is based on a high performance sparse-set Entity Component System (ECS) model.

The key to Quantum's performance is the use of pointer-based C# code combined with its custom heap allocator. The simulation code allocates no C# heap memory at runtime and creates no garbage.

Even including the re-simulations caused by input mispredictions, which are inherent to the predict/rollback approach, the goal is to leave most of the CPU budget for the view/rendering code, i.e. Unity.

### Code Generation

In Quantum, all gameplay data (game state) is kept either in the sparse-set ECS data structures (entities and components) or in our custom heap-allocator (dynamic collections and custom data), as blittable memory-aligned C# structs.

Pointer-based C# is used by Quantum, for performance. Complexity and boiler-plate code is hidden from the developer by using a simple domain specific language (DSL), called Qtn, to automatically generate C# code for game data.

Qtn

```cs
// components define reusable game state data groups
component Resources
{
  Int32 Mana;
  FP Health;
}

```

Structs, C-style unions, enums, flags, bitsets, lists, and dictionaries, etc, can also be defined and used.

The auto-generated API makes it easy to query and modify the game state with comprehensive functions to iterate, modify, create or destroy entities (based on their components):

C#

```csharp
var es = frame.Filter<Transform3D, Resources>();

// sets the entity ref and pointers to the components
while (es.NextUnsafe(out var entity, out var transform, out var resources)) {
  transform->Position += FPVector3.Forward * frame.DeltaTime;
}

```

### Stateless Systems

While Quantum's DSL covers data definition, there needs to be a way to organize the custom game logic that will update this game state.

Customer game logic is written by implementing Systems, which are stateless pieces of logic that will be executed every tick update by Quantum's client simulation loop. All the game state is stored in the Frame, which is passed into the Update() function.

C#

```csharp
public unsafe class LogicSystem : SystemMainThread
{
  public override void Update(Frame frame)
  {
    // customer game logic here
    // (frame is a reference for the generated game state container).
  }
}

```

## Quantum and Unity

Since Quantum and Unity are decoupled, communication into and out of the Quantum Simulation is well defined:

![Quantum Inputs and Outputs](https://doc.photonengine.com/docs/img/quantum/v3/getting-started/quantum-intro/quantum-inputs-outputs.jpg)### Asset Database

Unity is known for its flexible editor and asset pipeline. [Assets are defined in Quantum](/quantum/current/manual/assets/assets-simulation), then may be [created within Unity](/quantum/current/manual/assets/assets-unity) specifically to be shared with Quantum, allowing game and level designers to work as flexibly as they normally would in Unity:

![Character Classes - Asset Linking](https://doc.photonengine.com/docs/img/quantum/v3/asset-linking.png)
Example of a Quantum Asset for a character which is created & modified in the Unity Editor.

These assets are also available in the Unity code, of course.

### Input

[Input](/quantum/current/manual/input) must be defined and is sent to the server and distributed to all game clients every tick. This would typically be the subset of keyboard or controller buttons and mouse or controller stick positions required by the game.

### Commands

[Commands](/quantum/current/manual/commands) are intended for occasional actions. They are only sent when required.

### Events

[Events](/quantum/current/manual/quantum-ecs/game-events) are a way to transfer information from the Quantum simulation to the Unity view.

### Full Simulation State

The full simulation state from Quantum is observable from Unity. Common cases, for example synchronizing the transforms of GameObjects to their corresponding Entities in Quantum, we support out of the box. [See Entity Prototypes](/quantum/current/manual/entity-prototypes). For game-specific items e.g. character health this can be simply read from the simulation state for the game to display.

## Where to go next

To get started with Quantum we strongly recommend beginning with the [Asteroids Tutorial](/quantum/current/tutorials/asteroids/1-overview). This tutorial teaches you all the necessary basics to get started with Quantum.

For a video tutorial there is also the [Complete Course to Quantum 3](/quantum/current/tutorials/video-tutorial) stream.

To dive right in at a more advanced level we recommend the first pages you should read in the manual are the ones covering the Quantum [DSL](/quantum/current/manual/quantum-ecs/dsl) (Domain Specific Language) and the [Entity Prototypes](/quantum/current/manual/entity-prototypes). These are at the core of programming the simulation and designing the view, respectively.

Alternatively there are many [Game Samples](/quantum/current/game-samples/platform-shooter-2d/overview) to download.

Back to top

- [Overview](#overview)

  - [Predict / Rollback Networking](#predict-rollback-networking)

- [Entity Component System](#entity-component-system)

  - [Code Generation](#code-generation)
  - [Stateless Systems](#stateless-systems)

- [Quantum and Unity](#quantum-and-unity)

  - [Asset Database](#asset-database)
  - [Input](#input)
  - [Commands](#commands)
  - [Events](#events)
  - [Full Simulation State](#full-simulation-state)

- [Where to go next](#where-to-go-next)