# pausing

_Source: https://doc.photonengine.com/quantum/current/concepts-and-patterns/pausing_

This page is a work in progress and could be pending updates.


# Pausing

Pausing is rarely utilized in online games due to it disrupting the experience for other players; however, in many games that offer local gameplay, the ability to pause is often expected, allowing players to change settings on the fly or even just to take a break. Due to how Quantum works, more traditional methods for pausing gameplay cannot be applied. This page will go over a few different methods, citing which are and are not recommended.

## Disabling The Quantum Runner

For local games, the most straightforward way to pause a game is to either disable the `QuantumRunnerBehaviour` or set `QuantumRunner.IsSessionUpdateDisabled` to `true`. These will only work, however, if the following two conditions are met:

- The `DeterministicGameMode` is NOT `DeterministicGameMode.Multiplayer`
- The `SimulationUpdateTime` is NOT set to `SimulationUpdateTime.Default`

Again, this will only work for local games because neither of the previously discussed values are not deterministic; if attempted during an online game, the frames at which the games are paused and unpaused, could be different between clients and cause desyncs. The reason it's important that the `SimulationUpdateTime` is not using `SimulationUpdateTime.Default` is because with the default setting, the simulation will continue to progress. The game will appear paused, but skip ahead to the correct frame when resuming.

The following is an example of how to pause the game by toggling the `QuantumRunnerBehaviour`:

C#

```csharp
public void OnTogglePause()
{
    // Checks to see if there is a QuantumRunner and if the QuantumRunner has a UnityObject associated with it.
    if (QuantumRunner.Default == null || QuantumRunner.Default.UnityObject == null)
        return;
    if (!QuantumRunner.Default.UnityObject.TryGetComponent(out QuantumRunnerBehaviour runnerBehaviour))
        return;
    runnerBehaviour.enabled = !runnerBehaviour.enabled;
}

```

And this is how it could be done by toggling `QuantumRunner.IsSessionUpdateDisabled`

C#

```csharp
public void OnTogglePause()
{
if (QuantumRunner.Default == null)
        return;
    QuantumRunner.Default.IsSessionUpdateDisabled = !QuantumRunner.Default.IsSessionUpdateDisabled;
}

```

Note, if testing a game that utilizes `QuantumRunnerLocalDebug`, the latter method will not work out of the box due to the `Update` method in this class:

C#

```csharp
/// <summary>
/// Unity update event. Will update the simulation if a custom <see cref="SimulationSpeedMultiplier" /> was set.
/// </summary>
public void Update()
{
  if (QuantumRunner.Default != null && QuantumRunner.Default.Session != null)
  {
    QuantumRunner.Default.IsSessionUpdateDisabled = SimulationSpeedMultiplier != 1.0f;
    if (QuantumRunner.Default.IsSessionUpdateDisabled)
    {
       switch (QuantumRunner.Default.DeltaTimeType)
       {
            case SimulationUpdateTime.Default:
            case SimulationUpdateTime.EngineUnscaledDeltaTime:
              QuantumRunner.Default.Service(Time.unscaledDeltaTime * SimulationSpeedMultiplier);
              QuantumUnityDB.UpdateGlobal();
              break;
            case SimulationUpdateTime.EngineDeltaTime:
              QuantumRunner.Default.Service(Time.deltaTime);
              QuantumUnityDB.UpdateGlobal();
              break;
       }
    }
  }
}

```

## Recommended: Disabling and Re-Enabling Systems

If pausing a game that is online is important, say for an online training mode in a game or a more casual, cooperative game, then the recommended approach to pausing is to disable and then re-enable various `Systems` that affect gameplay. Unlike the previously discussed method, `Frame.SystemEnable` and `Frame.SystemDisable` will work within a simulation deterministically. [As specified in the documentation](/quantum/current/manual/quantum-ecs/systems#activating-and-deactivating-systems), When a system is disabled or enabled, the `Update` methods of these systems will not be executed nor will `Signals`. One way to achieve something like this could be to create a `System` that handles disabling and re-enabling all systems by responding to a pause action from a player's input. The following is one example of how a system like this could work:

C#

```csharp
public unsafe class PauseSystem : SystemMainThreadFilter<PauseSystem.Filter>
{
    public override void Update(Frame frame, ref Filter filter)
    {
        Input* input = frame.GetPlayerInput(filter.Avatar->PlayerRef);
        if (input->Paused.WasPressed)
        {
            frame.Global->Paused = !frame.Global->Paused;
            foreach (var system in frame.SystemsAll)
            {
                if (system is PauseSystem)
                {
                    continue;
                }
                if (frame.Global->Paused)
                {
                    frame.Signals.OnPause();
                    frame.SystemDisable(system);
                }
                else
                {
                    frame.SystemEnable(system);
                    frame.Signals.OnUnpause();
                }
            }
        }
    }
    public struct Filter
    {
        public EntityRef Entity;
        public Avatar* Avatar;
    }
}

```

This example is of a `PauseSystem` that toggles a global `Paused` bool. When the player presses the pause button, the system iterates through all of the other systems and disables or enables them based on the value of `Paused`. It also sends a pause or an unpause `Signal`, which would be defined in a [DSL](/quantum/current/manual/quantum-ecs/dsl) file. Note, the pause `Signal` is called before disabling the systems; otherwise, any system that inherits from it will not trigger it; likewise, the unpause `Signal` is called after the `Systems` are enabled. The reason all of the systems are iterated through is because some `Systems` that come with Quantum such as physics systems, do not inherit from a user-defined `Signal` by default. Additionally, the `PauseSystem` itself is not paused; if this was done, there would be no way for the game to be unpaused.

NOTE

While `Systems` are disabled, the simulation will continue and the frame number will continue to increase. This can cause issues with some frame-dependent elements such as `FrameTimers`, which will be discussed later in this documentation.

## Not Recommend: Adjusting Frame.DeltaTime

Similar to how some Unity tutorials suggest pausing a game by setting `Time.timeScale` to 0, a similar concept _could_ be done in Quantum with `Frame.DeltaTime` and settings its value to `FP.\_0`. This, however, is **NOT** recommended. Because of the way time is used in some systems, setting `Frame.DeltaTime` cause various issues:

- Any values that are divided by `Frame.DeltaTime` will cause a dividing by zero error.
- Some systems that increments values that do not `Frame.DeltaTime` will continue to progress in `Systems` that are still active.

Again, in theory you could try to pause a game by doing this, but it is **NOT** recommended and could cause many bugs; instead, disabling and enabling `Systems` is the preferred method.

## Pausing FrameTimers

[`FrameTimers`](frame-timer) are a useful tool for creating timers that don't require constant management and updating, but a trade off of them is that they cannot be paused since they only store their target frame. When pausing using either of the two previously mentioned methods, the simulation will continue, so when unpausing, the `FrameTimers` will have progressed and may have even finished. There are various ways this could be resolved. The following is one such solution. First, within the DSL file, defining a new type of struct.

```
struct PausibleFrameTimer
{
    NullableFP RemainingSeconds;
    FrameTimer Timer;
}

```

This `struct` uses a `NullableFP` to store the remaining seconds of the `FrameTimer` that will be paused. If the timer was not running, then this `NullableFP` will not have a value, which will indicate it was not paused.

Then, additional methods could be added to this `PausibleFrameTimer`.

C#

```csharp
public partial struct PausibleFrameTimer
{
    public void Pause(Frame frame)
    {
        // If the timer is running, the remaining seconds are stored
        if (Timer.IsRunning(frame))
        {
            RemainingSeconds = Timer.RemainingSeconds(frame);
        }
        else
        {
            RemainingSeconds = default;
        }
    }
    public void Unpause(Frame frame)
    {
        // If the remaining seconds has a value, a new FrameTimer is defined with this value.
        if (RemainingSeconds.HasValue)
        {
            Timer = FrameTimer.FromSeconds(frame, RemainingSeconds.Value);
            RemainingSeconds = default;
        }
    }
}

```

When the `OnPause` and `OnUnpause``Signals` are called, a system could iterate through `Components` that utilize `PausibleFrameTimers` and execute `PausibleFrameTimer.Pause` and `PausibleFrameTimer.Unpause` accordingly. This is just one possible implementation. Again, the need to pause an online multiplayer game is rare, but there are multiple solutions and implementations that could be done to do so, but they need to be approached with caution.

Back to top

- [Disabling The Quantum Runner](#disabling-the-quantum-runner)
- [Recommended: Disabling and Re-Enabling Systems](#recommended-disabling-and-re-enabling-systems)
- [Not Recommend: Adjusting Frame.DeltaTime](#not-recommend-adjusting-frame.deltatime)
- [Pausing FrameTimers](#pausing-frametimers)