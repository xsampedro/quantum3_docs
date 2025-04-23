# 2-project-setup

_Source: https://doc.photonengine.com/quantum/current/tutorials/asteroids/2-project-setup_

# 2 - Project Setup

## Overview

Quantum 101 explains the initial steps required to setup a Quantum project.

## Step 0 - Create an Account

As a first step create a Photon account [HERE](https://id.photonengine.com/account/signup) if you don't have an account yet.

## Step 1 - Download SDK

To start using Quantum, download the latest SDK from the [SDK & Release Notes](/quantum/current/getting-started/initial-setup) page.

The requirements to run the SDK are listed on the download page.

## Step 2 - Create an Empty Project

Create an empty project.

**N.B.:** Quantum is agnostic to the Rendering Pipeline chosen; it works with any of them for the sake of this tutorial create an empty 2d project.

![Create an Empty Project](/docs/img/quantum/v3/tutorials/asteroids/2-create-project.png)
Create an Empty Project.
## Step 3 - Importing the Quantum SDK

The SDK is provided as a .unitypackage file and can be imported with the `Assets > Import Package > Custom Package` tool or by opening it in the file explorer. Simply navigate to the location where the SDK was downloaded and trigger the import.

## Step 4 - Quantum Hub

After importing the Quantum SDK, the Quantum Hub window opens. If the window does not open automatically it can be manually opened in the Unity menu under Quantum > Quantum Hub.

![Quantum Hub](/docs/img/quantum/v3/tutorials/asteroids/2-quantum-hub.png)
Quantum Hub.


On the first page press the `Install Quantum User Files` button. This generates default files for the configuration and settings files Quantum uses.

## Step 5 Create and Link a Quantum AppID

To be able to run your Quantum application online and not just for local offline development, you need to link it to an App running with a Quantum Plugin. For this you simply go to your [Dashboard](https://dashboard.photonengine.compubliccloud) and hit the `Create a New App` button in the top right corner.

![Photon Dashboard](/docs/img/quantum/v3/tutorials/asteroids/2-dashboard-create-app-id.png)
Photon Dashboard


In the creation menu, select `Quantum` in the drop-down menu called `Select Photon SDK`.

![Select Quantum](/docs/img/quantum/v3/tutorials/asteroids/2-dashboard-chose-quantum-type.png)
Select Quantum


Select `Quantum 3` on the `Select SDK Version` dropdown that appeared and fill out the rest of the form and click on `Create`.

![Select Quantum 3](/docs/img/quantum/v3/tutorials/asteroids/2-dashboard-chose-quantum-type2.png)
Select Quantum 3


The new Quantum App will now appear in your dashboard. In order to connect your application, you will need to copy the _App ID_ of the App. Clicking on the field will reveal your unique _App ID_ and allow you to copy it.

![](/docs/img/quantum/v3/tutorials/asteroids/2-quantum-app-example.png)

![](/docs/img/quantum/v3/tutorials/asteroids/2-quantum-app-id.png)

Return to the Quantum Hub and paste your _APP ID_ into the App Id field. With that done, the project is now ready for development!

![App Id Field](/docs/img/quantum/v3/tutorials/asteroids/2-app-id.png)
The App Id field in the Quantum hub.
## Step 6 - Navigating the Quantum Project

Before starting the development let's have a quick look at the files we imported.

Quantum is a game engine and thus more intertwined that most SDKs. All files are imported into the `Photon` and `QuantumUser` folder and follow a strict folder structure. It is advised to not change the folder structure or move any Quantum files to a different place.

Find [HERE](/quantum/current/manual/quantum-project) more information on the folders contains in the Quantum SDK.

Back to top

- [Overview](#overview)
- [Step 0 - Create an Account](#step-0-create-an-account)
- [Step 1 - Download SDK](#step-1-download-sdk)
- [Step 2 - Create an Empty Project](#step-2-create-an-empty-project)
- [Step 3 - Importing the Quantum SDK](#step-3-importing-the-quantum-sdk)
- [Step 4 - Quantum Hub](#step-4-quantum-hub)
- [Step 5 Create and Link a Quantum AppID](#step-5-create-and-link-a-quantum-appid)
- [Step 6 - Navigating the Quantum Project](#step-6-navigating-the-quantum-project)