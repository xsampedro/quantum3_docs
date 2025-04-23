# input-flags

_Source: https://doc.photonengine.com/quantum/current/manual/player/input-flags_

# Input & Connection Flags

## Introduction

The `DeterministicInputFlags` are used by Quantum to:

- detect whether a player is _present_ , i.e. connected, to the simulation;
- decide how to _predict_ the next tick's input for a given player; and,
- know whether the input on a verified frame was provided by a client or was _replaced_ by the server.

It is possible to automate the checks by implementing `PlayerConnectedSystem`, for more information [please refer to its entry on the Player page](/quantum/current/manual/player/player).

## Types

C#

```csharp
public enum DeterministicInputFlags : byte {
  Repeatable = 1 << 0,
  PlayerNotPresent = 1 << 1,
  ReplacedByServer = 1 << 2,
  Command = 1 << 3
}

```

- `PlayerNotPresent` = means there is no client connected for this player index.
- `ReplacedByServer` = means the player index is controlled by a client, but the client did not send the input in time which resulted in the server repeating or replacing/zeroing out the input.
- `Repeatable` = tells both the server and other clients to copy this input data into the next tick (on server when replacing input due to timeout, and on other clients for the local prediction algorithm). This can be set by the developer from Unity when injecting player input and should be used on direct-control-like input such as movement; it is not meant for command-like input (e.g. buy item).

## Implementation example

**IMPORTANT:**`DeterministicInputFlags` can only be trusted on _verified_ frames.

The code snippet below was extracted from the `PlayerConnectedSystem` provided in the SDK.

C#

```csharp
public unsafe class PlayerConnectedSystem : SystemMainThread {
    public override void Update(Frame frame) {
      // only trustable in Verified frames
      if (frame.IsVerified == false) {
        return;
      }

      for (int p = 0; p < f.PlayerCount; p++) {
        var isPlayerConnected = (frame.GetPlayerInputFlags(p) & Photon.Deterministic.DeterministicInputFlags.PlayerNotPresent) == 0;
        // extra logic based on the player connectivity state
      }
    }
  }
}

```

Back to top

- [Introduction](#introduction)
- [Types](#types)
- [Implementation example](#implementation-example)