# encryption

_Source: https://doc.photonengine.com/realtime/current/reference/encryption_

# Encryption

Photon enables you to encrypt messages between client and server on demand, which is called Payload Encryption.

As it takes some performance, encryption is typically optional and sparingly used.

If a message can be shared with everyone in a room, it can probably be sent without encryption, too.

Any message which clearly contains a userID is always sent encrypted. This includes authentication values.

Encryption is always between a client and the server and Photon does not offer end-to-end encryption between clients.

## Technical Details

## UDP and TCP

Payload Encryption is the default for UDP and TCP connections. It uses a 256 bit key for **AES encryption** on demand.

Recent client SDKs connect to the Name Server via WSS and switch to UDP or TCP on the following connections. This way, the encryption key for the Payload Encrytion is securely exchanged on a TLS connection. This is done when the AuthMode is `AuthOnceWss` for a client.

If `AuthOnceWss`is not used, the key exchange is automatically done via **Diffie-Hellman Key Exchange** when a client connects to the Name Server via UDP or TCP.

If an operation is sent encrypted, any data you send and all operation parameters are serialized and then encrypted.

A lean header per message (length, etc.) remains unencrypted.

### WSS Transport Protocol

When the WSS Transport Protocol is used, all communication gets encrypted, including headers.

In that case, the option to send individual operations with Payload Encryption is safely ignored.

WSS is the default protocol for all WebGL exports from Unity, independent from the protocol selection in the PhotonServerSettings.

## Photon Tokens

Photon's Tokens are typically handled behind the scenes.

It's not very visible in the client APIs but it makes sense to be aware of them.

Once a client is authenticated, the server will issue a Token, an encrypted summary of the client's authentication values to be used on other servers.

The Token is only used by the Photon servers and can't be read by the clients.

It's possible to "inject" some data into the token, to share it between servers (via the client).

A Server Plugin may read this shared data from the Token.

### Token Refresh

By default Photon tokens expire after 1 hour but in most cases they are refreshed for the client automatically.

The refresh happens in two cases:

1. when switching from Master Server to Game Server, if you create or join a room.
2. as long as the client keeps raising events inside the room.

If a client stays inside a room for more than 1 hour without raising any event, the token will not be refreshed and expires.

This will not disconnect the client from the Game Server but when it leaves the room, connection to the Master Server will fail and the client has to reconnect and authenticate again.

## C# SDKs Encryption of Operations

In all C# client SDKs, we have a class called `PhotonPeer`.

To send an operation encrypted, call `PhotonPeer.SendOperation` with the `sendOptions` parameter that has `sendOptions.Encrypt = true`.

But usually you do not need to use that class or call that method yourself.

It is done internally in a lower level.

In the C# APIs we provide high level classes for Photon networking clients.

- In PUN1, the `PhotonNetwork.networkingPeer` is a `PhotonPeer`.
- In Photon Realtime, `LoadBalancingClient.LoadBalancingPeer` is a `PhotonPeer`.
- In Photon Voice, `LoadBalancingTransport` extends `LoadBalancingClient`.
- In PUN2, the `PhotonNetwork.NetworkingClient` is a `LoadBalancingClient`.
- In Photon Chat, `ChatClient.chatPeer` is a `PhotonPeer`.

Back to top

- [Technical Details](#technical-details)
- [UDP and TCP](#udp-and-tcp)

  - [WSS Transport Protocol](#wss-transport-protocol)

- [Photon Tokens](#photon-tokens)

  - [Token Refresh](#token-refresh)

- [C# SDKs Encryption of Operations](#c-sdks-encryption-of-operations)