# steam-auth

_Source: https://doc.photonengine.com/realtime/current/connection-and-authentication/authentication/steam-auth_

# Steam Authentication

## Application Setup

Adding Steam as an authentication provider is easy and it could be done in few seconds from your [Photon Applications' Dashboard](https://dashboard.photonengine.com/).

Go to the "Manage" page of an application and scroll down to the "Authentication" section.

If you add a new authentication provider for Steam or edit an existing one, here the mandatory settings:

- **apiKeySecret**: Steam Publisher Web API key. Do not confuse it with Steam User Key. Read more about how to get one [here](https://partner.steamgames.com/doc/webapi_overview/auth#create-publisher-key).
- **appid**: ID of the Steam game. You can get one after going through [Steam Direct](https://partner.steamgames.com/steamdirect) process (formerly known as Steam Greenlight).
- **verifyOwnership**: Can be `true` or `false`: Whether or not to enable Ownership Verification during authentication.


This allows you to verify if the user really owns (purchased the game and has it in his library) the game.


This step, if enabled, will be performed just after validating the user's session ticket.


Enabling this may add extra delay in authentication, so enable it only if you really need it.
- **verifyVacBan**: Can be `true` or `false`: Whether or not to check if the user has been banned using Valve's Anti-Cheat (VAC) during authentication.


Read more [here](https://partner.steamgames.com/doc/features/anticheat).


Enabling this may add extra delay in authentication, so enable it only if you really need it.
- **verifyPubBan**: Can be `true` or `false`: Whether or not to check if the user has been banned using a Publisher Ban during authentication.


Read more [here](https://partner.steamgames.com/doc/features/anticheat).


Enabling this may add extra delay in authentication, so enable it only if you really need it.
- **version**: Can be 1 or 2 (default value 1). If version is set to 2 the identity parameter is used when verifying the session ticket. Leave at 1 if you don't use identity.
- **identity**: Used if version is set to 2. Steamworks SDK 1.57 added GetAuthTicketForWebAPI which requires an identity parameter. Can be any string identifier (default "photon").

## Client Code (Unity)

The client must use Valve's Steamworks API to get a session ticket.

This ticket is proof that the client is a valid Steam user.

### Steamworks.NET

Steamworks.NET is a popular free and open source Steamworks API wrapper.

Follow [the instructions listed on this page](https://Steamworks.github.io/installation/#unity-instructions) to import a Unity version of Steamworks.NET.

#### Get Ticket

Use the following code to get a session ticket using the Steamworks API and convert it to a hex encoded UTF-8 string:

C#

```csharp
// hAuthTicket should be saved so you can use it to cancel the ticket as soon as you are done with it
public string GetSteamAuthTicket(out HAuthTicket hAuthTicket)
{
    byte[] ticketByteArray = new byte[1024];
    uint ticketSize;
    hAuthTicket = SteamUser.GetAuthSessionTicket(ticketByteArray, ticketByteArray.Length, out ticketSize);
    System.Array.Resize(ref ticketByteArray, (int)ticketSize);
    StringBuilder sb = new StringBuilder();
    for(int i=0; i < ticketSize; i++)
    {
        sb.AppendFormat("{0:x2}", ticketByteArray[i]);
    }
    return sb.ToString();
}

```

See version and identity parameters description above, if version is set to 2 GetAuthTicketForWebAPI has to be used instead of GetAuthSessionTicket.

#### Send Ticket

The client must send the user's session ticket (after converting it to a hex encoded UTF-8 string) as a value of a query string key "ticket".

C#

```csharp
loadBalancingClient.AuthValues = new AuthenticationValues();
loadBalancingClient.AuthValues.UserId = SteamUser.GetSteamID().ToString();
loadBalancingClient.AuthValues.AuthType = CustomAuthenticationType.Steam;
loadBalancingClient.AuthValues.AddAuthParameter("ticket", SteamAuthSessionTicket);
// connect

```

#### Cancel Ticket

It is recommended to cancel or revoke the ticket once authentication is done.

C#

```csharp
SteamUser.CancelAuthTicket(hAuthTicket);

```

### Facepunch.Steamworks

Facepunch.Steamworks is yet another alternative free and open source implementation of Steamworks API.

Follow [the instructions listed on this page](https://wiki.facepunch.com/steamworks/Installing_For_Unity) to import Facepunch.Steamworks.

#### Get Ticket

Use the following code to get a session ticket and convert it to a hex encoded UTF-8 string:

C#

```csharp
// authTicket should be saved so you can use it to cancel the ticket as soon as you are done with it
public string GetSteamAuthTicket(out AuthTicket authTicket)
{
    authTicket = SteamUser.GetAuthSessionTicket();
    StringBuilder ticketString = new StringBuilder();
    for (int i = 0; i < authTicket.Data.Length; i++)
    {
        ticketString.AppendFormat("{0:x2}", authTicket.Data[i]);
    }
    return ticketString.ToString();
}

```

#### Send Ticket

The client must send the user's session ticket (after converting it to a hex encoded UTF-8 string) as a value of a query string key "ticket".

C#

```csharp
loadBalancingClient.AuthValues = new AuthenticationValues();
loadBalancingClient.AuthValues.UserId = SteamClient.SteamId.ToString();
loadBalancingClient.AuthValues.AuthType = CustomAuthenticationType.Steam;
loadBalancingClient.AuthValues.AddAuthParameter("ticket", SteamAuthSessionTicket);
// connect

```

#### Cancel Ticket

It is recommended to cancel or revoke the ticket once authentication is done.

C#

```csharp
ticket.Cancel();

```

## Change History

**June 20, 2023:**

- Added: description for Steam identity usage

Back to top

- [Application Setup](#application-setup)
- [Client Code (Unity)](#client-code-unity)

  - [Steamworks.NET](#steamworks.net)
  - [Facepunch.Steamworks](#facepunch.steamworks)

- [Change History](#change-history)