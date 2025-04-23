# installation

_Source: https://doc.photonengine.com/quantum/current/addons/animator/installation_

# Installation

1. Import the .unitypackage in the Unity project;

2. Go to root `EntityPrototype` Game Object and add the script `AnimatorMecanim`:


![Adding Animator to Prototype](/docs/img/quantum/v3/addons/animator/animator-mecanim.png)

3. Add an empty Unity Animator component into the same object or in a child


object:

![Adding Unity Animator ](/docs/img/quantum/v3/addons/animator/animator-unity-animator.png)

4. In a folder of your preference, create a new `AnimatorGraph` asset:

![AnimatorGraph creation ](/docs/img/quantum/v3/addons/animator/animator-graph-creation.png)

5. Select it and, in the `Controller` field, make a reference to the Unity


animator controller that should be baked:

![Unity controller reference](/docs/img/quantum/v3/addons/animator/animator-controller-reference.png)

6. Click on the `Bake Animator Graph` button and the states, transitions, parameters, etc, will be baked into the asset:

![Baking Animator Graph](/docs/img/quantum/v3/addons/animator/animator-importing-controller.png)

7. In the Entity Prototype, add the `AnimatorComponent` and, in the `Animator Graph` field, reference the asset of preference:

![Adding AnimatorComponent](/docs/img/quantum/v3/addons/animator/animator-adding-component.png)

8. In the `SystemConfig` asset, add the AnimatorSystem, AnimatorBehaviourSystem and AnimatorTriggersSystem systems:

![Adding Systems](/docs/img/quantum/v3/addons/animator/animator-adding-systems.png)

9. Create a new Game Object in the game scene and add the `AnimatorViewUpdater` component to it:

![Adding AnimatorViewUpdater](/docs/img/quantum/v3/addons/animator/animator-adding-view-updater.png)

10. This is the initial setup. Using the Animator API on the simulation, such as setting animator parameter values in runtime is already enough to start the deterministic animations.

The basic API is similarly to how it is done on Unity, use Getters and Setters in order to read/write to the Animator:

C#

```csharp
// Getters
AnimatorComponent.GetBoolean(frame, filter.AnimatorComponent, "Defending");
AnimatorComponent.GetFixedPoint(frame, filter.AnimatorComponent, "Direction");
AnimatorComponent.GetInteger(frame, filter.AnimatorComponent, "State");
// Setters
AnimatorComponent.SetBoolean(frame, filter.AnimatorComponent, "Defending", true);
AnimatorComponent.SetInteger(frame, filter.AnimatorComponent, "Direction", 25);
AnimatorComponent.SetFixedPoint(frame, filter.AnimatorComponent, "Speed", FP._1);
AnimatorComponent.SetTrigger(frame, filter.AnimatorComponent, "Shoot");

```

## Replacing the old Custom Animator

1. Make sure you have a backup of the project.
2. Delete `QuantumUser/Simulation/QuantumCustomAnimator`.
3. Delete `QuantumUser/View/CustomAnimator`.
4. Delete `Scripts/QuantumCustomAnimator`.
5. Import `QuantumAnimator.unitypackage`.
6. On all scripts that uses `CustomAnimator` replace with `AnimatorComponent`.
7. Remove `Custom.Animator.AnimatorUpdater` from `Frame.User`.

### Replacing CustomAnimatorGraph

In case the `CustomAnimatorGraph` assets is not working follow this steps:

1. Change the Inspector to Debug mode.
2. Select the asset and remove the attached `Controller`.
3. Click on `Import Mecanim Controller`.
4. Reattach the `Controller` and click on `Import Mecanim Controller` again;

Back to top

- [Replacing the old Custom Animator](#replacing-the-old-custom-animator)
  - [Replacing CustomAnimatorGraph](#replacing-customanimatorgraph)