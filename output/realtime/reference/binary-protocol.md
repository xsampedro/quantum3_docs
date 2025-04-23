# binary-protocol

_Source: https://doc.photonengine.com/realtime/current/reference/binary-protocol_

# Binary Protocol

Photon and its clients are using a highly optimized binary protocol to communicate.

It's compact, yet easy to parse.

You don't have to know its details, as it is handled behind the scenes.

Just in case you want to know it anyways, this page describes it.

## Communication Layers

The Photon binary protocols are organized in several layers.

On its lowest level, UDP is used to carry the messages you want to send.

Within those standard datagrams, several headers reflect the properties you expect of Photon:

optional reliability, sequencing, aggregation of messages, time synchronization and various others.

The following chart shows the individual layers and their nesting:

![Photon Server: Binary Protocol UDP Layers](/docs/img/BinaryProtocol-udp-layers.png)
Photon Server: Binary Protocol UDP Layers


In words: Any UDP package contains an eNet header and at least one eNet command.

Each eNet command carries our messages: An operation, result or event.

Those in turn consist of the operation header and the data you provided.

Below, you will find tables listing each part of the protocols and its size.

This should give you an idea of the traffic needed by certain content you send.

If you wanted to, you could calculate how big any given message becomes.

To save you this work, we have an example, too.

## Example: Join "somegame"

Let's take a look at the `JoinRoom` operation (implemented by the Lite or LoadBalancing applications).

Without properties, this is an operation with a single parameter.

The "operation parameters" require: 2 + count(parameter-keys) + size(parameter-values) bytes.

The string "somegame" uses 8 bytes in UTF8 encoding and strings require another 3 bytes (length and type information). **Sum: 14 bytes.**

The parameters are wrapped into an operation-header, which is 8 bytes.

**Sum: 22 bytes.**

The operation code is not encoded as parameter.

Instead it's in the operation-header.

The operation and its header are wrapped into eNet's "send reliable" command as payload.

The command is 12 bytes. **Sum: 34 bytes.**

Let us assume the command is sent when no other commands are queued.

The eNet package header is 12 bytes. **Sum: 46 bytes for the reliable, sequenced operation.**

Last but not least, the eNet package is put into a UDP/IP datagram, which adds 28 bytes headers.

Quite a bit, compared to the 46 bytes so far.

Sadly, this can't be avoided completely but our aggregation of commands can save you a lot of traffic, sharing those headers.

**The complete operation takes up 74 bytes.**

The server will have to acknowledge this command. The ACK command is 20

bytes. If this is sent alone in a package, it will take up 40 bytes in

return.

![Binary Protocol Hex Bytes](/docs/img/BinaryProtocol-HexBytes.jpg)## Enet Channels

Enet Channels allow you to use several independent command-sequences.

In a single channel, all commands are sent and dispatched in order.

If a reliable command is missing, the receiving side can't continue to dispatch updates and things lag.

When some events (and operations) are independent from one-another, you can put them im separate channels.

Example: In-room chat messages should be reliable but if they arrive late, the position updates should not be delayed by (temporarily) missing messages.

By default Photon has two channels and channel zero is the default to send operations.

The operations join and leave are always sent in channel zero (for simplicity).

There is a "background" channel 255 used internally for connect and disconnect messages.

This is ignored for the channel count.

Channels are prioritized: the lowest channel number is put first.

Data in a higher channel might be sent later when a UDP package is already full.

## Operation Content - Serializable Types

You can find detailed table of Photon serializable types and their respective sizes at [this link](/realtime/current/reference/serialization-in-photon#photon-supported-types).

On the client side, operation parameters and their values are aggregated within a hashtable.

As each parameter is resembled by a byte key, operation requests are streamlined and use less bytes than any "regular" hashtable.

Results for operations and events are encoded in the same way as operation parameters.

An operation result will always contain: "operation code", "return code" and a "debug string".

Events always contain: "event code" and the "event data" hashtable.

## Enet Commands

| name | size | sent by | description |
| --- | --- | --- | --- |
| connect | 44 | client | reliable, per connection |
| verify connect | 44 | server | reliable, per connect |
| init message | 57 | client | reliable, per connection (choses the application) |
| init response | 19 | server | reliable, per init |
| ping | 12 | both | reliable, called in intervals (if nothing else was reliable) |
| fetch timestamp | 12 | client | a ping which is immediately answered |
| ack | 20 | both | unreliable, per reliable command |
| disconnect | 12 | both | reliable, might be unreliable in case of timeout |
| send reliable | 12 + payload | both | reliable, carries a operation, response or event |
| send unreliable | 16 + payload | both | unreliable, carries a operation, response or event |
| fragment | 32 + payload | both | reliable, used if the payload does not fit into a single datagram |

## UDP Packet Content - Headers

| name | size \[bytes\] | description |
| --- | --- | --- |
| udp/ip | 28 + size(optional headers) | IP + UDP Header. <br>Optional headers are less than 40 bytes. |
| enet packet header | 12 | contains up to byte.MaxValue commands <br>(see above) |
| operation header | 8 | in a command (reliable, unreliable or fragment) <br>contains operation parameters <br>(see above) |

Back to top

- [Communication Layers](#communication-layers)
- [Example: Join "somegame"](#example-join-somegame)
- [Enet Channels](#enet-channels)
- [Operation Content - Serializable Types](#operation-content-serializable-types)
- [Enet Commands](#enet-commands)