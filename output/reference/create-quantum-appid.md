# create-quantum-appid

_Source: https://doc.photonengine.com/quantum/current/reference/create-quantum-appid_

# Create A Photon Quantum AppId

Go to the [Photon Engine Dashboard](https://dashboard.photonengine.com) and log in or create an account.

On the dashboard select `Applications > Public Cloud` and press `CREATE A NEW APP`.

![](https://doc.photonengine.com/docs/img/quantum/v3/reference/dashboard.png)

- Select **Quantum** as Photon SDK
- Select **Quantum 3** as SDK Version

![](https://doc.photonengine.com/docs/img/quantum/v3/reference/dashboard-create-new-appid.png)

Copy the AppId from the freshly created Quantum app.

![](https://doc.photonengine.com/docs/img/quantum/v3/reference/dashboard-appid.png)

In the Unity Editor paste the AppId into the Quantum Hub text field or into the `PhotonServerSettings.asset` directly.

A) Paste into the Quantum Hub mask

![](https://doc.photonengine.com/docs/img/quantum/v3/reference/hub-appid.png)

B) Select the `PhotonServerSettings` asset and paste as `AppSettings > App Id Quantum`

![](https://doc.photonengine.com/docs/img/quantum/v3/reference/quantum-menu-photonserversettings.png)![](https://doc.photonengine.com/docs/img/quantum/v3/reference/quantum-photonserversettings.png)Back to top