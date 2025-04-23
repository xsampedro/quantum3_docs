# assets-simulation

_Source: https://doc.photonengine.com/quantum/current/manual/assets/assets-simulation_

# Assets in Simulation

## Data Asset Classes

Quantum assets are C# classes that will act as immutable data containers during runtime. A few rules define how these assets must be designed, implemented and used in Quantum.

Here is a minimal definition of an asset class (for a character spec) with some simple deterministic properties:

C#

```
```csharp
public class CharacterSpec : AssetObject {
public FP Speed;
public FP MaxHealth;
}

```

```

In Unity, \`AssetObject\` is a subclass of \`UnityEngine.ScriptableObject\`, so the definition should be saved in a file that matches the class name (e.g. \`CharacterSpec.cs\`). Also, adding methods \`Update\` and \`Start\` may cause a compile error, depending on Unity version.

Creating and loading instances of asset classes into the database (editing from Unity) will be covered later in this chapter.

## Using and Linking Assets

Asset instances are immutable objects that must be carried as references. Because normal C# object references are not allowed to be included into our memory aligned ECS structs, the _asset\_ref_ special type must be used inside the DSL to declare properties inside the game state (from entities, components or any other transient data structure):

Qtn

```
```cs
component CharacterData {
// reference to an immutable instance of CharacterSpec (from the Quantum asset database)
asset\_ref<CharacterSpec> Spec;
// other component data
}

```

```

To assign an asset reference when creating a Character entity, one option is to obtain the instance directly from the frame asset database and set it to the property:

C#

```
```csharp
// assuming cd is a pointer to the CharacterData component
// using the SLOW string path option (fast data driven asset refs will be explained next)
cd->Spec = frame.FindAsset<CharacterSpec>("path-to-spec");

```

```

The basic use of assets is to read data in runtime and apply it to any computation inside systems. The following example uses the _Speed_ value from the assigned _CharacterSpec_ to compute the corresponding character velocity (physics engine):

C#

```
```csharp
// consider cd a CharacterData\*, and body a PhysicsBody2D\* (from a component filter, for example)
var spec = frame.FindAsset(cd->Spec);
body->Velocity = FPVector2.Right \* spec.Speed;

```

```

### A Note On Determinism

Notice that the above code only **reads** the _Speed_ property to compute the desired velocity for the character during runtime, but its value (speed) is never changed.

It is completely safe and valid to switch a game state asset reference in runtime from inside an Update (as asset\_ref is a rollback-able type which hence can be part of the game state).

However, changing the **values** of properties of a data asset is NOT DETERMINISTIC (as the internal data on assets is not considered part of the game state, so it is never rolled back).

The following snippet shows examples of what is safe (switching refs) and not safe (changing internal data) during runtime:

C#

```
```csharp
// cd is a CharacterData\*

// this is VALID and SAFE, as the CharacterSpec asset ref is part of the game state
cd->Spec = frame.FindAsset<CharacterSpec>("anotherCharacterSpec-path");

// this is NOR valid NEITHER deterministic, as the internal data from an asset is NOT part of the transient game state:
var spec = frame.FindAsset<CharacterSpec>("anotherCharacterSpec-path");
// (DO NOT do this) changing a value directly in the asset object instance
spec.Speed = 10;

```

```

## Asset Inheritance

It is possible to use inheritance in data assets, which gives much more flexibility to the developer (specially when used together with polymorphic methods).

The basic step for inheritance is to create an abstract base asset class (we'll continue with our _CharacterSpec_ example):

C#

```
```csharp
public abstract class CharacterSpec : AssetObject {
public FP Speed;
public FP MaxHealth;
}

```

```

Concrete sub-classes of _CharacterSpec_ may add custom data properties of their own, and must be marked as _Serializable_ types:

C#

```
```csharp
public class MageSpec : CharacterSpec {
public FP HealthRegenerationFactor;
}

public class WarriorSpec : CharacterSpec {
public FP Armour;
}

```

```

### Data-Driven Polymorphism

Having gameplay logic to direct evaluate (in if or switch statements) the concrete _CharacterSpec_ class would be very bad design, so asset inheritance makes more sense when coupled with polymorphic methods.

Notice that adding logic to data assets means implementing logic in the quantum.state project and this logic still have to consider the following restrictions:

- Operate on transient game state data: that means logic methods in data assets must receive transient data as parameters (either entity pointers or the Frame object itself);
- Only read, never modify data on the assets themselves: assets must still be treated as _immutable_ read-only instances;

The following example adds a virtual method to the base class, and a custom implementation to one of the subclasses (notice we use the _Health_ field defined for the _Character_ entity more to the top of this document):

C#

```
```csharp
public unsafe abstract class CharacterSpec : AssetObject {
public FP Speed;
public FP MaxHealth;
public virtual void UpdateCharacter(Frame frame, EntityRef entity, CharacterData\* data) {
if (data->Health < 0)
frame.Destroy(entity);
}
}

public unsafe class MageSpec : CharacterSpec {
public FP HealthRegenerationFactor;
// reads data from own instance and uses it to update transient health of Character pointer passed as param
public override void UpdateCharacter(Frame frame, EntityRef entity, CharacterData\* data) {
data->Health += HealthRegenerationFactor \* frame.DeltaTime;
base.UpdateCharacter(frame, entity, data);
}
}

```

```

To use this flexible method implementation independently of the concrete asset assigned to each _CharacterData_, this could be executed from any System:

C#

```
```csharp
// Assuming data is the pointer to a specific entity's CharacterData component, and entity is the corresponding EntityRef:

var spec = frame.FindAsset(data->Spec);
// Updating Health using data-driven polymorphism (behavior depends on the data asset type and instance assigned to character
spec.UpdateCharacter(frame, entity, data);

```

```

### Using DSL Generated Structs In Assets

```
Structs
```

defined in the DSL can also be used on assets. The DSL struct must be annotated with the ```
\[Serializable\]
```

 attribute, otherwise the data is not inspectable in Unity.

```
```
\[Serializable\]
struct Foo {
 int Bar;
}

```

```

Using the DSL ```
struct
```

in a Quantum asset.

C#

```
```csharp
public class FooUser : AssetObject {
public Foo F;
}

```

```

If a struct is not ```
\[Serializable\]
```

-friendly (e.g. because it is an union or contains a Quantum collection), prototype can be used instead:

C#

```
```csharp
using Quantum.Prototypes;

 public class FooUser : AssetObject {
 public FooPrototype F;
 }

```

```

The prototype can be materialized into the simulation struct when needed:

C#

```
```csharp
Foo f = new Foo();
fooUser.F.Materialize(frame, ref f, default);

```

```

## Adding Static Assets at Runtime

It is possible to add static assets to the asset database at runtime, before the Quantum simulation starts. This can be useful for scenarios such as downloading maps from a backend or procedurally generating content. When adding assets at runtime, it's crucial to ensure that each asset has a deterministic GUID to maintain consistency across all clients.

There are two ways to generate deterministic ```
AssetGuid
```

:

1. Using a constant:
   - This is the simplest way to have a deterministic ```
     AssetGuid
     ```

     .
   - You must ensure that another asset is not assigned the same ```
     AssetGuid
     ```

     .
   - This is fine as long as you have a fixed amount of assets to add, and you are sure that the same asset will always be assigned the same ```
     AssetGuid
     ```

     .
2. Generation:
   - The ```
     QuantumUnityDB.CreateRuntimeDeterministicGuid
     ```

      method provides an API to generate ```
     AssetGuid
     ```

     s.
   - It uses the asset object's name as a seed to generate a deterministic ```
     AssetGuid
     ```

     .
   - You should use this approach if the amount of assets you need to add is not fixed.

Usage:

C#

```
```csharp
// create any asset
var assetObject = AssetObject.Create<MyAssetObjectType>();

// set its name
assetObject.name = "My Unique Asset Object Name";

// get a deterministic GUID
var guid = QuantumUnityDB.CreateRuntimeDeterministicGuid(assetObject);

// add the asset to the asset database
QuantumUnityDB.Global.AddAsset(assetObject);

// set the GUID
assetObject.Guid = guid;

```

```

After the asset is added to the asset database, it is effectively the same as if the asset was created and added in the editor. However, this is only valid if the game has not started yet. Once the game has started, the static asset database should not be mutated.

If you need to add or mutate assets during gameplay, you must use the ```
DynamicDB
```

 API instead.

## Dynamic Assets

Assets can be created at runtime, by the simulation. This feature is called _DynamicAssetDB_.

C#

```
```csharp
var mageSpec = AssetObject.Create<MageSpec>();
mageSpec.Speed = 1;
mageSpec.MaxHealth = 100;
frame.AddAsset(mageSpec);

```

```

Such asset can be loaded and disposed of just like any other asset:

C#

```
```csharp
MageSpec asset = frame.FindAsset<MageSpec>(assetGuid);
frame.DisposeAsset(assetGuid);

```

```

Dynamic assets are not synced between peers. Instead, the code that creates new assets needs to be deterministic and ensure that each peer will generate an asset using the same values.

The only exception to the rule above is when there is a late-join - the new client will receive a snapshot of the _DynamicAssetDB_ along with the latest frame data. Unlike serialization of the frame, serialization and deserialization of dynamic assets is delegated outside of the simulation, to ```
IAssetSerializer
```

interface. When run in Unity, ```
QuantumUnityJsonSerializer
```

is used by default: it is able to serialize/deserialize any Unity-serializable type.

### Initializing DynamicAssetDB

Simulation can be initialized with preexisting dynamic assets. Similar to adding assets during the simulation, these need to be deterministic across clients.

First, an instance of ```
DynamicAssetDB
```

needs to be created and filled with assets:

C#

```
```csharp
var initialAssets = new DynamicAssetDB();
initialAssets.AddAsset(mageSpec);
initialAssets.AddAsset(warriorSpec);
...

```

```

Second, ```
QuantumGame.StartParameters.InitialDynamicAssets
```

 needs to be used to pass the instance to a new simulation. In Unity, since it is the ```
QuantumRunner
```

behaviour that manages a ```
QuantumGame
```

, ```
QuantumRunner.StartParamters.InitialDynamicAssets
```

is used instead.

## Built in Assets

Quantum also comes shipped with several built-in assets, such as:

- **SimulationConfig** \- defines many specifications for a Quantum simulation, from scene management setup, heap configuration, thread count and Physics/Navigation settings;
- **DeterministicConfig** \- specifies details on the game session, such as it's simulation rate, the checksum interval and lots of configuration regarding Input related to both the client and the server;
- **QuantumEditorSettings** \- has the definition for many editor-only details, like the folder the DB should be based in, the color of the Gizmos and auto build options for automatically baking maps, nav meshes, etc;
- **BinaryData** \- an asset that allows the user to reference arbitrary binary information (in the form of a ```
byte\[\]
```

). For example, by default the Physics and the Navigation engines uses binary data assets to store information like the static triangles data. This asset also has built in utilities to compress and decompress the data using gzip.
- **CharacterController3DConfig** \- config asset for the built in 3D KCC.
- **CharacterController2DConfig** \- config asset for the built in 2D KCC.
- **PhysicsMaterial** \- defines a ```
Physics Material
```

for Quantum's 3D physics engine.
- **PolygonCollider** \- defines a ```
Polygon Collider
```

for Quantum's 2D physics engine.
- **NavMesh** \- defines a ```
NavMesh
```

used by Quantum's navigation system.
- **NavMeshAgentConfig** \- defines a ```
NavMesh Agent Config
```

for Quantum's navigation system.
- **Map** \- stores many static per-scene information such as Physics settings, colliders, NavMesh settings, links, regions and also the Scene Entity Prototypes on that Map. Every Map is correlated with a single Unity scene.

Back to top

- [Data Asset Classes](#data-asset-classes)
- [Using and Linking Assets](#using-and-linking-assets)

  - [A Note On Determinism](#a-note-on-determinism)

- [Asset Inheritance](#asset-inheritance)

  - [Data-Driven Polymorphism](#data-driven-polymorphism)
  - [Using DSL Generated Structs In Assets](#using-dsl-generated-structs-in-assets)

- [Adding Static Assets at Runtime](#adding-static-assets-at-runtime)
- [Dynamic Assets](#dynamic-assets)

  - [Initializing DynamicAssetDB](#initializing-dynamicassetdb)

- [Built in Assets](#built-in-assets)