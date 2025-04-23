# serialization-in-photon

_Source: https://doc.photonengine.com/realtime/current/reference/serialization-in-photon_

# Serialization in Photon

This page is not relevant for the JavaScript SDKs since they use JSON serialization and not the binary protocol.
This makes the JavaScript SDKs have limited data types, at least compared to the other SDKs.
Consider this especially if you want clients built with the JavaScript SDK to communicate with clients built with other SDKs.

Photon and its clients are using a highly optimized binary protocol to communicate.

It's compact, yet easy to parse.

Photon must convert all data into this binary protocol before it can be sent.

This is done automatically for a range of commonly used data types.

On most client APIs, you can also register your custom serialization methods for other classes you might need.

See below.

## Photon Supported Types

Each Photon supported type needs some reserved bytes for the type\_info.

- Primitive types need 1 byte for type\_info.
- Most collections need 1 extra byte for the type of the collection.


This is not the case for Hashtable, Object array and Byte array as the collection type is part of the type\_info.
- Strongly typed collections send the element type once instead of sending type\_info per element.
- All collections need 2 bytes to store their length. This is because the type of the length is `short`.


Byte array is the exception to this rule. Its length is of type `int`


and it needs 4 bytes to store the length.
- A string\[\] can't contain null values. Use empty instead.

The following types are generally supported and known by Photon's binary protocol.

As some languages don't provide every listed type, some SDKs support less types.

| Type (C#) | Size \[bytes\] (photon\_sizeof) | Description |
| --- | --- | --- |
| byte | 2 | 8 bit unsigned <br> 2 = type\_info(byte) + sizeof(byte) |
| bool (boolean) | 2 | true or false <br> 2 = type\_info(bool) + sizeof(bool) |
| short | 3 | 16 bit <br> 3 = type\_info(short) + sizeof(short) |
| int (integer) | 5 | 32 bit <br> 5 = type\_info(int) + sizeof(int) |
| long | 9 | 64 bit <br> 9 = type\_info(long) + sizeof(long) |
| float | 5 | 32 bit <br> 5 = type\_info(float) + sizeof(float) |
| double | 9 | 64 bit <br> 9 = type\_info(double) + sizeof(double) |
| String | 3 + sizeof( UTF8.GetBytes(string\_value) ) | length ≤ short.MaxValue <br> 3 = type\_info(String) + length\_size; <br> length\_size = sizeof(short) |
| Object\[\] (Object-array) | 3 + photon\_sizeof(elements) | length ≤ short.MaxValue <br> 3 = type\_info(Object\[\]) + length\_size; <br> length\_size = sizeof(short) |
| byte\[\] (byte-array) | 5 + length | length ≤ int.MaxValue <br> 5 = type\_info(byte\[\]) + length\_size; <br> length\_size = sizeof(int) |
| array (array of type T, T\[\]) | 4 + photon\_sizeof(elements) - length \* type\_info(T) | length ≤ short.MaxValue <br> T-type can be any of the types listed in this table except byte.<br> <br> 4 = type\_info(array) + type\_info(T) + length\_size; <br> length\_size = sizeof(short) |
| Hashtable | 3 + photon\_sizeof(keys) + photon\_sizeof(values) | pairs count ≤ short.MaxValue <br> 3 = type\_info(Hashtable) + length\_size; <br> length\_size = sizeof(short) |
| Dictionary<Object,Object> | 5 + photon\_sizeof(keys) + photon\_sizeof(values) | pairs count ≤ short.MaxValue <br> 5 = type\_info(Dictionary) + 2 \* type\_info(Object) + length\_size; <br> length\_size = sizeof(short)<br> <br> Dictionary keys should not be of type Dictionary. |
| Dictionary<Object,V> | 5 + photon\_sizeof(keys) + photon\_sizeof(values) - count(keys) \* type\_info(V) | pairs count ≤ short.MaxValue <br> V-type can be any of the types listed in this table.<br> <br> 5 = type\_info(Dictionary) + type\_info(Object) + type\_info(V) + length\_size; <br> length\_size = sizeof(short)<br> <br> Dictionary keys should not be of type Dictionary. |
| Dictionary<K,Object> | 5 + photon\_sizeof(keys) + photon\_sizeof(values) - count(keys) \* type\_info(K) | pairs count ≤ short.MaxValue <br> K-type can be any of the types listed in this table.<br> <br> 5 = type\_info(Dictionary) + type\_info(K) + type\_info(Object) + length\_size; <br> length\_size = sizeof(short)<br> <br> Dictionary keys should not be of type Dictionary. |
| Dictionary<K,V> | 5 + photon\_sizeof(keys) + photon\_sizeof(values) - count(keys) \* (type\_info(K) + type\_info(V)) | pairs count ≤ short.MaxValue <br> K- and V-types can be any of the types listed in this table.<br> <br> 5 = type\_info(Dictionary) + type\_info(K) + type\_info(V) + length\_size; length\_size = sizeof(short)<br> <br> Dictionary keys should not be of type Dictionary. |

## Custom Types

For any type not listed above, Photon will need your help to de/serialize important values.

The basic idea is that you write two methods to convert your class to a byte-array and back, then register those with the Photon API.

When that's done you can include instances of that type in any message you send.

Custom types have 2 bytes for type\_info:

one byte to tell that it's a custom type plus one more for the custom type code.

Photon supports up to 256 custom types. We recommend choosing custom type codes from 255 and downward.

Photon will call the serialization methods for a registered type and automatically prefix the created byte-arrays with 4 bytes:

2 bytes for the necessary type information and 2 bytes for payload length.

Due to the 4 bytes overhead, you might want to avoid registering types that are just a few bytes of data.

The **Photon Server is able to forward unknown custom types** "as is".

This is why you don't need to register your types in the Photon Cloud.

Make sure to register your custom types on all communicating clients.

Register custom types on server side or on plugins when needed.

`RegisterType` method returns a boolean result, which tells you if the type could be registered.

If any error occurs during custom type registration the method will return `false` and nothing will be changed.

Otherwise registration should be successful and the returned value is `true`.

If the custom code is already used the registration will fail and the method will return false.

Overriding registered serialization and deserialization methods for a same custom type will fail and the old ones will still be used.

## Custom Types in C#

All of our C#-based APIs (.NET, Unity, Xamarin, etc) provide the same way to register classes.

There are two ways of doing this which depends on whether or not you use our custom `Steam` class or just a "plain" byte array.

### Byte Array Method

The static method to call is:

C#

```csharp
    PhotonPeer.RegisterType(Type customType, byte code, SerializeMethod serializeMethod, DeserializeMethod deserializeMethod)

```

`SerializeMethod` and `DeserializeMethod` are defined delegates with the following respective signatures:

C#

```csharp
    public delegate byte[] SerializeMethod(object customObject);
    public delegate object DeserializeMethod(byte[] serializedCustomObject);

```

#### Example

As example, we implemented a simple basic `MyCustomType`:

C#

```csharp
public class MyCustomType
{
  public byte Id { get; set; }
  public static object Deserialize(byte[] data)
  {
    var result = new MyCustomType();
    result.Id = data[0];
    return result;
  }
  public static byte[] Serialize(object customType)
  {
    var c = (MyCustomType)customType;
    return new byte[] { c.Id };
  }
}

```

And to register it:

C#

```csharp
    PhotonPeer.RegisterType(typeof(MyCustomType), myCustomTypeCode, MyCustomType.Serialize, MyCustomType.Deserialize);

```

### StreamBuffer Method

`StreamBuffer` is our custom implementation of the `Stream` class.

It gives you all the benfits of a byte array wrapper and supports Photon's built-in serializable types.

The static method to call is:

C#

```csharp
    RegisterType(Type customType, byte code, SerializeStreamMethod serializeMethod, DeserializeStreamMethod deserializeMethod)

```

`SerializeStreamMethod` and `DeserializeStreamMethod` are defined delegates with the following respective signatures:

C#

```csharp
    public delegate short SerializeStreamMethod(StreamBuffer outStream, object customobject);
    public delegate object DeserializeStreamMethod(StreamBuffer inStream, short length);

```

#### Example

Let's have a look at how PUN implements support for Unity's `Vector2`.

A `Vector2` has 2 floats: `Vector2.x` and `Vector2.y`.

While floats are supported by Photon, `Vector2` is not.

Basically, using any C#-way to get the 4 bytes that resemble the float is ok.

Photon's Protocol class has a few Serialize methods that effectively write several types to byte-arrays.

C#

```csharp
    public static readonly byte[] memVector2 = new byte[2 * 4];
    private static short SerializeVector2(StreamBuffer outStream, object customobject)
    {
        Vector2 vo = (Vector2)customobject;
        lock (memVector2)
        {
            byte[] bytes = memVector2;
            int index = 0;
            Protocol.Serialize(vo.x, bytes, ref index);
            Protocol.Serialize(vo.y, bytes, ref index);
            outStream.Write(bytes, 0, 2 * 4);
        }
        return 2 * 4;
    }

```

Note that SerializeVector2 gets an object and has to cast this to the expected Vector2 type first.

Conversely, we also return just an object in DeserializeVector2:

C#

```csharp
    private static object DeserializeVector2(StreamBuffer inStream, short length)
    {
        Vector2 vo = new Vector2();
        lock (memVector2)
        {
            inStream.Read(memVector2, 0, 2 * 4);
            int index = 0;
            Protocol.Deserialize(out vo.x, memVector2, ref index);
            Protocol.Deserialize(out vo.y, memVector2, ref index);
        }
        return vo;
    }

```

And finally, we have to register the Vector2:

C#

```csharp
    PhotonPeer.RegisterType(typeof(Vector2), (byte)'W', SerializeVector2, DeserializeVector2);

```

### Additional information

When doing custom serializion of complex data, it's important to understand how data is transported and how to monitor traffic to keep it under control.

- Learn more about the [binary protocols](/realtime/current/reference/binary-protocol)
- Learn more about checking Network Traffic using [WireShark](/realtime/current/troubleshooting/analyzing-disconnects#wireshark)

Back to top

- [Photon Supported Types](#photon-supported-types)
- [Custom Types](#custom-types)
- [Custom Types in C#](#custom-types-in-c)
  - [Byte Array Method](#byte-array-method)
  - [StreamBuffer Method](#streambuffer-method)
  - [Additional information](#additional-information)