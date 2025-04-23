# webhooks-faq

_Source: https://doc.photonengine.com/realtime/current/gameplay/web-extensions/webhooks-faq_

# Webhooks FAQ

### How to rejoin rooms?

Actors can rejoin rooms before `RoomOptions.PlayerTTL` milliseconds of their deactivation time by calling `OpReJoinRoom()` from client SDK.

Deactivation time is when an actor becomes inactive inside the room.

An actor becomes inactive inside a room when it unexpectedly disconnects or leaves the room temporarily (without abandoning).

### What if I just set `EmptyRoomTtl` to `int.MaxValue` during new room creation?

This is not possible and the Photon Cloud will return an error indicating that `EmptyRoomTtl` should be modified to a lower value.

The threshold of `EmptyRoomTtl` in Photon Cloud is set to 300000 milliseconds (5 minutes) while this value could be changed for privately hosted Photon servers or Enterprise Cloud.

For realtime games, a best practice value of `EmptyRoomTtl` is 12 seconds (12,000 milliseconds).

### Can webhooks modify the room's state?

The short answer is no.

However when you load a saved room state you can change few of its _non binary top level_ properties.

### Is there any interesting use case for webhooks other than storing and retrieving the room's state?

Well, there are many things possible with webhooks :)

Here's a list of few :

- Server side authoritative code:


This can be useful when some specific game logic need to be executed at server side only.
- Protection against malicious hackers:


Webhooks can also be handy when it comes to detect cheaters.


PathProperties and PathEvent webhooks can be used to authorize players' actions from server side code.
- Sending Push Notifications from server side:


With PathEvent webhook sending push notifications can be optimized by sending it from server side to inactive players only.
- Save game data other than Photon's room state:


It is not an obligation to use only Photon's room properties or events cache as game data.
- Analytics:


Webhooks can be a powerful and free of charge analytics service.


The statistics you can track are numerous like room creations, actors joins, room events, ...

### Can webhooks modify room events data or send new Photon events?

No, webhooks cannot alter the data of received events nor send any other events.

### Is it possible to access custom room properties from webhooks?

The short answer is no.

However, the `State` object, available as an argument in "PathEvent", "PathProperties" and "PathClose" of `Type=='Save'` webhooks contains a readable form of public custom properties (i.e. visible to the lobby).

### When should an actor be considered as left or flagged as inactive?

Any Photon actor disconnected from a GameServer has left a room somehow.

The question is whether the same actor left the room intentionally - which means that he/she explicitly "abandoned" or "quit" or "exited" the room - or not.

If left intentionally, the actor will no longer be part of the room's actors and thus removed from the room's ActorList.

On the contrary, if the actor temporary left the room willing to rejoin it later then he/she will be marked as `Inactive` or `IsActive = false` or `IsComingBack`.

Any actor joined to a room and still connected to its respective GameServer will always be considered as `Active`.

Actors status is exposed in the room's ActorList in the PathEvent webhook.

See its respective section for more information.

To abandon a room, you should call `OpLeave(false)` and to leave it for a while just use `Disconnect() or OpLeave(true)`.

For possible scenarios like "resignation" or "rage quit" specific to some game types, the user could choose one of the two approaches.

See the [PathLeave webhook](/realtime/current/gameplay/web-extensions/webhooks#leave) section for additional information.

### Does it make sense to have `0 <= PlayerTTL <= EmptyRoomTTL` with `IsPersistent = true` ?

`PlayerTTL` is the amount of time that an actor can stay inactive inside a room. Meaning, how long a player can stay disconnected from a room. If you create rooms with `PlayerTTL = 0` then you should not expect players to come back. And if they don't come back then there is no need to save the room state for them. This is by design, a room state should be saved only if it contains at least one inactive player.

See the [PathLeave webhook](/realtime/current/gameplay/web-extensions/webhooks#leave) section for additional information.

Back to top

- [How to rejoin rooms?](#how-to-rejoin-rooms)
- [What if I just set EmptyRoomTtl to int.MaxValue during new room creation?](#what-if-i-just-set-emptyroomttl-to-int.maxvalue-during-new-room-creation)
- [Can webhooks modify the room's state?](#can-webhooks-modify-the-rooms-state)
- [Is there any interesting use case for webhooks other than storing and retrieving the room's state?](#is-there-any-interesting-use-case-for-webhooks-other-than-storing-and-retrieving-the-rooms-state)
- [Can webhooks modify room events data or send new Photon events?](#can-webhooks-modify-room-events-data-or-send-new-photon-events)
- [Is it possible to access custom room properties from webhooks?](#is-it-possible-to-access-custom-room-properties-from-webhooks)
- [When should an actor be considered as left or flagged as inactive?](#when-should-an-actor-be-considered-as-left-or-flagged-as-inactive)
- [Does it make sense to have 0 <= PlayerTTL <= EmptyRoomTTL with IsPersistent = true ?](#does-it-make-sense-to-have-0-playerttl-emptyroomttl-with-ispersistent-true)