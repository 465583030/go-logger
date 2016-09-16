package logger

import (
	"fmt"
	"os"
	"sync"

	workQueue "github.com/blendlabs/go-workqueue"
)

var (
	// DefaultDiagnosticsAgentQueueWorkers is the number of consumers
	// for the diagnostics agent work queue.
	DefaultDiagnosticsAgentQueueWorkers = 1 //runtime.NumCPU()
)

var (
	_diagnosticsAgent     *DiagnosticsAgent
	_diagnosticsAgentLock sync.Mutex
)

// InitializeDiagnostics initializes the Diagnostics() agent with a given verbosity
// and optionally a targeted writer (only the first variadic writer will be used).
func InitializeDiagnostics(verbosity uint64, writers ...Logger) {
	_diagnosticsAgentLock.Lock()
	defer _diagnosticsAgentLock.Unlock()

	_diagnosticsAgent = NewDiagnosticsAgent(verbosity, writers...)
}

// InitializeDiagnosticsFromEnvironment initializes the Diagnostics() agent with a given verbosity
// and optionally a targeted writer (only the first variadic writer will be used).
func InitializeDiagnosticsFromEnvironment() error {
	_diagnosticsAgentLock.Lock()
	defer _diagnosticsAgentLock.Unlock()

	eventFlag, err := ParseEventFlagNameSet(os.Getenv("LOG_VERBOSITY"))
	if err != nil {
		return err
	}
	_diagnosticsAgent = NewDiagnosticsAgent(eventFlag, NewLogWriterFromEnvironment())
	return nil
}

// Diagnostics returnes a default DiagnosticsAgent singleton.
func Diagnostics() *DiagnosticsAgent {
	if _diagnosticsAgent == nil {
		_diagnosticsAgentLock.Lock()
		defer _diagnosticsAgentLock.Unlock()
		if _diagnosticsAgent == nil {
			_diagnosticsAgent = NewDiagnosticsAgent(EventNone) // do this so .Diagnostics() calls don't panic
		}
	}
	return _diagnosticsAgent
}

// NewDiagnosticsAgent returns a new diagnostics with a given bitflag verbosity.
func NewDiagnosticsAgent(verbosity uint64, writers ...Logger) *DiagnosticsAgent {
	diag := &DiagnosticsAgent{
		verbosity:  verbosity,
		eventQueue: workQueue.NewQueueWithBufferSize(4096),
	}
	diag.eventQueue.Start(DefaultDiagnosticsAgentQueueWorkers)
	if len(writers) > 0 {
		diag.writer = writers[0]
	} else {
		diag.writer = NewLogWriter(os.Stdout, os.Stderr)
	}
	return diag
}

// DiagnosticsAgent is a handler for various logging events with descendent handlers.
type DiagnosticsAgent struct {
	writer         Logger
	verbosity      uint64
	eventListeners map[uint64][]EventListener
	eventQueue     *workQueue.Queue
}

// AddEventListener adds a listener for errors.
func (da *DiagnosticsAgent) AddEventListener(eventFlag uint64, listener EventListener) {
	if da.eventListeners == nil {
		da.eventListeners = map[uint64][]EventListener{}
	}
	da.eventListeners[eventFlag] = append(da.eventListeners[eventFlag], listener)
}

// OnEvent fires the currently configured event listeners.
func (da *DiagnosticsAgent) OnEvent(eventFlag uint64, state ...interface{}) {
	da.eventQueue.Enqueue(da.fireEvent, append([]interface{}{eventFlag}, state...)...)
}

// OnEvent fires the currently configured event listeners.
func (da *DiagnosticsAgent) fireEvent(actionState ...interface{}) error {
	if len(actionState) < 1 {
		return nil
	}

	eventFlag, err := stateAsEventFlag(actionState[0])
	if err != nil {
		return err
	}

	listeners := da.eventListeners[eventFlag]
	for x := 0; x < len(listeners); x++ {
		listener := listeners[x]
		listener(da.writer, eventFlag, actionState[1:]...)
	}

	return nil
}

// SetVerbosity sets the agent verbosity synchronously.
func (da *DiagnosticsAgent) SetVerbosity(verbosity uint64) {
	da.verbosity = verbosity
}

// CheckVerbosity asserts if a flag value is set or not.
func (da *DiagnosticsAgent) CheckVerbosity(flagValue uint64) bool {
	return EventFlagAny(da.verbosity, flagValue)
}

// Eventf checks an event flag and writes a message with a given label and color.
func (da *DiagnosticsAgent) Eventf(eventFlag uint64, label string, labelColor AnsiColorCode, format string, args ...interface{}) {
	if da.CheckVerbosity(eventFlag) && len(format) > 0 {
		defer da.OnEvent(eventFlag)
		da.eventQueue.Enqueue(da.writeEventMessage, append([]interface{}{label, labelColor, format}, args...)...)
	}
}

func (da *DiagnosticsAgent) writeEventMessage(actionState ...interface{}) error {
	if len(actionState) < 3 {
		return nil
	}
	label, err := stateAsString(actionState[0])
	if err != nil {
		return err
	}
	labelColor, err := stateAsAnsiColorCode(actionState[1])
	if err != nil {
		return err
	}
	format, err := stateAsString(actionState[2])
	if err != nil {
		return err
	}
	da.writer.Printf("%s %s", da.writer.Colorize(label, labelColor), fmt.Sprintf(format, actionState[3:]...))
	return nil
}

// Infof logs an informational message to the output stream.
func (da *DiagnosticsAgent) Infof(format string, args ...interface{}) {
	da.Eventf(EventInfo, "Info", ColorWhite, format, args...)
}

// Debugf logs a debug message to the output stream.
func (da *DiagnosticsAgent) Debugf(format string, args ...interface{}) {
	da.Eventf(EventDebug, "Debug", ColorLightYellow, format, args...)
}

// DebugDump dumps an object and fires a debug event.
func (da *DiagnosticsAgent) DebugDump(object interface{}) {
	da.Eventf(EventDebug, "Debug Dump", ColorLightYellow, "%v", object)
}

// Warningf logs a debug message to the output stream.
func (da *DiagnosticsAgent) Warningf(format string, args ...interface{}) {
	da.Eventf(EventWarning, "Warning", ColorYellow, format, args...)
}

// Warning logs a warning error to std out.
func (da *DiagnosticsAgent) Warning(err error) error {
	if err != nil {
		if da.CheckVerbosity(EventWarning) {
			da.Warningf(err.Error())
		}
	}
	return err
}

// Errorf writes an event to the log and triggers event listeners.
func (da *DiagnosticsAgent) Errorf(format string, args ...interface{}) {
	da.Eventf(EventError, "Error", ColorRed, format, args...)
}

// Fatal logs an error to std out.
func (da *DiagnosticsAgent) Error(err error) error {
	if err != nil {
		if da.CheckVerbosity(EventError) {
			da.Errorf(err.Error())
		}
	}
	return err
}

// Fatalf writes an event to the log and triggers event listeners.
func (da *DiagnosticsAgent) Fatalf(format string, args ...interface{}) {
	da.Eventf(EventFatalError, "Fatal Error", ColorRed, format, args...)
}

// Fatal logs an error to std out.
func (da *DiagnosticsAgent) Fatal(err error) error {
	if err != nil {
		if da.CheckVerbosity(EventError) {
			da.Fatalf(err.Error())
		}
	}
	return err
}
