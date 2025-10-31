# overview

_Source: https://doc.photonengine.com/quantum/current/addons/flow-fields/overview_

# Overview

Level

ADVANCED

## Introduction

Flow Fields is a navigation technique suitable for large amounts of units with a shared destination on grid based games (for example RTS, Tower Defense, etc). This addon supports map sizes up to 256x256 tiles.

## Technical Requirement

- Quantum SDK 3.0 or up
- Unity 2021.3 LTS or up

## Download

| Version | Release Date | Download |
| --- | --- | --- |
| 3.0.0 | Aug 12, 2025 | [Quantum Flow Fields Addon 3.0.0](https://downloads.photonengine.com/download/quantum/quantum-flow-fields-addon-3.0.0.unitypackage?pre=sp) | [Release Notes](/quantum/current/addons/flow-fields/overview#release-history) |

## Known Issues/Limitations

- The Flow Fields Map data can only be modified in Verified frames;
- Lack of built-in unit avoidance (though it is possible to perform avoidance with using the physics systems - see the Advanced Example for more info).

## The three core concepts

- **Flow Field Map**: the grid in which units can move. This map is subdivided into smaller parts - the controllers;
- **Flow Field Controller**: these are smaller parts of the map. It contains portals (connectors between neighboring controllers) and precomputed flows to these portals. Controllers are subdivided into smaller parts - the tiles;
- **Flow Field Pathfinder**: the unit component which provides the direction in which it should move.

## Screenshots

![Overview](/docs/img/quantum/v2/addons/flow-fields/overview-1.gif)![Overview](/docs/img/quantum/v2/addons/flow-fields/overview-2.gif)![Overview](/docs/img/quantum/v2/addons/flow-fields/overview-3.gif)![Overview](/docs/img/quantum/v2/addons/flow-fields/overview-4.gif)
2000 units, each calculating new path every second.
## Release History

### 3.1.0

- Added support for multiple flow field maps which can be baked individually and referenced by flow fields agents;

- Added possibility to customize the baking of flow field maps which can be used, for example, to create different maps for ground and flying units;

- Added possibility to add an Offset to a Flow Field map so it can be baked out of the world central position;

- Fixed memory leak when settings paths to agents;

- Fixed GC Alloc;

- Breaking Change: FrameContextUser now has an array of flow field maps rather than only one. The map now has to be manually set again on the Quantum Runner;


### 3.0.0

- Initial Release;

Back to top

- [Introduction](#introduction)
- [Technical Requirement](#technical-requirement)
- [Download](#download)
- [Known Issues/Limitations](#known-issueslimitations)
- [The three core concepts](#the-three-core-concepts)
- [Screenshots](#screenshots)
- [Release History](#release-history)
  - [3.1.0](#section)
  - [3.0.0](#section-1)