# systems

_Source: https://doc.photonengine.com/quantum/current/manual/quantum-ecs/systems_

# Systems (game logic)

## Introduction

Systems are the entry points for all gameplay logic in Quantum.

They are implemented as normal C# classes, although there are a few restrictions for a System to be compliant with the predict/rollback model. Systems must:

- Be stateless: no mutable fields should be declared in systems. All the mutable game data has to be declared in .qtn files which then becomes part of the rollbackable game state inside the `Frame` class;
- Implement and/or use only deterministic libraries and algorithms (Quantum comes with libraries for fixed point math, vector math, physics, random number generation, path finding, etc);

There are a few base system classes one can inherit from:

- `SystemMainThread`: has `OnInit` and `Update` callbacks. Update is executed once per system and, when needing to iterate thgouth entities and their components, the user has to create their own filters. Can also be used to subscribe and react to Quantum signals;
- `SystemMainThreadFilter<Filter>`: works similar to `SystemMainThread`, except that it takes a filter that defines components layout and `Update` is called once for every entity which has all the components defined in the Filter;
- `SystemSignalsOnly`: does _not_ provide an `Update` callback and is commonly used only for reacting to Quantum signals. It has reduced overhead as it does not have task scheduling for it);
- `SystemBase`: advanced uses only, for scheduling parallel jobs into the task graph (not covered in this basic manual).

## Core Systems

The Quantum SDK includes all _Core_ systems in the default `SystemsConfig`.

- `Core.CullingSystem2D()`: Culls entities with a `Transform2D` component in predicted frames.
- `Core.CullingSystem3D()`: Culls entities with a `Transform3D` component in predicted frames.
- `Core.PhysicsSystem2D()`: Runs physics on all entities with a `Transform2D` AND a `PhysicsCollider2D` component.
- `Core.PhysicsSystem3D()`: Runs physics on all entities with a `Transform3D` AND a `PhysicsCollider3D` component.
- `Core.NavigationSystem()`: Used for all NavMesh related components.
- `Core.EntityPrototypeSystem()`: Creates, Materializes and Initializes `EntityPrototypes`.
- `Core.PlayerConnectedSystem()`: Used to trigger the `ISignalOnPlayerConnected` and `ISignalOnPlayerDisconnected` signals.
- `Core.DebugCommand.CreateSystem()`: Used by the state inspector to send data to instantiate/remove/modify entities on the fly ( _Only available in the Editor!_).

All systems are included by default for the user's convenience. Core systems can be selectively added/removed based on the game's required functionalities; e.g. only keep the `PhysicsSystem2D` or `PhysicsSystem3D` based on what the game needs.

## Basic Systems

On Unity, it is possible to create Quantum systems using script templates using the right-click menu:

![System Templates](/docs/img/quantum/v3/manual/quantum-system-templates.png)

The corresponding code snippets generated are:

**System**

C#

```csharp
namespace Quantum {
  using Photon.Deterministic;
  using UnityEngine.Scripting;
  [Preserve]
  public unsafe class NewQuantumSystem : SystemMainThread {
    public override void Update(Frame frame) {
    }
  }
}

```

Overridable API:

- `OnInit(Frame frame)`;
- `Update(Frame frame)`;
- `OnDisabled(Frame frame)`/`OnEnabled(Frame frame)`
- `StartEnabled`;

**System Filter**

C#

```csharp
namespace Quantum {
  using Photon.Deterministic;
  using UnityEngine.Scripting;
  [Preserve]
  public unsafe class NewQuantumSystem : SystemMainThreadFilter<NewQuantumSystem.Filter> {
    public override void Update(Frame frame, ref Filter filter) {
    }
    public struct Filter {
      public EntityRef Entity;
    }
  }
}

```

Overridable API is the same `SystemMainThread`, plus:

- `Any`;
- `Without`;

**System Signals Only**

C#

```csharp
namespace Quantum {
  using Photon.Deterministic;
  using UnityEngine.Scripting;
  [Preserve]
  public unsafe class NewQuantumSystem : SystemSignalsOnly {
  }
}

```

Overridable API is the same `SystemMainThread`, except for `Update`;

These are some of the main callbacks that can be overridden in a System class:

- `OnInit`: executed only once, when game starts. Commonly used to setup initial game data;
- `Update`: used to advance the game state;
- `OnDisabled(Frame frame)` and `OnEnabled(Frame frame)`: called when a system is directly disabled/enabled or when a parent system state is toggled;
- `UseCulling` defines if the System should excluded culled entities.

**PS:** it is mandatory for any Quantum system to use the Attribute `\[UnityEngine.Scripting.Preserve\]`.

Notice that all available callbacks include an instance of `Frame`. The Frame class is the container for all the mutable and static game state data, including entities, physics, navigation and others like immutable asset objects (which will be covered in a separate chapter).

The reason for this is that Systems must be _stateless_ to comply with Quantum's predict/rollback model. Quantum only guarantees determinism if all (mutable) game state data is fully contained in the Frame instance.

It is valid to create read-only constants or private methods (that should receive all need data as parameters).

The following code snippet shows some basic examples of valid and not valid (violating the stateless requirement) in a System:

C#

```csharp
namespace Quantum
{
  public unsafe class MySystem : SystemMainThread
  {
    // This is ok
    private const int _readOnlyData = 10;
    // This is NOT ok (this data will not be rolled back, so it would lead to instant drifts between game clients during rollbacks)
    private int _mutableData = 10;
    public override void Update(Frame frame)
    {
        // it is ok to use a constant to compute something here
        var temporaryData = _readOnlyData + 5;
        // it is NOT ok to modify transient data that lives outside of the Frame object:
        _transientData = 5;
    }
  }
}

```

## SystemsConfig

In Quantum 3, the way systems configuration is handled has changed. Instead of embedding configurations directly within the code, configuration is encapsulated within an asset named `SystemsConfig`.

This config is passed into the `RuntimeConfig` and Quantum will automatically instantiate the requested systems.

Notice that Quantum includes a few pre-built Systems (entry point for the physics engine updates, navmesh and entity prototype instantiations).

To guarantee determinism, the order in which Systems are inserted will be the order in which all callbacks will be executed by the simulator on all clients. So, to control the sequence in which your updates occur, just insert your custom systems in the desired order.

### Creating a new SystemsConfig

A `SystemsConfig` is a normal Quantum asset. Meaning, you can create a new one by right clicking the project window -> Quantum -> SystemsConfig.

The asset has a serialized list of systems. You can interact with it like any normal unity list.

![Systems Config](/docs/img/quantum/v3/manual/config-files/systems-config.png)### Activating and Deactivating Systems

All injected systems are active by default, but it is possible to control their status in runtime by calling these generic functions from any place in the simulation (they are available in the Frame object):

C#

```csharp
public override void OnInit(Frame frame)
{
  // deactivates MySystem, so no updates (or signals) are called in it
  frame.SystemDisable<MySystem>();
  // (re)activates MySystem
  frame.SystemEnable<MySystem>();
  // possible to query if a System is currently enabled
  var enabled = frame.SystemIsEnabled<MySystem>();
}

```

Any System can deactivate (and re-activate) another System, so a common pattern is to have a main controller system that manages the active/inactive lifecycle of more specialized Systems using a simple state machine (one example is to have an in-game lobby first, with a countdown to gameplay, then normal gameplay, and finally a score state).

To make a system start disabled by default override this property:

C#

```csharp
public override bool StartEnabled => false;

```

### System Groups

Systems can be grouped, which allows them to be enabled and disabled together.

Select the `SystemsConfig` asset, add a new system of type `SystemGroup`, then append child systems to it.

![System Group](/docs/img/quantum/v3/manual/ecs/system-setup-groups.png)

**N.B.:** The `Frame.SystemEnable<T>()` and `Frame.SystemDisable<T>()` methods identify systems by type; thus if there are to be several system groups, they each need their own implementation to allow enabling / disabling multiple system groups independently. In this case, it is possible to declare a new system group type as shown below, which can then be used in the systems config asset.

C#

```csharp
namespace Quantum
{
  public class MySystemGroup : SystemGroup
  {
  }
}

```

## Entity Lifecycle API

This section uses the direct API methods for entity creation and composition. Please refer to the chapter on entity prototypes for the the data-driven approach.

To create a new entity instance, just use this (method returns an EntityRef):

C#

```csharp
var e = frame.Create();

```

Entities do not have pre-defined components any more, to add a Transform3D and a PhysicsCollider3D to this entity, just type:

C#

```csharp
var t = Transform3D.Create();
frame.Set(e, t);
var c =  PhysicsCollider3D.Create(f, Shape3D.CreateSphere(1));
frame.Set(e, c);

```

These two methods are also useful:

C#

```csharp
// destroys the entity, including any component that was added to it.
frame.Destroy(e);
// checks if an EntityRef is still valid (good for when you store it as a reference inside other components):
if (frame.Exists(e)) {
  // safe to do stuff, Get/Set components, etc
}

```

Also possible to check dynamically if an entity contains a certain component type, and get a pointer to the component data directly from frame:

C#

```csharp
if (frame.Has<Transform3D>(e)) {
    var t = frame.Unsafe.GetPointer<Transform3D>(e);
}

```

With ComponentSet, you can do a single check if an entity has multiple components:

C#

```csharp
var components = ComponentSet.Create<CharacterController3D, PhysicsBody3D>();
if (frame.Has(e, components)) {
  // do something
}

```

Removing components dynamically is as easy as:

C#

```csharp
frame.Remove<Transform3D>(e);

```

### The EntityRef Type

Quantum's rollback model maintains a variable sized frame buffer; in other words several copies of the game state data (defined from the DSL) are kept in memory blocks at separate locations. This means any pointer to either an entity, component or struct is only valid within a single Frame object (updates, etc).

Entity refs are safe-to-keep references to entities (temporarily replacing pointers) which work across frames, as long as the entity in question still exists. Entity refs contain the following data internally:

- Entity index: entity slot, from the DSL-defined maximum number for the specific type;
- Entity version number: used to render old entity refs obsolete when an entity instance is destroyed and the slot can be reused for a new one.

### Filters

Quantum does not have _entity types_. In the sparse-set ECS memory model, entities are indexes to a collection of components; the _EntityRef_ type holds some additional information such as versioning. These collections are kept in dynamically allocated sparse sets.

Therefore, instead of iterating over a collection of entities, filters are used to create a set of components the system will work on.

C#

```csharp
public unsafe class MySystem : SystemMainThread
{
    public override void Update(Frame frame)
    {
        var filtered = rame.Filter<Transform3D, PhysicsBody3D>();
        while (filtered.Next(out var e, out var t, out var b)) {
          t.Position += FPVector3.Forward * frame.DeltaTime;
          frame.Set(e, t);
        }
    }
}

```

For a comprehensive view on how filters are used, please refer to the _Components_ page.

### Pre-Built Assets and Config classes

Quantum contains a few pre-built data assets that are always passed into Systems through the Frame object.

These are the most important pre-built asset objects (from Quantum's Asset DB):

- `Map` and `NavMesh`: data about the playable area, static physics colliders, navigation meshes, etc... . Custom player data can be added from a data asset slot (will be covered in the data assets chapter);
- `SimulationConfig`: general configuration data for physics engine, navmesh system, etc.
- default `PhysicsMaterial` and `agent configs` (KCC, navmesh, etc):

The following snippets show how to access current Map and NavMesh instances from the Frame object:

C#

```csharp
// Map is the container for several static data, such as navmeshes, etc
Map map = f.Map;
var navmesh = map.NavMeshes["MyNavmesh"];

```

### Assets Database

All Quantum data assets are available inside Systems through the database API from the Frame. Find [here](/quantum/current/manual/assets/assets-simulation) more information about assets in the simulation. Find [here](/quantum/current/manual/assets/assets-unity) more information handling Quantum assets in the view (Unity editor) and [here](/quantum/current/manual/assets/extending-assets) details about extending assets with view-specific data.

## Signals

As explained in the previous chapter, signals are function signatures used to generate a publisher/subscriber API for inter-systems communication.

The following example in a DSL file (from the previous chapter):

C#

```csharp
signal OnDamage(FP damage, entity_ref entity);

```

Would lead to this trigger signal being generated on the Frame class (f variable), which can be called from "publisher" Systems:

C#

```csharp
// any System can trigger the generated signal, not leading to coupling with a specific implementation
f.Signals.OnDamage(10, entity)

```

A "subscriber" System would implement the generated "ISignalOnDamage" interface, which would look like this:

C#

```csharp
namespace Quantum
{
  class CallbacksSystem : SystemSignalsOnly, ISignalOnDamage
  {
    public void OnDamage(Frame frame, FP damage, EntityRef entity)
    {
      // this will be called everytime any other system calls the OnDamage signal
    }
  }
}

```

Notice signals always include the Frame object as the first parameter, as this is normally needed to do anything useful to the game state.

### Generated and Pre-Built Signals

Besides explicit signals defined directly in the DSL, Quantum also includes some pre-built ("raw" physics collision callbacks, for example) and generated ones based on the entity definitions (entity-type-specific create/destroy callbacks).

The collision callback signals will be covered in the specific chapter about the physics engine, so here's a brief description of other pre-built signals:

- `ISignalOnPlayerAdded`: called when a game client sends an instance of RuntimePlayer to server (and the data is confirmed/attached to one tick).
- `ISignalOnComponentAdded<T>`, `ISignalOnComponentRemoved<T>`: called when a component type T is added/removed to/from an entity.

## Triggering Events

Similar to what happens to signals, the entry point for triggering events is the Frame object, and each (concrete) event will result in a specific generated function (with the event data as the parameters).

C#

```csharp
// taking this DSL event definition as a basis
event TriggerSound
{
    FPVector2 Position;
    FP Volume;
}

```

This can be called from a System to trigger an instance of this event (processing it from Unity will be covered on the chapter about the bootstrap project):

C#

```csharp
// any System can trigger the generated events (FP._0_5 means fixed point value for 0.5)
f.Events.TriggerSound(FPVector2.Zero, FP._0_50);

```

Important to reinforce that events MUST NOT be used to implement gameplay itself (as the callbacks on the Unity side are not deterministic). Events are just a one-way fine-grained API to communicate the rendering engine of detailed game state updates, so the visuals, sound and any UI-related object can be updated on Unity.

## Extra Frame API Items

The Frame class also contains entry points for several other deterministic parts of the API that need to be treated as transient data (so rolled back when needed).

The following snippet shows the most important ones:

C#

```csharp
// RNG is a pointer.
// Next gives a random FP between 0 and 1.
// There are also bound options for both FP and int
f.RNG->Next();
// any property defined in the global {} scope in the DSL files is accessed through the Global pointer
var d = f.Global->DeltaTime;
// input from a player is referenced by its index (i is a pointer to the DSL defined Input struct)
var i = f.GetPlayerInput(0);

```

## Optimization By Scheduling

To optimize systems identified as performance hotspots a simple modulo-based entity scheduling can help. Using this only a subset of entities are updated while iterating through them each tick.

C#

```csharp
public override void Update(Frame frame) {
  foreach (var (entity, c) in f.GetComponentIterator<Component>()) {
    const int schedulePeriod = 5;
    if (entity.Index % schedulePeriod == frame.Number % schedulePeriod) {
      // it is time to update this entity
    }
}

```

Choosing a `schedulePeriod` of `5` will make the entity only be updated every 5th tick. Choosing `2` would mean every other tick.

This way the total number of updates is significantly reduced. To avoid updating all entities in **one** tick adding `entity.Index` will make the load be spread over multiple frames.

Deferring the entity update like this has requirements on the user code:

- The deferred update code has to be able to handle different delta times.
- The entity lazy "responsiveness" may be visually noticeable.
- Using `entity.Index` may add to the laziness because new information is processed sooner or later for different entities.

The [Quantum Navigation system](/quantum/current/manual/navigation/workflow-agents#update_interval "Quantum Navigation system") has this feature build-in.

Back to top

- [Introduction](#introduction)
- [Core Systems](#core-systems)
- [Basic Systems](#basic-systems)
- [SystemsConfig](#systemsconfig)

  - [Creating a new SystemsConfig](#creating-a-new-systemsconfig)
  - [Activating and Deactivating Systems](#activating-and-deactivating-systems)
  - [System Groups](#system-groups)

- [Entity Lifecycle API](#entity-lifecycle-api)

  - [The EntityRef Type](#the-entityref-type)
  - [Filters](#filters)
  - [Pre-Built Assets and Config classes](#pre-built-assets-and-config-classes)
  - [Assets Database](#assets-database)

- [Signals](#signals)

  - [Generated and Pre-Built Signals](#generated-and-pre-built-signals)

- [Triggering Events](#triggering-events)
- [Extra Frame API Items](#extra-frame-api-items)
- [Optimization By Scheduling](#optimization-by-scheduling)