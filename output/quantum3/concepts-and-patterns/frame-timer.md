# frame-timer

_Source: https://doc.photonengine.com/quantum/current/concepts-and-patterns/frame-timer_

# Frame Timer

The Quantum `FrameTimer` struct can be used to implement gameplay based on a timer. Like cooldown for skills or respawn or to run periodically checks in the simulation.

Start a new timer from seconds or from ticks.

C#

```csharp
FrameTimer timerA = FrameTimer.FromFrames(frame, 60)
FrameTimer timerB = FrameTimer.FromSeconds(frame, 2)
FrameTimer timerC = default(FrameTimer);

```

The timers _timerA_ and _timerB_ will be in running state, while _timerC_ is not set, yet.

Checking the state of a FrameTimer is exemplified below, considering the sample timers from above:

C#

```csharp
timerA.IsRunning(frame); // returns TRUE
timerB.IsRunning(frame); // returns TRUE
timerC.IsRunning(frame); // returns FALSE
timerA.IsSet; // returns TRUE
timerB.IsSet; // returns TRUE
timerC.IsSet; // returns FALSE

```

The `FrameTimer` keeps a memory of the target frame, which can be used to query if a timer stopped at an exact time. It can also return the elapsed and remaining time which can be handy to be displayed in the game UI for example.

C#

```csharp
if (timer.HasStoppedThisFrame(frame)) {
    // Only if the timer ran out this exact frame.
}
else {
    var ticksRemaining = timer.RemainingFrames(frame);
    var secondsRemaining = timer.RemainingSeconds(frame);
    var ticksElapsed = timer.ElapsedFrames(frame);
    var elapsedSeconds = timer.ElapsedSeconds(frame);
}

```

`FrameTimers` can be added to components in the DSL or it can be added to the Frame's global variables as exemplified below.

Qtn

```cs
component Character {
    FrameTimer SkillCooldown;
}
global {
    FrameTimer GameTimer;
}

```

This example uses the `FrameTimer` on the component to control a skill cooldown. The same logic could also be applied to the `FrameTimer` in global variables.

C#

```csharp
namespace Quantum
{
  using Photon.Deterministic;
  using UnityEngine.Scripting;

  [Preserve]
  public unsafe class CharacterSystem : SystemMainThreadFilter<CharacterSystem.Filter> {
    public struct Filter {
      public EntityRef Entity;
      public Character* Character;
    }
    public override void Update(Frame frame, ref Filter filter)
        var character = filter.Character;
        if (character->SkillCooldown.IsRunning(frame) == false) {
            // Execute skill and reset the cooldown
            // Reset the cooldown by creating a new timer
            character->SkillCooldown = FrameTimer.FromSeconds(frame, 2);
            // Reset the cooldown by restarting the timer.
            // BUT this requires that the timer was set at a previous time, for example when adding the component.
            character->SkillCooldown.Restart(frame);
        }
    }
}

```

Back to top