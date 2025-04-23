# interestgroups

_Source: https://doc.photonengine.com/realtime/current/gameplay/interestgroups_

# Interest Groups

You can imagine Photon's "Interest Groups" as sub-channels for conversations in a room:

Clients only get the messages of Interest Groups they subscribed to (and group 0).

They can send any event to any group they want to.

This simple feature can be used for basic Interest Management or any way you see fit.

See the example use cases below.

## Available Groups

Clients do not need to explicitly create interest groups.

Interest groups are created on demand; when an actor subscribes to a new group number, the server will create it.

Photon offers up to 256 interest groups.

The group number 0 is reserved and meant for broadcast: All actors (clients) inside a room are subscribed to group 0 and cannot unsubscribe from it.

The other 255 groups are available to the developer to use freely.

Any event assigned to a group > 0 will only be transmitted to clients that are interested in that group and in the room when the server relays the event.

**Important:** You can only cache events sent to interest group 0! No other group has an event cache.

## Global Setup

The list of subscribed groups can be updated at any time inside the room by adding or removing their numbers using:

C#

```csharp
bool LoadBalancingPeer.OpChangeGroups(byte[] groupsToRemove, byte[] groupsToAdd)

```

Notes:

- Priority is always to group addition: if the same group number is added to both arrays then the group will be added.
- A `null` array acts as "no group" and an empty array (`new byte\[0\]`) acts as "all groups".
- `LoadBalancingPeer.OpChangeGroups(new byte\[0\], groupsToAdd)` will remove all groups except the ones in `groupsToAdd`.
- `LoadBalancingPeer.OpChangeGroups(groupsToRemove, new byte\[0\])` will add all groups no matter what the value of `groupsToRemove` is.
- Since the `ChangeGroups` operation does not return a response or trigger an event, each client should cache his interest groups locally if needed.

## Example Use Cases

Interest Groups to subscribe to, per client, can be defined dynamically at runtime or statically at compile time.

Interest Groups setup can be the same on all clients or different per client.

Interest Groups are useful to lower the number of messages per second inside rooms.

By lowering traffic, you stay below the message/second limit, cut costs and sometimes it can help you increase the number of maximum players per room.

But you can come up with other clever ways of using Interest Groups in your game.

### Network Culling

Most common use case of Interest Groups is Network Culling.

Interest Groups could be mapped to areas of interest in your game.

For instance if you have a "big fat world" you can separate it into smaller chunks, let's call them "areas" and assign an interest group per "area".

### Team Events

If you have teams in your game and want to implement team exclusive events, you can assign an interest group per team.

All team members should subscribe to the team's interest group.

Inter team events should be sent using the team's own group number.

Intra team events should be sent using an opponent team's group number.

Back to top

- [Available Groups](#available-groups)
- [Global Setup](#global-setup)
- [Example Use Cases](#example-use-cases)
  - [Network Culling](#network-culling)
  - [Team Events](#team-events)