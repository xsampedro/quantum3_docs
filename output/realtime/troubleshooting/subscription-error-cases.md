# subscription-error-cases

_Source: https://doc.photonengine.com/realtime/current/troubleshooting/subscription-error-cases_

# Subscription Error Cases

When you use the Photon Cloud instead of hosting yourself, there are a few additional error cases to keep in mind.

These are explained below.

Additionally, check the documentation included in the client platform SDKs for implementation details.

In all cases, the [Photon Cloud Dashboard](https://dashboard.photonengine.com) is useful to check the AppId and subscription.

## Authentication Errors

The following errors happen during authentication with the cloud.

The client is already connected but identifies the respective game by your AppId.

Authentication is done automatically when you call `Connect`.

### Unknown or Archived AppId

This happens if the provided AppId is not known to the server at all.

Check for typos and make sure you use the correct application type.

### CCU Limit Reached

Unless you have a plan with "CCU Burst", clients might fail the authentication step during connect.

Affected clients are unable to call operations.

Please note that players who end a game and return to the master server will disconnect and re-connect, which means that they just played and are rejected in the next minute / re-connect.

This is a temporary measure.

Once the CCU is below the limit, players will be able to connect an play again.

OpAuthenticate is part of connection workflow but only on the Photon Cloud, this error can happen.

Self-hosted Photon servers with a CCU limited license won't let a client connect at all.

### Plugin Mismatch

This error happens when the AppId used to connect to Photon Cloud is not of the same type as the SDK in use.

For example, trying to start a Photon Fusion SDK Client using a Realtime SDK AppId type or a Voice SDK AppId.

The AppId type can be verified directly via the [Photon Cloud Dashboard](https://dashboard.photonengine.com).

If you are experiencing this type of error, review the AppId type and/or create a new AppId to check.

Back to top

- [Authentication Errors](#authentication-errors)
  - [Unknown or Archived AppId](#unknown-or-archived-appid)
  - [CCU Limit Reached](#ccu-limit-reached)
  - [Plugin Mismatch](#plugin-mismatch)