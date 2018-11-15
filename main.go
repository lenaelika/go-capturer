package capturer

import (
	"bytes"
	"io"
	"log"
	"os"
)

type Instance struct {
	files   []*file
	loggers []*log.Logger
	Stdlog  bool
}

type file struct {
	ptr  **os.File
	orig *os.File
}

// New creates and configures a new capturer.
// It expects **os.File or *log.Logger, for instance &os.Stdout.
// To catch std log from log package, turn on Stdlog flag on the returned capturer.
func New(writers ...interface{}) *Instance {
	files := make([]*file, 0, len(writers))
	loggers := make([]*log.Logger, 0, len(writers))
	for _, writer := range writers {
		switch writer.(type) {
		case **os.File:
			ptr := writer.(**os.File)
			files = append(files, &file{ptr: ptr, orig: *ptr})
		case *log.Logger:
			logger := writer.(*log.Logger)
			if logger == nil {
				panic("capturer: logger is not initialized")
			}
			loggers = append(loggers, logger)
		default:
			panic("capturer: unsupported writer type")
		}
	}
	return &Instance{files: files, loggers: loggers}
}

// On sets a new receiver for the writers that the capturer should catch.
func (self *Instance) on(w *os.File) {
	for _, writer := range self.files {
		*writer.ptr = w
	}
	for _, logger := range self.loggers {
		logger.SetOutput(w)
	}
	if self.Stdlog {
		log.SetOutput(w)
	}
}

// Off resets file writers and std log.
func (self *Instance) off() {
	for _, writer := range self.files {
		*writer.ptr = writer.orig
	}
	if self.Stdlog {
		log.SetOutput(os.Stderr)
	}
}

// Output captures function's writes to all writers, specified to the capturer.
// It resets file writers and standard log to their original states after capturing,
// but does not reset custom loggers (since loggers' "out" field is unexported).
func (self *Instance) Output(f func()) (output string, err error) {
	r, w, err := os.Pipe()
	if err != nil {
		return
	}

	self.on(w)
	defer self.off()

	f()
	w.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, r)
	if err == nil {
		output = buf.String()
	}
	return
}
