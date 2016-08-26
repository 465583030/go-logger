package logger

import "io"

// Logger is the basic interface to a logger implementation.
type Logger interface {
	Verbosity() int64
	SetVerbosity(verbosity int64)
	Printf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fprintf(w io.Writer, format string, args ...interface{})
	Colorize(value string, color AnsiColorCode)
}
