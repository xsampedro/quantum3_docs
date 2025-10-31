# whats-new

_Source: https://doc.photonengine.com/quantum/current/getting-started/preview-3-1/whats-new_

# What's New

Quantum SDK 3.1 complements the Quantum 3 generation with performance and scalability improvements.

The [Roadmap](#roadmap) section list more upcoming features and improvements currently in development.

## Table Components

Until now, Quantum's ECS exclusively used sparse sets to store component values. Every component type had an associated sparse set (a kind of dictionary) linking entities to their component values.

Sparse set components can be added and removed quickly, but iterating several at once slows down when they store different entities in different orders. Typically, a filter can only sequentially access the first component's values; the rest require random lookups.

Quantum 3.1 adds another kind of component storage called **tables**. A table stores component arrays for all entities that have a specific combination of table components, and each entity in the table occupies the same "row" (index) in every "column" (array). Changing which table components an entity has requires moving it to a different table.

Table components flip the benefits around. They're faster to iterate, but slower to add and remove.

Components will continue to use sparse set storage by default. To use table storage, simply add the `table` keyword to a component's declaration in the DSL.

```
table component Health {
    FP Value;
}

```

You can mix and match both table and sparse set components in the same entities, so you can make the best of both worlds. That's it, API semantics work the same way as before, so no need to modify any code. Notice that the built-in components (Transforms, Colliders, Bodies, etc) are all stored as tables now.

## New Physics Solver

Quantum Physics is becoming **stateful** and now comes with a Projected Gauss-Seidel solver with warm-start (you can still select the legacy solvers if you want to keep the old behaviour).

This new solver makes use of time coherence to improve on two main aspects relative to solvers used in previous versions:

- **Better stacking**: able to handle more complex systems of constraints like stacks of objects and compound joints (e.g. ragdolls) by converging towards a suitable solution across multiple frames.
- **Less overshooting**: iterations are no longer agnostic to each other and can correct intermediary overshooting, leading to more correct solutions.

Stacks of objects in the new PGS Solver:

Your browser does not support the video tag.

Previous version solver (HGS) in the same scenario:

Your browser does not support the video tag.
## Heap and Memory Storage Improvements

Quantum 3.1 brings several improvements to the frame heap:

- **Large blocks**: Blocks are no longer limited to the page size. It's even possible to allocate an entire chunk of backing memory (called a "segment") as a single block.
- **Lazy expansion**: When the frame heap exhausts its existing memory segments, it will allocate another.

This release also introduces a new frame heap variant that holds individual variable-size blocks instead of pages of same-size blocks, potentially lowering heap fragmentation. Select "block-based" heap management in the `SimulationConfig` to use this variant.

Sparse set component blocks were also re-implemented with a lazy allocated sparse-buffer chunks, which greatly reduces the time spent in frame `CopyFrom` calls, specially on cases with large number of component types and large number of entities.

## Simulator Callbacks

New callbacks now expose some of the simulator internals through the `SimulatorContext`, informing about and giving control over the simulation flow.

Such controls can be useful in projects that require tight control over the simulation budget and include:

- **Simulator Stats**: how many Verified/Predicted frames were or still will be simulated, target ticks, current stage, Stalling status, etc.
- **Skip Prediction-rollback**: selectively skip prediction-rollback in order to exchange prediction accuracy for performance. Notice rollbacks still eventually need to always be executed. This allows you to alternate between a) advancing multiple verified frames without rollback and full resimulations and b) rolling back and fixing mispredictions with full resim when there are not many verified frames to move forward.
- **Limit Simulations**: limit the amount of Verified and/or Predicted ticks in each Update. This is the feature to be used in combination with the above to achieve the smoothest performance possible in these advanced cases.

The new callbacks are raised before and after the verification and prediction stages of the simulation.

Additionally, existing callbacks like `SimulationFinished` now also carry the simulator context.

| Callback | Description |
| --- | --- |
| CallbackBeforeSimulationStage | Is called once before the start of each simulation stage (verification and prediction). |
| CallbackSimulationStageFinished | Is called once at the end of each simulation stage (verification and prediction). |
| CallbackSimulateFinished | Is called when a frame simulation (Verified or Predicted) has completed. |

## Roadmap

These are features either in current development, or planned to be still included during the Early Access period. We may decide to postpone some of them to a future major release, but in principle these will be finished for the stable version of the new 3.1 SDK

- **3D Character Joints**: (in tests) enable Ragdolls and other ball-socket joint applications
- **Stateful Physics Entries**: (in tests) unlocks larger number of runtime objects
- **Additive Quantum Maps**: (in development) load multiple maps in runtime
- **Compound Entity Prototypes**: (in development) easily spawn and link groups of entities
- **Deterministic Navmesh baker**: (in development) deterministic drop-in replacement for Unity navmesh baking
- **More Performance improvements for Frame copies**: (planned) avoiding the extra 2 frame copies used only for interpolation and error correction (predicted-previous and previous-update-predicted) by storing only the necessary transform data (for visible entity views) instead of full frame data.

Back to top

- [Table Components](#table-components)
- [New Physics Solver](#new-physics-solver)
- [Heap and Memory Storage Improvements](#heap-and-memory-storage-improvements)
- [Simulator Callbacks](#simulator-callbacks)
- [Roadmap](#roadmap)