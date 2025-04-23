# matchmaking-and-lobby

_Source: https://doc.photonengine.com/realtime/current/lobby-and-matchmaking/matchmaking-and-lobby_

# Matchmaking Guide

Getting into a room to play with (or against) someone is very easy with Photon.

There are basically three approaches:

Either tell the server to find a matching room, follow a friend into her room, or fetch a list of rooms to let the user pick one.

All three variants are supported by Photon and you can even roll your own.

We think, for most games it's best to use a quick and simple matchmaking, so we suggest to use Random Matchmaking and maybe filters for skills, levels and such.

## Matchmaking Checklist

If you run into issues matching players, here is a quick checklist:

- Verify clients are connected to the same `Region`. Note: [Best Region results are not deterministic.](/realtime/current/connection-and-authentication/regions#best-region-considerations)
- Use the same `AppId` and [AppVersion](/realtime/current/reference/glossary#appversion) in all clients and builds.
- Verify that clients use unique userIDs. Photon avoids matching the same userID into rooms.
- Before trying to join a room by name, make sure that this room is already created. Alternatively use `JoinOrCreateRoom`.
- For random matchmaking, use the same lobby (name and type) for creating and joining rooms.
- For random matchmaking with a properties filter, create rooms with fitting `Custom Properties in Lobby` string-array.
- For random matchmaking with SQL filters, use the SQL lobby type and create rooms with the fitting pre-defined keys (C0 and up) as `Custom Properties in Lobby` string-array.

## Quick Match

Most players nowadays just want to jump directly into matches.

They want to dive into the game right away.

That is why most games offer Quick Match as the first mode.

The suggested workflow described here gets players into rooms without asking them to pick one (randomly) from a long list of rooms.

If you just want to get players into a room quickly, do the following:

### JoinRandomOrCreateRoom

Simply call `JoinRandomOrCreateRoom` (exact function or method name may vary depending on your client SDK).

If a room is found it will be joined, otherwise a new room will be created.

C#

```csharp
using Photon.Realtime;
using System.Collections.Generic;
public class QuickMatchExample : IMatchmakingCallbacks
{
    private LoadBalancingClient loadBalancingClient;
    private void QuickMatch()
    {
        loadBalancingClient.OpJoinRandomOrCreateRoom(null, null);;
    }
    // do not forget to register callbacks via loadBalancingClient.AddCallbackTarget
    // also deregister via loadBalancingClient.RemoveCallbackTarget
    #region IMatchmakingCallbacks
    void IMatchmakingCallbacks.OnJoinedRoom()
    {
        // joined a room successfully
    }
    // [..] Other callbacks implementations are stripped out for brevity, they are empty in this case as not used.
    #endif
}

```

### JoinRandomRoom or CreateRoom

- Try `JoinRandomRoom` (exact function or method name may vary depending on your client SDK).

  - In best case, that's it.


    Your client will join a room successfully.
  - In worst case, no room is existing or no room can be joined (closed, invisible or full).
- If this doesn't find a room instantly, create one.
  - If you never show room names (and why should you), don't make up a name.


    Let the server do this.


    Set null or empty string as "room name" when creating the room.


    The room gets a GUID which is unique.
  - Apply a value for "max players".


    This way, the server eventually stops adding players when the room is full.
- If your client is alone in the room (players count == 1):


Wait.


Show a screen you're waiting for opponents.
- When enough players are in the room, you might "start" the game.


To keep new players out, "close" the room.


The server stops filling up the room, even if it's not full yet.

  - Note: When you close the room, there is a short time where players maybe are already on the way in.


    Don't be surprised if someone joins even 'after' a client closes the room.

C#

```csharp
using Photon.Realtime;
using System.Collections.Generic;
public class QuickMatchExample : IMatchmakingCallbacks
{
    [SerializeField]
    private maxPlayers = 4;
    private LoadBalancingClient loadBalancingClient;
    private void CreateRoom()
    {
        RoomOptions roomOptions = new RoomOptions();
        roomOptions.MaxPlayers = maxPlayers;
        EnterRoomParams enterRoomParams = new EnterRoomParams();
        enterRoomParams.RoomOptions = roomOptions;
        loadBalancingClient.OpCreateRoom(enterRoomParams);
    }
    private void QuickMatch()
    {
        loadBalancingClient.OpJoinRandomRoom();
    }
    // do not forget to register callbacks via loadBalancingClient.AddCallbackTarget
    // also deregister via loadBalancingClient.RemoveCallbackTarget
    #region IMatchmakingCallbacks
    void IMatchmakingCallbacks.OnJoinRandomFailed(short returnCode, string message)
    {
        CreateRoom();
    }
    void IMatchmakingCallbacks.OnJoinedRoom()
    {
        // joined a room successfully
    }
    // [..] Other callbacks implementations are stripped out for brevity, they are empty in this case as not used.
    #endif
}

```

Using this workflow, joining a game is a breeze for your players.

## Random Matchmaking

Sometimes players want more than a Quick Match, they want to play a certain map or mode (e.g. two versus two, etc.).

You can set arbitrary "Custom Room Properties" and use them as filter in `JoinRandomRoom` request (exact function or method name may vary depending on your client SDK).

### Exposing Some Properties In The Lobby

Custom room properties are synced to all players in the room and can be useful to keep track of the current map, game mode, difficulty, turn, round, start-time, etc.

They are handled as Hashtable with string keys.

By default, to keep things lean, these properties are accessible only inside rooms and are not sent to the Master Server (where lobbies exist).

You can choose some custom room properties to expose in the lobby.

Those properties will be used as filters for random matchmaking and they will be visible in the lobby (sent as part of the room info in the rooms list, only lobbies of default type send rooms lists).

**Example**:

To make "map" and "gm" available for for matchmaking, you can set a list of "room properties visible in the lobby" when you create a room.

Tip: Brief names are better, so use "gm" instead of "GameMode".

C#

```csharp
using Photon.Realtime;
using System.Collections.Generic;
using Hashtable = ExitGames.Client.Photon.Hashtable;
public class CreateRoomWithLobbyPropertiesExample : IMatchmakingCallbacks
{
    public const string MAP_PROP_KEY = "map";
    public const string GAME_MODE_PROP_KEY = "gm";
    public const string AI_PROP_KEY = "ai";
    private LoadBalancingClient loadBalancingClient;
    private void CreateRoom()
    {
        RoomOptions roomOptions = new RoomOptions();
        roomOptions.CustomRoomPropertiesForLobby = { MAP_PROP_KEY, GAME_MODE_PROP_KEY, AI_PROP_KEY };
        roomOptions.CustomRoomProperties = new Hashtable { { MAP_PROP_KEY, 1 }, { GAME_MODE_PROP_KEY, 0 } };
        EnterRoomParams enterRoomParams = new EnterRoomParams();
        enterRoomParams.RoomOptions = roomOptions;
        loadBalancingClient.OpCreateRoom(enterRoomParams);
    }
    // do not forget to register callbacks via loadBalancingClient.AddCallbackTarget
    // also deregister via loadBalancingClient.RemoveCallbackTarget
    #region IMatchmakingCallbacks
    void IMatchmakingCallbacks.OnCreateRoomFailed(short returnCode, string message)
    {
       // log error message and code
    }
    void IMatchmakingCallbacks.OnCreatedRoom()
    {
    }
    void IMatchmakingCallbacks.OnJoinedRoom()
    {
        // joined a room successfully, OpCreateRoom leads here on success
    }
    // [..] Other callbacks implementations are stripped out for brevity, they are empty in this case as not used.
    #endif
}

```

Note that "ai" has no value initially.

It won't show up in the lobby until it's set in the room (in C# SDKs, this is done via `Room.SetCustomProperties`).

When you change the values for "map" or "gm" or "ai", they will be updated in the lobby with a short delay, too.

Later (post room creation), you could also change the room properties keys visible to the lobby (add or remove) (in C# SDKs, this is done via `Room.PropertiesListedInLobby`)

Tip: Keep the list of lobby properties short to make sure your clients performance doesn't suffer from loading them either when joining the room or joining the lobby (only lobbies of default type send rooms lists).

Again: You don't have to join the lobby (and get the awfully long room list) to make use of this.

When you set some for the lobby, they become available as filter, too.

### Filtering Room Properties in Join Random

This section excludes SQL Lobby type.
Read more about SQL matchmaking [here](#sql_lobby_type).

When trying to find a random room, you could optionally choose the expected room properties or the expected max players.

These work as filters when the server selects a "fitting" room for you.

**Example**:

C#

```csharp
using Photon.Realtime;
using System.Collections.Generic;
using Hashtable = ExitGames.Client.Photon.Hashtable;
public class RandomMatchmakingExample : IMatchmakingCallbacks
{
    public const string MAP_PROP_KEY = "map";
    private LoadBalancingClient loadBalancingClient;
    public void JoinRandomRoom(byte mapCode, byte expectedMaxPlayers)
    {
        Hashtable expectedCustomRoomProperties = new Hashtable { { MAP_PROP_KEY, mapCode } };
        OpJoinRandomRoomParams opJoinRandomRoomParams = new OpJoinRandomRoomParams();
        opJoinRandomRoomParams.ExpectedMaxPlayers = expectedMaxPlayers;
        opJoinRandomRoomParams.ExpectedCustomRoomProperties = expectedCustomRoomProperties:
        loadBalancingClient.OpJoinRandomRoom(opJoinRandomRoomParams);
    }
    // do not forget to register callbacks via loadBalancingClient.AddCallbackTarget
    // also deregister via loadBalancingClient.RemoveCallbackTarget
    #region IMatchmakingCallbacks
    void IMatchmakingCallbacks.OnJoinRandomFailed(short returnCode, string message)
    {
        // log error code and message
        // here usually you create a new room
    }
    void IMatchmakingCallbacks.OnJoinedRoom()
    {
        // joined a room successfully, OpJoinRandomRoom leads here on success
    }
    // [..] Other callbacks implementations are stripped out for brevity, they are empty in this case as not used.
    #endif
}

```

If you pass more filter properties, chances are lower that a room matches them.

Better limit the options.

Make sure you always filter using properties visible to the lobby as shown [here](#exposing-some-properties-in-the-lobby).

## Play with Your Friends

If your users communicate with friends (e.g. with Photon Chat), they can easily make up a room name and everyone just uses `JoinOrCreateRoom` (exact function or method name may vary depending on your client SDK) to get into that room.

**Example**:

A unique room name could be composed (e.g.) as: "friendName1 + friendName2 + randomInteger".

To avoid anyone else joining, create the room invisible like so:

C#

```csharp
using Photon.Realtime;
using System.Collections.Generic;
public class PrivateRoomExample : IMatchmakingCallbacks
{
    private LoadBalancingClient loadBalancingClient;
    public void JoinOrCreatePrivateRoom(string nameEveryFriendKnows)
    {
        RoomOptions roomOptions = new RoomOptions();
        roomOptions.IsVisible = false;
        EnterRoomParams enterRoomParams = new EnterRoomParams();
        enterRoomParams.RoomName = nameEveryFriendKnows;
        enterRoomParams.RoomOptions = roomOptions;
        loadBalancingClient.OpJoinOrCreateRoom(enterRoomParams);
    }
    // do not forget to register callbacks via loadBalancingClient.AddCallbackTarget
    // also deregister via loadBalancingClient.RemoveCallbackTarget
    #region IMatchmakingCallbacks
    void IMatchmakingCallbacks.OnJoinRoomFailed(short returnCode, string message)
    {
      // log error code and message
    }
    void IMatchmakingCallbacks.OnJoinedRoom()
    {
        // joined a room successfully, OpJoinOrCreateRoom leads here on success
    }
    // [..] Other callbacks implementations are stripped out for brevity, they are empty in this case as not used.
    #endif
}

```

You can also look for your friends using `FindFriends` (exact function or method name may vary depending on your client SDK) if you use unique UserIDs which you should.

### Publishing UserIDs in a Room

Photon uses a UserID in various places.

For example, you can find friends only with a suitable UserID per player.

We added an option to Photon, which makes the UserID of players known per room.

In C# SDKs, set `RoomOptions.PublishUserId` to `true`, when you create a room.

The server will then provide the UserID and you can access it on the client.

In C# SDKs, it's done via `Player.UserId`.

Notes:

- UserIDs are broadcasted, with player properties, in the Photon join event.
- The UserID for a client can be set in three ways:
1. Client sets UserID before connecting.
2. Returned by an external web service using Custom Authentication.


     It will override the value sent by the client.
3. Photon will make up UserIDs (GUID) for users that don't explicitly set theirs.
- Generally, UserIDs, are not intended to be displayed.

### Matchmaking Slot Reservation

Sometimes, a player joins a room, knowing that a friend should join as well.

With Slot Reservation, Photon can block a slot for specific users and take that into account for matchmaking.

To reserve slots there is an `expectedUsers` parameter (exact parameter or argument name may vary depending on your client SDK) in the methods that get you in a room (`JoinRoom`, `JoinOrCreateRoom`, `JoinRandomRoom` and `CreateRoom`. Exact functions or methods names may vary depending on your client SDK).

{% if PUN %}

C#

```csharp
// create room example
PhotonNetwork.CreateRoom(roomName, roomOptions, typedLobby, expectedUsers);
// join room example
PhotonNetwork.JoinRoom(roomName, expectedUsers);
// join or create room example
PhotonNetwork.JoinOrCreateRoom(roomName, roomOptions, typedLobby, expectedUsers);
// join random room example
PhotonNetwork.JoinRandomRoom(expectedProperties, maxPlayers, expectedUsers, matchmakingType, typedLobby, sqlLobbyFilter, expectedUsers);

```

C#

```csharp
EnterRoomParams enterRoomParams = new EnterRoomParams();
enterRoomParams.ExpectedUsers = expectedUsers;
// create room example
loadBalancingClient.OpCreateRoom(enterRoomParams);
// join room example
loadBalancingClient.OpJoinRoom(enterRoomParams);
// join or create room example
loadBalancingClient.OpJoinOrCreateRoom(enterRoomParams);
// join random room example
OpJoinRandomRoomParams opJoinRandomRoomParams = new OpJoinRandomRoomParams();
opJoinRandomRoomParams.ExpectedUsers = expectedUsers;
loadBalancingClient.OpJoinRandomRoom(opJoinRandomRoomParams);

```

When you know someone should join, pass an array of UserIDs.

For `JoinRandomRoom`, the server will attempt to find a room with enough slots for you and your expected players (plus all active and expected players already in the room).

The server will update clients in a room with the current `expectedUsers`, should they change.

You can update the list of expected users inside a room (add or remove one or more users), this is done via a well known room property.

(In C# SDKs, you can get and set `Room.ExpectedUsers`).

To support Slot Reservation, you need to enable publishing UserIDs inside rooms.

#### Example Use Case: Teams Matchmaking

You can use this to support teams in matchmaking.

The leader of a team does the actual matchmaking.

He/She can join a room and reserve slots for all members:

Try to find a random room:

C#

```csharp
OpJoinRandomRoomParams opJoinRandomRoomParams = new OpJoinRandomRoomParams();
opJoinRandomRoomParams.ExpectedUsers = teamMembersUserIds;
loadBalancingClient.OpJoinRandomRoom(opJoinRandomRoomParams);

```

Create a new one if none found:

C#

```csharp
EnterRoomParams enterRoomParams = new EnterRoomParams();
enterRoomParams.ExpectedUsers = teamMembersUserIds;
loadBalancingClient.OpCreateRoom(enterRoomParams);

```

The others don't have to do any matchmaking but instead repeatedly call ('periodic poll', every few frames/(milli)seconds):

C#

```csharp
loadBalancingClient.OpFindFriends(new string[1]{ leaderUserId });

```

When the leader arrives in a room, the `FindFriends` operation will reveal that room's name and everyone can join it:

C#

```csharp
EnterRoomParams enterRoomParams = new EnterRoomParams();
enterRoomParams.RoomName = roomNameWhereTheLeaderIs;
loadBalancingClient.OpJoinRoom(enterRoomParams);

```

## Lobbies

Photon is organizing your rooms in so called "lobbies".

So all rooms belong to lobbies.

Lobbies are identified using their name and type.

The name can be any string, however there are only 3 types of lobbies: [Default](#default-lobby-type), [SQL](#sql-lobby-type) and [Async](#asynchronous-random-lobby-type).

Each one has a unique capability which suits specific use cases.

All applications start with a preexisting lobby: [The Default Lobby](#the-default-lobby).

Most applications won't need other lobbies.

However, clients can create other lobbies on the fly.

Lobbies begin to exist when you specify a new lobby definition in operation requests: `JoinLobby`, `CreateRoom` or `JoinOrCreateRoom`.

Like rooms, lobbies can be joined and you can leave them.

In a lobby, the clients only get the room list of that lobby when applicable.

Nothing else.

There is no way to communicate with others in a lobby.

When a client is joined to a lobby and tries to create (or `JoinOrCreate`) a room without explicitly setting a lobby, if the creation succeeds/happens, the room will be added to the currently joined lobby.

When a client is not joined to a lobby and tries to create (or `JoinOrCreate`) a room without explicitly setting a lobby, if the creation succeeds/happens, the room will be added to [the default lobby](#the-default-lobby).

When a client is joined to a lobby and tries to create (or `JoinOrCreate`) a room by explicitly setting a lobby, if the creation succeeds/happens:

- if the lobby name is null or empty: the room will be added to the currently joined lobby.


This means you cannot create rooms in [the default lobby](#the-default-lobby) when you are joined to a custom/different one.
- if the lobby name is not null nor empty: the room will be added to the lobby specified by the room creation request.

When a client is joined to a lobby and tries to join a random room without explicitly setting a lobby, the server will look for the room in the currently joined lobby.

When a client is not joined to a lobby and tries to join a random room without explicitly setting a lobby, the server will look for the room in [the default lobby](#the-default-lobby).

When a client is joined to a lobby and tries to join a random room a room by explicitly setting a lobby:

- if the lobby name is null or empty: the server will look for the room in the currently joined lobby.


This means you cannot join random rooms in [the default lobby](#the-default-lobby) when you are joined to a custom/different one.
- if the lobby name is not null nor empty: the server will look for the room in the lobby specified by the room creation request.

When a client is joined to a lobby and wants to switch to a different one, you can call JoinLobby directly and no need to leave the first one by calling LeaveLobby explicitly.

### Default Lobby Type

The most suited type for _synchronous_ **random matchmaking**.

Probably the less sophisticated and most used type.

While joined to a default lobby type, the client will receive periodic room list updates.

When the client joins a lobby of default type, it instantly gets an initial list of available rooms.

After that the client will receive periodic room list updates.

The list is sorted using two criteria: open or closed, full or not.

So the list is composed of three groups, in this order:

- first group: open and not full (joinable).
- second group: full but not closed (not joinable).
- third group: closed (not joinable, could be full or not).

In each group, entries do not have any particular order (random).

The list of rooms (or rooms' updates) is also limited in number, see [Lobby Limits](#lobby-limits).

C#

```csharp
using Photon.Realtime;
using System.Collections.Generic;
public class RoomListCachingExample : ILobbyCallbacks, IConnectionCallbacks
{
    private TypedLobby customLobby = new TypedLobby("customLobby", LobbyType.Default);
    private LoadBalancingClient loadBalancingClient;
    private Dictionary<string, RoomInfo> cachedRoomList = new Dictionary<string, RoomInfo>();
    public void JoinLobby()
    {
        loadBalancingClient.JoinLobby(customLobby);
    }
    private void UpdateCachedRoomList(List<RoomInfo> roomList)
    {
        for(int i=0; i<roomList.Count; i++)
        {
            RoomInfo info = roomList[i];
            if (info.RemovedFromList)
            {
                cachedRoomList.Remove(info.Name);
            }
            else
            {
                cachedRoomList[info.Name] = info;
            }
        }
    }
    // do not forget to register callbacks via loadBalancingClient.AddCallbackTarget
    // also deregister via loadBalancingClient.RemoveCallbackTarget
    #region ILobbyCallbacks
    void ILobbyCallbacks.OnJoinedLobby()
    {
        cachedRoomList.Clear();
    }
    void ILobbyCallbacks.OnLeftLobby()
    {
        cachedRoomList.Clear();
    }

    void ILobbyCallbacks.OnRoomListUpdate(List<RoomInfo> roomList)
    {
        // here you get the response, empty list if no rooms found
        UpdateCachedRoomList(roomList);
    }
    // [..] Other callbacks implementations are stripped out for brevity, they are empty in this case as not used.
    #endif
    #region IConnectionCallbacks
    void IConnectionCallbacks.OnDisconnected(DisconnectCause cause)
    {
        cachedRoomList.Clear();
    }
    // [..] Other callbacks implementations are stripped out for brevity, they are empty in this case as not used.
    #endregion
}

```

#### The Default Lobby

It has a `null` name and its type is [Default Lobby Type](#default-lobby-type).

In C# SDKs, it's defined in `TypedLobby.Default`.

The default lobby's name is reserved:

only the default lobby can have a `null` name, all other lobbies need to have a name string that is not null nor empty.

If you use a string empty or null as a lobby name it will point to the default lobby nomatter the type specified.

#### Recommended Flow

We encourage everyone to skip joining lobbies unless abolutely necessary.

If needed, when you want rooms to be added to specific or custom lobbies, the client can specify the lobby when creating new rooms.

Joining lobbies of default type will get you the list of rooms, but it's not useful in most cases:

- there is no difference in terms of ping between the entries of the list
- usually players are looking for a quick match
- receiving rooms list adds an extra delay and consumes traffic
- a long list with too much information can have a bad effect on the user experience

Instead, to give your players more control over the matchmaking, use filters for random matchmaking.

Multiple lobbies can still be useful, as they are also used in (server-side) random matchmaking and you could make use of lobby statistics.

### SQL Lobby Type

In SQL lobby type, string filters in `JoinRandomRoom` replace the default expected lobby properties.

Also, in SQL lobby type, only one MatchmakingMode is supported: `FillRoom` (default, 0).

Besides " [Custom Room Listing](#custom-room-listing)" replaces the automatic periodic rooms listing which exists only in the default lobby type.

This lobby type adds a more elaborate matchmaking filtering which could be used for a server-side [skill-based matchmaking](#skill-based-matchmaking) that's completely client-driven.

Internally, SQL lobbies save rooms in a SQLite table with up to 10 special "SQL filtering properties".

The naming of those SQL properties is fixed as: "C0", "C1" up to "C9".

Only integer-typed and string-typed values are allowed and once a value was assigned to any column in a specific lobby, this column is locked to values of that type.

Despite the static naming, clients have to define which ones are needed in the lobby.

Be careful as SQL properties are case sensitive when you define them as lobby properties or set their values but are not case sensitive inside SQL filters.

You can still use custom room properties other than the SQL properties, visible or invisible to the lobby, during room creation or after joining it.

Those will not be used for matchmaking however.

Queries can be sent in `JoinRandomRoom` operation.

The filtering queries are basically SQL WHERE conditions based on the "C0" .. "C9" values.

Find the list of all SQLite supported operators and how to use them [here](https://sqlite.org/lang_expr.html#binaryops).

Take into consideration the [excluded keywords](#excluded-sql-keywords).

**Example:**

C#

```csharp
using Photon.Realtime;
using System.Collections.Generic;
using Hashtable = ExitGames.Client.Photon.Hashtable;
public class RandomMatchmakingExample : IMatchmakingCallbacks
{
    public const string ELO_PROP_KEY = "C0";
    public const string MAP_PROP_KEY = "C1";
    private TypedLobby sqlLobby = new TypedLobby("customSqlLobby", LobbyType.SqlLobby);
    private LoadBalancingClient loadBalancingClient;
    private void CreateRoom()
    {
        RoomOptions roomOptions = new RoomOptions();
        roomOptions.CustomRoomProperties = new Hashtable { { ELO_PROP_KEY, 400 }, { MAP_PROP_KEY, "Map3" } };
        roomOptions.CustomRoomPropertiesForLobby = { ELO_PROP_KEY, MAP_PROP_KEY }; // makes "C0" and "C1" available in the lobby
        EnterRoomParams enterRoomParams = new EnterRoomParams();
        enterRoomParams.RoomOptions = roomOptions;
        enterRoomParams.Lobby = sqlLobby;
        loadBalancingClient.OpCreateRoom(enterRoomParams);
    }
    private void JoinRandomRoom()
    {
        string sqlLobbyFilter = "C0 BETWEEN 345 AND 475 AND C1 = 'Map2'";
        //string sqlLobbyFilter = "C0 > 345 AND C0 < 475 AND (C1 = 'Map2' OR C1 = \"Map3\")";
        //string sqlLobbyFilter = "C0 >= 345 AND C0 <= 475 AND C1 IN ('Map1', 'Map2', 'Map3')";
        OpJoinRandomRoomParams opJoinRandomRoomParams = new OpJoinRandomRoomParams();
        opJoinRandomRoomParams.SqlLobbyFilter = sqlLobbyFilter;
        loadBalancingClient.OpJoinRandomRoom(opJoinRandomRoomParams);
    }
    // do not forget to register callbacks via loadBalancingClient.AddCallbackTarget
    // also deregister via loadBalancingClient.RemoveCallbackTarget
    #region IMatchmakingCallbacks
    void IMatchmakingCallbacks.OnJoinRandomFailed(short returnCode, string message)
    {
        CreateRoom();
    }
    void IMatchmakingCallbacks.OnCreateRoomFailed(short returnCode, string message)
    {
        Debug.LogErrorFormat("Room creation failed with error code {0} and error message {1}", returnCode, message);
    }
    void IMatchmakingCallbacks.OnJoinedRoom()
    {
        // joined a room successfully, both JoinRandomRoom or CreateRoom lead here on success
    }
    // [..] Other callbacks implementations are stripped out for brevity, they are empty in this case as not used.
    #endif
}

```

#### Chained Filters

You can send up to 3 comma separated filters at once in a single `JoinRandomRoom` operation.

These are called chained filters.

Photon servers will try to use the filters in order.

A room will be joined if any of the filters matches a room.

Otherwise a NoMatchFound error will be returned to the client.

Chained filters could help save matchmaking requests and speed up its process.

It could be useful especially for [skill-based matchmaking](#skill-based-matchmaking) where you need to 'relax' the filter after failed attempt.

Possible filters string formats:

- 1 (min) filter value: `{filter1}` (or `{filter1};`)
- 2 filter values: `{filter1};{filter2}` (or `{filter1};{filter2};`)
- 3 (max) filter values: `{filter1};{filter2};{filter3}` (or `{filter1};{filter2};{filter3};`)

**Examples**:

- `C0 BETWEEN 345 AND 475`
- `C0 BETWEEN 345 AND 475;C0 BETWEEN 475 AND 575`
- `C0 BETWEEN 345 AND 475;C0 BETWEEN 475 AND 575;C0 >= 575`

#### Custom Room Listing

Client can also request a custom list of rooms from an SqlLobby using SQL-like queries.

This method will return up to 100 rooms that fit the conditions.

The returned rooms are joinable (i.e. open and not full) and visible.

C#

```csharp
using Photon.Realtime;
using System.Collections.Generic;
public class GetCustomRoomListExample : ILobbyCallbacks
{
    private TypedLobby sqlLobby = new TypedLobby("customSqlLobby", LobbyType.SqlLobby);
    public void GetCustomRoomList(string sqlLobbyFilter)
    {
      loadBalancingClient.OpGetGameList(sqlLobby, sqlLobbyFilter);
    }
    // do not forget to register callbacks via loadBalancingClient.AddCallbackTarget
    // also deregister via loadBalancingClient.RemoveCallbackTarget
    #region ILobbyCallbacks
    void ILobbyCallbacks.OnRoomListUpdate(List<RoomInfo> roomList)
    {
        // here you get the response, empty list if no rooms found
    }
    // [..] Other callbacks implementations are stripped out for brevity, they are empty in this case as not used.
    #endif
}

```

#### Skill-based Matchmaking

You can use lobbies of the SQL-type to implement your own skill-based matchmaking.

First of all, each room gets a fixed skill that players should have to join it.

This value should not change, or else it will basically invalidate any matching the players in it did before.

As usual, players should try to get into a room by `JoinRandomRoom`.

The filter should be based on the user's skill.

The client can easily filter for rooms of "skill +/- X".

`JoinRandomRoom` will get a response as usual but if it didn't find a match right away, the client should wait a few seconds and then try again.

You can do as many or few requests as you like.

If you use SQL lobby type, you could make use of [Chained Filters](#chained-sql-filters).

Best of all: The client can begin to relax the filter rule over time.

It's important to relax the filters after a moment.

Granted: A room might be joined by a player with not-so-well-fitting skill but obviously no other room was a better fit and it's better to play with someone.

You can define a max deviation and a timeout.

If no room was found, this client has to open a new room with the skill this user has.

Then it has to wait for others doing the same.

Obviously, this workflow might take some time when few rooms are available.

You can rescue your players by checking the "application stats" which tell you how many rooms are available.

See [Matchmaking For Low CCU](#matchmaking-for-low-ccu).

You can adjust the filters and the timing for "less than 100 rooms" and use different settings for "100 to 1000 rooms" and again for "even more".

#### Excluded SQL Keywords

SQL filters will not accept the following keywords:

- ALTER
- CREATE
- DELETE
- DROP
- EXEC
- EXECUTE
- INSERT
- INSERT INTO
- MERGE
- SELECT
- UPDATE
- UNION
- UNION ALL

If you use any of these words in the SQL filter string the corresponding operation will fail.

### Lobby Types Comparison

| LobbyType | Periodic Rooms List Updates | SQL Filter | Max Players Filter | Custom Room Properties Filter | Matchmaking Modes | Removed Rooms Entries TTL (minutes) |
| --- | --- | --- | --- | --- | --- | --- |
| Default |  |  |  |  |  | 0 |
| SQL |  |  |  |  |  | 0 |

## Matchmaking For Low CCU

For really good matchmaking, a game needs a couple hundred players online.

With less players online, it will become harder to find a worthy opponent and at some point it makes sense to just accept almost any match.

You have to take this into account when you build a more elaborate matchmaking on the client side.

To do so, the Photon Master Server provides the count of connected users, rooms and players (in a room), so you can adjust the client-driven matchmaking at runtime.

The number of rooms should be a good, generic indicator of how busy the game currently is.

You could obviously also fine tune the matchmaking by on how many players are not in a room.

Whoever is not in a room might be looking for one.

For example, you could define a low CCU situation as less than 20 rooms.

So, if the count of rooms is below 20, your clients use no filtering and instead run the [Quick Match](#quick-match) routine.

### Testing Matchmaking Early In Development

If you are testing matchmaking early in development phase and try to join a random room from two clients at about the same time, there is a chance both clients end up on different rooms: this happens because join random room will not return a match for both clients and each one will probably create a new room as none found.

So this is expected and OK.

To avoid this use JoinRandomOrCreateRoom (see [Quick Match](#quick-match)) instead of JoinRandomRoom then CreateRoom.

Otherwise, a possible workaround (for development purposes only) would be to add a random delay before (or after) attempting to join a room or retry again.

You could also listen for application or lobby statistics to make sure a room at least exists or has been created.

## External Matchmaking

While the built in matchmaking of Photon is simple and capable, there are of course use cases it can not deliver.

For example there is no global matchmaking and it is tricky to organize tournaments without some backend to control it.

No matter how your external matchmaking is implemented in detail, for the clients it just has to provide a unique room name.

The clients have to communicate directly with the backend to ask for a room and then they use `JoinRoomOrCreate` to enter it.

Unless only a single Photon Region is used, make sure clients are connected to the correct region for the room.

Photon can call WebHooks on your backend to signal when users join or leave rooms.

There is even a "Before Join" webhook, which can be used to verify that the user is assigned to the room.

To increase security, an external matchmaking backend could issue so called "Matchmaking Tickets" instead of sending a room name.

Matchmaking Tickets are encrypted and hide the matchmaking information from the clients.

They can be made mandatory on the Photon side to secure the matchmaking entirely from client side exploits.

## Lobby Versions

There are two versions of Photon Lobbies. v2 is the default for new Apps being created.

The differences are:

- v1 SQL lobbies used to send rooms lists.


In v2 they don't.
- v1 room lists were not limited.


In v2 they are capped.
- In v1 Room lists were not sorted by joinability.


In v2 the joinable ones are first on the list.

In the dashboard, you can check per app which lobby version it uses. Click "manage" and check the "Details".

### Lobby v2 Limits

Photon has the following lobbies related default limits:

- Maximum number of lobbies per application: 10000.
- Maximum number of room list entries in GameList events (initial list when you join lobbies of type Default): 500.
- Maximum number of updated rooms entries in GameListUpdate events (when joined to lobbies of type Default): 500.


This limit does not account for removed rooms entries (corresponding to rooms no longer visible or simply gone).
- Maximum number of room list entries in GetGameList operation response (SQL Lobby): 100.
- Room list updates have no limit for removed rooms.

Notes:

Limiting the number of entries per update makes sure clients can handle the incoming traffic well.

One resulf ot this is that clients can not expect to get the full list of rooms (for bigger titles).

These limits are purely about lists and updates sent to clients. In both versions, the lobbies on the server can contain a much larger number of room entries.

Back to top

- [Matchmaking Checklist](#matchmaking-checklist)
- [Quick Match](#quick-match)

  - [JoinRandomOrCreateRoom](#joinrandomorcreateroom)
  - [JoinRandomRoom or CreateRoom](#joinrandomroom-or-createroom)

- [Random Matchmaking](#random-matchmaking)

  - [Exposing Some Properties In The Lobby](#exposing-some-properties-in-the-lobby)
  - [Filtering Room Properties in Join Random](#filtering-room-properties-in-join-random)

- [Play with Your Friends](#play-with-your-friends)

  - [Publishing UserIDs in a Room](#publishing-userids-in-a-room)
  - [Matchmaking Slot Reservation](#matchmaking-slot-reservation)

- [Lobbies](#lobbies)

  - [Default Lobby Type](#default-lobby-type)
  - [SQL Lobby Type](#sql-lobby-type)
  - [Lobby Types Comparison](#lobby-types-comparison)

- [Matchmaking For Low CCU](#matchmaking-for-low-ccu)

  - [Testing Matchmaking Early In Development](#testing-matchmaking-early-in-development)

- [External Matchmaking](#external-matchmaking)
- [Lobby Versions](#lobby-versions)
  - [Lobby v2 Limits](#lobby-v2-limits)