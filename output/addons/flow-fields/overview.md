# overview

_Source: https://doc.photonengine.com/quantum/current/addons/flow-fields/overview_

# Overview

![Level 4](/v2/img/docs/levels/level04-advanced_1.5x.png)

## Introduction

Flow Fields is a navigation technique suitable for large amounts of units with a shared destination on grid based games (for example RTS, Tower Defense, etc). This addon supports map sizes up to 256x256 tiles.

## Technical Requirement

- Quantum SDK 3.0 or up
- Unity 2021.3 LTS or up

## Download

| Version | Release Date | Download |
| --- | --- | --- |
| 3.0.0 | Aug 29, 2024 | [Flow Fields 3.0.0 Build 499](https://dashboard.photonengine.com/download/quantum/quantum-flow-fields-addon-3.0.0.unitypackage) |  |

## Known Issues/Limitations

- The Flow Fields Map data can only be modified in Verified frames;
- Lack of built-in unit aviodance (though it is possible to perform avoidance with using the physics systems - see the Advanced Example for more info).

## The three core concepts

- **Flow Field Map**: the grid in which units can move. This map is subdivided into smaller parts - the controllers;
- **Flow Field Controller**: these are smaller parts of the map. It contains portals (connectors between neighbouring controllers) and precomputed flows to these portals. Controllers are subdivided into smaller parts - the tiles;
- **Flow Field Pathfinder**: the unit componenent which provides the direction in which it should move.

## Screenshots

![Overview](https://doc.photonengine.com/docs/img/quantum/v2/addons/flow-fields/overview-1.gif)![Overview](https://doc.photonengine.com/docs/img/quantum/v2/addons/flow-fields/overview-2.gif)![Overview](https://doc.photonengine.com/docs/img/quantum/v2/addons/flow-fields/overview-3.gif)![Overview](https://doc.photonengine.com/docs/img/quantum/v2/addons/flow-fields/overview-4.gif)
2000 units, each calculating new path every second.
## Release Notes

Based on Quantum version 3.0 Stable 1523

Back to top

- [Introduction](#introduction)
- [Technical Requirement](#technical-requirement)
- [Download](#download)
- [Known Issues/Limitations](#known-issueslimitations)
- [The three core concepts](#the-three-core-concepts)
- [Screenshots](#screenshots)
- [Release Notes](#release-notes)