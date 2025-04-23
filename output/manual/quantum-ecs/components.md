# components

_Source: https://doc.photonengine.com/quantum/current/manual/quantum-ecs/components_

# Components

## Introduction

Components are special structs that can be attached to entities, and used for filtering them (iterating only a subset of the active entities based on its attached components).

Aside from custom components, Quantum comes with several pre-built ones:

- Transform2D/Transform3D: position and rotation using Fixed Point (FP) values;
- PhysicsCollider, PhysicsBody, PhysicsCallbacks, PhysicsJoints (2D/3D): used by Quantum's stateless physics engines;
- PathFinderAgent, SteeringAgent, AvoidanceAgent, AvoidanceObstacle: navmesh-based path finding and movement.

## Component

This is a basic example definition of a component in the DSL:

Qtn

```cs
component Action
{
    FP Cooldown;
    FP Power;
}

```

Labeling them as components (like above), instead of structs, will generate the appropriate code structure (marker interface, id property, etc). Once compiled, these will also be available in the Unity Editor for use with the Entity Prototype. In the editor, custom components are named _Entity Component ComponentName_.

The API to work on components is presented via the _Frame_ class.

You have the option of working on copies on the components, or on them components via pointers. To distinguish between the access type, the API for working on copies is accessible directly via`Frame` and the API for accessing pointers is available under `Frame.Unsafe` \- as the latter modifies the memory.

The most basic functions you will require to add, get and set components are the functions of the same name.

`Add<T>` is used to add a component to an entity. Each entity can only carry one copy of a certain component. To aid you in debugging, `Add<T>` returns an _AddResult_ Enum.

C#

```csharp
public enum AddResult {
    EntityDoesNotExist     = 0, // The EntityRef passed in is invalid.
    ComponentAlreadyExists = 1, // The Entity in question already has this component attached to it.
    ComponentAdded         = 2  // The component was successfully added to the entity.
}

```

Once an entity has a component, you can retrieve it with `Get<T>`. This will return a copy of the component value. Since you are working on a copy, you will need to save the modified values on the component using `Set<T>`. Similarly to the _Add_ method, it returns a _SetResult_ which can be used to verify the operation's result or react to it.

C#

```csharp
public enum SetResult {
    EntityDoesNotExist = 0, // The EntityRef passed in is invalid.
    ComponentUpdated   = 1, // The component values were successfully updated.
    ComponentAdded     = 2  // The Entity did not have a component of this type yet, so it was added with the new values.
}

```

For example if you were to set the starting value of a health component, you would do the following:

C#

```csharp
private void SetHealth(Frame frame, EntityRef entity, FP value){
    var health = frame.Get<Health>(entity);
    health.Value = value;
    frame.Set(entity, health);
}

```

This table recaps the methods already presented and the others offered to you to manipulate components and their values are:

| Method | Return | Additional Info |
| --- | --- | --- |
| Add<T>(EntityRef entityRef) | `AddResult` enum, see above. | Allows an invalid `EntityRef`. |
| Get<T>(EntityRef entityRef) | A copy of `T` with the current values. | Does **not** allow an invalid `EntityRef`.<br> Throws an exception if the component `T` is not present on the entity. |
| Set<T>(EntityRef entityRef) | `SetResult` enum, see above. | Allows an invalid `EntityRef`. |
| Has<T>(EntityRef entityRef) | `true` if the entity exists and the component is attached.<br> <br>`false` if the entity does not exist or the component is not attached. | Allows invalid `EntityRef` and component to not exist. |
| TryGet<T>(EntityRef entityRef, out T value) | `true` if the entity exists and component is attached.<br> <br>`false` if the entity does not exist, or component not attached to it. | Allows an invalid `EntityRef`. |
| TryGetComponentSet(EntityRef entityRef, <br>out ComponentSet componentSet) | `true` = entity exists and all components of the components are attached.<br> <br>`false` = entity does not exist, or one or more components of the set are<br>not attached. | Allows an invalid `EntityRef`. |
| Remove<T>(EntityRef entityRef) | No return value. <br>Will remove component if the entity exists and carries the component. <br>Otherwise does nothing. | Allows an invalid `EntityRef`. |

To facilitate working on components directly and avoid the -small- overhead from using Get/Set, `Frame.Unsafe` offers unsafe versions of Get and TryGet (see table below).

| Method | Return | Additional Info |
| --- | --- | --- |
| GetPointer<T>(EntityRef entityRef) | `T`\* | Does NOT allow invalid entity ref.<br>Throws an exception if the component `T` is not present on the entity. |
| TryGetPointer<T>(EntityRef entityRef<br>out T\* value) | `true` if the entity exists and component is attached to it.<br> <br>`false` if the entity does not exist, or component not attached to it. | Allows an invalid `EntityRef`. |
| AddOrGet<T>(EntityRef entityRef, out <T>\* result) | `true` if the entity exists and the component is attached or has been attached.<br> <br>`false` if the entity does not exist. | Allows an invalid `EntityRef`. |

Monolithic structs should be avoided and split up into multiple structs. They can causes `bracket nesting level exceeded maximum` errors when compiling IL2CPP.

## Singleton Component

A _Singleton Component_ is a special type of component of which only one can exist at any given time. There can ever only be one instance of a specific T singleton component, on _any_ entity in the entire game state - this is enforced deep in the core of the ECS data buffers. This is strictly enforced by Quantum.

A custom _Singleton Component_ can be defined in the DSL using `singleton component`.

C#

```csharp
singleton component MySingleton{
    FP Foo;
}

```

Singletons inherit an interface called `IComponentSingleton` which itself inherits from `IComponent`. It can therefore do all the common things you would expect from regular components:

- It can be attached to any entity.
- It can be managed with all the regular safe & unsafe methods (e.g. Get, Set, TryGetPointer, etc...).
- It can be put on entity prototypes via the Unity Editor, or instantiated in code on an entity.

In addition to the regular component related methods, there are several special methods dedicated to singletons. Just like for regular components, the methods are separated in _Safe_ and _Unsafe_ based on whether they return a value type or a pointer.

| Method | Return | Additional Info |
| --- | --- | --- |
| API - Frame |
| SetSingleton<T> (T component, <br>EntityRef optionalAddTarget = default) | void | Sets a singleton IF the singleton does not exist.<br>\-\-\-----<br>EntityRef (optional), specifies which entity to add it to. <br>IF none is given, a new entity will be created to add the singleton to. |
| GetSingleton<T>() | T | Throws exception if singleton does not exist.<br>No entity ref is needed, it will find that automatically. |
| TryGetSingleton<T>(out T component) | bool<br>true = singleton exists<br>false = singleton does NOT exist | Does NOT throw an exception if singleton does not exist.<br>No entity ref is needed, it will find that automatically. |
| GetOrAddSingleton<T>(EntityRef optionalAddTarget = default) | T | Gets a singleton and returns it.<br>IF the singleton does not exist, it will be created like in SetSingleton.<br>\-\-\---<br>EntityRef (optional), specifies which entity to add it to if it has to be created.<br>A new entity will be created to add the singleton to if no EntityRef is passed in. |
| GetSingletonEntityRef<T>() | EntityRef | Returns the entity which currently holds the singleton.<br>Throws if the singleton does not exist. |
| TryGetSingletonEntityRef<T>(out EntityRef entityRef) | bool<br>true = singleton exists.<br>false = singleton does NOT exist. | Get the entity which currently holds the singleton.Does NOT throw if the single does not exist. |
| API - Frame.Unsafe |
| Unsafe.GetPointerSingleton<T>() | T\* | Gets a singleton pointer.<br>Throws exception if it does not exist. |
| TryGetPointerSingleton<T>(out T\* component) | bool<br>true = singleton exists.<br>false = singleton does NOT exist. | Gets a singleton pointer. |
| GetOrAddSingletonPointer<T>(EntityRef optionalAddTarget = default) | T\* | Gets or Adds a singleton and returns it.<br>IF the singleton does not exist, it will be created.<br>\-\-\---<br>EntityRef (optional), specifies which entity to add it to if it has to be created.<br>A new entity will be created to add the singleton to if no EntityRef is passed in. |

## ComponentTypeRef

The `ComponentTypeRef` struct provides a way for referencing a component by its type during runtime. This is useful if you are dynamically adding a component via polymorphism.

C#

```csharp
// set in an asset or prototype for example
ComponentTypeRef componentTypeRef;

var componentIndex = ComponentTypeId.GetComponentIndex(componentTypeRef);

frame.Add(entityRef, componentIndex);

```

## Adding Functionality

Since components are special structs, you can extend them with custom methods by writing a _partial_ struct definition in a C# file.

For example, if we could extend our Action component from before as follows:

C#

```csharp
namespace Quantum
{
    public partial struct Action
    {
        public void UpdateCooldown(FP deltaTime){
            Cooldown -= deltaTime;
        }
    }
}

```

## Reactive Callbacks

There are two component specific reactive callbacks:

- `ISignalOnComponentAdd<T>`: called when a component type T is added to an entity.
- `ISignalOnComponentRemove<T>`: called when a component type T is removed from an entity.

These are particularly useful in case you need to manipulate part of the component when it is added/removed - for instance allocate and deallocate a list in a custom component.

To receive these signals, simply implement them in a system.

## Components Iterators

If you were to require a single component only, _ComponentIterator_ (safe) and _ComponentBlockIterator_ (unsafe) are best suited.

C#

```csharp
foreach (var pair in frame.GetComponentIterator<Transform3D>())
{
    var component = pair.Component;
    component.Position += FPVector3.Forward * frame.DeltaTime;
    frame.Set(pair.Entity, component);
}

```

Component block iterators give you the fastest possible access via pointers.

C#

```csharp
// This syntax returns an EntityComponentPointerPair struct
// which holds the EntityRef of the entity and the requested Component of type T.
foreach (var pair in frame.Unsafe.GetComponentBlockIterator<Transform3D>())
{
    pair.Component->Position += FPVector3.Forward * frame.DeltaTime;
}

// Alternatively, it is possible to use the following syntax to deconstruct the struct
// and get direct access to the EntityRef and the component
foreach (var (entityRef, transform) in frame.Unsafe.GetComponentBlockIterator<Transform3D>())
{
    transform->Position += FPVector3.Forward * frame.DeltaTime;
}

```

## Filters

Filters are a convenient way to filter entities based on a set of components, as well as grabbing only the necessary components required by the system. Filters can be used for both Safe (Get/Set) and Unsafe (pointer) code.

### Generic

To create a filter simply use the **Filter()** API provided by the frame.

C#

```csharp
var filtered = frame.Filter<Transform3D, PhysicsBody3D>();

```

The generic filter can contain up to 8 components.

If you need to more specific by creating _without_ and _any_ **ComponentSet** filters.

C#

```csharp
var without = ComponentSet.Create<CharacterController3D>();
var any = ComponentSet.Create<NavMeshPathFinder, NavMeshSteeringAgent>();
var filtered = frame.Filter<Transform3D, PhysicsBody3D>(without, any);

```

A _ComponentSet_ can hold up to 8 components.

The _ComponentSet_ passed as the _without_ parameter will exclude all entities carrying at least one of the components specified in the set. The _any_ set ensures entities have at least one or more of the specified components; if an entity has none of the components specified, it will be excluded by the filter.

Iterating through the filter is as simple as using a while loop with `filter.Next()`. This will fill in all copies of the components, and the `EntityRef` of the entity they are attached to.

C#

```csharp
while (filtered.Next(out var e, out var t, out var b)) {
  t.Position += FPVector3.Forward * frame.DeltaTime;
  frame.Set(e, t);
}

```

**N.B.:** You are iterating through and working on **copies** of the components. So you need to set the new data back on their respective entity.

The generic filter also offers the possibility to work with component pointers.

C#

```csharp
while (filtered.UnsafeNext(out var e, out var t, out var b)) {
  t->Position += FPVector3.Forward * frame.DeltaTime;
}

```

In this instance you are modifying the components' data directly.

### FilterStruct

In addition to regular filters, you may use the _FilterStruct_ approach.

For this you need to first define a struct with **public** properties for each component type you would like to receive.

Qtn

```cs
struct PlayerFilter
{
    public EntityRef Entity;
    public CharacterController3D* KCC;
    public Health* Health;
    public FP AccumulatedDamage;
}

```

Just like a _ComponentSet_, a _FilterStruct_ can filter up to 8 different component pointers.

**N.B.:** A struct used as a _FilterStruct_ is **required** to have an _EntityRef_ field!

The **component type** members in a _FilterStruct_ **HAVE TO BE** pointers; only those will be filled by the filter. In addition to component pointers, you can also define other variables, however, these will be ignored by the filter and are left to you to manage.

C#

```csharp
var players = f.Unsafe.FilterStruct<PlayerFilter>();
var playerStruct = default(PlayerFilter);

while (players.Next(&playerStruct))
{
    // Do stuff
}

```

`Frame.Unsafe.FilterStruct<T>()` has an overload utilizing the optional ComponentSets _any_ and _without_ to further specify the filter.

### Note on Count

A filter does not know in advance how many entities it will touch and iterate over. This is due to the way filters work in _Sparse-Set_ ECS:

1. the filter finds which among the components provided to it has the least entities associated with it (smaller set to check for intersection); and then,
2. it goes through the set and discards any entity that does not have the other queried components.

Knowing the exact number in advance would require traversing the filter once; as this is an (O(n) operation, it would not be efficient.

## Components Getter

Should you want to get a specific set of components from a _known_ entity, use a filter struct in combination with the `Frame.Unsafe.ComponentGetter`. **N.B.:** This is only available in an unsafe context!

C#

```csharp
public unsafe class MySpecificEntitySystem : SystemMainThread

    struct MyFilter {
        public EntityRef      Entity; // Mandatory member!
        public Transform2D*   Transform2D;
        public PhysicsBody2D* Body;
    }

    public override void Update(Frame frame) {
        MyFilter result = default;

        if (frame.Unsafe.ComponentGetter<MyFilter>().TryGet(frame, frame.Global->MyEntity, &result)) {
            // Do Stuff
        }
    }

```

If this operation has to performed often, you can cache the look-up struct in the system as shown below (100% safe).

C#

```csharp
public unsafe class MySpecificEntitySystem : SystemMainThread

    struct MyFilter {
        public EntityRef      Entity; // Mandatory member!
        public Transform2D*   Transform2D;
        public PhysicsBody2D* Body;
    }

    ComponentGetter<MyFilter> _myFilterGetter;

    public override void OnInit(Frame frame) {
      _myFilterGetter = frame.Unsafe.ComponentGetter<MyFilter>();
    }

    public override void Update(Frame frame) {
      MyFilter result = default;

      if (_myFilterGetter.TryGet(frame, frame.Global->MyEntity, &result)) {
        // Do Stuff
      }
    }

```

## Filtering Strategies

Often times you will be running into a situation where you will have many entities, but you only want a subset of them. Previously we introduced the components and tools available in Quantum to filter them; in this section, we will present some strategies that utilize these.

**N.B.:** The _best_ approach will depend on your own game and its systems. We recommend taking the strategies below as a jumping off point to create a fitting one to your unique situation.

_Note: All terminology used below has been created in-house to encapsulate otherwise wordy concepts._

### Micro-component

Although many entities may be using the same component types, few entities use the same component composition. One way to further specialize their composition is by the use of **micro-components** . **Micro-components** are highly specialized components with data for a specific system or behaviour. Their uniqueness will allow you to create filters that can quickly identify the entities carrying it.

### Flag-component

One common way to identify entities is by adding a **flag-component** to them. In ECS the concept of _flags_ does not exist per-se, nor does Quantum support _entity types_; so what exactly are **flag-components** ? They are components holding little to no data and created for the exclusive purpose of identifying entities.

For instance, in a team based game you could have:

1. a "Team" component with an enum for TeamA and TeamB; or
2. a "TeamA" and "TeamB" component.

Option 1. is helpful when the main purpose is polling the data from the View, while option 2. will enable you to benefit from the filtering performance in the relevant simulation systems.

_Note:_ Sometimes a flag-component are also referred to as tag-component because tagging and flagging entities is used interchangeably.

#### Count

The amount of a components T currently existing in the simulation can be retrieved using `Frame.ComponentCount<T>()`. When used in conjunction with flag components it enables a quick count of, for instance, a certain type of units.

#### Add / Remove

In case you only need to _temporarily_ attach a flag-component or micro-component to an entity, they remain a suitable options as both the `Add` and `Remove` operations are O(1).

### Global Lists

An alternative to flag-components, albeit a "less" ECS-ish one, is to keep global lists in `FrameContext.User.cs`. While this does not necessarily scale if you need to keep track of N teams, it is convenient for sets where subsets are limited.

If you wanted to highlight all players with less than 50% health, you could hold a global list and do the following:

- Have a system at the beginning of the simulation that add/removes entity\_refs to the list;
- Use that same list in all subsequent systems.

**N.B.:** If you only need to identified these types of conditions sporadically, we would advise to dynamically calculated it when needed rather than keeping global lists.

## Max Component Count

By default, the Quantum solution supports the definition of up to 256 different component types.

For user-defined components this number is smaller (236), since the Core DLL already comes with 20 component types pre-defined (Transforms, Colliders, etc).

Although this has proven to be enough for most games, it is possible to increase this maximum count to 512 by adding this compiler define to a QTN file:

```
#pragma max_components 512

```

Increasing the component count can result in an increase of average simulation time for games with high entity count that rely heavily on filtering entities based on their component set, so profiling tests are recommended to measure the performance impact on your specific scenario using the instructions shared in the [Profiling](/quantum/current/manual/profiling "Profiling documentation") page.

## Importing Self-Defined Components

Quantum 3 allows to define components outside of the DSL and import them manually. This is useful if you need to define components in outside DLLs for example. It is **usually not necessary at all** as the the regular path of defining components in the DSL itself is safer.

To import a component, add `import FooComponent;` or `import singleton FooComponent;` to any DSL file.

The component definition itself has to follow a few guidelines before being properly imported.

**PS:** please, be very careful with this definition. Implementing it correctly is very important for the SDK functioning. Declaring proper component SIZE and FieldOffsets is very important.

The requirements are that the definition:

1. Implements `IComponent` interface;
2. Has `const int SIZE` field, which defines the size of the component;
3. Has a `Serialize` method, which is used to serialize the component (see signature below);
4. Has `ComponentChangedDelegate OnAdded` static property (may return null) or static `OnAdded` method matching the delegate signature.
5. Has `ComponentChangedDelegate OnRemoved` static property (may return null) or static `OnRemoved` method matching the delegate signature.

One **safer** alternative is to first define the component in the DSL, copy the generated code of it, then remove it again just so all the important details are handled.

Here is an example of the basic structure needed for a component definition:

C#

```csharp
[StructLayout(LayoutKind.Explicit)]
public unsafe struct Example : IComponent {
  public const int SIZE = sizeof(int);

  public static ComponentChangedDelegate OnAdded;
  public static ComponentChangedDelegate OnRemoved;

  [FieldOffset(0)]
  public int _number;

  public static void Serialize(void* ptr, IDeterministicFrameSerializer serializer) {
    serializer.Stream.Serialize(&((Example*)ptr)->_number);
  }

  public override int GetHashCode() {
    return _number;
  }
}

```

Back to top

- [Introduction](#introduction)
- [Component](#component)
- [Singleton Component](#singleton-component)
- [ComponentTypeRef](#componenttyperef)
- [Adding Functionality](#adding-functionality)
- [Reactive Callbacks](#reactive-callbacks)
- [Components Iterators](#components-iterators)
- [Filters](#filters)

  - [Generic](#generic)
  - [FilterStruct](#filterstruct)
  - [Note on Count](#note-on-count)

- [Components Getter](#components-getter)
- [Filtering Strategies](#filtering-strategies)

  - [Micro-component](#micro-component)
  - [Flag-component](#flag-component)
  - [Global Lists](#global-lists)

- [Max Component Count](#max-component-count)
- [Importing Self-Defined Components](#importing-self-defined-components)