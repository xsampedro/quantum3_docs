# processors

_Source: https://doc.photonengine.com/quantum/current/addons/kcc/processors_

# Processors

## Interaction

KCC processors play a major role when interacting with KCC.

KCC supports two interaction types:

- **Collision**
  - Physics based interaction.
  - Regular colliders and triggers are supported.
  - Triggered when KCC starts/stops colliding with a scene object.
    - _Static collider_ \- The processor asset must be linked as User Asset on Quantum static collider component.
    - _Entity_ \- The ```
      KCC Processor Link
      ```

       component must be added to entity prototype and Processor asset assigned.
- **Modifier**
  - Manually registered ```
    KCCProcessor
    ```

     asset via ```
    KCC.AddModifier()
    ```

     and ```
    KCC.RemoveModifier()
    ```

    .

* * *

The ```
KCCProcessor
```

 declares two important methods:

C#

```csharp
public virtual bool OnEnter(KCCContext context, KCCProcessorInfo processorInfo, KCCOverlapHit overlapHit) => true;

```

- Invoked when the KCC starts colliding with a collider.
- Return value controls start of the interaction
  - The interaction can be deferred until certain conditions are met (it is not registered within KCC until true is returned)
  - For example when the player enters an area and gets a speed buff at full health only

C#

```csharp
public virtual bool OnExit(KCCContext context, KCCProcessorInfo processorInfo) => true;

```

- Invoked when the KCC stops colliding with a collider.
- Return value controls end of the interaction
  - Stopping interaction can be deferred until certain conditions are met (it keeps registered in KCC until true is returned)
  - For example when the player enters an area and gets slowdown debuff, removal is deferred until the player is at full health (no matter if the player still stands within the area or not)

## Stages

Stage is a sequence of specific method calls executed on processors. It is uniquely identified by an interface type.

KCC supports these stages executed during update:

- ```
  IBeforeMove
  ```

   \- executed at the beginning of the move. Used to configure KCC, enable/disable features, add forces, …
- ```
  IAfterMoveStep
  ```

   \- executed after each move step (physics query + depenetration from colliders) and before updating collision hits and firing ```
  OnEnter()
  ```

  /```
  OnExit()
  ```

   callbacks. Used for minor position corrections, vector projections, … This method can be called multiple times in a row if the KCC moves too fast (CCD is applied).
- ```
  IAfterMove
  ```

   \- executed on the end of the move. Use to apply any post processing.

C#

```csharp
public unsafe class StepUpProcessor : KCCProcessor, IAfterMoveStep
{
public void AfterMoveStep(KCCContext context, KCCProcessorInfo processorInfo, KCCOverlapInfo overlapInfo)
{
// 1\. Detect blocking geometry
// 2\. Push character upwards based on unapplied movement.
}
}

```

## Default Processors

Following features have separate implementation which makes them very easy to strip or replace by a different implementation.

#### Environment Processor

- Default processor packed with KCC addon.
- Defines behavior for grounded and floating state.
- Projected movement along ground tangent.
- Simple acceleration and friction model.
- Custom Gravity.
- Jump multiplier.
- Implements ```
  IBeforeMove
  ```

   stage - used to calculate desired ```
  KCCData.DynamicVelocity
  ```

   and ```
  KCCData.KinematicVelocity
  ```

  .
- Implements ```
  IAfterMoveStep
  ```

   stage - used to recalculate properties after each move step (for example projection of kinematic velocity on ground).
- The ```
  EnvironmentProcessor
  ```

   asset is located at ```
  Assets\\Photon\\QuantumAddons\\KCC\\AssetDB\\Processors
  ```

  .

![Environment Processor](/docs/img/quantum/v3/addons/kcc/environment-processor.jpg)
 Environment Processor.


Speed of the character is defined by this processor and can be used on scene objects to simulate behavior in various environments like water, mud, ice.

#### Step Up Processor

- Allows detection of steps (geometry which blocks horizontal movement, with walkable surface at some height).
- If a step is detected, the KCC moves upwards until it gets grounded.
- Maximum step height is defined by ```
  Step Height
  ```

  .
- Forward step check distance is controlled by ```
  Step Depth
  ```

  .
- The upward movement is equal to unapplied horizontal movement.
- The speed of upward movement is multiplied by ```
  Step Speed
  ```

   to compensate for loss of horizontal velocity.
- The ```
  StepUpProcessor
  ```

   prefab is located at ```
  Assets\\Photon\\QuantumAddons\\KCC\\AssetDB\\Processors
  ```

  .

![Step-Up processor](/docs/img/quantum/v3/addons/kcc/step-up-processor.jpg)
 Step-Up processor.


Following image shows process of detecting steps:

![Step up process](/docs/img/quantum/v3/addons/kcc/step-up-feature.jpg)
 Step up process.


1. Upward check when the character is blocked by an obstacle in horizontal direction.
2. Forward check to detect if the space in front of the character is collision free.
3. Ground snap check to detect if the ground is walkable.
4. The character moves upwards as long as all 3 checks pass.

#### Ground Snap Processor

- Allows keeping grounded state when the contact with ground is lost (stairs, uneven terrain).
- Pushes character closer to ground if gravity is not enough.
- Maximum snap distance is defined by ```
  Snap Distance
  ```

  .
- The movement speed towards ground is defined by ```
  Snap Speed
  ```

  .
- The ```
  GroundSnapProcessor
  ```

   prefab is located at ```
  Assets\\Photon\\QuantumAddons\\KCC\\AssetDB\\Processors
  ```

  .

![Ground snap processor](/docs/img/quantum/v3/addons/kcc/ground-snap-processor.jpg)
 Ground snap processor.


Following image shows process of ground snapping:

![Ground snap process](/docs/img/quantum/v3/addons/kcc/ground-snap-feature.jpg)
 Ground snap process.


The character is **virtually** pushed downwards until one of following conditions are met:

1. The character hits walkable ground => KCC keeps grounded and moves towards the ground.
2. The character hits non-walkable surface and cannot slide along it => KCC loses grounded state.
3. ```
   Snap Distance
   ```

    is reached => KCC loses grounded state.

Back to top

- [Interaction](#interaction)
- [Stages](#stages)
- [Default Processors](#default-processors)