# rngsession

_Source: https://doc.photonengine.com/quantum/current/manual/rngsession_

# RNG Session

## Introduction

In a deterministic setting, the output will always be the same given the same input. However, there are some situations where one would want to introduce some randomness to simulate the effects of unpredictable events or to provide some degree of variability to a system. Random number generators, also known as ```
RNG
```

s, are an important component in achieving this inside of Quantum.

## RNGSession

Quantum provides the developer with ```
RNGSession
```

. This can be used in the simulation to generate pseudorandom numbers deterministically. The frame ```
globals
```

come preloaded with a global session. A session will ALWAYS produce the same sequence of numbers unless seeded differently.

Usage:

C#

```csharp

// inside of a system
public override void OnInit(Frame frame)
{
int randomNum = frame.RNG->Next(0, 100);
}

```

The session found in ```
globals
```

is seeded with the ```
RuntimeConfig.Seed
```

field. It is up to the user to set a specific Seed if desired. When running an offline session with ```
QuantumRunnerLocalDebug
```

, the seed never changes, thus the sequence of numbers generated will always be the same.

While running through the example Menu provided, the seed will always be re-randomized if the users does not change it themselves (e.g leave it equals zero).

## Changing Seed at Runtime

If required, the seed can be changed at runtime by overwriting the session with a new one. Example:

C#

```csharp

// inside of simulation
public void ResetSeed(Frame frame)
{
 int newSeed = 100;
 frame.Global->RngSession = new Photon.Deterministic.RNGSession(newSeed);
}

```

This can be useful in some cases, such as frequently resetting the session in order to make the generation even more unpredictable.

## Component Usage

In development, it may be favorable to have your RNG on a per-component basis, instead of just using the global session. This could be useful for having each entity behave slightly differently, for example; having plants in a farming game grow at different intervals instead of all at the same time.

C#

```csharp
// DSL component
component MyComponent
{
 RNGSession Session;
}

```

You can also set the seed in the same way as the global session:

C#

```csharp
public void InitComponentWithSeed(MyComponent\* component)
{
 int newSeed = 100;
 component->Session = new RNGSession(newSeed);
}

```

## Avoiding Prediction Issues

Having the RNG on a per component basis also guarantees culling will not affect the final positions of predicted entities unless the rollback actually required it.

for more information, see: [Avoiding RNG Issues](/quantum/current/manual/prediction-culling#avoiding_rng_issues)

## Cheating

Due to the nature of determinism, predicting randomness is quite trivial. For example, a hacker could read the simulation locally, then duplicate the ```
RNGSession
```

 outside of the simulation and know the given sequence of numbers before anyone else. A common way to combat this is to frequently reset the seed to a number that cannot be easily predicted. For example, you could hash the player's input and use the result as the seed. This makes it very hard for anyone to control OR predict what the seed will be.

## A note on determinism

**Important:** Only advance RNG sessions from within Quantum Simulation code and never from Unity/view code. This is important due to the fact that the RNG session internally holds its state value which determines the sequence of random numbers.

This means that using ```
frame.RNG->Next()
```

(or a session contained within a component) and similar API advances the internal state, thus wrongly using it on view code would cause the game to desync (checksum error).

Back to top

- [Introduction](#introduction)
- [RNGSession](#rngsession)
- [Changing Seed at Runtime](#changing-seed-at-runtime)
- [Component Usage](#component-usage)
- [Avoiding Prediction Issues](#avoiding-prediction-issues)
- [Cheating](#cheating)
- [A note on determinism](#a-note-on-determinism)