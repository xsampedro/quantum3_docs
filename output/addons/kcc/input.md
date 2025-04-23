# input

_Source: https://doc.photonengine.com/quantum/current/addons/kcc/input_

# Input

To get perfectly smooth movement and camera, it is mandatory to apply input signal correctly. Following sections focus on camera look rotation based on data received from mouse device.

## Getting Input

The mouse delta is obtained using:

- ```
UnityEngine.Input.GetAxis("Mouse X")
```

- ```
UnityEngine.Input.GetAxisRaw("Mouse X")
```

- ```
UnityEngine.InputSystem.Mouse.current.delta.ReadValue()
```


All methods above return position delta polled from mouse without any post-processing. Depending on mouse polling rate (common is 125Hz, gaming mouses usually 1000Hz), you get these values for ```
Mouse X
```

when moving with the mouse:

![Mouse input with 125Hz polling and 120 FPS](/docs/img/quantum/v3/addons/kcc/input-smoothing-125hz-120fps.jpg)Mouse input with 125Hz polling and 120 FPS.

![Mouse input with 125Hz polling and 360 FPS](/docs/img/quantum/v3/addons/kcc/input-smoothing-125hz-360fps.jpg)Mouse input with 125Hz polling and 360 FPS.

![Mouse input with 1000Hz polling and 120 FPS](/docs/img/quantum/v3/addons/kcc/input-smoothing-1000hz-120fps.jpg)Mouse input with 1000Hz polling and 120 FPS.

![Mouse input with 1000Hz polling and 360 FPS](/docs/img/quantum/v3/addons/kcc/input-smoothing-1000hz-360fps.jpg)Mouse input with 1000Hz polling and 360 FPS.

As you can see, each case produces different results. In almost all cases propagation directly to camera rotation results in looking not being smooth.

## Smoothing

To get perfectly smooth look, there are 2 possible solutions:

1. Mouse polling, engine update rate and monitor refresh rate must be aligned - for example 360Hz mouse polling, 360 FPS engine update rate, 360Hz monitor refresh rate. Not realistic.
2. Input smoothing - gives almost perfectly smooth results (the difference is noticeable also on high-end gaming setup) at the cost of increased input lag by few milliseconds.

The KCC provides utility script ```
Vector2Accumulator
```

. Following code shows its usage to calculate smooth mouse position delta from values in last 20 milliseconds.

C#

```
```csharp
public class PlayerInput : MonoBehaviour
{
// Creates an accumulator with 20ms smoothing window.
private Vector2Accumulator \_lookRotationAccumulator = new Vector2Accumulator(0.02f, true);

private void Update()
{
Vector2 mouseDelta = Mouse.current.delta.ReadValue();

\_lookRotationAccumulator.Accumulate(mouseDelta);
}

private void PollInput(CallbackPollInput callback)
{
Quantum.Input input = new Quantum.Input();

// 1\. Option - consume whole smoothed mouse delta which is aligned to render time.
// input.LookRotationDelta = \_lookRotationAccumulator.Consume();

// 2\. Option (better) - consume smoothed mouse delta which is aligned to Quantum frame time.
// This variant ensures smooth interpolation when look rotation propagates to transform.
input.LookRotationDelta = \_lookRotationAccumulator.ConsumeFrameAligned(callback.Game);

callback.SetInput(input, DeterministicInputFlags.Repeatable);
}
}

```

```

Following images show mouse delta being propagated to **accumulated look rotation** (roughly 90°) with various mouse polling rates, engine update rates and smoothing windows:

- ```
Purple
```

\- No smoothing applied.
- ```
Cyan
```

\- 10ms smoothing window.
- ```
Green
```

\- 20ms smoothing window.
- ```
Yellow
```

\- 30ms smoothing window.
- ```
Blue
```

\- 40ms smoothing window.

![Accumulated look rotation with 125Hz mouse polling and 120 FPS](/docs/img/quantum/v3/addons/kcc/accumulated-input-125hz-120fps.jpg)Accumulated look rotation with 125Hz mouse polling and 120 FPS.

![Accumulated look rotation with 125Hz mouse polling and 360 FPS](/docs/img/quantum/v3/addons/kcc/accumulated-input-125hz-360fps.jpg)Accumulated look rotation with 125Hz mouse polling and 360 FPS.

![Accumulated look rotation with 1000Hz mouse polling and 120 FPS](/docs/img/quantum/v3/addons/kcc/accumulated-input-1000hz-120fps.jpg)Accumulated look rotation with 1000Hz mouse polling and 120 FPS.

![Accumulated look rotation with 1000Hz mouse polling and 360 FPS](/docs/img/quantum/v3/addons/kcc/accumulated-input-1000hz-360fps.jpg)Accumulated look rotation with 1000Hz mouse polling and 360 FPS.

Generally for regular users with 125Hz mouse a 10-20ms smoothing window is a good balance between gain in smoothness and loss in responsivity.

For users with proper gaming hardware (500+Hz mouse, 120+Hz monitor) reaching high engine update rates 3-5ms smoothing window is recommended + option to disable smoothing completely.

Following image shows a detail of look rotation after accumulation of 90°.

![Look rotation after accumulation of 90°](/docs/img/quantum/v3/addons/kcc/accumulated-input-detail.jpg)Look rotation after accumulation of 90°.

You can see that applying smoothing adds 30-50% of smoothing window length to input lag (+3.6ms input lag using 10ms smoothing window).

Back to top

- [Getting Input](#getting-input)
- [Smoothing](#smoothing)