# overview

_Source: https://doc.photonengine.com/quantum/current/addons/bot-sdk/overview_

# Overview

Level

ADVANCED

## Overview

**Bots for Single player, Local & Online Multiplayer**

Bots can be tremendously important for the success of a multiplayer game. Two of the main advantages of having Bots are:

1. **Bots can fill up rooms when there are not enough players connected**. This can greatly improve the experience players in early stages of the a game release as they can play with bots and not depend entirely on the game already having enough players for quickly filling up rooms;
2. **Bots can also be used to replace players who got disconnected**, and players can get back the control of their entities again if they reconnect to the game.

Bot SDK can be a huge time saver, especially when used on early stages of a game's development as it is built with a custom editor in Unity which allows for fast prototyping which also gives game/level designers a lot of power when tweaking gameplay directly within the editor.

### Download Stable

| Version | Release Date | Download |
| --- | --- | --- |
| 3.9.4 | Oct 08, 2025 | [Quantum 3.1 Bot SDK](https://downloads.photonengine.com/download/quantum/photon-quantum-bot-sdk-preview-v3.1.zip?pre=sp) | [Release Notes](/quantum/current/addons/bot-sdk/overview#release-history) |
| 3.9.4 | Oct 08, 2025 | [Quantum 3.0 Bot SDK](https://downloads.photonengine.com/download/quantum/photon-quantum-bot-sdk-stable-v3.0.zip?pre=sp) | [Release Notes](/quantum/current/addons/bot-sdk/overview#release-history) |

## Installation and Migration Notes

To install Bot SDK, simply import the Unity package in the game project. The same steps apply for upgrading the SDK to a newer version.

When migrating a Quantum project from version 3.0 to version 3.1, there will be a step in which the user should fix all the compilation errors before proceeding.

At this point:

1. Close all AI documents;
2. Close the Bot SDK editor window;
3. Download Bot SDK for Quantum 3.1;
4. Import it into the project. All compilation errors specific to the addon should be solved now and you are able to proceed with the next migration steps.

## Introduction

Bot SDK is mainly divided in two parts:

1. The **Visual Editor** created for easily defining and tweaking the Bots behaviours;
2. The **deterministic AI code**, implemented as part of the Quantum simulation, directly integrated with the editor.

It currently supports the following AI models:

1. Hierarchical Finite State Machine ( **HFSM**);
2. Behaviour Tree ( **BT**);
3. Utility Theory ( **UT**);

When using this addon, the choice of the AI technique(s) that will be used is the project is completely up to the user:

- It is possible to design all of the bots using only one technique;
- Or design some bots using HFSM, others using UT, others using BT;

Find out the pros and cons of each AI model on their specific pages.

For understanding the basic addon setup and testing it out, take a look at the [Bot SDK Sample](/quantum/current/technical-samples/bot-sdk-sample).

## Opening the Editor

Go to the Unity's top menu and navigate to `Window > Bot SDK > Open Editor`, which opens the following window:

![Main Window](/docs/img/quantum/v3/addons/bot-sdk/bot-sdk-main-window.png)
Bot SDK initial editor window.
## Release History

### 3.9

#### Bot SDK 3.9.4 (Oct 06, 2025)

General

- Added better support to running multiple AI agents, of any Bot SDK type, referencing the same Entity;
- Improved the API of HFSMManager, UTManager and BTManager to better adapt to the use of AIContext;

Utility Theory

- Fixed issue which would make the UT compiler fail to bake Ranking and Commitment slots;

Behaviour Tree

- Added "ref AIContext" parameter to NodeOnChildCompletedRunning();

#### Bot SDK 3.9.3 (Sept 19, 2025)

Visual Editor:

- Fixed issue in which creating new HFSM documents while no other document was open would break the document;

#### Bot SDK 3.9.2 (Sept 12, 2025)

Visual Editor:

- Fixed issue in which the editor would fail to create SettingsDatabase;
- Fixed issue in which closing the AI document would throw errors if the debugger was disabled;
- Changed the initial position of the first node of each AI document;

Behaviour Tree:

- Fixed issue in which dynamic composite nodes abort would cause its parent composite to also be removed from the dynamics list;

#### Bot SDK 3.9.1 (Aug 08, 2025)

Debugger:

- Fixed issues when running replays;
- Fixed issue which would break the use of Debug Points;
- Fixed issue in which the HFSM transitions debugger would colors in the reversed order;

Blackboard:

- Added null checks when using the main Blackboard API to check if the key provided is valid;

#### Bot SDK 3.9.0 (Jul 17, 2025)

Behaviour Tree

- (Experimental) Added new node `SharedBT` which references a BT AI document, making it possible to create re-usable BT files;

### 3.8

#### Bot SDK 3.8.1 (May 16, 2025)

Behaviour Tree

- Fixed issue in which destroying the BT Agent during its own BT execution would cause null pointer exceptions;

#### Bot SDK 3.8.0 (Apr 17, 2025)

Behaviour Tree:

- \[BREAKING CHANGE\] Removed `BTAgent.AddFPData` and `BTAgent.AddIntData` methods. Delete any call to these methods and the BT agent should still behave the same;
- Added `BTNegate` decorator node, which executes the next Decorator node but negates its `CheckConditions` result;

Visual Editor:

- Fixed issue in which compiling an AI document would always generate all the folders, even when not necessary, e.g when the user does not need the AIConfigs folder;
- Upgrading Bot SDK does not shows up a dialog window informing anymore;

Debugger:

- Added new overload to the `AddToDebugger` method which receives a pointer to the component as a parameter;

Blackboard:

- Fixed issue in which the variables values would not show up;

### 3.7

#### Bot SDK 3.7.8 (Mar 20, 2025)

Circuit Editor

- Added support to registering AI agents on the debugger while the Bot SDK window is closed;
- On the debugger window, added checkbox to filter in/out Utility Theory agents;

#### Bot SDK 3.7.7 (Mar 13, 2025)

Visual Editor

- Fixed issue on the `Duplicate Document` button which would cause the new document and the old one to reference the same XML asset;

#### Bot SDK 3.7.6 (Mar 07, 2025)

Visual Editor

- Fixed issue in which closing and re-opening the Bot SDK window would cause debugger callbacks to be executed more than once, such as wrongly registering multiple AI agents;
- Fixed issue in which closing and re-opening the Bot SDK window would cause runtime errors when opening AI documents due to issues in its cache;

Debugger

- Fixed issue in which runtime errors would be thrown due to issues in DebugPoints debugging;

Behaviour Tree

- Fixed issue on the debugger which aborted branches would be colored as Running nodes;
- Fixed issue in which not all aborted nodes would have the OnAbort callback executed;
- Added black color to the debugger on aborted branches;

#### Bot SDK 3.7.5 (Feb 18, 2025)

State Machine

- Fixed issue which would cause the TimerDecision to always use a fixed time;

CodeGen

- Fixed issue with Constants name used on Bot SDK which would cause conflicts with Quantum's Constants;

#### Bot SDK 3.7.4 (Feb 12, 2025)

Visual Editor

- Fixed issue in which using HFSM compound agents would break the Debugger;

#### Bot SDK 3.7.3 (Jan 16, 2025)

Behaviour Tree

- Fixed issue on debugger of BT Services when their `Interval In Sec` were set to zero;

#### Bot SDK 3.7.2 (Jan 13, 2025)

Behaviour Tree

- Fixed issue on serialization of BTDataIndex fields;

Bot SDK 3.7.1 (Jan 08, 2025)

Behaviour Tree

- Fixed issue on compilation of BTDataIndex slots;

Bot SDK 3.7.0 (Jan 06, 2025)

Compiler

- Quicker compilation time when an AI document is compiled and there is no previous compilation of it found;

State Machine

- Fixed the logic of the hfsm data Time variable to be updated before executing any AI Action;
- Renamed variables: `Time` to `StateTimer` and `Times` to `Timers`;

Behaviour Tree

- Made BTDataIndex 0-based and added a protective layer to warn the use of un-initialized values on it;
- Added overload methods for adding and setting BT nodes data using the BTDataIndex type directly, rather than using an integer;

Blackboard

- On the left panel, renamed text from `Has Default` to `Has Initial Value`;
- On the left panel, renamed text from `Default Value` to `Initial Value`;
- Fixed issue when setting `Has Initial Value` to true, compiling the AI document, and then setting it to false;

Visual Editor

- Separated the XML text of AI Documents into a separate, standalone asset;
- Fixed NRE when debugging HFSM AI documents;
- Fixed issue when opening an AI document, closing the Bot SDK window and opening an AI document again;
- Fixed size of the top bar which was cropping the breadcrumb text buttons;
- Fixed width of the tooltips box which was cropping the texts;

### 3.6

#### Bot SDK 3.6.1 A1 (Nov 25, 2024)

Blackboard

- Added getters and setters overloads for Blackboard variables which receives a hash code instead of the variable name as parameter;
- Fixed issue on the Blackboard which used GetHashCode() that is non deterministic when running against server side simulation;

Folders

- Marked `QuantumBotSDK/Simulation` folder with `QuantumDotNetInclude` tag for it to be included in .Net standalone exports;

#### Bot SDK 3.6.0 A1 (Nov 05, 2024)

Visual Editor

- Performance improvements on all the AI document editors;
- Left side panels can now be collapsed;

### 3.5

#### Bot SDK 3.5.1 A1 (Oct 29, 2024)

Compiler

- Fixed issue in which AI documents compiler would always generate assets with `Override Guid`;

#### Bot SDK 3.5.0 A1 (Sept 30, 2024)

Visual Editor

- Changed top menu to `Tools/Quantum Bot SDK`;
- Fixed issue in which double cliking an AI document asset without having the editor window opened would break the editor;
- Removed usage of Assets/Gizmos folder;

AI Config

- Fixed issue when using the AI Config panel with AssetRef entries;

Utility Theory

- Added debugger to Utility Theory AI documents;
- Fixed issue in which Commitment and Rank fields would not be set to default after disconnecting its slots;

Behaviour Tree

- BT Debugger now shows the progress of BTServices;

### 3.4

#### Bot SDK 3.4.8 A1 (Sept 05, 2024)

Core

- Changed the Quantum.BotSDK.Core dll to make it possible to use Bot SDK in standalone projects;

#### Bot SDK 3.4.7 A1 (Aug 13, 2024)

Visual Editor

- Fixed issue in which compiling an AI document whilst the export folder doesn't exist would break the compilation process;

Behaviour Tree

- Fixed issue in which overriding BTNode.Init() and not calling base.Init() would break the internal BT state;
- Fixed issue in which the BTComposite nodes' Loaded callbacks would try to access possibly null collections;

State Machine

- Fixed issue in which logical decision nodes (AND, OR, NOT) would not cleanup references to nodes baked on the previous compilation;

#### Bot SDK 3.4.6 A1 (July 30, 2024)

Visual Editor

- Fixed sharing violation issue when opening the Bot SDK window for the first time;
- Moved the Bot SDK toolbar options from `Window` to `Tools`;

Bot SDK Core

- On Quantum.BotSDK.Core.dll, removed dependencies on UnityEngine.dll;

AIConfig

- Added log message when trying to set a variable in AIConfig but the key is empty or null;

#### Bot SDK 3.4.5 A1 (July 03, 2024)

AIConfig

- Fixed issue with AIConfig assets in which override versions could not be properly created;

#### Bot SDK 3.4.4 A1 (July 02, 2024)

Visual Editor

- Fixed issue in which slots with AssetRef type would not be visible;

#### Bot SDK 3.4.3 A1 (July 01, 2024)

Visual Editor

- Fixed issue in which the editor would not create nor open AI documents on Unity 2023 due to broken reference to dll;
- Changed default export folder to `Assets/QuantumUser/Resources/DB`;

#### Bot SDK 3.4.2 A1 (Jun 21, 2024)

Systems

- \[BREAKING CHANGE\] Changed the namespaces of the Bot SDK systems;
- Added \[Preserve\] attribute which is a requirement for Quantum 3 RC version;

BT

- Fixed issue when trying to update a BTAgent without passing a pointer to a Blackboard component;
- Fixed issue in the debugger which would throw NRE exceptions when trying to debug a BTAgent;

UT

- Fixed issue in the debugger which would throw NRE exceptions when trying to debug a BTAgent;
- Disabled the basic UT debugger;

Unity Editor

- Removed outdated ways of creating AI documents from the context menu;

#### Bot SDK 3.4.1 A1 (Jun 20, 2024)

Visual Editor

- Fixed issue in which the Bot SDK version file would not be properly created on the first time the addon is used;
- Changed the default export folder for compiled AI assets;

#### Bot SDK 3.4.0 A1 (Jun 19, 2024)

GOAP

- \[BREAKING CHANGE\] Removed the GOAP simulation and editor code;

BT

- \[BREAKING CHANGE\] Changed BTNode.Init() to receive the following parameters: BTParams and AIContext;

Folders

- Moved Bot SDK location to `Assets/Photon/QuantumAddons/QuantumBotSDK`;

### 3.3

#### Bot SDK 3.3.0 A1 (May 22, 2024)

Visual Editor

- Cleaning up the editor selection when changing the graph being edited;
- Comment bubbles: pressing `Enter` does not save the text anymore, it creates a line breank instead (press Esc to save);
- Added possibility to create Response Curve nodes in more types of graph;

Compiler

- Better error handling when a Blackboard or Constant Node is not present on the variables panel;
- Fixed issue when changing nodes types from the editor context menu which would not properly create a new asset with the new type (e.g when converting BT composite nodes);

Utility Theory

- Fixed issue in assets caching which could lead to a desync on late joiners;
- Fixed issue in which every Consideration node would always be created from scratch and not re-use previously baked assets;

HFSM

- Fixed issue in assets caching which could lead to a desync on late joiners;
- Changed from \[HideInInspector\] to \[ExcludeFromPrototype\] on the HFSMData struct;
- Removed unused `Prerequisite` variable from TransitionSets;

Blackboard

- Fixed issue when trying to compile an AI document which referenced a missing asset in a Blackboard variable;
- Added TryGet methods

BT

- Baking BT Service nodes nicknames into the created asset's Label field;
- Improved null pointer check when trying to clear the reactive decorators list;
- \[BREAKING CHANGE\] Renamed BTDecorator's method name from `DryRun` to `CheckConditions`;

### 3.2

#### Bot SDK 3.2.1 A1 (May 08, 2024)

Visual Editor

- Added \[BotSDKTooltip\] attribute which can be used in classes and fields;

#### Bot SDK 3.2.0 A1 (Apr 23, 2024)

BT

- Fixed issue in Dynamic Decorators which would cause multiple nodes to run at the same time;
- Added BTSelectorRandom which randomly picks a single child node to be executed with chances evenly distributed;
- Fixed issue in which BTAgent code would try to de-allocate lists which were not yet allocated;

UT

- Added back the OnComponentAdded/Removed signals to the BotSDKSystem;

HFSM

- Fixed issue with the back-porting and baking of ANY Transitions, which would cause some transitions to not be baked
- Added import to HFSMData in the DSL so it can be used in other qtn files;

Compiler

- Fixed issues when compiling AssetRefs;
- Added better error handling compiling AssetRefs;

Visual Editor

- Fixed issue which would throw NRE when trying to handle Debug Points;
- Changed the referenced Unity dll to version 2019.4.28f1;
- Fixed issue which would cause the SettingsDatabase asset to be created from scratch;

Debugger

- Fixed issues on the debugger enabled/disabled state which would causes issues when trying to run the game with Bot SDK open;
- Improved how debugger hierarchies are handled internally;

### 3.1

#### Bot SDK 3.1.0 A1 (Mar 19, 2024)

Behaviour Tree

- Fixed issue in which BTServices would not be executed;
- Fixed issue in which BTServices would not be properly baked into assets when compiling the AI document;
- Changed BTParams.Frame variable type from FrameThreadSafe to FrameBase;
- Added \[ExcludeFromPrototype\] in most of BTAgent component fields;

Bot SDK Systems

- Fixed issue in the BotSDKTimerSystem in which it would not write the FP current game time;

AIConfig

- Changed its extension methods to use FrameBase instead of Frame;

GOAP

- Added `import singleton component` on Unity, to import GOAPData;

### 3.0

#### Bot SDK 3.0.1 A1 (Mar 14, 2024)

Compiler

- Fixed issue in which AIBlackboard and AIConfig assets would always be generated, even if not used;

#### Bot SDK 3.0.0 A1 (Feb 28, 2024)

Bot SDK Core

- Added a new dll which contains all the Bot SDK core types and logic;
- Adapted all the assets and AssetRefs to function with the new Quantum 3 API;
- Changed how the Bot SDK internal Timer was located, moving it from Frame.Global to a Singleton Component;
- \[BREAKING CHANGE\] All the core methods on Bot SDK now use FrameBase instead of Frame;
- \[BREAKING CHANGE\] AIContext: added a new field called `UserData` which should be used for inserting the custom context struct;

Visual Editor

- Reworked all the document compilers;
- Reworked the document debuggers;
- Adapted other algorithms to use new Quantum 3 API;

Bot SDK Runtime

- Added, in Unity, only the necessary types;
- Adapted all the custom Unity code to work with the new Bot SDK Core dll;

Back to top

- [Overview](#overview)

  - [Download Stable](#download-stable)

- [Installation and Migration Notes](#installation-and-migration-notes)
- [Introduction](#introduction)
- [Opening the Editor](#opening-the-editor)
- [Release History](#release-history)
  - [3.9](#section)
  - [3.8](#section-1)
  - [3.7](#section-2)
  - [3.6](#section-3)
  - [3.5](#section-4)
  - [3.4](#section-5)
  - [3.3](#section-6)
  - [3.2](#section-7)
  - [3.1](#section-8)
  - [3.0](#section-9)