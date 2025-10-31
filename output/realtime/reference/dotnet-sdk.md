# dotnet-sdk

_Source: https://doc.photonengine.com/realtime/current/reference/dotnet-sdk_

# .NET Platform SDK

This is the documentation and reference for the Photon Client Library for .NET and Unity.

## Overview

Photon is a development framework to build real-time multiplayer games and applications for various platforms.

It consists of a Server SDK and Client SDKs for several platforms.

Photon provides a low-latency communication-layer based on UDP (or alternatively TCP).

It enables reliable and unreliable transfer of data in "commands".

On top of this, an operation- and event-framework is established to ease development of your own games.

Photon Cloud is a hosted service for your Photon games.

To use it, register [here](https://id.photonengine.com/account/signup) and use the LoadBalancing API and demos to start your development.

### Photon Workflow

To get an impression of how to work on the client, we will use the server's LoadBalancing API logic.

This application defines rooms which are created when users try to join them.

Each user in a room becomes an actor with her own number.

A simplified client workflow looks like this:

- create a `LoadBalancingClient` instance
- from now on: regularly call `Service` to get events and send commands (e.g. ten times a second)
- call one of the `Connect` methods to connect the server
- wait until the library calls `OnStatusChanged`
- the returned status int should equal `StatusCode.Connect`
- call `OpJoinRoom` to get into a game
- wait until the library calls `OnOperationResponse` with `OperationCode.Join`
- send data in the game by calling `OpRaiseEvent`
- receive events in `OnEvent`
- The LoadBalancing API defines several useful events for common situations: Someone joins or leaves the room.
- In LoadBalancing API, events created by calling `OpRaiseEvent` will be received by others in the same room in this method.
- when you are done: call `LoadBalancingClient.OpLeave` to quit/leave the game
- wait for "leave" return in `OnOperationResponse` with `OperationCode.Leave`
- disconnect with `Disconnect`
- check "disconnect" return in `OnStatusChanged` with statusCode: `StatusCode.Disconnect`


Combined with the server's LoadBalancing application, this simple workflow would allow you to use rooms and send your game's events.


The methods used could be broken down into three layers:
- Low Level: `Service, Connect, Disconnect` and the `OnStatusChanged` are directly referring to the connection to the server.


This level works with UDP/TCP packets which transport commands (which in turn carry your operations).


It keeps your connection alive and organizes your RPC calls and events into packages.
- Logic Level: Operations, results and events make up the logical level in Photon.


Any operation defined on the server (think RPC call) and can have a result.


Events are incoming from the server and update the client with some data.
- Application Level: Made up by a specific application and its features.


In this case we use the operations and logic of the LoadBalancing API.


In this specific case, we have rooms and actors and more.

The LoadBalancingClient is matching the server side implementation and wraps it up for you.

You don't have to manage the low level communication in most cases.

However, it makes sense to know that everything that goes from client to server (and the other way round) is put into "commands".

Internally, commands are also used to establish and keep the connection between client and server alive (without carrying additional data).

All methods that are operations (RPC calls) are prefixed with "Op" to tell them apart from anything else.

Other server side applications (like MMO or your own) will define different operations.

These will have different parameters and return values.

These operations are not part of the client library but can be implemented by calling `OpCustom`.

The interface IPhotonPeerListener must be implemented for callbacks.

They are:

- `OnStatusChanged` is for peer state-changes (connect, disconnect, errors, compare with `StatusCode` Enumeration)
- `OnOperationResponse` is the callback for operations (join, leave, etc.)
- `OnEvent` as callback for events coming in
- `DebugReturn` as callback to debug output (less frequently used by release builds)


The following properties in PhotonPeer are of special interest:
- `TimePingInterval` sets the time between ping-operations
- `RoundTripTime` of reliable operations to the server and back
- `RoundTripTimeVariance` shows the variability of the roundtrip time
- `ServerTimeInMilliSeconds` is the continuously approximated server's time

### Operations

Operation is our term for remote procedure calls (RPC) on Photon.

This in turn can be described as methods that are implemented on the server-side and called by clients.

As any method, they have parameters and return values.

The Photon development framework takes care of getting your RPC calls from clients to server and results back.

Server-side, operations are part of an application running on top of Photon.

The default application provided by Exit Games is called "LoadBalancing API".

The `LoadBalancingPeer` class extends the `PhotonPeer` by methods for each of the LoadBalancing API Operations.

The `LoadBalancingClient` is a wrapper of the `LoadBalancingPeer` class.

Examples for LaodBalancing API Operations are "Join" and "RaiseEvent".

On the client side, they can be found in the `LoadBalancingClient` class as methods: `OpJoinRoom` and `OpRaiseEvent`.

They can be used right away with the default implementation of Photon and the LoadBalancing API.

#### Custom Operations

Photon is extendable with features that are specific to your game.

You could persist world states or double check information from the clients.

Any operation that is not in Lite, LoadBalancing API or the MMO application logic is called Custom Operation.

Creating those is primarily a server-side task, of course, but the clients have to use new functions / operations of the server.

So Operations are methods that can be called from the client side.

They can have any number of parameters and any name.

To preserve traffic, we assign byte-codes for every operation and each parameter.

The definition is done server side.

Each Operation has its own, unique number to identify it, known as the operation code (`OperationCode`).

An operation class defines the expected parameters and assigns a parameter code for each.

With this definition, the client side only has to fill in the values and let the server know the `OperationCode` of the Operation.

Photon uses Dictionaries to aggregate parameters for operation requests, responses and events.

Use `OpCustom` to call any operation, providing the parameters in a Dictionary.

Client side, `OperationCode` and parameter-codes are currently of type byte (to minimize overhead).

They need to match the definition of the server side to successfully call your operation.

Recommended for further reading:

- [Calling Operations](/server/v4/reference/calling-operations)

### Events

Unlike operations, events are "messages" that are rarely triggered by the client that receives them.

Events come from outside: the server or other clients.

They are created as side effect of operations (e.g. when you join a room) or raised as main purpose of the operation Raise Event.

Most events carry some form of data but in rare cases the type of event itself is the message.

Events are (once more) Dictionaries with arbitrary content.

In the "top-level" of an event, bytes are used as keys for values.

The values can be of any serializable type.

The LoadBalancing API, e.g., uses a `Hashtable` for custom event content in its operation `RaiseEvent`.

Recommended for further reading:

- [Receiver Groups](/realtime/current/gameplay/cached-events)

- [Serialization in Photon](/realtime/current/reference/serialization-in-photon)

### Fragmentation and Channels

#### Fragmentation

Bigger data chunks of data won't fit into a single package, so they are fragmented and reassembled automatically.

Depending on the data size, a single operation or event can be made up of multiple packages.

Be aware that this might stall other commands.

Call `Service` or `SendOutgoingCommands` more often than absolutely necessary.

You should check that `PhotonPeer.QueuedOutgoingCommands` is becoming zero regularly to make sure everything gets out.

You can also check the debug output for "UDP package is full", which can happen from time to time but should not happen permanently.

#### Maximum Transfer Unit

The maximum size for any UDP package can be configured by setting `PhotonPeer.MaximumTransferUnit` Property.

By default, this is 1200 bytes.

Some routers will fragment even this UDP package size.

If you don't need bigger sizes, go for 512 bytes per package, which is more overhead per command but potentially safer.

This setting is ignored by TCP connections, which negotiate their MTU internally.

#### Sequencing

The sequencing of the protocol makes sure that any receiving client will Dispatch your actions in the order you sent them.

Unreliable data is considered replaceable and can be lost.

Reliable events and operations will be repeated several times if needed but they will all be dispatched in order without gaps.

Unreliable actions are also related to the last reliable action and not dispatched before that reliable data was dispatched first.

This can be useful, if the events are related to each other.

Example: Your FPS sends out unreliable movement updates and reliable chat messages.

A lost package with movement updates would be left out as the next movement update is coming fast.

On the receiving end, this would maybe show as a small jump.

If a package with a chat message is lost, this is repeated and would introduce lag, even to all movement updates after the message was created.

In this case, the data is unrelated and should be put into different channels.

#### Channels

The .NET clients and server are now supporting "channels".

This allows you to separate information into multiple channels, each being sequenced independently.

This means, that Events of one channel will not be stalled because events of another channel are not available.

By default an PhotonPeer has two channels and channel zero is the default to send operations.

The operations join and leave are always sent in channel zero (for simplicity).

There is a "background" channel 255 used internally for connect and disconnect messages.

This is ignored for the channel count.

Channels are prioritized: the lowest channel number is put into a UDP package first.

Data in a higher channel might be sent later when a UDP package is already full.

Example: The chat messages can now be sent in channel one, while movement is sent in channel zero.

They are not related and if a chat message is delayed, it will no longer affect movement in channel zero.

Also, channel zero has higher priority and is more likely to be sent (in case packages get filled up).

### Using TCP

A PhotonPeer could be instanced with TCP as underlying protocol if necessary.

This is not best practice but some platforms don't support UDP sockets.

This is why Silverlight (e.g.) uses TCP in all cases.

The Photon Client API is the same for both protocols but there are some differences in what goes on under the hood.

Everything sent over TCP is always reliable, even if you call your operations as unreliable!

If you use only TCP clients simply send any operation unreliable.

It saves some work (and traffic) in the underlying protocols.

If you have TCP and UDP clients anything you send between the TCP clients will always be transferred reliable.

But as you communicate with some clients that use UDP these will get your events reliable or unreliable.

Example:

A Silverlight client might send unreliable movement updates in channel # 1.

This will be sent via TCP, which makes it reliable.

Photon however also has connections with UDP clients (like a 3D downloadable game client).

It will use your reliable / unreliable settings to forward your movement updates accordingly.

### Network Simulation

During development, most tests will be done in a local network.

Once released, the clients will communicate through the internet, which has a higher delay per message and in some cases even drops messages entirely.

To prepare a game for real-life conditions, the Photon client libraries let you simulate some effects of

internet-communication: lag, jitter and packet loss.

- Lag / Latency: a more or less constant delay of messages between client and server.


Either direction can be affected in a


different way but usually the values are close to another.


Affects the roundtrip time.
- Jitter: Is randomizes the Lag in the simulation.


This affects the variance of the roundtrip time.


Udp packages can get out


of order this way, which also is simulated.


The new lag will be: Lag + \[-JitterValue..+JitterValue\].


This keeps the mean Lag


at the setting and some packages are actually faster than the Lag value implies.
- Packet Loss: UPD packages can become lost.


In the Photon protocol, commands that are flagged as reliable will be repeated while other commands (operations) might get lost this way.

The lag simulation is running in its own Thread which tries to meet delays defined in the settings.

In most cases, they can be

met but actual delays will have a variance of up to +/- 20ms.

#### Using Network Simulation

By default, Network Simulation is turned off.

It can be turned on by setting `PhotonPeer.IsSimulationEnabled` and the settings are aggregated into a NetworkSimulationSet, also accessible by the peer class (e.g. the LoadBalancingClient).

Code Sample:

C#

```csharp
//Activate / Deactivate:
this.peer.IsSimulationEnabled = true;
//Raise Incoming Lag:
this.peer.NetworkSimulationSettings.IncomingLag = 300; //default is 100ms
//add 10% of outgoing loss:
this.peer.NetworkSimulationSettings.OutgoingLossPercentage = 10; //default is 1
//this property counts the actual simulated loss:
this.peer.NetworkSimulationSettings.LostPackagesOut;

```

### The Photon Server

The Photon Server is the central hub for communication for all your clients.

It is a service that can be run on any Windows machine, handling the UDP and TCP connections of clients and hosting a .NET runtime layer with your own business logic, called application.

The Photon Server SDK includes several applications in source and pre-built.

You can run them out of the box or develop your own server logic.

Get the Photon Server SDK [here](https://www.photonengine.com/sdks#server-sdkserverserver)

### LoadBalancing API

The LoadBalancing API is the example application for room-based games on Photon and (hopefully) a flexible basis for your own games.

It offers rooms, joining and leaving them, sending events to the other players in a room and handles properties.

The LoadBalancing API is tightly integrated with the client libraries and used as example throughout most documentation.

#### Properties on Photon

The LoadBalancing API implements a general purpose mechanism to set and fetch key/value pairs on the server side (in memory).

They are associated to a room/game or a player within a room and can be fetched or updated by anyone in that game.

Each entry in the properties Hashtable is considered a separate property and can be overwritten independently.

The value of a property can be of any serializable datatype.

The keys must be either of type string or byte.

Bytes are preferred, as they mean the less overhead.

To avoid confusion, don't mix string and byte as key-types.

Mixed types of keys, require separate requests to fetch them.

Property broadcasting and events Property changes in a game can be "broadcasted", which triggers events for the other players to update them.

The player who changed the property does not get the update (again).

Any change that uses the broadcast option will trigger a property update event.

This event carries the changed properties (only), which changed the properties and where the properties belong to.

Your clients need to "merge" the changes (if properties are cached at all).

Properties can be set by these methods:

- `LoadBalancingClient.OpSetPropertiesOfActor` sets actor properties.
- `LoadBalancingClient.OpSetPropertiesOfRoom` sets room properties.
- `LoadBalancingClient.OpJoinOrCreateRoom` sets initial room properties if the room did not exist before and sets joining actor properties.
- `LoadBalancingClient.OpCreateRoom` sets initial room properties and sets joining actor properties.
- `LoadBalancingClient.OpJoinRoom` sets joining actor properties.

Any change that uses the broadcast option will trigger a property update event `EventCode.PropertiesChanged`.

This event carries the properties as value of key `ParameterCode.Properties`.

Additionally, there is information about who changed the properties in key `ParameterCode.ActorNr`.

The key `ParameterCode.TargetActorNr` will only be available if the property-set belongs to a certain player.

If it's not present, the properties are room properties.

Back to top

- [Overview](#overview)
  - [Photon Workflow](#photon-workflow)
  - [Operations](#operations)
  - [Events](#events)
  - [Fragmentation and Channels](#fragmentation-and-channels)
  - [Using TCP](#using-tcp)
  - [Network Simulation](#network-simulation)
  - [The Photon Server](#the-photon-server)
  - [LoadBalancing API](#loadbalancing-api)