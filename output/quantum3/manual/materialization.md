# materialization

_Source: https://doc.photonengine.com/quantum/current/manual/materialization_

# Materialization

## Introduction

The process of creating an entity or component instance from a ```
Component Prototype
```

 or ```
Entity Prototype
```

is called **Materialization** .

The materialization of scene prototypes baked into the map asset follow the same rules and execution flow as the materialization of instances created by code using the ```
Frame.Create
```

 API.

## Prototype vs Instance

The component instances and entity instances are part of the game state; in other words they can be manipulated at runtime. Components declared in the DSL are used to generate their corresponding ```
Component Prototypes
```

. The code generated prototypes follow the naming convention ```
MyComponentPrototype
```

.

```
Component Prototypes
```

and ```
Entity Prototypes
```

are both **assets**; this means they are not part of the game state, they are immutable at runtime and have to be identical for all clients at all time. Each ```
Component Prototype
```

has a ```
ComponentPrototypeRef
```

which can be used to find it the corresponding asset using the ```
Frame.FindPrototype<MyComponentNamePrototype>(MyComponentPrototypeRef)
```

.

## Component Prototypes

It is possible to extend a ```
Component Prototype
```

to include data which may not be directly used in materialization. This allows, for example, to have shared data between instances of a particular component or exclude read-only data from the frame to keep the game state slim.

Code generated ```
Component Prototypes
```

 are partial classes which can be easily extended:

1. Create a C# file called ```
   MyComponentNamePrototype.Partial.cs
   ```

   ;
2. Place the body of the script into the ```
   Quantum.Prototypes
   ```

    namespace;

It is then possible to add extra data to the ```
Component Prototype
```

asset and implement the partial ```
MaterializeUser()
```

method to add custom materialization logic.

If a component prototype needs to have an additional Unity prototype adapter being generated, it will **not** be emit as partial by default. For a workaround, please check [Unity prototype adapters](#unity-prototype-adapters) section.

### Example

The following example presents the materialization of the ```
Vehicle
```

component as found in the **Arcade Racing Template**.

The ```
Vehicle
```

 component holds mainly dynamic values computed at runtime. Since, per design choice, these variables should not be initialized in the Untiy Editor, the component definition in the DSL uses the ```
ExcludeFromPrototype
```

attribute on those parameters to exclude them from the ```
VehiclePrototype
```

 asset designers can manipulate in the Unity editor. The ```
Nitro
```

parameter is only part that can be edited to allow designers to decide with how much nitro a specific ```
Vehicle
```

 is initialized.

Qtn

```cs
component Vehicle
{
\[ExcludeFromPrototype\]
ComponentPrototypeRef Prototype;

\[ExcludeFromPrototype\]
Byte Flags;
\[ExcludeFromPrototype\]
FP Speed;
\[ExcludeFromPrototype\]
FP ForwardSpeed;
\[ExcludeFromPrototype\]
FPVector3 EngineForce;
\[ExcludeFromPrototype\]
FP WheelTraction;

\[ExcludeFromPrototype\]
FPVector3 AvgNormal;

\[ExcludeFromPrototype\]
array<Wheel>\[4\] Wheels;

FP Nitro;
}

```

The ```
VehiclePrototype
```

asset is extended to provide designers with customizable read-only parameters. The ```
VehiclePrototype
```

asset can thus hold shared values for all instances of a specific vehicle entity prototype "type". The ```
Prototype
```

parameter in the ```
Vehicle
```

component is of type ```
ComponentPrototypeRef
```

which is the component specific equivalent to ```
AssetRef
```

. To populate it, the partial ```
MaterializeUser()
```

method is used to assign the reference of the ```
VehiclePrototype
```

.

C#

```csharp
using Photon.Deterministic;
using Quantum.Inspector;
using System;

namespace Quantum.Prototypes
{
public unsafe partial class VehiclePrototype
{
// PUBLIC METHODS

\[Header("Engine")\]
public FP EngineForwardForce = 130;
public FP EngineBackwardForce = 120;
public FPVector3 EngineForcePosition;
public FP ApproximateMaxSpeed = 20;

\[Header("Hand Brake")\]
public FP HandBrakeStrength = 10;
public FP HandBrakeTractionMultiplier = 1;

\[Header("Resistances")\]
public FP AirResistance = FP.\_0\_02;
public FP RollingResistance = FP.\_0\_10 \* 6;
public FP DownForceFactor = 0;
public FP TractionGripMultiplier = 10;
public FP AirTractionDecreaseSpeed = FP.\_0\_50;

\[Header("Axles")\]
public AxleSetup FrontAxle = new AxleSetup();
public AxleSetup RearAxle = new AxleSetup();

\[Header("Nitro")\]
public FP MaxNitro = 100;
public FP NitroForceMultiplier = 2;

// PARTIAL METHODS
partial void MaterializeUser(Frame frame, ref Vehicle result, in PrototypeMaterializationContext context)
{
result.Prototype = context.ComponentPrototypeRef;
}

\[Serializable\]
public class AxleSetup
{
public FPVector3 PositionOffset;
public FP Width = 1;
public FP SpringForce = 120;
public FP DampingForce = 175;
public FP SuspensionLength = FP.\_0\_10 \* 6;
public FP SuspensionOffset = -FP.\_0\_25;
}
}
}

```

The parameters in the ```
VehiclePrototype
```

hold values necessary to compute the dynamic values found in the component instance which impact the behaviour of the entity to which the ```
Vehicle
```

component is attached. For example, when a player picks up additional ```
Nitro
```

, the value held in the ```
Vehicle
```

component is clamped to the ```
MaxNitro
```

value found in the ```
VehiclePrototype
```

. This enforces the limits under penality of desynchronization and keeps the game state slim.

C#

```csharp
namespace Quantum
{
public unsafe partial struct Vehicle
{
public void AddNitro(Frame frame, EntityRef entity, FP amount)
{
var prototype = frame.FindPrototype<Vehicle\_Prototype>(Prototype);
Nitro = FPMath.Clamp(Nitro + amount, 0, prototype.MaxNitro);
}
}
}

```

## Materialization Order

Every ```
Entity Prototype
```

's materialization, including the scene prototypes, executes the following steps in order:

1. An empty entity is created.
2. For each ```
Component Prototype
```

    contained in the ```
Entity Prototype
```

:

1. the component instance is created on the stack;
2. the ```
      Component Prototype
      ```

       is materialized into the component instance;
3. ```
      MaterializeUser()
      ```

       is called (though implementing this is _optional_) ; and,
4. the component is added to the entity which triggers the ```
      ISignalOnComponentAdded<MyComponent>
      ```

       signal.
3. ```
ISignalOnEntityPrototypeMaterialized
```

    is invoked for each materialized entity.

   - Load Map / Scene: the signal is invoked for all entity & ```
     Entity Prototype
     ```

      pair after all scene prototypes have been materialized.
   - Created with ```
     Frame.Create()
     ```

     : the signal is invoked immediately after the prototype has been materialized.

The ```
Component Prototype
```

materialization step materializes default components in a predetermined order.

C#

```csharp
Transform2D
Transform3D
Transform2DVertical
PhysicsCollider2D
PhysicsBody2D
PhysicsCollider3D
PhysicsBody3D
PhysicsJoints2D
PhysicsJoints3D
PhysicsCallbacks2D
PhysicsCallbacks3D
CharacterController2D
CharacterController3D
NavMeshPathfinder
NavMeshSteeringAgent
NavMeshAvoidanceAgent
NavMeshAvoidanceObstacle
View
MapEntityLink

```

Once all default components have been materialized, the user defined components are materialized in alphabetically order.

C#

```csharp
MyComponentAA
MyComponentBB
MyComponentCC
...

```

## Unity Prototype Adapters

If any field of a component uses ```
ReplaceTypeHintAttribute
```

or is one of following types:

- ```
EntityRef
```

- ```
EntityPrototypeRef
```

- ```
ComponentPrototypeRef
```

- ```
ComponentPrototypeRef<T>
```


Then Quantum will generate an additional component prototype adapter type. Adapter types are placed in ```
Quantum.Prototypes.Unity
```

namespace and have the same fields as their source prototypes, with an exception of following type replacements taking place:

- ```
EntityRef
```

   -\> ```
QuantumEntityPrototype
```

- ```
EntityPrototypeRef
```

   -\> ```
QUnityEntityPrototypeRef
```

- ```
ComponentPrototypeRef
```

   -\> ```
QUnityComponentPrototypeRef
```

- ```
ComponentPrototypeRef<T>
```

   -\> ```
QUnityComponentPrototypeRef<T>
```

- ```
\[ReplaceTypeHintAttribute\]
```

argument used instead of the field type

If an adapter is generated, Unity MonoBehaviour-based prototype wrappers will use it instead of the prototype. That way prototypes appear to work with Unity Object references, while still being simulation compatible.

Most of the time this is completely transparent to user. However, this process has a side effect of preventing component prototype from being emit as partial. The reason is the adapter's fields and conversion should be in sync with the source prototype, but once partial classes are involved, the code generation has no way of knowing what has been added. ```
\[CodeGen(ForcePartialPrototype)\]
```

attribute can be used to enforce partial component prototype, but make sure to implement the adapter's ```
ConvertUser
```

method in such case.

Back to top

- [Introduction](#introduction)
- [Prototype vs Instance](#prototype-vs-instance)
- [Component Prototypes](#component-prototypes)

  - [Example](#example)

- [Materialization Order](#materialization-order)
- [Unity Prototype Adapters](#unity-prototype-adapters)