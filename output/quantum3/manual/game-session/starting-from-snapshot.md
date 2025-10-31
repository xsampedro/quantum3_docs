# starting-from-snapshot

_Source: https://doc.photonengine.com/quantum/current/manual/game-session/starting-from-snapshot_

# Starting From Snapshot

Quantum supports promoting a local simulation to an online simulation (SDK 3.0.5+).

## Requirements

Local games, that are later promoted to online games, have to be starting with the final target `SessionRunner.Arguments.PlayerCount`, because the ECS layout needs to match the player count of the online game later created from the snapshot.

The new online simulation will always be started from beginning and without any players. They have to be re-added after the start. For simplicity we suggest to remove the local player from the local simulation before taking the snapshot. Otherwise player related entities like avatars have to be reclaimed when adding the player again on the new online simulation.

When starting the online game from a snapshot `SessionRunner.Arguments.FrameData` has to be set to the serialized frame and `InitialTick` has to be set to `RollbackWindow - 1`. This will configure the server to accept the snapshot as initial state for a new simulation that starts from the tick `RollbackWindow`.

Because the simulation restarts from the beginning, timers that are based on ticks will be compromised.

## How to Save a Snapshot

The following snippet uses the default runner `QuantumRunner.Default` to access the running local game.

The process is `async` for reading simplicity. It waits for the player removed confirmation before saving the snapshot.

Stop and destroy the runner before continuing with `LoadAndStart()`.

C#

```csharp
using Photon.Deterministic;
using Quantum;
using System;
using System.IO;
using UnityEngine;
public class Foo : QuantumMonoBehaviour
{
  [EditorButton("Save")]
  public async void Save() {
    // Remove the player and wait for the remove confirmation callback.
    var completionSource = new System.Threading.Tasks.TaskCompletionSource<bool>();
    using (QuantumCallback.SubscribeManual<CallbackLocalPlayerRemoveConfirmed>(c => completionSource.TrySetResult(true)))
    using (QuantumCallback.SubscribeManual<CallbackLocalPlayerRemoveFailed>(c => completionSource.TrySetException(new Exception(c.Message)))) {
      QuantumRunner.Default.Game.RemovePlayer();
      await completionSource.Task;
    }
    // Save the snapshot to a file.
    var snapshot = QuantumRunner.Default.Game.Frames.Verified.Serialize(DeterministicFrameSerializeMode.Serialize);
    File.WriteAllBytes(Path.Combine(Application.dataPath, "savegame.quantum"), snapshot);
    // Save the runtime config to a file.
    File.WriteAllText(Path.Combine(Application.dataPath, "config.quantum"), JsonUtility.ToJson(QuantumRunner.Default.Game.Configurations.Runtime, true));
    // Shutdown the runner, the simulation has to be restarted.
    await QuantumRunner.Default.ShutdownAsync();
  }
}

```

## How to Start From a Snapshot

The player that created the snapshot now connects to a new Photon online room using a secret and unique room name. `IsRoomVisible` is false to prevent random matchmaking to add players to this room. The room name is later send to the friendly player that is supposed to join.

Load the snapshot from a file or from a `byte\[\]` member.

Start Quantum using the snapshot as `FrameData` and `RollbackWindow - 1` as `InitialTick`.

Wait for the simulation to start, then add the player, wait for the confirmation and finally invite the friend to join the online game.

There is a known limitation that requires other players to wait for a short delay (~5 seconds) before they can successfully join the game. This delay must be taken into account to avoid potential desynchronization issues.

C#

```csharp
using Photon.Deterministic;
using Photon.Realtime;
using Quantum;
using System;
using System.IO;
using UnityEngine;
public class Foo : QuantumMonoBehaviour
{
  [EditorButton("LoadAndStart")]
  public async void LoadAndStart() {
    // Delete the debug runner if it exists, it interferes with adding players.
    var debugRunner = FindAnyObjectByType<QuantumRunnerLocalDebug>();
    if (debugRunner != null) {
      Destroy(debugRunner);
    }
    // Connect to a room with a secret room name.
    var arguments = new MatchmakingArguments {
      PhotonSettings = PhotonServerSettings.Global.AppSettings,
      MaxPlayers = Quantum.Input.MAX_COUNT,
      // This name has be unique and must be shared with the other players that are expected to join this room
      RoomName = "my secret room name",
      PluginName = "QuantumPlugin",
      AuthValues = new AuthenticationValues(),
      // Don't use for random matchmaking
      IsRoomVisible = false
    };
    var client = await MatchmakingExtensions.ConnectToRoomAsync(arguments);
    // Load the snapshot from file, or just from a member variable.
    byte[] snapshot = default;
    if (File.Exists(Path.Combine(Application.dataPath, "savegame.quantum"))) {
      snapshot = File.ReadAllBytes(Path.Combine(Application.dataPath, "savegame.quantum"));
    }
    // Use global session config
    var sessionConfig = QuantumDeterministicSessionConfigAsset.DefaultConfig;
    // Set and use custom runtime config instead
    var runtimeConfig = JsonUtility.FromJson<RuntimeConfig>(File.ReadAllText(Path.Combine(Application.dataPath, "config.quantum")));
    // Start and wait for the game, if snapshot is set, it will start from that snapshot.
    // Initial tick does not have to be set explicitly, the simulation will always restart from the beginning.
    var sessionRunnerArguments = new SessionRunner.Arguments {
      FrameData = snapshot,
      RunnerFactory = QuantumRunnerUnityFactory.DefaultFactory,
      GameParameters = QuantumRunnerUnityFactory.CreateGameParameters,
      ClientId = "client secret",
      PlayerCount = Quantum.Input.MaxCount,
      GameMode = DeterministicGameMode.Multiplayer,
      RuntimeConfig = runtimeConfig,
      SessionConfig = sessionConfig,
      Communicator = new QuantumNetworkCommunicator(client),
    };
    var runner = (QuantumRunner)await SessionRunner.StartAsync(sessionRunnerArguments);
    // Add player back to the game and wait for the confirmation.
    var completionSource = new System.Threading.Tasks.TaskCompletionSource<bool>();
    using (QuantumCallback.SubscribeManual<CallbackLocalPlayerAddConfirmed>(c => completionSource.TrySetResult(true)))
    using (QuantumCallback.SubscribeManual<CallbackLocalPlayerAddFailed>(c => completionSource.TrySetException(new Exception(c.Message)))) {
      runner.Game.AddPlayer(0, new RuntimePlayer());
      await completionSource.Task;
    }
    // Invite other player to the room
    Debug.Log($"Invite to room {client.CurrentRoom.Name} on region {client.CurrentRegion}");
  }
}

```

Back to top

- [Requirements](#requirements)
- [How to Save a Snapshot](#how-to-save-a-snapshot)
- [How to Start From a Snapshot](#how-to-start-from-a-snapshot)