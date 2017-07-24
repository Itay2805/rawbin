# rawbin

rawbin is a library for encoding and decoding Go structs and valyes to raw binary data.

# Supported Types

  - Anything implementing `rawbin.Umarshaller` (decoding) and `rawbin.Marshaller` (encoding)
  - Anything implementing `encoding.BinaryMarshaler` (encoding)
  - Most integer types (`int8`, `int16`, `int32`, `int64`, `uint8`, `uint16`, `uint32`, `uint64`) (`int` and `uint` are not supported because they are not size-fixed)
  - `string` as length(`int16`) + data(`[]byte`)
  - floating point types (`float32`, `float64`)
  - `bool` (`true` - 1, `false` - 0)
  
 There are also some basic types which implements the `rawbin.Umarshaller` and `rawbin.Marshaller`

- `rawbin.Varint` and `rawbin.Varlong` 
- `rawbin.String` (string which uses `rawbin.Varlong` for the length instead of `int16`)

    
### Todos

 - Support `encoding.BinaryUnmarshaler`

License
----

MIT

**Free Software, Hell Yeah!**