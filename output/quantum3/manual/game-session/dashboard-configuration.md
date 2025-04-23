# dashboard-configuration

_Source: https://doc.photonengine.com/quantum/current/manual/game-session/dashboard-configuration_

# Dashboard Configuration

## Introduction

Quantum provides extra configuration via the Photon Dashboard for specific game and server related variables.

## Lobby Property Restrictions

Maximum number of properties: `20` (since 3.0.2)

Property types are restricted to: `bool`, `byte`, `short`, `int`, `long`, `string`, `byte\[\]`

Maximum `string``byte\[\]` length: `128` (since 3.0.2)

## Dashboard Properties

### AllowedLobbyProperties

| Property Name | Type |
| --- | --- |
| AllowedLobbyProperties | String |

Set a list of properties that are allowed for the client to send as Lobby Properties as a protection for the matchmaking performances on the master servers. If this property is set, then non-listed properties send by clients will be stripped.

### BlockNonProtocolMessages

| Property Name | Type |
| --- | --- |
| BlockNonProtocolMessages | Boolean |

Default is `false`.

We strongly recommend settings this to `true`.

Cancels all non-protocol messages sent to the server, logging "Blocked non-protocol message" once when the condition is triggered. Active only when set to true.

### BlockPlayerProperties

| Property Name | Type |
| --- | --- |
| BlockPlayerProperties | Boolean |

Default is `false`.

We strongly recommend settings this to `true`.

Cancels all player property set requests from clients, logging "Blocked player properties" once when triggered. Active only when set to true.

### BlockRoomProperties

| Property Name | Type |
| --- | --- |
| BlockRoomProperties | Boolean |

Default is `false`.

We recommend settings this to `true`.

Blocks all room properties set by clients post-creation, except for the special `StartQuantum` room property.

Caveat: This affects `IsOpen` and `IsVisible` of the room as well. Use `HideRoomAfterStartSec` and `CloseRoomAfterStartSec` to control them instead.

### ClientInputExceptionTolerance

| Property Name | Type |
| --- | --- |
| ClientInputExceptionTolerance | Integer |

Default is `2`.

This property sets the number of input serialization errors that will be ignore before a client is disconnected. Quantum 3.0.1 and earlier used to disconnect clients after the first occurrence. The setting makes the server slightly more robust against random network corruptions.

Since version: 3.0.2.

### HideRoomAfterStartSec

| Property Name | Type |
| --- | --- |
| HideRoomAfterStartSec | Integer |

Default is `-1` (disabled).

If set to a number greater than zero, `n` seconds after Quantum start the room will be removed from public or search listings.

This can help manage room visibility and ensure that new players do not join games that are already in progress.

### CloseRoomAfterStartSec

| Property Name | Type |
| --- | --- |
| CloseRoomAfterStartSec | Integer |

Default is `-1` (disabled).

If set to a number greater than zero, `n` seconds after Quantum start the room will be closed.

Closing a room prevents any new players from joining and can be used to manage the lifecycle of the game session.

### MaxPlayerSlots

| Property Name | Type |
| --- | --- |
| MaxPlayerSlots | Integer |

Limits the number of player slots one client can activate. Default is `2`, which allows a client to create two local players.

Setting the property will restrict max player slots for all games running under this AppId.

The value can also be set by the webhooks (`CreateGame` and `JoinGame` responses) for individual clients.

### ServerUpdateRateMs

| Property Name | Type |
| --- | --- |
| ServerUpdateRateMs | Integer |

If set this property will set the server update rate on the server for all games. By default the update rate is derived from the `UpdateFPS` set in `SessionConfig` of a particular game.

Regardless, the update rate will be clamped between `16` and `60` ms.

Since version: 3.0.2.

### SessionConfig

| Property Name | Type |
| --- | --- |
| SessionConfig | Json |

Allows setting a global SessionConfig in the dashboard, overwriting client-sent configs to protect against malicious players. More options to control this config are provided by webhooks.

Caveat: For testing a local Photon Server use escape the JSON using `&quot;`

### StartPropertyBlockedTimeSec

| Property Name | Type |
| --- | --- |
| StartPropertyBlockedTimeSec | Integer |

Default is `-1` (disabled).

If set to a number greater than zero, the starting of Quantum is blocked until the `n` seconds has passed since the room has been created.

Can be used to ensure that players have enough time to join before the game begins.

### StartPropertyForcedTimeSec

| Property Name | Type |
| --- | --- |
| StartPropertyForcedTimeSec | Integer |

Default is `-1` (disabled).

If set to a number greater than zero, it specifies the maximum amount of seconds that can elapse before starting Quantum inside a room after the room has been created. If the specified time is exceeded, the game will set the "StartQuantum" property in the room's game properties to true if it hasn't already.

## Dashboard Properties (Enterprise)

### ServerSimulationEnabled

| Property Name | Type |
| --- | --- |
| ServerSimulationEnabled | Boolean |

Default is `true`.

If set to `false` server simulation is disabled. Will be overwritten when explicitly set by the `CreateQuantumGameResponse` webhook (`RunServerSimulation`).

### ServerSimulationPercent

| Property Name | Type |
| --- | --- |
| ServerSimulationPercent | Integer |

Default is `100`.

If the `ServerSimulationEnabled` is set to `true` this property can be used to only simulation a percentage of games.

## Security and Optimization

Quantum 3 introduces webhooks as a new method of protecting the plugin, moving away from the configuration-based security measures of Quantum 2.1. This shift aims to make it easier for developers to implement checks and detections into their own backend.

## Recommended Settings - Security Checklist

- Only disable `BlockRoomProperties` if the application requires to `RaiseEvent()`.
- Only disable `BlockPlayerProperties` is the application required player properties.
- Consider to enable `BlockRoomProperties` and never use room properties.
- Consider setting `HideRoomAfterStartSec` to allow setting `IsVisible` by the server.
- Consider setting `CloseRoomAfterStartSec` to allow setting `IsOpen` by the server.
- Set `MaxPlayerSlots` to `1` to never allow multiple players per client.
- If clients control when the Quantum session is started (e.g. lobby inside a Photon room) consider using `StartPropertyBlockedTimeSec` and `StartPropertyForcedTimeSec` to use server timeouts to finally start the game.

Back to top

- [Introduction](#introduction)
- [Lobby Property Restrictions](#lobby-property-restrictions)
- [Dashboard Properties](#dashboard-properties)

  - [AllowedLobbyProperties](#allowedlobbyproperties)
  - [BlockNonProtocolMessages](#blocknonprotocolmessages)
  - [BlockPlayerProperties](#blockplayerproperties)
  - [BlockRoomProperties](#blockroomproperties)
  - [ClientInputExceptionTolerance](#clientinputexceptiontolerance)
  - [HideRoomAfterStartSec](#hideroomafterstartsec)
  - [CloseRoomAfterStartSec](#closeroomafterstartsec)
  - [MaxPlayerSlots](#maxplayerslots)
  - [ServerUpdateRateMs](#serverupdateratems)
  - [SessionConfig](#sessionconfig)
  - [StartPropertyBlockedTimeSec](#startpropertyblockedtimesec)
  - [StartPropertyForcedTimeSec](#startpropertyforcedtimesec)

- [Dashboard Properties (Enterprise)](#dashboard-properties-enterprise)

  - [ServerSimulationEnabled](#serversimulationenabled)
  - [ServerSimulationPercent](#serversimulationpercent)

- [Security and Optimization](#security-and-optimization)
- [Recommended Settings - Security Checklist](#recommended-settings-security-checklist)