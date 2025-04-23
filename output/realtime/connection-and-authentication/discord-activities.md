# discord-activities

_Source: https://doc.photonengine.com/realtime/current/connection-and-authentication/discord-activities_

# Discord Activities

Discord Activities let you [integrate your WebAssembly game directly into Discord](https://discord.com/blog/server-activities-games-voice-watch-together). In this environment, multiplayer games have to meet additional requirements, which are covered below. Please refer to [Discord's Development Guide](https://discord.com/developers/docs/activities/development-guides) for everything else.

From Photon perspective, the most important part is: You need to setup an URL Mapping for Photon's addresses and the clients must rewrite the server URLs accordingly.

This is supported with the Photon Realtime .Net SDKs and Unity SDKs v4.1.8.11 and up, as well as v5.1.1 and up.

## Requirements and Setup

Applications running as "Discord Activity" are hosted on a secure page (https), which means they can not communicate with just any server (address) they like. All communication must go to discordsays.com.

Discord solves this via "URL Mapping", which allows you to sets up an address on discordsays.com for arbitrary addresses needed by an app.

The URL Mapping for a Discord Activity looks like this (incomplete):

![Discord URL Mapping](/docs/img/Discord-Url-Mapping-v2.png)
Discord URL Mapping


Photon servers run on the domains `exitgames.com` and `photonengine.io`. Apps of the Industries Circle can use `\*.photonindustries.io` exclusively.

Each machine has it's own subdomain (e.g. `ns.photonengine.io` for the Name Servers). In some cases, Photon uses a path (typically "Master" or "Game") which should be preserved.

To cover the variable range addresses of the Photon Cloud, use "Parameter Matching" in the map:

The prefix `/photonengine/{subdomain}` can map to target: `{subdomain}.photonengine.io`.

## Address Rewriter

Setting up the URL Mapping for Photon opens a route through discordsays.com but the clients will still get regular Photon URLs for every server. On the client, Address Rewriting must be done via code to turn Photon URLs into the mapped URLs.

To modify the server addresses on a Photon client, a `client.AddressRewriter` function is used. Once registered, your `Func<string, ServerConnection, string>` gets called before any connection is made and can modify the address as needed.

For Discord Activities, the rewriting could look like this:

C#

```csharp
string clientId = "12345678"; // your app's discord client id
private string AddressRewriterDiscord(string address, ServerConnection serverType)
{
    bool isUri = Uri.TryCreate(address, UriKind.Absolute, out Uri uri);
    if (isUri)
    {
        string host = uri.Host;
        string[] hostSplit = host.Split('.');
        if (hostSplit.Length != 3)
        {
            Debug.LogError($"Host address could not be split into 3 parts (subdomain, domain and tld).");
            return address;
        }
        string subdomain = hostSplit[0];
        string domain = hostSplit[1];
        string discordAddress = $"{uri.Scheme}://{clientId}.discordsays.com/.proxy/{domain}/{subdomain}{uri.path}";
        //Debug.Log($"discordAddress: {discordAddress}");
        return discordAddress;
    }
    return address;
}

```

To apply this address rewriting, register the method before the client connects:

`client.AddressRewriter = this.AddressRewriterDiscord;`

Make sure to use this only on builds for Discord Activities.

## Related Topics

To successfully run the app as Discord Activity, you will also have to adjust the data paths in the html Unity generates to also use the ".proxy" address.

![Discord Data Paths](/docs/img/Discord-Data-Paths.png)
Discord Data Paths
Back to top

- [Requirements and Setup](#requirements-and-setup)
- [Address Rewriter](#address-rewriter)
- [Related Topics](#related-topics)