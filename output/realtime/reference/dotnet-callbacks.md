# dotnet-callbacks

_Source: https://doc.photonengine.com/realtime/current/reference/dotnet-callbacks_

# C# Callbacks

### C# SDKs provides callback interfaces that you can implement in your classes:

**Photon Realtime C# SDK callbacks interfaces:**

- `IConnectionCallbacks`: connection related callbacks.
- `IInRoomCallbacks`: callbacks that happen inside the room.
- `ILobbyCallbacks`: lobby related callbacks.
- `IMatchmakingCallbacks`: matchmaking related callbacks.
- `IOnEventCallback`: a single callback for any received event. This is 'equivalent' to the C# event `LoadBalancingClient.EventReceived`.
- `IWebRpcCallback`: a single callback for receiving WebRPC operation response.
- `IOnErrorInfoCallback`: a single callback for receiving ErrorInfo event.

You can implement one or more interfaces per class.

You can also implement the same interface by more than one class.

The price to pay is that classes that implement these interfaces may become too long or have unused methods.

All classes implementing callback interfaces must be registered and unregistered.

Call `LoadBalancingClient.AddCallbackTarget(this)` and `LoadBalancingClient.RemoveCallbackTarget(this)`.

For example in Unity, you could use the `MonoBehaviour`'s `OnEnable()` and `OnDisable()` or `Start()`/`Awake()` and `OnDestroy()`.

Implementing these interfaces is optional but recommended as we think it can make your code more readable and maintainable.

It also makes the Photon flow and states easier to manage by providing the exact timing to execute some logic.

Other alternatives may require the usage of state flags fields, polling to check client networking state or subscribing to all networking client's status changes or received events or operation responses.

This requires deep knowledge about some internals or low-level Photon details, which you can avoid and focus on your game.

If an unhandled (uncaught) exception occurs in one of the implemented interfaces' callbacks' methods, all other implemented ones for the same interface, same signature and not already called, won't be called.
This is due to the fact that we call the implemented interface callback methods of the same signature in a loop in the order of their registration (which could be random in Unity if you register in `MonoBehaviour` methods).

The reasons behind choosing interfaces over other ways of implementing a callbacks system:

- making sure callbacks' methods signatures are respected which is guaranteed by the compiler when implementing interfaces
- grouping callbacks logically related into a single class
- compared to other methods of providing callbacks, it has less garbage overhead and avoids memory leaks

If you happen to have a method that has the exact signature as one of the callbacks' interfaces' methods or you wish to hide the callbacks methods (unless a cast is made) you could choose [explicit interface implementation](https://docs.microsoft.com/en-us/dotnet/csharp/programming-guide/interfaces/explicit-interface-implementation).

Back to top

- [C# SDKs provides callback interfaces that you can implement in your classes:](#c-sdks-provides-callback-interfaces-that-you-can-implement-in-your-classes)