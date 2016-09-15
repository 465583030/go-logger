package logger

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/blendlabs/go-workqueue"
)

const (
	// DefaultBufferSize is the default inner buffer size used in Fprintf.
	DefaultBufferSize = 1 << 8
)

var (
	_writer     *LogWriter
	_writerLock sync.Mutex
)

// NewLogWriter returns a new writer.
func NewLogWriter(output io.Writer, errorOutputs ...io.Writer) *LogWriter {
	agent := &LogWriter{
		Output:        NewSyncWriter(output),
		UseAnsiColors: true,
		ShowTimestamp: true,
		ShowLabel:     false,
		bufferPool:    NewBufferPool(DefaultBufferSize),
	}
	if len(errorOutputs) > 0 {
		agent.ErrorOutput = errorOutputs[0]
	}
	return agent
}

// LogWriter handles outputting logging events to given writer streams.
type LogWriter struct {
	ShowTimestamp bool
	ShowLabel     bool
	UseAnsiColors bool

	Output      io.Writer
	ErrorOutput io.Writer

	Label string

	writeQueue *workQueue.Queue
	bufferPool *BufferPool
}

// GetErrorOutput returns an io.Writer for the error stream.
func (wr *LogWriter) GetErrorOutput() io.Writer {
	if wr.ErrorOutput != nil {
		return wr.ErrorOutput
	}
	return wr.Output
}

// Colorize (optionally) applies a color to a string.
func (wr *LogWriter) Colorize(value string, color AnsiColorCode) string {
	if wr.UseAnsiColors {
		return color.Apply(value)
	}
	return value
}

// ColorizeByStatusCode colorizes a string by a status code (green, yellow, red).
func (wr *LogWriter) ColorizeByStatusCode(statusCode int, value string) string {
	if wr.UseAnsiColors {
		if statusCode >= http.StatusOK && statusCode < 300 { //the http 2xx range is ok
			return ColorGreen.Apply(value)
		} else if statusCode == http.StatusInternalServerError {
			return ColorRed.Apply(value)
		} else {
			return ColorYellow.Apply(value)
		}
	}
	return value
}

// GetTimestamp returns a new timestamp string.
func (wr *LogWriter) GetTimestamp() string {
	return wr.Colorize(time.Now().UTC().Format(time.RFC3339), ColorGray)
}

// GetLabel returns the app name.
func (wr *LogWriter) GetLabel() string {
	return wr.Colorize(wr.Label, ColorBlue)
}

// Printf writes to the output stream.
func (wr *LogWriter) Printf(format string, args ...interface{}) {
	wr.Fprintf(wr.Output, format, args...)
}

// Errorf writes to the error output stream.
func (wr *LogWriter) Errorf(format string, args ...interface{}) {
	wr.Fprintf(wr.GetErrorOutput(), format, args...)
}

// Write writes a binary blob to a given writer.
func (wr *LogWriter) Write(binary []byte) (int64, error) {
	buf := wr.bufferPool.Get()
	defer wr.bufferPool.Put(buf)

	if wr.ShowTimestamp {
		buf.WriteString(wr.GetTimestamp())
		buf.WriteRune(RuneSpace)
	}

	if wr.ShowLabel && len(wr.Label) > 0 {
		buf.WriteString(wr.GetLabel())
		buf.WriteRune(RuneSpace)
	}

	buf.Write(binary)
	buf.WriteRune(RuneNewline)
	return buf.WriteTo(wr.Output)
}

// Fprintf writes a given string and args to a writer.
func (wr *LogWriter) Fprintf(w io.Writer, format string, args ...interface{}) {
	if w == nil {
		return
	}
	if len(format) == 0 {
		return
	}
	message := fmt.Sprintf(format, args...)
	if len(message) == 0 {
		return
	}

	buf := wr.bufferPool.Get()
	defer wr.bufferPool.Put(buf)

	if wr.ShowTimestamp {
		buf.WriteString(wr.GetTimestamp())
		buf.WriteRune(RuneSpace)
	}

	if wr.ShowLabel && len(wr.Label) > 0 {
		buf.WriteString(wr.GetLabel())
		buf.WriteRune(RuneSpace)
	}

	buf.WriteString(message)
	buf.WriteRune(RuneNewline)
	buf.WriteTo(w)
}

// GetBuffer returns a leased buffer from the buffer pool.
func (wr *LogWriter) GetBuffer() *bytes.Buffer {
	return wr.bufferPool.Get()
}

// PutBuffer adds the leased buffer back to the pool.
// It Should be called in conjunction with `GetBuffer`.
func (wr *LogWriter) PutBuffer(buffer *bytes.Buffer) {
	wr.bufferPool.Put(buffer)
}
