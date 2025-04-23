# dynamic-map

_Source: https://doc.photonengine.com/quantum/current/manual/maps/dynamic-map_

# Dynamic Map

Quantum 3 comes with a new map type called ```
Dynamic Map
```

 in order to make the customization process of maps at runtime easier.

## Features

- Provides an API for adding and removing static colliders at runtime.
  - Example:
  - C#
    ```csharp
    // add
    var meshColliderIndex = dynamicMap.AddMeshCollider(frame, meshCollider);

    // remove
    dynamicMap.RemoveMeshCollider(frame, meshColliderIndex);

    ```
- Clone existing maps and modify them safely.
  - Example:
  - C#
    ```csharp
    var dynamicMap = DynamicMap.FromStaticMap(f.Map);

    ```
- Modify static collider triangles at runtime.
  - Example:
  - C#
    ```csharp
    using (var scope = dynamicMap.EditMeshTriangles(frame, colliderIndex)) {
     TriangleCCW triangle = new TriangleCCW { A = v0, B = v1, C = v2 };

     if(scope.TryAddTriangle(ref triangle)) {
     Log.Info("Triangle added successfully.");
     }
     else {
     Log.Error("Failed to add triangle.");
     }
    }

    ```
- Optimizes Triangle data for late joiners.
  - Previously, when creating a dynamic map, the triangle data had to be serialized into a separate binary data asset and linked to the map manually.
  - With the new dynamic map system, this is done automatically and does not create any additional assets.
- Overridable methods for customizing the triangle serialization process.
  - It is now possible to override how the triangles are serialized and deserialized in order to allow further optimization of the data depending on the game context.

Back to top

- [Features](#features)