# optimizations

_Source: https://doc.photonengine.com/quantum/current/concepts-and-patterns/optimizations_

# Optimizations

## CPU Performance Optimizations

### Prediction Culling

Quantum, by default runs a full prediction and rollback loop for the entire game state. This means multiple simulation updates are simulated for each entity which can be very CPU expensive.

Prediction Culling can be applied to run the predicated simulation only for entities in a specified culling area which can reduce CPU usage drastically.

For more information see [Prediction Culling](/quantum/current/manual/prediction-culling)

### Object Pooling

By default, the `QuantumEntityViewUpdater` creates new instances of the `QuantumEntityView` prefabs whenever an entity gets created and destroys the view GameObjects respectively. This can lead to CPU performance spikes and increased garbage collection. For best performance `QuantumEntityViews` should be pooled.

For more information see [Pooling](/quantum/current/manual/entityview#pooling)

### Simulation Rate

The Simulation Rate in Quantum determines how many simulation ticks per second Quantum executes and can be configured in the `SessionConfig`.

This is the equivalent to `FixedUpdate` in Unity. Running at a lower tick rate frees up CPU budget. It is recommended to lower the Simulation Rate as far as possible as long as it does not affect the quality of the gameplay experience.

### TargetFrameRate / vSync

By default, on some platforms, Unity tries to render as many frames per second as possible. This additional CPU load can negatively impact the gameplay quality and make lag spikes more noticeable. It is recommended to enable `vSync` or set an `Application.TargetFrameRate` or to give players configuration options for these values.

## Bandwidth Optimization

Optimizing bandwidth is important for a variety of reasons:

1. It reduces the costs of your application by reducing bandwidth costs.
2. It allows a larger percentage of your player base to play the game without encountering networking issues.
3. It allows for larger player counts per session.

### Optimizing Inputs

The majority of the data that Quantum sends over the network are inputs. It is crucial that the input type you use is optimized.

Learn more about input optimization [here](/quantum/current/manual/input#optimization)

Back to top

- [CPU Performance Optimizations](#cpu-performance-optimizations)

  - [Prediction Culling](#prediction-culling)
  - [Object Pooling](#object-pooling)
  - [Simulation Rate](#simulation-rate)
  - [TargetFrameRate / vSync](#targetframerate-vsync)

- [Bandwidth Optimization](#bandwidth-optimization)
  - [Optimizing Inputs](#optimizing-inputs)