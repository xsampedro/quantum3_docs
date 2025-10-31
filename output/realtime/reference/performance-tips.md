# performance-tips

_Source: https://doc.photonengine.com/realtime/current/reference/performance-tips_

# Performance Tips

Performance is a vital part for providing a fluid and seamless integration of multiplayer components into your application.

So we assembled a list of tips you should keep in mind when developing with Photon.

### Call Service Regularly

The client libraries are built to send any messages only when the app logic triggers that. This way, the clients can aggregate several operations and avoid network overhead.

To trigger sending any data, a main loop must call `LoadBalancingClient.Service()` or `PhotonPeer.SendOutgoingCommands()` frequently. The bool return value is true if some data is still queued. If so, call SendOutgoingCommands again (but not more than three times in a row).

Service and SendOutgoingCommands also send acknowledgements and pings, which are important to keep a connection alive. You should avoid longer pauses between calling either. Especially, make sure that Service is called despite loading.

If overlooked, this problem is hard to identify and reproduce. The C# library has a ConnectionHandler class, which can help.

To avoid local lag, you could call SendOutgoingCommands after the game loop wrote network updates.

### Updates vs. Traffic

Ramping up the number of updates per second makes a game more fluid and up-to-date. On the other hand, traffic might increase dramatically. Also, random lag and loss can not be avoided, so receivers of updates should always be capable of interpolating important values.

Keep in mind that many operations you call will create events for other players and that it might in fact be faster to send fewer updates per second.

### Optimizing Traffic

You can usually send less to avoid traffic issues.

Doing so has a lot of different approaches:

#### Don't Send More Than What's Needed

Exchange only what's totally necessary.

Send only relevant values and derive as much as you can from them.

Optimize what you send based on the context.

Try to think about what you send and how often.

Non critical data should be either recomputed on the receiving side based on the data synchronized or with what's happening in game instead of forced via synchronization.

Examples:

- In an RTS, you could send "orders" for a bunch of units when they happen.


This is much leaner than sending position, rotation and velocity for each unit ten times a second.


Good read: [1500 archers](https://www.gamasutra.com/view/feature/3094/1500_archers_on_a_288_network_.php).

- In a shooter, send a shot as position and direction.


Bullets generally fly in a straight line, so you don't have to send individual positions every 100 ms.


You can clean up a bullet when it hits anything or after it travelled "so many" units.


- Don't send animations. Usually you can derive all animations from input and actions a player does.


There is a good chance that a sent animation gets delayed and playing it too late usually looks awkward anyways.

- Use delta compression. Send only values when they changes since last time they were sent.


Use interpolation of data to smooth values on the receiving side.


It's preferable over brute force synchronization and will save traffic.


#### Don't Send Too Much

Optimize exchanged types and data structures.

Examples:

- Make use of bytes instead of ints for small ints, make use of ints instead of floats where possible.
- Avoid exchanging strings at all costs and prefer enums/bytes instead.
- Avoid exchanging custom types unless you are totally sure about what get sent.

Use another service to download static or bigger data (e.g. maps).

Photon is not built as content delivery system.

It's often cheaper and easier to maintain to use HTTP-based content systems.

Anything that's bigger than the Maximum Transfer Unit (MTU) will be fragmented and sent as multiple reliable packages (they have to arrive to assemble the full message again).

#### Don't Send Too Often

- Lower the send rate, you should go under 10 if possible.


This depends on your gameplay of course.


This has a major impact on traffic.


You can also use adaptive or dynamic send rate based on the user's activity or the exchanged data, this is also helping a lot.

- Send unreliable when possible.


You can use unreliable messages in most cases if you have to send another update as soon as possible.


Unreliable messages never cause a repeat.


Example: In an FPS, player position can usually be sent unreliable.


### Producing and Consuming Data

Related to the "traffic" topic is the problem of producing only the amount of data that can be consumed on the receiving end.

If performance or frame rate don't keep up with incoming events they are outdated before they are executed.

In the worst case, one side produces so much data that it breaks the receiving end.

Keep an eye on the queue length of your clients while developing.

### Limiting Execution of Unreliable Commands

Even if a client doesn't dispatch incoming messages for a while (e.g. while loading), it will still receive and buffer everything.

Depending on the activity of the other players, a client might have a lot to catch up with.

To keep things lean, a client will automatically cut the unreliable messages to a certain length.

The idea is that you get the latest info faster and missing updates will be replaced by new, up-to-date messages soon.

This limit is set via `LoadbalancingPeer.LimitOfUnreliableCommands` which has a default of 20 (in PUN, too).

### Datagram Size

The content size of datagrams is limited to 1200 bytes by default.

These 1200 bytes include all the overhead from headers (see " [Binary Protocol](/realtime/current/reference/binary-protocol)"), size and type information (see " [Serialization in Photon](/realtime/current/reference/serialization-in-photon)"), so that the number for actual pure payload is somewhat lower.

In fact, even if it varies depending on how data is structured, we can safely assume that pure payload data lower than 1kb can fit into a single datagram.

Operations and events that are bigger than 1200 bytes get fragmented and are sent in multiple commands.

These become _reliable_ automatically and the receiving side can only reassemble and dispatch those bigger data chunks when all fragments are received.

Bigger data "streams" can considerably affect latency as they need to be reassembled from many packages before they are dispatched.

They can be sent in a separate _channel_, so they don't affect the "live" position updates of a (lower) channel number.

## Reduce Allocations With Pooled ByteArraySlice

By default, Photon clients in C# SDKs serializes `byte\[\]` and `ArraySegment<byte>` as `byte\[\]`.

On the receiving side, this allocates a new `byte\[\]` of the same length, which is passed to the `OnEvent` callbacks.

ByteArraySlice is a non-alloc / non-boxing alternative to these options.

`ByteArraySlice` is a wrapper class for a `byte\[\]` very similar to `ArraySegment<byte>`, except that it is a recyclable class.

As a class it can be cast to object (which all Photon messages are cast to) without creating an allocation from boxing.

The fields/properties of `ByteArraySlice` are:

- `Buffer`: The wrapped `byte\[\]` array.
- `Offset`: The starting byte the transport will read from.
- `Count`: The number of bytes that were written past the Offset.

### Serialization Usage

This can be done in of two ways:

#### Acquire from ByteArraySlicePool

C#

```csharp
void Serialization()
{
    // Get a pooled Slice.
    var pool = loadBalancingClient.LoadBalancingPeer.ByteArraySlicePool;
    var slice = pool.Acquire(256);
    // Write your serialization to the byte[] Buffer.
    // Set Count to the number of bytes written.
    slice.Count = MySerialization(slice.Buffer);
    loadBalancingClient.OpRaiseEvent(MSG_ID, slice, opts, sendOpts);
    // The ByteArraySlice that was Acquired is automatically returned to the pool
    // inside of the OpRaiseEvent
}

```

#### Maintain your own ByteArraySlice

C#

```csharp
private ByteArraySlice slice = new ByteArraySlice(new byte[1024]);
void Serialization()
{
    // Write your serialization to the byte[] Buffer.
    // Set Count to the number of bytes written.
    slice.Count = MySerialization(slice.Buffer);
    loadBalancingClient.OpRaiseEvent(MSG_ID, slice, opts, sendOpts);
}

```

### Deserialization Usage

By default `byte\[\]` data is deserialized to `new byte\[x\]`.

We must set `LoadBalancingPeer.UseByteArraySlicePoolForEvents = true` to enable the non-alloc conduit.

Once enabled, we cast incoming objects to `ByteArraySlice` rather than `byte\[\]`.

C#

```csharp
// By default byte arrays arrive as byte[]
// UseByateArraySlicePoolForEvents must be enabled to use this feature
private static void EnableByteArraySlicePooling()
{
    loadBalancingPeer.UseByteArraySlicePoolForEvents = true;
}
private void OnEvent(EventData photonEvent)
{
    // Rather than casting to byte[], we now cast to ByteArraySlice
    ByteArraySlice slice = photonEvent.CustomData as ByteArraySlice;
    // Read in the contents of the byte[] Buffer
    // Your custom deserialization code for byte[] will go here.
    Deserialize(slice.Buffer, slice.Count);
    // Be sure to release the slice back to the pool
    slice.Release();
}

```

## Reusing EventData

The C# clients receive events via `OnEvent(EventData ev)`. By default, each EventData is a new instance, which causes some extra work for the garbage collector.

In many cases, it is easily possible to reuse the EventData and avoid the overhead. This can be enabled via the `PhotonPeer.ReuseEventInstance` setting.

Back to top

- [Call Service Regularly](#call-service-regularly)
- [Updates vs. Traffic](#updates-vs.traffic)
- [Optimizing Traffic](#optimizing-traffic)
- [Producing and Consuming Data](#producing-and-consuming-data)
- [Limiting Execution of Unreliable Commands](#limiting-execution-of-unreliable-commands)
- [Datagram Size](#datagram-size)

- [Reduce Allocations With Pooled ByteArraySlice](#reduce-allocations-with-pooled-bytearrayslice)

- [Serialization Usage](#serialization-usage)
- [Deserialization Usage](#deserialization-usage)

- [Reusing EventData](#reusing-eventdata)