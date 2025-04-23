# webhooks

_Source: https://doc.photonengine.com/quantum/current/manual/webhooks_

# Webhooks

## Overview

Webhooks are primarily used to have a trusted source of game configurations and room compositions and can greatly contribute to the safety of an online application.

The default cloud plugin has support for different hooks defined via the AppID Dashboard.

Once configured and enabled, the Photon Cloud will send WebRequests (```
HTTP POST
```

) to a **custom backend** and will use the response data (Json) for various configurations on the Photon rooms and game sessions.

## Setup

Webhooks are enabled per AppId on the Photon Dashboard.

1. Navigate to the [Photon Dashboard](https://dashboard.photonengine.com/publiccloud) and log in.
2. Find the AppId and click ```
   Manage
   ```

   .
3. Scroll down to Plugins and click ```
   Edit
   ```

   .
4. Click ```
   Add New Pair
   ```

    and add ```
   keys
   ```

    and ```
   values
   ```

    (the maximum length allowed for each setting value is 1024 chars).
5. Press ```
   Save
   ```

    and wait up to one minute for changes to take effect.

![Photon Dashboard Properties](/docs/img/_shared/sdk/dashboard-properties.png)### Dashboard Configurations

| Key | Type | Example | Description |
| --- | --- | --- | --- |
| WebHookBaseUrl | ```<br>string<br>``` | ```<br>https://localhost:3581<br>``` | The base url of the custom backend. Webhook paths will be appended resulting in this format of paths: ```<br>{WebHookBaseUrl}/game/create<br>```<br>**Must always be set.** |
| WebHookIntegration | ```<br>string<br>``` | ```<br>Default<br>``` | The selected webhook integration (default is ```<br>Default<br>```<br>).<br>```<br>Default<br>```<br>, ```<br>PlayFab<br>``` |
| WebHookSecret | ```<br>string<br>``` | ```<br>\\*\\*\\*\*\*\*\*\*\*\*<br>``` | This will be send with each web request and can be used to authenticate the request. Will set the header named ```<br>X-SecretKey<br>```<br>. |
| WebHookCustomHttpHeaders | ```<br>Dictionary <string, string><br>``` | ```<br>{"A": "Foo", "B": "100" }<br>``` | A ```<br>JSON<br>```<br> dictionary. All entries will be added to the custom web request headers. Make sure the use double quotes. |
| WebHookEnableOnCreate | ```<br>bool<br>``` | ```<br>true<br>``` | If set to ```<br>true<br>```<br>, the [CreateGame](#creategame) webhook will be fired when a client is creating a game session. |
| WebHookEnableOnClose | ```<br>bool<br>``` | ```<br>false<br>``` | If set to ```<br>true<br>```<br>, the [CloseGame](#closegame) webhook will be fired when the game session was closed. |
| WebHookEnableOnJoin | ```<br>bool<br>``` | ```<br>true<br>``` | If set to ```<br>true<br>```<br>, the [JoinGame](#joingame) webhook will be fired when any client tries to join a game session. |
| WebHookEnableOnLeave | ```<br>bool<br>``` | ```<br>false<br>``` | If set to ```<br>true<br>```<br>, the [LeaveGame](#leavegame) webhook will be fired when any client leaves a Game Session |
| WebHookUrl{KEY} | ```<br>string<br>``` | ```<br>https://foo.net/path<br>``` | This pattern is optional and can be used to overwrite the default paths based on ```<br>WebHookBaseUrl<br>```<br> for individual webhooks. Valid replacements for ```<br>{KEY}<br>```<br> are ```<br>OnCreate<br>```<br>, ```<br>OnJoin<br>```<br>, ```<br>OnLeave<br>```<br>, etc. |

### Quantum Dashboard Configurations

| Key | Type | Example | Description |
| --- | --- | --- | --- |
| WebHookEnableGameConfigs | ```<br>bool<br>``` | ```<br>false<br>``` | If set to ```<br>true<br>```<br>, the [GameConfigs](#gameconfigs) webhook will be fired when a client uploaded the game configs ```<br>RuntimeConfig<br>```<br> and ```<br>SessionConfig<br>```<br> using the ```<br>StartRequest<br>```<br> operation. |
| WebHookEnableGameResult | ```<br>bool<br>``` | ```<br>false<br>``` | If set to ```<br>true<br>```<br>, the [GameResult](#gameresult) webhook will be invoked. |
| WebHookEnableAddPlayer | ```<br>bool<br>``` | ```<br>false<br>``` | If set to ```<br>true<br>```<br>, the client operation ```<br>AddPlayer<br>```<br> will invoke the [AddPlayer](#addplayer) web request. |
| WebHookEnablePlayerAdded | ```<br>bool<br>``` | ```<br>false<br>``` | If set to ```<br>true<br>```<br>, the [PlayerAdded](#playeradded) web request is invoked after a player has been successfully added to the Quantum online game. |
| WebHookEnablePlayerRemoved | ```<br>bool<br>``` | ```<br>false<br>``` | If set to ```<br>true<br>```<br>, the [PlayerRemoved](#playerremoved) web request is invoked after a player was removed from the Quantum game. |
| WebHookEnableReplay | ```<br>bool<br>``` | ```<br>false<br>``` | Setting this to ```<br>true<br>```<br> will enable the replay streaming based on the [ReplayStart](#replaystart) and [ReplayChunk](#replaychunk) web requests. |

## Webhook API

Webhooks send Json content and only accept Json content as a response. ```
UTF-8
```

charset for Json is mandatory.

Webhooks expect HTTP response codes ```
200
```

 or ```
400
```

:

- ```
200
```

: successful request.
- ```
400
```

: error or operation has been denied (e.g. a client cannot create a room/game session).

The Photon server retries webhooks after transportation errors three times as well as after receiving StatusCode ```
503
```

 (Service Unavailable) with delay of 400ms, 1600ms and 6400ms. The timeout before the request fails and is not retried is 10 seconds.

Fill out the [WebhookError](#webhookerror) definition to return information of the nature of an error back to the Photon plugin. It can be used for:

- Logging;
- Returning information to the client;
- A custom plugin to add further custom error handling;

### Http Request Retries

The Photon Server reacts to retry-able error responses from a backend by resending the request three additional times.

Subsequent request have different headers to identify them.

```
EGRepeatId
```

\- the repeated number for example ```
0
```

, ```
1
```

, ```
2
```

or ```
3
```

(int)

```
EGInvokeId
```

\- request id (int)

### Common Request Headers

These common request headers are added to **every** web request.

| Name | Type | Content | Description |
| --- | --- | --- | --- |
| ```<br>Accept<br>``` | ```<br>string<br>``` | ```<br>application/json<br>``` | Webhooks only accept ```<br>JSON<br>```<br> as a response body |
| ```<br>Accept-Charset<br>``` | ```<br>string<br>``` | ```<br>utf-8<br>``` | Webhooks only accept utf-8 as response body charset |
| ```<br>Content-Type<br>``` | ```<br>string<br>``` | ```<br>application/json<br>``` | Webhooks all send ```<br>JSON<br>```<br> body data |
| ```<br>X-SecretKey<br>``` | ```<br>string<br>``` | ```<br>\\*\\*\\*\*\*\*\*\*\*\*<br>``` | This key should only be know to the custom backend and should be used to authenticate the incoming web request. This is set on the Photon Dashboard as ```<br>WebHookSecret<br>```<br>. |
| ```<br>X-Origin<br>``` | ```<br>string<br>``` | ```<br>Photon<br>``` | Will always be set to "Photon" |

### CreateGame

This webhook is called before the room/game session is created on the Photon Server. The creation is blocked until the webhook receives a response, which will affect the time clients require to create a connection.

A ```
CreateGame
```

webhook is always a ```
JoinGame
```

request for the user that initiates the room/game session creation. There will be no subsequent ```
JoinGame
```

webhook for this user. This webhook shares the data from the ```
JoinGame
```

webhook.

Requires ```
WebHookBaseUrl
```

and ```
WebHookEnableCreateGame
```

to be set on the Photon dashboard.

JavaScript

```javascript
POST https://{WebHookBaseUrl}/game/create

```

**CreateGame Request**

| Name | Type | Example | Description |
| --- | --- | --- | --- |
| ```<br>AppId<br>``` | ```<br>string<br>``` | ```<br>d1f67eec-51fb-45c1<br>``` | The Photon AppId. |
| ```<br>AppVersion<br>``` | ```<br>string<br>``` | ```<br>1.0-live<br>``` | The AppVersion used when creating the room/game session. |
| ```<br>Region<br>``` | ```<br>string<br>``` | ```<br>eu<br>``` | The Region code of the Game Server that the room/game session was created in. |
| ```<br>Cloud<br>``` | ```<br>string<br>``` | ```<br>1<br>``` | The ```<br>Cloud Id<br>```<br> of that the Game Server is running on. |
| ```<br>UserId<br>``` | ```<br>string<br>``` | ```<br>db757806-8570-45aa<br>``` | The ```<br>UserId<br>```<br> of the client that creates the room/game session. |
| ```<br>AuthCookie<br>``` | ```<br>Dictionary<string, object><br>``` | ```<br>db757806-8570-45aa<br>``` | The Photon custom authentication cookie set by the backend. |
| ```<br>RoomName<br>``` | ```<br>string<br>``` | ```<br>e472a861-a1e2-49f7<br>``` | The room/game session Name. |
| ```<br>GameId<br>``` | ```<br>string<br>``` | ```<br>0:eu:e472a861-a1e2-49f7<br>``` | A unique ```<br>GameId<br>```<br> which is composed of ```<br>{Cloud:}{Region:}RoomName<br>```<br>. Can be overwritten in the response. |
| ```<br>EnterRoomParams<br>``` | ```<br>[EnterRoomParams](#enterroomparams)<br>``` | JSON: See ```<br>EnterRoomParam<br>```<br> section | The Photon room/game session Options sent by the client. |

Json Example:

JSON

```json
{
"AppId": "d1f67eec-51fb-45c1",
"AppVersion": "1.0-live",
"Region": "eu",
"Cloud": "1",
"UserId": "db757806-8570-45aa",
"AuthCookie": {
"Secret": "\*\*\*\*\*\*\*\*\*\*"
}
"RoomName": "e472a861-a1e2-49f7",
"GameId": "0:eu:e472a861-a1e2-49f7",
"EnterRoomParams": {
"RoomOptions": {
"IsVisible": true,
"IsOpen": true
}
}
}

```

**HTTP Response Codes**

| Name | Type | Description |
| --- | --- | --- |
| ```<br>200 OK<br>``` | ```<br>[CreateGame Response](#creategame_response)<br>``` | room/game session creation can commence, config data from the response will overwrite data sent by the client. |
| ```<br>400 Bad Request<br>``` | ```<br>[WebhookError](#webhookerror)<br>``` | room/game session creation is not allowed and will be canceled. The client will receive an error. |

**CreateGame Response**

| Name | Type | Description |
| --- | --- | --- |
| ```<br>GameId<br>``` | ```<br>string<br>``` | Overwrites the ```<br>GameId<br>```<br> used in subsequent web requests. Can be ```<br>null<br>```<br> or omitted. |
| ```<br>EnterRoomParams<br>``` | ```<br>[EnterRoomParams](#enterroomparams)<br>``` | Enforce selected Room Options during its creation. The ```<br>JSON<br>```<br> object does not have to include all members, only the ones that should be overwritten.<br>Only the initial options are protected by sending ```<br>EnterRoomParams<br>```<br>. Most of them can be changed by clients sending Photon Room properties. To block this enable the Photon dashboard property ```<br>BlockRoomProperties<br>```<br>. Can be ```<br>null<br>```<br> or omitted. |

| Name | Type | Description |
| --- | --- | --- |
| ```<br>SessionConfig<br>``` | ```<br>[SessionConfig](#sessionconfig)<br>``` | Return a ```<br>SessionConfig<br>```<br> object that is used by the game. Game configs sent by clients will be ignored. |
| ```<br>RuntimeConfig<br>``` | ```<br>[RuntimeConfig](#runtimeconfig)<br>``` | Return a ```<br>RuntimeConfig<br>```<br> object that is used by the game. Game configs sent by clients will be ignored. |
| ```<br>RuntimePlayer<br>``` | ```<br>[RuntimePlayer](#runtimeplayer)<br>``` | Return a ```<br>RuntimePlayer<br>```<br> object that will be used for the client that created this room/session.<br>This only overwrites the first ```<br>AddPlayer<br>```<br> data sent for player slot ```<br>0<br>```<br>. ```<br>MaxPlayerSlots<br>```<br> should be set to ```<br>1<br>```<br>. |
| ```<br>MaxPlayerSlots<br>``` | ```<br>int<br>``` | The maximum number of player slots this client can acquire:<br>```<br>0<br>```<br> = only spectating<br>```<br>1..255<br>```<br> = specific number<br>```<br>-1<br>```<br> = unlimited<br>If this response is sent but this value is not set ```<br>MaxPlayerSlots<br>```<br> will default to 1.<br>Players requesting an invalid player slot number or more slots than allowed will be disconnected. |
| ```<br>SnapshotsBlocked<br>``` | ```<br>bool<br>``` | This player will not be selected for sending game snapshots to other players if there are other clients available. |
| ```<br>StartPropertyBlockedTimeSec<br>``` | ```<br>int<br>``` | Minimum delay in seconds before starting Quantum inside a room after its creation, ensuring players have enough time to join. A value greater than zero activates this feature. |
| ```<br>StartPropertyForcedTimeSec<br>``` | ```<br>int<br>``` | Maximum delay in seconds before starting Quantum inside a room after its creation. Exceeding this time auto-activates "StartQuantum" if not already set. A value greater than zero activates this feature. |
| ```<br>HideRoomAfterStartSec<br>``` | ```<br>int<br>``` | Number of seconds after which the room will be hidden from public listings once Quantum starts within the room. A value greater than zero activates this feature. |
| ```<br>CloseRoomAfterStartSec<br>``` | ```<br>int<br>``` | Number of seconds after which the room will be closed following the start of Quantum inside the room, preventing new players from joining. A value greater than zero activates this feature. |

Json Example:

JSON

```json
{
"AppId": "d1f67eec-51fb-45c1",
"GameId": "0:eu:db757806-8570-45aa",
"EnterRoomParams": {
"RoomOptions": {
"IsVisible": true,
"IsOpen": true
}
},
"SessionConfig": {
"PlayerCount": 8,
"ChecksumCrossPlatformDeterminism": false,
"UpdateFPS": 30
},
"RuntimeConfig": {
"Map": {
"Id": {
"Value": 94358348534
}
}
},
"RuntimePlayer": {
"Name": "player1"
},
"MaxPlayerSlots": 2,
"SnapshotsBlocked": true
}

```

### JoinGame

The ```
JoinGame
```

webhook is send before a client joins an existing room/game session. Return ```
200
```

to allow or ```
400
```

to cancel the joining.

Requires ```
WebHookBaseUrl
```

 and ```
WebHookEnableOnJoin
```

to be set on the Photon AppId Dashboard.

JavaScript

```javascript
POST https://{WebHookBaseUrl}/game/join

```

**JoinGame Request**

| Name | Type | Example | Description |
| --- | --- | --- | --- |
| ```<br>AppId<br>``` | ```<br>string<br>``` | ```<br>d1f67eec-51fb-45c1<br>``` | The Photon AppId |
| ```<br>GameId<br>``` | ```<br>string<br>``` | ```<br>0:eu:db757806-8570-45aa<br>``` | Unique ```<br>GameId<br>``` |
| ```<br>UserId<br>``` | ```<br>string<br>``` | ```<br>db757806-8570-45aa<br>``` | Photon ```<br>UserId<br>``` |
| ```<br>AuthCookie<br>``` | ```<br>Dictionary<string, object><br>``` | ```<br>db757806-8570-45aa<br>``` | The Photon custom authentication cookie set by the backend. |

Json Example:

JSON

```json
{
 "AppId": "\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*",
 "GameId": "0:eu:db757806-8570-45aa",
 "UserId": "db757806-8570-45aa",
 "AuthCookie": {
 "Secret": "\*\*\*\*\*\*\*\*\*\*"
 }
}

```

**HTTP Response Codes**

| Name | Type | Description |
| --- | --- | --- |
| ```<br>200 OK<br>``` | ```<br>[JoinGame Response](#joingame_response)<br>``` | The client will join the room. |
| ```<br>400 Bad Request<br>``` | ```<br>[WebhookError](#webhookerror)<br>``` | Joining the room will fail. |

**JoinGame Response**

| Name | Type | Description |
| --- | --- | --- |
| ```<br>RuntimePlayer<br>``` | ```<br>[RuntimePlayer](#runtimeplayer)<br>``` | Return a ```<br>RuntimePlayer<br>```<br> object that will be used for the client that created this room/session.<br>This only overwrites the first ```<br>AddPlayer<br>```<br> data sent for player slot ```<br>0<br>```<br>. ```<br>MaxPlayerSlots<br>```<br> should be set to ```<br>1<br>```<br>. |
| ```<br>MaxPlayerSlots<br>``` | ```<br>int<br>``` | The maximum number of player slots this client can acquire:<br>```<br>0<br>```<br> = only spectating<br>```<br>1..255<br>```<br> = specific number<br>```<br>-1<br>```<br> = unlimited<br>If this response is send but this value is not set ```<br>MaxPlayerSlots<br>```<br> will default to 1.<br>Players requesting an invalid player slot number or more slots than allowed will be disconnected. |

Json Example:

JSON

```json
{
 "RuntimePlayer": {
 "Name": "player1"
 },
 "MaxPlayerSlots": 1
}

```

_Hint_

Rather use ```
AddPlayer
```

 to return the ```
RuntimePlayer
```

data, because its request will includes a ```
RuntimePlayer
```

 client object sent by the client.

Also consider using the ```
CreateGame
```

response to reserve player slots for users (if they are already known at that time).

### LeaveGame

The ```
LeaveGame
```

 webhook is send after the client leaves an existing room/game session.

Requires ```
WebHookBaseUrl
```

and ```
WebHookEnableOnLeave
```

to be set on the Photon AppId Dashboard.

JavaScript

```javascript
POST https://{WebHookBaseUrl}/game/leave

```

**LeaveGame Request**

| Name | Type | Example | Description |
| --- | --- | --- | --- |
| ```<br>AppId<br>``` | ```<br>string<br>``` | ```<br>d1f67eec-51fb-45c1<br>``` | The Photon AppId |
| ```<br>GameId<br>``` | ```<br>string<br>``` | ```<br>0:eu:db757806-8570-45aa<br>``` | Unique ```<br>GameId<br>``` |
| ```<br>UserId<br>``` | ```<br>string<br>``` | ```<br>db757806-8570-45aa<br>``` | Photon ```<br>UserId<br>``` |
| ```<br>ActorNr<br>``` | ```<br>int<br>``` | ```<br>db757806-8570-45aa<br>``` | The Photon actor number, a incrementing runtime id for clients. |
| ```<br>AuthCookie<br>``` | ```<br>Dictionary<string, object><br>``` | ```<br>db757806-8570-45aa<br>``` | Photon ```<br>UserId<br>``` |
| ```<br>IsInactive<br>``` | ```<br>bool<br>``` | ```<br>db757806-8570-45aa<br>``` | Is set to true when the player left the room but is still marked inactive, e.g. when Player TTL was set. Additional LeaveGame requests can follow in this case. |

Json Example:

JSON

```json
{
"AppId": "\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*",
"GameId": "0:eu:db757806-8570-45aa",
"UserId": "db757806-8570-45aa",
"ActorNr": 1,
"AuthCookie": {
"Secret": "\*\*\*\*\*\*\*\*\*\*"
},
"IsInactive": false
}

```

**HTTP Response Codes**

| Name | Type | Description |
| --- | --- | --- |
| ```<br>200 OK<br>``` | ```<br>[LeaveGame Response](#leavegame_response)<br>``` | Just a confirmation of receipt. |
| ```<br>400 Bad Request<br>``` | ```<br>[WebhookError](#webhookerror)<br>``` | Error is ignored, it will just be logged on the Photon Cloud. |

**LeaveGame Response**

Json Example:

JSON

```json
{
// empty
}

```

### CloseGame

The ```
CloseGame
```

webhook is sent when the room/game session is closed after all clients have left.

Requires ```
WebHookBaseUrl
```

 and ```
WebHookEnableOnClose
```

to be set on the Photon dashboard.

JavaScript

```javascript
POST https://{WebHookBaseUrl}/game/close

```

**CloseGame Request**

| Name | Type | Example | Description |
| --- | --- | --- | --- |
| ```<br>AppId<br>``` | ```<br>string<br>``` | ```<br>d1f67eec-51fb-45c1<br>``` | The Photon AppId |
| ```<br>GameId<br>``` | ```<br>string<br>``` | ```<br>0:eu:db757806-8570-45aa<br>``` | Unique game id |
| ```<br>CloseReason<br>``` | ```<br>int ( [CloseReason](#closereason)<br>```<br>) | ```<br>0<br>``` | The reason why this room/session has been closed. |

Json Example:

JSON

```json
{
 "GameId": "0:eu:db757806-8570-45aa",
 "CloseReason": 0
}

```

**HTTP Response Codes**

| Name | Type | Description |
| --- | --- | --- |
| ```<br>200 OK<br>``` |  | Confirmation of receipt. |

#### CloseReason

| Name | Value | Description |
| --- | --- | --- |
| ```<br>Ok<br>``` | ```<br>0<br>``` | Session was closed without errors. |
| ```<br>FailedOnCreate<br>``` | ```<br>1<br>``` | Session was closed because it failed to Create. |

### GameConfigs

The ```
GameConfigs
```

 webhook is sent when a player sent a ```
StartRequest
```

operation attached with the game configs ```
RuntimeConfig
```

 and ```
SessionConfig
```

. Both config files are attached to the request body as Json objects.

This webhook is only send once per room upon the first arrival of any clients start operation.

Requires ```
WebHookBaseUrl
```

 and ```
WebHookEnableGameConfigs
```

to be set on the Photon dashboard.

JavaScript

```javascript
POST https://{WebHookBaseUrl}/game/configs

```

**GameConfigs Request**

| Name | Type | Example | Description |
| --- | --- | --- | --- |
| ```<br>AppId<br>``` | ```<br>string<br>``` | ```<br>d1f67eec-51fb-45c1<br>``` | The Photon AppId |
| ```<br>GameId<br>``` | ```<br>string<br>``` | ```<br>0:eu:db757806-8570-45aa<br>``` | Unique game id |
| ```<br>UserId<br>``` | ```<br>string<br>``` | ```<br>db757806-8570-45aa<br>``` | Photon ```<br>UserId<br>``` |
| ```<br>ActorNr<br>``` | ```<br>int<br>``` | ```<br>db757806-8570-45aa<br>``` | The Photon actor number, a incrementing runtime id for clients. |
| ```<br>RuntimeConfig<br>``` | ```<br>[RuntimeConfig](#runtimeconfig)<br>``` | ```<br>{ "Level": 1 }<br>``` | The ```<br>RuntimeConfig<br>```<br> object sent by the client to the plugin. Can be null. |
| ```<br>SessionConfig<br>``` | ```<br>[SessionConfig](#sessionconfig)<br>``` | ```<br>{ "PlayerCount": 8, .. }<br>``` | The ```<br>SessionConfig<br>```<br> object sent by the client to the plugin. Can be null. |
| ```<br>AuthCookie<br>``` | ```<br>Dictionary<string, object><br>``` | ```<br>db757806-8570-45aa<br>``` | The Photon custom authentication cookie set by the backend. |

Json Example:

JSON

```json
{
 "AppId": "d1f67eec-51fb-45c1",
 "GameId": "0:eu:db757806-8570-45aa",
 "UserId": "db757806-8570-45aa",
 "ActorNr": 1,
 "RuntimeConfig": {
 "Map": {
 "Id": {
 "Value": 94358348534
 }
 }
 },
 "SessionConfig": {
 "PlayerCount": 8,
 "ChecksumCrossPlatformDeterminism": false,
 "LockstepSimulation": false,
 //..
 },
 "AuthCookie": {
 "Secret": "\*\*\*\*\*\*\*\*\*\*"
 }
}

```

**HTTP Response Codes**

| Name | Type | Description |
| --- | --- | --- |
| ```<br>200 OK<br>``` | ```<br>[GameConfigs Response](#gameconfigs_response)<br>``` | The game start sequence can continue. Game configs attached to response should be overwritten. |
| ```<br>400 Bad Request<br>``` | ```<br>[WebhookError](#webhookerror)<br>``` | The game will be terminated and all clients are disconnected. |

**GameConfigs Response**

Both objects on the response can be ```
null
```

 to accept the configs that the client sent. Otherwise they will be overwritten.

| Name | Type | Description |
| --- | --- | --- |
| ```<br>RuntimeConfig<br>``` | ```<br>[RuntimeConfig](#runtimeconfig)<br>``` | The ```<br>RuntimeConfig<br>```<br> object to overwrite the one the client sent. |
| ```<br>SessionConfig<br>``` | ```<br>[SessionConfig](#sessionconfig)<br>``` | The ```<br>SessionConfig<br>```<br> object to overwrite the one the client sent. |

Json Example:

JSON

```json
{
"RuntimeConfig": {
"Map": {
"Id": {
"Value": 94358348534
}
}
},
"SessionConfig": {
"PlayerCount": 8,
"ChecksumCrossPlatformDeterminism": false,
"LockstepSimulation": false,
//..
}
}

```

### AddPlayer

The ```
AddPlayer
```

webhook is sent when a client tries to add a player to the Quantum online game using the ```
AddPlayer
```

operation.

Adding a player to the online game can still fail after this webhook returns when no player slot is free then. Use the ```
PlayerAdded
```

webhook to track the player online status.

Requires ```
WebHookBaseUrl
```

 and ```
WebHookEnableAddPlayer
```

to be set on the Photon dashboard.

JavaScript

```javascript
POST https://{WebHookBaseUrl}/player/add

```

**AddPlayer Request**

| Name | Type | Example | Description |
| --- | --- | --- | --- |
| ```<br>AppId<br>``` | ```<br>string<br>``` | ```<br>d1f67eec-51fb-45c1<br>``` | The Photon AppId |
| ```<br>GameId<br>``` | ```<br>string<br>``` | ```<br>0:eu:db757806-8570-45aa<br>``` | Unique game id |
| ```<br>UserId<br>``` | ```<br>string<br>``` | ```<br>db757806-8570-45aa<br>``` | Photon UserId or ClientId |
| ```<br>ActorNr<br>``` | ```<br>int<br>``` | ```<br>db757806-8570-45aa<br>``` | The Photon actor number, a incrementing runtime id for clients. |
| ```<br>PlayerSlot<br>``` | ```<br>int<br>``` | ```<br>0<br>``` | Player slot requested. Usually is 0. |
| ```<br>RuntimePlayer<br>``` | ```<br>[RuntimePlayer](#runtimeplayer)<br>``` | ```<br>{ "Foo": 222 }<br>``` | The ```<br>RuntimePlayer<br>```<br> object sent by the client to the plugin. Can be null. |
| ```<br>AuthCookie<br>``` | ```<br>Dictionary<string, object><br>``` | ```<br>db757806-8570-45aa<br>``` | The Photon custom authentication cookie set by the backend. |

Json Example:

JSON

```json
{
 "AppId": "d1f67eec-51fb-45c1",
 "GameId": "0:eu:db757806-8570-45aa",
 "UserId": "db757806-8570-45aa",
 "ActorNr": 1,
 "PlayerSlot": 0,
 "RuntimePlayer": {
 "Name": "player1"
 },
 "AuthCookie": {
 "Secret": "\*\*\*\*\*\*\*\*\*\*"
 }
}

```

**HTTP Response Codes**

| Name | Type | Description |
| --- | --- | --- |
| ```<br>200 OK<br>``` | ```<br>[AddPlayer Response](#addplayer_response)<br>``` | The client can add a player to the selected player slot and optionally received a ```<br>RuntimePlayer<br>```<br> object from the backend. |
| ```<br>400 Bad Request<br>``` | ```<br>[WebhookError](#webhookerror)<br>``` | The client cannot add the player and will receive an error protocol message and callback: ```<br>OnLocalPlayerAddFailed<br>```<br>. |

**AddPlayer Response**

| Name | Type | Description |
| --- | --- | --- |
| ```<br>RuntimePlayer<br>``` | ```<br>object<br>``` | The ```<br>RuntimePlayer<br>```<br> object to overwrite the ```<br>RuntimePlayer<br>```<br> that the client sent. If ```<br>null<br>```<br> the clients ```<br>RuntimePlayer<br>```<br> will be accepted. |

Json Example:

JSON

```json
{
 "RuntimePlayer": {
 "Name": "player1"
 }
}

```

### PlayerAdded

The ```
PlayerAdded
```

 webhook is send after a client successfully added a player to the online game.

Requires ```
WebHookBaseUrl
```

and ```
WebHookEnablePlayerAdded
```

to be set on the Photon dashboard.

JavaScript

```javascript
POST https://{WebHookBaseUrl}/player/added

```

**PlayerAdded Request**

| Name | Type | Example | Description |
| --- | --- | --- | --- |
| ```<br>AppId<br>``` | ```<br>string<br>``` | ```<br>d1f67eec-51fb-45c1<br>``` | The Photon AppId |
| ```<br>GameId<br>``` | ```<br>string<br>``` | ```<br>0:eu:db757806-8570-45aa<br>``` | Unique game id |
| ```<br>UserId<br>``` | ```<br>string<br>``` | ```<br>db757806-8570-45aa<br>``` | Photon UserId or ClientId |
| ```<br>ActorNr<br>``` | ```<br>int<br>``` | ```<br>db757806-8570-45aa<br>``` | The Photon actor number, a incrementing runtime id for clients. |
| ```<br>PlayerSlot<br>``` | ```<br>int<br>``` | 0 | The (local) player slot that the client reserved. |
| ```<br>Player<br>``` | ```<br>int<br>``` | 21 | The (global) Player id that the client received. |
| ```<br>AuthCookie<br>``` | ```<br>Dictionary<string, object><br>``` | ```<br>db757806-8570-45aa<br>``` | The Photon custom authentication cookie set by the backend. |

Json Example:

JSON

```json
{
"AppId": "d1f67eec-51fb-45c1",
"GameId": "0:eu:db757806-8570-45aa",
"UserId": "db757806-8570-45aa",
"ActorNr": 1,
"PlayerSlot": 0,
"Player": 21,
"AuthCookie": {
"Secret": "\*\*\*\*\*\*\*\*\*\*"
}
}

```

**HTTP Response Codes**

| Name | Type | Description |
| --- | --- | --- |
| ```<br>200 OK<br>``` |  | Confirmation of receipt. |

### PlayerRemoved

The ```
PlayerRemoved
```

webhook is sent when a client was removed from the online game.

Requires ```
WebHookBaseUrl
```

 and ```
WebHookEnablePlayerRemoved
```

to be set on the Photon dashboard.

JavaScript

```javascript
POST https://{WebHookBaseUrl}/player/removed

```

**PlayerRemoved Request**

| Name | Type | Example | Description |
| --- | --- | --- | --- |
| ```<br>AppId<br>``` | ```<br>string<br>``` | ```<br>d1f67eec-51fb-45c1<br>``` | The Photon AppId |
| ```<br>GameId<br>``` | ```<br>string<br>``` | ```<br>0:eu:db757806-8570-45aa<br>``` | Unique game id |
| ```<br>UserId<br>``` | ```<br>string<br>``` | ```<br>db757806-8570-45aa<br>``` | Photon UserId or ClientId |
| ```<br>ActorNr<br>``` | ```<br>int<br>``` | ```<br>db757806-8570-45aa<br>``` | The Photon actor number, a incrementing runtime id for clients. |
| ```<br>PlayerSlot<br>``` | ```<br>int<br>``` | 0 | The (local) player slot. |
| ```<br>Player<br>``` | ```<br>int<br>``` | 21 | The (global) Player id. |
| ```<br>Reason<br>``` | ```<br>int<br>``` | 0 | The reason why the player was removed from the game.<br>0 = ```<br>Requested<br>```<br>1 = ```<br>ClientDisconnected<br>```<br>2 = ```<br>Error<br>``` |
| ```<br>AuthCookie<br>``` | ```<br>Dictionary<string, object><br>``` | ```<br>db757806-8570-45aa<br>``` | The Photon custom authentication cookie set by the backend. |

Json Example:

JSON

```json
{
 "AppId": "d1f67eec-51fb-45c1",
 "GameId": "0:eu:db757806-8570-45aa",
 "UserId": "db757806-8570-45aa",
 "ActorNr": 1,
 "PlayerSlot": 0,
 "Player": 21,
 "Reason": 0,
 "AuthCookie": {
 "Secret": "\*\*\*\*\*\*\*\*\*\*"
 }
}

```

**HTTP Response Codes**

| Name | Type | Description |
| --- | --- | --- |
| ```<br>200 OK<br>``` |  | Confirmation of receipt. |

### GameResult

Using the ```
GameResult
```

 Quantum event in the simulation will trigger the upload of a [GameResult](#gameresult-event) instance by clients to the Quantum server once per game.

The results are aggregated and forwarded as the ```
GameResult
```

webhook when the game/room is closed.

Requires ```
WebHookBaseUrl
```

 and ```
WebHookEnableGameResult
```

to be set on the Photon dashboard.

If the server is running the simulation the webhooks are executed immediately from the server simulation and can be used as a reliable and trusted source of game results.

JavaScript

```javascript
POST https://{WebHookBaseUrl}/game/result

```

**GameResult Request**

| Name | Type | Example | Description |
| --- | --- | --- | --- |
| ```<br>AppId<br>``` | ```<br>string<br>``` | ```<br>d1f67eec-51fb-45c1<br>``` | The Photon AppId |
| ```<br>GameId<br>``` | ```<br>string<br>``` | ```<br>0:eu:db757806-8570-45aa<br>``` | Unique game id |
| ```<br>Results<br>``` | ```<br>[GameResultInfo](#gameresultinfo)\[\]<br>``` | see below | The aggregated game results |

Json Example:

JSON

```json
{
 "AppId": "d1f67eec-51fb-45c1",
 "GameId": "0:eu:db757806-8570-45aa",
 "Results": \[
 {
 "Clients": \[
 {
 "UserId": "FJEH43FL56FSDR",
 "Players": \[
 0
 \],
 "GameTime": 63.3636703
 }
 \],
 "Result": {
 "$type": "Quantum.GameResult, Quantum.Simulation",
 "Frame": 12010
 "Winner": 2
 },
 "IsServerResult": false
 }
 \],
 "UserId": "0"
}

```

**HTTP Response Codes**

| Name | Description |
| --- | --- |
| ```<br>200 OK<br>``` | Confirmation of receipt |

**GameResultInfo**

| Name | Type | Description |
| --- | --- | --- |
| ```<br>Result<br>``` | ```<br>[GameResult](#gameresult-event)<br>``` | The game specific game result Json object that the listed clients sent |
| ```<br>Clients<br>``` | ```<br>[GameResultClientInfo](#gameresultclientinfo)\[\]<br>``` | The list of clients that generated this result |
| ```<br>IsServerResult<br>``` | ```<br>bool<br>``` | Has this result been generated by server simulation |

**GameResultClientInfo**

| Name | Type | Description |
| --- | --- | --- |
| ```<br>UserId<br>``` | ```<br>string<br>``` | The Photon user id |
| ```<br>Players<br>``` | ```<br>int\[\]<br>``` | The list of players that this user controls |
| ```<br>GameTime<br>``` | ```<br>float<br>``` | The timestamp when the result reached the server |

### ReplayStart

The ```
ReplayStart
```

 webhook is sent when the simulation and the input recording is starting on the server. It's is a trusted source for capturing the game replay directly from the server instead of relying on clients to send it to a developers backend.

This webhook has to be answered with the ```
ReplayStartResponse
```

which can signal the replay streaming to be skipped for this particular game session.

A response has to be received by the Quantum server before it is sending its first replay slice or the replay recording will be canceled.

Requires ```
WebHookBaseUrl
```

 and ```
WebHookEnableReplay
```

to be set on the Photon dashboard.

JavaScript

```javascript
POST https://{WebHookBaseUrl}/replay/start

```

**ReplayStart Request**

| Name | Type | Description |
| --- | --- | --- |
| ```<br>AppId<br>``` | ```<br>string<br>``` | The Photon AppId. |
| ```<br>AppVersion<br>``` | ```<br>string<br>``` | The AppVersion used when creating the room/game session. |
| ```<br>Region<br>``` | ```<br>string<br>``` | The Region code of the Game Server that the room/game session was created in. |
| ```<br>Cloud<br>``` | ```<br>string<br>``` | The ```<br>Cloud Id<br>```<br> of that the Game Server is running on. |
| ```<br>RoomName<br>``` | ```<br>string<br>``` | The room/game session Name. |
| ```<br>GameId<br>``` | ```<br>string<br>``` | A unique ```<br>GameId<br>```<br> which is composed of ```<br>{Cloud:}{Region:}RoomName<br>```<br>. |
| ```<br>SessionConfig<br>``` | ```<br>[SessionConfig](#sessionconfig)<br>``` | The ```<br>SessionConfig<br>```<br> that the simulation started with. |
| ```<br>RuntimeConfig<br>``` | ```<br>byte\[\]<br>``` | The GZipped Json of the ```<br>[RuntimeConfig](#runtimeconfig)<br>```<br> that the simulation started with. |

Json Example:

JSON

```json
{
 "AppId": "d1f67eec-51fb-45c1",
 "AppVersion": "1.0-live",
 "Region": "eu",
 "Cloud": "0:",
 "RoomName": "1.2-party-2349535735",
 "GameId": "0:eu:e472a861-a1e2-49f7",
 "SessionConfig": { },
 "RuntimeConfig": "H4sIAAAAAAAACnWNPQvCMBCG/4ocjkWuye
 X6sXZycNCCe8AogSYtNBlK6X/3UHQR4Zb35X2eW2GflslBC+ds
 Y8rhcMkx+eC6Md79o9h96t6HPNjkxwgF9M7doMUCTnaCdoWjpB
 WudshiUkxYsVHacE11KWe2TZiv4K3+4YjQkECKNJcVMuoXtszJ
 hfkPI1tUVc1sqFGN1g3Kr+0J+3ktedUAAAA="
}

```

**HTTP Response Codes**

| Name | Type | Description |
| --- | --- | --- |
| ```<br>200 OK<br>``` | ```<br>[ReplayStart Response](#replaystart_response)<br>``` | The replay streaming can start or be disabled when the ```<br>Skip<br>```<br> property is set. |
| ```<br>400 Bad Request<br>``` | ```<br>[WebhookError](#webhookerror)<br>``` | In this case an error is logged on the server and the replay streaming is stopped. |

**ReplayStart Response**

| Name | Type | Description |
| --- | --- | --- |
| ```<br>Skip<br>``` | ```<br>bool<br>``` | The replay streaming is disabled for this game session. |

Json Example:

JSON

```json
{
 "Skip": true
}

```

### ReplayChunk

The ```
ReplayChunk
```

 webhook is sent in intervals and it contains the delta compressed input history of a part of the simulation.

Requires ```
WebHookBaseUrl
```

and ```
WebHookEnableReplay
```

to be set on the Photon dashboard.

JavaScript

```javascript
POST https://{WebHookBaseUrl}/replay/chunk

```

There are additional dashboard variables that configure the replay input streaming.

| Dashboard Variable | Type | Description | Default |
| --- | --- | --- | --- |
| ```<br>WebHookReplayUseBinaryRequests<br>``` | ```<br>bool<br>``` | The web request are not send as JSON content but as binary data. | ```<br>false<br>``` |
| ```<br>WebHookReplayUseCompression<br>``` | ```<br>bool<br>``` | The input data on the chunks is GZip compressed. | ```<br>false<br>``` |
| ```<br>WebHookReplaySendIntervalInSec<br>``` | ```<br>int<br>``` | The chunk send interval in seconds (min = 2, max = 40). | ```<br>20<br>``` |

**ReplayChunk Request**

If ```
WebHookReplayUseBinaryRequests
```

is selected, then the following properties are part of the web request ```
Headers
```

instead of a JSON body. The body then will include the binary ```
Input
```

.

| Name | Type | Description |
| --- | --- | --- |
| ```<br>AppId<br>``` | ```<br>string<br>``` | The Photon AppId |
| ```<br>GameId<br>``` | ```<br>string<br>``` | The game id to identify the replay chunk. |
| ```<br>ChunkNumber<br>``` | ```<br>int<br>``` | The incrementing chunk number. |
| ```<br>IsLast<br>``` | ```<br>bool<br>``` | A flag that indicates that the last chunk has been send. Usually when a room closes. |
| ```<br>LastTick<br>``` | ```<br>int<br>``` | The oldest tick of input on this chunk. |
| ```<br>TickCount<br>``` | ```<br>int<br>``` | The number of ticks of input on this chunk. |
| ```<br>TickCountTotal<br>``` | ```<br>int<br>``` | The total incrementing number of ticks in the whole replay. |
| ```<br>IsCompressed<br>``` | ```<br>bool<br>``` | A flag indicating the input is GZip compressed. |
| ```<br>Input<br>``` | ```<br>byte\[\]<br>``` | The binary delta compressed input that needs to be appended to the complete input stream for this replay.<br> The input for each tick has a leading int describing the data length. It can be stored in the ```<br>ReplayFile.InputHistoryRaw<br>```<br> field to be readable with the ```<br>QuantumRunnerLocalReplay<br>```<br> script in Unity. |

Json Example:

JSON

```json
{
 "AppId": "d1f67eec-51fb-45c1",
 "GameId": ":eu:e472a861-a1e2-49f7",
 "ChunkNumber": 0,
 "IsLast": false,
 "LastTick": 302,
 "TickCount": 243,
 "TickCountTotal": 243,
 "IsCompressed": false,
 "Input": "JQAAADwAAAAIAAMKFsCUwYDggOB/UCAAAPgfDMhgEoAHA////4PRBwATAAAAPQAAAAgAA2PaSK"
}

```

**HTTP Response Codes**

| Name | Type | Description |
| --- | --- | --- |
| ```<br>200 OK<br>``` |  | Confirmation of receipt. |

### WebhookError

| Name | Type | Description |
| --- | --- | --- |
| ```<br>Status<br>``` | ```<br>int<br>``` | HTTP status code |
| ```<br>Error<br>``` | ```<br>string<br>``` | Error name |
| ```<br>Message<br>``` | ```<br>string<br>``` | Error message |

Json example:

JSON

```json
{
 "Status": 400,
 "Error": "PlayerNotAllowed",
 "Message": "LoremIpsum"
}

```

## Quantum Classes

### RuntimeConfig

Quantum 3 runtime configuration files ```
RuntimeConfig
```

 and ```
RuntimePlayer
```

are uploaded by the clients using Json serialization. This way it is possible to send configurations to the Quantum public cloud game servers while it does not know the implementation and serialization details.

When ```
RuntimeConfig
```

 or ```
RuntimePlayer
```

are used in a webhook response they need to be **complete** because the configs send by clients are completely replaced.

Json data sent by the custom backend has to be deterministically deserializable on every client.

Quantum internal classes usually operate with fields, make sure that the Json serialization and deserialization code needs to ```
IncludeFields = true
```

 and ```
IgnoreReadOnlyProperties = true
```

.

A simple example of the ```
RuntimeConfig
```

 class. Additionally to the partial declaration the base class adds a couple properties as well (Seed, Map, SimulationConfig, SystemsConfig).

The ```
$type
```

property is required for deserialization on the standalone or custom plugin.

C#

```csharp
namespace Quantum {
 public partial class RuntimeConfig {
 public int GameMode;
 }
}

```

Json Example:

JSON

```json
{
 "$type": "Quantum.RuntimeConfig, Quantum.Simulation",
 "GameMode": 1,
 "Seed": 0,
 "Map": {
 "Id": {
 "Value":2640765235684814815
 }
 },
 "SimulationConfig": {
 "Id": {
 "Value":440543562436170603
 }
 },
 "SystemsConfig": {
 "Id": {
 "Value":2430278665492933905
 }
 }
}

```

### RuntimePlayer

C#

```csharp
namespace Quantum {
 public partial class RuntimePlayer {
 public AssetRef<GearConfig> Loadout;
 }
}

```

Json Example:

JSON

```json
{
 "$type":"Quantum.RuntimePlayer, Quantum.Simulation",
 "Loadout": {
 "Id": {
 "Value": 440543562436170603
 }
 },
 "PlayerAvatar": {
 "Id": {
 "Value": 2430278665492933905
 }
 },
 "PlayerNickname": "foo"
}

```

### SessionConfig

```
SessionConfig
```

is the abbreviation of the ```
DeterministicSessionConfig
```

 class.

When a ```
SessionConfig
```

is returned by an webhook it needs to be **complete** because single values are not replaced.

The current ```
SessionConfig
```

 asset can be exported in Unity using this menu entry:

```
Unity Editor > Quantum > Export > SessionConfig (Json)

```

When serializing the class on ```
netcoreapp3.1
```

either use ```
Newtonsoft
```

or use ```
DeterministicSessionConfigJsonConverter
```

for ```
Text.Json
```

. Because "including fields" is a feature of ```
net5
```

.

Json Example:

JSON

```json
{
 "PlayerCount": 8,
 "ChecksumCrossPlatformDeterminism": false,
 "LockstepSimulation": false,
 "InputDeltaCompression": true,
 "UpdateFPS": 60,
 "ChecksumInterval": 60,
 "RollbackWindow": 60,
 "InputHardTolerance": 8,
 "InputRedundancy": 3,
 "InputRepeatMaxDistance": 10,
 "SessionStartTimeout": 1,
 "TimeCorrectionRate": 4,
 "MinTimeCorrectionFrames": 1,
 "MinOffsetCorrectionDiff": 1,
 "TimeScaleMin": 100,
 "TimeScalePingMin": 100,
 "TimeScalePingMax": 300,
 "InputDelayMin": 0,
 "InputDelayMax": 60,
 "InputDelayPingStart": 100,
 "InputFixedSizeEnabled": true,
 "InputFixedSize": 24
}

```

### GameResult Event

The ```
GameResult
```

 class can be extended by adding fields using the partial class declaration inside the ```
GameResult.User.cs
```

script. It will be Json serialized on the client and send to the Quantum server.

C#

```csharp
namespace Quantum {
 public partial class GameResult {
 public int Winner;
 }
}

```

Json Example:

JSON

```json
{
 "$type":"Quantum.GameResult, Quantum.Simulation",
 "Frame": 200,
 "Winner": 2
}

```

The game result event can be invoked **once** per game by each client, which is triggered from inside the simulation by calling the GameResult Quantum event. The actual webhook is launched when the room is disbanded and closed.

C#

```csharp
f.Events.GameResult(new GameResult { Winner = 3 });

```

When running an enterprise Photon Quantum cloud with server simulation this will be transparently send by the server for a trusted source of game results.

## Photon Realtime Classes

### EnterRoomParams

This definitions is designed to be similar to the ```
EnterRoomParams
```

 class in Photon Realtime. It includes all options that can be set by the ```
CreateGame
```

webhook. When composing the Json Response every member is optional and can be null or not set.

| Name | Type | Description |
| --- | --- | --- |
| ```<br>RoomOptions<br>``` | ```<br>[RoomOptions](#roomoptions)<br>``` | RoomOptions object |
| ```<br>ExpectedUsers<br>``` | ```<br>string\[\]<br>``` | A list of ```<br>UserIds<br>```<br> that are permitted to enter the room/session (additionally to the user creating the room). If ```<br>MaxPlayers<br>```<br> is greater than the number ```<br>ExpectedUsers<br>```<br> listed, any player may join and fill the unreserved slots.<br>Works only for ```<br>RoomJoin()<br>```<br> not ```<br>JoinRandom()<br>```<br>. |

Json example:

JSON

```json
{
 "RoomOptions": {
 "IsVisible": true,
 "IsOpen": true,
 "MaxPlayers": 8,
 "PlayerTtl": null,
 "EmptyRoomTtl": 10000,
 "CustomRoomProperties": {
 "Foo": "bar",
 "PlayerClass": 1
 },
 "CustomRoomPropertiesForLobby": \[
 "Foo"
 \],
 "SuppressRoomEvents": null,
 "SuppressPlayerInfo": null,
 "PublishUserId": null,
 "DeleteNullProperties": null,
 "BroadcastPropsChangeToAll": null,
 "CleanupCacheOnLeave": null,
 "CheckUserOnJoin": null
 },
 "ExpectedUsers": \[
 "A",
 "B",
 "C"
 \]
}

```

### RoomOptions

All values are nullable types and can be set to ```
null
```

 or be omitted when sending back to the Quantum server, in which case this response will not alter the nulled or omitted room property and will remain the default or the value send by the client when creating the room.

| Name | Type | Description |
| --- | --- | --- |
| ```<br>IsVisible<br>``` | ```<br>bool<br>``` | Defines if this room is listed in the Photon matchmaking. |
| ```<br>IsOpen<br>``` | ```<br>bool<br>``` | Defines if this room can be joined by other clients. |
| ```<br>MaxPlayers<br>``` | ```<br>byte<br>``` | Max number of players that can be in the room at any time. 0 means "no limit". |
| ```<br>PlayerTtl<br>``` | ```<br>int<br>``` | Time To Live (TTL) for an 'actor' in a room. If a client disconnects, this actor is inactive first and removed after this timeout. In milliseconds. |
| ```<br>EmptyRoomTtl<br>``` | ```<br>int<br>``` | Time To Live (TTL) for a room when the last player leaves. Keeps room in memory for case a player re-joins soon. In milliseconds. |
| ```<br>CustomRoomProperties<br>``` | ```<br>Dictionary <string, object><br>``` | The room's custom properties to set during creation. |
| ```<br>CustomRoomPropertiesForLobby<br>``` | ```<br>string\[\]<br>``` | Defines which of the custom room properties get listed in the lobby.<br>Value type of the properties has to be ```<br>bool<br>```<br>, ```<br>byte<br>```<br>, ```<br>short<br>```<br>, ```<br>int<br>```<br>, ```<br>long<br>```<br>or ```<br>string<br>```<br>.<br>Max number of properties is **3**.<br>The max length for string value is **64**.<br>Key restrictions can also be enforced by the Photon dashboard property: ```<br>AllowedLobbyProperties<br>```<br>. |
| ```<br>SuppressRoomEvents<br>``` | ```<br>bool<br>``` | Tells the server to skip room events for joining and leaving players.<br>Default is ```<br>false<br>```<br>. |
| ```<br>SuppressPlayerInfo<br>``` | ```<br>bool<br>``` | Disables events join and leave from the server as well as property broadcasts in a room (to minimize traffic).<br>Default is ```<br>false<br>```<br>. |
| ```<br>PublishUserId<br>``` | ```<br>bool<br>``` | Defines if the UserIds of players get "published" in the room. Useful for FindFriends, if players want to play another game together.<br>Default is ```<br>false<br>```<br>. |
| ```<br>DeleteNullProperties<br>``` | ```<br>bool<br>``` | Optionally, properties get deleted, when null gets assigned as value.<br>Default is ```<br>false<br>```<br>. |

## PlayFab

To enable the PlayFab integration add the ```
WebHookIntegration
```

dashboard variable and set it to ```
PlayFab
```

.

- All slashes from all paths will be substituted with an underscore: ```
  game/create
  ```

   =\> ```
  game\_create
  ```

- All webhook requests will automatically add the ```
  AppId
  ```

   property (if missing) controlled by ```
  WebHooksConfig.AppId
  ```

- All webhook requests will automatically add the ```
  UserId
  ```

   property (if missing) ```
  "0"
  ```

- Webhook responses will handle the following Json result body and cause the webhooks to fail on ```
  ResultCode != 0
  ```

  . The ```
  Message
  ```

   is copied on errors to ```
  WebHookError.Message
  ```

  .

JSON

```json
{
"ResultCode": 0,
"Message": "success"
}

```

## Photon Cloud Web Request Limitations

By default Photon Server manages HTTP request **serially** and a new request are not started until the current request completes. New requests are queued. This limitation is per room.

This could lead to unexpected wait times for clients when using OnJoin or AddPlayer webhooks, while having a high player influx and a high roundtrip time from the Photon cloud to a custom backend.

Parallel request can be enabled for enterprise clouds. And we working on a solution to have them opt-in for public could apps as well.

Other limitations are:

- ```
  HttpRequestTimeout
  ```

  : 30000
- ```
  LimitHttpResponseMaxSize
  ```

  : 200000
- ```
  MaxQueuedRequests
  ```

  : 5000

Back to top

- [Overview](#overview)
- [Setup](#setup)

  - [Dashboard Configurations](#dashboard-configurations)
  - [Quantum Dashboard Configurations](#quantum-dashboard-configurations)

- [Webhook API](#webhook-api)

  - [Http Request Retries](#http-request-retries)
  - [Common Request Headers](#common-request-headers)
  - [CreateGame](#creategame)
  - [JoinGame](#joingame)
  - [LeaveGame](#leavegame)
  - [CloseGame](#closegame)
  - [GameConfigs](#gameconfigs)
  - [AddPlayer](#addplayer)
  - [PlayerAdded](#playeradded)
  - [PlayerRemoved](#playerremoved)
  - [GameResult](#gameresult)
  - [ReplayStart](#replaystart)
  - [ReplayChunk](#replaychunk)
  - [WebhookError](#webhookerror)

- [Quantum Classes](#quantum-classes)

  - [RuntimeConfig](#runtimeconfig)
  - [RuntimePlayer](#runtimeplayer)
  - [SessionConfig](#sessionconfig)
  - [GameResult Event](#gameresult-event)

- [Photon Realtime Classes](#photon-realtime-classes)

  - [EnterRoomParams](#enterroomparams)
  - [RoomOptions](#roomoptions)

- [PlayFab](#playfab)
- [Photon Cloud Web Request Limitations](#photon-cloud-web-request-limitations)