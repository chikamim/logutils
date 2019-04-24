# logutils

logltsv is a Go package that augments the standard library "log" package
to make logging a bit more modern, without fragmenting the Go ecosystem
with new logging packages. Based on hashicorp's [logutils](https://github.com/hashicorp/logutils).

## Key features

- Good human readability, writability and computer parsability.
- Using standard package won't break your log code.
- Set log level and apply level filter.
- Structured log using [LTSV](http://ltsv.org) with plain text.
- Output [zap](https://github.com/uber-go/zap) compatible development and production format.

## Log format

```
[INFO] message\tkey1:value1\tkey2:value2
```

## The simplest thing that could possibly work

Presumably your application already uses the default `log` package. To switch, you'll want your code to look like the following:

```go
package main

import (
	"log"
	"os"

	"github.com/chikamim/logutils"
)

func main() {
	log.SetOutput(logutils.NewProduction())
	log.Print("[INFO] purchase items\tname:apple\tcount:3")
	// Output: {"level":"info","ts":1556071424.0936801,"caller":"main.go:6","msg":"purchase items","name":"apple","count":3}
}
```
