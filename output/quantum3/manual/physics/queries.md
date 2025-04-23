# queries

_Source: https://doc.photonengine.com/quantum/current/manual/physics/queries_

# Queries

## Introduction

Queries may take into account dynamic entities and static colliders. The API for rays, lines and shape overlaps is very similar in that it always returns a collection of hits (with the same kind of data in their fields).

Every Physics query can be customized with optional flags which makes them very flexible and optimized. For example, by default most of the queries _does not return the hit Point or Normal_ and which has to be explicitly defined by the user (explained in the Options topic later in this document).

## Queries

### Linecast and Raycast

C#

```csharp
// For 2D
var hits = f.Physics2D.LinecastAll(FPVector2.Zero, FPVector2.One);
for (int i = 0; i < hits.Count; i++) {
    var hit = hits[i];
}
// For 3D
var hits = f.Physics3D.LinecastAll(FPVector3.Zero, FPVector3.One);
for (int i = 0; i < hits.Count; i++){
    var hit = hits[i];
}

```

The resulting _HitCollection_ object, contains the following properties:

- Each item in the _HitCollection_ holds an EntityRef or Static collider info. They mutually exclusive - one will be valid, and the other _null_;
- _Count_ should always be used to iterate over the _HitCollection_; and,
- Hits are not sorted. You can sort them by calling `Sort()` and passing in a _FPVector2_, this will result in the hits being sorted according to their distance to the reference point provided to the function.

Raycasts are syntax-sugar for Linecasts. They work the same and simply require a _start_, _direction_ and _max-distance_ instead of _start_ and _end_. Additionally, it is possible to pass optional parameters to the linecast and raycast:

- LayerMask, to specify which physics layers to perform the cast against; and,
- QueryOptions, to specify the type of collider to consider in the cast.

### Shape Queries

Quantum supports two different types of shape queries:

- ShapeOverlap; and,
- ShapeCasts.

These can be used with all dynamic shapes supported in Quantum.

**Note:**`CompoundShapes` can be used for performing shape queries. For more information, please read the _Shape Config_ page.

#### ShapeOverlaps

`OverlapShape()` returns a _HitCollection_. The required parameters are:

- a center position ( _FPVector2_ or _FPVector3_);
- a rotation ( _FP_ or _FPQuaternion_ for the 3D equivalent); and,
- a shape ( _Shape2D_ or _Shape3D_ \- either from a _PhysicsCollider_, or created at the time of calling).

C#

```csharp
// For 2D
var hits = f.Physics2D.OverlapShape(FPVector2.Zero, FP._0, Shape2D.CreateCircle(FP._1));
for (int i = 0; i < hits.Count; i++){
    var hit = hits[i];
}
// For 3D
var hits = f.Physics3D.OverlapShape(FPVector3.Zero, FPQuaternion.Identity, Shape3D.CreateSphere(1));
for (int i = 0; i < hits.Count; i++){
    var hit = hits[i];
}

```

#### ShapeCasts

`ShapeCastAll()` returns a _HitCollection_. The required parameters are:

- the center position ( _FPVector2_ or _FPVector3_);
- the rotation of the shape ( _FP_ or _FPQuaternion_ for the 3D equivalent);
- the shape pointer ( \_Shape2D\* \_ or _Shape3D\*_ \- either from a _PhysicsCollider_, or created at the time of calling); and,
- the distance and direction expressed as a vector ( _FPVector2_ or _FPVector3_ ).

C#

```csharp
// For 2D
var shape = Shape2D.CreateCircle(FP._1);
var hits = f.Physics2D.ShapeCastAll(FPVector2.Zero, FP._0, &shape, FPVector2.One);
for (int i = 0; i < hits.Count; i++){
    var hit = hits[i];
}
// For 3D
var shape = Shape3D.CreateSphere(1);
var hits = f.Physics3D.ShapeCastAll(FPVector3.Zero, FPQuaternion.Identity, &shape, FPVector3.One);
for (int i = 0; i < hits.Count; i++){
    var hit = hits[i];
}

```

The shape casts uses a custom GJK-based algorithm. If the initial shape cast position is contained within a collider, the shape cast cannot detect a hit with such collider unless the `DetectOverlapsAtCastOrigin` flag in the `QueryOptions` parameter is enabled, which makes it perform an extra check at the start position.

The GJKConfig settings are available in the `SimulationConfig` asset's `Physics > GJKConfig` section. The settings allow to balance accuracy and performance as both come with their trade-offs. The default values are balanced to compromise for regular sized shapes.

- `Simplex Min/Max Bit Shift`: Allows better precision for points in the Voronoy Simplex by progressively shifting their raw values, avoiding degenerate cases without compromising the valid range of positions in the Physics space. Consider increasing the values if the scale of the shapes involved and/or the distance between them is very small.
- `Shape Cast Max Iterations`: The max number of iterations performed by the algorithm while searching for a solution below the hard tolerance. Increasing it might result in more accurate results, at the cost of performance in worst-case scenarios, and vice-versa.
- `Shape Cast Hard Tolerance`: An iteration result (closest distance between the shapes) below this threshold is acceptable as a finishing condition. Decreasing it might result in more accurate results, at the cost of more iterations, and vice-versa.
- `Shape Cast Soft Tolerance`: A shape cast resolution that fails to find an acceptable result below the defined Hard Tolerance within the Max Iterations allowed will still return positive if the best result found so far is below this soft threshold. In these cases, increasing this threshold enhances the probability of false-positives, while decreasing it enhances false-negatives.

### Sorting Hits

All queries returning a `HitCollection` can be sorted.

- `Sort()`: takes a FPVector2 in 2D and a FPVector3 in 3D and sorts the collection according to the hits' respective distance to the point provided.
- `SortCastDistance()`: used for sorting the results of `ShapeCast` query. It takes no arguments and orders the hits based on the cast distance.

## Options

All queries, including their broadphase version, can use `QueryOptions ` to customize the operation and its results.

`QueryOptions` create a mask that filters which types of objects are taken into account and what information will be computed. Is is possible to combined these by using the binary `\|` operator.

### Hit Normals

To offer the most performant query, all default queries only check whether the two shapes are overlapping. By default, most of the queries does not retrieve the hit Point or Normal, for example.

In order to receive additional information, more computation is needed which in turn creates additional overhead; it is therefore necessary to explicitly specify it by passing `ComputeDetailedInfo` as the QueryOptions parameter. This will enable the computation of the hit's:

- point
- normal
- penetration

For ray-triangle checks the normal is always the triangle's normal. Since this is cached in the triangle data, there is no additional computation in this case.

### Filtering Hits

The following `QueryOptions` allow you to define the mask used by the query. If an object does not match the `QueryOptions` specified as the parameter, it will be skipped; only objects matching the `QueryOptions` will be evaluated and returned in the result.

- **HitStatics** : will only hit static colliders

- **HitKinematics** : will hit entities who meet any of the following conditions:

  - entities with a PhysicsCollider and **no** PhysicsBody
  - entities with a PhysicsCollider and a **disabled** PhysicsBody
  - entities with a PhysicsCollider and a **kinematic** PhysicsBody
- **HitDynamics** : will only hit entities with an _enabled_ and _non-kinematic_ PhysicsBody

- **HitTriggers** : has to be used **in combination with** other flags to hit trigger colliders.

- **HitAll** : will hit all entities that have a PhysicsCollider


By default, a query will use the `HitAll` option. Choosing any other option will save computation.

## Broadphase Queries

Quantum comes with an option for injecting physics queries (ray-casts and overlaps) to be resolved during the physics systems. For this you need to:

1. Create a system;
2. Add it to the uses Systems Config asset, before the `Core.PhysicsSystem`;
3. Retrieve the information in any system running after `Core.PhysicsSystem`.

This setup benefits from the parallel resolution on physics steps which makes it significantly faster than normal querying after physics.

![Broadphase Queries System](/docs/img/quantum/v3/manual/physics/broadphase-queries-system.png)

_Note:_ Sometimes broadphase queries are also referred to as injected queries or scheduled queries, because they are scheduled/injected into the physics engine before the solver runs.

### Injecting Queries

It is possible to inject a query from any main thread system running before physics. Injecting a query will return a `PhysicsQueryRef`, which can be stored and used to retrieve the results after the physics system ran. Broad-phase query results are meant to be retrieve within the same frame the query was injected, therefore the `PhysicsQueryRef` can be stored anywhere - including outside the rollback-able frame data.

C#

```csharp
namespace Quantum
{
    public unsafe struct ProjectileFilter
    {
        public EntityRef EntityRef;
        public Transform3D* Transform;
        public Projectile* Component;
    }

    public unsafe class ProjectileHitQueryInjectionSystem : SystemMainThread
    {
        public override void Update(Frame frame)
        {
            var projectileFilter = frame.Unsafe.FilterStruct<ProjectileFilter>();
            var projectile = default(ProjectileFilter);

            while (projectileFilter.Next(&projectile))
            {
                projectile.Component->PathQueryRef = frame.Physics3D.AddRaycastQuery(
                    projectile.Transform->Position,
                    projectile.Transform->Forward,
                    projectile.Component->Speed * frame.DeltaTime);
                var spec = frame.FindAsset<WeaponSpec>(projectile.Component->WeaponSpec.Id);

                projectile.Component->DamageZoneQueryRef = frame.Physics3D.AddOverlapShapeQuery(
                    projectile.Transform->Position,
                    projectile.Transform->Rotation,
                    spec.AttackShape.CreateShape(frame),
                    spec.AttackLayers);
            }
        }
    }
}

```

**IMPORTANT:** The `PhysicsQueryRef` returned by `AddXXXQuery` is absolutely necessary to retrieve the results of the query later on. It is thus advisable to save in a component attached to the entity who will need to process the hits down the line.

### Retrieving Query Results

The query results can be retrieved from any system that runs after the core physics system. To retrieve the results (HitCollection\*), pass the index previous saved into either `Frame.Physics.GetQueryHits()` or `.TryGetQueryHits()`.

Attempting to retrieve the results with an invalid `PhysicsQueryRef` (e.g. from a query injected in a different frame) will throw an exception in `GetQueryHits` or return false in `TryGetQueryHits`.

C#

```csharp
using Photon.Deterministic;
namespace Quantum
{
    public unsafe class ProjectileHitRetrievalSystem : SystemMainThread
    {
        public override void Update(Frame frame)
        {
            var projectileFilter = frame.Unsafe.FilterStruct<ProjectileFilter>();
            var projectile = default(ProjectileFilter);
            while (projectileFilter.Next(&projectile))
            {
                if (frame.Physics3D.TryGetQueryHits(projectile.Component->PathQueryRef, out var hitsOnTrajectory) == false || hitsOnTrajectory.Count <= 0)
                {
                    projectile.Transform->Position =
                        projectile.Transform->Rotation *
                        projectile.Transform->Forward *
                        projectile.Component->Speed * frame.DeltaTime;
                    continue;
                }

                if (frame.Physics3D.TryGetQueryHits(projectile.Component->DamageZoneQueryRef, out var damageZoneHits))
                {
                    for (int i = 0; i < damageZoneHits.Count; i++)
                    {
                        // Apply damage logic
                    }
                }
            }
        }
    }
}

```

In addition to that, it is possible to grab all broadphase results via the `public bool GetAllQueriesHits(out HitCollection\* queriesHits, out int queriesCount)` call which is also available via `Frame.Physics`.

### Note

A few important points to keep in mind when using broadphase queries:

- The performance is around 20x better for large numbers (e.g. projectiles).
- They are based on the frame state before the Physics system kicks in.
- Broadphase queries do not carry over between frames; i.e. they need to be injected at the start of a frame before the Physics. A broadphase query injected after the Physics has run will never return a result. This is because Quantum's Physics are stateless.

## Emulating CCD

Quantum's physics engine is stateless. Continuous collision detection would be excessively expensive in such a system. The solutions to emulate the CCD behaviour in stateless physics engines normally involve raycasts or shape overlaps that extend to the expected movement for one frame.

This topic usually comes up in combination with fast moving entities such as projectiles. Depending on the size of the fast moving object, we recommend using one of the following approaches:

- a short ray in the movement direction with length being `velocity \* deltaTime`; or,
- a single overlap; note that a shape overlap can also be done with compound shapes.

Either of these solutions would replicate a 100% accurate CCD and result in a much better performance overall. To further improve performance, it is possible to combine this with broadphase queries.

Back to top

- [Introduction](#introduction)
- [Queries](#queries)

  - [Linecast and Raycast](#linecast-and-raycast)
  - [Shape Queries](#shape-queries)
  - [Sorting Hits](#sorting-hits)

- [Options](#options)

  - [Hit Normals](#hit-normals)
  - [Filtering Hits](#filtering-hits)

- [Broadphase Queries](#broadphase-queries)

  - [Injecting Queries](#injecting-queries)
  - [Retrieving Query Results](#retrieving-query-results)
  - [Note](#note)

- [Emulating CCD](#emulating-ccd)