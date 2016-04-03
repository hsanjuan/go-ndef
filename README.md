ndef
====

| Master/stable | Unstable |
|:-------------:|:--------:|
| [![Build Status](https://travis-ci.org/hsanjuan/ndef.svg?branch=master)](https://travis-ci.org/hsanjuan/ndef) [![Coverage Status](https://coveralls.io/repos/github/hsanjuan/ndef/badge.svg?branch=master)](https://coveralls.io/github/hsanjuan/ndef?branch=master) | [![Build Status](https://travis-ci.org/hsanjuan/ndef.svg?branch=unstable)](https://travis-ci.org/hsanjuan/ndef) [![Coverage Status](https://coveralls.io/repos/github/hsanjuan/ndef/badge.svg?branch=unstable)](https://coveralls.io/github/hsanjuan/ndef?branch=unstable) |

A Go implementation of the NFC Data Exchange Format (NDEF).

`ndef` allows to easily work with NDEF Messages in Go, providing an easy way to parse bytes and generate bytes, ensuring the NDEF specification is followed.

Usage and documentation
-----------------------

```
$ go get github.com/hsanjuan/ndef
```


```go
import (
	"github.com/hsanjuan/ndef"
)
```

`ndef` uses godoc for documentation and examples. You can read it at https://godoc.org/github.com/hsanjuan/ndef .
