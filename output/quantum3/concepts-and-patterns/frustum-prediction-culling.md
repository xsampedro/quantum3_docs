# frustum-prediction-culling

_Source: https://doc.photonengine.com/quantum/current/concepts-and-patterns/frustum-prediction-culling_

# Frustum Prediction Culling

The snippet shows how to apply simple frustum culling to Quantum entities.

Culled entities in Quantum will not be updated when simulating predicted frames which is a performance improvement.

Read more about [Profiling](/quantum/current/manual/prediction-culling "Prediction Culling") in Quantum.

Implementing a custom culling system involves three main steps:

## Custom culling callback definition

Declare the callback that will be lately used for defining the actual culling logic:

C#

```csharp
namespace Quantum
{
  using Photon.Deterministic;
  partial class FrameContextUser
  {
    public delegate bool CullingDelegate(FPVector3 position);
    public CullingDelegate CullingCallback;
  }
}

```

## Custom culling System

Create a new system which triggers the callback.

P.S: this system should _replace_ the regular culling systems which are provided by default on the SystemsConfig asset.

C#

```csharp
namespace Quantum
{
  using UnityEngine.Scripting;
  using Quantum.Task;
  [Preserve]
  public unsafe class CustomCullingSystem : SystemBase
  {
    private TaskDelegateHandle _updateTask;
    public override void OnInit(Frame f)
    {
      f.Context.TaskContext.RegisterDelegate(Update, "CustomCullingSystem Update", ref _updateTask);
    }
    private void Update(FrameThreadSafe frame, int start, int count, void* arg)
    {
      var f = (Frame)frame;
      var context = f.Context;
      var filter = f.Filter<Transform3D>();
      while (filter.NextUnsafe(out var entity, out var transform))
      {
        if (context.CullingCallback(transform->Position))
        {
          f.Cull(entity);
        }
      }
    }
    protected override TaskHandle Schedule(Frame f, TaskHandle taskHandle)
    {
      if (f.IsVerified)
      {
        return taskHandle;
      }
      if (f.Context.CullingCallback != null)
      {
        var handle = f.Context.TaskContext.AddMainThreadTask(_updateTask, null);
        handle.AddDependency(taskHandle);
        return handle;
      }
      return taskHandle;
    }
  }
}

```

## Frustum culling logic

Create a new script in which the culling logic is implemented and set in the callback previusly created:

C#

```cs
namespace Quantum
{
  using Photon.Deterministic;
  using UnityEngine;
  public class FrustumPredictionCulling : QuantumSceneViewComponent
  {
    public float ProximityRadiusOverhaul = 20;
    public override void OnActivate(Frame frame)
    {
      Game.Frames.Verified.Context.CullingCallback = Callback;
    }
    private bool Callback(FPVector3 position)
    {
      var unityPosition = position.ToUnityVector3();
      var distance = (unityPosition - Camera.main.transform.position).sqrMagnitude;
      if (distance < ProximityRadiusOverhaul * ProximityRadiusOverhaul) return false;
      var normalizedPos = Camera.main.WorldToViewportPoint(unityPosition);
      if (normalizedPos.z <= 0) return true;
      if (normalizedPos.y < 0 || normalizedPos.y > 1) return true;
      if (normalizedPos.x < 0 || normalizedPos.x > 1) return true;
      return false;
    }
  }
}

```

Back to top

- [Custom culling callback definition](#custom-culling-callback-definition)
- [Custom culling System](#custom-culling-system)
- [Frustum culling logic](#frustum-culling-logic)