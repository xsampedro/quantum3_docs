# ipv6

_Source: https://doc.photonengine.com/realtime/current/connection-and-authentication/ipv6_

# IPv6

## Concept

The transition to IPv6 is inevitable and already frequently used in mobile networks. For this reason, Apple for example no longer accepts submissions to the App Store if IPv6 is not supported.

All recent Photon SDKs support IPv6 networks as long as a DNS64/NAT64 service is available (which is what Apple requires).

So far, all our Photon Cloud server are running on IPv4 (which is no problem with the DNS64/NAT64 services).

If you are hosting your own Photon Server, please visit [this page](/server/current/operations/ipv6) for more instructions.

**Option A: IPv6-ready clients, IPv4-only servers**

In this case, clients are in a network that only uses IPv6 addresses or hostnames. Hostnames that are only available via IPv4 are translated into IPv6 addresses for these networks (via DNS64/NAT64).

As long as you avoid using IPv4 addresses directly, servers can be reached and don't need a IPv6 address themselves. This option can be tested with the help of a Mac (see below).

Photon Server should be configured using domain names (FQDN) instead of IP addresses. See [how to](#hostnames) do this in Photon Server.

**Option B: both client and server are running on IPv6**

This option is for an IPv6-only network and is supported by Photon Server v4 and up.

### Non-breaking strategy

As required, server addresses will be sent to the client in a suitable form:

- IPv4 address: when both server and client use IPv4.
- IPv6 address: when both server and client use IPv6.
- DNS hostname: when client uses IPv6 and server uses IPv4.

### Testing IPv6 with DNS64/NAT64

We recommend testing on a local environment using the methods proposed by Apple.

You can easily [setup a IPv6 WiFi using a Mac](https://developer.apple.com/library/ios/documentation/NetworkingInternetWeb/Conceptual/NetworkingOverview/UnderstandingandPreparingfortheIPv6Transition/UnderstandingandPreparingfortheIPv6Transition.html#//apple_ref/doc/uid/TP40010220-CH213-SW16).

Back to top

- [Concept](#concept)
  - [Non-breaking strategy](#non-breaking-strategy)
  - [Testing IPv6 with DNS64/NAT64](#testing-ipv6-with-dns64nat64)