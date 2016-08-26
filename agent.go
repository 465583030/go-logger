package logger

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	exception "github.com/blendlabs/go-exception"
)

const (
	// DefaultBufferSize is the default inner buffer size used in Fprintf.
	DefaultBufferSize = 1 << 8

	// EventNone is effectively logging disabled.
	EventNone = uint64(0)
	// EventAll represents every flag being enabled.
	EventAll = ^EventNone
	// EventError enables logging errors
	EventError = 1 << iota
	// EventDebug enables logging for debug messages.
	EventDebug = 1 << iota
	// EventInfo enables logging for informational messages.
	EventInfo = 1 << iota
)

var (
	_agent     *Agent
	_agentLock sync.Mutex
)

// Default returns the default agent.
func Default() *Agent {
	if _agent == nil {
		_agentLock.Lock()
		defer _agentLock.Unlock()
		if _agent == nil {
			_agent = NewAgent(EventNone, nil, nil) // do this so .Default() calls don't panic
		}
	}
	return _agent
}

// SetDefault sets the default agent.
func SetDefault(agent *Agent) {
	_agentLock.Lock()
	defer _agentLock.Unlock()
	_agent = agent
}

// NewAgent returns a new agent.
func NewAgent(verbosity uint64, output io.Writer, errorOutputs ...io.Writer) *Agent {
	agent := &Agent{
		Output:        NewSyncWriter(output),
		UseAnsiColors: true,
		ShowTimestamp: true,
		ShowLabel:     false,
		verbosity:     verbosity,
		verbosityLock: &sync.RWMutex{},
		bufferPool:    NewBufferPool(DefaultBufferSize),
	}
	if len(errorOutputs) > 0 {
		agent.ErrorOutput = NewSyncWriter(errorOutputs[0])
	}
	return agent
}

// Agent handles outputting logging events.
type Agent struct {
	ShowTimestamp bool
	ShowLabel     bool
	UseAnsiColors bool

	Output      io.Writer
	ErrorOutput io.Writer

	Label string

	verbosity     uint64
	verbosityLock *sync.RWMutex
	bufferPool    *BufferPool

	eventListeners map[uint64][]EventListener
	errorListeners []ErrorListener
}

// AddEventListener adds a listener for errors.
func (a *Agent) AddEventListener(eventFlag uint64, listener EventListener) {
	if a.eventListeners == nil {
		a.eventListeners = map[uint64][]EventListener{}
	}
	a.eventListeners[eventFlag] = append(a.eventListeners[eventFlag], listener)
}

// OnEvent runs the currently configured event listeners.
func (a *Agent) OnEvent(eventFlag uint64, state ...interface{}) {
	listeners := a.eventListeners[eventFlag]
	for x := 0; x < len(listeners); x++ {
		listener := listeners[x]
		go listener(eventFlag, state...)
	}
}

// AddErrorListener adds a listener for errors.
func (a *Agent) AddErrorListener(listener ErrorListener) {
	a.errorListeners = append(a.errorListeners, listener)
}

// OnError runs the currently configured error listeners.
func (a *Agent) OnError(err error) {
	for x := 0; x < len(a.errorListeners); x++ {
		listener := a.errorListeners[x]
		go listener(err)
	}
}

// Printf writes to the output stream.
func (a *Agent) Printf(format string, args ...interface{}) {
	a.Fprintf(a.Output, format, args...)
}

// Errorf writes to the error stream (if present).
func (a *Agent) Errorf(format string, args ...interface{}) {
	if a.ErrorOutput != nil {
		a.Fprintf(a.ErrorOutput, format, args...)
		return
	}
	a.Fprintf(a.Output, format, args...)
}

// Fprintf writes a given string and args to a writer.
func (a *Agent) Fprintf(w io.Writer, format string, args ...interface{}) {
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

	buf := a.bufferPool.Get()
	defer a.bufferPool.Put(buf)

	if a.ShowTimestamp {
		buf.WriteString(a.GetTimestamp())
		buf.WriteRune(RuneSpace)
	}

	if a.ShowLabel && len(a.Label) > 0 {
		buf.WriteString(a.GetLabel())
		buf.WriteRune(RuneSpace)
	}

	buf.WriteString(message)
	buf.WriteRune(RuneNewline)
	buf.WriteTo(w)
}

// SetVerbosity sets the agent verbosity synchronously.
func (a *Agent) SetVerbosity(verbosity uint64) {
	a.verbosityLock.Lock()
	defer a.verbosityLock.Unlock()

	a.verbosity = verbosity
}

// CheckVerbosity asserts if a flag value is set or not.
func (a *Agent) CheckVerbosity(flagValue uint64) bool {
	a.verbosityLock.RLock()
	defer a.verbosityLock.RUnlock()

	return BitFlagAny(a.verbosity, flagValue)
}

// Colorize (optionally) applies a color to a string.
func (a *Agent) Colorize(value string, color AnsiColorCode) string {
	if a.UseAnsiColors {
		return color.Apply(value)
	}
	return value
}

// ColorizeByStatus colorizes a string by a status code (green, yellow, red).
func (a *Agent) ColorizeByStatus(value string, statusCode int) string {
	if a.UseAnsiColors {
		if statusCode == http.StatusOK {
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
func (a *Agent) GetTimestamp() string {
	return a.Colorize(time.Now().UTC().Format(time.RFC3339), ColorGray)
}

// GetLabel returns the app name.
func (a *Agent) GetLabel() string {
	return a.Colorize(a.Label, ColorBlue)
}

// ErrorLogger returns a generic error logger for use in other services.
func (a *Agent) ErrorLogger() *log.Logger {
	if a.ErrorOutput != nil {
		return log.New(a.ErrorOutput, "", 0)
	}
	return log.New(a.Output, "", 0)
}

// Infof logs an informational message to the output stream.
func (a *Agent) Infof(format string, args ...interface{}) {
	if a.CheckVerbosity(EventInfo) && len(format) > 0 {
		defer a.OnEvent(EventInfo)
		a.Printf("%s %s", a.Colorize("Info", ColorLightWhite), fmt.Sprintf(format, args...))
	}
}

// Debugf logs a debug message to the output stream.
func (a *Agent) Debugf(format string, args ...interface{}) {
	if a.CheckVerbosity(EventDebug) && len(format) > 0 {
		defer a.OnEvent(EventInfo)
		a.Printf("%s %s", a.Colorize("Debug", ColorLightWhite), fmt.Sprintf(format, args...))
	}
}

// Writef checks an event flag and writes a message with a given label and color.
func (a *Agent) Writef(eventFlag uint64, label string, labelColor AnsiColorCode, format string, args ...interface{}) {
	if a.CheckVerbosity(eventFlag) && len(format) > 0 {
		defer a.OnEvent(eventFlag)
		a.Printf("%s %s", a.Colorize(label, labelColor), fmt.Sprintf(format, args...))
	}
}

// WriteBinary writes a message from bytes.
func (a *Agent) WriteBinary(eventFlag uint64, label string, labelColor AnsiColorCode, bytes []byte) {
	if a.CheckVerbosity(eventFlag) && len(bytes) > 0 {
		defer a.OnEvent(eventFlag)

		buf := a.bufferPool.Get()
		defer a.bufferPool.Put(buf)

		if a.ShowTimestamp {
			buf.WriteString(a.GetTimestamp())
			buf.WriteRune(RuneSpace)
		}

		if a.ShowLabel && len(a.Label) > 0 {
			buf.WriteString(a.GetLabel())
			buf.WriteRune(RuneSpace)
		}

		buf.WriteString(a.Colorize(label, labelColor))
		buf.WriteRune(RuneSpace)
		buf.Write(bytes)
		buf.WriteRune(RuneNewline)
		buf.WriteTo(a.Output)
	}
}

// WriteErrorf logs an error message with a given format and args.
// It also triggers event handlers for errors.
func (a *Agent) WriteErrorf(format string, args ...interface{}) {
	if a.CheckVerbosity(EventError) {
		defer a.OnEvent(EventError)
		message := fmt.Sprintf(format, args...)
		a.Errorf("%s %s", a.Colorize("Error", ColorRed), message)
	}
}

// WriteError logs an error to std out and the db.
func (a *Agent) WriteError(err error) {
	if err != nil {
		if a.CheckVerbosity(EventError) {
			defer a.OnError(err)
			if ex, isException := err.(*exception.Exception); isException {
				defer a.OnEvent(EventError)
				a.Errorf("%s %s", a.Colorize("Exception", ColorRed), ex.Message())
				a.Errorf("%s %s", a.Colorize("Stack Trace", ColorRed), ex.StackString())
				return
			}
			a.WriteErrorf(err.Error())
		}
	}
}

// WriteErrorIfNotNil writes and returns an error
func (a *Agent) WriteErrorIfNotNil(err error) error {
	if err != nil {
		a.WriteError(err)
	}
	return err
}
