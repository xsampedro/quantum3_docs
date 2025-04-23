# manual

_Source: https://doc.photonengine.com/quantum/current/addons/animator/manual_

This page is a work in progress and could be pending updates.

# Manual

## Animator View Updater

The `QuantumEntityViewUpdater` script is responsible for triggering the `Animate` callback on each entity with an `AnimatorComponent` attached. The entity must have a `QuantumEntityViewComponent` that implements the `IAnimatorEntityViewComponent` interface. The Animator Addon includes the `AnimatorMecanim` and `AnimatorPlayables`, which can be used as-is or extended for a more customized workflow.

### Animator Mecanim

Provides integration between Quantum Animator and Unity's Mecanim system, allowing for smooth transitions and frame synchronization of animations. It provides frame-rate control and enables the usage of [Animation Layers](https://docs.unity3d.com/Manual/AnimationLayers.html), which is not supported by `AnimatorPlayables`.

![Debug.](https://doc.photonengine.com/docs/img/quantum/v3/addons/animator/animator-mecanim-debug.png)

Muting Transitions

To prevent transitions from occuring automatically in Unity, make sure to check `Mute Graph Transitions On Export` in the `Animator Graph` settings.

### Animator Playables

It handles animations using Playable Graphs instead of Mecanim, providing a more customizable approach to animation updates with the Quantum Animator. However, not all Mecanim features are supported in this setup, such as Animation Layers.

## Animator Events

Quantum Animator supports both `Instant` events and `Time-Window` events. The event callback is triggered when the specified clip is played at the configured moment during event creation. This means two states that make use of the same clip will share events.

### Instant Event

C#

```csharp
public abstract class AnimatorInstantEventAsset : AnimatorEventAsset, IAnimatorEventAsset
  {
    /// <inheritdoc cref=&#34;AnimatorEventAsset.OnBake&#34;/>
    public new AnimatorEvent OnBake(AnimationClip unityAnimationClip, AnimationEvent unityAnimationEvent)
    {
      var quantumAnimatorEvent = new AnimatorInstantEvent();
      quantumAnimatorEvent.AssetRef = Guid;
      quantumAnimatorEvent.Time = FP.FromFloat_UNSAFE(unityAnimationEvent.time);
      return quantumAnimatorEvent;
    }
  }

```

Instant Events will trigger only one time the `Execute()` function once the Animator plays the clip that contains such event baked. A good pattern is to inherit from the base `AnimatorInstantEventAsset` and override the base class methods with different procedures, like custom `Frame.Siganls`.

Events assets can be created using `Create->Quantum->Assets...`

![Events.](https://doc.photonengine.com/docs/img/quantum/v3/addons/animator/animator-events-01.png)### Time-Window Event

The procedure is similar to instant event creation, but in this case, you must create a second event of the same type in the Unity Clip to mark the end of the event execution. For this event type, the Execute callback is triggered every frame between the start and end frames defined in the Unity Clip.

![Time-Window event.](https://doc.photonengine.com/docs/img/quantum/v3/addons/animator/animator-events-05.png)

In case you want to create new events that have a similar behavior with `OnEnter()`, `Execute()` and `OnExit()` methods, by inherit the `AnimatorTimeWindowEventAsset` class:

C#

```csharp
  /// <summary>
  /// This is a sample of how to use SampleTimeWindowEvent events. Use it as a base to create a new class inheriting from AnimatorInstantEventAsset and
  /// implement a custom logic on Execute method
  /// </summary>
  [Serializable]
  public class ExampleTimeWindowEventAsset : AnimatorTimeWindowEventAsset
  {
    public override unsafe void OnEnter(Frame frame, AnimatorComponent* animatorComponent, LayerData* layerData)
    {
      Debug.Log($&#34;[Quantum Animator ({frame.Number})] OnEnter animator time window event.&#34;);
    }

    public override unsafe void Execute(Frame frame, AnimatorComponent* animatorComponent, LayerData* layerData)
    {
      Debug.Log($&#34;[Quantum Animator ({frame.Number})] Execute animator time window event.&#34;);
    }

    public override unsafe void OnExit(Frame frame, AnimatorComponent* animatorComponent, LayerData* layerData)
    {
      Debug.Log($&#34;[Quantum Animator ({frame.Number})] OnExit animator time window event.&#34;);
    }
  }

```

### Event Setup

Follow these steps to define the event on an animation clip:

1. Add the `AnimationEventData` component to the same GameObject the contains `Animator` component, in your `QuantumEntityPrototype`:

![EventAsset creation.](https://doc.photonengine.com/docs/img/quantum/v3/addons/animator/animator-events-02.png)

2. On the Animation window select the clip and add a new `Unity AnimationEvent`:

![EventAsset creation.](https://doc.photonengine.com/docs/img/quantum/v3/addons/animator/animator-events-03.png)

3. Select the created event and specify the `Function` and `Object` as the image below:

![EventAsset creation.](https://doc.photonengine.com/docs/img/quantum/v3/addons/animator/animator-events-04.png)## Animator State Signals

Animator State Signals are now called when changing state both through transitions as well as [`FadeTo`](#fadeto). These are implemented by the `AnimatorBehaviourSystem` and can be implemented by other systems as well.

C#

```csharp
// Called in the first update when the Animator begins entering a state.
void OnAnimatorStateEnter(Frame frame, EntityRef entity, AnimatorComponent* animator, AnimatorGraph graph, LayerData* layerData, AnimatorState state, FP time);

// Called every frame while the Animator is updating a state.
void OnAnimatorStateUpdate(Frame frame, EntityRef entity, AnimatorComponent* animator, AnimatorGraph graph, LayerData* layerData, AnimatorState state, FP time, AnimatorStateType stateType);

// Called when the Animator updates for the last time, before fading to another state.
void OnAnimatorStateExit(Frame frame, EntityRef entity, AnimatorComponent* animator, AnimatorGraph graph, LayerData* layerData, AnimatorState state, FP time);

```

These signals provide the following:

- `Frame`: the current `Frame` that the signal was called
- `EntityRef`: the `EntityRef` associated with the `AnimatorComponent` provided to the signal
- A pointer to the `AnimatorComponent`
- References to the current `AnimatorGraph` and `AnimatorState`
- A pointer to the `LayerData`
- The current, corresponding time of the `LayerData` when the signal was called

`OnAnimatorStateUpdate` also provides the `AnimatorStateType`, an `enum` with the following:

- `None`
- `FromState`
- `CurrentState`
- `ToState`

During transitions, both the `ToState` and `FromState` send `OnAnimatorStateUpdate` signals, so `AnimatorStateType` is provided to help distinguish which one is being updated. For example, this could be used to prevent `AnimatorStateBehaviours` from executing their updates while being transitioned out of.

## Animator State Behaviours

In contrast to [Animator Events](#animator-events), `Animator State Behaviour` can be independent between states using the same clip and have support for `OnStateEnter`, `OnStateUpdate` and `OnStateExit`.

### Setup

1- In the Unity Animator window, select the state where you want to add the behavior, then attach the `AnimatorStateBehaviourHolder` script to it.

![State Behaviours.](https://doc.photonengine.com/docs/img/quantum/v3/addons/animator/animator-state-behaviour-holder.png)

2- Add a new element to the `AnimatorStateBehaviourAssets` list by using the `DebugBehavior` example included with the Animator package.

![State Behaviours.](https://doc.photonengine.com/docs/img/quantum/v3/addons/animator/animator-state-behaviour-asset.png)

3- Bake the graph asset using the menu option: `Tool->Quantum Animator->Bake All Graph Assets`.

4- Add the `AnimatorBehaviourSystem` to `Default Config Systems` asset.

5- The debug message should be logged every time the state is played.

`AnimatorStateBehaviour` is abstract, so classes that inherit from it, must implement the following methods:

C#

```csharp
bool OnStateEnter(Frame frame, EntityRef entity, AnimatorComponent* animator, AnimatorGraph graph, LayerData* layerData, AnimatorState state, FP time);
bool OnStateUpdate(Frame frame, EntityRef entity, AnimatorComponent* animator, AnimatorGraph graph, LayerData* layerData, AnimatorState state, FP time, AnimatorStateType stateType);
bool OnStateExit(Frame frame, EntityRef entity, AnimatorComponent* animator, AnimatorGraph graph, LayerData* layerData, AnimatorState state, FP time);

```

These are nearly identical to the [Animator State Signals](#animator-state-signals) mentioned previously; however, they return a `bool`. If using the provided `AnimatorBehaviourSystem`, returning `true` will stop the system from checking the remaining `AnimatorStateBehaviours` that may belong to the `AnimatorState`. For example, if an `AnimatorStateBehaviour` transitions to a new state using `FadeTo`, returning `true` will prevent the remaining `AnimatorStateBehaviours` from executing their `OnAnimatorStateUpdate` calls.

## FadeTo

The FadeTo can be used to transition to another state independent of a condition. Enable `AllowFadeToTransitions` to enable this feature.

C#

```csharp
  // public void FadeTo(Frame frame, AnimatorComponent* animatorComponent, string stateName, bool setIgnoreTransitions, FP deltaTime)
  //usage example
  if (input->Run.WasPressed)
  {
    graph.FadeTo(frame, filter.AnimatorComponent, &#34;Running&#34;, true, 0);
  }

```

## Root Motion

The Quantum Animator uses the the Unity's `AnimationClip` information to bake the motion curves data into the `AnimatorGraph` asset. The baked information can later be used in the simulation to modify the entity position and rotation according to the currently playing clip.

![Baked Root Motion data.](https://doc.photonengine.com/docs/img/quantum/v3/addons/animator/animator-baked-root-motion-data.png)#### How to use

The `Root Motion` option needs to be enabled in order to allow the motion data to be processed. It is also important to setup correctally the root motion configuration in the `Animation` tab of your source model, more information on [Unity's Root Motion documentation](https://docs.unity3d.com/Manual/RootMotion.html).

![Baked Root Motion data.](https://doc.photonengine.com/docs/img/quantum/v3/addons/animator/animator-enable-root-motion.png)#### Disable State Root Motion

By default, the motion data is enabled for all states in the `AnimatorGraph` as long s the `Root Motion` option is enabled. Adding the `Animator Disable Root Motion Behaviour` to a Mecanim state ensures that the state is excluded from root motion calculations.

![Baked Root Motion data.](https://doc.photonengine.com/docs/img/quantum/v3/addons/animator/animator-removing-root-motion-from-state.png)#### Usage Example

The QuantumAnimation is only responsible for computing the current motion delta for each Quantum Frame, allowing the user to use the generated data to update the entity's transform as needed. The code below shows an example application of `OnAnimatorRootMotion3D`and `OnAnimatorRootMotion2D`.

C#

```csharp
  // Handles 3D root motion
  public void OnAnimatorRootMotion3D(Frame frame, EntityRef entity, AnimatorFrame deltaFrame, AnimatorFrame currentFrame)
    {
      //Return in case there is no motion delta
      if (deltaFrame.Position == FPVector3.Zero && deltaFrame.RotationY == FP._0) return;

      if (frame.Unsafe.TryGetPointer<Transform3D>(entity, out var transform))
      {
        // Create a quaternion representing the inverse of the current frame&#39;s Y-axis rotation
        var currentFrameRotation = FPQuaternion.CreateFromYawPitchRoll(currentFrame.RotationY, 0, 0);
        currentFrameRotation = FPQuaternion.Inverse(currentFrameRotation);

        // Rotate the delta position by the inverse current rotation to align movement
        var newPosition = currentFrameRotation * deltaFrame.Position;

        // Apply the transform&#39;s rotation to the new position to get the world displacement
        var displacement = transform->Rotation * newPosition;

        var kccSettings = frame.FindAsset<KCCSettings>(frame.Unsafe.GetPointer<KCC>(entity)->Settings);

        // Compute an adjusted target hit position for raycasting
        var targetHitPosition =(displacement.XOZ.Normalized * FP._0_33 * 2 ) + displacement;

        // Perform a raycast in the direction of the intended motion to detect potential collisions with statics
        var hits = frame.Physics3D.RaycastAll(transform->Position, targetHitPosition.XOZ, targetHitPosition.Magnitude, -1,
          QueryOptions.HitStatics);

        if (hits.Count <= 0)
        {
          // If no collision, disable the character controller temporarily
          if (frame.Unsafe.TryGetPointer<KCC>(entity, out var kcc))
          {
            kcc->SetActive(false);
          }

          // Apply the motion and rotation to the transform
          transform->Position += displacement;
          transform->Rotate(FPVector3.Up, deltaFrame.RotationY * FP.Rad2Deg);
        }
        else
        {
          // If there is collision, enable the character controller
          if (frame.Unsafe.TryGetPointer<KCC>(entity, out var kcc))
          {
            kcc->SetActive(true);
          }
        }
      }
    }

    public void OnAnimatorRootMotion2D(Frame frame, EntityRef entity, AnimatorFrame deltaFrame, AnimatorFrame currentFrame)
    {
      //Return in case there is no motion delta
      if (deltaFrame.Position == FPVector3.Zero && deltaFrame.RotationY == FP._0) return;

      if (frame.Unsafe.TryGetPointer<Transform2D>(entity, out var transform))
      {
        // Calculate new rotation by applying delta
        FP newRotation = transform->Rotation + deltaFrame.RotationY;

        // Normalize rotation to keep it within [-π, π]
        while (newRotation < -FP.Pi) newRotation += FP.PiTimes2;
        while (newRotation > FP.Pi) newRotation += -FP.PiTimes2;

        // Rotate delta movement vector based on new orientation
        var deltaMovement = FPVector2.Rotate(deltaFrame.Position.XZ, newRotation);

         // Apply movement and rotation to the transform
        transform->Position += deltaMovement;
        transform->Rotation = newRotation;
      }
    }

```

Back to top

- [Animator View Updater](#animator-view-updater)

  - [Animator Mecanim](#animator-mecanim)
  - [Animator Playables](#animator-playables)

- [Animator Events](#animator-events)

  - [Instant Event](#instant-event)
  - [Time-Window Event](#time-window-event)
  - [Event Setup](#event-setup)

- [Animator State Signals](#animator-state-signals)
- [Animator State Behaviours](#animator-state-behaviours)

  - [Setup](#setup)

- [FadeTo](#fadeto)
- [Root Motion](#root-motion)