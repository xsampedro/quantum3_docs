# userids-and-friends

_Source: https://doc.photonengine.com/realtime/current/lobby-and-matchmaking/userids-and-friends_

# UserIDs and Friends

## UserIDs

In Photon, a player is identified using a unique UserID. In best case, across multiple sessions.

UserIDs are case sensitive. **Example**: " **Photonian**" and " **photonian**" are two different UserIDs for two different players.

This UserID is useful inside and outside rooms.

Photon clients with the same UserID can be connected to the same server but you can't join the same Photon room from two separate clients using the same UserID.

Each actor inside the room should have a unique UserID.

In old C# client SDKs, this was enabled using `RoomOptions.CheckUserOnJoin`.

### Unique UserIDs

Unlike usernames, displaynames or nicknames, the UserIDs are not intended to be shown. They are often not very readable or memorable.

It is not uncommon to use a [GUID](http://guid.one/guid) as a UserID.

The advantages of keeping a unique UserID per player:

- You preserve your data between game sessions and across multiple devices.


You can rejoin rooms and resume playing where you stopped.
- You can become known to all players you meet and easily identifiable by everyone.


You can [play with your friends](#friends), send them invitations and challenges, make online parties, form teams and guilds, etc.


You can add user profiles (e.g. experience, statistics, achievements, levels, etc.) and make games more challenging (also using tournaments and leaderboards).
- You could make use of another service to bind Photon UserID to an external unique identifier.


For instance, Photon UserID could be set to Facebook ID, Google ID, Steam ID, PlayFab ID, etc.
- You can prohibit malicious users from connecting to your applications by keeping a blocklist of their UserIDs and making use of Custom Authentication.

### Setting UserIDs

Once authenticated, a Photon client will keep the same UserID until disconnected.

The UserID for a client can be set in three ways:

1. Client sends its UserID before connecting by setting `AuthenticationValues.UserId`.


This option is useful when you do not use Custom Authentication and just want to set a UserID.
2. An external authentication provider returns the UserID on successful authentication. This overrides any value sent by the client.
3. Photon Server will assign GUIDs as IDs for users that did not get UserIDs using 1 or 2. So even anonymous users will have a temporary UserID.

More about [Custom Authentication](/realtime/current/connection-and-authentication/authentication/custom-authentication#returning-data-to-client).

### Publish UserIDs

Players can share their UserIDs with each other inside rooms.

In C# SDKs, to enable this and make the UserID visible to everyone, set `RoomOptions.PublishUserId` to `true`, when you create a room.

The server will then broadcast this information on each new join and you can access each player's UserID using `Player.UserId`.

### Matchmaking Slot Reservation

Sometimes, a player joins a room, knowing that a friend should join as well.

With Slot Reservation, Photon can block a slot for specific users and take that into account for matchmaking.

To reserve slots there is an `expectedUsers` parameter (exact parameter or argument name may vary depending on your client SDK) in the methods that get you in a room (`JoinRoom`, `JoinOrCreateRoom`, `JoinRandomRoom` and `CreateRoom`. Exact functions or methods names may vary depending on your client SDK).

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

## Friends

You can find out if your friends who are playing the same game are online and if so, which room they are joined to.

Like users, friends are also identified using their UserID.

So the FriendID is the same thing as the UserID and in order to find some friends you need to know their UserIDs first.

Friends can only find one another, if they are connected to the same Region with the same AppId and AppVersion.

Friends make use of UserIDs, which are case sensitive.

FindFriends is only available on the Master Server. Clients can not find friends while playing a session in a room.

Examples:

C#

```csharp
using System.Collections.Generic;
using Photon.Realtime;
public class FindFriendsExample : IMatchmakingCallbacks
{
    private LoadBalancingClient loadBalancingClient;
    public bool FindFriends(string[] friendsUserIds)
    {
        return loadBalancingClient.OpFindFriends(friendsUserIds);
    }
    // do not forget to register callbacks via loadBalancingClient.AddCallbackTarget
    // also deregister via loadBalancingClient.RemoveCallbackTarget
    #region IMatchmakingCallbacks
    public override void IMatchmakingCallbacks.OnFriendListUpdate(List<FriendInfo> friendsInfo)
    {
        for(int i=0; i < friendsInfo.Count; i++)
        {
            FriendInfo friend = friendsInfo[i];
            Debug.LogFormat("{0}", friend);
        }
    }
    // [..] Other callbacks implementations are stripped out for brevity, they are empty in this case as not used.
    #endif
}

```

Back to top

- [UserIDs](#userids)

  - [Unique UserIDs](#unique-userids)
  - [Setting UserIDs](#setting-userids)
  - [Publish UserIDs](#publish-userids)
  - [Matchmaking Slot Reservation](#matchmaking-slot-reservation)

- [Friends](#friends)