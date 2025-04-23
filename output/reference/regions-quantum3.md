# regions-quantum3

_Source: https://doc.photonengine.com/quantum/current/reference/regions-quantum3_

# Regions

The Photon Cloud provides low latency gaming globally by hosting servers in various regions.

Clients get the list of regions from our Photon Name Servers. Over the lifetime of a project, new regions may be added or old ones may be deprecated and removed.

Each region is completely separate from the others and consists of a Master Server (for matchmaking) and Game Servers (hosting rooms).

Below is an outline of the client workflow:

(1) Clients retrieve a list of available regions from the Photon Name Servers.

(2) The client selects one of the available regions — either automatically based on the lowest ping or via manual player choice — and connects to the region's Master Server. The Master Server is responsible for matchmaking and creating or joining a room on a Game Server within that region.

(3) Finally, the client connects to the designated Game Server in order to join the room. This Game Server is located in the same region as the Master Server the client connected to in step (2).

![Photon Cloud Regions' Connect Flows](https://doc.photonengine.com/docs/img/photon-cloud-connect-name-server.png)
Connect to Photon Cloud regions

With the Region Allowlist, you can define which regions should be available per AppId (see below).

## Available Regions

The Photon Cloud consists of servers in several regions, distributed across multiple hosting centers worldwide.

Each Photon Cloud region is identified by a "region code", which is a case insensitive, short string.

For example, "EU" or "eu" are both accepted and refer to the same Europe region.

The lists below indicate the location of the hosting center and the region code for each region.

### Photon Cloud for Gaming

The Photon Products Quantum, Fusion, Voice, Realtime and PUN are available to the **Photon Cloud for Gaming** in the following regions.

| Region | Hosted in | Code |
| --- | --- | --- |
| Asia | Singapore | asia |
| Australia | Sydney | au |
| Canada, East | Montreal | cae |
| Chinese Mainland ( [See Instructions](#using-the-chinese-mainland-region)) | Shanghai | cn |
| Europe | Amsterdam | eu |
| Hong Kong | Hong Kong | hk |
| India | Chennai | in |
| Japan | Tokyo | jp |
| South Africa | Johannesburg | za |
| South America | Sao Paulo | sa |
| South Korea | Seoul | kr |
| Turkey | Istanbul | tr |
| United Arab Emirates | Dubai | uae |
| USA, East | Washington D.C. | us |
| USA, West | San José | usw |
| USA, South Central | Dallas | ussc |

Photon Chat is available in the following regions:

| Region | Hosted in | Code |
| --- | --- | --- |
| Asia | Singapore | asia |
| Europe | Amsterdam | eu |
| USA, East | Washington D.C. | us |
| Chinese Mainland ( [See Instructions](#using-the-chinese-mainland-region)) | Shanghai | cn |

### Photon Industries Premium Cloud

The Photon Products Quantum, Fusion, Voice, Realtime and PUN are available to the **Photon Industries Premium Cloud** in the following regions.

| Region | Hosted in | Code |
| --- | --- | --- |
| Asia | Singapore | asia |
| Australia | Sydney | au |
| Europe | Amsterdam | eu |
| India | Chennai | in |
| Japan | Tokyo | jp |
| South Africa | Johannesburg | za |
| South America | Sao Paulo | sa |
| South Korea | Seoul | kr |
| USA, East | Washington D.C. | us |
| USA, West | San José | usw |

Photon Chat is available in the following regions:

| Region | Hosted in | Code |
| --- | --- | --- |
| USA, East | Washington D.C. | us |

### Regions for China for Gaming/Industries

There are special conditions for using the Photon Cloud Region Chinese Mainland:

- Access must be unlocked ( [see below](#using-the-chinese-mainland-region))
- Photon Voice is not available in China
- 20CCU for development is free (non commercial use)
- Only a 500CCU subscription available on Photon Cloud
- Large setups need custom agreements

The Photon Products Quantum, Fusion, Realtime, PUN and Chat are available to the Photon Cloud in the following regions:

| Region | Hosted in | Code |
| --- | --- | --- |
| China Mainland | Shanghai | cn |

## Region Allowlist

The Region Allowlist enables you to customize the available regions per application directly from the dashboard. Clients using the Best Region feature, will adapt automatically.

By using using more or less regions, you balance the quality of service (roundtrip times are better, when there is a region close to players) versus the matchmaking experience (less regions mean more players per region).

To define the regions per app, [open the dashboard](https://dashboard.photonengine.com), click "Manage" for a chosen application and then click "Edit Allowlist".

You will find an input field to enter the list of allowed regions as follows:

- the available regions are listed above per SDK and sometimes separately for the Industries Circle.
- the allowlist must be a string of region codes separated by semicolons. e.g. "eu;us".
- region codes are case insensitive.
- undefined or unrecognized region codes will be ignored from the list.
- empty ("") or malformed string (e.g. ";;;") means all available regions are allowed.

Within 10 minutes of a change (confirm and save), the Name Servers will send the filtered list to connecting clients.

To avoid conflicts on the client side, connect to the "Best Region" by ping or make sure to pick a region received with the regions list.

**Note**: Changing the available regions for a popular app will affect the Peak CCUs in multiple regions, which is the basis for subscription fees. Adjust the subscription plan as needed to avoid the more expensive overage fees. Reducing the subscription is perfectly fine when the switch settled down.

## How To Choose A Region

Users in the US have the lowest latency if connected to the Photon Cloud US region. Easy.

But what if you have users from **all over the world**?

Options are..

- **a)** let the game client ping the different Photon Cloud regions and pre-select the one with the best ping, read our [how to](#howto).
- **b)** distribute client builds tied to a region, so users from different regions connect to different Photon Cloud regions or
- **c)** let the user choose a matching region from within your game\`s UI.
- **d)** let all users connect to the same region if the higher latency is acceptable for your gameplay.

All Photon Cloud apps are working in all available regions without any extra charge.

[See pricing.](https://www.photonengine.com/quantum/pricing)

Photon Cloud's dashboard lets you monitor the usage of your game in each region and easily upgrade or downgrade your subscription plan.

[Go to your dashboard.](https://dashboard.photonengine.com)

### C# Realtime API

Photon Realtime (used by most Photon SDKs) can detect the Best Region to connect to and enables you to stick to that region.

To do so, clients always fetch the list of available regions from the Name Server on connect.

The servers response is used to setup the `LoadBalancingClient.RegionHandler` which is also provided via the callback `OnRegionListReceived(RegionHandler regionHandler)`, as defined in the `IConnectionCallbacks`.

Typically, the next step is to call `regionHandler.PingMinimumOfRegions()` to detect the current ping to each region. You need to pass a method to call on completion and in best case you can also pass the "best region summary" from a previous run (explained below).

After pinging the servers, the (new) results are summarized in the `regionHandler.SummaryToCache` which should be saved on the device for later use.

Without the `SummaryToCache` from a previous session, all regions will be pinged, which takes a moment longer.

If a previous result is available, the client will check:

a. if the region list changed (covers the case if the "previous best region" is still available)

b. if the ping is no longer acceptable (>= 1.5x slower than previously saved reference value)

If either applies, all regions are pinged and a new result gets picked.

Using Best Region works well with the server-side Region Filter in the Dashboard.

It enables you to change the list of regions available to players on demand.

To access the list of regions or to override previous results, refer to the API Reference for regions.

### Best Region Considerations

"Best Region" option is not deterministic.

Sometimes it may be "random" due to little variations or exact same ping calculations.

Theoretically, you could:

- have the same exact ping to multiple regions from the same device. So it is random, if you end up with different regions on clients connected to the same network.
- different ping values for the same region on different devices (or different retries on the same device) connected to the same network.

For instance, in the case of "us" and "usw" (or "ru" and "rue"), you could either make use of the online regions allowlist to select the ones you want and drop the others or connect to an explicit region.

To debug, set the logging level to "Info" and clear the "current best region" (in PUN: PhotonNetwork.BestRegionSummaryInPreferences = null). Have a look at the details or send us the log via mail.

## Connect to a specific Master Server

To connect your clients to a specific region, set the `AppSettings.FixedRegion` to a valid Region code and call `ConnectUsingSettings(settings)`.

The SDK will get the master server address for the requested region from the Name Server (1 in the figure "Connect to Photon Cloud regions") and automatically connect you to the master server in the chosen region (2 in the figure "Connect to Photon Cloud regions").

With a FixedRegion, the client will not fetch the region list and skip pinging regions for a Best Region result. This speeds up the connection time.

If you compile a FixedRegion into your build, it can not be changed without an update. In best case, use Best Region and the Region Allowlist instead.

## How To Show A Region List

If you want to select the region at runtime - e.g. by showing a list of available regions to your players and let them choose - you need to connect to the name server first.

This will automatically fetch a list of currently available region master server addresses (1 in the figure "Connect to Photon Cloud regions").

While we write about "the name server", the name server is geographically load-balanced across available regions.
That keeps the time to request the master servers' addresses as low as possible.

### C# Client SDKs

C#

```csharp
    loadBalancingClient.ConnectToNameServer()

```

After a successful connection, `LoadBalancingClient.OpGetRegions()` gets called internally. The result of this sets up the `loadBalancingClient.RegionHandler` and calls `OnRegionListReceived` if your code implements it and registered for callbacks.

With the list of master servers, you could now ping all to figure out the best region to connect to for lowest latency gameplay, or let your players choose a region. This can be done with `RegionHandler.PingMinimumOfRegions()`.

When your client has determined a region, connect to the master server for that region (2 in the figure "Connect to Photon Cloud regions").

C#

```csharp
    loadBalancingClient.ConnectToRegionMaster(&#34;us&#34;)

```

Finally, join or create a room for your game (3 in the figure "Client Connect to Photon Cloud").

## Using The Chinese Mainland Region

You need to request access to the Chinese Mainland region for your Photon application. [Send us an email so we could unlock it for your AppID.](https://www.photonengine.com/contact)

On our Dashboard, you can not subscribe to paid plans to be used in the Chinese Mainland region.
Reach out to us by email to receive a quote for any subscription: [\[email protected\]](/cdn-cgi/l/email-protection).

The Photon Name Server has to be local to China, as the firewall might block the traffic otherwise.

The Chinese Photon Name Server is "ns.photonengine.cn".

Connecting with clients from outside of China mainland will most likely not produce good results.

Also, connecting from the Photon servers to servers outside of China mainland (e.g. for Custom Authentication, WebHooks, WebRPCs) might not be reliable.

**Important**: in the current phase, changes you make to your app via your dashboard are not automatically reflected in the app caches for China.
Let us know by email if you have an update request there.

For legal reasons, you need a separate build for China and we recommend using a separate AppId with it.

For example, use a compile condition (of your choice) to change the AppId and the Photon Name Server depending on the build.

Follow the instructions corresponding to your client SDK to make a special build for the Chinese market.

Back to top

- [Available Regions](#available-regions)

  - [Photon Cloud for Gaming](#photon-cloud-for-gaming)
  - [Photon Industries Premium Cloud](#photon-industries-premium-cloud)
  - [Regions for China for Gaming/Industries](#regions-for-china-for-gamingindustries)

- [Region Allowlist](#region-allowlist)
- [How To Choose A Region](#how-to-choose-a-region)

  - [C# Realtime API](#c-realtime-api)
  - [Best Region Considerations](#best-region-considerations)

- [Connect to a specific Master Server](#connect-to-a-specific-master-server)
- [How To Show A Region List](#how-to-show-a-region-list)

  - [C# Client SDKs](#c-client-sdks)

- [Using The Chinese Mainland Region](#using-the-chinese-mainland-region)