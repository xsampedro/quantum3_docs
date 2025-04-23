# player-replacement

_Source: https://doc.photonengine.com/quantum/current/manual/player/player-replacement_

# Replacement Bots

## Introduction

It is often useful to let AI get the control over players' characters in one of two situations:

1. To replace players who got disconnected from the game during an ongoing match. This helps creating fairer matches as Bots can help while the player tries to reconnect to the game, or even to compensate for a player who did a rage quit.
2. To fill a room with fake players when the minimum amount of necessary players to start a game session has not been reached. This is particularly important during the early stages of a game's release cycle when the playerbase is still small.

## The Setup

In Quantum, the AI logic for such a feature is executed locally by every client's machine; meaning there is no concept of a "master client who simulates Bots input".

Although it is usually simple to run an AI to control some game entities, doing so is game specific. Of course, the complexity of the AI implementation itself can range from very simple to very complex.

The easiest way to start doing this is to signalize whether an entity is, at any point in time, controlled by AI. This can be achieved in different ways:

- Adding a "flag component", like a ```
component AI {}
```

which is added / removed from entities when needed. A system can then iterate over every entity which has an ```
AI
```

component to perform the controlling logic;
- Using a Boolean in a component to turn on / off the AI controls, like ```
component MyCharacter { bool ControlledByAI; }
```

;
- Adding more AI-specific components with lots of extra data, such as the Bot SDK's agent components (HFSM, BT, etc), or a custom one;

Now, when should it be done? It depends on the chosen use cases mentioned before; these will be explored throughout the next sections.

### Replacing a real player during a game match

This can be achieved by activating the ```
PlayerConnectedSystem
```

and then reacting to the ```
ISignalOnPlayerConnected
```

and ```
ISignalOnPlayerDisconnected
```

signals. The system and its accompanying signals are explained in the [Player documentation](/quantum/current/manual/player/player).

Once the player disconnects: find the entity or entities controlled by that player, and setup the AI for them as explained above.

When the player connects again, check if there are entities which were controlled by that player and remove the AI setup so the player takes back the control from the AI.

It is worth mentioning the system above uses the ```
PlayerInputFlags
```

in order to work, which can also be used independently of the ```
PlayerConnectedSystem
```

, if desired. Find more information about [Player Input Flags HERE](/quantum/current/manual/player/input-flags).

### Filling a room with Bots

In this case, _there are no actual players involved_. In other words, entities are created which are never meant to be controlled by an actual person.

Since no players are involved, no connectivity logic is needed either. It is possible for the custom game logic to fill the room with entities, like in this _sample algorithm / snippet_:

- In a Quantum system, wait an interval of time after the game started so Players have time to connect and send their player data;
- When Players arrive, using the ```
OnPlayerDataSet
```

callback, save in the game state (e.g. in a variable in frame.Global) the amount of players who have successfully connected and joined to the game;
- After the interval, subtract this amount from the expected player count from the frame API, like so:

```
int fillAmount = frame.PlayerCount - frame.Global->ConnectedPlayersCount;
```

- Use the result to perform a ```
for
```

loop where the bot entities will be created:

C#

```csharp
for(int i = 0; i < fillAmount; i++)
{
// Create a new Entity here
// Setup it as a Bot as explained earlier on this document
}

```

The snippet above is very simple and should be adjusted to the game's and game design's requirements; for instance, it may be useful to assign special information to the bot entities such as faked player information, team data, etc.

## Selecting a Bot to create

Depending on the game type, it can be useful to create a new Bot based on some already known data. For example pick a Bot for a character which was not yet chosen or with varying levels of difficulty.

**TIP:** The ```
RuntimeConfig
```

asset can hold some references to Entity Prototypes (i.e ```
AssetRefEntityPrototype
```

) so you can reference a variety of characters to pick from. Alternatively, there could be a single type of character with a reference to different AI assets to control it (e.g. different State Machines based on the difficulty level).

## Players and Bots Architecture

Characters are controlled by Quantum Systems. These systems usually know how to read player inputs to change their character's game state, such as moving them, rotating them and triggering attacks.

Now, controlling these same characters with AI logic can be done in multiple ways. Here is an example _code architecture_ which usually works well:

1. Players naturally have an Input which can be polled with ```
frame.GetPlayerInput(playerIndex)
```

, which returns you a pointer to a struct of type ```
Input
```

;
2. Bots can also have the same struct in a custom component - ```
component Bot { Input Input }
```

    \- , and the AI logic itself might be used just to fill the data inside of it;
3. Fill the input data _before any character system runs_. This way, if systems know how to get the input regardless of who is filling it, then no additional special checks in the systems are needed to know if that entity is a player or a Bot;
4. This means that the AI system might (almost) never directly influence the entity state, but rather it generates fake inputs based on its decision making logic.

The advantage of using such architecture is a clear separation of: Inputs \| Players and Bots \| Characters, by providing decoupled systems.

Remember: _this is just a suggestion_. This architecture is not at all mandatory and the same result can be achieve in many other ways.

Here is a visualisation of this strategy used in the [Twin Stick Shooter Sample](/quantum/current/game-samples/twin-stick-shooter):

![Input Polling](/docs/img/quantum/v2/game-samples/twin-stick-shooter/Polling Input.jpg)

**Attention:** the Twin Stick Shooter Sample is an **ADVANCED** sample, so an understanding of Quantum's basics is necessary prior to analysing it.

## AI with Bot SDK

In case you just started thinking about the AI for your project, you may want to look into Quantum's Bot SDK. The Bot SDK is a set of Editor tools and Quantum code which supports the creation of AI Agents such as State Machines, Behaviour Trees and other commonly used solutions.

Find out more about Bot SDK [HERE](/quantum/current/addons/bot-sdk/overview).

Back to top

- [Introduction](#introduction)
- [The Setup](#the-setup)

  - [Replacing a real player during a game match](#replacing-a-real-player-during-a-game-match)
  - [Filling a room with Bots](#filling-a-room-with-bots)

- [Selecting a Bot to create](#selecting-a-bot-to-create)
- [Players and Bots Architecture](#players-and-bots-architecture)
- [AI with Bot SDK](#ai-with-bot-sdk)