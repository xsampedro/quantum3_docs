# unreal-engine

_Source: https://doc.photonengine.com/realtime/current/getting-started/unreal-engine_

# Unreal Engine

All C++ based Photon SDKs are compatible with the UNREAL ENGINE off the shelf.

In detail the multiplayer SDKs you can gear up with the UNREAL ENGINE SDK are

- [Android (NDK)](https://www.photonengine.com/sdks#android-cpp)
- [Emscripten (HTML5)](https://www.photonengine.com/sdks#emscripten-cpp)
- [iOS](https://www.photonengine.com/sdks#ios-cpp)
- [Linux](https://www.photonengine.com/sdks#linux-cpp)
- [Mac OS X](https://www.photonengine.com/sdks#macosx-cpp)
- [Nintendo Switch)](https://www.photonengine.com/sdks#nintendoswitch-cpp)
- [PS4](https://www.photonengine.com/sdks#playstation4-cpp)
- [PS5](https://www.photonengine.com/sdks#playstation5-cpp)
- [tvOS](https://www.photonengine.com/sdks#tvos-cpp)
- [UWP](https://www.photonengine.com/sdks#windowsstore-cpp)
- [Windows Desktop](https://www.photonengine.com/sdks#windows-cpp)
- [Xbox One / Xbox Series S/X GDK](https://www.photonengine.com/sdks#gamecore-cpp)
- [Xbox One XDK](https://www.photonengine.com/sdks#xboxone-cpp)

[You get all Photon SDKs from the Download page.](https://www.photonengine.com/sdks)

## Get Started

Proceed as follows to integrate any of the compatible Photon multiplayer game SDKs with the UNREAL ENGINE SDK.

1. [Download the UNREAL ENGINE SDK](https://www.unrealengine.com/).

2. [Download the Photon SDK for your targeted Platform](https://www.photonengine.com/sdks#unrealengine).

3. Only Unreal 'C++' projects are supported.

4. Unpack the Photon SDK of your choice (Windows, Android or iOS) in the 'Photon' folder inside of the 'Source' folder of your Unreal project.


Only the header files and pre-built libraries are required. You may want to add libraries for multiple different platforms.


Sample folders layout:

text
```text
    \---Source
           +---Photon
           |    +---Common-cpp
           |    |    \---inc
           |    |        (*.h)
           |    +---LoadBalancing-cpp
           |    |    \---inc
           |    |        (*.h)
           |    |---Photon-cpp
           |    |    \---inc
           |    |        (*.h)
           |    +---lib
           |    |    +---Android
           |    |        (*.a)
           |    |    +---iOS
           |    |        (*.a)
           |    |    \---Windows
           |    |        (*.lib)

```

5. Edit the "\*.Build.cs" file of the Unreal project to load the libraries for a given platform and to set Photons platform defines.


See "Source/PhotonDemoParticle/PhotonDemoParticle.Build.cs" in the demo that is linked below and Unreal documentation:

C#
```csharp
private string PhotonPath
{
       get { return Path.GetFullPath(Path.Combine(ModulePath, "..", "Photon")); }
}
//
if ( Target.Platform == UnrealTargetPlatform.Android)
{
       // Set _EG_WINDOWS_PLATFORM for Windows, _EG_IPHONE_PLATFORM for iOS and _EG_IMAC_PLATFORM for OS X
       Definitions.Add("_EG_ANDROID_PLATFORM");
       //
       PublicAdditionalLibraries.Add(Path.Combine(PhotonPath, "lib", "Android", "libcommon-cpp-static_debug_android_armeabi_no-rtti.a"));
       PublicAdditionalLibraries.Add(Path.Combine(PhotonPath, "lib", "Android", "libphoton-cpp-static_debug_android_armeabi_no-rtti.a"));
       PublicAdditionalLibraries.Add(Path.Combine(PhotonPath, "lib", "Android", "libloadbalancing-cpp-static_debug_android_armeabi_no-rtti.a"));
}

```


Set \_EG\_WINDOWS\_PLATFORM for Windows, \_EG\_IPHONE\_PLATFORM for iOS, \_EG\_ANDROID\_PLATFORM for Android, and so on.

The define for each platform can be found in Common-cpp/inc/platform\_definition.h inside the Client SDK for the respective platform

6. Use the imported Photon API in the source code of your project.

7. Build your Unreal project for your platform of choice.


## Notes

Some hints regarding Unreal iOS builds can be found [here](https://answers.unrealengine.com/questions/21222/steps-for-ios-build-with-unrealremotetool.html).

## Ready-to-run Demo

Find a ready-to-run proof of concept for download [here](https://dashboard.photonengine.com/download/photon-unreal-sdk_demoparticle-ue.zip)

- Follow steps 1 and 2 from above
- Unzip the downloaded package
- Follow step 4 from above
- Open the context menu for "./PhotonDemoParticle.uproject" and choose "Generate Visual Studio project files"
- If you have multiple different versions of unreal engine installed, choose the desired engine version and click OK
- UE is now generating the project files. This may take a couple of seconds and UE will indicate that it has finished doing so only by letting the "Generating" message box disappear
- Open "PhotonDemoParticle.sln" with Visual Studio
- Choose "Win64" as solution platform
- Choose "DebugGame\_Editor" as solution configuration
- In the solution explorer navigate to Games/PhotonDemoParticle and build that project
- Debug or Run that VS project - this will start the UE editor with PhotonDemoParticle as loaded UE project
- In the UE Editor "World Outliner" tab navigate to "PhotonLBClient" -> "Demo" -> "App ID"
- Replace the content of that field with the App ID from your Dashboard on our Website
- Press Play

Back to top

- [Get Started](#get-started)
- [Notes](#notes)
- [Ready-to-run Demo](#ready-to-run-demo)