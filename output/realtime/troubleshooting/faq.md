# faq

_Source: https://doc.photonengine.com/realtime/current/troubleshooting/faq_

# Frequently Asked Questions

#### Which Photon product is the right one for me?

The answer to this depends mostly on your project and team.

Generally, we suggest to use either [Fusion](https://www.photonengine.com/slides/fusion) or [Quantum](https://www.photonengine.com/slides/quantum), which are our most advanced client solutions.

For a quick overview, both product sheets contain the product picker "Quadrant":

- [Quadrant in Fusion slides](https://www.photonengine.com/slides/fusion#page-7)
- [Quadrant in Quantum slides](https://www.photonengine.com/slides/quantum#page-6)

Feel free to reach out to us for any questions.

## Photon Cloud

### Is Photon Cloud down?

Our **[Photon Cloud Status Page](https://www.photonengine.com/status)** shows the current and past status per product.

We also announce status updates on our Discord servers.

### Is there a default Photon Cloud region?

Actually, there is no default region. Clients know the Name Server address for the Photon Cloud. They are global and used to provide an up to date region list for the given AppId.

Clients will ping each region and identify the "Best Region", which has the lowest latency.

If none of the regions can be pinged successfully, the first region of the list is used.

#### Is it possible to disable some regions?

Yes.

It works by defining a "Region Allow List" without the "disabled" regions.

Read more about the " [Dashboard Regions Filtering](/realtime/current/connection-and-authentication/regions#region-allowlist)".

### Can we get a list of all Cloud servers / IPs?

Such a list does not exist as the Photon Cloud is changing too frequently. Servers get added or removed and even new regions show up from time to time. This means it is impossible to add the Photon Cloud (as a whole) to an allow-list.

It is a different topic for an **Enterprise Cloud**. We'd discuss this via mail.

Apps within the **Photon Industries Circle** can rely on the host name for allow-listing: **\*.photonindustries.io**. Make sure to use `ns.photonindustries.io` as "Server" to connect to.

## Load Balancing

#### What is the maximum number of players supported by Photon rooms?

Most Photon multiplayer games have 2-16 players, but the theoretical limit of players/peers per room can be quite high.

There are Photon games live with 32 or even 64 players and in virtual conferencing scenarios it can be in the hundreds.

However, sending too many messages per second (msg/s per room) can cause performance issues depending on the client's processing power coping with data.

While high player numbers in e.g. turnbased games are totally fine, more than 32 players in a fast-paced action game likely will require you to implement interest management.

This way not every player receives every message from all the other players.

The number of players per room is the main factor for increasing data traffic inside the game room:

This is why we recommend keeping your msg/s per room count below 500.

Photon does not enforce this limit, but relies on a fair use policy.

Keeping an eye on your bandwidth usage is always important and it helps to ensure you stay within your plans included traffic range of 3GB per CCU.

#### Is there a limit for Photon strings?

Photon uses strings for lots of purposes: room name, lobby name, UserID, Nickname, custom property key, etc.

Technically, the Photon binary protocol can serialize strings up to 32767 one-byte-characters but it is strongly recommended to use the shortest strings possible.

For names and UserIDs 36 characters should be enough (e.g. a GUID is 36 characters).

For custom property keys, you should use shorter strings to minimize their overhead.

This is especially important for properties which are visible in the lobby, as those are part of the room lists and get sent to everyone in the lobby, not just the couple of clients in a room.

#### Is there a limit for the number of custom properties?

Yes, there is a limit: Each user can set 13k key-values and their total size in the room can't go over 500kB.

The more Custom Properties are set, the longer the clients take to join a room, as clients receive all of the properties.

#### Can I send a huge message using Photon?

We do not recommend transferring large data (i.e. files) using Photon unless you know what you are doing.

We advise you to optimize the data you exchange and get in touch with us if you really have to send a very big message.

Photon Cloud has a server side limit for client buffers which is 500KB.

So depending on the context a message could be considered:

- "too big" for our per client buffer size on Photon Cloud > 500KB. If a client hits this limit in a short period of time, it will be disconnected by the server.
- "too big" to be sent with UDP without resulting in an amount of fragments that might cause problems > 100KB.
- "too big" to be sent without splitting it up into multiple UDP packets > 1.2KB (including protocols overhead).

For messages that get sent very regularly (10 times a second or even more often) we would recommend to keep their size below 1KB.

If a message only get sent rarely (for example once at the start of a match), then a size of multiple KB is still fine, but we still would recommend to keep it below 10KB.

Larger messages can be sent in a separate Enet channel (in UDP), which affects the other channels (state sync) less.

For large files, consider using a separate backend.

#### Which data should be sent reliable and which should be sent unreliable?

First of all, you should know reliability is an option only when the protocol used is UDP.

TCP has its own "reliability" mechanism(s) which are not covered here.

Sending something reliable means that Photon makes sure it arrives at the target client(s).

So in case clients do not receive an acknowledgment within time, they will repeat sending until the acknowledgment is received or the number of resends is exceeded.

Also, repeating reliable events may cause extra latency and make subsequent events delayed.

Examples for not using reliability:

- player position updates in realtime games
- voice or video chat (streaming)

Example for using reliability:

- turn events in turn-based games
- score updates that don't change rarely

#### Why do I have so many disconnects in my game?

The disconnects could be due to various reasons.

We already have this documentation page that can help you investigate the related issues: " [Analyzing Disconnects](/realtime/current/troubleshooting/analyzing-disconnects)".

#### How messages per second per room are calculated?

Photon server counts total inbound and outbound messages every second and divide it by the total number of rooms (on the same Master Server).

Any operation request or operation response or event is considered a message.

Photon operations return an optional operation response and trigger zero or more events.

Cached events are also counted as messages.

**Messages cost per in-room operation:**

| Operation | Success: Best Case | Success: Average Case | Success: Worst Case |
| --- | --- | --- | --- |
| Create | **2**<br>(SuppressRoomEvents=true) | **3**<br>\+ Join event (SuppressRoomEvents=false, default) | **3** |
| Join | **2 + k**<br>(SuppressRoomEvents=true)<br>\+ k \* cached custom event | **2 + n + k**<br>\+ n \* Join event (SuppressRoomEvents=false, default) | **2 + 2 \* n + k**<br>\+ n \* ErroInfo event (HasErrorInfo=true) |
| Leave | **2**<br>(SuppressRoomEvents=true) | **1 + n**<br>\+ (n - 1) \* Leave event (SuppressRoomEvents=false, default) | **2 + (n - 1) \* 2**<br>\+ (n - 1) \* ErroInfo event (HasErrorInfo=true) |
| RaiseEvent | **1**<br>(no operation response)<br>(target: interest group with no subscribers) | **1 + n**<br>\+ n \* custom event<br>(target: all/broadcast) | **2 + 2 \* n**<br>\+ n \* ErroInfo event (HasErrorInfo=true)<br>\+ Auth event (token refresh) |
| SetProperties | **2**<br>Broadcast=false | **2 + n**<br>\+ n \* PropertiesChanged event (default: Broadcast=true, BroadcastPropertiesChangeToAll=true) | **2 + 2 \* n**<br>\+ n \* ErrorInfo event (HasErrorInfo=true)<br>\+ 1 in case of CAS or BroadcastPropsChangeToAll |

#### How do I calculate traffic consumed by a user?

This is a complex topic.

First you need to know that any calculation done is just a theoretical estimation and may not reflect reality.

We recommend building a Proof-of-Concept and use it to gather real data.

That being said here is how to estimate traffic generated by a single user inside a room:

Let's assume the following:

- a room has N players.
- a player sends F messages per second (message send rate in Hz)
- average message size is X (in bytes, payload (P) + protocol overhead (O))
- an average player spends H hours per month on your game

If we do not consider ACKs, connection handling (establishment, keep alive, etc.) commands and resends.

Then we say that on average, a CCU consumes C (in bytes/month) in your game as follows:

C = X \* F \* N \* H \* 60 (minutes per hour) \* 60 (seconds per minute)

#### How to quickly rejoin a room after disconnection?

To recover from an unexpected disconnection while joined to a room, the client can attempt to reconnect and rejoin the room.

We call this a "quick rejoin".

A quick rejoin will succeed only if:

- the room still exits on the same server or can be loaded:


If a player leaves a room, the latter can stay alive on Photon server if other players are still joined.


If the player is the last one to leave and the room becomes empty then the EmptyRoomTTL is how long it will remain alive waiting for players to join or rejoin.


If after EmptyRoomTTL the room is still empty and no one joined then it will be removed from Photon server.


If persistence conditions are met and webhooks are setup, the room state can be saved on the configured web service to be loaded later.
- the actor is marked as inactive inside: an actor with the same UserId exists inside the actors' list but not currently joined to the room.


This requires PlayerTTL to be different from 0.
- the PlayerTTL for the inactive actor did not expire: when an actor leaves a room with the option to come back we save his Deactivation timestamp.


As long as the room is alive, if PlayerTTL milliseconds are elapsed after Deactivation time the respective actor is removed from the actors' list.


Otherwise, when the actor tries to rejoin, if the difference in milliseconds between the time of the rejoin attempt and Deactivation time exceeds PlayerTTL then the actor is removed from the actors' list and the rejoin fails.


So an inactive actor can rejoin a room only for PlayerTTL milliseconds after Deactivation time.

A "quick rejoin" is composed of two steps:

- Reconnect: simply call the appropriate connect method once disconnected.
- Rejoin: call `loadBalancingClient.OpRejoin(roomName)`.

## Billing

#### Do you have special offers for students, hobbyists or indies?

No, but all our products have a free tier and a one-off entry-plan.

We also usually take part in Unity's asset store sales and occasionally give vouchers to lucky ones.

#### Can I combine more than one 100 CCU plan for a single Photon application?

No.

The 100 CCU plans are not stackable with each other. Only one can be applied per AppId.

They can be combined with paid subscriptions and contribute to the CCU for their duration.

They do not automatically renew.

In case you need more CCU for a single app, the next higher plan is the 500 CCU one.

If you subscribe to a monthly or yearly plan, then you will still keep the 100 CCUs for 12 months on top of / in addition to the CCU from your monthly/yearly plan.

One exception to this rule are the Free100 plans for Quantum and Fusion. These can be combined with 100 CCU coupons from the Unity Asset store (Quantum Plus and Fusion Plus).

#### How much traffic is included in my Photon plan and what happens if my app generates traffic beyond the included limit?

Photon Public Cloud and Premium Cloud plans include 3GB per CCU.

For example, a monthly plan with 1,000 CCUs includes 3 TB of traffic per month.

If your app generates more traffic, we will automatically send a heads up via email. You will receive an automatically generated overage invoice at the end of the month at your Photon account email address. The invoice amount is based on the following calculation:

> Total traffic - included traffic = overage traffic (in GB)

Traffic is calculated with $0.05 / $0.10 per GB depending on the Photon Cloud region used. This invoice is automatically charged to your credit card on file.

#### What happens if the Peak CCU exceeds the booked CCU of my Photon Cloud plan?

If you subscribed to a 500 CCU / 1,000 CCU / 2,000 CCU plan, “CCU burst” is automatically activated for your application. The Photon Cloud will allow more CCU than you booked to deliver the best possible experience for your users.

Once the burst kicks in, you are obliged to upgrade to the required subscription tier within 48 hours, according to the terms agreed.

If you do not upgrade, we will send an “overage invoice” to your Photon account email address and charge each CCU above your subscribed plan with a fee of $0.75 / $1.00 (based on the used SDK) per CCU. This invoice is automatically charged to your credit card on file.

Photon charges the “Peak CCU”, which is the sum of the peak CCU per region added up in a given month. Please make sure to upgrade even if the usage decreases after reaching its peak to avoid overage charges.

Back to top

- [Photon Cloud](#photon-cloud)

  - [Is Photon Cloud down?](#is-photon-cloud-down)
  - [Is there a default Photon Cloud region?](#is-there-a-default-photon-cloud-region)
  - [Can we get a list of all Cloud servers / IPs?](#can-we-get-a-list-of-all-cloud-servers-ips)

- [Load Balancing](#load-balancing)
- [Billing](#billing)