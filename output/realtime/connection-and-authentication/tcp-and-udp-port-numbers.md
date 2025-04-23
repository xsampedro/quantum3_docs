# tcp-and-udp-port-numbers

_Source: https://doc.photonengine.com/realtime/current/connection-and-authentication/tcp-and-udp-port-numbers_

# Port Numbers

## Photon Cloud

In most cases you don't have to think about port numbers, as the client SDKs use a default port to reach the Photon Cloud Name Servers.

Every time a Photon server forwards a client to another server, the assigned address includes the port. This means that the port number of the Name Server is the most important one as this is not fetched.

Port numbers are specific per transport protocol and often per server type. For reference, the table below lists commonly used ports for Photon. These are available on the Photon Cloud.

### Default Ports by Protocol and Purpose

| Port Number | Protocol | Purpose |
| --- | --- | --- |
| 5058 or 27000 | UDP | Client to Nameserver (UDP) |
| 5055 or 27001 | UDP | Client to Master Server (UDP) |
| 5056 or 27002 | UDP | Client to Game Server (UDP) |
| 4533 | TCP | Client to Nameserver (TCP) |
| 4530 | TCP | Client to Master Server (TCP) |
| 4531 | TCP | Client to Game Server (TCP) |
| 80 or 9090 | TCP | Client to Master Server (WebSockets) |
| 80 or 9091 | TCP | Client to Game Server (WebSockets) |
| 80 or 9093 | TCP | Client to Nameserver (WebSockets) |
| 443 or 19090 | TCP | Client to Master Server (Secure WebSockets) |
| 443 or 19091 | TCP | Client to Game Server (Secure WebSockets) |
| 443 or 19093 | TCP | Client to Nameserver (Secure WebSockets) |

Back to top

- [Photon Cloud](#photon-cloud)
  - [Default Ports by Protocol and Purpose](#default-ports-by-protocol-and-purpose)