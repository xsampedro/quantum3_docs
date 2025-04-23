# 8-multiplayer

_Source: https://doc.photonengine.com/quantum/current/tutorials/asteroids/8-multiplayer_

# 8 - Multiplayer

## Overview

At this point the game is playable in single player mode by entering play mode in the ```
AsteroidsGameplay
```

 scene. Turning gameplay like this into a high quality multiplayer game is usually quite the challenge. But as you can see this is already the last section of this tutorial. So how will we go about adding multiplayer?

The answer is that multiplayer support has already been there all along! With Quantum you simply focus on writing gameplay code and the

Quantum engine turns it into a multiplayer game! All that is left is to hook up a menu scene quickly that allows players to join the game online.

## Menu Scene

Quantum comes with a built-in sample menu for prototyping. This menu can be extracted by double-clicking the ```
Quantum-Menu
```

Unity package in the ```
Photon/QuantumMenu
```

folder. Extract the package into the project.

Open the ```
QuantumSampleMenu
```

scene in the ```
Photon/QuantumMenu
```

folder. The menu scene is a sample scene that comes with all functionality needed to run a Quantum game online including a lobby system. As you open the scene a popup will appear asking to import TMP Essentials. Press the button to import TMP essentials into your project.

Once the import is finished go to ```
File > Save As
```

and save the scene as ```
Menu
```

in the ```
Scenes
```

folder to create a copy.

![The Menu scene](/docs/img/quantum/v3/tutorials/asteroids/8-menu.png)
The Menu scene.


In this example we will only use the Quick Play functionality of the menu.

## Configuring the Menu

The menu is mainly configured way the ```
Menu Config
```

 asset.

Create a new menu config asset (right click on ```
Resources
```

and selelect ```
Quantum > Menu > Menu Config
```

). Name the config ```
AsteroidsMenuConfig
```

.

Add a new entry to the ```
Available Scenes
```

 list and link the ```
AsteroidsGamplay
```

and the ```
AsteroidsMap
```

 asset that belongs to the scene. Name the scene "Asteroids Gameplay". Expand the ```
Runtime Config
```

field in the inspector and set the ```
AsteroidsMap
```

, ```
AsteroidsSimulationConfig
```

, ```
AsteroidsSystemConfig
```

 and ```
AsteroidsGameConfig
```

in their respective fields.

Also select the default ```
Machine Id
```

 and ```
Code Generator
```

assets.

You can also configure available regions and app versions in the menu config.

![The Menu Config](/docs/img/quantum/v3/tutorials/asteroids/8-menu-config.png)
The Menu Config.


Inside the ```
Menu
```

 scene find the ```
Canvas/QuantumMenu
```

GameObject. Start by expanding the ```
Default Connection Args
```

 field and adding a player to ```
Runtime Players
```

and drag in the ```
AsteroidsShipEntityPrototype
```

 into the ```
Player Avatar
```

field.

Next, replace the ```
Config
```

 with the newly created ```
AsteroidsMenuConfig
```

.

The ```
Server Settings
```

 and ```
SessionConfig
```

are currently not populated. The default values in ```
Assets/QuantumUser/Resources
```

 are used. When creating copies of the config files at a different location they need to be manually linked.

![The Menu Object Setup](/docs/img/quantum/v3/tutorials/asteroids/8-menu-object-setup.png)
 The Menu Object Setup.


The Quantum Menu loads the game scene using Unity's scene management. For that it needs to be registered in the build settings. Open the Build Settings and add the ```
Menu
```

scene to the ```
Scenes In Build
```

list. Below it add the ```
AsteroidsGameplay
```

scene.

![Adding the scenes to build settings](/docs/img/quantum/v3/tutorials/asteroids/8-build-settings.png)
Adding the scenes to build settings.
## Online Play

Next, create a game build to test the game online.

First, go to ```
Edit > ProjectSettings > Player > Resolution and Presentation
```

 and change Fullscreen Mode to ```
Windowed
```

. Then go to ```
File > Build Settings
```

 and create a build.

Start the built application and enter play mode in the editor to start 2 instances of the game. Press ```
Quick Play
```

. Each client joins the game scene and is in control of their respective avatars/entities.

![Two clients playing the game](/docs/img/quantum/v3/tutorials/asteroids/8-multiplayer-gameplay.gif)Back to top

- [Overview](#overview)
- [Menu Scene](#menu-scene)
- [Configuring the Menu](#configuring-the-menu)
- [Online Play](#online-play)