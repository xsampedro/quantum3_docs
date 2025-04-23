# playfab

_Source: https://doc.photonengine.com/realtime/current/reference/playfab_

# PlayFab Integration

## Introduction

In this document we help you integrate PlayFab with Photon.

With this approach, both systems will get used in parallel to their full potential.

Billing is separate for each service.

Read [PlayFab's Photon Quickstart here](https://docs.microsoft.com/en-us/gaming/playfab/sdks/photon/quickstart).

Below are the setup instructions for Photon.

## Custom Authentication

### Dashboard Configuration

Here are the steps to setup custom authentication with PlayFab:

1. Go to Photon dashbaord.
2. Choose an application or create a new one.
3. Click "Manage".
4. Under "Authentication" section, click "Custom Server".
5. \[Required\] Set Authentication URL to `https://{PlayFabTitleId}.playfabapi.com/photon/authenticate`.


Make sure to replace `{PlayFabTitleId}` placeholder with your actual PlayFab TitleId.


Do not keep the open (`{`) and closing (`}`) braces or curly brackets characters in the URL.

**Example**: if your PlayFab TitleId is `AB12`: `https://AB12.playfabapi.com/photon/authenticate`.
6. Save by hitting "Create".
7. \[Recommended\] Untick "Allow anonymous clients to connect, independently of configured providers".

### Client Code

Client is expected to send a pair of key/values as credentials:

- PlayFab UserId of the logged in user.
- Photon token (retrieved using [GetPhotonAuthenticationToken](https://learn.microsoft.com/en-us/rest/api/playfab/client/authentication/get-photon-authentication-token?view=playfab-rest) Client API method).

C#

```csharp
lbClient.AuthValues = new AuthenticationValues();
lbClient.AuthValues.AuthType = CustomAuthenticationType.Custom;
lbClient.AuthValues.AddAuthParameter("username", PlayFabUserId);
lbClient.AuthValues.AddAuthParameter("token", PlayFabPhotonToken);
// do not set AuthValues.Token or authentication will fail
// connect

```

C++

```cpp
ExitGames::Common::JString params = "username="+PlayFabUserId+"&token="+PlayFabPhotonToken;
ExitGames::LoadBalancing::AuthenticationValues playFabAuthenticationValues;
playFabAuthenticationValues.setType(ExitGames::LoadBalancing::CustomAuthenticationType::CUSTOM).setParameters(params);
// pass playFabAuthenticationValues as parameter on connect

```

JavaScript

```javascript
var queryString = "username="+playFabUserId+"&token="+playFabPhotonToken;
loadBalancingClient.setCustomAuthentication(queryString, Photon.LoadBalancing.Constants.CustomAuthenticationType.Custom);
// connect

```

## Realtime Webhooks and WebRPC Configuration

If you do not need Realtime Webhooks nor WebRPCs for your application, you can skip this part.

It is recommended to setup custom authentication first.

Realtime Webhooks and WebRPCs may not work otherwise.

PlayFab has deprecated CloudScript (Classic) as announced in [this blogpost](https://blog.playfab.com/index.php/blog/announcing-cloudscript-using-azure-functions-is-now-ga).
Only paying customers who were actively using CloudScript (Classic) could still make use of the Webhooks integration as described in this document.
Read more [here](https://community.playfab.com/answers/57746/view.html).
If you are getting the following **Error response ResultCode = '1' Message = 'Not authorized'**, it means that you can't create or join rooms anymore with Webhooks configured to use CloudScript (Classic) as your PlayFab title can't access that anymore.
Get in touch with PlayFab or start looking into other options for implementing Photon Webhooks endpoints.
For instance, PlayFab replaced CloudScript Classic with [CloudScript Functions which is using Azure Functions](https://docs.microsoft.com/en-us/gaming/playfab/features/automation/cloudscript-af).

Here is how to setup Realtime Webhooks and WebRPCs to work with your PlayFab title.

1. Go to Photon dashbaord.
2. Choose an application or create a new one.
3. Click "Manage".
4. Under "Webhooks" section, click "Create a new Webhook".
5. From the dropdown list "Select Type" choose Webhooks v1.2.
6. \[Required\] Set "BaseUrl" to `https://{PlayFabTitleId}.playfablogic.com/webhook/1/prod/{PhotonSecretKey}`.

   - Make sure to replace `{PlayFabTitleId}` placeholder with your actual PlayFab TitleId.


     Do not keep the open (`{`) and closing (`}`) braces or curly brackets characters in the URL.

     **Example**: if your PlayFab TitleId is `AB12` and your Photon Secret Key is `3z4ujpikyp4hbufno3jfm6fzcxsw5zxjaberfe4zkf4hzuen3`:

     `https://AB12.playfablogic.com/webhook/1/prod/3z4ujpikyp4hbufno3jfm6fzcxsw5zxjaberfe4zkf4hzuen3`.
   - The "Photon Secret Key" is a string that you find on PlayFab's Photon add-on page once you add a Photon Realtime application.
   - You can replace `prod` with `test` if you want to target latest uploaded / pushed CloudScript revision as opposed to the latest deployed / active one.
7. \[Optional\] Configure the Webhooks Paths you need.


Remove those you do not need.


Read more [here](/realtime/current/gameplay/web-extensions/webhooks#paths).
8. \[Optional\] Configure "IsPersistent", "HasErrorInfo" and "AsyncJoin".


Remove their keys if you want to keep default values.


Read more [here](/realtime/current/gameplay/web-extensions/webhooks#options).

### Notes

- "CustomHttpHeaders" setting is not supported as those are not exposed in the CloudScript handlers. Any value you set will not be useful.

- PlayFab suggests the following names for Realtime Webhooks CloudScript handlers functions:


  - "PathCreate": "RoomCreated"
  - "PathClose": "RoomClosed" (do not change this handler name)
  - "PathJoin": "RoomJoined"
  - "PathLeave": "RoomLeft"
  - "PathEvent": "RoomEventRaised"
  - "PathGameProperties": "RoomPropertyUpdated"

Other than "PathClose", you are free to choose whatever handler name you want for the other paths.

PlayFab expects a valid PlayFabId as UserId argument in all Webhooks or WebRPCs handlers.

The only exception to this is "RoomClosed".

- All Photon Realtime Webhooks except "PathClose" and all WebRPCs should set PlayFab's CloudScript global variable "currentPlayerId" to the value of the "UserId" argument.

- Configured Realtime Webhooks or WebRPCs called by client will not work unless there is an explicit CloudScript handler with the same name.

**Examples:**


If you configured "PathJoin" to "GameJoined", you need to have this function in the target CloudScript revision:

JavaScript
```javascript
handlers.GameJoined = function(args)
{
      // your custom code goes here
      return { ResultCode : 0, Message: "Success" };
};

```


If client calls "foo" WebRPC method, you need to have this function in the target CloudScript revision:

JavaScript
```javascript
handlers.foo = function(args)
{
      // your custom code goes here
      return { ResultCode : 0, Message: "Success" };
};

```


Back to top

- [Introduction](#introduction)
- [Custom Authentication](#custom-authentication)

  - [Dashboard Configuration](#dashboard-configuration)
  - [Client Code](#client-code)

- [Realtime Webhooks and WebRPC Configuration](#realtime-webhooks-and-webrpc-configuration)
  - [Notes](#notes)