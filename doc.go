/*
Capture output to stdout/stderr, logs and other file writers (useful in tests).
Originally made for testing packages that use fmt/log to print data.

Example

	package somepkg

	import (
	"github.com/mailru/easyjson/tests"
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

The capturer resets stdout/stderr/stdlog after capturing.
But since log package does not provide a way to access the current writer without unsafe reflect,
it does not reset writers for custom loggers. If necessary, do it manually by AppLog.SetOutput().
*/
package capturer
