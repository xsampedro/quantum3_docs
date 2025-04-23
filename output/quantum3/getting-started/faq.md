# faq

_Source: https://doc.photonengine.com/quantum/current/getting-started/faq_

# Frequently Asked Questions

Most of the FAQ content is derived from questions asked by developers in the ```
#quantum-sdk-v2
```

 and ```
#beginner-questions
```

Discord chat. Also use the search functionality there to find more help.

## Quantum Unity Framework

### Why is the Game scene not loading when starting the game via the demo menu?

When the AppId is rejected by the server there is no proper error message, yet. Make sure you have created your AppId correctly:

1. Create a Quantum App Id on [your Photon dashboard](https://dashboard.photonengine.com)
2. Click CREATE A NEW APP
3. Set Photon Type to **Photon Quantum**
4. Fill out the name field
5. Scroll down and press **CREATE**

### Why is there a 1 second delay before the Quantum game starts?

Check the ```
Room Wait Time (seconds)
```

 in your ```
Deterministic Config
```

asset. There you can tweak it. The _Room Wait Time_ is used to cope with fluctuating ping times.

It will always be used in full to help sync all client sessions at start. If you set it to 1 second, it will always wait the full 1 second.

**N.B.:** This should NOT be used to synchronize scene loading times! If you need to coordinate this, load the scene before starting the Quantum session and coordinate through Photon Realtime directly. In this scenario, the value can be set to 0 safely.

### Why is the connection timing out while trying to connect to ns.exitgames.com?

To get more details increase the connection log level ```
Network Logging
```

 in your ```
PhotonServerSettings
```

. Most often, the timeout is due to your UDP traffic being blocked.

For guide on how to analyze your situation, please refer to the **Analyzing Disconnects** page in the Photon Realtime documentation: [Analyzing Disconnects](/realtime/current/troubleshooting/analyzing-disconnects "Analyzing Disconnects")

### How can I debug a frozen game?

Adding the ```
QUANTUM\_STALL\_WATCHER\_ENABLED
```

 define to ```
Project Settings > Player > Scripting Define Symbols
```

will enable a watcher script to spin up a thread which will watch the Update loop. If a stall is detected (i.e. the Update takes more than X seconds to be called again), it will create a crash. This is useful when debugging freezes on the simulation as the generated crash should have the call stack from all threads running.

### Is there an effective way to simulate network latency when running the game in Unity Editor?

The Quantum performance profiler has a latency simulation on board.

Download here: [AddOns \| Profiler](/quantum/current/manual/profiling#quantum_graph_profiler)

Or use an external network throttling tool for example like Clumsy (Windows) and filter out the game server ports:

- UDP 5056
- TCP 4531

```
Clumsy Filter: (udp.DstPort == 5056 or udp.SrcPort == 5056) or (tcp.DstPort == 4531 or tcp.SrcPort == 4531)

```

### Why is the game simulation running faster after pausing or debugging a breakpoint?

By default the time is measured internally and does not compensate for halting the simulation. When ```
DeltaTimeType
```

 on the SimulationConfig is changed to ```
EngineDeltaTime
```

the gameplay will resume in regular speed after a pause. Caveat: Changing it will make every client use the setting which might be undesirable when only used for debugging. Although some games with very tight camera controls (e.g. flight simulation) will profit from setting it to ```
EngineDeltaTime
```

.

C#

```csharp
public enum SimulationUpdateTime {
Default = 0, // internal clock
EngineDeltaTime = 1, // Time.deltaTime (Unity)
EngineUnscaledDeltaTime = 2 // Time.unscaledDeltaTime
}

```

### Why are Navmesh islands being generated inside my geometry?

This will at some point be fixed by Unity (see forum post: [forum.unity.com/threads/nav-generating-inside-non-walkable-objects](https://forum.unity.com/threads/nav-generating-inside-non-walkable-objects.445177))

The workaround with NavmeshModifierVolume will mitigate this (requires [NavMeshComponents](https://github.com/Unity-Technologies/NavMeshComponents)).

Triangle culling during the import step is another alternative which we could provide in the future.

![Navmesh Island](/docs/img/quantum/v2/getting-started/faq/navmesh-island.png)### Why does my game gets a timeout disconnection when the Scene loading takes too long?

When loading a Unity scene, even if it is done with ```
LoadSceneAsync
```

, the main thread can freeze for some time accordingly to the size and complexity of such scene. This can then result in a disconnect error due to timeout as the communication doesn't happen while the game is frozen.

To prevent this from happening, you can use some of the API provided in the ```
ConnectionHandler
```

 class. Here is a step-by-step on how to setup and use it:

- Check if there is any Game Object with the ```
  ConnectionHandler
  ```

   component. If there is none, please add one;

- On the component, you'll be able to see a field named ```
  KeepAliveInBackground
  ```

  , which you can use to increase the time that the connection will be kept. The value is informed in milliseconds;

- You should now inform what is the ```
  QuantumLoadBalancingClient
  ```

  , to which there is a static getter on ```
  UIMain
  ```

  , in case you use it (it comes by default with Quantum). Once you have done this, you can start the ```
  StartFallbackSendAckThread
  ```

  . Here is a sample snippet on how to achieve that:


C#

```csharp
// Before starting loading the scene
if (\_connectionHandler != null)
{
\_connectionHandler.Client = UIMain.Client;
\_connectionHandler.StartFallbackSendAckThread();
}

```

### IL2CPP compilation fails with bracket nesting level exceeded maximum

IL2CPP compilation throws the following error:

```
bracket nesting level exceeded maximum
```

Which can be caused by for example by large classes or structs like ```
RuntimeConfig
```

or ```
GetHashCode()
```

methods of DSL generated components.

To workaround this issue split up the large data into smaller structs.

## Developing A Simulation In Quantum

### Why is the simulation starting with frame 60?

This is the number of rollback-able frames that are allocated at the start of the simulation. The number will match the value set up in DeterministicConfig->Rollback Window.

### Why is Update() called multiple times with the same frame number?

When running Quantum online, Update() on systems are called multiple times for the same frame number in case of a rollback. This happens when the prediction of a remote players input was detected to be incorrect and the simulation has to re-run a frame with the correct input data to get back into a deterministic state.

It is possible though to check/log "frame.IsVerified" in order to check if a frame is verified. When running Quantum in offline mode, there are no duplicates due to rollbacks not happening in such situation.

### What's the difference between FP.MaxValue and FP.UseableMax?

The fixed point math only uses 16+16 bits of its 64-bit value. This makes part of the math faster because we don't have to check for overflows. That said: ```
FP.MinValue
```

and ```
FP.MaxValue
```

are using all 64 bits and should **never** be used for calculations. Use ```
FP.UseableMax
```

and ```
FP.UseableMin
```

instead (to initialize a distance variable with the min FP value for example).

**N.B.:** FP can represent values from -32,768 to 32,768 (-2¹⁵ to 2¹⁵).

### Why does the pointer to a new struct point to stale data?

Inside the loop the pointer to the struct gets the same stack pointer and will contain stale data if **neither**```
new
```

or ```
default
```

was used.

Qtn

```cs
struct Bar {
public bool Foo;
}

static unsafe void Main(string\[\] args) {
for (int i = 0; i < 2; i++) {

Bar bar;
//Bar bar = default(Bar); // <---- Fixes the stale data

Bar\* barPt = &bar;
if (barPt->Foo)
Console.WriteLine("Stuff and Things");

barPt->Foo = true;
}

Console.ReadKey();
}

```

### Why is my simulation desync-ing?

When ```
DeterministicConfig.ChecksumInterval
```

is > 0 a checksum of a verified frame is computed, sent to the server and compared with checksums that other clients have sent.

Most common causes are:

#### Writing to Quantum data assets

C#

```csharp
var characterSpecAsset = frame.FindAsset<CharacterSpec>("WhiteFacedBarghast");
characterSpecAsset.RemainigLifetime = 21;

```

**Never** write anything to the Quantum assets. They contain read-only data.

#### Writing to Quantum from Unity thread

All scripts in Unity have read-only access to everything that is exposed through the Quantum Frame. Only influence the simulation by Input and/or Commands.

#### Caching data

The code snippet presented below will desync eventually when the simulation is rolled-back.

C#

```csharp
public class CleaningSystem : SystemBase {
 public Boolean HasShoweredToday; // <----- Error
 public override void Update(Frame frame) {
 if (!HasShoweredToday && frame.Global->ElapsedTime > 100) {
 Shower();
 HasShoweredToday = true;
 }
 }
}

```

Instead save non-transient data on the Frame or on Entity Components.

C#

```csharp
// Frame
unsafe partial class Frame {
 public Boolean HasShoweredToday;
 partial void CopyFromUser(Frame frame) {
 // Implement copy of the custom parameters.
 }
}

public class CleaningSystem : SystemBase {
 public override void Update(Frame frame) {
 if (!frame.HasShoweredToday && frame.Global->ElapsedTime > 100) {
 Shower();
 frame.HasShoweredToday = true;
 }
 }
}

```

#### Floating point math

Refrain from using floats inside the simulation and exclusively use ```
FP
```

 math.

Handle ```
FP.FromFloat\_UNSAFE()
```

with care. Using it "offline" for balancing asset generation on one machine is fine; however, be aware this can return different results on different platforms. If you need to import floats during run-time and cannot use integers or FPs (e.g. downloading balancing data), convert from **String to FP**.

#### Data Created During AssetObject.Loaded()

```
AssetObject.Loaded()
```

is called once per asset during loading. It is totally fine to store, calculate and store new data inside the asset members at this time - **N.B.:** If you are running the simulation on the server, all games will share this one asset.

If your assets are loaded from Resources and you are either restarting the Unity Editor or resetting the Unity DB at runtime, be aware that Unity does not unload the Unity asset.

C#

```csharp
public partial class FooAsset {
 public Foo Settings;
 public int RawData;

 \[NonSerialized\]
 public List<int> Bar = new List<int>();

 public override AssetObject AssetObject => Settings;

 public override void Loaded(IResourceManager resourceManager, Native.Allocator allocator)
 {
 base.Loaded(resourceManager, allocator);
 // This will break on the second run (see above) because Bar needs to be reset by either Bar.Clear() or Bar = new List<int>()
 Bar.Add(RawData);
 }
}

```

### Can I reuse the Photon Room for new Quantum sessions with the same clients?

**No**, you should not reuse the room for another Quantum session.

**But** you can keep the players (or a portion of them) inside the room and soft-restart the simulation. This is how Stumble Guys progresses between with the three stages.

- Keep the Quantum session running even if your game round has ended.
- Add code to your gameplay systems (e.g. a game state machine) that handles the starting of a game rounds.
- Deterministically disable Quantum systems and/or load a new Quantum map.
- Reset or destroy all Quantum entities, reset your game view.
- Players that lost a round can keep the connection to spectate the match but can't influence the game anymore.

**Alternatively** all players can switch the room. Because this involves server transitions and with distributed systems a lot of things can fail it is recommended to try to soft-restart the simulation and only use room transition as a last resort.

- Share the id of the new room between the clients using (e.g. Photon room properties or a Quantum command). Or create the new room programmatically/deterministically (be careful that this cannot be guessed easily).
- Every client stops the Quantum session and runs ```
LeaveRoom()
```

but does not disconnect.
- All clients use ```
JoinOrCreate()
```

to connect to the new room.
- Mitigate connection problems like clients restarting their app during the process, other connection errors, waiting for players to join, etc

[Photon Realtime Matchmaking Guide](/realtime/current/lobby-and-matchmaking/matchmaking-and-lobby)

### Can I use DSL generated union structs in Quantum assets?

Not directly. Unity does not support overlapping fields.

Instead ```
<UnionName>\_Prototype
```

 can be used as it is Unity-serializable and has a drawer.

To convert to a ```
union
```

-struct like this:

C#

```csharp
UnionName result = default;
prototype.Materialize(frame, ref result, default);

```

## Billing

#### Do you have special offers for students, hobbyists or indies?

No, but all our products have a free tier and a one-off entry-plan.

We also usually take part in Unity's asset store sales and occasionally give vouchers to lucky ones.

#### Can I combine more than one 100 CCU plan for a single Photon application?

No.

The 100 CCU plans are not stackable with each other. Only one can be applied per AppId.

They can be combined with paid subscriptions and contribute to the CCU for their duration.

They do not automatically renew.

In case you need more CCU for a single app, the next higher plan is the 500 CCU one.

If you subscribe to a monthly or yearly plan, then you will still keep the 100 CCUs for 12 months on top of / in addition to the CCU from your monthly/yearly plan.

One exception to this rule are the Free100 plans for Quantum and Fusion. These can be combined with 100 CCU coupons from the Unity Asset store (Quantum Plus and Fusion Plus).

#### How much traffic is included in my Photon plan and what happens if my app generates traffic beyond the included limit?

Photon Public Cloud and Premium Cloud plans include 3GB per CCU.

For example, a monthly plan with 1,000 CCUs includes 3 TB of traffic per month.

If your app generates more traffic, we will automatically send a heads up via email. You will receive an automatically generated overage invoice at the end of the month at your Photon account email address. The invoice amount is based on the following calculation:

> Total traffic - included traffic = overage traffic (in GB)

Traffic is calculated with $0.05 / $0.10 per GB depending on the Photon Cloud region used. This invoice is automatically charged to your credit card on file.

#### What happens if the Peak CCU exceeds the booked CCU of my Photon Cloud plan?

If you subscribed to a 500 CCU / 1,000 CCU / 2,000 CCU plan, “CCU burst” is automatically activated for your application. The Photon Cloud will allow more CCU than you booked to deliver the best possible experience for your users.

Once the burst kicks in, you are obliged to upgrade to the required subscription tier within 48 hours, according to the terms agreed.

If you do not upgrade, we will send an “overage invoice” to your Photon account email address and charge each CCU above your subscribed plan with a fee of $0.75 / $1.00 (based on the used SDK) per CCU. This invoice is automatically charged to your credit card on file.

Photon charges the “Peak CCU”, which is the sum of the peak CCU per region added up in a given month. Please make sure to upgrade even if the usage decreases after reaching its peak to avoid overage charges.

Back to top

- [Quantum Unity Framework](#quantum-unity-framework)

  - [Why is the Game scene not loading when starting the game via the demo menu?](#why-is-the-game-scene-not-loading-when-starting-the-game-via-the-demo-menu)
  - [Why is there a 1 second delay before the Quantum game starts?](#why-is-there-a-1-second-delay-before-the-quantum-game-starts)
  - [Why is the connection timing out while trying to connect to ns.exitgames.com?](#why-is-the-connection-timing-out-while-trying-to-connect-to-ns.exitgames.com)
  - [How can I debug a frozen game?](#how-can-i-debug-a-frozen-game)
  - [Is there an effective way to simulate network latency when running the game in Unity Editor?](#is-there-an-effective-way-to-simulate-network-latency-when-running-the-game-in-unity-editor)
  - [Why is the game simulation running faster after pausing or debugging a breakpoint?](#why-is-the-game-simulation-running-faster-after-pausing-or-debugging-a-breakpoint)
  - [Why are Navmesh islands being generated inside my geometry?](#why-are-navmesh-islands-being-generated-inside-my-geometry)
  - [Why does my game gets a timeout disconnection when the Scene loading takes too long?](#why-does-my-game-gets-a-timeout-disconnection-when-the-scene-loading-takes-too-long)
  - [IL2CPP compilation fails with bracket nesting level exceeded maximum](#il2cpp-compilation-fails-with-bracket-nesting-level-exceeded-maximum)

- [Developing A Simulation In Quantum](#developing-a-simulation-in-quantum)

  - [Why is the simulation starting with frame 60?](#why-is-the-simulation-starting-with-frame-60)
  - [Why is Update() called multiple times with the same frame number?](#why-is-update-called-multiple-times-with-the-same-frame-number)
  - [What's the difference between FP.MaxValue and FP.UseableMax?](#whats-the-difference-between-fp.maxvalue-and-fp.useablemax)
  - [Why does the pointer to a new struct point to stale data?](#why-does-the-pointer-to-a-new-struct-point-to-stale-data)
  - [Why is my simulation desync-ing?](#why-is-my-simulation-desync-ing)
  - [Can I reuse the Photon Room for new Quantum sessions with the same clients?](#can-i-reuse-the-photon-room-for-new-quantum-sessions-with-the-same-clients)
  - [Can I use DSL generated union structs in Quantum assets?](#can-i-use-dsl-generated-union-structs-in-quantum-assets)

- [Billing](#billing)