# blackboard

_Source: https://doc.photonengine.com/quantum/current/addons/bot-sdk/blackboard_

# Blackboard

### Introduction

The Blackboard is a very useful tool built as part of the Bot SDK core.

It is a Quantum component which has the purpose of serving as a data storage which any other Quantum logic can read from and write into. It is split in three main pieces:

1. The ```
AIBlackboardComponent
```

    a Quantum component which can be added to entities. It has the runtime storage with the data that can change during the game simulation;
2. The ```
AIBlackboard
```

    asset which is created in Unity and has the data layout that will be applied to the component (the entries Types and Keys);
3. The ```
AIBlackboardInitializer
```

, and asset that stores the initial values the Blackboard entries should have, when the blackboard component is first initialised. Using the initializer is optional;

### Why use the Blackboard

Even though storing data in the Blackboard is very similar to storing data into a regular Quantum component, using the Blackboard alongside with Bot SDK comes with some extra advantages:

1. The Bot SDK visual editor has an integration with Blackboard assets and it allows users to define its entries and initial values directly on the editor window, which allows for quick changes for developers and can even be used directly by Game Designers;
2. The editor also comes with the capability of drag-and-dropping Blackboard nodes to the AI graph and linking them to other nodes such as HFSM Actions and BT Leaf nodes (and much more). It makes it very easy then to reference a data entry directly via the editor;

### Blackboard disadvantages

When using the Blackboard, take into consideration:

1. Every entry on the Blackboard is a ```
union
```

, so the size of every entry, in runtime memory, is equivalent to the biggest type supported by the union, which is the ```
AssetRef
```

    type of size ```
8
```

    bytes. When using _many_ Blackboards at the same time, it doesn't matter if most of the entries are booleans (size = ```
1 byte
```

), the actual consumed size will be ```
1
```

    byte. If frame size becomes a bottleneck, it is possible to consider moving some variables out of the Blackboard;
2. The Blackboard is mostly useful for variables that change in runtime. If any of the variables used is not supposed to change, it could be better to have such data stored in any other Quantum asset instead, as to avoid such data from being created as part of the Blackboard runtime data in the frame;

### The cardinality

The Blackboard asset is read-only and is meant to be re-used by any amount of entities in order to define how their Blackboard component layout will be defined.

So the asset itself should be referenced by the entities, which will each one have their own runtime, changeable data, created after the components initialisation.

The cardinality then is \[1..n\]: for each Blackboard asset, there are ```
n
```

 agents which will reference such asset.

### A simple usage example

Consider an entity which runs around in a map collecting coins and, once it collects ```
3
```

coins, it stops and idle. For such entity's AI, it is clearly helpful to keep track of the collectibles amount.

For this, a new ```
int
```

 entry can be created in the Blackboard with the name ```
CollectedCoinsAmonut
```

.

Then whenever the entity collects a new coin, it increases the value on the blackboard by ```
1
```

 and checks if the new value stored is already ```
3
```

. This check could be done, for example, every frame as part of the AI decision making process, as polling data from the Blackboard is a fast operation, based on a Dictionary lookup.

As every entity with the Blackboard component has its own copy of the runtime data, one entity can even read the Blackboard of another one, or a separate System can read the values from every agent and so on.

### The Blackboard on the Visual Editor

The Visual Editor already comes with a sub-menu for the Blackboard.

It is named Blackboard Variables and it is located on the left side panel:

![Blackboard Visual Editor Empty](/docs/img/quantum/v2/addons/bot-sdk/blackboard4.png)

Press the **+** button in order to create new blackboard variables.

**Double Click** an entry in order to edit it.

There are also alternatives that can be found in the **right click menu**.

When creating/editing a new entry, define:

- The variable ```
Name
```

, which is internally used to generate the Key, used when storing creating the variable in the dictionary;
- The variable ```
Type
```

, which can is a set of pre-defined supported Types, selectable from the dropdown menu;
- The ```
HasDefault
```

checkbox is used to inform if the variable will or not be initialized to a default value upon setup;
- The ```
Default
```

value to be used.

![Blackboard Visual Editor Creating Entry](/docs/img/quantum/v2/addons/bot-sdk/blackboard5.png)

Now it is up to the user to decide, during development, which blackboard entries will be necessary for the entities AI.

![Blackboard Visual Editor Many Entries](/docs/img/quantum/v2/addons/bot-sdk/blackboard6.png)

**PS:** it is possible to invoke the context menu by with the **Right Mouse Button** on a blackboard entry.

The types currently supported by the Blackboard are:

1. Boolean;
2. Byte;
3. Integer;
4. FP;
5. Vector2;
6. Vector3;
7. Entity Ref;
8. Asset Ref.

With the Blackboard defined, two extra Quantum assets are generated when the AI document is compiled and, by default, it is output to the folder ```
Assets/Resources/DB/CircuitExport/Blackboard\_Assets
```

. The assets are:

- The ```
  AIBlackboardAsset
  ```

   which has the types and keys layout;
- The ```
  AIBlackboardInitializer
  ```

   that also stores the initial values for the entries if applicable.

**PS:** both assets reference eachother, which makes the code for initializing the component simpler (examples in the code section).

![Blackboard Assets](/docs/img/quantum/v2/addons/bot-sdk/blackboard7.png)### Blackboard Nodes

With the variables defined, it is possible to drag and drop them to the graph in order to create Blackboard Nodes. These nodes always have two outbound slots: the ```
Key
```

and the ```
Value
```

.

![Blackboard Assets](/docs/img/quantum/v2/addons/bot-sdk/blackboard-node.png)

The ```
Key
```

slot is used to be linked to fields of type ```
AIBlackboardValueKey
```

, which can be used to replace hardcoded strings when informing the key of the variable to get/set. And, of course, by removing the hardcoded keys, the code gets more flexible and reliable.

Let's analyze, in terms of quantum code, how the Get/Set methods should be used, depending on whether a hardcoded string key, or the key from a Blackboard node are used:

C#

```csharp
// \-\- Using hardcoded keys --
var bbComponent = f.Unsafe.GetPointer<AIBlackboardComponent>(entityRef);
// Reading
var value = bbComponent->GetInteger(frame, "someKey");
// Writing
bbComponent->Set(frame, "someKey", value);

// \-\- Using keys from Blackboard nodes --
public AIBlackboardValueKey PickupsKey;
// Reading
var value = bbComponent->GetInteger(frame, PickupsKey.Key);
// Writing
bbComponent->Set(frame, PickupsKey.Key, value);

```

With that new field declared, it is possible to link it directly with a Blackboard Node's ```
Key
```

slot in the Visual Editor:

![Blackboard Assets](/docs/img/quantum/v2/addons/bot-sdk/blackboard-node-key-connected.PNG)

Besides of using the ```
Key
```

 slot as explained, it is also possible to link the ```
Value
```

slot to define fields of that same type (an integer blackboard variable linked to an integer field on some nodes such as Actions, BTLeaf, Response Curves, etc). The ```
Default
```

 value defined on the left panel is the value which will be used (baked into the target node asset).

![Blackboard Assets](/docs/img/quantum/v2/addons/bot-sdk/blackboard-node-value-connected.PNG)### Blackboard Quantum code

**Initializing the Blackboard component**

The important part for initializing the blackboard component is to have a reference to the ```
AIBlackboardInitializer
```

asset which was created after compiling the AI document on the Visual Editor.

Then, on the simulation code, use this to initialize the component:

C#

```csharp
// \-\- Blackboard setup
// First, create the blackboard component (or have it created on the Entity Prototype)
var blackboardComponent = new AIBlackboardComponent();
// Find the Blackboard Initializer asset
var bbInitializerAsset = f.FindAsset<AIBlackboardInitializer>(blackboardAsset.BlackboardInitializer.Id);
// Call the static initialization method passing the blackboard component and the asset
AIBlackboardInitializer.InitializeBlackboard(f, &blackboardComponent, bbInitializerAsset);
// Set the blackboard into to the entity
f.Set(littleGuyEntity, blackboardComponent);

```

That's it. Once the initialization is done, it is ready to read/write from the Blackboard using its API:

C#

```csharp
// There is one method for each specific blackboard type (int, byte, FP, boolean, FP vectors and entityRef)
blackboardComponent->GetInteger(frame, key);

// For the setter method, there are different overrides depending on the type of data passed as the value
blackboardComponent->Set(frame, key, value);

```

**Using Entity Prototypes**

It is also possible to reference the blackboard asset via the Unity editor by using components in the Entity Prototypes. Add the ```
AIBlackboardComponent
```

 and define the ```
Board
```

field. Then, use it on the simulation on the initialisation step if needed.

**Disposing the Blackboard component**

When destroying an entity which uses the Blackboard component, it is important **to free the component as to avoid memory leaks.**

To dispose the blackboard memory, use ```
blackboardComponent->Free(frame);
```

**Initializing and Clearing with Component Callbacks**

The ```
BotSDKSystem
```

that comes with the SDK has an example on how to init/free the ```
AIBlackboardComponent
```

by using the ```
OnComponentAdded
```

and ```
OnComponentRemoved
```

callbacks.

**PS.:** Bot SDK already comes with some sample code that demonstrates reading/writing with the Blackboard, located at: ```
BotSDK/Samples
```

and check the files ```
IncreaseBlackboardInt.cs
```

, ```
SetBlackboardInt.cs
```

and ```
HFSM.CheckBlackboardInt.cs
```

Back to top

- [Introduction](#introduction)
- [Why use the Blackboard](#why-use-the-blackboard)
- [Blackboard disadvantages](#blackboard-disadvantages)
- [The cardinality](#the-cardinality)
- [A simple usage example](#a-simple-usage-example)
- [The Blackboard on the Visual Editor](#the-blackboard-on-the-visual-editor)
- [Blackboard Nodes](#blackboard-nodes)
- [Blackboard Quantum code](#blackboard-quantum-code)