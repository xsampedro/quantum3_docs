# commands

_Source: https://doc.photonengine.com/quantum/current/manual/commands_

# Commands

## Introduction

Quantum Commands are an alternative to sending data to the Quantum simulation other than using Inputs. Commands are similar to Quantum Inputs but **are not** required to be sent every tick and can be triggered in specific situations instead.

Quantum Commands are fully reliable. By default, the server will _always_ accept them and confirm it, regardless of the time at which they are sent. A command is locally included in predicted frames and will be quickly executed on the local machine. But for remote clients there is a trade-off; remote clients cannot predict the tick in which the command is received by the simulation, so there is a delay until the Command is received.

Commands are implemented as regular C# classes that inherit from ```
Photon.Deterministic.DeterministicCommand
```

. They can contain any serializable data.

C#

```csharp
namespace Quantum
{
using Photon.Deterministic;

public class CommandSpawnEnemy : DeterministicCommand
{
public AssetRefEntityPrototype EnemyPrototype;

public override void Serialize(BitStream stream)
{
stream.Serialize(ref EnemyPrototype);
}

public void Execute(Frame frame)
{
frame.Create(EnemyPrototype);
}
}
}

```

## Commands Setup in the Simulation

With the Command class defined, it now needs to be registered into the ```
DeterministicCommandSetup
```

's factories. Navigate to ```
Assets/QuantumUser/Simulation
```

and open the script ```
CommandSetup.User.cs
```

. Add, the desired commands to the factory as follows:

C#

```csharp
// CommandSetup.User.cs

namespace Quantum {
 using System.Collections.Generic;
 using Photon.Deterministic;

 public static partial class DeterministicCommandSetup {
 static partial void AddCommandFactoriesUser(ICollection<IDeterministicCommandFactory> factories, RuntimeConfig gameConfig, SimulationConfig simulationConfig) {
 // user commands go here
 // new instances will be created when a FooCommand is received (de-serialized)
 factories.Add(new FooCommand());

 // BazCommand instances will be acquired from/disposed back to a pool automatically
 factories.Add(new DeterministicCommandPool<BazCommand>());
 }
 }
}

```

## Sending Commands From The View

Commands can be sent from anywhere inside Unity.

C#

```csharp
namespace Quantum
{
 using UnityEngine;

 public class EnemySpawnerUI : MonoBehaviour
 {
 \[SerializeField\] private AssetRefEntityPrototype \_enemyPrototype;

 public void SpawnEnemy()
 {
 CommandSpawnEnemy command = new CommandSpawnEnemy()
 {
 EnemyPrototype = \_enemyPrototype,
 };
 QuantumRunner.Default.Game.SendCommand(command);
 }
 }
}

```

### Overloads

```
SendCommand()
```

has two overloads.

C#

```csharp
void SendCommand(DeterministicCommand command);
void SendCommand(Int32 player, DeterministicCommand command);

```

Specify the player index (PlayerRef) if multiple players are controlled from the same machine. Games with only one local player can ignore the player index field.

## Polling Commands From The Simulation

To receive and handle Commands inside the simulation poll the frame for a specific player:

C#

```csharp
using Photon.Deterministic;
namespace Quantum
{
 public class PlayerCommandsSystem : SystemMainThread
 {
 public override void Update(Frame frame)
 {
 for (int i = 0; i < f.PlayerCount; i++)
 {
 var command = frame.GetPlayerCommand(i) as CommandSpawnEnemy;
 command?.Execute(frame);
 }
 }
 }
}

```

### Note

The API does neither enforce, nor implement, a specific callback mechanism or design pattern for Commands. It is up to the developer to chose how to consume, interpret and execute Commands; for example by encoding them into signals, using a Chain of Responsibility, or implementing the command execution as a method in them.

## Examples for Collections

### List

C#

```csharp
namespace Quantum
{
 using System.Collections.Generic;
 using Photon.Deterministic;

 public class ExampleCommand : DeterministicCommand
 {
 public List<EntityRef> Entities = new List<EntityRef>();

 public override void Serialize(BitStream stream)
 {
 var count = Entities.Count;
 stream.Serialize(ref count);
 if (stream.Writing)
 {
 foreach (var e in Entities)
 {
 var copy = e;
 stream.Serialize(ref copy.Index);
 stream.Serialize(ref copy.Version);
 }
 }
 else
 {
 for (int i = 0; i < count; i++)
 {
 EntityRef readEntity = default;
 stream.Serialize(ref readEntity.Index);
 stream.Serialize(ref readEntity.Version);
 Entities.Add(readEntity);
 }
 }
 }
 }
}

```

### Array

C#

```csharp
namespace Quantum
{
 using Photon.Deterministic;

 public class ExampleCommand : DeterministicCommand
 {
 public EntityRef\[\] Entities;

 public override void Serialize(BitStream stream)
 {
 stream.SerializeArrayLength(ref Entities);
 for (int i = 0; i < Cars.Length; i++)
 {
 EntityRef e = Entities\[i\];
 stream.Serialize(ref e.Index);
 stream.Serialize(ref e.Version);
 Entities\[i\] = e;
 }
 }
 }
}

```

## Compound Command Example

Only one command can be attached to an input stream per tick. Even though a client can send multiple Deterministic Commands per tick, the commands will not reach the simulation at the same tick, rather they will arrive separately on consecutive ticks. To go around this limitation, is it possible to pack multiple Deterministic Commands into a single ```
CompoundCommand
```

, which is provided by the SDK.

Instantiating and sending compound commands from the View:

C#

```csharp
var compound = new Quantum.Core.CompoundCommand();
compound.Commands.Add(new FooCommand());
compound.Commands.Add(new BazCommand());

QuantumRunner.Default.Game.SendCommand(compound);

```

Intercepting compound commands:

C#

```csharp
public override void Update(Frame frame) {
for (var i = 0; i < frame.PlayerCount; i++) {
var compoundCommand = frame.GetPlayerCommand(i) as CompoundCommand;
if (compoundCommand != null) {
foreach (var cmd in compoundCommand.Commands) {
// execute individual commands logic
}
}
}
}

```

Back to top

- [Introduction](#introduction)
- [Commands Setup in the Simulation](#commands-setup-in-the-simulation)
- [Sending Commands From The View](#sending-commands-from-the-view)

  - [Overloads](#overloads)

- [Polling Commands From The Simulation](#polling-commands-from-the-simulation)

  - [Note](#note)

- [Examples for Collections](#examples-for-collections)

  - [List](#list)
  - [Array](#array)

- [Compound Command Example](#compound-command-example)