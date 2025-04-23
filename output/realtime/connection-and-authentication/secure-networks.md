# secure-networks

_Source: https://doc.photonengine.com/realtime/current/connection-and-authentication/secure-networks_

# Secure Networks

## Network Security Overview

As long as a device has access to the Internet, a Photon client on should be able to connect to Photon servers, as the client initiates any connection and replies are expected.

For security reasons, there are a few variations where the network:

- blocks UDP (and or other protocols).
- allows only some port ranges (per protocol).
- blocks some host name(s).

Below are some options to help get clients connected in those cases.

## WebSockets

Most of the networks with restricted internet access still allow a secure connection via Secure WebSockets, a well known standard internet protocol.

Configure your client to use the **protocol WebSocketSecure** and **set the port to 443** for it.

Apps that belong to an Industries Circle account should also **set the Name Server address to: ns.photonindustries.io**. This is the `AppSettings.Server` value. With that set, allow-lists in firewalls only need to include \*.photonindustries.io.

## Proxy Support

Photon clients support using proxies for WebSocket connections. In WebGL, proxy support is defined by the browser and system. All other platforms need some setup work, described below.

- Make sure your project contains the folder `PhotonLibs\\WebSocket\` containing **websocket-sharp.dll** among other files.
- Use the `PhotonWebSocket.asmdef` to include the WebSocket assembly on your target platforms.
- Add the define WEBSOCKET to your project (this is usually set per platform).

The Realtime API level now uses the `SocketWebTcp` class for WS/WSS connections, instead of the default `PhotonClientWebSocket`. To verify this, the SupportLogger (if available and enabled) will log a line including "Socket: SocketWebTcp" once connected.

Set the AppSettings.Server to: ns.photonindustries.io.

Make sure the port is 80 for ws or 443 for wss. Some client SDKs use these as default. Others need it configured as AppSettings.Port.

The `AppSettings.ProxyServer` field can be empty or contain one of the following options:

- `<proxyhostname>:<proxyport>` // just a proxy address
- `<user>:<pass>@<proxyhostname>:<proxyport>` // authentication and proxy address

Parts in < and > are variable info. Don't include the < and > signs in the proxy config.

To verify your proxy is used:

- increase the log level of the Photon networking layer to Info (e.g. via `AppSettings.NetworkLogging`).
- run Wireshark and monitor the specific port.

Back to top

- [Network Security Overview](#network-security-overview)
- [WebSockets](#websockets)
- [Proxy Support](#proxy-support)