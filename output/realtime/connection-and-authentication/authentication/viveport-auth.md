# viveport-auth

_Source: https://doc.photonengine.com/realtime/current/connection-and-authentication/authentication/viveport-auth_

# VIVEPORT Authentication

## Overview

Adding VIVEPORT as a Photon Authentication Provider is easy.

You will need a VIVEPORT AppId and a VIVEPORT AppSecret. Contact the VIVEPORT content team ( [\[email protected\]](/cdn-cgi/l/email-protection)) to apply for an AppSecret. As this one step is not automated, start the process ahead of time.

## Application Setup

As first step, you should setup VIVEPORT Authentication in the [Photon Applications' Dashboard](https://dashboard.photonengine.com).

Go to the "Manage" page of an application and scroll down to the "Authentication" section.

With one click, you can add or edit the "HTC Vive" authentication provider. You will need:

- **appid**: ID of your VIVEPORT app.
- **appsecret**: Secret for your VIVEPORT app.

The VIVEPORT AppId can be found in the [VIVEPORT Developer Console](https://developer.viveport.com/console). Select your application and open the "VIVEPORT Listing" page.

You have to contact the VIVEPORT content team ( [\[email protected\]](/cdn-cgi/l/email-protection)) to apply for an AppSecret.

## Client Code

Photon verifies VIVEPORT a user with a temporary token, which is provided by the [VIVEPORT API](https://developer.viveport.com).

Download the VIVEPORT SDK and import the Unity package from the zip into your project.

### Get VIVEPORT Token

Before connecting to Photon, the client must to login to VIVEPORT and get a session token.

The following workflow is a summary of the essential steps. The entire code can be seen in VIVEPORTDemo.cs available in the VIVEPORT SDK.

As usual with VIVEPORT, the first step is to init the API with the VIVEPORT AppId.

C#

```csharp
//...
Api.Init(InitStatusHandler, APP_ID);
//...
private static void InitStatusHandler(int nResult)
{
    if (nResult == 0)
    {
        bInit = true;
        bIsReady = false;
        ViveSessionToken = string.Empty;
        bArcadeIsReady = false;
        Viveport.Core.Logger.Log("InitStatusHandler is successful");
    }
    else
    {
        // Init error, close your app and make sure your app ID is correct or not.
        bInit = false;
        Viveport.Core.Logger.Log("InitStatusHandler error : " + nResult);
    }
}
//...

```

After init (e.g. in InitStatusHandler()), check if the token is ready. Call `Viveport.Token.IsReady`:

C#

```csharp
//...
Token.IsReady(IsTokenReadyHandler);
//...
private static void IsTokenReadyHandler(int nResult)
{
    if (nResult == 0)
    {
        bTokenIsReady = true;
        Viveport.Core.Logger.Log("IsTokenReadyHandler is successful");
    }
    else
    {
        bTokenIsReady = false;
        Viveport.Core.Logger.Log("IsTokenReadyHandler error: " + nResult);
    }
}
//...

```

Now the client can get a session token, which is the proof of a valid VIVEPORT user.

C#

```csharp
//...
Token.GetSessionToken(GetSessionTokenHandler);
//...
private static void GetSessionTokenHandler(int nResult, string message)
{
    if (nResult == 0)
    {
        Viveport.Core.Logger.Log("GetSessionTokenHandler is successful, token:" + message);
        // Photon:
        // With the viveport token, we can set the auth values for Photon and connect / auth.
        // We store the token for later use.
        ViveSessionToken = message;
    }
    else
    {
        if (message.Length != 0)
        {
            Viveport.Core.Logger.Log("GetSessionTokenHandler error: " + nResult + ", message:" + message);
        }
        else
        {
            Viveport.Core.Logger.Log("GetSessionTokenHandler error: " + nResult);
        }
    }
}

```

### Authenticate

The ViveSessionToken is the only value that a client needs for Photon VIVEPORT Authentication.

Make sure to use the `CustomAuthenticationType.Viveport` and to set the "userToken" AuthParameter.

C#

```csharp
loadBalancingClient.AuthValues = new AuthenticationValues();
loadBalancingClient.authValues.AuthType = CustomAuthenticationType.Viveport;
loadBalancingClient.AuthValues.AddAuthParameter("userToken", ViveSessionToken);
// do not set loadBalancingClient.AuthValues.Token or authentication will fail
// connect

```

Back to top

- [Overview](#overview)
- [Application Setup](#application-setup)
- [Client Code](#client-code)
  - [Get VIVEPORT Token](#get-viveport-token)
  - [Authenticate](#authenticate)