# logltsv

logltsv is a Go package that augments the standard library "log" package
to make logging a bit more modern, without fragmenting the Go ecosystem
with new logging packages. Based on hashicorp's [logutils](https://github.com/hashicorp/logutils).

## Key features
* Good human readablity, writability and computer parsablility.
* Using standard package won't break your log code.
* Set log level and apply level filter.
* Structured log using [LTSV](http://ltsv.org) with plain text.
* Can output to single line JSON.

## The simplest thing that could possibly work

Presumably your application already uses the default `log` package. To switch, you'll want your code to look like the following:

```go
package main

import (
	"log"
	"os"

	"github.com/chikamim/logltsv"
)

func main() {
	log.SetOutput(logltsv.NewJSONOutput())
	log.Print("level:DEBUG\tmessage:Debugging")
	// this will not print
	log.Print("level:WARN\tmessage:Warning")
	// Outout: {"time":"2009/01/23 01:23:23", "level":"WARN", "message":"Warning"}
	log.Print("level:ERROR\tmessage:Erring")
	// Outout: {"time":"2009/01/23 01:23:23", "level":"ERROR", "message":"Erring"}
}
```
