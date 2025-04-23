# epic-eos-auth

_Source: https://doc.photonengine.com/realtime/current/connection-and-authentication/authentication/epic-eos-auth_

# Epic EOS Authentication

## Application Setup

Adding Epic / EOS as authentication provider is easy and it could be done in few seconds from your [Photon Applications' Dashboard](https://dashboard.photonengine.com).

Go to the "Manage" page of an application and scroll down to the "Authentication" section.

### Mandatory Settings

- **clientid**: Client ID used to authenticate the user with Epic Account Services (used for ID Token validation). Hint: this value can be (currently) found in the EOS dashboard > Product Settings > Clients ( > Details > Client ID)
- **catalogitemids** (optional): Catalog items which have to be owned, entries separated by semicolon. NOTE: key can’t be removed and value can’t be empty in dashboard. If ownership check is not required use “none” or “0” as value.

## Client Side

Clients have to send:

- **token**: ID Token (see ”Retrieving an ID Token For User” below)
- **ownershipToken** (optional): if catalogItemIds are configured in dashboard client has to send ownershipToken, NS verifies that items are owned

### Retrieving an ID Token

Clients have to use the Epic Account Service API to fetch an ID Token, which is described in [Epic's "Auth Interface - Retrieving an ID Token For User" doc.](https://dev.epicgames.com/docs/services/en-US/EpicAccountServices/AuthInterface/index.html#retrievinganidtokenforuser)

Excerpt from Epic's documentation: “Game clients can obtain ID Tokens for local users by calling the `EOS\_Auth\_CopyIdToken` SDK API after the user has been logged in, passing in a `EOS\_Auth\_CopyIdTokenOptions` structure containing the `EOS\_EpicAccountId` of the user.

The outputted `EOS\_Auth\_IdToken` structure contains the `EOS\_EpicAccountId` of the user, and a JWT representing the ID Token data. Note, you must call `EOS\_Auth\_IdToken\_Release` to release the ID Token structure when you are done with it.

Once retrieved, the game client can then provide the ID Token to another party. An ID Token is always readily available for a logged in local user.”

#### Sample Code

This can be used to retrieve an ID Token:

C#

```csharp
// Call this method after login.
private bool GetLocalIdToken(out IdToken? a_IdToken)
{
    var options = new CopyIdTokenOptions()
    {
        AccountId = LocalUserId
    };

    // NOTE: Make sure to use the EOSAuthInterface to get the IdToken instead of the EOSConnectInterface.
    var result = EOSManager.Instance.GetEOSAuthInterface().CopyIdToken(ref options, out a_IdToken);

    if (result != Result.Success)
    {
        Debug.LogError("Failed to copy the IdToken.");
        return false;
    }
    return true;
}

```

And the `AuthValues` can be setup like this:

C#

```csharp
var authValues = new AuthenticationValues();
authValues.AuthType = CustomAuthenticationType.Epic;
var idToken = /* token retrieved by GetLocalIdToken */;
authValues.AddAuthParameter("token", idToken.Value.JsonWebToken.ToString());
Client.AuthValues = authValues;

```

### Requesting an Ownership Verification Token

Epic's documentation describes [how to get the optional Ownership Verification Token](https://dev.epicgames.com/docs/web-api-ref/ecom-web-apis?sessionInvalidated=true#requesting-an-ownership-verification-token).

Back to top

- [Application Setup](#application-setup)

  - [Mandatory Settings](#mandatory-settings)

- [Client Side](#client-side)
  - [Retrieving an ID Token](#retrieving-an-id-token)
  - [Requesting an Ownership Verification Token](#requesting-an-ownership-verification-token)