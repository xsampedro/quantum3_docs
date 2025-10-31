# change-worker-thread-count

_Source: https://doc.photonengine.com/quantum/current/concepts-and-patterns/change-worker-thread-count_

# Change Worker Thread Count

## Change Amount of Worker Threads per Platform

By default, the simulation will attempt to use a total number of threads equal to the `SimulationConfig.ThreadCount` setting, including the main-thread.

During its update loop, the simulation will request the provided `IDeterministicPlatformTaskRunner` to schedule a number of delegates equal to the number of extra worker threads, i.e.

`SimulationConfig.ThreadCount - 1`. The main-thread will _always_ execute the simulation task graph, so, even if no delegates are scheduled, the simulation will be able to proceed with its update.

It is possible to effectively use a varying number of worker threads (or none at all) based, for instance, on the platform in which the application is running.

This might be useful when targeting both mobile and PC, as PCs will usually have more fast-performing cores than mobile devices. In those cases, PC clients could use the full range of `SimulationConfig.ThreadCount`, while mobile devices could be further limited to a maximum of 2 threads, for instance, as a broad range of devices only has 2 fast cores.

In Unity, the default `IDeterministicPlatformTaskRunner` is implemented by `QuantumTaskRunnerJobs`. It uses Unity Jobs to spin up worker threads that will execute the simulation task graph in a work-stealing fashion alongside the main simulation thread.

By default, this Task Runner will clamp the number of effectively scheduled worker threads (in this case Unity Jobs) to the current `JobsUtility.JobWorkerCount - 1`.

This ensures that the simulation will not starve Unity from Job threads while its update loop is running, which could potentially cause deadlocks.

It is possible to further customize and override this value by extending `QuantumTaskRunnerJobs` as shown in the snippet below and adding the custom MonoBehaviour to a GameObject in the scene.

A common approach is to have such object added to a hub or menu scene and marked as `DontDestroyOnLoad`.

### Snippet

C#

```csharp
public class CustomQuantumTaskRunnerJobs : QuantumTaskRunnerJobs {
  private void Awake() {
    // force Custom mode to ensure that the Custom Count will be used
    quantumJobsMaxCountMode = QuantumJobsMaxCountMode.Custom;
  }
  protected override Int32 CustomQuantumJobsMaxCount {
    get {
#if !UNITY_EDITOR && UNITY_IOS
      // example: forcing a single-threaded simulation
      // (the main thread will always execute the task graph, this is the number of WORKER threads)
      return 0;
#elif !UNITY_EDITOR && UNITY_ANDROID
      // example: limiting worker threads to 1 in case there are more threads available
      return Math.Min(1, DefaultQuantumJobsMaxCount);
#else
      return DefaultQuantumJobsMaxCount;
#endif
    }
  }
}

```

Back to top

- [Change Amount of Worker Threads per Platform](#change-amount-of-worker-threads-per-platform)
  - [Snippet](#snippet)