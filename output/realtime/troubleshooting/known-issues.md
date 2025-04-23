# known-issues

_Source: https://doc.photonengine.com/realtime/current/troubleshooting/known-issues_

# Known Issues

On this page, we will list known issues with Photon on the various platforms.

**The focus here is on issues which we can't fix or workaround.**

In some cases, this means we will simply list broken versions per platform and guide you to other versions.

### Running in Background

On mobile platforms, if the app moves to background it pauses the main message loop which is responsible for keeping the client connected among other things.

Typical causes of this are:

- Player hits "home button".
- Phone call received.
- Video ads.
- Third party plugin that introduces an overlay view in the app (e.g. Facebook, Google, etc.).

On iOS, applications can not keep a connection, while in background (see [Background Execution](https://developer.apple.com/library/ios/documentation/iPhone/Conceptual/iPhoneOSProgrammingGuide/BackgroundExecution/BackgroundExecution.html#//apple_ref/doc/uid/TP40007072-CH4-SW4) on the Apple dev pages). It might make sense to Disconnect the client, when the app switches to the background.

On WebGL, it might also make sense to set a PlayerTTL and reconnect to the session when the tab is back in focus.

Often, browsers will not run JS and WebAssembly in the background. In some cases, this can be worked aorund by playing audio (even inaudible one) while in the background.

If the app is paused for longer than the client disconnect timeout (10 seconds by default) the client will be disconnected and you need to reconnect as soon as the app is "unpaused".

A "timeout disconnect" callback will be triggered only after the app back from being in the background.

If the game design enables a player to return after several seconds or minutes, then you could also reconnect and rejoin the game again.

If you want to rejoin the same room with same actor number when the app is unpaused you need to take few things into consideration:

- PlayerTTL: the room needs to be created with a PlayerTTL value high enough that permits a player to return after a while.
- EmptyRoomTTL: the room needs to be created with an EmptyRoomTTL value high enough to keep the room alive for a while when the last joined player's app is in the background.

## Unity

### Endless Compile Errors On Import

Some Unity Editor versions create incorrect .sln and .csproj files. It may be temporary. You can regenerate the project via a button in the Preferences "External Tools" panel.

You may want to update or reinstall the Visual Studio Editor package from the Unity Package Manager. This could also solve the project creation issues.

Sometimes assets are not updated properly from the Assset Store as old packages are stuck in the local offline cache.

To fix this, first remove the Photon asset package locally, then try downloading and importing again.

The [paths for the local Unity assets store cache folder are listed here](https://answers.unity.com/questions/690729/how-do-you-remove-packages-already-downloaded-from.html?childToView=793621#answer-793621).

### ArgumentException in Socket.SetSocketOption

There was a [known Unity issue](https://issuetracker.unity3d.com/issues/system-dot-net-dot-socket-objects-throw-argumentexception-in-il2cpp-after-installing-windows-sdk-2004) which caused IL2CPP builds to fail to connect with an "ArgumentException: Value does not fall within the expected range.

at System.Net.Sockets.Socket.SetSocketOption". This happened when the Windows 10 SDK 10.0.19041.0 was installed.

Affected Unity versions: 2018.4.23f1, 2019.4.0f1, 2020.1.0b11, 2020.2.0a13. Many more minor Unity releases have been affected, too.

Fixed versions are: 2020.1.1f1 and 2019.4.5f1 and up. Presumably, 2018.4.27 is also fixed. A [workaround for 2018.4.23 exists as described here](https://forum.unity.com/threads/il2cpp-failing-in-windows-machine.891436/#post-5944052).

### Unity 2018.2 Sockets Freeze with .Net 4.x

Unity 2018.2 used a Mono version, which could freeze communication via sockets. Depending on the message size and frequency, this happened sooner or later.

Eventually 2019.2 got a fix for this and 2018.3 should also have it at some point.

When using Mono and .Net 4.x or .Net Standard 2.0, we recommend to use the 2018.4.x or 2019.4.x releases.

### RunInBackground

Unity's `Application.runInBackground` is not supported on mobile platforms.

Instead, the `OnApplicationPause` method is called whenever the app moves to or back from background:

C#

```csharp
void OnApplicationPause( bool pauseStatus )
{
    if (pauseStatus)
    {
        // app moved to background
    } else
    {
        // app is foreground again
    }
}

```

### IOS App Store Submission Rejects

Sometimes, submissions for the App Store are being rejected by the Apple team, due to connection issues.

We tried to sort this out with Apple and now this is rare but can still happen. Typically, UDP is being blocked in those cases.

Newer Photon clients can automatically fallback to using TCP or WSS, if UDP is not connecting.

You will have to appeal the rejection and in doubt, have to ask for Apple's Developer Support to take over. Their setup usually supports UDP.

### UWP / Windows Store Capabilities

If you target Windows Store (UWP) and you are having exceptions when you try to connect or you have this error:

A network capability is required to access the network!

Make sure to enable the required capability from Unity's "Player Settings" -> "Publisher Settings" -> "Capabilities -> "InternetClient"

![Required Capability for Windows Store Apps](/docs/img/voice/wsa_capabilities.png)
Required Capability for Windows Store Apps. You also need the 'Microphone' capability, if you use Photon Voice.


UWP applications are isolated from others and by design can not connect to a server running on the same Windows instance.

See Microsoft documentation for "AppContainer Isolation".

Back to top

- [Running in Background](#running-in-background)

- [Unity](#unity)
- [Endless Compile Errors On Import](#endless-compile-errors-on-import)
- [ArgumentException in Socket.SetSocketOption](#argumentexception-in-socket.setsocketoption)
- [Unity 2018.2 Sockets Freeze with .Net 4.x](#unity-2018.2-sockets-freeze-with.net-4.x)
- [RunInBackground](#runinbackground)
- [IOS App Store Submission Rejects](#ios-app-store-submission-rejects)
- [UWP / Windows Store Capabilities](#uwp-windows-store-capabilities)