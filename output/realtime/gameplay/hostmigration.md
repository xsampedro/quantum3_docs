# hostmigration

_Source: https://doc.photonengine.com/realtime/current/gameplay/hostmigration_

# Master Client and Host Migration

"Host Migration" is a known concept in online multiplayer games.

It covers questions about how to do a smooth and seamless transition of the "host" peer.

The "host" is a peer that has more control over the game and on which the game depends the most.

Usually the "host" is the client who starts ("hosts") the game and other players need to connect to him or join him to be able to play.

In Photon's universe, we do not really have a "host" per se.

Instead we have a "special" client per room, called "Master Client".

By default it is a normal client like all others.

It stays that way until you decide to make him do more tasks than the others.

### What Photon Does For You

There is no generic solution in Photon for "Host Migration".

However, we make things much easier for you.

Photon Server detects when a Master Client disconnects and assigns another actor in the room as the new Master Client.

By default, the active actor with the lowest actor number is chosen.

Whenever the Master Client changes, the following callback will be called:

C#

```csharp
void IInRoomCallbacks.OnMasterClientSwitched(Player newMasterClient)

```

Photon also offers a way of explicitly changing the Master Client from clients, server and plugins SDKs.

From client you can do this using:

C#

```csharp
loadBalancingClient.CurrentRoom.SetMasterClient(photonPlayer);

```

The method will return if the operation could be sent or not.

### What You Need To Do

Photon does not hand all information that the former Master Client had about the room's state to the newly elected one.

It is not an obvious thing to do and sometimes it is too late to do it as the Master Client is already gone (disconnected, not responding, etc.).

Photon does not move player properties or cached events from one Master Client to another.

Also Photon does not resend events meant to the old Master Client to the newly elected one.

That is why this is the developer's responsibility to make sure that no room state is lost when the switch happens.

## Recommendations

- Master Client concept is not suitable for all kind of games.


For instance for fast paced games, you should rely on custom server side code instead (self-hosted server applications or plugins).
- The more players you have in your game the less you need to rely on a single client to do extra tasks.
- You should avoid sending events to Master Client only. Replicating data is better than losing it.


When you send to all actors you make any eventual Master Client switch easier.


Since all actors will keep track of the room's state, the newly assigned Master Client will have the full room state already.
- Custom room properties are the best reliable option if you want to persist game data in case a Master Client switch happens.
- If your Master Client is constantly sending events in the room, you can implement a Master Client timeout detection mechanism by saving timestamp of last event received from Master Client and checking that value constantly in your game loop.


When you detect that the Master Client is not sending events as expected you can then switch Master Client explicitly.


The tricky part is to choose the timeout value.


A low value may give false positives and make the Master Client switch happen too often.


Also a choice needs to be made on whether this check needs to be done from a single actor (next Master Client candidate maybe) or all actors.

Back to top

- [What Photon Does For You](#what-photon-does-for-you)
- [What You Need To Do](#what-you-need-to-do)

- [Recommendations](#recommendations)