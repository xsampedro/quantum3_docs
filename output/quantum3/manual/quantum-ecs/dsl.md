# dsl

_Source: https://doc.photonengine.com/quantum/current/manual/quantum-ecs/dsl_

# DSL (game state)

## Introduction

Quantum requires components and other runtime game state data types to be declared with its own DSL (domain-specific-language).

These definitions are written into text files with the ```
.qtn
```

 extension. The Quantum compiler will parse them into an AST, and generate partial C# struct definitions for each type (definitions can be split across as many files if needed, the compiler will merge them accordingly).

The goal of the DSL is to abstract away from the developer the complex memory alignment requirements imposed by Quantum's ECS sparse set memory model, required to support the deterministic predict/rollback approach to simulation.

This code generation approach also eliminates the need to write "boiler-plate" code for type serialization (used for snapshots, game saves, killcam replays), checksumming and other functions, like printing/dumping frame data for debugging purposes.

In order to create a new ```
.qtn
```

file to the project, open the context menu on Unity's Project tab and click in ```
Create/Quantum/Qtn
```

, or simply create a new file with the ```
.qtn
```

extension.

## Components

Components are special structs that can be attached to entities, and used for filtering them (iterating only a subset of the active entities based on its attached components). This is a basic example definition of a component:

Qtn

```cs
component Action
{
 FP Cooldown;
 FP Power;
}

```

These will be turned into regular C# structs. Labelling them as components (like above) will generate the appropriate code structure (marker interface, id property, etc).

Aside from custom components, Quantum comes with several pre-built ones:

- Transform2D/Transform3D: position and rotation using Fixed Point (FP) values;
- PhysicsCollider, PhysicsBody, PhysicsCallbacks, PhysicsJoints (2D/3D): used by Quantum's stateless physics engines;
- PathFinderAgent, SteeringAgent, AvoidanceAgent, AvoidanceObstacle: navmesh-based path finding and movement.

## Structs

Structs can be defined in both the DSL and C#.

### DSL Defined

The Quantum DSL also allows the definition of regular structs (just like components, memory alignment, and helper functions will be taken care of):

Qtn

```cs
struct ResourceItem
{
 FP Value;
 FP MaxValue;
 FP RegenRate;
}

```

The generated struct will have fields declared in the same order, but with memory offsets adjusted to avoid paddings and provide optimal packing.

This would let you use the "Resources" struct as a type in all other parts of the DSL, for example using it inside a component definition:

Qtn

```cs
component Resources
{
 ResourceItem Health;
 ResourceItem Strength;
 ResourceItem Mana;
}

```

The generated struct is partial and can be extended in C# if so desired.

### CSharp Defined

You can define structs in C# as well; however, in this case you will have to manually define:

- The memory layout of the struct (using LayoutKind.Explicit)
- Add a const int ```
SIZE
```

to the struct containing its size.
- Implement the ```
Serialize
```

function.

C#

```csharp
\[StructLayout(LayoutKind.Explicit)\]
public struct Foo {
 public const int SIZE = 12; // the size in bytes of all members in bytes.

 \[FieldOffset(0)\]
 public int A;

 \[FieldOffset(4)\]
 public int B;

 \[FieldOffset(8)\]
 public int C;

 public static unsafe void Serialize(void\* ptr, FrameSerializer serializer)
 {
 var foo = (Foo\*)ptr;
 serializer.Stream.Serialize(&foo->A);
 serializer.Stream.Serialize(&foo->B);
 serializer.Stream.Serialize(&foo->C);
 }
}

```

When using C# defined structs in the DSL (e.g. inside components), you will have to manually import the struct definition.

```
import struct Foo(12);

```

**N.B.:** The _import_ does not support constants in the size; you will have to specify the exact numerical value each time.

### Components vs. Structs

An important question is why and when should components be used instead of regular structs (components, in the end, are also structs).

Components contain generated meta-data that turns them into a special type with the following features:

- Can be attached directly to entities;
- Used to filter entities when traversing the game state (next chapter will dive into the filter API);

Components can accessed, used or passed as parameters as either pointers or as value types, just like any other struct.

## Dynamic Collections

Quantum's custom allocator exposes blittable collections as part of the rollback-able game state. Collections only support support blittable types (i.e. primitive and DSL-defined types).

To manage collection, the Frame API offers 3 methods for each:

- ```
Frame.AllocateXXX
```

: To allocate space for the collection on the heap.
- ```
Frame.FreeXXX
```

: To free/deallocate the collection's memory.
- ```
Frame.ResolveXXX
```

: To access the collection by resolving the pointer it.

**Note:** After freeing a collection, it **HAS TO** be nullified by setting it to ```
default
```

. This is required for serialization of the game state to work properly. Omitting the nullification will result in indeterministic behavior and de-synchronization. As alternative to freeing a collection and nullifying its Ptrs manually, it possible to use the ```
FreeOnComponentRemoved
```

attribute on the field in question.

### Important Notes

- Several components can reference the same collection instance.
- Dynamic collections are stored as references inside components and structs. They therefore **have to** to be _allocated_ when initializing them, and more importantly, _freed_ when they are not needed any more. If the collection is part of a component, two options are available:

  - implement the reactive callbacks ```
    ISignalOnAdd<T>
    ```

     and ```
    ISignalOnRemove<T>
    ```

     and allocate/free the collections there. (For more information on these specific signals, see the Components page in the ECS section of the Manual); or,
  - use the ```
    \[AllocateOnComponentAdded\]
    ```

     and ```
    \[FreeOnComponentRemoved\]
    ```

     attributes to let Quantum handle the allocation and deallocation when the component is added and removed respectively.
- Quantum do **NOT** pre-allocate collections from prototypes, unless there is at least value. If the collection is empty, the memory has to be manually allocated.
- Attempting to _free_ a collection more than once will throw an error and puts the heap in an invalid state internally.

### Lists

Dynamic lists can be defined in the DSL using ```
list<T> MyList
```

.

Qtn

```cs
component Targets {
list<EntityRef> Enemies;
}

```

The basic API methods for dealing with these Lists are:

- ```
Frame.AllocateList<T>()
```

- ```
Frame.FreeList(QListPtr<T> ptr)
```

- ```
Frame.ResolveList(QListPtr<T> ptr)
```


Once resolved, a list can be iterated over or manipulated with all the expected API methods of a list such as Add, Remove, Contains, IndexOf, RemoveAt, \[\], etc... .

To use the list in the component of type _Targets_ defined in the code snippet above, you could create the following system:

C#

```csharp
namespace Quantum
{
public unsafe class HandleTargets : SystemMainThread, ISignalOnComponentAdded<Targets>, ISignalOnComponentRemoved<Targets>
{
public override void Update(Frame frame)
{
foreach (var (entity, component) in frame.GetComponentIterator<Targets>()) {
// To use a list, you must first resolve its pointer via the frame
var list = frame.ResolveList(component.Enemies);

// Do stuff
}
}

public void OnAdded(Frame frame, EntityRef entity, Targets\* component)
{
// allocating a new List (returns the blittable reference type - QListPtr)
component->Enemies = frame.AllocateList<EntityRef>();
}

public void OnRemoved(Frame frame, EntityRef entity, Targets\* component)
{
// A component HAS TO de-allocate all collection it owns from the frame data, otherwise it will lead to a memory leak.
// receives the list QListPtr reference.
frame.FreeList(component->Enemies);

// All dynamic collections a component points to HAVE TO be nullified in a component's OnRemoved
// EVEN IF is only referencing an external one!
// This is to prevent serialization issues that otherwise lead to a desynchronisation.
component->Enemies = default;
}
}
}

```

### Dictionaries

Dictionaries can be declared in the DSL like so ```
dictionary<key, value> MyDictionary
```

.

Qtn

```cs
component Hazard {
dictionary<EntityRef, Int32> DamageDealt;
}

```

The basic API methods for dealing with these dictionaries are:

- ```
Frame.AllocateDictionary<K,V>()
```

- ```
Frame.FreeDictionary(QDictionaryPtr<K,V> ptr)
```

- ```
Frame.ResolveDictionary(QDictionaryPtr<K,V> ptr)
```


Just like with any other dynamic collection it is mandatory to allocate it before using it, as well as de-allocate it from the frame data and nullified it once the dictionary is no longer used. See the example provided in the section about lists here above.

### HashSet

HashSets can be declared in the DSL like so ```
hash\_set<T> MyHashSet
```

.

Qtn

```cs
component Nodes {
hash\_set<FP> ProcessedNodes;
}

```

The basic API methods for dealing with these dictionaries are:

- ```
Frame.AllocateHashSet(QHashSetPtr<T> ptr, int capacity = 8)
```

- ```
Frame.FreeHashSet(QHashSetPtr<T> ptr)
```

- ```
Frame.ResolveHashSet(QHashSetPtr<T> ptr)
```


Just like with any other dynamic collection it is mandatory to allocate it before using it, as well as de-allocate it from the frame data and nullified it once the hash set is no longer used. See the example provided in the section about lists here above.

## Enums, Unions and Bitsets

### Enums

Enums can be used to define a set of named constant values. Use the "enum" keyword to define its name and data. This example defines a simple enum, ```
EDamageType
```

, and how to use it as field in a struct:

Qtn

```cs
enum EDamageType {
None, Physical, Magic
}

struct StatsEffect {
EDamageType DamageType;
}

```

Enums are treated as integer constants on CodeGen by default, starting from 0. However, it is possible to explicitly assign integer values to Enum members if needed. Also, in scenarios where memory usage is a concern, it is possible to reduce the memory footprint of enum values by using a diffrent type as the underlying type, such as ```
byte
```

. Here's an example:

Qtn

```cs
enum EModifierOperation : Byte
{
None = 0,
Add = 1,
Subtract = 2
}

```

The ```
flags
```

keyword is used with Enum types to indicate that each value in the Enum represents a distinct bit flag. This allows you to combine multiple Enum values using bitwise operations, it is a way to represent sets of related options or states.

Qtn

```cs
flags ETeamStatus : Byte
{
None,
Winning,
SafelyWinning,
LowHealth,
MidHealth,
HighHealth,
}

```

Using the ```
flags
```

keyword also code-generates utility methods for the enum type which are more performant than regular ```
System.Enum
```

methods such as ```
IsFlagSet()
```

, which is more performance than ```
System.Enum.HasFlag()
```

as it avoids the need of value type boxing.

### Unions

C-like unions can be generated as well. The union type overlaps in memory the data layout of all the involved structs. Here is an example of how to declare a union in the DSL:

Qtn

```cs
struct DataA
{
FPVector2 Foo;
}

struct DataB
{
FP Bar;
}

union Data
{
DataA A;
DataB B;
}

```

Unions can be declared as part of a component:

Qtn

```cs
component ComponentWithUnion {
Data ComponentData;
}

```

Internally, the union type ```
Data
```

contains the logic necessary for switching between the union types as they are accessed. Here are a few usage examples:

C#

```csharp
private void UseWarriorAttack(Frame frame)
{
var character = frame.Unsafe.GetPointer<Character>(entity);
character->Data.Warrior->ImpulseDirection = FPVector3.Forward;
}

private void ResetSpellcasterMana(Frame frame)
{
var character = frame.Unsafe.GetPointer<Character>(entity);
character->Data.Spellcaster->Mana = FP.\_10;
}

```

When using unions, it is possible to use only one of the multiple internal types it contains, or it can be switched dynamically, in runtime, by simply accessing specific internal types. In the code snippets above, accessing the Warrior and Spellcaster pointers already changed the union type internally.

It is also possible to check what is the currently used union type by using the ```
Field
```

property and some of code-gen'd constants:

C#

```csharp
private bool IsWarrior(CharacterData data)
{
return data.Field == CharacterData.WARRIOR;
}

```

### Bitset

Bitsets can be used to declared fixed-size memory blocks for any desired purpose (for example fog-of-war, grid-like structures for pixel perfect game mechanics, etc.):

Qtn

```cs
struct FOWData
{
bitset\[256\] Map;
}

```

## Input

In Quantum, the runtime input exchanged between clients is also declared in the DSL. This example defines a simple movement vector and a Fire button as input for a game:

Qtn

```cs
input
{
FPVector2 Movement;
button Fire;
}

```

The input struct is polled every tick and sent to the server (when playing online).

For more information about input, such as best practices and recommended approaches to optimization, refer to this page: [Input](/quantum/current/manual/input)

## Signals

Signals are function signatures used as a decoupled inter-system communication API (a form of publisher/subscriber API). This would define a simple signal (notice the special type **entity\_ref** \- these will be listed at the end of this chapter):

Qtn

```cs
signal OnDamage(FP damage, entity\_ref entity);

```

This would generate the following interface (that can be implemented by any System):

C#

```csharp
public interface ISignalOnDamage
{
public void OnDamage(Frame frame, FP damage, EntityRef entity);
}

```

Signals are the only concept which allows the direct declaration of a pointer in Quantum's DSL, so passing data by reference can be used to modify the original data directly in their concrete implementations:

Qtn

```cs
signal OnBeforeDamage(FP damage, Resources\* resources);

```

Notice this allows the passing of a component pointer (instead of the entity reference type).

## Events

Events are a fine-grained solution to communicate what happens inside the simulation to the rendering engine / view (they should never be used to modify/update part of the game state). Use the "event" keyword to define its name and data:

Find detailed information about events in the [Frame Events Manual](/quantum/current/manual/quantum-ecs/game-events#frame_events).

Define an event using the Quantum DSL

Qtn

```cs
event MyEvent{
int Foo;
}

```

Trigger the event from the simulation

C#

```csharp
f.Events.MyEvent(2022);

```

And subscribe and consume the event in Unity

C#

```csharp
QuantumEvent.Subscribe(listener: this, handler: (MyEvent e) => Debug.Log($"MyEvent {e.Foo}"));

```

## Globals

It is possible to define globally accessible variables in the DSL. Globals can be declared in any .qtn file by using the ```
global
```

scope.

Qtn

```cs
global {
// Any type that is valid in the DSL can also be used.
FP MyGlobalValue;
}

```

Like all things DSL-defined, global variables are part of the state and are fully compatible with the predict-rollback system.

Variables declared in the global scope are made available through the Frame API. They can be accessed (read/write) from any place that has access to the frame - see the _Systems_ document in the ECS section.

**N.B.:** An alternative to global variables are the Singleton Components; for more information please refer to the _Components_ page in the ECS section of the manual.

## Special Types

Quantum has a few special types that are used to either abstract complex concepts (entity reference, player indexes, etc.), or to protect against common mistakes with unmanaged code, or both. The following special types are available to be used inside other data types (including in components, also in events, signals, etc.):

- ```
player\_ref
```

: represents a runtime player index (also cast to and from Int32). When defined in a component, can be used to store which player controls the associated entity (combined with Quantum's player-index-based input).
- ```
entity\_ref
```

: because each frame/tick data in quantum resides on a separate memory region/block (Quantum keeps a a few copies to support rollbacks), pointers cannot be cached in-between frames (nor in the game state neither in Unity scripts). An entity ref abstracts an entity's index and version properties (protecting the developer from accidentally accessing deprecated data over destroyed or reused entity slots with old refs).
- ```
asset\_ref<AssetType>
```

: rollback-able reference to a data asset instance from the Quantum asset database (please refer to the data assets chapter).
- ```
list<T>
```

, ```
dictionary<K,T>
```

: dynamic collection references (stored in Quantum's frame heap). Only supports blittable types (primitives + DSL-defined types).
- ```
array<Type>\[size\]
```

: fixed sized "arrays" to represent data collections. A normal C# array would be a heap-allocated object reference (it has properties, etc.), which violates Quantum's memory requirements, so the special array type generates a pointer based simple API to keep rollback-able data collections inside the game state;

### A Note On Assets

Assets are a special feature of Quantum that let the developer define data-driven containers (normal classes, with inheritance, polymorphic methods, etc.) that end up as immutable instances inside an indexed database. The "asset" keyword is used to assign an (existing) class as a data asset that can have references assigned inside the game state (please refer to the Data Assets chapter to learn more about features and restrictions):

Qtn

```cs
asset CharacterData; // the CharacterData class is partially defined in a normal C# file by the developer

```

The following struct show some valid examples of the types above (sometimes referencing previously defined types):

Qtn

```cs
struct SpecialData
{
player\_ref Player;
entity\_ref Character;
entity\_ref AnotherEntity;
asset\_ref<CharacterData> CharacterData;
array<FP>\[10\] TenNumbers;
}

```

## Available Types

When working in the DSL, you can use a variety of types. Some are pre-imported by the parsers, while others need to be manually imported.

### By default

Quantum's DSL parser has a list of pre-imported cross-platform deterministic types that can be used in the game state definition:

- Boolean / bool - internally gets wrapped in QBoolean which works identically (get/set, compare, etc...)
- Byte
- SByte
- UInt16 / Int16
- UInt32 / Int32
- UInt64 / Int64
- FP
- FPVector2
- FPVector3
- FPMatrix
- FPQuaternion
- PlayerRef / player\_ref in the DSL
- EntityRef / entity\_ref in the DSL
- LayerMask
- NullableFP / FP? in the DSL
- NullableFPVector2 / FPVector2? in the DSL
- NullableFPVector3 / FPVector3? in the DSL
- QString is for UTF-16 (aka Unicode in .NET)
- QStringUtf8 is always UTF-8
- Hit
- Hit3D
- Shape2D
- Shape3D
- Joint, DistanceJoint, SpringJoint and HingeJoint

**Note on QStrings**: ```
N
```

represents the total size of the string in bytes minus 2 bytes used for bookkeeping. In other words ```
QString<64>
```

will use 64 bytes for a string with a max byte length of 62 bytes, i.e. up to 31 UTF-16 characters.

### Manual Import

If a type that is not listed in the previous section is needed, it has to be imported manually when using it in QTN files.

#### Namespaces / Types outside of Quantum

#### Importing specific types

To import types defined in other namespaces and use them directly in components or global, on the DSL, use the following syntax:

Qtn

```cs
import MyInterface;
or
import MyNameSpace.Utils;

```

For an enum the syntax is as follows:

Qtn

```cs
import enum MyEnum(underlying\_type);

// This syntax is identical for Quantum specific enums
import enum Shape3DType(byte);

```

#### Including namespaces

In some cases it might also be necessary to add ```
using MyNamespace;
```

to any QTN file for such namespace to be included on the generated class.

#### Built-In Quantum Type and Custom Type

When importing a Quantum built-in type or a custom type, the struct size is predefined in their C# declaration. It is therefore important to add some safety measures.

C#

```csharp
namespace Quantum {
\[StructLayout(LayoutKind.Explicit)\]
public struct Foo {
public const int SIZE = sizeof(Int32) \* 2;
\[FieldOffset(0)\]
public Int32 A;
\[FieldOffset(sizeof(Int32))\]
public Int32 B;
}
}

```

Qtn

```cs
#define FOO\_SIZE 8 // Define a constant value with the known size of the struct
import struct Foo(8);

```

To ensure the expected size of the struct is equal to the actual size, it is recommended to add an ```
Assert
```

as shown below in one of your systems.

C#

```csharp
public unsafe class MyStructSizeCheckingSystem : SystemMainThread{
public override void OnInit(Frame frame)
{
Assert.Check(Constants.FOO\_SIZE == Foo.SIZE);
}
}

```

If the size of the built-in struct changes during an upgrade, this ```
Assert
```

will throw and allow you to update the values in the DSL.

## Attributes

Quantum supports several attributes to present parameters in the Inspector.

The attributes are contained within the ```
Quantum.Inspector
```

namespace.

Field: Excludes field from a the prototype generated for the component.


Component: No prototype will be generated for this component.


| Attribute | Parameters | Description |
| --- | --- | --- |
| **DrawIf** | ```<br>string<br>```<br> fieldName<br> <br>```<br>long<br>```<br> value<br> <br>```<br>CompareOperator<br>```<br> compare<br> <br>```<br>HideType<br>```<br> hide | Displays the property only if the condition evaluates to true.<br> <br>_fieldName_ = the name of the property to evaluate.<br> <br>_value_ = the value used for comparison.<br> <br>_compare_ = the comparison operation to be performed ```<br>Equal<br>```<br>, ```<br>NotEqual<br>```<br>, ```<br>Less<br>```<br>, ```<br>LessOrEqual<br>```<br>, ```<br>GreaterOrEqual<br>```<br> or ```<br>Greater<br>```<br>.<br> <br>_hide_ = the field's behavior when the expression evaluates to ```<br>False<br>```<br>:```<br>Hide<br>```<br> or ```<br>ReadOnly<br>```<br>.<br> <br> For more information on compare and hide, see below. |
| **Header** | ```<br>string<br>```<br> header | Adds a header above the property.<br> <br>_header_ = the header text to display. |
| **HideInInspector** |  | Serializes the field and hides the following property in the Unity inspector. |
| **Layer** |  | Can only be applied to type **int**.<br> <br>Will call ```<br>EditorGUI.LayerField<br>```<br> on the field. |
| **Optional** | ```<br>string<br>```<br> enabledPropertyPath | Allows to turn the display of a property on/off.<br> <br>_enabledPropertyPath_ = the path to the ```<br>bool<br>```<br> used to evaluate the toggle. |
| **Space** |  | Adds a space above the property |
| **Tooltip** | ```<br>string<br>```<br> tooltip | Displays a tool tip when hovering over the property.<br> <br>_tooltip_ = the tip to display. |
| **ArrayLength**<br> ONLY FOR CSharp | ```<br>int<br>```<br> length | Using _length_ allows to define the size of a an array. |
| **ArrayLength**<br> ONLY FOR CSharp | ```<br>int<br>```<br> minLength<br> <br>```<br>int<br>```<br> maxLength | Using _minLength_ and _maxLength_ allows to define a range for the size in the Inspector. <br> <br>The final size can then be set in the Inspector.<br> <br>( _minLength_ and _maxLength_ are inclusive) |
| **ExcludeFromPrototype** |  | Can be applied a component or component fields. |
| **OnlyInPrototype** |  | Can be applied to a field which then will be ignored in the object state and only be added to the prototype. |
| **OnlyInPrototype** | ```<br>string<br>```<br> fieldName<br>```<br>string<br>```<br> fieldType | Can be applied to a component which then will add the field to the prototype but ignore it in the object state (similar to the above). |
| **PreserveInPrototype** |  | Added to a type marks it as usable in prototypes and prevents prototype class from being emit.<br> Added to a field only affects a specific field. Useful for simple ```<br>\[Serializable\]<br>```<br> structs as it avoids having to use ```<br>\_Prototype<br>```<br> types on Unity side. |
| **AllocateOnComponentAdded** |  | Can be applied to dynamic collections. <br> <br>This will allocate memory for the collection if it has not already been allocated when the component holding the collection is added to an entity. |
| **FreeOnComponentRemoved** |  | Can be applied to dynamic collections and ```<br>Ptrs<br>```<br>. <br> <br>This will deallocate the associated memory and nullify the ```<br>Ptr<br>```<br> held in the field when the component is removed.<br> <br> IMPORTANT: Do **NOT** use this attribute in combination with cross-referenced collections as it only nullifies the ```<br>Ptr<br>```<br> held in that particular field and the others will be pointing to invalid memory. |

The _Attributes_ can be used in both C# and qtn files unless otherwise specified; however, there are some syntactic differences.### Use in CSharp

In C# files, attributes can be used and concatenated like any other attribute.

C#

```csharp
// Multiple single attributes
\[Header("Example Array")\]\[Tooltip("min = 1\\nmax = 20")\] public FP\[\] TestArray = new FP\[20\];

// Multiple concatenated attributes
\[Header("Example Array"), Tooltip("min = 1\\nmax = 20")\] public FP\[\] TestArray = new FP\[20\];

```

### Use in qtn

In qtn files, the usage of single attributes remains the same as in C#.

Qtn

```cs
\[Header("Example Array")\] array<FP>\[20\] TestArray;

```

When combining multiple attributes, they **have** to be concatenated.

Qtn

```cs
\[Header("Example Array"), Tooltip("min = 1\\nmax = 20")\] array<FP>\[20\] TestArray;

```

## Compiler Options

The following compiler options are currently available to be used inside Quantum's DSL files (more will be added in the future):

Qtn

```cs
// pre defining max number of players (default is 6, absolute max is 64)
#pragma max\_players 16

// increase the component count from 256 to 512
#pragma max\_components 512

// numeric constants (useable inside the DSL by MY\_NUMBER and useable in code by Constants.MY\_NUMBER)
#define MY\_NUMBER 10

// overriding the base class name for the generated constants (default is "Constants")
#pragma constants\_class\_name MyFancyConstants

```

## Custom FP Constants

You can also define custom ```
FP
```

constants inside Quantum's DSL files. ex:

Qtn

```cs
// in a DSL file
#define Pi 3.14

```

Then, Quantum codegen will generate the corresponding constant in the ```
FP
```

struct:

Qtn

```cs
// 3.14
FP constant = Constants.Pi;

```

It will also generate the corresponding raw value as well:

Qtn

```cs
// 3.14 Raw
var rawConstant = Constants.Raw.Pi;

```

Back to top

- [Introduction](#introduction)
- [Components](#components)
- [Structs](#structs)

  - [DSL Defined](#dsl-defined)
  - [CSharp Defined](#csharp-defined)
  - [Components vs. Structs](#components-vs.structs)

- [Dynamic Collections](#dynamic-collections)

  - [Important Notes](#important-notes)
  - [Lists](#lists)
  - [Dictionaries](#dictionaries)
  - [HashSet](#hashset)

- [Enums, Unions and Bitsets](#enums-unions-and-bitsets)

  - [Enums](#enums)
  - [Unions](#unions)
  - [Bitset](#bitset)

- [Input](#input)
- [Signals](#signals)
- [Events](#events)
- [Globals](#globals)
- [Special Types](#special-types)

  - [A Note On Assets](#a-note-on-assets)

- [Available Types](#available-types)

  - [By default](#by-default)
  - [Manual Import](#manual-import)

- [Attributes](#attributes)

  - [Use in qtn](#use-in-qtn)

- [Compiler Options](#compiler-options)
- [Custom FP Constants](#custom-fp-constants)