# overview

_Source: https://doc.photonengine.com/quantum/current/addons/animator/overview_

# Overview

Level

INTERMEDIATE

## Introduction

Quantum’s deterministic animator works by baking information from Unity’s Mecanim Controller. It imports many configurations such as the States, the Transitions, the Motion clips and so on.

**The main advantage of using it** is that it can control animations to be 100% in sync across machines. Animations will snap to the correct state and clip time accordingly to the predict-rollback. So it provides very accurate tick based animations, which is something that is usually needed on Fighting and some Sports games because the animations should be perfectly in sync between all client simulations.

## Known issues/Limitations

The Animator, as it comes by default, has these known limitations:

1. No support for hierarchical states;
2. When using 2D Blend Trees, it is mandatory to have the blending value between \[-1, 1\];
3. No support for Root Motion on 3D Transforms/PhysicsBody.

## Download - Addon

| Version | Release Date | Download |
| --- | --- | --- |
| 3.0.8 | Sep 29, 2025 | [Quantum Animator 3.0.8](https://downloads.photonengine.com/download/quantum/quantum-animator-3.0.8-alpha.unitypackage?pre=sp) |
| 3.0.7 | Aug 11, 2025 | [Quantum Animator 3.0.7](https://downloads.photonengine.com/download/quantum/quantum-animator-3.0.7-alpha.unitypackage?pre=sp) |
| 3.0.6 | Jul 04, 2025 | [Quantum Animator 3.0.6](https://downloads.photonengine.com/download/quantum/quantum-animator-3.0.6-alpha.unitypackage?pre=sp) |
| 3.0.5 | Apr 22, 2025 | [Quantum Animator 3.0.5](https://downloads.photonengine.com/download/quantum/quantum-animator-3.0.5-alpha.unitypackage?pre=sp) |

Back to top

- [Introduction](#introduction)
- [Known issues/Limitations](#known-issueslimitations)
- [Download - Addon](#download-addon)