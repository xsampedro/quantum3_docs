# cheat-protection

_Source: https://doc.photonengine.com/quantum/current/manual/cheat-protection_

# Cheat Protection

## Introduction

Security and cheat-protection are important aspects of online multiplayer gaming. Quantum's determinism offers unique features to address these issues. This page delves into the details of the built-in protection and offers best practices for creating production-ready online games with Photon Quantum.

It is vital that developers are aware of any security issues, the steps to mitigate them and when they should be taken. Although it is possible to run the 100% of the entire simulation on the server, it is rare that this is practical or cost effective.

- Servers are expensive, especially if the game is not yet generating revenue.
- In most cases, cheaters make up a very small proportion of the user base.
- Making a game 100% cheat-proof is a utopian idea, and even staying one step ahead of them is a daunting task.
- There are genres of games that need to be as safe as possible.

The single most important advice is: write the game **now** and add these more complex safety checks **incrementally**. It is perfectly viable to go live without a custom server and be successful.

In the documentation the following terms are used:

- A **game backend** refers to online services created and hosted by the customer.
- A **custom server plugin** refers to a customer created Quantum plugin hosted by Photon.
- The **Quantum Public Cloud** refers to non-customized Quantum servers hosted by Photon.
- **Webhooks** are standardised HTTP requests called from Quantum 3 games running in the Public Cloud.

## Cheat Protection By Determinism

The great advantage of a deterministic game, even without the server running the simulation, is that it is cheat-protection robust; if a player modifies its client, for example by changing its character's speed, it will **not** affect other players. They may notice that the cheating player is behaving strangely (e.g. constantly bumping into a walls) but otherwise their experience is unaffected.

The reason for this is simple: each client runs the entire simulation deterministically locally and only shares input with other clients.

## Trustable Match Results

A match result can be used to update the player progression in a game backend. In the most secure scenarios this is done from the server where the game logic was running.

However, there are a number of iterations that can be used to validate the results in a cost-effective manner before moving to a dedicated game server.

| Prototyping | Clients push their individual results to the developer's game backend. This is good for prototyping and games can be launched with this setup. Having a data structure that can be filled with the results to be sent to a backend is the first thing to have. |
| Game Result Webhook | Clients send results to the Quantum server from within the simulation all at the same time. The GameResult webhook is called when the online game room is finally disbanded and contains a list of all the results sent by clients. This solution is similar to the Prototyping suggestion but already uses the GameResult API from the Custom Quantum Server suggestion. |
| Resimulation Replays | If the game results are inconclusive or are generally untrustworthy the session can be resimulated using the captured input stream and a non-Unity session runner application. See the Replay manual for information on capturing and replaying online game sessions. |
| Custom Quantum Server | Run the simulation on our Photon Enterprise servers, which will automatically call the GameResult webhook from the server without having to wait for clients to send their data. |

## Protect Client-Controlled Game Configs

Many start parameters of the actual online Quantum simulation and game are client controlled by default.

Each client uploads ```
SessionConfig
```

 (Quantum settings) and ```
RuntimeConfig
```

(game settings) files before the session starts. Clients also upload their ```
RuntimePlayer
```

, which usually describes their game loadout or progression when a player is added to the game.

To secure the data sent by clients, in games running on the Public Quantum Cloud, [Quantum Webhooks](/quantum/current/manual/webhooks) can be used.

| Public Quantum Cloud | The server chooses the first ```<br>SessionConfig<br>```<br> and ```<br>RuntimeConfig<br>```<br> it receives from clients, which is more or less random. ```<br>RuntimePlayer<br>```<br> is completely unprotected. |
| Webhooks | ```<br>SessionConfig<br>```<br>, ```<br>RuntimeConfig<br>```<br> and ```<br>RuntimePlayer<br>```<br> sent by clients can be checked using HTTP webhooks configured on the Photon dashboard to call the developer's game backend. |
| Custom Quantum Server | All configs can be intercepted and replaced after being retrieved from the developer's game backend. |

## Custom Authentication

We do not provide an authentication service or a player database ourselves but we strongly recommend to add a proprietary or third party **authentication service**.

[Photon Realtime Custom Authentication](/realtime/current/connection-and-authentication/authentication/custom-authentication)

## Determinism As A Drawback

While determinism has many advantages, there are some notable drawbacks inherent in this type of technology.

### Perfect Information Problem

With Quantum, each client has access to all the information it needs to simulate the game locally (except other players' input). This means that **client-controlled secrets** used in a card game and Fog Of War-like features **are easily hackable**.

There are also fringe cases that allow clients to "guess" a **next random number** or the ability to **create bots**.

### Detecting Cheaters Using Checksums

It is **not** recommended to use the Quantum checksum detection for live games as a way to detect cheaters.

- Checksum calculation is expensive and could lead to hiccups; and,
- The built-in mechanism stops the simulation for **every** client in the game session immediately.

Back to top

- [Introduction](#introduction)
- [Cheat Protection By Determinism](#cheat-protection-by-determinism)
- [Trustable Match Results](#trustable-match-results)
- [Protect Client-Controlled Game Configs](#protect-client-controlled-game-configs)
- [Custom Authentication](#custom-authentication)
- [Determinism As A Drawback](#determinism-as-a-drawback)
  - [Perfect Information Problem](#perfect-information-problem)
  - [Detecting Cheaters Using Checksums](#detecting-cheaters-using-checksums)