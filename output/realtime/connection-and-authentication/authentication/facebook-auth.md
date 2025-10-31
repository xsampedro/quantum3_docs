# facebook-auth

_Source: https://doc.photonengine.com/realtime/current/connection-and-authentication/authentication/facebook-auth_

# Facebook Authentication

For Facebook Authentication, clients will first use Meta's API to authenticate the user, then send the resulting secret token to Photon to confirm and use their identity.

This page describes the setup and workflow.

## Server Side

### Facebook App Setup

First, you need to create a Facebook application if you don't have one already.

- Go to the [Facebook Developers](https://developers.facebook.com) website and login.
- Click _Apps_ -\> _Create a New App_ , enter the name of your app and press _Create App_ button.
- Choose _Apps_ -\> _\[your\_app\]_ and copy **App ID** and **App Secret**.

### Photon Configuration

Open the [Photon Cloud Dashboard](https://dashboard.photonengine.com/), and find the application to setup.

Click "Manage" and find the "Authentication" section, which shows the currently configured authentication providers.

Add a provider for "Facebook" and fill in the two required parameters:

- appid = this refers to _your Facebook App ID_
- secret = _your Facebook App Secret_

Unlike for Custom Authentication, the _Authentication URL_ is fixed and set automatically.

Save changes.

## Client Side

On the client side, apps first have to authenticate via Facebook to get a token, which is sent to the Photon server in the AuthenticationValues.

In Unity, you will need to:

- Import [Facebook SDK for Unity](https://developers.facebook.com/docs/unity/).
- In Unity's main menu go to _Facebook_ -\> _Edit Settings_, enter the _App Name_ and _App Id_ for your Facebook Application.

## Implementation

Create a new MonoBehaviour, attach it to an object on scene and add the following code for Facebook initialization and login:

C#

```csharp
// Include Facebook namespace
using Facebook.Unity;
// [..]
private void Awake()
{
    if (!FB.IsInitialized)
    {
        // Initialize the Facebook SDK
        FB.Init(InitCallback);
    }
    else
    {
        FacebookLogin();
    }
}
private void InitCallback()
{
    if (FB.IsInitialized)
    {
        FacebookLogin();
    }
    else
    {
        Debug.Log("Failed to initialize the Facebook SDK");
    }
}
private void FacebookLogin()
{
    if (FB.IsLoggedIn)
    {
        OnFacebookLoggedIn();
    }
    else
    {
        var perms = new List<string>(){"public_profile", "email", "user_friends"};
        FB.LogInWithReadPermissions(perms, AuthCallback);
    }
}
private void AuthCallback(ILoginResult result)
{
    if (FB.IsLoggedIn)
    {
        OnFacebookLoggedIn();
    }
    else
    {
        Debug.LogErrorFormat("Error in Facebook login {0}", result.Error);
    }
}

```

`OnFacebookLoggedIn()` gets called when the authentication with Facebook is complete and the user's Facebook token is available.

The example below shows how to set up the AuthenticationValues:

C#

```csharp
private void OnFacebookLoggedIn()
{
    // Facebook's AccessToken class will have session details
    string aToken = AccessToken.CurrentAccessToken.TokenString;
    string facebookId = AccessToken.CurrentAccessToken.UserId;
    var authValues = new AuthenticationValues();
    authValues.AuthType = CustomAuthenticationType.FacebookGaming;  // specifically for gaming apps
    authValues.UserId = facebookId;
    authValues.AddAuthParameter("token", aToken);
    // set the new authValues on the client / runner and connect
    // example for Realtime API
    client.AuthValues = authValues;
    client.ConnectUsingSettings();
}

```

If your app does **not use the Facebook Gaming APIs** (described on https://developers.facebook.com/docs/games), the `AuthType` need to be set to `CustomAuthenticationType.Facebook`.

From here on, the usual connection workflow happens. This differs a little, depending on the used Photon SDK.

The Realtime API in C# has the `IConnectionCallbacks` and you want to check `OnConnectedToMaster()` as well as `OnCustomAuthenticationFailed(string debugMessage)`.

C#

```csharp
    public void OnConnectedToMaster()
    {
        Debug.Log("Successfully connected to Photon!");
    }
    // something went wrong, check the setup of the Facebook app, login, token, etc.
    public void OnCustomAuthenticationFailed(string debugMessage)
    {
        Debug.LogErrorFormat("Error authenticating to Photon using facebook: {0}", debugMessage);
    }

```

Back to top

- [Server Side](#server-side)

  - [Facebook App Setup](#facebook-app-setup)
  - [Photon Configuration](#photon-configuration)

- [Client Side](#client-side)
- [Implementation](#implementation)