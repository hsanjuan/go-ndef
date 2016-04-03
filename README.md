ndef
====

[![Build Status](https://travis-ci.org/hsanjuan/ndef.svg?branch=master)](https://travis-ci.org/hsanjuan/ndef) [![Coverage Status](https://coveralls.io/repos/github/hsanjuan/ndef/badge.svg?branch=master)](https://coveralls.io/github/hsanjuan/ndef?branch=master)

A Go implementation of the NFC Data Exchange Format (NDEF).

`ndef` allows to easily work with NDEF Messages in Go, providing an easy way to parse bytes and generate bytes, ensuring the NDEF specification is followed.

Usage
-----

```go
import (
	"github.com/hsanjuan/ndef"
)
```

`ndef` provides a `ndef.Message` type which can be used to create or parse NDEF Messages. A NDEF Message is formed by NDEF Records. A `ndef.Record` type is also available, allowing to produce NDEF Messages by providing their NDEF Records directly (to produce a chunked NDEF message bytestream for example).

`ndef.Message` and `ndef.Record` implement methods `.ParseBytes(bytes []byte) error` and `.Bytes() ([]byte, error)`. `ParseBytes()` takes a byte slice and parses it into the Type struct fields. `Bytes()` does the opposite and returns a byte slice produced from the NDEF Message or Record. In both cases, if the bytes parsed or the given type is not following the NDEF specification, errors are returned.

Some examples are below:

### Parsing an NDEF Message from a byte slice

```go
ndef_message := &ndef.Message{} // Create uninitialized message
err := ndef_message.ParseBytes(some_bytes) // Parse bytes into it
if err != nil { // Your bytes don't look good
    return err
}
fmt.Println(ndef_message) // Print some info
```

### Creating an NDEF Message and converting it to a byte slice

```go
ndef_message := &ndef.Message{
    TNF: ndef.NFC_FORUM_WELL_KNOWN_TYPE,
    Type: []byte("T"),
    Payload: []byte("This is a message of T[ext] type"),
}
ndef_bytes, err := ndef_message.Bytes()
```

### Extracting a URL from an NDEF Message

```go
fmt.Printf("%s%s\n", ndef.URIProtocols(ndef_message.Payload[0]),
    string(ndef_message.Payload[1:]))
```
