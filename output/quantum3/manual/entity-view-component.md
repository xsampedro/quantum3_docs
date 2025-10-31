# entity-view-component

_Source: https://doc.photonengine.com/quantum/current/manual/entity-view-component_

# Entity View Component

## Introduction

The Entity View framework was designed to quickly bind the simulation internals of an entity to view-side code. It's a simple Entity Component system like MonoBehaviours for example but for Quantum Entity Views.>

By adding a script derived from the `QuantumEntityViewComponent` class to any game object under the EntityView the script immediately has access to the `EntityRef`, `Game`, `PredictedFrame` and `VerifiedFrame` and various virtual method to override.

All view component are bundled on the `QuantumEntityViewUpdater` (EVU) and are updated from there.

The view components work well with the `QuantumEntityViewPool`.

## Lifetime

The following callbacks define the life-cycle of (pooled) entity view and its `QuantumViewComponents`.

|     |     |
| --- | --- |
| **`OnInitialize()`** | Is called when the entity view component is created for the first time. Its view context is already set, but access to `Game`, `VerifiedFrame`, `PredictedFrame` and `PredictedPreviousFrame` is not available yet. |
| **`OnActivate(Frame frame)`** | Is called when the entity view has been created or enabled (also from pool) for example when a new Quantum entity was created and the EVU created the associated EntityView. |
| **`OnDeactivate()`** | Is called before the entity view is deactivated either destroyed or returned to the pool. |
| **`OnUpdateView()`** | Is called on Unity updated originating from the EVUs `OnObservedGameUpdated()`. |
| **`OnLateUpdateView()`** | Is called on the Unity LateUpdate() method inside EVU. |
| **`OnGameChanged()`** | Is called after the observed game in the EVU changed. |

## Contexts

A `QuantumEntityViewComponent` can be defined to have one optional generic type for a context object. A context object is a simple class (can be a MonoBehaviour, singleton, etc) that derives from `IQuantumViewContext`. It's an option to share data between different parts of the game.

The context has to be in the same game object as the `QuantumEntityViewUpdater` or in a child game object and it must be added during `Awake()` to be automatically loaded and made available for view components.

C#

```csharp
namespace Quantum {
  using UnityEngine;
  public class MyGameContext : QuantumMonoBehaviour, IQuantumViewContext {
    public GameObject Template;
  }
}

```

The context is then accessible by the `ViewContext` property.

C#

```csharp
namespace Quantum {
  using UnityEngine;
    public class MyViewScript : QuantumEntityViewComponent<MyGameContext> {
      GameObject _go;
      public override void OnInitialize() {
        _go = Instantiate(ViewContext.Template);
      }
  }
}

```

A dictionary of each loaded context type can be accessed on the EVU: `Dictionary<Type, IQuantumViewContext> Context`

## Scene View Components

A `QuantumSceneViewComponent` is a view component that does not have an associated entity. It can be be added to any object in the scene to get access to the view component properties. But it has to be added to the `QuantumEntityViewUpdater` explicitly.

Either use the `Updater` field to reference the EVU directly.

Or toggle on `UseFindUpdater` which causes a `FindFirstObjectByType()` during OnEnable() which may be considered to be slow.

The ViewUpdater also allows to dynamically add and remove scene view components. Although `OnInitialize()` and `OnDeactivate()` are called right away, `OnActivate()` is deferred to the next Update call.

`QuantumEntityViewUpdater.AddViewComponent(IQuantumViewComponent viewComponent)`

`QuantumEntityViewUpdater.RemoveViewComponent(IQuantumViewComponent viewComponent)`

## Example

Setting character animation based on character controller state.

C#

```csharp
namespace Quantum {
  using UnityEngine;
  public class CharacterViewAnimations : QuantumEntityViewComponent {
    private Animator _animator;
    public override void OnInitialize() {
      _animator = GetComponentInChildren<Animator>();
    }
    public override void OnUpdateView() {
      // probably should use RealSpeed, but the variable isn't been written to in the KCC code currently
      var kcc = PredictedFrame.Get<KCC>(EntityRef);
      var kinematicSpeed = kcc.Data.KinematicVelocity.Magnitude;
      _animator.SetFloat("Speed", kinematicSpeed.AsFloat * 10);
      _animator.SetBool("Jump", kcc.Data.HasJumped);
      _animator.SetBool("FreeFall", !kcc.Data.IsGrounded);
      _animator.SetBool("Grounded", kcc.Data.IsGrounded);
    }
  }
}

```

Camera follow behaviour example:

C#

```csharp
// Context, added to QuantumEntityViewUpdater game object
namespace Quantum {
  using UnityEngine;
  public class CustomViewContext : MonoBehaviour, IQuantumViewContext {
    public Camera MyCamera;
  }
}
// View component, added to entity prefab (QuantumEntityView)
namespace Quantum {
  using UnityEngine;
  public class QuantumCameraFollow : QuantumEntityViewComponent<CustomViewContext> {
    public Vector3 Offset;
    public float LerpSpeed = 4;
    private bool _isPlayerLocal;
    public override void OnActivate(Frame frame) {
      var playerLink = frame.Get<PlayerLink>(EntityRef);
      _isPlayerLocal = Game.PlayerIsLocal(playerLink.Player);
    }
    public override void OnUpdateView() {
      if (_isPlayerLocal == false) {
        return;
      }
      var myPosition = transform.position;
      var desiredPos = myPosition + Offset;
      var currentCameraPos = ViewContext.MyCamera.transform.position;
      ViewContext.MyCamera.transform.position = Vector3.Lerp(currentCameraPos, desiredPos, Time.deltaTime * LerpSpeed);
      ViewContext.MyCamera.transform.LookAt(transform);
    }
  }
}

```

Back to top

- [Introduction](#introduction)
- [Lifetime](#lifetime)
- [Contexts](#contexts)
- [Scene View Components](#scene-view-components)
- [Example](#example)