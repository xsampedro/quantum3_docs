# cached-events

_Source: https://doc.photonengine.com/realtime/current/gameplay/cached-events_

# Cached Events

Photon events have a central role for game logic and players communication inside a room.

Clients call `RaiseEvent` (exact function or method name may vary depending on your client SDK) to send data to one or more players in the room.

All joined players will be able to receive those events immediately.

However, late joining players will miss those events unless you use the caching option.

When you cache an event, Photon server will save it in the room state to keep it available for players who will join the room later.

Photon server will then send it to the joining player first, before forwarding any eventual "live" events.

**Use case 1:**

You can cache a "start time" event containing the time stamp of a round creation as data.

This way, anyone who joins the game will know when the round has started (and probably when it should end).

## Ordered Delivery Is Guaranteed

Joining players get cached events in the same order as those events have arrived on the server.

The client first reads the room and player properties (so each player is known) and then gets the cached events.

All events that have been sent since joining will be delivered afterwards.

So you can consider the cached events reception order identical to the order of their transmission.

Exceptions may only occur when UDP gets used and events are sent unreliable.

## Understand How Events Are Cached

Without going too much into details, here is what you need to know about how events are cached:

Each cached event can be defined by its code, its data and the actor number of the sender.

The event cache can also be grouped into two "logical partitions":

- "global cache": the cache associated to the room (ActorNr == 0).


It contains all the cached events that have been sent with the `EventCaching.AddToRoomCacheGlobal` flag and which have not explicitly been removed.


These cached events can no longer be traced back to their original sender. That is why these events are sent to any joining actor.
- "actor cache": the cache associated to a particular actor number (ActorNr != 0).


It contains all the cached events that have been sent by that actor and which have not explicitly been added to the global cache or removed from the cache.


These events are sent to any joining actor except their original respective sender.

## Take Control of the Event Cache

In the C# LoadBalancing API and in PUN users need to pass a `RaiseEventOptions` object to any `RaiseEvent` (exact function or method name may vary depending on your client SDK) call.

This parameter includes a `CachingOption` property.

Let's discover how each of the possible `EventCaching` values affect the events cache:

- `DoNotCache`: this is the default value. It indicates that the event to be sent will not be added to the room's cache.
- `AddToRoomCache`: this indicates that the event to be sent will be added to the room cache for the sending actor.


It will be marked by the number of the sending actor.
- `AddToRoomCacheGlobal`: this indicates that the event to be sent will be added to the "global cache".


Be careful when using this value as with this code the event only gets removed when you explicitly request its removal. Otherwise, it will have the same lifetime as the room.
- `RemoveFromRoomCache`: this indicates that all previously cached events that match the specified "filter pattern" will be removed from the cache.


The "filter pattern" is a combination of three parameters: event code, sender number and event data. You can use one, two or all three filters.

  - An event code equal to 0 is a wild card for event codes.
  - Use target actors option (in C# client SDKs, this is done using `RaiseEventOptions.TargetActors` array) to specify the sender number in this case.


    The sender number can be 0, which means you can remove from "global cache".


    So to remove from global cache you could specify 0 as sender number.


    Since the target actors is an array, you could filter by multiple actors, even combining global (ActorNr == 0) and non-global cached events (ActorNr != 0).
  - Also if you filter by event data, you are free to send only a part of the data.


    For instance, if you use a `Hashtable` as event data, you can remove events only by specifying one or more key/value pairs.
- `RemoveFromRoomCacheForActorsLeft`: this indicates that the cached events sent by removed actors should also be removed.


This does not make any difference if you create rooms with `RoomOptions.CleanupCacheOnLeave` set to true which is the default value.

Events will not be added to cache if any of the following conditions is met:

- `RaiseEventOptions.Receivers == ReceiverGroups.MasterClient`.
- `RaiseEventOptions.TargetActors != null`.
- `RaiseEventOptions.InterestGroups != 0`.

**Use case 2:**

If you use a `Hashtable` as event content, you can mark all events belonging to some object with an "oid" key (short form of "ObjectID") and some value.

When you want to clean up the cache related to a specific object, you can just send a `Hashtable(){"oid", <objectId>}` as event data filter and `EventCaching.RemoveFromRoomCache`.

## Cache Cleaning

When a player quits a game, usually his/her actions are no longer relevant for new players.

To avoid congestion on join, Photon server by default automatically cleans up events that have been cached by a player, that has left the room for good.

If you want to manually clean up the rooms' event cache you can create rooms with `RoomOptions.CleanupCacheOnLeave` set to false.

## Special Considerations

- The Photon client cannot start calling operations (RaiseEvent, SetProperties, etc.) inside a room only when the actor is fully joined.


An actor is not considered fully joined only when the server finishes sending all cached events.


The client will receive an error (`OperationNotAllowedInCurrentState (-3)`) when trying to call operations inside a room before receiving all cached events first.
- Cached events also delay live events that are sent after the actor has joined the room.
- Photon limits the total number of cached events per room (actors cache and global cache combined).


If you hit the limit, the server will broadcast an ErrorInfo event and close the room so no new actor can join.


Actors already joined can remain while inactive actors will not be able to rejoin.

Back to top

- [Ordered Delivery Is Guaranteed](#ordered-delivery-is-guaranteed)
- [Understand How Events Are Cached](#understand-how-events-are-cached)
- [Take Control of the Event Cache](#take-control-of-the-event-cache)
- [Cache Cleaning](#cache-cleaning)
- [Special Considerations](#special-considerations)