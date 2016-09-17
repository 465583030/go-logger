package logger

import (
	"bytes"
	"io"
)

// Logger is the basic interface to a logger implementation.
type Logger interface {
	Printf(format string, args ...interface{})
	PrintfWithTimingSource(timingSource TimingSource, format string, args ...interface{})
	Errorf(format string, args ...interface{})
	ErrorfWithTimingSource(timingSource TimingSource, format string, args ...interface{})
	Fprintf(w io.Writer, format string, args ...interface{})
	FprintfWithTimingSource(timingSource TimingSource, w io.Writer, format string, args ...interface{})
	Write(data []byte) (int64, error)
	WriteWithTimingSource(timingSource TimingSource, data []byte) (int64, error)

	Colorize(value string, color AnsiColorCode) string
	ColorizeByStatusCode(statusCode int, value string) string

	GetBuffer() *bytes.Buffer
	PutBuffer(*bytes.Buffer)
}
