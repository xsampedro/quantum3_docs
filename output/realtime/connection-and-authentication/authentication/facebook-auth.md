# facebook-auth

_Source: https://doc.photonengine.com/realtime/current/connection-and-authentication/authentication/facebook-auth_

# Facebook Authentication

## Server Side

### Facebook App Setup

First, we need to create a Facebook application if you don't have one already.

- Go to [Facebook Developers](https://developers.facebook.com) website.
- Click _Apps_ -\> _Create a New App_ , enter the name of your app and press _Create App_ button.
- Choose _Apps_ -\> _\[your\_app\]_ and copy **App ID** and **App Secret**.

### Photon Configuration

Go to the _Details_ page of your application via your [Photon Cloud Dashboard](https://dashboard.photonengine.com).

Expand the Custom Authentication section.

The _Authentication URL_ for Facebook authentication is set by Exit Games.

Set the values for these two parameters:

- appid = _your Facebook App ID_
- secret = _your Facebook App Secret_

Save changes.

## Client Code

The client needs to set the correct authentication type (Facebook, 2) and send a valid Facebook token as a query string parameter named "token".

### PUN

- Open Unity.
- [Import](http://u3d.as/1cMM).
- [Setup PUN](/pun/v2/getting-started/initial-setup).
- Import [Facebook SDK for Unity](https://developers.facebook.com/docs/unity/).
- In Unity's main menu go to _Facebook_ -\> _Edit Settings_, enter the _App Name_ and _App Id_ for your Facebook Application.

## Implementation

Create a new MonoBehaviour, attach it to an object on scene and then open.

Use the following code for Facebook initialization and login:

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

To use Facebook Authentication in PUN, add:

C#

```csharp
private void OnFacebookLoggedIn()
{
    // AccessToken class will have session details
    string aToken = AccessToken.CurrentAccessToken.TokenString;
    string facebookId = AccessToken.CurrentAccessToken.UserId;
    PhotonNetwork.AuthValues = new AuthenticationValues();
    PhotonNetwork.AuthValues.AuthType = CustomAuthenticationType.Facebook;
    PhotonNetwork.AuthValues.UserId = facebookId; // alternatively set by server
    PhotonNetwork.AuthValues.AddAuthParameter("token", aToken);
    PhotonNetwork.ConnectUsingSettings("1.0");
}

```

The PUN callbacks for success and error that you can implement are:

C#

```csharp
public class FacebookAuthTest : MonoBehaviourPunCallbacks
{
    public override void OnConnectedToMaster()
    {
        Debug.Log("Successfully connected to Photon!");
    }
    // something went wrong
    public override void OnCustomAuthenticationFailed(string debugMessage)
    {
        Debug.LogErrorFormat("Error authenticating to Photon using facebook: {0}", debugMessage);
    }
}

```

Back to top

- [Server Side](#server-side)

  - [Facebook App Setup](#facebook-app-setup)
  - [Photon Configuration](#photon-configuration)

- [Client Code](#client-code)

  - [PUN](#pun)

- [Implementation](#implementation)