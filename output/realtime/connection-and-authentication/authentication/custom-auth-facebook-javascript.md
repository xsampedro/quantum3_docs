# custom-auth-facebook-javascript

_Source: https://doc.photonengine.com/realtime/current/connection-and-authentication/authentication/custom-auth-facebook-javascript_

# Facebook Authentication in JS

This tutorial gives a walk through on how to authenticate a user in a Photon Cloud application with his Facebook credentials using Photon Custom Authentication.

## Overview

Photon Cloud Custom Authentication involves calling an external service (Authentication Provider) to check user authentication.

For an outline of the authentication flow, please check our [general documentation](/realtime/current/connection-and-authentication/authentication/custom-authentication).

Usually several parameters are required for the authentication check, depending on the service you want to authenticate against.

These are static parameters set in your Photon Cloud Application settings and dynamic parameters passed in the Photon Cloud JavaScript Client API call.

This tutorial uses the [Facebook Login plugin](https://developers.facebook.com/docs/reference/plugins/login/)

on the client side to get the user's Facebook token and to check the token validity from Photon Servers using our Facebook Authentication Provider.

For Facebook authentication, the Authentication Provider requires _appid_, _secret_ and _token_ parameters and checks if the given token is valid for the application with the given appid at Facebook.

## Implementation

Enabling Custom Authentication for an existing application involves a few simple steps:

### 1\. Facebook App Setup

Reference the Facebook documentation for details on Facebook App setup.

For this tutorial you need your **Facebook App ID** and **App Secret**.

### 2\. Photon Cloud Application

Go to the _Details_ page of your application via your [Photon Cloud Dashboard](https://dashboard.photonengine.com/).

Expand the Custom Authentication section.

The _Authentication URL_ for Facebook authentication is set by Exit Games.

Set the values for these two parameters:

- appid = _your Facebook App ID_
- secret = _your Facebook App Secret_

Save changes.

### 3\. HTML5 / JavaScript Client

The client needs to pass

Add a Facebook Login button div anywhere inside your html document body.

HTML

```markup
<div class="fb-login-button" data-show-faces="true" data-width="400" data-max-rows="1" data-registration-url=""></div>

```

Add code for Facebook JavaScript SDK loading, Facebook Login button creation and Facebook authentication status handling.

Ideally right after the opening <body> tag.

HTML

```markup
<div id="fb-root"></div>
<script>
    window.fbAsyncInit = function () {
        FB.Event.subscribe('auth.authResponseChange', function (response) {
            if (response.status === "connected") {
                loadBalancingClient.setCustomAuthentication("token=" + response.authResponse.accessToken, 2);
                loadBalancingClient.connect(true);
            }
        });
    };
    (function (d, s, id) {
        var js, fjs = d.getElementsByTagName(s)[0];
        if (d.getElementById(id)) return;
        js = d.createElement(s); js.id = id;
        js.src = "https://connect.facebook.net/en_US/all.js#xfbml=1&appId=" + "<YOUR_APP_ID>";
        fjs.parentNode.insertBefore(js, fjs);
    }(document, 'script', 'facebook-jssdk'));
</script>

```

This code loads the Facebook JavaScript SDK asynchronously and creates a Facebook Login button.

The function `window.fbAsyncInit` is called after loading completed.

It subscribes a handler on Facebook `auth.authResponseChange` event, which gets fired with `response.status = "connected"` after the user successfully logged in to Facebook.

The user's token gets passed to the handler as parameter.

Having this token, the handler sets Custom Authentication mode for the client and adds the additional parameter _token_ to the authentication URL.

Then the handler starts a connection to the Master server.

### 4\. Deploy and Run

Place your HTML5 / JavaScript client code on one of the App Domains specified in your Facebook App settings.

This is required for proper functionality of the Facebook Login button.

User accessing your application page will log in to Photon Cloud after pressing the Facebook Login button and entering their username/password.

If already logged in at Facebook, they will be signed in to Photon Cloud immediately.

Back to top

- [Overview](#overview)
- [Implementation](#implementation)
  - [1\. Facebook App Setup](#facebook-app-setup)
  - [2\. Photon Cloud Application](#photon-cloud-application)
  - [3\. HTML5 / JavaScript Client](#html5-javascript-client)
  - [4\. Deploy and Run](#deploy-and-run)