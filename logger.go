package logger

import (
	"bytes"
	"io"
)

// Logger is the basic interface to a logger implementation.
type Logger interface {
	Printf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fprintf(w io.Writer, format string, args ...interface{})
	Write(data []byte) (int64, error)

	Colorize(value string, color AnsiColorCode) string
	ColorizeByStatusCode(statusCode int, value string) string

	GetBuffer() *bytes.Buffer
	PutBuffer(*bytes.Buffer)
}
