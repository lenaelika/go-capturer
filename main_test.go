package capturer_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lenaelika/go-capturer"
)

func TestCapturer(t *testing.T) {
	// original stdout & stderr
	stdout := os.Stdout
	stderr := os.Stderr
	// suppressive writer
	devnull, err := os.Open(os.DevNull)
	require.Nil(t, err)

	// log to stderr w/o timestamp
	log.SetFlags(0)
	// log to stdout with custom prefix
	custom := log.New(stdout, "custom", 0)

	capture := capturer.New()
	check := func(expected string) {
		got, err := capture.Output(func() {
			fmt.Print("stdout")
			fmt.Fprint(os.Stderr, "stderr")
			log.Print("stdlog")
			custom.Print("log")
		})
		if assert.Nil(t, err) {
			assert.Equal(t, expected, got)
		}
	}

	// empty case
	os.Stdout = devnull
	os.Stderr = devnull
	log.SetOutput(devnull)
	custom.SetOutput(devnull)
	check("")

	// std log
	capture.Stdlog = true
	check("stdlog\n")

	// custom log
	capture = capturer.New(custom)
	check("customlog\n")

	// write to os.File
	os.Stdout = stdout
	capture = capturer.New(&os.Stdout)
	check("stdout")

	// multi write
	os.Stderr = stderr
	capture = capturer.New(&os.Stdout, &os.Stderr, custom)
	capture.Stdlog = true
	check("stdoutstderrstdlog\ncustomlog\n")

	// reset os.File
	assert.Equal(t, stdout, os.Stdout)
	assert.Equal(t, stderr, os.Stderr)

	// how to check std log reset?
	// custom log cannot be reset

	// error cases
	msg := "capturer: logger is not initialized"
	assert.PanicsWithValue(t, msg, func() {
		var logger *log.Logger
		capturer.New(logger)
	})
	msg = "capturer: unsupported writer type"
	assert.PanicsWithValue(t, msg, func() {
		capturer.New("unsupported")
	})
}
