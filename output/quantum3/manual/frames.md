# frames

_Source: https://doc.photonengine.com/quantum/current/manual/frames_

# Frames

## Introduction

Quantum's predict-rollback architecture allows to mitigate latency. Quantum always rolls-back and re-simulates frames. It is a necessity for determinism and involves the validation of player input by the server. Once the server has either confirmed the player input or overwritten/replaced it (only in cases were the input did not reach the server in time), the validated input of all players for a given frame is sent to the clients. Once the validated input is received, the last verified frame advances using the confirmed input.

**N.B.:** A player's own input will be rolled back if it has not reached the server in time or could not be validated.

## Types of Frame

Quantum differentiates between two types of frame:

- verified; and,
- predicted.

### Verified

A _Verified_ frame is a _trusted_ simulation frame. A verified frame is guaranteed to be deterministic and identical on all client simulations. The verified simulation only simulates the next verified frame once it has received the server-confirmed inputs; as such, it moves forward proportional to RTT/2 from the server.

A frame is verified if both of the following condition are both true:

- the input from **ALL** players is confirmed by the server for this tick; and,
- all previous ticks it follows are verified.

A partial tick confirmation where the input from _only a subset_ of player has been validated by the server will not result in a verified tick/frame.

### Predicted

Contrary to _verified_ frames, _predicted_ frames do not require server-confirmed input. This means the predicted frame advances with prediction as soon as the simulation has accumulated enough delta time in the local session.

The Unity-side API offers access to various versions of the predicted frame, see the API explanation below.

- ```
Predicted
```

: the simulation "head", based on the synchronised clock.
- ```
PredictedPrevious
```

(predicted - 1): used for main clock-aliasing interpolation (most views will use this to stay smooth, as Unity's local clock may slightly drift from the main server clock. Quantum runs from a separate clock, in sync with the server clock - smoothly corrected).
- ```
PreviousUpdatePredicted
```

: this is the exact frame that was the "Predicted/Head" the last time ```
Session.Update
```

was called (with the "corrected" data in it). Used for error correction interpolation (most of the time there will be no error).

## API

The concept of _Verified_ and _Predicted_ frames exists in both the simulation and the view, albeit with a slightly different API.

### Simulation

In the simulation, one can access the state of the currently simulated frame via the ```
Frame
```

 class.

| Method | Return Value | Description |
| --- | --- | --- |
| IsVerified | bool | Returns true if the frame is deterministic across all clients and uses server-confirmed input. |
| IsPredicted | bool | Returns true if the frame is a locally predicted one. |

### View

In the view, the _verified_ and _predicted_ frames are made available via ```
QuantumRunner.Default.Game.Frames
```

.

| Method | Description |
| --- | --- |
| Verified | Trusted simulation frame, identical across all clients. |
| Predicted | The local simulation "head" based on the synced Quantum clock. Can differ between clients. |
| PredictedPrevious | Predicted - 1<br>Used for main clock-aliasing interpolation, most views will use this to stay smooth. As Unity's local clock may slightly drift from the main server clock, Quantum runs from a separate clock which is in sync with the server clock - smoothly corrected |
| PreviousUpdatePredicted | The re-simulated version of the frame that had been the "Predicted/Head" frame when the last time Session.Update was called. This is necessary in case of rollbacks in order to "correct" data held by it. It is used by the View for error-correction in the interpolation - this is a safety measure and rarely ever necessary. |

## Using Frame.User

It is possible to extend the Frame by adding data to ```
Frame.User.cs
```

. However, in doing so it will also be necessary to implement the corresponding initialization, allocation and serialization methods used by the frame.

C#

```csharp
partial void InitUser() // Initialize the Data

partial void SerializeUser(FrameSerializer serializer) // De/Serialize the Data
partial void CopyFromUser(Frame frame) // Copy to next Frame

partial void AllocUser() // Allocate space
partial void FreeUser() // Free allocated space

```

**NOTE**: Adding an excessive amount of data to the frame will impact performance (de/serialization), as well as affect late joins.

### Example

This is a very simple example which does not require manual memory allocation.

C#

```csharp
namespace Quantum {

unsafe partial class Frame {
public byte\[\] Grid => \_grid;
private byte\[\] \_grid;

partial void InitUser() {
\_grid = new byte\[RuntimeConfig.GridSize\];
}

partial void SerializeUser(FrameSerializer serializer)
{
serializer.Stream.SerializeArrayLength<Byte>(ref \_grid);
for (int i = 0; i < Grid.Length; i++)
{
serializer.Stream.Serialize(ref Grid\[i\]);
}
}

partial void CopyFromUser(Frame frame)
{
Array.Copy(frame.\_grid, \_grid, frame.\_grid.Length);
}
}
}

```

Back to top

- [Introduction](#introduction)
- [Types of Frame](#types-of-frame)

  - [Verified](#verified)
  - [Predicted](#predicted)

- [API](#api)

  - [Simulation](#simulation)
  - [View](#view)

- [Using Frame.User](#using-frame.user)
  - [Example](#example)