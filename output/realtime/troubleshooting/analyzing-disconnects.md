# analyzing-disconnects

_Source: https://doc.photonengine.com/realtime/current/troubleshooting/analyzing-disconnects_

# Analyzing Disconnects

Even though the internet offers amazing connectivity overall, it is neither entirely stable nor reliable. A lot of hardware and software is involved to connect clients with servers and other devices. The result is that connections stall or fail entirely from time to time.

A low amount of lag, loss and timeouts is to be expected and handled.

This page names the most common cases and helps narrow down the cause for disconnects. There are workarounds for some disconnect causes and best practices to handle the remaining cases.

## Shortlists

### Possible Causes

The following is a non-complete list of interruptions that cause connection issues.

- **Server deployment.** We update and deploy Photon Cloud servers without interrupting games where possible. Sometimes, specific servers become unavailable for some minutes.
- **Hosting issues.** Photon Cloud regions are hosted with backup and alternative routes to them. These kind of issues are usually detected by our monitoring and fixed asap.
- **Cut cables.** We depend on landlines for most of the distances and digging operations sometimes cut them by accident. This could be local, affect your ISP or even a complete region when a backbone cable fails. Even deep sea cables get damaged quite frequently, affecting longer distance connections for a longer duration.
- **DDOS attacks** (and other). Even if malicious traffic can be re-routed, the quality of service during an attack can vary. We got several layers of protection for this and detect them quick. Still, this might be noticed for some hours.
- **Firewalls and web security** software (virus scanners). These may deny some protocols or block connections to specific ports or IPs. Clients may be able to connect to (e.g.) the Name Server but never to a Master (using a different ip and port).
- **App permissions.** The system firewall and permissions can block specific apps from connecing at all.
- **Block lists for ads and tracking.** Such lists randomly include some of our addresses as false positive. Such lists typically affect browser games but blocking via DNS services is also common.
- **VPN.** It may be misconfigured or just not support UDP as effectively as the direct connection would.
- **Different protocols.** We are moving towards using WSS for the initial connection and Authentication. There are several cases, where more than one protocol is used over time and some may be affected by problems.
- **Device incompatibility.** Sometimes two network stacks just do not work with one another properly. One specific mobile may disconnect in just one WiFi network.
- **Many clients in one network.** Some local network hardware or software runs badly when lots of clients run in parallel.
- **Lots of traffic.** Even if clients and servers buffer some of the traffic, connections will fail if the bandwidth need is over the available capabilities. These depend on the device and network. If the server runs out of buffer for sending, this gets logged (SendBufferFull).
- **Performance issues.** The clients need some time to process the networking messages, to acknowledge they still react to messages and to not run out of network buffers.
- **Main Loop pauses.** In some cases Unity will pause the main loop when loading assets. Browsers typically do not run JS/WebGL games in inactive tabs. iOS pauses apps when they go to the background (even to answer calls).
- **Engine bugs.** Photon clients often rely on the engine to implement socket code for target platforms. In some versions, this code failed.

### Debugging Disconnects

Outages aside, individual connection problems are more likely on the player / game side than on the hosting side. The following steps narrow down the cause for problems.

- Check our [Photon Cloud Status Page](https://www.photonengine.com/status)
- Run the project on different device
- Run the project on another network
- Run a demo instead
- Increase logging and look out for relevant messages
- Does it happen to some users more often? Which platform? Where are they?
- What changed in the client?
- [Wireshark](#wireshark)

### Potential fixes

- Update SDK
- Enable Protocol Fallback
- Use WSS port 443 (industries can use ns.photonindustries.io)
- Alternative ports
- RunInBackgroud
- WebGL: Play audio
- Send less
- Fragments (large data)
- Use channel 1++ (large data)
- Less MaxPlayers
- Enable CRC Check

### Avoiding Disconnects

- Test worst case (most players, most activity, longest games)
- PlayerTTL / EmptyRoomTTL
- ReconnectAndRejoin
- Deliberately disconnect when going to background
- Use Interest Group 1+ more and unsubscribe when going to background

### Get Help

If you reach out to us, include this info:

- AppId
- SDK and SDK Version
- Engine and Engine Version
- Platform, device and device OS version (on mobile)
- Send a log file (use high logging level, file preferred)
- Try to identify if the client never connects or fails after connecing
- When did it start
- How common is it (x out of y attempts / percentage)

## Disconnect Reasons

There are additional cases when client can't connect at all (e.g. server unreachable, bad server address, no DNS available, self-hosted server not started, etc.).
In this context, those are not considered disconnects but '(initial) connection failures'.

Client SDKs provide disconnection callbacks and a disconnect cause.

Use those to investigate the unexpected disconnects you are having.

Here we list the major disconnect causes and whether they are caused on the client or the server side.

### Disconnects by Client

- Client-side timeout: no/too late ACKs from the server. See " [Timeout Disconnect](#timeout_disconnect)" for more details.
- Client socket exception (connection loss).
- Client connection fails on receive (buffer full, connection loss). See " [Traffic Issues and Buffer Full](#traffic_issues_and_buffer_full)".
- Client connection fails on send (buffer full, connection loss). See " [Traffic Issues and Buffer Full](#traffic_issues_and_buffer_full)".

Note: If the client connects to localhost, there can be exceptions about failed connections even for UDP (usually UDP just times out).

### Disconnects by Server

- Server-side timeout: no/too late ACKs from client. See " [Timeout Disconnect](#timeout_disconnect)" for more details.
- Server send buffer full (too many messages). See " [Traffic Issues and Buffer Full](#traffic_issues_and_buffer_full)".
- License or Subscription CCU limit hit.

## Timeout Disconnect

Unlike plain UDP, Photon's reliable UDP protocol establishes a connection between server and clients:

Commands within a UDP package have sequence numbers and a flag if they are reliable.

If so, the receiving end has to acknowledge the command.

Reliable commands are repeated in short intervals until an acknowledgement arrives.

If it does not arrive, the connection times out.

Both sides monitor this connection independently from their perspective.

Both sides have their rules to decide if the other is still available.

If a timeout is detected, a disconnect happens on that side of the connection.

As soon as one side thinks the other side does not respond anymore, no message is sent to it.

This is why timeout disconnects are one sided and not synchronous.

The timeout disconnect is the most frequent issue, aside from problems to connect "at all".

There is no single point of failure when you run into frequent timeouts but there are a few common scenarios that cause issues and some ways to fix them.

Here is a quick checklist:

- Is the game [Running in Background](/realtime/current/troubleshooting/known-issues#running_in_background).
- A Unity app is not always running the main loop, [when loading](#unity)
- Check if you can reproduce the issue on other hardware and on another network.


See " [Try Another Connection](#try_another_connection)".
- Check the amount of data you are sending.


If there are spikes or if your messages/sec rate is very high, this can affect the connection quality.


Read " [Send Less](#send_less)"
- You can adjust the number and timing of resends.


See " [Tweak Resends](#tweak_resends)".
- If you want to debug your game using breakpoints and all, [read this](/server/current/operations/debugging).

## Traffic Issues and Buffer Full

Photon servers and clients usually buffer some commands before they are actually put into a package and sent via the internet.

This allows us to aggregate multiple commands into (fewer) packages.

If some either produces a lot of commands (e.g. by sending lots of big events), then the buffers might run out.

Filling buffers will also cause additional Lag:

You will notice that events take longer to arrive on the other side.

Operation responses are not as quick as usual.

Read " [Send Less](#send_less)".

## First Aid

### Check The Logs

This is the first check you need to do.

All clients have some callback to provide log messages about internal state changes and issues.

You should log these messages and access them in case of problems.

You can usually increase the logging to some degree, if nothing useful shows up.

Check the API reference how to do this.

If you customized the server, check the logs there.

### Enable The SupportLogger

The `SupportLogger` is a tool that logs the most commonly needed info to debug problems with Photon, like the (abbreviated)

AppId, version, region, server IPs and some callbacks.

For Unity, the SupportLogger is a MonoBehaviour.

When not using PUN, you can add this component to any GameObject. In your code, find the component (or reference it) to assign

the `LoadBalancingClient` as `Client`. Call `DontDestroyOnLoad()` for the GameObject, if you switch scenes.

Outside of Unity, the SupportLogger is a regular class.

Instantiate it and set the LoadBalancingClient, to make it register for callbacks.

The `Debug.Log` method(s) get mapped to `System.Diagnostics.Debug` respectively.

### Try Another Project

All client SDKs for Photon include some demos.

Use one of those on your target platform.

If the demo fails too, an issue with the connection is more likely.

### Try Another Server or Region

Using the Photon Cloud, you can also use another region easily.

Hosting yourself? Prefer physical over virtual machines.

Test minimum lag (round-trip time) with a client near the server (but not on the same machine or network).

Think about adding servers close to your customers.

### Try Another Connection

In some cases, specific hardware can make the connection fail.

Try another WiFi, router, etc.

Check if another device runs better.

### Try Alternative Ports

As there are some cases where our default UDP port range (5055 to 5058) is blocked, we support an alternative range in all Photon Cloud deployments:

Instead of using 5055 to 5058, the ports start at 27000.

Switching to these "Alternative Ports" is easy.

In the Realtime API you can assign `LoadBalancingClient.ServerPortOverrides = PhotonPortDefinition.AlternativeUdpPorts` before you connect.

The Name Server has port 27000 (was 5058), the Master Server 27001 (was 5055) and the Game Server becomes 27002 (was 5056).

### Enable CRC Checks

Sometimes, packages get corrupted on the way between client and server.

This is more likely when a router or network is especially busy.

Some hardware or software is outright buggy corruption might happen anytime.

Photon has an optional CRC Check per package.

As this takes some performance, we didn't activate this by default.

You enable CRC Checks in the client but the server will also send a CRC when you do.

C#

```csharp
loadBalancingClient.LoadBalancingPeer.CrcEnabled = true

```

Photon clients track how many packages get dropped due to enabled CRC checks.

Check:

C#

```csharp
LoadBalancingPeer.PacketLossByCrc

```

## Fine Tuning

### Check Traffic Stats

On some client platforms, you can enable `Traffic Statistics` directly in Photon.

Those track various vital performance indicators and can be logged easily.

In C#, the Traffic Stats are available in the LoadBalancingPeer class as `TrafficStatsGameLevel` property.

This provides an overview of the most interesting values.

As example, use `TrafficStatsGameLevel.LongestDeltaBetweenDispatching` to check the longest time between to consecutive `DispatchIncomginCommands` calls.

If this time is more than a few milliseconds, you might have some local lag.

Check `LongestDeltaBetweenSending` to make sure your client is frequently sending.

The `TrafficStatsIncoming` and `TrafficStatsOutgoing` properties provide more statistics for in- and outgoing bytes, commands and packages.

### Tweak Resends

C#/.Net Photon library has two properties which allow you to tweak the resend timing:

#### QuickResendAttempts

The `LoadBalancingPeer.QuickResendAttempts` speed up repeats of reliable commands that did not get acknowledged by the receiving end.

The result is a bit more traffic for a shorter delays if some message got dropped.

#### SentCountAllowance

By default, Photon clients send each reliable command up to 6 times.

If there is no ACK for it after the 5th re-send, the connection is shut down.

The `LoadBalancingPeer.SentCountAllowance` defines how often the client will repeat an individual, reliable message.

If the client repeats faster, it should also repeat more often.

In some cases, you see a good effect when setting `QuickResendAttempts` to 3 and `SentCountAllowance` to 7.

More repeats don't guarantee a better connection though and definitely allow longer delays.

### Check Resent Reliable Commands

You should begin to monitor `ResentReliableCommands`.

This counter goes up for each resend of a reliable command (because the acknowledgement from the server didn't arrive in time).

C#

```csharp
LoadBalancingPeer.ResentReliableCommands

```

If this value goes through the roof, the connection is unstable and UDP packets don't get through properly (in either direction).

### Send Less

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


### Try Lower MTU

With a setting on the client-side, you can force server and client to use an even smaller maximum package size than usual.

Lowering the MTU means you need more packages to send some messages but if nothing else helped, it makes sense to try this.

The results of this are unverified and we would like to hear from you if this improved things.

C#

```csharp
loadBalancingClient.LoadBalancingPeer.MaximumTransferUnit = 520;

```

## Tools

### Wireshark

When debugging issues, we usually rely on our client and server logs. In rare cases, it makes sense to look at the actual network traffic.

Wireshark is a protocol analyzer and logger and extremely useful to find out what is actually happening on the network layer of your game.

On start Wireshark should show a list of network interfaces and the "Capture Filter" input field.

![](/docs/img/Wireshark-03-InterfaceSettings.png)
Wireshark 3.6 Interfaces, simple capture filter and main menu.


A small graph per network device indicates current network activity. Click the one device Wireshark should capture.

The "Capture Filter" is used to narrow down which traffic will be captured. Restrict this to the ports used by your Photon clients.

Use `portrange 5055-5058 or portrange 27000-27003` to cover our default UDP port ranges.

To capture WebSocket traffic, you may want to explicitly use the Photon specific port ranges: `portrange 9090-9093 or portrange 19090-19093`. You may have to configure the client to use these ports. If the clients run WS on port 80 or WSS on port 443, Wireshark will pick up a lot of other network traffic, which can be distracting.

Click the blue menu icon on the top left to "Start capturing packages", reproduce the issue and stop the capture (third toolbar button). Only after you stopped capturing data, Wireshark allows you to save it as file.

Mail the `.pcapng` and other files to us and we take a look.

In best case, you also include a description of what you did, if the error happens regularly, how often and when it happened in this case (there are timestamps in the log).

Attach a client console log, too.

### Clumsy

Clumsy is a minimalistic tool to mess with the quality of network connections. Is it very useful to test the game behavior with more lag and if connection loss is handled properly.

![](/docs/img/Tools-Clumsy.jpg)
Clumsy Window


Clumsy uses a protocol and port filter to define which traffic gets affected. To affect only the gameplay on a game server (connected via UDP), set `udp.DstPort == 5056 or udp.SrcPort == 5056 or udp.DstPort == 27002 or udp.SrcPort == 27002`.

Lag and Loss are useful effects to apply in moderation. All effects may apply to outgoing _and_ incoming traffic. Additional lag of 50ms to each will increase the rountrip times by 100ms.

Set the "Drop" values to 100% and click "start" to test a complete connection loss.

## Platform Specific Info

### Unity

PUN automatically keeps the connection for you. To do so, it calls `DispatchIncomingCommands` and `SendOutgoingCommands` on Unity's main loop.

However, Unity won't call `Update` while it's loading scenes and assets or while you drag a standalone-player's window.

To keep the connection while loading scenes, you should set `PhotonNetwork.IsMessageQueueRunning = false`.

Pausing the message queue has two effects:

- A background thread will be used to call `SendOutgoingCommands` while `Update` is not called.


This keeps the connection alive, sending acknowledgements only but no events or operations (RPCs or sync updates).


Incoming data is not executed by this thread.
- All incoming updates are queued. Neither RPCs are called, nor are observed objects updated.


While you change the level, this avoids calling RPCs in the previous one.

If you use our Photon Unity SDK, you probably do the `Service` calls in some MonoBehaviour `Update` method.

To make sure Photon client's `SendOutgoingCommands` is called while you load scenes, implement a background thread.

This thread should pause 100 or 200 ms between each call, so it does not take away all performance.

## Recover From Unexpected Disconnects

Disconnects will happen, they can be reduced but they can't be avoided.

So it's better to implement a recovery routine for when those unexpected disconnects occur especially mid-game.

### When To Reconnect

First you need to make sure that the disconnect cause can be recovered from.

Some disconnects may be due to issues that cannot be resolved or bypassed by a simple reconnect.

Instead those cases should be treated separately and handled case by case.

### Quick Rejoin (ReconnectAndRejoin)

A "Quick Rejoin" can be used when a client got disconnected while playing (in a session/room). In that case, the Photon client uses an existing Photon Token, room name and game server address and can get back into the room, even if the server did not notice the absence of this client (the player still being active).

In C# SDKs, this is done using `LoadBalancingClient.ReconnectAndRejoin()`.

Check the return value of this method to make sure the quick rejoin process is initiated.

In order for the reconnect and rejoin to succeed, the room needs to have PlayerTTL != 0.

But this is not a guarantee that the rejoin will work.

If the reconnection successful, rejoin can fail with one of the following errors:

- GameDoesNotExist (32758): the room was removed from the server while disconnected. This probably means that you were the last actor leaving the room when disconnected and that 0 <= EmptyRoomTTL < PlayerTTL or PlayerTTL < 0 <= EmptyRoomTTL.
- JoinFailedWithRejoinerNotFound (32748): the actor was removed from the room while disconnected. This probably means that PlayerTTL is too short and expired, we suggest at least a value of 12000 milliseconds to allow a quick rejoin.
- PluginReportedError (32752): this probably means that you use webhooks and that PathCreate returns ResultCode other than 0.
- JoinFailedFoundActiveJoiner (32746): this is very unlikely to happen but it may. It means that another client using the same UserId but a different Photon token joined the room while a client was disconnected.

You can catch these in the `OnJoinRoomFailed` callback.

### Reconnect

If the client got disconnected outside of a room or if quick rejoin failed (`ReconnectAndRejoin` returned false) you could still do a Reconnect only.

The client will reconnect to the master server and reuse the cached authentication token there.

In C# SDKs, this is done using `LoadBalancingClient.ReconnectToMaster()`.

Check the return value of this method to make sure the quick rejoin process is initiated.

It could be useful in some cases to add:

- check if connectivity is working as expected (internet connection available, servers/network reachable, services status)
- reconnect attempts counter: max. retries
- backoff timer between retries

### Sample (C#)

C#

```csharp
using System;
using Photon.Realtime;
public class RecoverFromUnexpectedDisconnectSample : IConnectionCallbacks
{
    private LoadBalancingClient loadBalancingClient;
    private AppSettings appSettings;
    public RecoverFromUnexpectedDisconnectSample(LoadBalancingClient loadBalancingClient, AppSettings appSettings)
    {
        this.loadBalancingClient = loadBalancingClient;
        this.appSettings = appSettings;
        this.loadBalancingClient.AddCallbackTarget(this);
    }
    ~RecoverFromUnexpectedDisconnectSample()
    {
        this.loadBalancingClient.RemoveCallbackTarget(this);
    }
    void IConnectionCallbacks.OnDisconnected(DisconnectCause cause)
    {
        if (this.CanRecoverFromDisconnect(cause))
        {
            this.Recover();
        }
    }
    private bool CanRecoverFromDisconnect(DisconnectCause cause)
    {
        switch (cause)
        {
            // the list here may be non exhaustive and is subject to review
            case DisconnectCause.Exception:
            case DisconnectCause.ServerTimeout:
            case DisconnectCause.ClientTimeout:
            case DisconnectCause.DisconnectByServerLogic:
            case DisconnectCause.DisconnectByServerReasonUnknown:
                return true;
        }
        return false;
    }
    private void Recover()
    {
        if (!loadBalancingClient.ReconnectAndRejoin())
        {
            Debug.LogError("ReconnectAndRejoin failed, trying Reconnect");
            if (!loadBalancingClient.ReconnectToMaster())
            {
                Debug.LogError("Reconnect failed, trying ConnectUsingSettings");
                if (!loadBalancingClient.ConnectUsingSettings(appSettings))
                {
                    Debug.LogError("ConnectUsingSettings failed");
                }
            }
        }
    }
    #region Unused Methods
    void IConnectionCallbacks.OnConnected()
    {
    }
    void IConnectionCallbacks.OnConnectedToMaster()
    {
    }
    void IConnectionCallbacks.OnRegionListReceived(RegionHandler regionHandler)
    {
    }
    void IConnectionCallbacks.OnCustomAuthenticationResponse(Dictionary<string, object> data)
    {
    }
    void IConnectionCallbacks.OnCustomAuthenticationFailed(string debugMessage)
    {
    }
    #endregion
}

```

Back to top

- [Shortlists](#shortlists)

  - [Possible Causes](#possible-causes)
  - [Debugging Disconnects](#debugging-disconnects)
  - [Potential fixes](#potential-fixes)
  - [Avoiding Disconnects](#avoiding-disconnects)
  - [Get Help](#get-help)

- [Disconnect Reasons](#disconnect-reasons)

  - [Disconnects by Client](#disconnects-by-client)
  - [Disconnects by Server](#disconnects-by-server)

- [Timeout Disconnect](#timeout-disconnect)
- [Traffic Issues and Buffer Full](#traffic-issues-and-buffer-full)
- [First Aid](#first-aid)

  - [Check The Logs](#check-the-logs)
  - [Enable The SupportLogger](#enable-the-supportlogger)
  - [Try Another Project](#try-another-project)
  - [Try Another Server or Region](#try-another-server-or-region)
  - [Try Another Connection](#try-another-connection)
  - [Try Alternative Ports](#try-alternative-ports)
  - [Enable CRC Checks](#enable-crc-checks)

- [Fine Tuning](#fine-tuning)

  - [Check Traffic Stats](#check-traffic-stats)
  - [Tweak Resends](#tweak-resends)
  - [Check Resent Reliable Commands](#check-resent-reliable-commands)
  - [Send Less](#send-less)
  - [Try Lower MTU](#try-lower-mtu)

- [Tools](#tools)

  - [Wireshark](#wireshark)
  - [Clumsy](#clumsy)

- [Platform Specific Info](#platform-specific-info)

  - [Unity](#unity)

- [Recover From Unexpected Disconnects](#recover-from-unexpected-disconnects)
  - [When To Reconnect](#when-to-reconnect)
  - [Quick Rejoin (ReconnectAndRejoin)](#quick-rejoin-reconnectandrejoin)
  - [Reconnect](#reconnect)
  - [Sample (C#)](#sample-c)