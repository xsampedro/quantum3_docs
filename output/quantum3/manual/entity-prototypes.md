# entity-prototypes

_Source: https://doc.photonengine.com/quantum/current/manual/entity-prototypes_

# Entity Prototypes

## Introduction

To facilitate data driven design, Quantum features entity prototypes.

A Quantum Entity Prototype is a serialized version of an entity that includes:

- Composition (i.e. which components it is made of); and,
- Data (i.e. the components' properties and their initial value).

This allows for a clean separation of data and behaviour, while enabling designers to tweak the former without programmers having to constantly edit the latter.

## Setting up a Prototype

Entity prototypes can be set up in the Unity Editor.

### Basic

To create it, add the ```
QuantumEntityPrototype
```

 component to any GameObject.

![Entity Prototype Script on an empty GameObjet](/docs/img/quantum/v3/manual/entityprototype-basic.png)
 Basic Entity Prototype (Empty GameObject + Entity Prototype Script).


The ```
QuantumEntityPrototype
```

script allows to set up and define the parameters for the most commonly used components for both 2D and 3D.

- Transform (including Transform2DVertical for 2D);
- PhysicsCollider;
- PhysicsBody;
- NavMeshPathFinder;
- NavMeshSteeringAgent;
- NavMeshAvoidanceAgent.

The dependencies for the Physics and NavMesh related agents are respected.

For more information, please read their respective documentation.

### Custom Components

Additional components can be added to an entity prototypes via either:

- The **+** button in the ```
Entity Components
```

list; or,
- The regular Unity _Add Component_ button by searching for the right ```
QPrototype
```

component.

#### Note on Collections

Dynamic collections in components are only automatically allocated **IF** there is at least one item in the prototype. Otherwise, the collection will have to be allocated manually. For more information on the subject, refer to the [Dynamics Collection entry on the DSL page](/quantum/current/manual/quantum-ecs/dsl#dynamic_collections).

### Hierarchy

In ECS the concept of entity/GameObject hierarchy does not exist. As such entity prototypes do not support hierarchies or nesting.

Although child prototypes are not supported directly, you can:

1. Create separate prototypes in the scene and bake them;
2. Link them by keeping a reference in a component;
3. Update the position of the "child" manually.

_Note:_ Prototypes that are not baked in scene will have to follow a different workflow where the entities are created and linked in code.

It is possible to have hierarchies in game objects (View), however hierarchies in entities (Simulation) will have to be handled by you.

## Creating/Instantiating a Prototype

Once an entity prototype has been defined in Unity, there are various ways to include it in the simulation.

### Baked in the Scene/Map

If the entity prototype is created as part of a Unity Scene, it will be baked into the corresponding Map Asset. The baked entity prototype will be loaded when the Map is initialized with the values it was baked with.

**N.B.:** If a Scene's entity prototype is edited or has its values changed, the Map Data has to be re-baked (which depending on the project setup, it might be done automatically during some editor actions like saving the project).

### In Code

To create a new entity from a ```
QuantumEntityPrototype
```

, follow these steps:

1. Create a Unity Prefab of the GameObject which has the ```
   QuantumEntityPrototype
   ```

    component;
2. Place the Prefab in any folder included in the ```
   QuantumEditorSettings
   ```

    asset in ```
   Asset Search Paths
   ```

   , which by default includes all the ```
   Assets
   ```

    folder:

![Entity Prototype Asset](/docs/img/quantum/v3/manual/entityprototype-asset.png)
 Entity Prototype Prefab + Generated Entity Prototype Asset.


This automatically generates an ```
EntityPrototype
```

asset which is associated with the prefab itself, as shown on the screenshot above.

3. In the editor, it is possible to reference such ```
EntityPrototype
```

    assets via fields of type ```
AssetRef<EntityPrototype>
```

. This is a way to reference the prototype of an entity to be created via the simulation code whilst having the referencing done similarly to the usual way it is done in Unity, performing a drag and drop or selecting from a list of assets.


Just to exemplify, the screenshot below has an example of referencing via a field already declared on the ```
RuntimePlayer
```

    class:

![Entity Prototype Asset GUID & Path](/docs/img/quantum/v3/manual/assetref-entityprototype.png)
Referencing an Entity Prototype asset in the editor.


4. Use ```
frame.Create()
```

    to create an entity from the prototype. It is most commonly done by passing a reference to the asset shown above, but it also has other overrides:

C#

```csharp
void CreateExampleEntity(Frame frame) {
 // Using a reference to the entity prototype asset
 var exampleEntity = frame.Create(myPrototypeReference);
}

```

### Note

Entity prototypes present in the Scene are baked into the **Map Asset**, while prefab-ed entity prototypes are individual **assets** that are part of the Quantum Asset DataBase.

## Renaming a Component/Prototype

When renaming a component generated entity prototypes Unity scripts will also change and could cause Unity prefabs to loose the connection to the prototype scripts because their script GUIDs changed.

The ```
FormerlyNamed
```

 attribute can be used to safely rename components:

- add the ```
  FormerlyNamed
  ```

   attribute to start the renaming

  - and the prefix ```
    QPrototype..
    ```

     when using Quantum SDK older than 3.0.3: ```
    FormerlyNamed("QPrototypeOldComponentName")
    ```
- CodeGen will automatically migrate the script GUIDs

Qtn

```cs
\[FormerlyNamed("OldComponentName")\]
component NewComponentName {
}

```

## Quantum Entity View

The ```
QuantumEntityView
```

corresponds to the visual representation of an entity.

In the spirit of data driven design, a Quantum Entity Prototype can either incorporate its ```
View
```

component or point to a separate ```
EntityView
```

Asset.

### Self

To set an entity prototype's view to itself, simply add the ```
QuantumEntityView
```

 component to it.

![Entity Prototype with Entity View](/docs/img/quantum/v3/manual/entityprototype-viewself.png)
 Entity Prototype with "Self" View.

Once the component has been added, the entity prototype script will list \*\*Self\*\* as the value for the \*View\* parameter. This will also create a nested \*Entity View\* \*\*asset\*\* in the same prefab.
![Entity Prototype Asset and ](/docs/img/quantum/v3/manual/entityprototype-viewselfasset.png)
 Entity Prototype Asset and "Self" View Asset.
 ### Separate from Prototype

To set up and link a view separate from the _Entity Prototype_ asset:

1. Add ```
   QuantumEntityView
   ```

    to the GameObject which should represent the view;
2. Create a Prefab of that GameObject;
3. This will create an _Entity View_ **Asset** nested in the prefab.

![Entity Prototype with Entity View](/docs/img/quantum/v3/manual/entityprototype-viewseparate.png)
 Entity Prototype Asset and separate Entity View Asset.


4. Link the _View_ field from an _Entity Prototype_ which does not have a ```
   QuantumEntityView
   ```

    associate with it. Reference the newly created _Entity View Asset_. This can be done via drag-and-drop or using the Unity context search menu.

![Linking an Entity Prototype with a separate Entity View Asset](/docs/img/quantum/v3/manual/entityprototype-viewseparateasset.png)
 Linking an Entity Prototype with a separate Entity View Asset.
 ### Important

For an _Entity View_ to be visible in Unity, the scene has to have a ```
QuantumEntityViewUpdater
```

script.

Back to top

- [Introduction](#introduction)
- [Setting up a Prototype](#setting-up-a-prototype)

  - [Basic](#basic)
  - [Custom Components](#custom-components)
  - [Hierarchy](#hierarchy)

- [Creating/Instantiating a Prototype](#creatinginstantiating-a-prototype)

  - [Baked in the Scene/Map](#baked-in-the-scenemap)
  - [In Code](#in-code)
  - [Note](#note)

- [Renaming a Component/Prototype](#renaming-a-componentprototype)
- [Quantum Entity View](#quantum-entity-view)
  - [Self](#self)
  - [Separate from Prototype](#separate-from-prototype)
  - [Important](#important)