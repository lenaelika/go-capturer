# Go Capturer

It can capture output to `stdout`/`stderr`, logs and other `os.File` writers.
It is useful for writing tests that print data using `fmt`/`log` packages.

## Usage

```go
package somepkg

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/lenaelika/go-capturer"
)

// custom app logger (optional)
var AppLog *log.Logger

func init() {
	AppLog = log.New(os.Stdout, "[like]", 0)
	log.SetFlags(0) // w/o timestamp
}

func TestPrint(t *testing.T) {

	capture := capturer.New(&os.Stdout, &os.Stderr, AppLog)
	capture.Stdlog = true

	expected := "there's\nno~test\n[like]production\n"
	output, err := capture.Output(print)
	if err != nil {
		t.Errorf("expected nil, got error: %v", err)
	} else if expected != output {
		t.Errorf("expected: %v\ngot: %v\n", expected, output)
	}
}

func print() {
	// os.File output
	fmt.Println("there's")
	fmt.Fprint(os.Stderr, "no~")

	// standard logger
	log.Print("test")

	// custom logger
	AppLog.Print("production")
}
```

_The capturer resets stdout/stderr/stdlog after capturing.
But since log package does not provide a way to access the current writer without unsafe reflect,
it does not reset writers for custom loggers. If necessary, do it manually by `AppLog.SetOutput()`._

## Installation

```
$ go get github.com/lenaelika/go-capturer
```

Please feel free to submit issues and send pull requests. 🇷🇺🇬🇧
