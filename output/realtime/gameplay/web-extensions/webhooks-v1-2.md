# webhooks-v1-2

_Source: https://doc.photonengine.com/realtime/current/gameplay/web-extensions/webhooks-v1-2_

# New in Webhooks v1.2

## "Username" is now "NickName"

This is a breaking change and must be applied when you upgrade to Webhooks 1.2!

Any previous occurrence of the `Username` property in the data of webhook requests will be renamed to `NickName`.

This change affects all webhooks.

- `Username` is a top level property in all webhooks data except PathClose.
- `Username` is a property of `ActorList` array elements.

`ActorList` is part of the serialized room state.


Webhooks that may include the state are PathClose (Type="Save"), PathProperties and PathEvent.

## Error Codes

For easier debugging find new error codes to replace the old `InternalServerError (-1)`.

These new codes cover all join related exceptions.

- `JoinFailedPeerAlreadyJoined (32750)`:


Indicates the current peer already called join and is joined to the room.
- `JoinFailedFoundInactiveJoiner (32749)`:


Indicates the list of InactiveActors already contains an actor with the requested ActorNr or UserId.
- `JoinFailedWithRejoinerNotFound (32748)`:


Indicates the list of actors (active and inactive) did not contain an actor with the requested ActorNr or UserId.
- `JoinFailedFoundActiveJoiner (32746)`:


Indicates the list of ActiveActors already contains an actor with the requested ActorNr or UserId.

## WebFlags

WebFlags are being phased out. We will keep them working for existing games but the .Net Realtime SDK v5 no longer has WebFlags. If this feature is a must for you, stay with the .Net Realtime SDK v4.
[Mail us for questions](https://www.photonengine.com/contact).

Webflags are optional flags to be used in Photon client SDKs with `OpRaiseEvent` and `OpSetProperties` methods.

They affect the `PathEvent` and `PathProperties` webhooks.

### HttpForward

This webflag is how Webhooks 1.2 will internally replace the previous boolean parameter in `OpRaiseEvent` and `OpSetProperties` methods with the same name.

The purpose of this flag is to decide whether or not to also forward those operations to web service.

All other webflags described below will depend on `HttpForward (0x01)`.

They cannot be considered enabled only if `HttpForward` is set.

### SendAuthCookie

Webhooks offers an option to securely transmit the encrypted object `AuthCookie` to the web service when available.

This can be done by setting the appropriate webflag (`SendAuthCookie (0x02)`).

The `AuthCookie` could be retrieved after a successful authentication against a custom authentication provider.

For more information please visit the [custom authentication documentation page](/realtime/current/connection-and-authentication/authentication/custom-authentication).

### SendSync

In Webhooks 1.2, it is now possible to choose if HTTP queries should be processed synchronously or asynchronously.

Webhooks are processed asynchronously by default.

Now with `SendSync (0x04)`, it is possible to make webhooks plugin block and wait for response before continuing normal process of operations.

### SendState

This is a breaking change and must be applied when you upgrade to Webhooks 1.2!

In Webhooks 1.0 the serialized room state was automatically sent with all _PathEvent_ and _PathProperties_ webhooks as `State` argument.

This increases traffic and adds delay.

This could be problematic for games with lots of forwarded events, frequent property updates and a big game state.

Besides not all developers need access to room state within those webhooks and it was not possible to change the default behavior.

That is why in Webhooks 1.2, `SendState (0x08)` webflag was introduced to give developers the ability to choose whether or not to send room state when forwarding events or property updates to web service.

`SendState` webflag works only if "IsPersistent" webhooks setting is enabled (set to `true`).

## Add your own HTTP headers

With Webhooks 1.2, it is now possible to configure custom headers to be sent with every webhook's HTTP request.

This can be done by adding the headers in JSON format as a value of `CustomHttpHeaders` key when configuring webhooks for your application.

The new feature will open a lot of interesting possibilities like providing a secret API key for your web service as a HTTP header.

Back to top

- ["Username" is now "NickName"](#username-is-now-nickname)
- [Error Codes](#error-codes)
- [WebFlags](#webflags)

  - [HttpForward](#httpforward)
  - [SendAuthCookie](#sendauthcookie)
  - [SendSync](#sendsync)
  - [SendState](#sendstate)

- [Add your own HTTP headers](#add-your-own-http-headers)