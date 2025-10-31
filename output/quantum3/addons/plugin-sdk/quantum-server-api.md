# quantum-server-api

_Source: https://doc.photonengine.com/quantum/current/addons/plugin-sdk/quantum-server-api_

# Quantum Server API

Available in the [Gaming Circle](https://www.photonengine.com/gaming) and [Industries Circle](https://www.photonengine.com/industries)

## Plugin Factory

The plugin factory is used by the Photon Server to instantiate the plugin code for individual rooms. The base Photon Server class `IPluginFactory2` is derived by the Quantum plugin factory `DeterministicPluginFactory`. Override and implement the `CreateDeterministicPlugin()` method to instantiate different plugins and/or server objects.

C#

```csharp
public override DeterministicPlugin CreateDeterministicPlugin(IPluginHost gameHost, String pluginName, Dictionary<String, String> config, ref String errorMsg)

```

## Basic Structure

Apart from the factory the Quantum server plugin consist of two additional, extended-able parts.

The `DeterministicPlugin` is build upon the Photon Server `PluginBase` class. Override the base class virtual methods when needed. The plugin class is instantiated per Photon room when the first client enters.

The `DeterministicServer` is instantiated when the Quantum simulation is started and controls all online aspects related to the simulation of the Quantum server (e.g. input, server simulation, etc).

## Logging

Caveat: server logging can easily result in floods of logs that are unusable and which can quickly result in destabilizing entire servers.

The `Quantum.Log` will be statically initialized using a Photon Plugin logger named `Quantum` before the first Plugin is instantiated. Quantum.Log debug methods use the `Conditional("DEBUG")` attribute which are only enabled when the using projects (Quantum Custom Server) have this define enabled.

Logging from a custom server can be done by calling `DeterministicPlugin.LogInfo()`, `DeterministicPlugin.LogWarning()`, or using the `DeterministicPlugin.Logger` directly. Also, new Photon plugin loggers can be created using `IPluginHost.CreateLogger()`.

The logs for a local server are located in this folder: `Photon.Server\\deploy\_win\\log`. The `GSGame.log` is where logs from the custom plugin show up. To retrieve logs from the Photon enterprise server contact us please.

## Server Simulation

For the Quantum server to run the simulation the `DeterministicServer` object needs to be created with an additional parameter of the type `IDeterministicSessionRunner`. Different to Quantum 2 starting and updating the simulation is now completely encapsulated. Check out the `DotNetSessionRunner` class in the Unity project.

C#

```csharp
var sessionRunner = new DotNetSessionRunner {
  AssetSerializer = new QuantumJsonSerializer()
};

```

When running the server simulation, the Quantum dependencies, the simulation dll and the assets **always** have to be in sync with the client builds. Run the sync process in Unity (e.g. by selecting the `QuantumDotnetBuildSettings` asset).

![Unity - Sync Server Simulation](/docs/img/quantum/v3/addons/plugin_sdk/unity-build-settings-sync.png)

The dependencies from the Unity SDK that the simulation was build with **take precedence** over the files from the plugin SDK inside the `Lib` folder (e.g. Quantum.Deterministic.dll, Quantum.Log.dll).

When the server simulation throws an exception it will log inside `GSGame.log`, the simulation is terminated but the **online game will continue** to run.

External dependencies to the `Newtonsoft.Json` package is required which works together with the `Quantum.Json` to support the asset db and RuntimeConfig, RuntimePlayer deserialization.

### Quantum Asset DB

The Quantum Asset DB must be available on the server to run the server simulation. By default an exported Json file is supported. Syncing the server simulation from Unity already exports and copies it to `Photon.Server\\deploy\_win\\Plugins\\QuantumPlugin3.0\\bin\\assets\\db.json`. It can also be exported using the Unity menu: `Tools > Quantum > Export > Asset Database`.

It's possible to replace the Json asset serialization by using a custom implementation of `IAssetSerializer`.

#### Embedded Quantum Asset DB

Optionally the exported asset file can be embedded into the Quantum.Simulation.dll.

XML

```xml
  <ItemGroup>
    <EmbeddedResource Include="db.json" />
  </ItemGroup>

```

The `EmbeddedDBFile` configuration property points to the name of the embedded resource. Always add the additional `Quantum.` prefix to the property.

The plugin will always try to load the Asset DB from an external file first. If none was found it will try to load from the embedded resource.

### Local Quantum Plugin Configuration

The Quantum plugin configuration values for a local Photon Server can be found in the `Photon.Server\\deploy\_win\\LoadBalancing\\GameServer\\bin\\plugin.config` file.

The default `plugin.config` contains a plugin XML configuration. Change or add new key-value variables accordingly.

These properties are replaced with the key-values settings from the Photon online dashboard for a particular AppId.

XML

```xml
<root>
  <PluginSettings Enabled="true">
    <Plugins>
      <Plugin
        Name="QuantumPlugin3.0"
        Version=""
        AssemblyName="Quantum.Plugin.Custom.dll"
        Type="Quantum.QuantumCustomPluginFactory"
        PathToLUTFolder="assets/LUT"
        PathToDBFile="assets/db.json"
        EmbeddedDBFile="Quantum.db.json"
      />
    </Plugins>
  </PluginSettings>
</root>

```

The asset paths have to be relative to `Photon.Server\\deploy\\Plugins\\DeterministicPlugin\\bin`.

When upgrading the Plugin SDK local changes should not be overwritten.

Add new properties by adding a line to the Plugin XML node, key and value are both strings and parsed as such.

XML

```xml
<Plugin
  NewStringProperty="foo"
  NewIntProperty="10"
/>

```

C#

```cs
if (config.TryGetString("NewStringProperty", out var newStringProperty, defaultValue: "default")) {
}
if (config.TryParseInt("NewIntProperty", out var newIntProperty, defaultValue: 0)) {
}

```

# Server Classes

## DeterministicPluginFactory

The default factory that instantiates the Quantum server plugins.

## DeterministicPlugin

This class is the default implementation of the Quantum server plugin.

## DeterministicServer

An instance of the DeterministicServer orchestrates a Quantum online game session within a Photon room.

It provides a variety of virtual callback methods to hook into the server flow.

### Methods

#### OnDeterministicServerSetup

C#

```cs
void OnDeterministicServerSetup(IHost host, IEventSender eventSender, IWebhookHost webhookHost, Dictionary<string, string> config, ref Boolean runServerSimulation)

```

This is called when the Realtime room is created.

The `config` dictionary includes the configurations set up in the Photon dashboard when testing online applications or inside the `Photon.Server\\deploy\_win\\LoadBalancing\\GameServer\\bin\\plugin.config` file when testing a local Photon Server.

For ease of use, you can use the `ConfigParsingExtensions` in order to parse additional values.

Example:

C#

```cs
config.TryParseBool("WebHookEnableReplay", out var isReplayStreamingEnabled, false);

```

Set `runServerSimulation` to enabled or disable the server simulation for this session. The bool has already been configured by the dashboard variables `ServerSimulationEnabled` and `ServerSimulationPercent`. The simulation will only start if the server object has been created with an `IDeterministicSessionRunner`.

#### OnDeterministicServerClose

C#

```cs
void OnDeterministicServerClose()

```

Is called when the Realtime room closes.

#### OnDeterministicUpdate

C#

```cs
void OnDeterministicUpdate()

```

Is called for running sessions during each server update after the input has been processed and the server simulation (if any) has been processed.

#### OnDeterministicStartRequest

C#

```cs
Boolean OnDeterministicStartRequest(Protocol.StartRequest startRequestData)

```

This is called when Quantum is requesting to start the simulation. You can return false to deny the request.

#### OnDeterministicGameConfigs

C#

```cs
void OnDeterministicGameConfigs(ref byte[] runtimeConfig, ref DeterministicSessionConfig sessionConfig)

```

This is called when Quantum receives the initial game configurations.

#### OnDeterministicPlayerAdd

C#

```cs
void OnDeterministicPlayerAdd(int playerSlot, ref byte[] runtimePlayer)

```

This is called when a new player is added to the Quantum simulation.

#### OnDeterministicPlayerRemove

C#

```cs
void OnDeterministicPlayerRemove(int playerSlot)

```

This is called when a player is removed from the Quantum simulation.

#### OnDeterministicSnapshotRequested

C#

```cs
Boolean OnDeterministicSnapshotRequested(ref Int32 tick, ref byte[] data)

```

This is called when Quantum requests a snapshot of the current simulation state in order to send it to a late joining client.

You can return false to deny the request.

#### OnDeterministicCommand

C#

```cs
Boolean OnDeterministicCommand(DeterministicPluginClient client, Protocol.Command cmd)

```

This is called when a player sends a command to the Quantum simulation. You can override this method to alter or validate the command.

You can return false to deny the command.

#### OnDeterministicLateStart

C#

```cs
void OnDeterministicLateStart(DeterministicPluginClient client, Protocol.SimulationStart startData)

```

Is called when the late-joining or reconnecting client is about to receive the simulation start event.

#### OnDeterministicInputReceived

C#

```cs
void OnDeterministicInputReceived(DeterministicPluginClient client, DeterministicTickInput input)

```

Is called when a player input is received by the Quantum simulation. You can override this method to alter the input.

#### OnDeterministicInputConfirmed

C#

```cs
void OnDeterministicInputConfirmed(DeterministicPluginClient client, Int32 tick, Int32 playerIndex, DeterministicTickInput input)

```

Is called when input for a client and tick is confirmed.

#### OnDeterministicServerInput

C#

```cs
Boolean OnDeterministicServerInput(DeterministicTickInput input)

```

Override to set the input data for a server controlled player. Tick and player index have already been set on the input object.

#### OnDeterministicStartSession

C#

```cs
void OnDeterministicStartSession()

```

Is called when the Quantum simulation is started.

#### OnDeterministicServerReplacedInput

C#

```cs
void OnDeterministicServerReplacedInput(DeterministicTickInput input)

```

Is called when the server replaces the input for a player.

#### OnDeltaCompressedInput

C#

```cs
void OnDeltaCompressedInput(int tick, byte[] data)

```

This callback is invoked when the delta compressed input if finalized and can be used to stream the replay.

#### SendDeterministicCommand

C#

```cs
void SendDeterministicCommand(Protocol.Command cmd)

```

Sends a command to the Quantum simulation from the server.

Back to top

- [Plugin Factory](#plugin-factory)
- [Basic Structure](#basic-structure)
- [Logging](#logging)
- [Server Simulation](#server-simulation)
  - [Quantum Asset DB](#quantum-asset-db)
  - [Local Quantum Plugin Configuration](#local-quantum-plugin-configuration)

- [Server Classes](#server-classes)
- [DeterministicPluginFactory](#deterministicpluginfactory)
- [DeterministicPlugin](#deterministicplugin)
- [DeterministicServer](#deterministicserver)
  - [Methods](#methods)