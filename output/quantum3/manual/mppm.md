# mppm

_Source: https://doc.photonengine.com/quantum/current/manual/mppm_

# Multiplayer Play Mode

Unity's new ```
Multiplayer Play Mode
```

 can be easily used with Quantum 3. The Quantum sample menu provides communication between the master and virtual players to start and connect all players with a single click in the master player. See the ```
QuantumMppm
```

and ```
QuantumMenuMppmJoinCommand
```

 class for more details.

Requires Unity Editor 6

## How to install and use MPPM with Quantum

1. Install the ```
   Multiplayer Play Mode
   ```

    package in the Unity package manager.

![](/docs/img/quantum/v3/manual/mppm/mppm-package-manager.png)

2. Open the ```
   Multiplayer Play Mode
   ```

    control window.

![](/docs/img/quantum/v3/manual/mppm/mppm-menu.png)

3. Enable at least on virtual player.

![](/docs/img/quantum/v3/manual/mppm/mppm-open-mppm-window.png)

4. Open a Quantum menu scene, for example from the Platform Shooter 2D sample
   - Enter Unity Editor play mode.
   - Start an online match ```
     QUICK PLAY
     ```

      or through ```
     PARTY MODE
     ```

      in the master player, it will communicate to other virtual players to connect to the same Photon room.

![](/docs/img/quantum/v3/manual/mppm/mppm-connecting.png)![](/docs/img/quantum/v3/manual/mppm/mppm-gameplay.png)

Disable the MPPM communication by toggling ```
EnableMppm
```

in the ```
QuantumMenuConnectionBehaviour
```

to start the connection on each player individually.

![](/docs/img/quantum/v3/manual/mppm/mppm-config.png)Back to top

- [How to install and use MPPM with Quantum](#how-to-install-and-use-mppm-with-quantum)