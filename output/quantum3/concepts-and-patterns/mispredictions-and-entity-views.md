# mispredictions-and-entity-views

_Source: https://doc.photonengine.com/quantum/current/concepts-and-patterns/mispredictions-and-entity-views_

# Mispredictions and Entity Views

## Introduction

Photon Quantum is a predict-rollback engine, which means that players' last inputs are usually repeated during Predicted frames, so that the local simulation advances based on such predictions.

When a Verified frame arrives with the verified input set, the data that was previously advanced during prediction may need to be changed due to possible mispredictions.

This kind of rollback can be noticeable in the game view, especially on the positions of the entities.

A common example of this is when changing entities transform without any kind of acceleration/braking, i.e. their positions change a abruptly. See the gif below:

The entity at the top, that moves first, is the local player's entity.

The entity at the bottom is controlled by a remote player and the rollback is noticeable.

![Mispredictions and Entity Views](/docs/img/quantum/v3/concepts-and-patterns/entity-view-rollback.gif)

_PS:_ this is an exaggerated example in which the entities move very fast.

What happens is:

1. The remote player begins to press the ```
Right Arrow
```

    button, causing the entity to move to the right;
2. During predicted frames, such player input is repeated, causing the remote player entity to continue moving to the right;
3. The remote player stops pressing the 'Right Arrow' button, but the local client is not yet aware of this change, so it continues to predict that the entity is moving to the right;
4. When the next Verified frame arrives, the ```
EntityView
```

    lerps from the mispredicted position to the correct one, which is a slightly backwards;

The lerping is done by the Entity View's Error Prediction setup.

NOTE: the example above uses changes directly into the entity's Transform, but a similar behaviour would be observed in any other character movement solution, such as moving with a KCC, a PhysicsBody, etc

## Mitigation

This is expected to happen in the predict-rollback model, but there are ways to reduce the impact.

These are the two most common pieces of advice (number 1 is usually enough):

1. Add acceleration and deceleration to the entity's movement instead of always moving it abruptly. Just a little bit of (de)acceleration will help greatly as the differences due to misprediction will be much smaller;

![Mispredictions and Entity Views](/docs/img/quantum/v3/concepts-and-patterns/misprediction-with-acceleration.gif)

2. Add a little bit of Input Delay as it reduces the amount of predicted frames and thus reduces the amount of mispredictions;

![Mispredictions and Entity Views](/docs/img/quantum/v3/concepts-and-patterns/misprediction-with-input-delay.gif)

_PS:_ **avoid** changing the Input Delay unless you know that it can really be used in your game without major drawbacks;

If you decide to try this alternative, the suggestion is to navigate to the ```
DeterministicConfig
```

 asset and set ```
OffsetMin
```

to a small value, such as ```
2
```

, and see how the game plays.

Back to top

- [Introduction](#introduction)
- [Mitigation](#mitigation)