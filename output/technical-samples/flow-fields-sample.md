# flow-fields-sample

_Source: https://doc.photonengine.com/quantum/current/technical-samples/flow-fields-sample_

# Flow Fields Sample

![Level 4](/v2/img/docs/levels/level04-advanced_1.5x.png)

## Download

| Version | Release Date | Download |
| --- | --- | --- |
| 3.0.0 | Sep 02, 2024 | [Quantum Flow Fields Sample 3.0.0 Build 500](https://dashboard.photonengine.com/download/quantum/quantum-flow-fields-sample-3.0.0.zip) |  |

## Basic Example

The simple example showcases only one moving unit in a static predefined map.

![Example Basic](https://doc.photonengine.com/docs/img/quantum/v2/addons/flow-fields/example-basic-1.png)## Where to look?

- The example is implemented in the 'ExampleBasic' scene;
- \\quantum\_unity\\Assets\\Photon\\FlowFields\\ExampleBasic;

## Controls

- Right Mouse Button - sets unit's new destination

## Map Definition

- `FlowFieldMap` is created when simulation starts from static data (ExampleBasicSystem.cs)

## Advanced Example

The advanced example showcases multiple units with avoidance and final destination grouping.

![Example Advanced](https://doc.photonengine.com/docs/img/quantum/v2/addons/flow-fields/example-advanced-1.png)## Where to look?

- The example is implemented 'ExampleAdvanced' scene;
- \\quantum\_unity\\Assets\\Photon\\FlowFields\\ExampleAdvanced;
- \\quantum\_code\\quantum.code\\ExampleAdvanced

## Controls

- Left Mouse Button - click/drag to select single/multiple units;
- Right Mouse Button - set a new destination for selected unit;
- WASD/arrows - camera movement;
- scroll wheel - camera zoom;
- Q - spawn a new unit on the cursor position;
- E - change the cost of the tile on cursor position (between 1 a 255)

## Map Definition

- The Map parameters are defined in `TileMapSetup` (\\Assets\\Photon\\FlowFields\\ExampleAdvanced\\Resources\\DB\\ExampleAdvancedTileMapSetup);
- The Map cost field is baked via `TileMapBaker` \- based on Static Box and Circle colliders placed in scene;
- `FlowFieldMap` is created when the simulation starts (TileMapSystem.cs)

## Movement

Units are moved by setting their velocity based on data provided by `FlowFieldPathfinder`.

## Avoidance

Avoidance between units is done with physics. The `Physics Solver Iterations` are set to 0. To tweak the behaviour you can play with `Penetration Allowance` and `Penetration Correction` in `Simulation Config`.

- Penetration Allowance - Allow a certain degree of penetration to improve the stability of the physics simulation.
- Penetration Correction - How much of the exceeding penetration (above allowance) should be corrected in a single frame. 0 = no correction, 1 = full correction.

## UnitGroup

When multiple units are controlled at the same time they are put into a UnitGroup. Units still move on their own (they can choose different paths towards the destination) but their final destination is stored in the formation.

## Release Notes

Based on Quantum version 3.0 Stable 1523

Back to top

- [Download](#download)
- [Basic Example](#basic-example)
- [Where to look?](#where-to-look)
- [Controls](#controls)
- [Map Definition](#map-definition)
- [Advanced Example](#advanced-example)
- [Where to look?](#where-to-look-1)
- [Controls](#controls-1)
- [Map Definition](#map-definition-1)
- [Movement](#movement)
- [Avoidance](#avoidance)
- [UnitGroup](#unitgroup)
- [Release Notes](#release-notes)