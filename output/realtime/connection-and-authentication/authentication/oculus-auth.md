# oculus-auth

_Source: https://doc.photonengine.com/realtime/current/connection-and-authentication/authentication/oculus-auth_

# Oculus Authentication

## Application Setup

Adding Oculus as an authentication provider is easy and it could be done in few seconds from your [Photon Applications' Dashboard](https://dashboard.photonengine.com).

Go to the "Manage" page of an application and scroll down to the "Authentication" section.

If you add a new authentication provider for Oculus or edit an existing one, here the mandatory settings:

- **appid**: ID of your Oculus App.
- **appsecret**: secret for your Oculus App.

## Client Code

Oculus verifies users based on their Oculus ID and a client-provided nonce.

In cryptography, a nonce is an arbitrary number that can only be used once.

Read more [here](https://developer.oculus.com/documentation/platform/latest/concepts/dg-ownership/).

### Get Credentials

Client needs to log in to Oculus then generate a nonce.

This nonce is proof that the client is a valid Oculus user.

You can get Oculus SDKs from [their website](https://developer.oculus.com).

#### Unity Instructions

Download Oculus Platform SDK for Unity and import it to your project.

From the Editor's menu bar, go to "Oculus Platform" -> "Edit Settings" and enter your Oculus AppId.

Use the following code to get the logged in user's Oculus ID and generate a nonce:

C#

```csharp
using UnityEngine;
using Oculus.Platform;
using Oculus.Platform.Models;
public class OculusAuth : MonoBehaviour
{
    private string oculusId;
    private void Start()
    {
        Core.AsyncInitialize().OnComplete(OnInitializationCallback);
    }

    private void OnInitializationCallback(Message<PlatformInitialize> msg)
    {
        if (msg.IsError)
        {
            Debug.LogErrorFormat("Oculus: Error during initialization. Error Message: {0}",
                msg.GetError().Message);
        }
        else
        {
            Entitlements.IsUserEntitledToApplication().OnComplete(OnIsEntitledCallback);
        }
    }
    private void OnIsEntitledCallback(Message msg)
    {
        if (msg.IsError)
        {
            Debug.LogErrorFormat("Oculus: Error verifying the user is entitled to the application. Error Message: {0}",
                msg.GetError().Message);
        }
        else
        {
            GetLoggedInUser();
        }
    }
    private void GetLoggedInUser()
    {
        Users.GetLoggedInUser().OnComplete(OnLoggedInUserCallback);
    }
    private void OnLoggedInUserCallback(Message<User> msg)
    {
        if (msg.IsError)
        {
            Debug.LogErrorFormat("Oculus: Error getting logged in user. Error Message: {0}",
                msg.GetError().Message);
        }
        else
        {
            oculusId = msg.Data.ID.ToString(); // do not use msg.Data.OculusID;
            GetUserProof();
        }
    }
    private void GetUserProof()
    {
        Users.GetUserProof().OnComplete(OnUserProofCallback);
    }
    private void OnUserProofCallback(Message<UserProof> msg)
    {
        if (msg.IsError)
        {
            Debug.LogErrorFormat("Oculus: Error getting user proof. Error Message: {0}",
                msg.GetError().Message);
        }
        else
        {
            string oculusNonce = msg.Data.Value;
            // Photon Authentication can be done here
        }
    }
}

```

### Authenticate

Client needs to send the Oculus ID and generated nonce as query string parameters with the respective keys "userid" and "nonce":

C#

```csharp
loadBalancingClient.AuthValues = new AuthenticationValues();
loadBalancingClient.AuthValues.UserId = oculusId;
loadBalancingClient.AuthValues.AuthType = CustomAuthenticationType.Oculus;
loadBalancingClient.AuthValues.AddAuthParameter("userid", oculusId);
loadBalancingClient.AuthValues.AddAuthParameter("nonce", oculusNonce);
// do not set loadBalancingClient.AuthValues.Token or authentication will fail
// connect

```

Back to top

- [Application Setup](#application-setup)
- [Client Code](#client-code)
  - [Get Credentials](#get-credentials)
  - [Authenticate](#authenticate)