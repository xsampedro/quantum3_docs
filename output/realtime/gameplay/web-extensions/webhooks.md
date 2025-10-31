# webhooks

_Source: https://doc.photonengine.com/realtime/current/gameplay/web-extensions/webhooks_

# Webhooks

Use webhooks with Photon to extend your application, allow players to rejoin games and to persist player and game data.

Photon webhooks are event-driven HTTP POST requests sent by the Photon servers to your specified service.

Each Photon webhook is defined by its own triggers, data and destination path.

We are phasing out "Persistent Games" support. New titles will not be able to rely on the \`IsPersistent\` setting, WebFlags and related features.

## Setup

To set up webhooks with Photon Cloud, follow the _Manage_ link from the application list on your [dashboard](https://dashboard.photonengine.com/).

Add new or change existing within the _Webhooks_ section.

From the webhooks configuration page you can select one of the predefined presets using the dropdown list.

The "Demo" preset should be used in development only.

Its purpose is to give you a quick preview of the different settings available and how to use them.

Configuration is done by defining key/value pairs of strings.

The maximum length allowed for each setting value is 1024 chars.

### Basic Settings

- **BaseUrl** (required)


The URL for the service hosting your hooks as well as any method that should be called via WebRPC.


It must not end with a forward slash.


Callbacks are received at the relative path URIs that are set up below.


Read [this](#querystring) if the value you want to configure will contain a query string.
- **CustomHttpHeaders**


JSON string of key/value pairs (string:string) that should be set as HTTP headers in any request made to the configured web service.


Read more [here](#httpheaders).

### Paths

Configure each path to receive events at the URIs each identifies on your host.

Any path that is left empty will not be hit, you will not receive any callbacks on the affected webhook.

- **PathCreate**


Called when a new room is created or when its state needs to be loaded from external service.
- **PathClose**


Called when a room is removed from Photon servers memory.


If _IsPersistent_ is set to `true` the room state is sent.
- **PathBeforeJoin**


Called when a player attempts to join a room.
- **PathJoin**


Called when a player joins a room when it is in Photon servers memory.
- **PathLeave**


Called when a player leaves the room.
- **PathEvent**


Called when the client raises an event in the room with the web flag `HttpForward` set.
- **PathGameProperties**


Called when the client sets a room or a player property with the web flag `HttpForward` set.

### Options

Finetune the behaviour of your webhooks configuration with these options.

An option is considered not configured when not set up or if its value is left empty.

The default value for each option is `false` and applies for each which is not set.

- **HasErrorInfo**


If set to `true`, clients will be notified with an `ErrorEvent` when a call to a hook fails.
- **IsPersistent**


If set to `true`, Photon sends the state before removing a room from memory.


This option is only valid when both _PathCreate_ and _PathClose_ are properly configured.


For more information read more about [GameCreate](#create) and [GameClose](#close) webhooks.
- **AsyncJoin**


This option is only valid when _IsPersistent_ is set to `true`.


By default this setting is set to `true` and join operations that fail to find room by name on server will be forwarded to web service as webhook.


To disable this behaviour set this to `false`.


Read more [here](#async).

### Query String Handling

Query string parameters can be included in the BaseUrl.

If you do use them you should know the following:

- Query string will be used as is. No URL encoding will be done. No duplicate keys check will be done.
- You can use URL tags in query string parameters.
- Do not add query string to Path settings.

**Example:**

- Configuration:

  - BaseUrl: `https://myawesomegame.com/chat/webhooks?clientver={AppVersion}&key=&keyA=valueA&keyA=valueB&keyB=valueB&=value`
  - PathCreate: `create`
  - PathClose: `close`
- Resulting URLs:

  - PathCreate: `https://myawesomegame.com/realtime/webhooks/create?clientver=1.0&key=&keyA=valueA&keyA=valueB&keyB=valueB&=value`
  - PathClose: `https://myawesomegame.com/realtime/webhooks/close?clientver=1.1&key=&keyA=valueA&keyA=valueB&keyB=valueB&=value`

### HTTP Headers Considerations

There are few things you need to consider when using custom HTTP headers:

- The value of CustomHttpHeaders configuration key needs to be a stringified JSON object that has properties with string values only.


The JSON object's properties' names will be used as HTTP request header field names and the properties' values will be used as their respective values.

**Example:**

  - CustomHttpHeaders value:


    ```
    `{'X-Secret': 'YWxhZGRpbjpvcGVuc2VzYW1l', 'X-Origin': 'Photon'}
    `
    ```

  - Webhooks HTTP request headers:


    ```
    `X-Secret: YWxhZGRpbjpvcGVuc2VzYW1l
    X-Origin: Photon
    `
    ```
- Custom HTTP headers field names are case sensitive.


- The supported formats for `Date` and `If-Modified-Since` can be found [here](https://docs.microsoft.com/en-us/dotnet/api/system.datetime.parse#StringToParse).
- `Content-Type` should not be used as as the webhooks plugin will set it to `application/json`.

- The following HTTP headers are restricted and will be ignored if set from the "CustomHttpHeaders" configuration value.

  - `Connection`
  - `Content-Length`
  - `Host`
  - `Range`
  - `Proxy-Connection`

### Turnkey Solutions

Find the latest available turnkey installations for Microsoft Azure and Heroku on [our github page](https://github.com/exitgames).

### URL Tags

You can optionally set one or more "dynamic variables".

Those variables will be replaced with their respective values in our backend before sending out any requests.

Empty space characters will be removed.

URL tags come in handy if you want to handle user segments of your application differently, each.

Photon supports the following URL tags:

- `{AppVersion}` will pass the application version as set by the client, e.g. "1.0".
- `{AppId}` will pass the ID of your application, e.g. "2afda618-e64f-4a85-b2a2-74e05fdf0b65".
- `{Region}` will pass the token for the cloud region the triggering client is connected to, e.g. "EU" or "USW".
- `{Cloud}` will pass the name of the cloud the triggering client is connected to, e.g. "public" or "enigmaticenterprise".

#### Examples of URL Tags use cases

1. `https://{Region}.mydomain.com/{AppId}?version={AppVersion}&cloud={Cloud}` to e.g. route each region to different hosts, version and cloud passed as query parameters.
2. `https://mydomain.com/{Cloud}/{Region}/{AppId}/{AppVersion}` passes all tags as well structured URI.

## Common Criteria

All Photon webhooks share the base URL to POST to, the ability to return an eventual error that could happen on the web server and baseline data that is sent with each.

### Common Arguments to All Webhooks

`AppId`:

AppId of your application as set from the game client and found in the [dashboard](https://dashboard.photonengine.com/).

`AppVersion`:

version of your application as set from the game client.

`Region`:

Has the region to which the game client is connected to and to which the room in question belongs to.

`GameId`:

ID or name of the room.

`Type`:

Type of the webhook, it's always "Event" in [PathEvent](#event) webhook, "Join" in [PathJoin](#join) webhook and can have multiple values in all the other webhooks.

### Common arguments to all Webhooks except PathClose

`ActorNr`:

number of the actor triggering the hook.

`UserId`:

UserId of the actor triggering the hook.

`NickName`:

name of the actor triggering the hook as set from the client SDK.

### Return Values

#### Expected Response

When sending the webhooks, Photon will in most cases expect only a response that has a JSON object with a `ResultCode` property set to 0.

The zero `ResultCode` is an acknowledgment for the reception of the webhook request from and for successful processing on the web server.

Any webhook call with a return object that has a `ResultCode` different than zero is considered as failed.

If `HasErrorInfo` is set to true in the configuration; an `ErrorInfo` event will be broadcast to clients still joined to the room.

Other than [PathCreate](#create), the server continues normal execution, despite this.

The event also contains `ParameterCode.Info`. [See the API doc for your SDK](https://www.photonengine.com/sdks) for additional info.

The only case where Photon expects more is when a player is trying to rejoin a Room that has been closed and already removed from memory.

For more details about this, please refer to the [PathCreate section](#create).

#### Best Practices

Always check the arguments needed in each webhook.

Return an error code with an optional, yet useful human readable message in case a required argument is missing.

The suggested format of the JSON return object of a webhook is `{ "ResultCode" : x, "Message" : "xxx" }`.

A very good way to handle return values is by implementing helper routines in the backends logic for your webhooks.

The following is just a non-exhaustive list of examples of webhooks return objects.

- Default return object for success

`{ "ResultCode" : 0 }`

`{ "ResultCode" : 0, "Message" : "OK" }`

- Error on protocol level, either on Photon side or on web server

`{ "ResultCode" : 1, "Message" : "Missing Webhook Argument: <argument name>." }`

- Application specific errors

`{ "ResultCode" : 2, "Message" : "Game with GameId=<gameId> already exists." }`

`{ "ResultCode" : 3, "Message" : "Could not load the State, Reason=<reason>." }`


## Paths in Details

### PathCreate

If a room is created as a result of a `JoinOrCreateRoom` call, "PathCreate" webhook will be triggered only if:

- "IsPersistent" is set to "true".

- "PathCreate" is configured.

- "PathClose" is configured.


Otherwise, the room is created and no webhook will be fired.

This webhook is triggered every time a room is created at server level.

In case the room is created using `OpCreateRoom`, the webhook will have its `Type` argument set to `Create`.

Most of the `RoomOptions` along with the `TypedLobby` used when creating the room will be sent as `CreateOptions`.

The value of the `ActorNr` will be 1 and the actor will auto join the room without firing [PathJoin](#join) webhooks.

If an actor tries to join or rejoin a room that was previously created but removed from the server's memory, the `Type` will be set to `Load`.

The web server should return the serialized `State` of the same room.

The backend logic will need to fetch the previously saved, serialized version of the room state.

In case the state is found, the web server should return a JSON object as `{ "State" : state, "ResultCode" : 0 }`.

If the state cannot be found for some reason, the web service should check the value of `CreateIfNotExists` to allow or disallow room creation with new state.

If `CreateIfNotExists` is true, you can return an empty room state (e.g. `{ "State" : "", "ResultCode" : 0 }`) to allow creation of a room the options provided by the client.

With its value being false, Photon should be informed that the web server failed to load the state by returning a `ResultCode` other than 0.

Consider adding a human readable error message with the cause.

If "IsPersistent" is set to `true` and both "PathCreate" and "PathClose" are configured, room creation, asynchronous join or rejoin will fail when:

- "PathCreate" is unreachable or returns HTTP error.

- "PathCreate" returns `ResultCode` other than `0`.

- Photon Server fails to parse or set received room state.


#### AsyncJoin Option

Photon applications support asynchronous operations like accepting an invitation to join a game hours after its creation to make playing with friends easier.

When both `IsPersistent` and `AsyncJoin` are enabled, any join operation where the room cannot be found in Photon servers memory, will trigger `PathCreate` (`Type="Load"`) webhook.

The goal is to try loading the room state from your web service, assuming it was previously created and saved.

`AsyncJoin` is enabled by default when `IsPersistent` is `true`.

To disable `AsyncJoin`, add its key if not present in webhooks settings, and set its value to `false`.

This table shows all client operations that might trigger the `PathCreate` webhook with `Type="Load"`, when `IsPersistent` is enabled (in case the room is not found in Photon servers memory).

| Client Operation | JoinMode | AsyncJoin = false | AsyncJoin = true | AsyncJoin = N/A |
| --- | --- | --- | --- | --- |
| `OpJoinRoom` | Default (Join) |  |  |  |
| `OpRejoinRoom` | RejoinOnly |  |  |  |
| `OpJoinOrCreateRoom` | CreateIfNotExists |  |  |  |

#### Specific Arguments

**PathCreate, Type "Create"**

`CreateOptions`:

The options used when creating the room.

It contains information from the `RoomOptions` class as set from the Client with details about the used `TypedLobby`.

All its properties can be retrieved later from the `State` as it will be copied there as is.

The following table has a comparison between `CreateOptions` and `RoomOptions`.

| Property or Field | CreateOptions | RoomOptions | Notes |
| --- | --- | --- | --- |
| IsVisible |  |  | Can be retrieved later from room State. |
| IsOpen |  |  | Can be retrieved later from room State. |
| MaxPlayers |  |  |  |
| PlayerTtl |  |  |  |
| EmptyRoomTtl |  |  |  |
| CheckUserOnJoin |  |  | This should be set to `true` if you have unique ID per player. |
| SuppressRoomEvents |  |  | If set to `true`, no room events are sent to the clients on join and leave.<br> Default is `false` and sent. |
| DeleteCacheOnLeave |  |  | This is called CleanupCacheOnLeave in RoomOptions. |
| LobbyType |  |  | The type of the room's lobby as set by game client in `TypedLobby.Type`.<br> [See the API doc for your SDK](https://www.photonengine.com/sdks) for additional info. |
| LobbyId |  |  | The name of the room's lobby as set by game client in `TypedLobby.Name`.<br> [See the API doc for your SDK](https://www.photonengine.com/sdks) for additional info. |
| CustomProperties |  |  | Contains only the initial custom properties of the room that should be public, i.e. visible to the lobby. |
| CustomRoomProperties |  |  | Contains all initial custom properties of the room. |
| CustomRoomPropertiesForLobby |  |  | Contains the keys of the custom properties of the room that should be public, i.e. visible to the lobby. |
| PublishUserId |  |  | Defines if the UserIDs of players get "published" in the room.<br> Can be retrieved later from room State.<br> Read more about it [here](/realtime/current/lobby-and-matchmaking/matchmaking-and-lobby#publish-userid). |

**PathCreate, Type "Load"**

`CreateIfNotExists`: A flag used to indicate if a new room should be created if its state could not be found from web service.

Properties of the room's `State` other than those found in `CreateOptions`:

- `ActorCounter`: The ActorNr of the last joined actor.
- `ActorList`: An array of information entries about each actor present in the room (active or inactive).


Each entry has the following properties:

  - `ActorNr`: Number of the actor in the room.
  - `UserId`: UserID of the actor.
  - `NickName`: NickName of the actor.
  - `IsActive`: Indicates if the actor is joined to the room or not. It is not sent in `State` argument of [PathClose](#close).
  - `DeactivationTime`: Timestamp of room leave event. Sent only in `State` argument of [PathClose](#close).
  - `Binary`\*: Base64 encoded actor properties.
  - `DEBUG\_BINARY`\*\*: readable form of the actor properties.
- `Slice`: the cache slice index.
- `Binary` \*: Base64 encoded room properties and cached events.
- `DebugInfo` \*\*: readable form of the binary data.

  - `DEBUG\_PROPERTIES\_18`\*\*: readable form of the room properties.
  - `DEBUG\_EVENTS\_19`\*\*: readable form of the cached events.
  - `DEBUG\_GROUPS\_20`\*\*: readable form of the interest groups.
- `ExpectedUsers`: Array of UserIDs of players expected to join the room. Read more about [Slot Reservation](/realtime/current/lobby-and-matchmaking/matchmaking-and-lobby#slot-reservation).

\\* : The `Binary` properties exist because of the nature of [the protocol used by Photon](/realtime/current/reference/binary-protocol).

\*\*: The `Debug` properties should be disabled for applications in production environment. Make sure to use these properties for debugging only.

**Sample Call**

JSON

```json
`{
    "ActorNr": 1,
    "AppVersion": "client-x.y.z",
    "AppId": "00000000-0000-0000-0000-000000000000",
    "CreateOptions": {
        "MaxPlayers": 4,
        "LobbyId": null,
        "LobbyType": 0,
        "CustomProperties": {
            "lobby3Key": "lobby3Val",
            "lobby4Key": "lobby4Val"
        },
        "EmptyRoomTTL": 0,
        "PlayerTTL": 2147483647,
        "CheckUserOnJoin": true,
        "DeleteCacheOnLeave": false,
        "SuppressRoomEvents": false
    },
    "GameId": "MyRoom",
    "Region": "EU",
    "Type": "Create",
    "UserId": "MyUserId1",
    "NickName": "MyPlayer1"
}
`
```

### PathBeforeJoin

This WebHook can be used to allow or deny access to rooms. If set, it is fired before an Actor is created for the user, so the ActorNr argument is unavailable. The UserId is defined server side for apps with authentication.

The `Type` argument is set to `BeforeJoin`.

Note: This path is missing in older WebHook configurations but can be added by us on demand.

Use a `ResultCode==0` to allow the join. Every other code denies the join.

#### Specific Arguments

`PathBeforeJoin` has no extra arguments.

**Sample Call**

JSON

```json
`{
    "Type": "BeforeJoin",
    "AppVersion": "client-x.y.z",
    "AppId": "00000000-0000-0000-0000-000000000000",
    "GameId": "MyRoom",
    "Region": "EU",
    "UserId": "MyUserId0",
    "NickName": "MyPlayer0",
    "AuthCookie": "data"
}
`
```

### PathJoin

If an actor joins or rejoins a room that was not removed from memory yet, this webhook will be fired.

The `Type` argument is set to `Join`.

#### Specific Arguments

`PathJoin` has no extra arguments.

**Sample Call**

JSON

```json
`{
    "ActorNr": 2,
    "AppVersion": "client-x.y.z",
    "AppId": "00000000-0000-0000-0000-000000000000",
    "GameId": "MyRoom",
    "Region": "EU",
    "Type": "Join",
    "UserId": "MyUserId0",
    "NickName": "MyPlayer0"
}
`
```

### PathGameProperties

This webhook is fired every time the user sets the custom properties of either the room or the player from the client side if the right overload method is used with the `HttpForward` web flag set.

The `Type` argument sent with the webhook will be set accordingly to `Game` or `Actor`.

#### Specific Arguments

`Properties`: a set of updated properties as sent from client SDK.

`State`: a serialized snapshot of the room's full state.

It's sent only if `SendState` webflag is set when calling `OpSetCustomProperties` and "IsPersistent" setting is set to `true`.

`AuthCookie`: an encrypted object invisible to client, optionally returned by web service upon successful custom authentication.

It's sent only if `SendAuthCookie` webflag is set when calling `OpSetCustomProperties`.

**PathGameProperties, Type="Actor"**

`TargetActor`: the number of the actor whose properties are updated.

**Sample Call**

JSON

```json
`{
    "ActorNr": 1,
    "AppVersion": "client-x.y.z",
    "AppId": "00000000-0000-0000-0000-000000000000",
    "Properties": {
        "turn": 1,
        "lobby3Key": "test1a",
        "lobby4Key": "test1b"
    },
    "GameId": "MyRoom",
    "Region": "EU",
    "State": {
        "ActorCounter": 2,
        "ActorList": [
            {
                "ActorNr": 1,
                "UserId": "MyUserId1",
                "NickName": "MyPlayer1",
                "IsActive": true,
                "Binary": "RGIAAAEBRAAAAAJzAAlz...",
                "DEBUG_BINARY": {
                    "1": {
                        "255": "MyPlayer1",
                        "player_id": "12345"
                    }
                }
            },
            {
                "ActorNr": 2,
                "UserId": "MyUserId0",
                "NickName": "MyPlayer0",
                "IsActive": true,
                "Binary": "RGIAAEBRAAAAAFi/3M15...",
                "DEBUG_BINARY": {
                    "1": {
                        "255": "MyPlayer0"
                    }
                }
            }
        ],
        "Binary": {
            "18": "RAAAAAdzAAhwcm9wMUtl..."
        },
        "CheckUserOnJoin": true,
        "CustomProperties": {
            "lobby4Key": "test1b",
            "lobby3Key": "test1a"
        },
        "DeleteCacheOnLeave": false,
        "EmptyRoomTTL": 0,
        "IsOpen": true,
        "IsVisible": true,
        "LobbyType": 0,
        "LobbyProperties": [
            "lobby3Key",
            "lobby4Key"
        ],
        "MaxPlayers": 4,
        "PlayerTTL": 2147483647,
        "SuppressRoomEvents": false,
        "Slice": 0,
        "DebugInfo": {
            "DEBUG_PROPERTIES_18": {
                "250": [
                    "lobby3Key",
                    "lobby4Key"
                ],
                "prop1Key": "prop1Val",
                "prop2Key": "prop2Val",
                "lobby4Key": "test1b",
                "lobby3Key": "test1a",
                "map_name": "mymap",
                "turn": 1
            }
        }
    },
    "Type": "Game",
    "UserId": "MyUserId1",
    "NickName": "MyPlayer1"
}
`
```

### PathEvent

Fired every time the user raises a custom event from the client side if the right overload method is used with the `HttpForward` web flag set.

The custom event code and the event data will be sent with the webhook.

To prevent abuse, we are phasing out the WebFlags and limit PathEvent calls to three per second.

#### Specific Arguments

`EvCode`: custom event code.

`Data`: custom event data as sent from client SDK.

`State`: a serialized snapshot of the room's full state.

It's sent only if `SendState` webflag is set when calling `OpRaiseEvent` and "IsPersistent" setting is set to `true`.

`AuthCookie`: an encrypted object invisible to client, optionally returned by web service upon successful custom authentication.

It's sent only if `SendAuthCookie` webflag is set when calling `OpRaiseEvent`.

**Sample Call**

JSON

```json
`{
    "ActorNr": 3,
    "AppVersion": "client-x.y.z",
    "AppId": "00000000-0000-0000-0000-000000000000",
    "Data": "data",
    "GameId": "MyRoom",
    "Region": "EU",
    "State": {
        "ActorCounter": 3,
        "ActorList": [
            {
                "ActorNr": 1,
                "UserId": "MyUserId1",
                "NickName": "MyPlayer1",
                "Binary": "RGIAAAEBRAAAAAJzAAlw...",
                "DEBUG_BINARY": {
                    "1": {
                        "255": "MyPlayer1",
                        "player_id": "12345"
                    }
                }
            },
            {
                "ActorNr": 3,
                "UserId": "MyUserId0",
                "NickName": "MyPlayer0",
                "IsActive": true,
                "Binary": "RGIAAAEBRAAAAAFi/3MAC...",
                "DEBUG_BINARY": {
                    "1": {
                        "255": "MyPlayer0"
                    }
                }
            }
        ],
        "Binary": {
            "18": "RAAAAAdzAAhwcm9wMUtl...",
            "19": "RGl6AAEAAAAAAAN6AANp..."
        },
        "CheckUserOnJoin": true,
        "CustomProperties": {
            "lobby4Key": "test1b",
            "lobby3Key": "test1a"
        },
        "DeleteCacheOnLeave": false,
        "EmptyRoomTTL": 0,
        "IsOpen": true,
        "IsVisible": true,
        "LobbyType": 0,
        "LobbyProperties": [
            "lobby3Key",
            "lobby4Key"
        ],
        "MaxPlayers": 4,
        "PlayerTTL": 2147483647,
        "SuppressRoomEvents": false,
        "Slice": 0,
        "DebugInfo": {
            "DEBUG_PROPERTIES_18": {
                "250": [
                    "lobby3Key",
                    "lobby4Key"
                ],
                "prop1Key": "prop1Val",
                "prop2Key": "prop2Val",
                "lobby4Key": "test1b",
                "lobby3Key": "test1a",
                "map_name": "mymap",
                "turn": 1
            },
            "DEBUG_EVENTS_19": {
                "0": [
                    [
                        3,
                        0,
                        "data"
                    ],
                    [
                        3,
                        0,
                        "data"
                    ],
                    [
                        3,
                        0,
                        "data"
                    ]
                ]
            }
        }
    },
    "Type": "Event",
    "UserId": "MyUserId0",
    "NickName": "MyPlayer0",
    "EvCode": 0
}
`
```

### PathLeave

This hook is triggered whenever an actor is disconnected from Photon game servers.

A disconnect could happen for several reasons.

The webhook itself tells you about that reason in a human readable form within its `Type` and in a coded way, see `Reason`.

#### Specific Arguments

`Type`: readable form of the reason that could be one of the following values:

- `ClientDisconnect`: Indicates that the client called `OpLeaveRoom()` or `Disconnect()`.
- `ClientTimeoutDisconnect`: Indicates that client has _timed-out_ server.


This is valid only when using UDP/ENET.
- `ManagedDisconnect`: Indicates client is too slow to handle data sent.
- `ServerDisconnect`: Indicates low level protocol error which can be caused by data corruption.
- `TimeoutDisconnect`: Indicates that the server has _timed-out_ client.
- `LeaveRequest`: Indicates that the client explicitly abandoned the room by calling `OpLeaveRoom()``OpLeaveRoom(false)`.
- `PlayerTtlTimedOut`: Indicates that the inactive actor _timed-out_, meaning the `PlayerTtL` of the room expired for that actor.


See the API doc for your SDK for additional info.
- `PeerLastTouchTimedOut`: Indicates a very unusual scenario where the actor did not send _anything_ to Photon Servers for 5 minutes.


Normally peers _timeout_ long before that but Photon does a check for every connected peer's _timestamp_ of the last exchange with the servers (called `LastTouch`) every 5 minutes.
- `PluginRequest`: Indicates that the actor was removed from ActorList by plugin.
- `PluginFailedJoin`: Indicates an internal error in the plugin implementation.

`IsInactive`: Refers to the state of the actor before leaving the Room. If set to `true` then the actor can rejoin the game. If set to `false`, then the actor left for good and is removed from ActorList and can't rejoin the game.

`Reason`: Code of the Reason

The tables below match each type to its code and explain how the webhook is produced and how the player presence in the room is affected.

##### `RoomOptions.PlayerTTL` set to 0 upon Room Creation

Actors of a room where `RoomOptions.PlayerTTL == 0` can never be in the inactive state and will be removed from the ActorList as soon as they leave.

`IsPersistent` option will be ignored as by design room state should be saved only if the room contains at least one inactive actor.

Any event that triggers the [PathLeave](#leave) webhook will result in removing the respective actor from the room.

Do not try and check if such actor is active or not ;)

`OpLeaveRoom(true)` will no longer be of use and `PlayerTtlTimedOut` type of [PathLeave](#leave) can not happen in this situation.

| Type | Reason | Trigger | Inactive |
| --- | --- | --- | --- |

| ClientDisconnect | 0 | a call to `Disconnect()` | false |
| ClientTimeoutDisconnect | 1 | called by Photon server | false |
| ManagedDisconnect | 2 | called by Photon server | false |
| ServerDisconnect | 3 | called by Photon server | false |
| TimeoutDisconnect | 4 | called by Photon server | false |
| LeaveRequest | 101 | a call to `OpLeaveRoom()` | false |
| PlayerTtlTimedOut | 102 | N/A | N/A |
| PeerLastTouchTimedout | 103 | called by Photon server | false |
| PluginRequest | 104 | called by plugin | false |
| PluginFailedJoin | 105 | called by Photon server | false |

##### `RoomOptions.PlayerTTL != 0` upon Room Creation

| Type | Reason | Trigger | IsInactive |
| --- | --- | --- | --- |

| ClientDisconnect | 0 | a call to `Disconnect()` | true |
| ClientTimeoutDisconnect | 1 | called by Photon server | true |
| ManagedDisconnect | 2 | called by Photon server | true |
| ServerDisconnect | 3 | called by Photon server | true |
| TimeoutDisconnect | 4 | called by Photon server | true |
| LeaveRequest | 101 | a call to `OpLeaveRoom()` | false |
| PlayerTtlTimedOut | 102 | called by Photon server | false |
| PeerLastTouchTimedout | 103 | called by Photon server | false |
| PluginRequest | 104 | called by plugin | false |
| PluginFailedJoin | 105 | called by Photon server | false |

**Sample Call**

JSON

```json
`{
    "ActorNr": 1,
    "AppVersion": "client-x.y.z",
    "AppId": "00000000-0000-0000-0000-000000000000",
    "GameId": "MyRoom",
    "IsInactive": true,
    "Reason": "0",
    "Region": "EU",
    "Type": "ClientDisconnect",
    "UserId": "MyUserId1",
    "NickName": "MyPlayer1"
}
`
```

### PathClose

This hook is fired just before removing a room instance from memory server side.

This happens only when the `EmptyRoomTTL` has expired.

`EmptyRoomTTL` is set when creating the room and it is a _timer_ that starts when the room becomes empty.

A room is considered empty when the last active actor in the room leaves it.

In case `IsPersistent` is set to `true`, Photon will send the `State` which is a serialized _snapshot_ of the room that contains its properties and information about the actors and cached events.

In this case, the `Type` of the webhook will have the `Save` value.

In case `IsPersistent` is set to `false`, which is the default state, the `Type` is `Close` and the `State` is not sent and the room is lost for ever.

#### Specific Arguments

`ActorCount`: the number of inactive actors. If 0 then `Type` should be "Close".

**PathClose, Type: "Save"**

`State`: a serialized snapshot of the room's full state

**Sample Call**

JSON

```json
`{
    "ActorCount": 2,
    "AppVersion": "client-x.y.z",
    "AppId": "00000000-0000-0000-0000-000000000000",
    "GameId": "MyRoom",
    "Region": "EU",
    "State": {
        "ActorCounter": 3,
        "ActorList": [
            {
                "ActorNr": 1,
                "UserId": "MyUserId1",
                "NickName": "MyPlayer1",
                "Binary": "RGIAAAEBRAAAAAJzAAlw...",
                "DEBUG_BINARY": {
                    "1": {
                        "255": "MyPlayer1",
                        "player_id": "12345"
                    }
                }
            },
            {
                "ActorNr": 3,
                "UserId": "MyUserId0",
                "NickName": "MyPlayer0",
                "Binary": "RGIAAAEBRAAAAAFi/3MA...",
                "DEBUG_BINARY": {
                    "1": {
                        "255": "MyPlayer0"
                    }
                }
            }
        ],
        "Binary": {
            "18": "RAAAAAdzAAhwcm9wMUtl...",
            "19": "RGl6AAEAAAAAAAN6AANp..."
        },
        "CheckUserOnJoin": true,
        "CustomProperties": {
            "lobby4Key": "test1b",
            "lobby3Key": "test1a"
        },
        "DeleteCacheOnLeave": false,
        "EmptyRoomTTL": 0,
        "IsOpen": true,
        "IsVisible": true,
        "LobbyType": 0,
        "LobbyProperties": [
            "lobby3Key",
            "lobby4Key"
        ],
        "MaxPlayers": 4,
        "PlayerTTL": 2147483647,
        "SuppressRoomEvents": false,
        "Slice": 0,
        "DebugInfo": {
            "DEBUG_PROPERTIES_18": {
                "250": [
                    "lobby3Key",
                    "lobby4Key"
                ],
                "prop1Key": "prop1Val",
                "prop2Key": "prop2Val",
                "lobby4Key": "test1b",
                "lobby3Key": "test1a",
                "map_name": "mymap",
                "turn": 1
            },
            "DEBUG_EVENTS_19": {
                "0": [
                    [
                        3,
                        0,
                        "data"
                    ],
                    [
                        3,
                        0,
                        "data"
                    ],
                    [
                        3,
                        0,
                        "data"
                    ]
                ]
            }
        }
    },
    "Type": "Save"
}
`
```

## Stripping Room State

If the size of the room state is a concern for you, you can do few tricks to shrink it and not lose any data.

When saving the room state in your web service you can reduce its size by:

1. Removing some fields that are made for read only debug purpose and will not be used during reconstruction of the room state when loading the room.


These fields are: `CustomProperties`, `DebugInfo` and every `DEBUG\_BINARY` and `NickName` of each actor inside `ActorList`.


These fields are safe to remove as they are ignored when deserializing the room state.

**Here is a full room state returned in PathCreate, Type="Load" webhook:**

JSON
```json
`{
       "ResultCode": 0,
       "Message": "Room State successfully loaded",
       "State": {
           "ActorCounter": 3,
           "ActorList": [
               {
                   "ActorNr": 1,
                   "UserId": "MyUserId1",
                   "NickName": "MyPlayer1",
                   "Binary": "RGIAAAEBRAAAAAJzAAlw...",
                   "DEBUG_BINARY": {
                       "1": {
                           "255": "MyPlayer1",
                           "player_id": "12345"
                       }
                   }
               },
               {
                   "ActorNr": 3,
                   "UserId": "MyUserId0",
                   "NickName": "MyPlayer0",
                   "Binary": "RGIAAAEBRAAAAAFi/3MA...",
                   "DEBUG_BINARY": {
                       "1": {
                           "255": "MyPlayer0"
                       }
                   }
               }
           ],
           "Binary": {
               "18": "RAAAAAdzAAhwcm9wMUtl...",
               "19": "RGl6AAEAAAAAAAN6AANp..."
           },
           "CheckUserOnJoin": true,
           "CustomProperties": {
               "lobby4Key": "test1b",
               "lobby3Key": "test1a"
           },
           "DeleteCacheOnLeave": false,
           "EmptyRoomTTL": 0,
           "IsOpen": true,
           "IsVisible": true,
           "LobbyType": 0,
           "LobbyProperties": [
               "lobby3Key",
               "lobby4Key"
           ],
           "MaxPlayers": 4,
           "PlayerTTL": 2147483647,
           "SuppressRoomEvents": false,
           "Slice": 0,
           "DebugInfo": {
               "DEBUG_PROPERTIES_18": {
                   "250": [
                       "lobby3Key",
                       "lobby4Key"
                   ],
                   "prop1Key": "prop1Val",
                   "prop2Key": "prop2Val",
                   "lobby4Key": "test1b",
                   "lobby3Key": "test1a",
                   "map_name": "mymap",
                   "turn": 1
               },
               "DEBUG_EVENTS_19": {
                   "0": [
                       [
                           3,
                           0,
                           "data"
                       ],
                       [
                           3,
                           0,
                           "data"
                       ],
                       [
                           3,
                           0,
                           "data"
                       ]
                   ]
               }
           }
       }
}
`
```


**Here is a stripped version of the room state that could be returned in the same webhook:**

JSON
```json
`{
       "ResultCode": 0,
       "Message": "Room State successfully loaded",
       "State": {
           "ActorCounter": 3,
           "ActorList": [
               {
                   "ActorNr": 1,
                   "UserId": "MyUserId1",
                   "Binary": "RGIAAAEBRAAAAAJzAAlw..."
               },
               {
                   "ActorNr": 3,
                   "UserId": "MyUserId0",
                   "Binary": "RGIAAAEBRAAAAAFi/3MA..."
               }
           ],
           "Binary": {
               "18": "RAAAAAdzAAhwcm9wMUtl...",
               "19": "RGl6AAEAAAAAAAN6AANp..."
           },
           "CheckUserOnJoin": true,
           "DeleteCacheOnLeave": false,
           "EmptyRoomTTL": 0,
           "IsOpen": true,
           "IsVisible": true,
           "LobbyType": 0,
           "LobbyProperties": [
               "lobby3Key",
               "lobby4Key"
           ],
           "MaxPlayers": 4,
           "PlayerTTL": 2147483647,
           "SuppressRoomEvents": false,
           "Slice": 0
       }
}
`
```

2. Removing some fields that are considered in your application logic as constants and should not be changed during room lifetime or can be injected by code before returning the state to Photon Servers when loading the room.


These fields should be added back to their original value before returning the state to Photon Servers. Otherwise the state may become broken.

Those fields could be any of these:

   - `ActorCounter`
   - `CheckUserOnJoin`
   - `DeleteCacheOnLeave`
   - `EmptyRoomTTL`
   - `IsOpen`
   - `IsVisible`
   - `MaxPlayers`
   - `LobbyId`
   - `LobbyType`
   - `LobbyProperties`
   - `PlayerTTL`
   - `SuppressRoomEvents`
   - `Slice`

## Securing Webhooks

Other than using HTTPS and newer TLS version, you can enforce webhooks security using [custom HTTP request headers](#httpheaders) or [query string parameters](#querystring).

The idea is that you set up one or more "secrets" ("token", "key", etc.) from Photon dashboard that can help you make sure that the incoming HTTP requests are originating from Photon servers.

Back to top

- [Setup](#setup)

  - [Basic Settings](#basic-settings)
  - [Paths](#paths)
  - [Options](#options)
  - [Query String Handling](#query-string-handling)
  - [HTTP Headers Considerations](#http-headers-considerations)
  - [Turnkey Solutions](#turnkey-solutions)
  - [URL Tags](#url-tags)

- [Common Criteria](#common-criteria)

  - [Common Arguments to All Webhooks](#common-arguments-to-all-webhooks)
  - [Common arguments to all Webhooks except PathClose](#common-arguments-to-all-webhooks-except-pathclose)
  - [Return Values](#return-values)

- [Paths in Details](#paths-in-details)

  - [PathCreate](#pathcreate)
  - [PathBeforeJoin](#pathbeforejoin)
  - [PathJoin](#pathjoin)
  - [PathGameProperties](#pathgameproperties)
  - [PathEvent](#pathevent)
  - [PathLeave](#pathleave)
  - [PathClose](#pathclose)

- [Stripping Room State](#stripping-room-state)
- [Securing Webhooks](#securing-webhooks)