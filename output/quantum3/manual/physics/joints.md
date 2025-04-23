# joints

_Source: https://doc.photonengine.com/quantum/current/manual/physics/joints_

# Physics Joints

## Introduction

A physics joint is a connection between physics bodies that restrict movement. They are valuable in games as they allow for realistic movement, interaction between physics bodies, and vehicle simulations.

## Configuration

Quantum's physics joint API provides several parameters that can be used to modify their behavior.

| Field | Description |
| --- | --- |
| Start Disabled | Whether or not the joint should start disabled. Disabled joints will be ignored by the physics systems. |
| Type | What type the joint should be. |
| User Tag | A numerical tag that can be used to define groups of joints. |
| Connected Entity | What entity this joint is connected to. This entity MUST have a transform component. |
| Connected Anchor | The anchor point of where the joint entity connects to. If there is no connected entity, this will be in world space, otherwise it will be in local space. |
| Anchor | The anchor offset in local space of the joint entity's transform. |
| Axis | The axis around which the joint rotates. If set to zero, \`FPVector3.Right\` is used instead. |
| Use Angle Limits | If the angle between the joint's transform and its anchor should be limited by the hinge joint. Only relevent if the joint is a hinge joint. |
| Lower Angle | The lower limiting angle of the arc of rotation around the anchor, in degrees. |
| Upper Angle | The upper limiting angle of the arc of rotation around the anchor, in degrees. |
| Use Motor | If the hinge joint should use a motor. Only relevant if the joint is a hinge joint. |
| Motor Speed | The speed at which the motor will attempt to rotate, in angles per second. |
| Max Motor Torque | The maximum torque the hinge joint can provide, if set to 0 then it will not be limited. |

### Hinge Joints

A ```
Hinge Joint
```

 connects two physics bodies together at a point, allowing them to rotate around a specified axis. These joints can be useful for creating doors, gates or rotating obstacles. Hinge Joints also have optional motors which will constantly try to rotate in a specified direction. To use a motor in the hinge joint, set ```
UseMotor
```

to true in the joint component.

### Spring Joints

A ```
Spring Joint
```

 keeps bodies apart and maintains a separation between them, but allows some stretching in the distance between them. They can be useful for creating complex character movement or enviroment interactions like bridges that respond to player movement.

### Distance Joints

A ```
Distance Joint
```

is a type of physics joint used to maintain a fixed distance between two objects. It restricts the relative motion of the connected objects along a straight line, keeping them at a specific distance from each other. This joint type is useful for simulating rigid connections, such as ropes and virtual springs. Generally useful for any joint where the distance between the objects remains constant. It allows for realistic interactions, stability, and control over connected objects in the simulation.

Back to top

- [Introduction](#introduction)
- [Configuration](#configuration)
  - [Hinge Joints](#hinge-joints)
  - [Spring Joints](#spring-joints)
  - [Distance Joints](#distance-joints)