# logutils

logutils is a Go package that augments the standard library "log" package
to make logging a bit more modern, without fragmenting the Go ecosystem
with new logging packages.

## The simplest thing that could possibly work

Presumably your application already uses the default `log` package. To switch, you'll want your code to look like the following:

## Changed in this fork

- [Labeled Tab-separated Values](http://ltsv.org) log level

```go
package main

import (
	"log"
	"os"

	"github.com/chikamim/logltsv"
)

func main() {
	output := &logltsv.Output{
		Levels:     []logltsv.LogLevel{"DEBUG", "INFO", "WARN", "ERROR"},
		MinLevel:   logltsv.LogLevel("WARN"),
		Writer:     os.Stderr,
		JSONOutput: true,
	}

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.LUTC)
	log.SetOutput(output)
	log.Print("level:DEBUG\tmessage:Debugging")
	log.Print("level:WARN\tmessage:Warning")
	log.Print("level:ERROR\tmessage:Erring")
}
```

This logs to standard error exactly like go's standard logger. Any log messages you haven't converted to have a level will continue to print as before.
