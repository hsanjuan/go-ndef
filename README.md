Go-ndef
=======

| Master/stable | Unstable | Reference |
|:-------------:|:--------:|:---------:|
| [![Build Status](https://travis-ci.org/hsanjuan/go-ndef.svg?branch=master)](https://travis-ci.org/hsanjuan/go-ndef) [![Coverage Status](https://coveralls.io/repos/github/hsanjuan/go-ndef/badge.svg?branch=master)](https://coveralls.io/github/hsanjuan/go-ndef?branch=master) | [![Build Status](https://travis-ci.org/hsanjuan/go-ndef.svg?branch=unstable)](https://travis-ci.org/hsanjuan/go-ndef) [![Coverage Status](https://coveralls.io/repos/github/hsanjuan/go-ndef/badge.svg?branch=unstable)](https://coveralls.io/github/hsanjuan/go-ndef?branch=unstable) | [![GoDoc](https://godoc.org/github.com/hsanjuan/go-ndef?status.svg)](http://godoc.org/github.com/hsanjuan/go-ndef) |

A Go implementation of the NFC Data Exchange Format (NDEF).

`go-ndef` allows to easily work with NDEF Messages in Go, providing an easy way to parse bytes and generate bytes, ensuring the NDEF specification is followed.

Usage and documentation
-----------------------

```
$ go get github.com/hsanjuan/go-ndef
```


```go
import (
	"github.com/hsanjuan/go-ndef"
)
```

`go-ndef` uses godoc for documentation and examples. You can read it at https://godoc.org/github.com/hsanjuan/go-ndef .
