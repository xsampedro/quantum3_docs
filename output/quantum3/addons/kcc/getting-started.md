# getting-started

_Source: https://doc.photonengine.com/quantum/current/addons/kcc/getting-started_

# Getting Started

## Creating character with KCC prefab variant

The easiest way to create custom character is to create a prefab variant from default ```
KCC
```

 prefab:

1. Right click on ```
   KCC
   ```

    prefab (in ```
   Assets/Photon/QuantumAddons/KCC/AssetDB/Entities
   ```

   ).
2. Select ```
   Create > Prefab Variant
   ```

   .

![Create prefab variant](/docs/img/quantum/v3/addons/kcc/create-prefab-variant.jpg)Create KCC prefab variant

3. Add your own visual, add custom components.
4. The character is ready for use. Continue with [Moving the character](/quantum/current/addons/kcc/getting-started#moving-the-character)

## Creating character from scratch

1. Create a new player prefab.

![Create player prefab](/docs/img/quantum/v3/addons/kcc/create-player-prefab.jpg)Create player prefab

2. Add ```
   Quantum Entity View
   ```

    and ```
   Q Prototype KCC
   ```

    components to the root game object. Optionally add ```
   Capsule Collider
   ```

   .

![Add components](/docs/img/quantum/v3/addons/kcc/add-kcc-components.jpg)Add components to player prefab

3. Configure the game object

- Set game object layer (optional).
- On ```
  Quantum Entity View
  ```

   set ```
  Bind Behaviour
  ```

   to ```
  Verified
  ```

  .
- On ```
  Quantum Entity Prototype
  ```

   set ```
  Transform
  ```

   to ```
  3D
  ```

  .
- Enable ```
  PhysicsCollider3D
  ```

   and link previously created ```
  Capsule Collider
  ```

   to ```
  SourceCollider
  ```

   (optional).
- On ```
  Q Prototype KCC
  ```

   set reference to a ```
  KCC Settings
  ```

   asset.

4. The character is ready for use. Continue with [Moving the character](/quantum/current/addons/kcc/getting-started#moving-the-character)

## Configuring KCC behavior

The KCC addon contains some pre-configured assets.

1. Select ```
   Assets/Photon/QuantumAddons/KCC/AssetDB/KCCSettings.asset
   ```

    to configure defaults (radius, height, collision layer mask). These should match values in capsule collider created above.

![Configure KCC settings](/docs/img/quantum/v3/addons/kcc/kcc-settings.jpg)Configure KCC settings

2. Link KCC processors if you are creating new ```
   KCC Settings
   ```

    asset. These are responsible for actual movement logic, more info can be found in [Processors](/quantum/current/addons/kcc/processors) section. Default processors are located in ```
   Assets/Photon/QuantumAddons/KCC/AssetDB/Processors
   ```

   .

## Moving the character

The movement is processed in ```
KCC
```

component update. This is managed by ```
KCC System
```

, therefore it must be added to your ```
Systems Config
```

.

![Add KCC System](/docs/img/quantum/v3/addons/kcc/kcc-system.jpg)Add KCC System

Following code example sets character look rotation and input direction for ```
KCC
```

 which is later processed by ```
EnvironmentProcessor
```

(or your own processor).

C#

```
```csharp
public unsafe class PlayerSystem : SystemMainThreadFilter<PlayerSystem.Filter>
{
public struct Filter
{
public EntityRef Entity;
public Player\* Player;
public KCC\* KCC;
}

public override void Update(Frame frame, ref Filter filter)
{
Player\* player = filter.Player;
if (player->PlayerRef.IsValid == false)
return;

KCC\* kcc = filter.KCC;
Input\* input = frame.GetPlayerInput(player->PlayerRef);

kcc->AddLookRotation(input->LookRotationDelta.X, input->LookRotationDelta.Y);
kcc->SetInputDirection(kcc->Data.TransformRotation \* input->MoveDirection.XOY);

if (input->Jump.WasPressed == true && kcc->IsGrounded == true)
{
kcc->Jump(FPVector3.Up \* player->JumpForce);
}
}
}

```

```

It is also possible to skip processors functionality completely and simply set the velocity using ```
kcc->SetKinematicVelocity();
```

.

More movement code examples can be found in [Sample Project](/quantum/current/addons/kcc/sample-project).

Back to top

- [Creating character with KCC prefab variant](#creating-character-with-kcc-prefab-variant)
- [Creating character from scratch](#creating-character-from-scratch)
- [Configuring KCC behavior](#configuring-kcc-behavior)
- [Moving the character](#moving-the-character)