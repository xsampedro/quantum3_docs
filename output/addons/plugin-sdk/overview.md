# overview

_Source: https://doc.photonengine.com/quantum/current/addons/plugin-sdk/overview_

# Overview

Available in the [Gaming Circle](https://www.photonengine.com/gaming) and [Industries Circle](https://www.photonengine.com/industries)

![Circle](/v2/img/docs/circles/icon-gaming_1x.png)

## Overview

The Custom Plugin SDK enables developers to create a custom Quantum server, test it locally and upload to an Photon Enterprise Cloud.

Creating a Photon Enterprise Cloud involves additional costs. Please contact us for more information and for the Plugin SDK download permissions.

Common features of a custom Quantum plugin include:

- Running the simulation on the Photon server cloud for trusted game results
- Sending server snapshots to late-joining or reconnecting clients

Starting a local Photon Server is only supported on Windows.

## Download

| Version | Release Date | Download |
| --- | --- | --- |
| 3.0.2 | Mar 04, 2025 | [Quantum Plugin SDK 3.0.2 Build 52](https://dashboard.photonengine.com/download/quantum/quantum-plugin-sdk-3.0.2.zip) | [Release Notes](#build-52-mar-04-2025) |
| 3.0.1 | Dec 13, 2024 | [Quantum Plugin SDK 3.0.1 Build 49](https://dashboard.photonengine.com/download/quantum/quantum-plugin-sdk-3.0.1.zip) | [Release Notes](#build-49-dec-13-2024) |
| 3.0.0 | Dec 09, 2024 | [Quantum Plugin SDK 3.0.0 Build 44](https://dashboard.photonengine.com/download/quantum/quantum-plugin-sdk-3.0.0-44.zip) | [Release Notes](#build-44-dec-09-2024) |

## Changelog

### 3.0.2

#### Build 52 (Mar 04, 2025)

- Input deserialization errors caused by clients are now ignored, set ```
ClientInputExceptionTolerance
```

(default is 2) to define how many are ignored before disconnecting the client
- Custom plugins running on the cloud are using a more efficient native memory allocator
- Replaced the scheduled plugin fiber timer with an incremental timer, started at the end of an update which contributes to server stability
- Fixed an issue that prevented ```
MaxPlayerSlots
```

to be set correctly when selecting ```
  -1
```

(unlimited) by webhooks
- The lobby restrictions are now less strict, max count from ```
3
```

to ```
20
```

and max string length from ```
64
```

to```
128
```

, violations now fail the game creation instead of silently removing the properties
- Added a mapping from ```
PlayerRef
```

to client object using ```
DeterministicServer.GetClientForPlayer(PlayerRef)
```


### 3.0.1

#### Build 49 (Dec 13, 2024)

- Fixed an issue that could cause the disconnect signal to be missing when replacing disconnected clients on the plugin
- Fixed an issue that would disconnect fast reconnecting clients before the session has started
- Removing ```
ActorNr
```

property from ```
Create
```

and ```
Join
```

webhooks because the ids are not know at this point

#### Build 39 (Dec 02, 2024)

- Added ```
ActorNr
```

and ```
AuthCookie
```

properties for all player related webhooks
- Added ```
UserData
```

to the ```
CreateGameResponse
```

webhook
- Added ```
IsInactive
```

to the ```
LeaveGameRequest
```

webhook
- Added a getter for ```
DeterministicPluginClient
```

objects on the ```
DeterministicServer
```

class
- Added the ```
RunServerSimulation
```

flag to the CreateGame webhook response to control running server simulation for particular games
- Added a new demo project that shows how to intercept and rewrite commands
- Changed the server simulation wrapper ```
DotNetSessionRunner
```

and moved responsibility to another interface implemented by ```
DotNetSessionContext
```

which is responsible to make the resource manager or command serialized accessible by custom plugins that do not run the simulation
- Renamed ```
DeterministicServer.ServerSimulation
```

to ```
SessionRunner
```

- Deprecated the ```
DeterministicServer.ClientCount
```

property
- Fixed an issue that caused the game to stall after reconnecting into a single client online game using a RoomTTL
- Fixed an issue with server snapshots that caused disconnects with ```
Error #13: Snapshot request failed to start
```


### 3.0.0

#### Build 44 (Dec 09, 2024)

- Fixed an issue that could cause the disconnect signal to be missing when replacing disconnected clients on the plugin

#### Build 38 (Oct 16, 2024)

- Fixed an issue that caused server commands to flip the player connected flag for disconnected players resulting in unwanted PlayerConnected signals

#### Build 37 (Oct 01, 2024)

- Fixed an issue that caused the GZipped replay chunks to be larger than necessary

#### Build 28 (Jul 16, 2024)

- Fixed an issue with missing dependencies

#### Build 26 (Jun 19, 2024)

- Remove starting Windows performance counter, will be replaced by something new

#### Build 24 (Jun 11, 2024)

- Added ```
DeterministicServer.OnDeterministicSessionCanStart
```

callback to stall the session start on the plugin

#### Build 21 (Jun 04, 2024)

- Initial release

Back to top

- [Overview](#overview)
- [Download](#download)
- [Changelog](#changelog)
  - [3.0.2](#section)
  - [3.0.1](#section-1)
  - [3.0.0](#section-2)