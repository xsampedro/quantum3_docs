# overview-auth

_Source: https://doc.photonengine.com/realtime/current/connection-and-authentication/authentication/overview-auth_

# Overview

When you start a new project with any Photon SDK, you might notice that users (a.k.a. players) are anonymous. There is no need to login as user and you may notice that each new session gets a random `userId` assigned.

This is great to get projects started but there are serious disadvantages:

- no security against impersonation, anyone can claim any `userId`
- any inventory or stats are provided client side and easy to manipulate
- users can not identify others or find friends
- everyone has unconditional access to the online component of your app
- malicious users can not be banned

Due to this, we strongly recommend to add authentication to all apps before a general release.

While Photon itself does not provide user accounts, third party services can easily be integrated. Once setup, Photon uses a server to server REST api to authenticate users. These services can grant or deny access to Photon.

## Anonymous Users

Even if your Photon App does not require server-side authentication, clients always have to send an authentication operation. The default authentication request uses `CustomAuthenticationType.None` and happens behind the scenes when a client connects.

Unless your client code sets any credentials (in .Net via `AuthenticationValues`), the server assigns a new GUID as `userId` which lasts until the session ends.

Clients are by default allowed to identify themselves by sending a `userId` but this is not checked in any way.

Storing and re-using a `userId` client side is a very simple way to "identify" users but also very vulnerable to identity theft. This should be replaced with proper server side authentication.

### Rejecting Anonymous Users

Even after setting up Authentication Provider(s) for your App, clients can still try to use `AuthenticationValues.CustomAuthenticationType = CustomAuthenticationType.None`. To reject these clients, make sure to uncheck "Allow anonymous clients to connect" per App in the Dashboard.

## Authentication Setup

Authentication requires some coordination between the server side and your clients.

Per application (and corresponding AppId), the setup is done in the Photon Dashboard first. Several popular services for user accounts are predefined and any other service can be added as "Custom Server". This setup defines the credentials for the server-to-server authentication calls.

Clients will always send an authentication operation when they connect. To request authentication with a specific Authentication Provider, the clients have to set their `AuthenticationValues` before they connect. The required values depend on the service that provides the account. Some services just require a username and password while others provide their own SDKs and APIs to login and fetch an authentication token to verify the user on Photon.

The client's `AuthenticationValues.CustomAuthenticationType` value defines which service to use, so you could mix and match.

### Predefined Providers

Photon implements serveral popular user account services directly to simplify their use. For each service, the Photon Dashboard will ask for a set of values to make the association between the systems. Client side, each service has a corresponding value in the `CustomAuthenticationType` enum.

These services are described in more detail in separate doc pages. Clients will have to meet the expected values per service.

### Custom Authentication

Custom Authentication can be used to integrate any service not covered by the predefined providers. If you have a user base and want to bring them into a Photon title Custom Authentication can be used.

A web service must be setup to answer Photon's requests for authorization. Via the response, the service grants access to Photon, sets the `userId` or additional values. Client side, `AuthType = CustomAuthenticationType.Custom` is used.

The [Custom Authentication doc page](/realtime/current/connection-and-authentication/authentication/custom-authentication) provides all the details to implement and configure the required REST api to authenticate users with any backend.

### PlayFab Integration

While PlayFab is a popular choice as "user backend service", it's authentication is implemented as Custom Authentication.

See the [PlayFab Integration doc](/realtime/current/reference/playfab) for more details.

### Authentication on Consoles

Photon provides Authentication Providers for most popular console platforms. Each is described on the related Console doc page.

Check the [Consoles doc category](/realtime/current/consoles/overview) for more details.

## Adding Voice and Chat

The Voice and Chat SDKs are separate solutions and have their own Authentication configurations.

When Voice and Chat are paired with other SDKs, we would recommend to configure each AppId to offer the same Authentication options as e.g. Fusion or Quantum do.

The Voice SDK for Unity even comes with an option to re-use the `AuthenticationValues` of Fusion.

Back to top

- [Anonymous Users](#anonymous-users)

  - [Rejecting Anonymous Users](#rejecting-anonymous-users)

- [Authentication Setup](#authentication-setup)

  - [Predefined Providers](#predefined-providers)
  - [Custom Authentication](#custom-authentication)
  - [PlayFab Integration](#playfab-integration)
  - [Authentication on Consoles](#authentication-on-consoles)

- [Adding Voice and Chat](#adding-voice-and-chat)