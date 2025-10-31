# setup

_Source: https://doc.photonengine.com/quantum/current/addons/plugin-sdk/setup_

# Setup

Available in the [Gaming Circle](https://www.photonengine.com/gaming) and [Industries Circle](https://www.photonengine.com/industries)

## Installation

1. _(Windows only)_ Unblock `Photon-Quantum-X.X.X.Plugin-SDK.zip` by `right-click > Properties > Unblock` before unzipping it.
2. Unzip the archive into a folder inside or next to the Quantum Unity project (Quantum tools will auto-detect the location).

|     |     |     |
| --- | --- | --- |
| ```<br>`<br>MyQuantumProject<br>├─Assets<br>├─Library<br>├─..<br>└─PluginSDK<br>    `<br>` | Or | `<br>`<br>MyQuantumProject<br>├─Assets<br>├─Library<br>└─..<br>PluginSDK<br>    `<br>``` |

3. _(Running local Photon Server)_ Download and copy a Photon Server license from the [Photon dashboard](https://dashboard.photonengine.com/selfhosted) into the `Photon.Server/deploy\_win/bin`.

**PS:** only Circle membership _owners_ have access to the license file.

## Testing The Server Plugin

- Open the `Quantum.Plugin.Custom.sln` in Visual Studio or Rider, set the `Quantum.Plugin.Custom.Stock` as start-up project and ignore that other project may file to compile;
- Press F5 to start the Photon-Server;
- Open the game server log file under `Photon.Server\\deploy\_win\\log\\GSGame.log`;
- Open the Unity project, select `PhotonServerSettings.asset` and press the `Local Name Server` button;
- Start an online game using the Quantum Menu.

## Testing Server Simulation

- Select and configure the `QuantumDotnetProjectSettings` asset in Unity

  - Add include paths to all simulation code or Qtn-file folders or mark the folders with the `QuantumDotNetInclude` Unity asset label;
  - Simulations folder for example are `Assets\\Photon\\QuantumAsteroids\\Simulation`;
- Select the `QuantumDotnetBuildSettings` asset and press `Detect Plugin SDK`. If this fails, set the folder (that contains the `Photon.Server` folder) manually;
- Select the `QuantumDotnetBuildSettings` asset and press `Sync Plugin SDK Server Simulation`. This step builds the non-Unity simulation dll, exports the Quantum Unity DB and copies the LUT files:

  - Libraries to the `Lib` folder;
  - Assets to `Photon.Server\\deploy\_win\\Plugins\\QuantumPlugin3.0\\bin\\assets`;
- Open the `Quantum.Plugin.Custom.sln` in Visual Studio or Rider;

  - Set the `Quantum.Plugin.Custom.Sample.ServerSimulation` project as the start-up project;
  - Press F5 to start the Photon-Server;
  - Open the game server log file under `Photon.Server\\deploy\_win\\log\\GSGame.log`;
- Open the Unity project, select `PhotonServerSettings.asset` and press the `Local Name Server` button;
- While selecting the `PhotonServerSettings.asset` change the `Auth Mode` to `Auth Once`;
- Start an online game using the Quantum Menu.

![Auth Once](/docs/img/quantum/v3/addons/plugin_sdk/auth-once.png)## SDK Content

### Lib Folder

The Lib folder includes all dependencies required to compile and run the Quantum Photon Server plugin.

The `PhotonHivePlugin.dll` for example is the interface for general Photon Server plugins. `Quantum.Deterministic.Plugin.dll`, `Quantum.Deterministic.Server.dll` and `Quantum.Deterministic.Server.Interface.dll` are the main libraries of the Quantum plugin. `Quantum.Deterministic.dll` and `Quantum.Log.dll` are usually replaced by the libraries of custom simulation build (from Unity).

### Photon.Server Folder

The folder includes the local Photon Server. The `bin` folder has the Photon.Server executables, `LoadBalancing` and `NameServer` has the server code and `Plugins` has the server plugins. The `log` folder contains the server logs (e.g. the game server log `GSGame.log`).

The plugin configuration file `LoadBalancing\\GameServer\\bin\\plugin.config` includes the local configuration key value store that online servers get from the Photon Dashboard.

The custom plugin libraries will be outputted to `Plugins\\QuantumPlugin3.0\\bin`. This is also the folder that is uploaded to the Photon Enterprise cloud in the end.

### Quantum.Plugin.Custom.\* Folders

These folders contain sample projects for a custom Quantum plugin (see next section).

## Plugin Samples

The `Quantum.Plugin.Custom.sln` comes with multiple sample projects. Select one as a start-up project to be able to start it with F5.

![Selecting a Project](/docs/img/quantum/v3/addons/plugin_sdk/selecting-project.png)

See in the `launchSettings.json` how F5 will start a Photon Server:

The shared csproj file `Quantum.Plugin.Shared.csproj.include` is used in each project to support basic integration tasks like:

- Copying the output to the destination folder: Photon.Server\\Plugins\\QuantumPlugin3.0\\bin;
- Zipping the plugin to prepare for Photon Cloud upload;
- Checking for a license file.

### Stock Sample

This sample is a blank slate which only consists of a custom `DeterministicPluginFactory` which is responsible to instantiate the required `DeterministicPlugin`. It will just instantiate the unmodified base plugin object.

C#

```csharp
public override DeterministicPlugin CreateDeterministicPlugin(IPluginHost gameHost, String pluginName, Dictionary<String, String> config, IPluginLogger logger, ref String errorMsg) {
  return new DeterministicPlugin();
}

```

### HTTP Sample

In some situations, you may need to query an external service to get some information or to perform some action (such as giving a player's account rewards).

This sample is a slightly more advanced version of the stock plugin. It includes a custom `DeterministicPluginFactory` and a custom `DeterministicServer` called `QuantumCustomServer`. This class overrides some of the base plugin methods to add custom behavior (HTTP calls in this case).

The custom server utilizes the `HTTPExample` class to perform a simple HTTP GET request to a public API. The `HTTPExample` class is a simple wrapper around the PluginHost's HTTP calls.

C#

```csharp
HttpExample.SendAsync(host, result => {
    if (result) {
        Log.Info("HTTP Asyncronous Response: Success");
    } else {
        Log.Info("HTTP Asyncronous Response: FAILED");
    }
});

```

### Server Simulation Sample

This sample is the same as the Stock sample but includes a Server Simulation.

In order to do this, the instance of `DeterministicServer` created is also provided with a `DotNetSessionRunner` in its constructor.

To compile this project all steps from the section _Testing Server Simulation_ have to be completed first.

C#

```csharp
public override DeterministicPlugin CreateDeterministicPlugin(IPluginHost gameHost, String pluginName, Dictionary<String, String> config, IPluginLogger logger, ref String errorMsg) {
  var sessionRunner = new DotNetSessionRunner {
    AssetSerializer = new QuantumJsonSerializer()
  };
  return new DeterministicPlugin(new DeterministicServer(sessionRunner));
}

```

The `DotNetSessionRunner` is a simple implementation of the `ISessionRunner` interface that allows the Quantum simulation to be run outside of Unity.

The `QuantumJsonSerializer` is used to serialize and deserialize the AssetDB.

#### Events and Callbacks from Server Simulation

It is possible to, from the custom plugin, react to Events (defined and triggered in the simulation code) and Callbacks.

This can be achieved that by getting the `EventDispatcher` and `CallbackDispatcher` from the `DotNetSessionRunner` and subscribing to the desired events and callbacks.

C#

```csharp
private void SubscribeToEvents(DotNetSessionRunner sessionRunner) {
    var eventDispatcher = (EventDispatcher)sessionRunner.EventDispatcher;
    eventDispatcher.Subscribe<EventFoo>(this, OnEventFoo);
}

private void OnEventFoo(EventFoo foo) {
    // Do something with the event
}

```

C#

```csharp
private void SubscribeToCallbacks(DotNetSessionRunner sessionRunner) {
    var callbackDispatcher = (CallbackDispatcher)sessionRunner.CallbackDispatcher;
    callbackDispatcher.Subscribe<CallbackGameStarted>(this, c => Log.Info("Game Started"));
}

```

## Further Readings

The Quantum custom plugin is based of Photon-Server V5 and follows the workflow described in the Photon Server docs. Dive into these docs for further reading: [Photon-Server V5 Manual](/server/v5/plugins/manual)

To upload the plugin to the Photon Enterprise cloud the following tutorial is helpful: [Enterprise Plugin Setup](/quantum/current/addons/plugin-sdk/enterprise-setup)

Back to top

- [Installation](#installation)
- [Testing The Server Plugin](#testing-the-server-plugin)
- [Testing Server Simulation](#testing-server-simulation)
- [SDK Content](#sdk-content)

  - [Lib Folder](#lib-folder)
  - [Photon.Server Folder](#photon.server-folder)
  - [Quantum.Plugin.Custom.\* Folders](#quantum.plugin.custom.folders)

- [Plugin Samples](#plugin-samples)

  - [Stock Sample](#stock-sample)
  - [HTTP Sample](#http-sample)
  - [Server Simulation Sample](#server-simulation-sample)

- [Further Readings](#further-readings)