# local-player-view

_Source: https://doc.photonengine.com/quantum/current/manual/player/local-player-view_

# Local Player View

## Introduction

It is often useful to get, in Unity scripts, a reference to the game object which represents the view of the local player (or _players_, as there can be multiple of them).

There are several ways to achieve this. This document introduces the core concepts to help you understand the available approaches.

### Basic simulation setup

The key concept here is the PlayerRef, a unique index assigned to each client that joins the game.

A PlayerRef is independent of entities or entity views--it exists solely to represent a player. By default, entities created at runtime have no inherent connection to a PlayerRef.

How this connection is established varies from game to game. In most cases, a player controls a single entity, but it's entirely possible for a player to control multiple entities.

Since this linking is not provided by default, the first step is to establish this it however prefered.

A very common approach is to define a custom `PlayerLink` component as explained below:

1. Create a Quantum component which is meant to store the `PlayerRef` of the player who controls an entity:

Qtn

```cs
component PlayerLink
{
    PlayerRef PlayerRef;
}

```

2. Populate such field whenever best fit. A common approach is to do so whenever a player joins the game for the first time and an entity is created for them:

C#

```cs
namespace Quantum
{
  using UnityEngine.Scripting;
  [Preserve]
  public unsafe class PlayerJoiningSystem : SystemSignalsOnly, ISignalOnPlayerAdded
  {
    public void OnPlayerAdded(Frame frame, PlayerRef player, bool firstTime)
    {
      // Get the player data in order to access its PlayerAvatar
      var playerData = frame.GetPlayerData(player);
      // Create the character entity
      var character = frame.Create(playerData.PlayerAvatar);
      // Get the PlayerLink component from the entity
      var playerLink = frame.Unsafe.GetPointer<PlayerLink>(character);
      // OR, add the playerLink "frame.Add<PlayerLink>(character, out var playerLink);"
      // Assign the player ref to the component
      playerLink->PlayerRef = player;
    }
  }
}

```

### Polling the Player Link on the view side

Since the objective is to find the entity view GameObject that represents the local player, it is first necessary to have access to the `EntityRef` that represents the view game objects. This can be done in a multitude of ways, but a common approach is to use a `QuantumEntityViewComponent` to access it easily.

For example, create a view component which, upon activation, gets the PlayerRef stored on the entity and checks if it is the local one, when the view object is activated:

C#

```cs
using Quantum;
public unsafe class CharacterViewComponent : QuantumEntityViewComponent
{
  public override void OnActivate(Frame frame)
  {
    if (frame.TryGet<PlayerLink>(EntityRef, out var playerLink))
    {
      var isLocalPlayer = Game.PlayerIsLocal(playerLink.PlayerRef);
      // Add any desired logic based on this
    }
  }
}

```

### Storing the reference to the local entity view

It is possible to then cache the local entity view for ease of access by other scripts. One alternative is to right away store a reference to the `QuantumEntityView`. Or just perform some logic immediately, such as activating game UI that should only be active for local players, etc.

**Storing it in a Quantum View Context**

Create a context MonoBehaviour with a field in which the local entity view will be stored:

C#

```cs
namespace Quantum
{
  public class GameContext : QuantumMonoBehaviour, IQuantumViewContext
  {
    public QuantumEntityView LocalEntityView;
  }
}

```

Have an entity view component that implements the base class with the context type and store the reference similarly:

C#

```cs
using Quantum;
public unsafe class CharacterViewComponent : QuantumEntityViewComponent<GameContext>
{
  public override void OnActivate(Frame frame)
  {
    if (frame.TryGet<PlayerLink>(EntityRef, out var playerLink))
    {
      var isLocalPlayer = Game.PlayerIsLocal(playerLink.PlayerRef);
      if (isLocalPlayer == true)
        {
            // Store the local player reference so it can be easily accessed later
            ViewContext.LocalEntityView = EntityView;
        }
    }
  }
  public override void OnDeactivate()
  {
    // Make sure to also clear the reference when, for example, the entity view is deactivated
    ViewContext.LocalEntityView = null;
  }
}

```

Add the `GameContext` component to a game object child of the `QuantumEntityViewUpdater` game object, or directly add the component to it.

Now any other view components just needs also have access to the context, for example:

C#

```cs
public class MyCustomViewComponent : QuantumEntityViewComponent<GameContex>
{
    public override void OnActivate(Frame frame)
    {
        var localPlayerView = ViewContext.LocalEntityView;
    }
}

```

**Storing a static reference to it.**

A quicker/simpler approach is to save a reference to it in a singleton. Of course, the caveats of having a singleton applies and this should be used with care.

C#

```cs
using Quantum;
public unsafe class CharacterViewComponent : QuantumEntityViewComponent
{
  public override void OnActivate(Frame frame)
  {
    if (frame.TryGet<PlayerLink>(EntityRef, out var playerLink))
    {
      var isLocalPlayer = Game.PlayerIsLocal(playerLink.PlayerRef);
      if (isLocalPlayer == true)
        {
            // Store the local player reference so it can be easily accessed later
            LocalPlayerViewReference.Instance.LocalEntityView = EntityView;
        }
    }
  }
  public override void OnDeactivate()
  {
    // Make sure to also clear the reference when, for example, the entity view is deactivated
    LocalPlayerViewReference.Instance.LocalEntityView = null;
  }
}

```

### Server Replays

The `PlayerIsLocal` work for most of cases, except when running replays in which the definition of which are the local players will not exist.

For the use case of replays, the user must define what is the PlayerRef that is being observed based on their own replay logic in order for this to properly work.

For more information, check the Replays documentation.

Back to top

- [Introduction](#introduction)
  - [Basic simulation setup](#basic-simulation-setup)
  - [Polling the Player Link on the view side](#polling-the-player-link-on-the-view-side)
  - [Storing the reference to the local entity view](#storing-the-reference-to-the-local-entity-view)
  - [Server Replays](#server-replays)