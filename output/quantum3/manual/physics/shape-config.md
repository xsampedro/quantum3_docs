# shape-config

_Source: https://doc.photonengine.com/quantum/current/manual/physics/shape-config_

# Shape Config

## Introduction

The shape config holds information on a shape. It can be used to easily expose configuration options in the editor and streamline shape initialization in code.

ShapeConfigs exist for both 2D and 3D - ```
Shape2DConfig
```

 and ```
Shape3DConfig
```

respectively.

Currently, ShapeConfigs are targeted for use with dynamic entities; as such, the shape types supported by it are:

// in 2D

- Circle
- Box
- Polygon
- Edge
- Compound (a combination of any of these)

// in 3D

- Sphere
- Box
- Capsule
- Compound (a combination of any of these)

All of the aforementioned shapes are compatible with _ShapeOverlaps_ and _PhysicsColliders_.

## Exposing ShapeConfig to the Editor

Including either a ```
Shape2DConfig
```

 or ```
Shape3DConfig
```

in a custom asset will automatically expose it in the Editor.

C#

```csharp
namespace Quantum
{
 public unsafe partial class WeaponSpec
 {
 public Shape3DConfig AttackShape;
 public LayerMask AttackLayers;
 public FP Damage;
 public FP KnockbackForce;
 }
}

```

The asset resulting from the snippet above will offer the following exposing all ShapeConfigs options in the inspector:

![ShapeConfig Setting of the example asset as shown in the Unity Editor](/docs/img/quantum/v3/manual/physics/physics-shapeconfig-editor.png)
ShapeConfig Setting of the WeaponSpec example asset as shown in the Unity Editor.
## Creating/Using a Shape from ShapeConfig

When using a ShapeConfig, call its ```
CreateShape
```

 method. This will automatically process the information held in the ShapeConfig asset, create the appropriate shape and its parameter with the data found in the config.

C#

```csharp
private static void Attack(in Frame frame, in EntityRef entity)
{
// A melee attack performed by using an OverlapShape on the attack area.

var transform = frame.Unsafe.GetPointer<Transform3D>(entity);
var weapon = frame.Unsafe.GetPointer<Weapon>(entity);
var weaponSpec = frame.FindAsset<WeaponSpec>(weapon->WeaponSpec.Id);

var hits = frame.Physics3D.OverlapShape(
transform->Position,
transform->Rotation,
weaponSpec.AttackShape.CreateShape(frame),
weaponSpec.AttackLayers);

// Game logic iterating over the hits.
}

```

The same can be done when initializing the Shape of a PhysicsCollider.

## Compound Shape

A Compound Shape is a shape made of several other shapes. The shapes that can be used to create a compound shape are the ones listed in the introduction section.

PhysicsColliders and ShapeOverlaps are fully compatible with compound shapes.

As of now, Quantum offers persistent compound shapes, i.e. a shape pointing to a buffer of other shapes in the heap. This buffer will persist between frames until it is manually disposed.

### Create New Compound Shape

A compound shape can be create from a ShapeConfig by simply calling the ```
CreateShape
```

method, or manually via the ```
Shape.CreatePersistentCompound
```

, and later dispose of it calling ```
FreePersistent
```

on. Here is an example of how the lifetime can be managed, the same applies to 3D Shapes:

C#

```csharp
 // creating a persistent compound. This does not allocate memory until actually adding shapes
 var compoundShape = Shape2D.CreatePersistentCompound();

 // adding shapes to a compound (shape1 and 2 can be of any type)
 compoundShape.Compound.AddShape(f, shape1);
 compoundShape.Compound.AddShape(f, shape2);

 (...) // Game logic

 // this compound persists until it is manually disposed
 compoundShape.Compound.FreePersistent(f);

```

The API also offers methods such as ```
RemoveShapes
```

, ```
GetShapes
```

, and ```
FreePersistent
```

; check the API documentation in the SDK's ```
docs
```

folder for more information on those and other methods.

### CopyFrom an Existing Compound Shape

You can also create a new compound shape by copying an existing one.

C#

```csharp
// Using the exising compoundShape from the example above.

var newCompoundShape = Shape2D.CreatePersistentCompound();
newCompoundShape.Compound.CopyFrom(f, ref oldCompoundShape);

```

Creating a new compound shape also means a new buffer. This will need to be freed by the developer manually like any other persistent compound shape.

### Accesing individual Shapes

An example of how to iterate through the Shapes is to use the ```
GetShapes
```

 method to get the pointers buffer, and use a simple ```
for
```

loop, where the integer index can be used to access all ```
Shape\*
```

 contained in the compound.

Do not surprass the ```
count
```

returned by the method as it is the boundary where the shape pointers are contained in memory.

C#

```csharp
if (shape->Compound.GetShapes(frame, out Shape3D\* shapesBuffer, out int count))
{
 for (var i = 0; i < count; i++)
 {
 Shape3D\* currentShape = shapesBuffer + i;
 // do something with the shape
 }
}

```

### Compound Collider

A _compound_ collider is a regular collider with a shape of type compound.

#### Create in Editor

In the editor, find the options to create a compound collider in the ```
Entity Component Physics Collider 2D/3D
```

 and the ```
Entity Prototype
```

script's respective section.

![ShapeConfig for a Compound Collider in the Entity Prototype script as shown in the Unity Editor](/docs/img/quantum/v3/manual/physics/physics-compound-shape.png)
Creating a Compound Collider (Sphere + Box) via the Entity Prototype script in the Unity Editor.


When creating a collider prototype with a compound shape, the memory management is already handled, i.e. the collider will have and manage its compound shape and it is not needed to manually dispose anything.

#### Create in Code

When creating a collider in code, simply pass a compound shape into its ```
Create()
```

 factory method. Once the compound shape has been created as shown in the code snippet from the previous section, it is possible to create a collider by replacing ```
(...)
```

with this:

C#

```csharp
 var collider = PhysicsCollider2D.Create(f, compoundShape);
 f.Set(entity, collider);

```

In the code snippet provided above ```
collider.Shape
```

 and ```
compoundShape
```

point to different buffers in the heap. When done using the compound shape - i.e. it was only needed for creating the collider - it is possible to dispose it right after that. The collider will dispose its own copy in memory when destroyed/removed.

#### Important Note about Memory

A collider only creates a copy of the compound shape buffer if it is used as part of its factory method ```
Create()
```

.

C#

```csharp
var compoundShape = Shape2D.CreatePersistentCompound();
compoundShape.Compound.AddShape(f, shape1);
compoundShape.Compound.AddShape(f, shape2);

// collider1 and collider2 each create a copy of the compoundShape buffer.
// collider1 and collider2 will each dispose of their copy on destroy/remove.
var collider1 = PhysicsCollider2D.Create(f, compoundShape);
f.Set(entity1, collider1);

var collider2 = PhysicsCollider2D.Create(f, compoundShape);
f.Set(entity2, collider2);

// Here we dispose of the compoundShape's buffer as it is no longer needed
compoundShape.Compound.FreePersistent(f);

```

In contrast, if creating a collider with a regular shape is needed, and later set its ```
collider.Shape = someCompound
```

it will not create a copy of the buffer, i.e. ```
collider.Shape
```

and ```
someCompound
```

will point to the same buffer. Doing this can be dangerous if when having multiple compound colliders and/or compound shapes pointing towards the same buffer. If one disposes of it, it will effectively break the others' reference to the buffer.

C#

```csharp
 var compoundShape = Shape2D.CreatePersistentCompound();
 compoundShape.Compound.AddShape(f, shape1);
 compoundShape.Compound.AddShape(f, shape2);

 var collider1 = PhysicsCollider2D.Create(f, default(Shape2D));
 collider1.Shape = compoundShape;
 f.Set(entity1, collider1);

 var collider2 = PhysicsCollider2D.Create(f, Shape2D.CreateCircle(1));
 collider2.Shape = compoundShape;
 f.Set(entity2, collider2);

 // collider1.Shape, collider2.Shape and compoundShape all point to the same buffer
 // dispose of compoundShape here will break collider1 and collider2
 compoundShape.Compound.FreePersistent(f);

```

However, when doing it conscientiously this will assign the shape and its memory management to the collider in question thus simplifying memory management by relinquishing the responsibilty to the collider which will dispose of it when destroyed/removed.

C#

```csharp
 var compoundShape = Shape2D.CreatePersistentCompound();
 compoundShape.Compound.AddShape(f, shape1);
 compoundShape.Compound.AddShape(f, shape2);

 var collider1 = PhysicsCollider2D.Create(f, Shape2D.CreateCircle(1));
 collider1.Shape = compoundShape;
 f.Set(entity1, collider1);

 // In this instance we do not need to dispose of the compoundShape buffer because
 // collider1 already points to it and will take care of it on destroy/remove.

```

### Compound Shape Query

Compound shape queries are fully supported for both broadphase and regular queries. The behaviour and performance impact is the same as performing multiple queries, except the results will be returned in the same HitCollection;

### Nested Compound Shape

Nested compound shapes are supported by the physics engines though with two limitations:

- a shape can only hold one reference to the same buffer in the heap in its hierarchy. Add an already referenced buffer will throw an error in debug mode. This is to avoid issues with cyclic references and invalid pointers on disposal.
- nested compound shapes are not supported in the editor (a warning message will be shown). This is due to the Unity serializer limitations which requires a more intricate structure and drawer for the shape config.

Back to top

- [Introduction](#introduction)
- [Exposing ShapeConfig to the Editor](#exposing-shapeconfig-to-the-editor)
- [Creating/Using a Shape from ShapeConfig](#creatingusing-a-shape-from-shapeconfig)
- [Compound Shape](#compound-shape)
  - [Create New Compound Shape](#create-new-compound-shape)
  - [CopyFrom an Existing Compound Shape](#copyfrom-an-existing-compound-shape)
  - [Accesing individual Shapes](#accesing-individual-shapes)
  - [Compound Collider](#compound-collider)
  - [Compound Shape Query](#compound-shape-query)
  - [Nested Compound Shape](#nested-compound-shape)