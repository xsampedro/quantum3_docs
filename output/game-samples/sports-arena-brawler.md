# sports-arena-brawler

_Source: https://doc.photonengine.com/quantum/current/game-samples/sports-arena-brawler_

# Sports Arena Brawler

![Level 4](/v2/img/docs/levels/level03-intermediate_1.5x.png)

## Overview

The Quantum Sport Arena Brawler sample is a top-down 3v3 sports arena brawler. Pass the ball, punch opponents off of the arena, and score against the enemy team in chaotic, lightning fast matches. It supports up to 4 local players via split screen. Input buffering and ability activation delays allow for a smooth multiplayer experience at higher pings.

_This sample was developed by the MicroverseLabs studio._

## Download

| Version | Release Date | Download |
| --- | --- | --- |
| 3.0.2 | Mar 17, 2025 | [Quantum Sports Arena Brawler 3.0.2 Build 598](https://dashboard.photonengine.com/download/quantum/quantum-sports-arena-brawler-3.0.2.zip) |

## Technical Info

- Unity: 2021.3.18f1.
- Platforms: PC (Windows / Mac)

## Highlights

### Technical

- Multiple Local Players leveraging the default Quantum features.
- Input encoding (Vector2 as Byte).
- Custom interpolation for fast moving ball in view.
- Splitscreen Multiplayer (Local + Online).

### Gameplay

- Different sets of abilities.
- Available abilities change depending on ball possession.
- Multiple Local Players.
- Coyote Time.

## Screenshots

![](https://doc.photonengine.com/docs/img/quantum/v2/game-samples/qball/punch.jpg)

![](https://doc.photonengine.com/docs/img/quantum/v2/game-samples/qball/passing.jpg)

![](https://doc.photonengine.com/docs/img/quantum/v2/game-samples/qball/score.jpg)

![](https://doc.photonengine.com/docs/img/quantum/v2/game-samples/qball/split-screen.jpg)

## Local Players

### UI And Matchmaking

The sample uses a modified version of the default Quantum Demo Manu. The main addition is a local players count dropdown on the connect screen.

![Game Start](https://doc.photonengine.com/docs/img/quantum/v3/game-samples/qball/qball-connect-ui.jpg)

A `SqlLobbyFilter` is used when starting the connection in order to limit the maximum number of players to 6 while also taking into account local players.

C#

```csharp
public const string LOCAL_PLAYERS_PROP_KEY = &#34;LP&#34;;
public const string TOTAL_PLAYERS_PROP_KEY = &#34;C0&#34;;

public static readonly TypedLobby SQL_LOBBY = new TypedLobby(&#34;customSqlLobby&#34;, LobbyType.SqlLobby);

```

C#

```csharp
protected override void OnConnect(QuantumMenuConnectArgs connectArgs, ref MatchmakingArguments args)
{
    args.RandomMatchingType = MatchmakingMode.FillRoom;
    args.Lobby = LocalPlayerCountManager.SQL_LOBBY;
    args.CustomLobbyProperties = new string[] { LocalPlayerCountManager.TOTAL_PLAYERS_PROP_KEY };
    args.SqlLobbyFilter = $&#34;{LocalPlayerCountManager.TOTAL_PLAYERS_PROP_KEY} <= {Input.MAX_COUNT - _localPlayersCountSelector.GetLastSelectedLocalPlayersCount()}&#34;;
}

```

After connecting to a room the master client keeps a custom property up to date with the total players count (including additional local players for all connected clients).

C#

```csharp
private void UpdateRoomTotalPlayers()
{
    if (_connection != null && _connection.Client.InRoom && _connection.Client.LocalPlayer.IsMasterClient)
    {
        int totalPlayers = 0;
        foreach (var player in _connection.Client.CurrentRoom.Players.Values)
        {
            if (player.CustomProperties.TryGetValue(LOCAL_PLAYERS_PROP_KEY, out var localPlayersCount))
            {
                totalPlayers += (int)localPlayersCount;
            }
        }

        _connection.Client.CurrentRoom.SetCustomProperties(new PhotonHashtable
        {
            { TOTAL_PLAYERS_PROP_KEY, totalPlayers }
        });
    }
}

```

### Local Players Initialization

When the gameplay starts a different configuration prefab is instantiated depending on the amount of local players. Inside the configuration prefab each local player has their own `Camera`, `UI` and `PlayerInput`. `PlayerInput` automatically takes care of assigning different input devices for each local player. If there are multiple local players the main player always gets assigned mouse and keyboard and any additional players each get a different controller (controllers need to be plugged in before the gameplay starts).

## Abilities

### Overview

The state data of each ability is stored inside an `Ability` struct that holds a few timers and an `AbilityData` asset reference.

Qtn

```cs
struct Ability
{
    [ExcludeFromPrototype] AbilityType AbilityType;

    [ExcludeFromPrototype] CountdownTimer InputBufferTimer;
    [ExcludeFromPrototype] CountdownTimer DelayTimer;
    [ExcludeFromPrototype] CountdownTimer DurationTimer;
    [ExcludeFromPrototype] CountdownTimer CooldownTimer;

    asset_ref<AbilityData> AbilityData;
}

```

The `Ability` structs are stored inside of an array in the `AbilityInventory` component.

Qtn

```cs
component AbilityInventory
{
    [ExcludeFromPrototype] ActiveAbilityInfo ActiveAbilityInfo;

    // Same order as AbilityType enum also used for activation priority
    [Header(&#34;Ability Order: Block, Dash, Attack, ThrowShort, ThrowLong, Jump&#34;)]
    array<Ability>[6] Abilities;
}

```

A single `AbilitySystem` takes care of updating all abilities in a data-driven way by passing all relevant state data to their corresponding `AbilityData` assets.

C#

```csharp
public override void Update(Frame frame, ref Filter filter)
{
    QuantumDemoInputTopDown input = *frame.GetPlayerInput(filter.PlayerStatus->PlayerRef);

    for (int i = 0; i < filter.AbilityInventory->Abilities.Length; i++)
    {
        AbilityType abilityType = (AbilityType)i;
        ref Ability ability = ref filter.AbilityInventory->Abilities[i];
        AbilityData abilityData = frame.FindAsset<AbilityData>(ability.AbilityData.Id);

        abilityData.UpdateAbility(frame, filter.EntityRef, ref ability);
        abilityData.UpdateInput(frame, ref ability, input.GetAbilityInputWasPressed(abilityType));
        abilityData.TryActivateAbility(frame, filter.EntityRef, filter.PlayerStatus, ref ability);
    }
}

```

The base `AbilityData` implementation takes care of activating the ability logic and updating its state while all immutable ability specific data and logic is implemented in derived `AbilityData` assets by using polymorphism. This setup allows for all ability logic to be self-contained and creating new abilities becomes as simple as writing their unique logic without the need of any boilerplate code.

### Input Buffering

When ability input is detected instead of trying to activate the ability right away an InputBufferTimer is started. The `AbilityData` then checks if the timer is running each frame in order to activate the ability. That allows both for a smoother player experience and helps mitigate high latency in some situations. E.g. If the player is in the middle of a dash and tries to throw the ball, their input normally would be consumed without anything happening - the input buffering queues the throw ability to be activated as soon as the dash ends and is also sent to the other remote players a bit earlier so it can arrive in time and prevent mispredictions.

### Activation Delay

When an ability is activated it first enters a delayed state that gives some time for the input to reach other remote players and prevent mispredictions. In order for the abilities to feel responsive for the local player their animations are triggered instantly and last for the whole delay + actual duration.

### Different Abilities When Holding The Ball

Without the ball the player has access to an offensive punch and a defensive block abilities. When holding the ball they are replaced by short and long throw abilities. Unavailable abilities still get updated each tick so their `InputBufferTimer` and `CooldownTimer` can be ticked down. This in combination with a reduced movement speed when holding the ball incentives passing it or relying on teammates for protection.

### Punch

The punch ability uses a compound hit detection shape with multiple growing in size spheres in order to create a cone-shaped hitbox. It applies both a knockback and a stun status effect. Knockbacks can be chained together by multiple players and being knockbacked into the void results in a short timeout followed by a respawn.

### Block Ability

The block ability completely prevents all attacks while it lasts.

### Throw Abilities

Since all abilities are targeted just by a single aim direction, a short and a long pass allow for better control.

### Dash Ability

Dashing allows for rapid movement driven by an animation curve. Any custom movement needs to be calculated relative to the current player position in order to allow multiple custom movements on top of each other and KCC collider penetration correction.

C#

```csharp
if (abilityState.IsActive)
{
    AbilityInventory* abilityInventory = frame.Unsafe.GetPointer<AbilityInventory>(entityRef);
    Transform3D* transform = frame.Unsafe.GetPointer<Transform3D>(entityRef);
    CharacterController3D* kcc = frame.Unsafe.GetPointer<CharacterController3D>(entityRef);

    FP lastNormalizedPosition = DashMovementCurve.Evaluate(lastNormalizedTime);
    FPVector3 lastRelativePosition = abilityInventory->ActiveAbilityInfo.CastDirection * DashDistance * lastNormalizedPosition;

    FP newNormalizedTime = ability.DurationTimer.NormalizedTime;
    FP newNormalizedPosition = DashMovementCurve.Evaluate(newNormalizedTime);
    FPVector3 newRelativePosition = abilityInventory->ActiveAbilityInfo.CastDirection * DashDistance * newNormalizedPosition;

    transform->Position += newRelativePosition - lastRelativePosition;
}

```

### Jump Ability

Jumping is also implemented as an ability so it can benefit from input buffering and activation delay. Input buffering is especially useful for it because it allows for the next jump to be queued shortly before becoming grounded. The activation delay for jumping is much lower than other abilities because it would feel unresponsive otherwise.

## Character Controller

### KCC Configuration

There are 3 different KCC configurations that get applied depending on the state of the player. The first one is just the default behavior and allows for normal movement. The second one is used when holding the ball and it reduces the movement speed and jump height of the player. The third one is applied during ability usage and while being knockbacked. It prevents all input based movement and gravity and allows for full control via code. The `KCC->Move()` method is still executed in order to prevent the player from going inside of obstacles when moved by code.

C#

```csharp
public unsafe void UpdateKCCSettings(Frame frame, EntityRef playerEntityRef)
{
    PlayerStatus* playerStatus = frame.Unsafe.GetPointer<PlayerStatus>(playerEntityRef);
    AbilityInventory* abilityInventory = frame.Unsafe.GetPointer<AbilityInventory>(playerEntityRef);
    CharacterController3D* kcc = frame.Unsafe.GetPointer<CharacterController3D>(playerEntityRef);

    CharacterController3DConfig config;

    if (playerStatus->IsKnockbacked || abilityInventory->HasActiveAbility)
    {
        config = frame.FindAsset<CharacterController3DConfig>(NoMovementKCCSettings.Id);
    }
    else if (playerStatus->IsHoldingBall)
    {
        config = frame.FindAsset<CharacterController3DConfig>(CarryingBallKCCSettings.Id);
    }
    else
    {
        config = frame.FindAsset<CharacterController3DConfig>(DefaultKCCSettings.Id);
    }

    kcc->SetConfig(frame, config);
}

```

### Coyote Time

In order to achieve a better game feel and to minimize player mistakes when jumping between platforms there is a "coyote time" mechanic. It allows the player to jump normally shortly after becoming airborne. Every tick while the player is grounded a `JumpCoyoteTimer` is started. When the player tries to jump instead of checking if grounded we check if the `JumpCoyoteTimer.IsRunning` instead.

## Ball

### View Interpolation

While the ball is held its real position is in the center of the player, its physics are disabled and it is not manipulated further in any way. That allows the view to temporarily take control over the ball and move its graphics via animations. As soon as the ball is caught or released by the player, its transform quickly gets interpolated between real space and animated space.

C#

```csharp
public unsafe class BallEntityView : QuantumEntityView
{
    private float _interpolationSpaceAlpha;

    public void UpdateSpaceInterpolation()
    {
        // . . .
        UpdateInterpolationSpaceAlpha(isBallHeldByPlayer);

        if (_interpolationSpaceAlpha > 0f)
        {
            Vector3 interpolatedPosition = Vector3.Lerp(_lastBallRealPosition, _lastBallAnimationPosition, _interpolationSpaceAlpha);
            Quaternion interpolatedRotation = Quaternion.Slerp(_lastBallRealRotation, _lastBallAnimationRotation, _interpolationSpaceAlpha);

            transform.SetPositionAndRotation(interpolatedPosition, interpolatedRotation);
        }
    }

    private void UpdateInterpolationSpaceAlpha(bool isBallHeldByPlayer)
    {
        float deltaChange = _spaceTransitionSpeed * Time.deltaTime;
        if (isBallHeldByPlayer)
        {
            _interpolationSpaceAlpha += deltaChange;
        }
        else
        {
            _interpolationSpaceAlpha -= deltaChange;
        }

        _interpolationSpaceAlpha = Mathf.Clamp(_interpolationSpaceAlpha, 0f, 1f);
    }
}

```

### Gravity Scale

When the ball is thrown in order to allow for low passing without a parabola and without drastically increasing the throw force, the ball temporarily is not affected by gravity. After the ball is thrown its `GravityScale` quickly gets interpolated from 0 to 1 using a curve in order to give the control back to the physics system and achieve more realistic results.

C#

```csharp
private void UpdateBallGravityScale(Frame frame, ref Filter filter, BallHandlingData ballHandlingData)
{
    if (filter.BallStatus->GravityChangeTimer.IsRunning)
    {
        FP gravityScale = ballHandlingData.ThrowGravityChangeCurve.Evaluate(filter.BallStatus->GravityChangeTimer.NormalizedTime);
        filter.PhysicsBody->GravityScale = gravityScale;

        filter.BallStatus->GravityChangeTimer.Tick(frame.DeltaTime);
        if (filter.BallStatus->GravityChangeTimer.IsDone)
        {
            ResetBallGravity(frame, filter.EntityRef);
        }
    }
}

```

### Custom Lateral Friction

Additional lateral friction is applied to the ball when it is bouncing / rolling on the ground in order to apply a more precise control over its travel distance when thrown and to prevent it from constantly rolling over the edge into the void.

C#

```csharp
public void OnCollisionEnter3D(Frame frame, CollisionInfo3D info)
{
    if (frame.Unsafe.TryGetPointer(info.Entity, out BallStatus* ballStatus))
    {
        ballStatus->HasCollisionEnter = true;
    }
}

public void OnCollision3D(Frame frame, CollisionInfo3D info)
{
    if (frame.Unsafe.TryGetPointer(info.Entity, out BallStatus* ballStatus))
    {
        ballStatus->HasCollision = true;
    }
}

private void HandleBallCollisions(Frame frame, ref Filter filter, BallHandlingData ballHandlingData)
{
    if (!filter.PhysicsBody->IsKinematic)
    {
        if (filter.BallStatus->HasCollisionEnter)
        {
            filter.PhysicsBody->Velocity.X *= ballHandlingData.LateralBounceFriction;
            filter.PhysicsBody->Velocity.Z *= ballHandlingData.LateralBounceFriction;

            frame.Events.OnBallBounced(filter.EntityRef);
        }

        if (filter.BallStatus->HasCollision)
        {
            filter.PhysicsBody->Velocity.X *= ballHandlingData.LateralGroundFriction;
            filter.PhysicsBody->Velocity.Z *= ballHandlingData.LateralGroundFriction;
        }
    }

    filter.BallStatus->HasCollisionEnter = false;
    filter.BallStatus->HasCollision = false;
}

```

## Input

Input is handled by Unity's Input System package. On the Quantum code side the `QuantumDemoInputTopDown` as base.

C#

```csharp
// DSL Definition
[ExcludeFromPrototype]
struct QuantumDemoInputTopDown {
    FPVector2 MoveDirection;
    FPVector2 AimDirection;
    button Left;
    button Right;
    button Up;
    button Down;
    button Jump;
    button Dash;
    button Fire;
    button AltFire;
    button Use;
}

```

C#

```csharp
// Example use case
public override void Update(Frame frame, ref Filter filter)
{
    QuantumDemoInputTopDown input = *frame.GetPlayerInput(filter.PlayerStatus->PlayerRef);
    // . . .
}

```

## Camera

The camera is controlled by `Cinemachine` via a `CinemachineTargetGroup` in order to focus on all actors using a higher weight for the local players and a larger radius for the ball so all the action can be framed with ease.

## Third Party Assets

This sample includes third-party free assets. The full packages can be acquired for your own projects at their respective site:

- [Pixabay Sound Effects](https://pixabay.com/service/license-summary/) by Pixabay
- [Controller Icon Pack](https://assetstore.unity.com/packages/2d/gui/icons/controller-icon-pack-128505) by NullSave
- [Fonts](https://fonts.google.com/specimen/Bangers/about) by Google Fonts
- [Fonts](https://www.fontsquirrel.com/license/liberation-sans) by Font Squirrel

Back to top

- [Overview](#overview)
- [Download](#download)
- [Technical Info](#technical-info)
- [Highlights](#highlights)

  - [Technical](#technical)
  - [Gameplay](#gameplay)

- [Screenshots](#screenshots)
- [Local Players](#local-players)

  - [UI And Matchmaking](#ui-and-matchmaking)
  - [Local Players Initialization](#local-players-initialization)

- [Abilities](#abilities)

  - [Overview](#overview-1)
  - [Input Buffering](#input-buffering)
  - [Activation Delay](#activation-delay)
  - [Different Abilities When Holding The Ball](#different-abilities-when-holding-the-ball)
  - [Punch](#punch)
  - [Block Ability](#block-ability)
  - [Throw Abilities](#throw-abilities)
  - [Dash Ability](#dash-ability)
  - [Jump Ability](#jump-ability)

- [Character Controller](#character-controller)

  - [KCC Configuration](#kcc-configuration)
  - [Coyote Time](#coyote-time)

- [Ball](#ball)

  - [View Interpolation](#view-interpolation)
  - [Gravity Scale](#gravity-scale)
  - [Custom Lateral Friction](#custom-lateral-friction)

- [Input](#input)
- [Camera](#camera)
- [Third Party Assets](#third-party-assets)