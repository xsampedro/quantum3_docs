# realtime-intro

_Source: https://doc.photonengine.com/realtime/current/getting-started/realtime-intro_

# Realtime Intro

## Overview

Photon Realtime is our base layer for multiplayer games and higher-level network solutions. It solves problems like authentication, matchmaking and fast communication with a scalable approach.

These pages are primarily the manual for the client-side Realtime APIs but also provide an overview of the workflow and structures involved. The term Photon Realtime relates in the same way to the client-side APIs, the server-side and the workflows and features both sides use.

Photon Realtime SDKs are available in various languages for most of the popular engines and enable cross platform communication for mobile, web, console, standalone and XR. ( [SDK Download Page](https://www.photonengine.com/sdks#realtime)).

If you develop with the Unity engine, consider using [**Fusion**](/fusion) or [**Quantum**](/quantum), which both help synchronize game state and simulation with a rich and deep integration into Unity.

### Users and Authentication

Photon Realtime does not store user information for your app. Instead, clients can authenticate with an external service and ask the Photon server to validate this. Photon servers rely on an "Authentication Provider" to verify the provided authentication values and set the userID for the client.

### Matchmaking

Photon Realtime offers a lean, client-driven approach to matchmaking. On demand, it works well with external matchmaking systems.

The primary goal should be to get players into rooms as fast as possible. In many cases, the server can simply pick a random room to join. Custom properties can be defined to describe rooms and in turn to narrow down which random room matches the player expectations.

Rooms are identified via string, so external matchmaking only needs to make this known to clients to bring players together. With Matchmaking Tokens, the external server can control what a client is able to do in matchmaking.

New rooms get created on demand. Each room can define a max player number and be visible or invisible for random matchmaking and room lists. An invisible room can still be joined by name. Of course, rooms can also be closed, which disables any further join. When all players left, rooms gets cleaned up.

### In Game Communication

Within a room, players can share state and data in several ways:

- Events can be raised with custom code and content. Perfect for frequent updates such as input and positions.
- Custom Room Properties act like a Hashtable which can be updated by either player. This can store general room state for pickup items but also matchmaking tags.
- Custom Player Properties provide a Hashtable per player. Store the players look, equipment and similar values.

## Concepts

### Connections

Photon Realtime always connects clients to a server, as opposed to connecting clients directly. Some SDKs (e.g. Fusion) may establish an additional, independent connection to a host.

The [Photon Cloud is the default option as scalable infrastructure.](/realtime/current/connection-and-authentication/regions) Creating an Exit Games Account is free and without any obligation to purchase any Exit Games products.

The server side is split into three distinct server types:

- **Name Server** s provide region lists and addresses and authenticate the users.
- **Master Server** s handle matchmaking per region.
- **Game Server** s host the actual gameplay in rooms.

Realtime clients are only connected to one server type at any time. This means, matchmaking is only available while the client is not active in a room (and vice versa). Most of the time, the client API will switch servers as needed but it makes sense to know them to understand the workflow of clients.

### Operation, Response, Event

Clients call `Operations` on the server side and get `Operation Responses` for most of those. Aside from Operation Responses, clients also receive `Events`, which are used independently of what the client asked for.

While in a room, the operation `RaiseEvent` is used to pass data to the others, which receive a custom event.

The Realtime client API summarizes common workflows and abstracts-away some of the required operations. For example, the initial Connect-call integrates OpGetRegions and OpAuthenticate. Still, these operation calls may show up in logs and sometimes need to be handled independently.

### Messages

Operations, Responses and Events are sent within `Messages`.

For our reliable UDP protocol, you can choose if these Messages are sequenced or not and if the messages are reliable or not. In TCP and WebSockets, reliability and order are mandatory (by the transport itself).

## Code Samples

Below are a few code samples to give you an idea of how the Realtime API is being used. Consider this an overiew but not a complete, working guide.

Use the Photon Cloud to skip setting up servers. You will need to set an AppId in the code.

Get your AppId from [the Realtime Dashboard](https://dashboard.photonengine.com/) after [**free** signup](https://id.photonengine.com/account/signup).

### Connect

The following code is a class that will connect to the Photon Cloud, if you fill in your appid.

C#

```csharp
using System;
using System.Collections.Generic;
using System.Threading;
using Photon.Realtime;
class GameClass : IConnectionCallbacks
{
    private readonly LoadBalancingClient client = new LoadBalancingClient();
    private bool quit;
    ~GameClass()
    {
        this.client.Disconnect();
        this.client.RemoveCallbackTarget(this);
    }
    public void StartClient()
    {
        this.client.AddCallbackTarget(this);
        this.client.StateChanged += this.OnStateChange;
        this.client.ConnectUsingSettings(new AppSettings() { AppIdRealtime = "<your appid>", FixedRegion = "eu" });
        Thread t = new Thread(this.Loop);
        t.Start();
        Console.WriteLine("Running until key pressed.");
        Console.ReadKey();
        this.quit = true;
    }
    private void Loop(object state)
    {
        while (!this.quit)
        {
            this.client.Service();
            Thread.Sleep(33);
        }
    }
    private void OnStateChange(ClientState arg1, ClientState arg2)
    {
        Console.WriteLine(arg1 + " -> " + arg2);
    }
    // from IConnectionCallbacks:
    public void OnConnectedToMaster()
    {
        Console.WriteLine("OnConnectedToMaster Server: " + this.client.LoadBalancingPeer.ServerIpAddress);
    }
    // ...
}

```

C++

```cpp
class SampleNetworkLogic
{
public:
    SampleNetworkLogic(const ExitGames::Common::JString& appID, const ExitGames::Common::JString& appVersion);
    void connect(void);
    void disconnect(void);
    void run(void);
private:
    ExitGames::LoadBalancing::Client mLoadBalancingClient;
    Listener mListener; // your implementation of the ExitGames::LoadBalancing::Listener interface
    ExitGames::Common::Logger mLogger; // accessed by EGLOG()
};
SampleNetworkLogic::SampleNetworkLogic(const ExitGames::Common::JString& appID, const ExitGames::Common::JString& appVersion)
    : mLoadBalancingClient(mListener, appID, appVersion)
{
}
void SampleNetworkLogic::connect(void)
{
    // connect() is asynchronous - the actual result arrives in the Listener::connectReturn() or the Listener::connectionErrorReturn() callback
    if(!mLoadBalancingClient.connect())
        EGLOG(ExitGames::Common::DebugLevel::ERRORS, L"Could not connect.");
}
int main(void)
{
    static const ExitGames::Common::JString appID = L"<no-app-id>"; // set your app id here
    static const ExitGames::Common::JString appVersion = L"1.0";
    SampleNetworkLogic networkLogic(appID, appVersion);
    networkLogic.connect();

```

#### Connect to self-hosted Photon Server

Connecting to a self-hosted Photon Server does not require an AppId. In best case clients only define a different `AppSettings.Server` with the address of your Name Server.

Read about the differences between Photon Cloud and Photon Server [here](/realtime/current/getting-started/onpremises-or-saas).

### Call Service

The LoadBalancing API is built to integrate well with any game logic. Internally, incoming and outgoing messages are buffered so the game logic can define when incoming messages are handled and when to send outgoing ones.

Calling `client.Service` will dispatch all available incoming messages and send anything outgoing. Calling Service 30x per seconds is common and some apps call it every frame.

While Service() is convenient, there are two methods which can be called instead for more control:

- Call `DispatchIncomingCommands` early in the game loop to dispatch a received response or event. The return bool signals if there are (likely) more events to be dispatched. This can be called in a tight loop.
- Call `SendOutgoingCommands` after the game wrote updates to send whatever outgoing messages were created and buffered. The returned bool signals if there is more to send.

`Service` simply calls `DispatchIncomingCommands` and `SendOutgoingCommands` as long as each returns true.

C#

```csharp
void GameLoop()
{
    while (!shouldExit)
    {
        this.loadBalancingClient.Service();
        Thread.Sleep(50); // wait for a few frames/milliseconds
    }
}

```

C++

```cpp
void SampleNetworkLogic::run(void)
{
    mLoadBalancingClient.service(); // needs to be called regularly!
}
int main(void)
{
    static const ExitGames::Common::JString appID = L"<no-app-id>"; // set your app id here
    static const ExitGames::Common::JString appVersion = L"1.0";
    SampleNetworkLogic networkLogic(appID, appVersion);
    networkLogic.connect();
    while(!shouldExit)
    {
        networkLogic.run();
        SLEEP(100);
    }

```

### Join Random Room

Typically, players should get into rooms as soon and as easy as possible. With the Realtime API, clients can ask the server to join a random room or create a new one if needs be.

This can be as simple as: `OpJoinRandomOrCreateRoom(null, null)` (in C#). The server checks if any rooms accept more players and if that fails, it creates a new one right away. The client gets into a room in both cases. New rooms can be found by the next client looking for a room.

There are plenty of options to refine this workflow. For new rooms, clients can define a maximum player count, the name (e.g. useful for matchmaking via another service), Custom Room Properties and more.

The Custom Room Properties are key-value pairs typically used to store room state. They can also be used as matchmaking filters.

The example below asks the server to find a room with a certain map type and pre-emptively configures the values for a new room, if none is found.

C#

```csharp
// ...
public class MyClient : IConnectionCallbacks, IMatchmakingCallbacks
{
    private LoadBalancingClient loadBalancingClient;
    public MyClient()
    {
        this.loadBalancingClient = new LoadBalancingClient();
        this.loadBalancingClient.AddCallbackTarget(this);
        // TODO: connect, call service, handle more error cases
    }
    // key of our "map type" room property
    private static string MapProperty = "m";
    // room properties available in matchmaking
    private static string[] RoomPropsInLobby = new string[] {"m"};
    // user choice, e.g. types 1 - 9
    private byte selectedMapType = 2;
    void MyJoinRandomOrCreateRoom()
    {
        // custom room properties to use when this client creates a room.
        Hashtable mapSelectionAsProperties = new Hashtable() { { MapProperty, selectedMapType } };

        // if a new room gets created, this sets the map property and makes it available in matchmaking
        RoomOptions propertiesForRoomCreation = new RoomOptions
                                    {
                                        CustomRoomProperties = mapSelectionAsProperties,
                                        CustomRoomPropertiesForLobby = RoomPropsInLobby
                                    };
        EnterRoomParams enterRoomParams = new EnterRoomParams
                                    {
                                        RoomOptions = propertiesForRoomCreation
                                    };
        // this defines the join random filter. rooms must match the key-values in this hashtable
        OpJoinRandomRoomParams joinRoomParams = new OpJoinRandomRoomParams()
                                    {
                                        ExpectedCustomRoomProperties = mapSelectionAsProperties
                                    };
        this.loadBalancingClient.OpJoinRandomOrCreateRoom(joinRoomParams, enterRoomParams);
    }

    public void OnConnectedToMaster()
    {
        this.MyJoinRandomOrCreateRoom();
    }
    void IMatchmakingCallbacks.OnJoinedRoom()
    {
        // ...
    }
    void IMatchmakingCallbacks.OnCreatedRoom()
    {
        // only called when the room got created in addition to OnJoinedRoom()
    }
    // ...

```

Any client can change Custom Room Properties of a room and matchmaking will reflect this automatically. In our example, the players could change the map they play and others can find this room accordingly.

There is an extensive [Matchmaking Guide in our docs](/realtime/current/lobby-and-matchmaking/matchmaking-and-lobby).

### Sending Events

Whatever happens on one client can be sent as an event to update everyone in the same room.

Update your players with positions, your current turn or state values.

Photon will send it as fast as possible (with optional reliability).

- **Send messages/events:** Send any type of data to other players.
- **Player/Room properties:** Photon updates and syncs these, even to players who join later.

C#

```csharp
byte eventCode = 1; // make up event codes at will
Hashtable evData = new Hashtable(); // put your data into a key-value hashtable
this.loadBalancingClient.OpRaiseEvent(eventCode, evData, RaiseEventOptions.Default, SendOptions.SendReliable);

```

C++

```cpp
nByte eventCode = 1; // use distinct event codes to distinguish between different types of events (for example 'move', 'shoot', etc.)
ExitGames::Common::Hashtable evData; // organize your payload data in any way you like as long as it is supported by Photons serialization
bool sendReliable = false; // send something reliable if it has to arrive everywhere
mLoadBalancingClient.opRaiseEvent(sendReliable, evData, eventCode);

```

Your event codes should stay below 200.

Each code defines the type of event and the content receivers can expect.

The event data in the example above is a `Hashtable`.

It can be a `byte\[\]` or any data type supported by Photon's serialization (a `string`, `float\[\]`, etc.).

See [Serialization in Photon](/realtime/current/reference/serialization-in-photon) for more information.

### Receiving Events

Whenever an event is dispatched a handler is called. An example is shown below.

C#

```csharp
using System.Collections.Generic;
using ExitGames.Client.Photon;
using Photon.Realtime;
// we add IOnEventCallback interface implementation
public class MyClient : IConnectionCallbacks, IMatchmakingCallabacks, IOnEventCallback
{
    private LoadBalancingClient loadBalancingClient;
    public MyClient()
    {
        this.loadBalancingClient = new LoadBalancingClient();
        this.loadBalancingClient.AddCallbackTarget(this);
    }
    ~MyClient()
    {
        this.loadBalancingClient.RemoveCallbackTarget(this);
    }
    void IOnEventCallback.OnEvent(EventData photonEvent)
    {
        // we have defined two event codes, let's determine what to do
        switch (photonEvent.Code)
        {
            case 1:
                // do something
                break;
            case 2:
                // do something else
                break;
        }
    }
    // ...

```

C++

```cpp
void NetworkLogic::customEventAction(int playerNr, nByte eventCode, const ExitGames::Common::Object& eventContent)
{
    // logging the string representation of the eventContent can be really useful for debugging, but use with care: for big events this might get expensive
    EGLOG(ExitGames::Common::DebugLevel::ALL, L"an event of type %d from player Nr %d with the following content has just arrived: %ls", eventCode, playerNr, eventContent.toString(true).cstr());
    switch(eventCode)
    {
    case 1:
        {
            // you can access the content as a copy (might be a bit expensive for really big data constructs)
            ExitGames::Common::Hashtable content = ExitGames::Common::ValueObject<ExitGames::Common::Hashtable>(eventContent).getDataCopy();
            // or you access it by address (it will become invalid as soon as this function returns, so (any part of the) data that you need to continue having access to later on needs to be copied)
            ExitGames::Common::Hashtable* pContent = ExitGames::Common::ValueObject<ExitGames::Common::Hashtable>(eventContent).getDataAddress();
        }
        break;
    case 2:
        {
            // of course, the payload does not need to be a Hashtable - how about just sending around for example a plain 64bit integer?
            long long content = ExitGames::Common::ValueObject<long long>(eventContent).getDataCopy();
        }
        break;
    case 3:
        {
            // or an array of floats?
            float* pContent = ExitGames::Common::ValueObject<float*>(eventContent).getDataCopy();
            float** ppContent = ExitGames::Common::ValueObject<float*>(eventContent).getDataAddress();
            short contentElementCount = *ExitGames::Common::ValueObject<float*>(eventContent).getSizes();
            // when calling getDataCopy() on Objects that hold an array as payload, then you must deallocate the copy of the array yourself using deallocateArray()!
            ExitGames::Common::MemoryManagement::deallocateArray(pContent);
        }
        break;
    default:
        {
            // have a look at demo_typeSupport inside the C++ client SDKs for example code on how to send and receive more fancy data types
        }
        break;
    }
}

```

Each event carries the code and data your clients define and send.

Your application knows which content to expect by the code passed (see above).

For an up-to-date list of default event codes look for the event codes constants in your SDK, e.g. within `ExitGames.Client.Photon.LoadBalancing.EventCode` for C#.

### Disconnect

When the application is quitting or when the user logs out do not forget to disconnect.

C#

```csharp
using System.Collections.Generic;
using Photon.Realtime;
public class MyClient : IConnectionCallbacks
{
    private LoadBalancingClient loadBalancingClient;
    public MyClient()
    {
        this.loadBalancingClient = new LoadBalancingClient();
        this.loadBalancingClient.AddCallbackTarget(this);
    }
    ~MyClient()
    {
        this.Disconnect();
        this.loadBalancingClient.RemoveCallbackTarget(this);
    }
    void Disconnect()
    {
        if (this.loadBalancingClient.IsConnected)
        {
            this.loadBalancingClient.Disconnect();
        }
    }
    void IConnectionCallbacks.OnDisconnected(DisconnectCause cause)
    {
        switch (cause)
        {
            // ...

```

C++

```cpp
void SampleNetworkLogic::disconnect(void)
{
    mLoadBalancingClient.disconnect(); // disconnect() is asynchronous - the actual result arrives in the Listener::disconnectReturn() callback
}
int main(void)
{
    static const ExitGames::Common::JString appID = L"<no-app-id>"; // set your app id here
    static const ExitGames::Common::JString appVersion = L"1.0";
    SampleNetworkLogic networkLogic(appID, appVersion);
    networkLogic.connect();
    while(!shouldExit)
    {
        networkLogic.run();
        SLEEP(100);
    }
    networkLogic.disconnect();
}

```

### Custom or Authoritative Server Logic

As is, without authoritative logic, Photon Cloud products already allow for a broad range of game types.

- First Person Shooters
- Racing Games
- Minecraft type of games
- Casual real-time games
- ...

Use [Photon Server](/server/current/getting-started/photon-server-in-5min) or [Photon Plugins](/server/current/plugins/manual) to implement your own custom logic.

Back to top

- [Overview](#overview)

  - [Users and Authentication](#users-and-authentication)
  - [Matchmaking](#matchmaking)
  - [In Game Communication](#in-game-communication)

- [Concepts](#concepts)

  - [Connections](#connections)
  - [Operation, Response, Event](#operation-response-event)
  - [Messages](#messages)

- [Code Samples](#code-samples)
  - [Connect](#connect)
  - [Call Service](#call-service)
  - [Join Random Room](#join-random-room)
  - [Sending Events](#sending-events)
  - [Receiving Events](#receiving-events)
  - [Disconnect](#disconnect)
  - [Custom or Authoritative Server Logic](#custom-or-authoritative-server-logic)