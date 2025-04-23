# sample-project

_Source: https://doc.photonengine.com/quantum/current/addons/kcc/sample-project_

# Sample Project

## Overview

The **KCC** comes with a full Unity project with examples, testing playground and handy code snippets to better understand how KCC works and to speed up onboarding process.

![KCC](/docs/img/quantum/v3/addons/kcc/overview-3.jpg)

The sample project can be downloaded in [Download](/quantum/current/addons/kcc/download) section.

## Features

- Playground scene - Stairs, Slopes, Corridors, Gaps.
- Interaction examples - Teleport, Jump Pad.
- Simple NPC navigation example.
- Input smoothing.
- Support for PC / Mobile.
- Photon Menu integration.

## Sample Controls

- ```
Mouse
```

\- Look
- ```
W
```

,```
S
```

,```
A
```

,```
D
```

\- Move
- ```
Space
```

\- Jump
- ```
Enter
```

\- Lock/unlock cursor

## Project structure

- ```
Assets/Photon/QuantumAddons/KCC
```

\- Base addon folder.

  - ```
    AssetDB
    ```

     \- Contains basic/minimal setup (KCC entity, KCC settings) and default KCC processor assets for movement (EnvironmentProcessor) and post-processing after move (GroundSnapProcessor, StepUpProcessor).
  - ```
    Example
    ```

     \- Contains simple game loop (Menu <=> Playground scene), player controller view from first-person and third-person perspective and other example gameplay elements related to character movement.
  - ```
    Simulation
    ```

     \- Core KCC addon scripts.

## Recommended walkthrough

1. Try ```
Playground
```

    scene.
2. Check player implementation - ```
PlayerInput
```

, ```
PlayerSystem
```

, ```
FirstPersonCamera
```

, ```
ThirdPersonCamera
```

.
3. Check world object interactions - ```
JumpPadProcessor
```

, ```
TeleportProcessor
```

.
4. Check movement implementation - ```
EnvironmentProcessor
```

, ```
SimpleMoveProcessor
```

.
5. To get more context it's good time to explore and learn about [Processors](/quantum/current/addons/kcc/processors).

Now you should have a basic understanding how to move with KCC, how the KCC interacts with other objects (processors) and how they modify player behavior.

6. Check rest of the sample project.
7. Take some scripts to your project or cleanup this and make a game 🚀

⚠️ The sample uses custom input smoothing get super smooth camera rotation, check [Input](/quantum/current/addons/kcc/input) section for more info.

Back to top

- [Overview](#overview)
- [Features](#features)
- [Sample Controls](#sample-controls)
- [Project structure](#project-structure)
- [Recommended walkthrough](#recommended-walkthrough)