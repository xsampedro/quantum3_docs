# animation-curves-baking

_Source: https://doc.photonengine.com/quantum/current/concepts-and-patterns/animation-curves-baking_

# Baking Unity Animation Curves

Even though Quantum's ```
FPAnimationCurve
```

 already does a conversion between Unity's non deterministic type under the hood through the curves editor, sometimes it might be useful to convert from an ```
AnimationCurve
```

when there is no automatic-in-editor conversion available.

One example is the conversion of curves which comes from Unity'a animation Clips, to create a deterministic version of them which could be used on the simulation.

Here are the snippets needed to create an ```
FPAnimationCurve
```

 from an ```
AnimationCurve
```

:

C#

```
```csharp
public FPAnimationCurve ConvertAnimationCurve(AnimationCurve animationCurve)
{
// Get UNITY keyframes
Keyframe\[\] unityKeys = animationCurve.keys;

// Prepare QUANTUM curves and keyframes to receive the info
FPAnimationCurve fpCurve = new FPAnimationCurve();
fpCurve.Keys = new FPAnimationCurve.Keyframe\[unityKeys.Length\];

// Get the Unity Start and End time for this specific curve
float startTime = animationCurve.keys.Length == 0 ? 0.0f : float.MaxValue;
float endTime = animationCurve.keys.Length == 0 ? 1.0f : float.MinValue;

// Set the resolution for the curve, which informs how detailed it is
fpCurve.Resolution = 32;

for (int i = 0; i < unityKeys.Length; i++)
{
fpCurve.Keys\[i\].Time = FP.FromFloat\_UNSAFE(unityKeys\[i\].time);
fpCurve.Keys\[i\].Value = FP.FromFloat\_UNSAFE(unityKeys\[i\].value);

if (float.IsInfinity(unityKeys\[i\].inTangent) == false)
{
fpCurve.Keys\[i\].InTangent = FP.FromFloat\_UNSAFE(unityKeys\[i\].inTangent);
}
else
{
fpCurve.Keys\[i\].InTangent = FP.SmallestNonZero;
}

if (float.IsInfinity(unityKeys\[i\].outTangent) == false)
{
fpCurve.Keys\[i\].OutTangent = FP.FromFloat\_UNSAFE(unityKeys\[i\].outTangent);
}
else
{
fpCurve.Keys\[i\].OutTangent = FP.SmallestNonZero;
}

fpCurve.Keys\[i\].TangentModeLeft = (byte)AnimationUtility.GetKeyLeftTangentMode(animationCurve, i);
fpCurve.Keys\[i\].TangentModeRight = (byte)AnimationUtility.GetKeyRightTangentMode(animationCurve, i);

startTime = Mathf.Min(startTime, animationCurve\[i\].time);
endTime = Mathf.Max(endTime, animationCurve\[i\].time);
}

fpCurve.StartTime = FP.FromFloat\_UNSAFE(startTime);
fpCurve.EndTime = FP.FromFloat\_UNSAFE(endTime);

fpCurve.PreWrapMode = (int)animationCurve.preWrapMode;
fpCurve.PostWrapMode = (int)animationCurve.postWrapMode;

// Actually save the many points of the unity curve into the quantum curve
SaveQuantumCurve(animationCurve, 32, ref fpCurve, startTime, endTime);
return fpCurve;
}

private void SaveQuantumCurve(AnimationCurve animationCurve, int resolution, ref FPAnimationCurve fpCurve, float startTime, float endTime)
{
if (resolution <= 0)
return;
fpCurve.Samples = new FP\[resolution + 1\];
var deltaTime = (endTime - startTime) / (float)resolution;
for (int i = 0; i < resolution + 1; i++)
{
var time = startTime + deltaTime \* i;
var fp = FP.FromFloat\_UNSAFE(animationCurve.Evaluate(time));
fpCurve.Samples\[i\].RawValue = fp.RawValue;
}
}

```

```

Back to top