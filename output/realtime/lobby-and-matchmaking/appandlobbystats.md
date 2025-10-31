# appandlobbystats

_Source: https://doc.photonengine.com/realtime/current/lobby-and-matchmaking/appandlobbystats_

# App and Lobby Stats

Photon servers can broadcast application and lobby statistics to clients.

You can make use of this data to enhance your client-driven matchmaking.

You can also brag about these statistics in your game to show how popular it is. :\]

## Application Statistics

While connected to a Master Server, clients receive application statistics for the current region.

The applications statistics are:

- Number of live rooms:

C#
```csharp
loadBalancingClient.RoomsCount

```


C++
```cpp
Client::getCountGamesRunning()

```

- Number of players not joined to rooms:

C#
```csharp
loadBalancingClient.PlayersOnMasterCount

```


C++
```cpp
Client::getCountPlayersOnline() - Client::getCountPlayersIngame()

```

- Number of players inside rooms:

C#
```csharp
loadBalancingClient.PlayersInRoomsCount

```


C++
```cpp
Client::getCountPlayersIngame()

```

- Total number of connected players:

C#
```csharp
loadBalancingClient.PlayersOnMasterCount + loadBalancingClient.PlayersInRoomsCount

```


C++
```cpp
Client::getCountPlayersOnline()

```


AppStats event is sent to client every five seconds.

In the native C++ SDK, the `Listener` class provides a callback every time to know that the statistics' getters are updated:

C++
```cpp
virtual void onAppStatsUpdate(void) {}

```


## Lobby Statistics

Lobby statistics can be useful if a game uses multiple lobbies and you want to show the activity.

Lobby statistics are per region.

Per typed lobby (name + type) you can get information about:

- Number of live rooms
- Total number of players joined to the lobby or joined to the lobby's rooms

### Automatically Get Lobby Statistics

Lobby statistics events are sent as soon as the client is authenticated to a master server.

Then they are sent every minute.

Lobby statistics events are disabled by default.

- C#

Before connecting, to enable lobby statistics:

C#
```csharp
loadBalancingClient.EnableLobbyStatistics = true;

```


Get the statistics from the `ILobbyCallbacks.OnLobbyStatisticsUpdate` callback which could be useful to update your UI.

- C++

To enable lobby statistics, one has to pass true for parameter `autoLobbyStats` to the constructor of class `Client`:

C++
```cpp
Client(LoadBalancing::Listener& listener, const Common::JString& applicationID, const Common::JString& appVersion, nByte connectionProtocol=Photon::ConnectionProtocol::DEFAULT, bool autoLobbyStats=false, nByte regionSelectionMode=RegionSelectionMode::DEFAULT);

```


The `Listener` class provides the following optional callback whenever a lobby stats event arrives (when enabled):

C++
```cpp
virtual void onLobbyStatsUpdate(const Common::JVector<LobbyStatsResponse>& lobbyStats) {}

```


### Explicitly Get Lobby Statistics

You can explicitly request lobby statistics using an operation call when not joined to a room:

- C#

This is currently not implemented.

- C++

C++
```cpp
Client::opLobbyStats()

```


The `Listener` class provides the following optional callback when the response arrives:

C++
```cpp
virtual void onLobbyStatsResponse(const Common::JVector<LobbyStatsResponse>& lobbyStats) {}

```


Back to top

- [Application Statistics](#application-statistics)
- [Lobby Statistics](#lobby-statistics)
  - [Automatically Get Lobby Statistics](#automatically-get-lobby-statistics)
  - [Explicitly Get Lobby Statistics](#explicitly-get-lobby-statistics)