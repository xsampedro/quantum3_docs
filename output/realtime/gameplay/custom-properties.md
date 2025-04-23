# custom-properties

_Source: https://doc.photonengine.com/realtime/current/gameplay/custom-properties_

# Custom Properties

Custom Properties are essentially Hashtables for your own key-value entries.

These Custom Properties are either related to the room or one of the players.

The values are synced and cached on the clients, so you don't have to fetch them before use.

Changes are pushed to the others by `SetCustomProperties()`.

How is this useful? Typically, rooms and players have some attributes that are not related to a GameObject:

The current map played in a session or the look of a player's character.

Those can be sent via `OpRaiseEvent` of course but it is often more convenient to use Custom Properties.

To set Custom Properties for a Player, use `Player.SetCustomProperties(Hashtable propsToSet)` and include the key-values to add or update.

A shortcut to the local player object is: `RealtimeClient.LocalPlayer`.

Similarly, use `RealtimeClient.CurrentRoom.SetCustomProperties(Hashtable propsToSet)` to update the room you are in.

All updates take a moment to distribute but all clients will update `CurrentRoom.CustomProperties` and `Player.CustomProperties` accordingly.

The callbacks `OnRoomPropertiesUpdate(Hashtable propertiesThatChanged)` or `OnPlayerPropertiesUpdate(Player targetPlayer, Hashtable changedProps)` are used respectively new updates arrived.

For convenience can also set properties when you get into a room. Joining a room allows you to define your own Custom Player Properties and creating a room also allows setting the Custom Room Properties.

This is especially useful because room properties can be used for matchmaking.

There is a `JoinRandomRoom()` overload which uses a properties-hashtable to filter acceptable rooms for joining.

When you create a room, make sure to define which room properties are available for filtering in the lobby by setting `RoomOptions.CustomRoomPropertiesForLobby` accordingly.

### Check And Swap for Properties (CAS)

When you use `SetCustomProperties`, the server usually accepts new values from any client, which can be tricky in some situations.

For example, a property could be used to store who picked up a unique item in a room.

So the key for the property would be the item and the value defines who picked it up.

Any client can set the property to his actorNumber anytime.

If all do it at about the same time, the _last_`SetCustomProperties` call will win the item (set the final value).

That's counter-intuitive and probably not what you want.

`SetCustomProperties` has an optional `expectedProperties` parameter, which can be used as condition.

With `expectedProperties`, the server will only update the properties, if its current key-values match the ones in `expectedProperties`.

Updates with outdated `expectedProperties` will be ignored (the clients get an error as a result, others won't notice the failed update).

In our example, the `expectedProperties` could contain the current owner from which you take the unique item.

Even if everyone tries to take the item, only the first will succeed, because every other update request will contain an outdated owner in the `expectedProperties`.

Using the `expectedProperties` as a condition in `SetCustomProperties`, is called Check and Swap (CAS).

It is useful to avoid concurrency issues but can also be used in other creative ways.

Note: As `SetCustomProperties` might fail with CAS, all clients update their custom properties by server-sent events only.
This includes the client which attempts to set new values.This is a different timing, compared to setting values without CAS.

You should know that initializing (i.e. first time creating a new property) using CAS is not supported.

Also, currently, there is no callback for SetProperties failures with CAS.

If you want to get notified about CAS failures, here is an example code to add to your MonoBehaviour:

This does not replace the properties update callbacks (`IInRoomCallbacks.OnPlayerPropertiesUpdate` and `IInRoomCallbacks.OnRoomPropertiesUpdate`) which should be triggered in case of success only, either with CAS or without it.

C#

```csharp
    private void OnEnable()
    {
        client.OpResponseReceived += NetworkingClientOnOpResponseReceived;
    }
    private void OnDisable()
    {
        client.OpResponseReceived -= NetworkingClientOnOpResponseReceived;
    }
    private void NetworkingClientOnOpResponseReceived(OperationResponse opResponse)
    {
        if (opResponse.OperationCode == OperationCode.SetProperties &&
            opResponse.ReturnCode == ErrorCode.InvalidOperation)
        {
            // CAS failure
        }
    }

```

### Properties Synchronization

Properties synchronize 'via the server' by default, meaning:

Setting properties will not take effect on the sender/setter client immediately.

Instead, the sender/setter client of the properties waits for the server event `PropertiesChanged` to apply/set changes locally.

So you need to wait until `OnPlayerPropertiesUpdate` or `OnRoomPropertiesUpdate` callbacks are triggered for the local client in order to access them.

The reason behind this is that properties can easily go out of synchronization if we set them locally first and then send the request to do so on the server and for other actors in the room.

The latter might fail and we may end up with properties of the sender/setter client (actor that sets the properties) different locally from what's on the server or on other clients.

If you want to have the old behaviour (set properties locally before sending the request to the server to synchronize them) set `roomOptions.BroadcastPropsChangeToAll` to `false` before creating rooms.

But we highly recommend against doing this.

The client can still cache properties of the local player outside of rooms.

Those properties will be sent when entering rooms.

Also, setting properties in offline mode happens right away.

Besides, by default, local actor's properties are not cleared between rooms, you should do that yourself.

Back to top

- [Check And Swap for Properties (CAS)](#check-and-swap-for-properties-cas)
- [Properties Synchronization](#properties-synchronization)