# input

_Source: https://doc.photonengine.com/quantum/current/manual/input_

# Input

## Introduction

Input is a crucial component of Quantum's core architecture. In a deterministic networking library, the output of the system is fixed and predetermined given a certain input. This means that as long as the input is the same across all clients in the network, the output will also be the same.

## Defining in DSL

Input can be defined in any [DSL](/quantum/current/manual/quantum-ecs/dsl) file. For example, an input struct where you have a movement direction and a singular jump button would look something like this:

Qtn

```cs
input
{
 button Jump;
 FPVector3 Direction;
}

```

The server is responsible for batching and sending down input confirmations for full tick-sets (all player's input). For this reason, this struct should be kept to a minimal size as much as possible.

## Commands

[Deterministic Commands](/quantum/current/manual/commands) are another input path for Quantum, and can have arbitrary data and size, which make them ideal for special types of inputs, like "buy this item", "teleport somewhere", etc.

## Polling in Unity

To send input to the Quantum simulation, poll for it inside of Unity. To do this, subscribe to the ```
PollInput
```

 callback inside of a MonoBehaviour in the gameplay scene.

C#

```csharp
private void OnEnable()
{
QuantumCallback.Subscribe(this, (CallbackPollInput callback) => PollInput(callback));
}

```

Then, in the callback, read from the input source and populate the input struct.

C#

```csharp
public void PollInput(CallbackPollInput callback)
{
Quantum.Input i = new Quantum.Input();

var direction = new Vector3();
direction.x = UnityEngine.Input.GetAxisRaw("Horizontal");
direction.y = UnityEngine.Input.GetAxisRaw("Vertical");

i.Jump = UnityEngine.Input.GetKey(KeyCode.Space);

// convert to fixed point.
i.Direction = direction.ToFPVector3();

callback.SetInput(i, DeterministicInputFlags.Repeatable);
}

```

NOTE: The float to fixed point conversion here is deterministic because it is done before it is shared with the simulation.

## Optimization

Although Quantum 3 delta-compresses input, it is still generally a good practice to make the raw ```
Input
```

data as compact as possible for optimal bandwidth. Below are a few ways to optimize it.

## Buttons

Instead of using booleans or similar data types to represent key presses, the ```
Button
```

 type is used inside the Input DSL definition. This is because it only uses one bit per instance, so it is favorable to use where possible. Although they only use one bit over the network, locally they will contain a bit more game state. This is because the single bit is only representative of whether or not the button was pressed during the current frame, the rest of the information is computed locally.

Buttons are defined as follows:

Qtn

```cs
input
{
button Jump;
}

```

An important detail when polling the values for a button from a Unity script is to poll the _current button state_, i.e whether if it is pressed or not at the current frame. With this, Quanutm automatically sets up internal properies which allows the user to, _on the simulation code_, poll for specific states such as ```
WasPressed
```

, ```
IsDown
```

and ```
WasReleased
```

.

This means that, in Unity, you do not need to set specific states such as ```
GetKeyUp()
```

or ```
GetKeyDown()
```

as using these would actually be problematic as Unity does not run at the same rate as Quantum, so some of these states would be lost and make the input feel less responsive.

So, when setting the value of a ```
button
```

 in the Input structure, always poll for the current button state as shown below:

C#

```csharp
// In Unity, when polling a player's input
input.Jump = UnityEngine.Input.GetKey(KeyCode.Space);

```

The state of the button can also be updated in Quantum simulation code, which is particulaly useful for simulating changes on a button for non-player entities such as bots, if the user chooses to also update bot entities by using the Input struct and the ```
button
```

type. To achieve this, it is necessary to set the state of the button in simulation code _every frame_ as shown below:

C#

```csharp
// In Quantum code
input.button.Update(frame, value);

```

This way, the specific states (Pressed, Down, Released) is also internally generated. Not updating the button state every frame results in those states being wrongly set.

## Encoded Direction

In a typical setting, movement is often represented using a direction vector, often defined in a ```
DSL
```

 file as such:

Qtn

```cs
input
{
FPVector2 Direction;
}

```

However, ```
FPVector2
```

is comprised of two 'FP', which takes up 16 bytes of data, which can be a lot of data sent, especially with many clients in the same room.

One such way of optimizing it, is by extending the ```
Input
```

struct and encoding the directional vector into a ```
Byte
```

instead of sending the full vector every single time. One such implemetation is as follows:

First, we define our input like normal, but instead of including an ```
FPVector2
```

 for direction, we replace it with a ```
Byte
```

where we will store the encoded version.

Qtn

```cs
input
{
 Byte EncodedDirection;
}

```

Next, extend the input struct the same way a component is extended (see: [Adding Functionality](/quantum/current/manual/quantum-ecs/components#adding_functionality)):

C#

```csharp
namespace Quantum
{
 partial struct Input
 {
 public FPVector2 Direction
 {
 get
 {
 if (EncodedDirection == default)
 return default;

 Int32 angle = ((Int32)EncodedDirection - 1) \* 2;

 return FPVector2.Rotate(FPVector2.Up, angle \* FP.Deg2Rad);
 }
 set
 {
 if (value == default)
 {
 EncodedDirection = default;
 return;
 }

 var angle = FPVector2.RadiansSigned(FPVector2.Up, value) \* FP.Rad2Deg;

 angle = (((angle + 360) % 360) / 2) + 1;

 EncodedDirection = (Byte) (angle.AsInt);
 }
 }
 }
}

```

This implementation allows for the same usage as before, but it only takes up a singular byte instead of 16 bytes. It does this by utilzing a ```
Direction
```

 property, which encodes and decodes the value from ```
EncodedDirection
```

automatically.

Back to top

- [Introduction](#introduction)
- [Defining in DSL](#defining-in-dsl)
- [Commands](#commands)
- [Polling in Unity](#polling-in-unity)
- [Optimization](#optimization)
- [Buttons](#buttons)
- [Encoded Direction](#encoded-direction)