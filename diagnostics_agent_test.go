package logger

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/blendlabs/go-assert"
)

func TestNewDiagnosticsEventQueue(t *testing.T) {
	assert := assert.New(t)

	eq := newDiagnosticsEventQueue()
	defer eq.Drain()

	assert.True(eq.IsDispatchSynchronous())
	assert.Zero(eq.Len())
	assert.Equal(DefaultDiagnosticsAgentQueueWorkers, eq.NumWorkers())
	assert.Equal(DefaultDiagnosticsAgentQueueLength, eq.MaxWorkItems())
	assert.False(eq.Running())
}

func TestNewDiagnosticsAgent(t *testing.T) {
	assert := assert.New(t)

	buffer := bytes.NewBuffer([]byte{})
	da := NewDiagnosticsAgent(EventAll, NewLogWriter(buffer))
	defer da.Close()

	assert.NotNil(da)
	assert.Equal(EventAll, da.Verbosity())
	assert.NotNil(da.eventListeners)
	assert.NotNil(da.eventQueue)
}

func TestNewDiagnosticsAgentFromEnvironment(t *testing.T) {
	assert := assert.New(t)

	oldLogVerbosity := os.Getenv("LOG_VERBOSITY")
	defer func() {
		os.Setenv("LOG_VERBOSITY", oldLogVerbosity)
	}()
	os.Setenv("LOG_VERBOSITY", "ALL")

	oldLogLabel := os.Getenv("LOG_LABEL")
	defer func() {
		os.Setenv("LOG_LABEL", oldLogLabel)
	}()
	os.Setenv("LOG_LABEL", "Testing Harness")

	da := NewDiagnosticsAgentFromEnvironment()
	defer da.Close()

	assert.Equal(EventAll, da.Verbosity())
	assert.True(da.Writer().UseAnsiColors())
	assert.True(da.Writer().ShowTimestamp())
	assert.True(da.Writer().ShowLabel())
	assert.False(da.EventQueue().Running())
	assert.True(da.EventQueue().IsDispatchSynchronous())
	assert.Equal("Testing Harness", da.Writer().Label())
}

func TestNewDiagnosticsAgentFromEnvironmentCustomVerbosity(t *testing.T) {
	assert := assert.New(t)

	oldLogVerbosity := os.Getenv("LOG_VERBOSITY")
	defer func() {
		os.Setenv("LOG_VERBOSITY", oldLogVerbosity)
	}()
	os.Setenv("LOG_VERBOSITY", "LOG_SHOW_ERROR,LOG_SHOW_INFO,LOG_SHOW_REQUEST")

	oldLogLabel := os.Getenv("LOG_LABEL")
	defer func() {
		os.Setenv("LOG_LABEL", oldLogLabel)
	}()
	os.Setenv("LOG_LABEL", "Testing Harness")

	da := NewDiagnosticsAgentFromEnvironment()
	defer da.Close()

	assert.True(da.CheckVerbosity(EventError))
	assert.True(da.CheckVerbosity(EventRequestComplete))
	assert.True(da.CheckVerbosity(EventInfo))
	assert.False(da.CheckVerbosity(EventWarning))
	assert.False(da.CheckVerbosity(EventFatalError))
	assert.True(da.Writer().UseAnsiColors())
	assert.True(da.Writer().ShowTimestamp())
	assert.True(da.Writer().ShowLabel())
	assert.False(da.EventQueue().Running())
	assert.True(da.EventQueue().IsDispatchSynchronous())
	assert.Equal("Testing Harness", da.Writer().Label())
}

func TestDiagnosticAgentVerbosity(t *testing.T) {
	assert := assert.New(t)

	buffer := bytes.NewBuffer([]byte{})
	da := NewDiagnosticsAgent(EventAll, NewLogWriter(buffer))
	defer da.Close()

	assert.Equal(EventAll, da.Verbosity())
	da.SetVerbosity(EventInfo)
	assert.Equal(EventInfo, da.Verbosity())
	assert.True(da.CheckVerbosity(EventInfo))
	assert.False(da.CheckVerbosity(EventRequest))
}

func TestDiagnosticsAgentAddEventListener(t *testing.T) {
	assert := assert.New(t)

	buffer := bytes.NewBuffer([]byte{})
	da := NewDiagnosticsAgent(EventAll, NewLogWriter(buffer))
	defer da.Close()

	assert.NotNil(da.eventListeners)
	da.AddEventListener(EventError, func(writer Logger, ts TimeSource, eventFlag uint64, state ...interface{}) {})
	assert.True(da.CheckVerbosity(EventError))
	assert.True(da.CheckHasHandler(EventError))
}

func TestDiagnosticsAgentOnEvent(t *testing.T) {
	assert := assert.New(t)

	buffer := bytes.NewBuffer([]byte{})
	da := NewDiagnosticsAgent(EventAll, NewLogWriter(buffer))
	defer da.Close()

	wg := sync.WaitGroup{}
	wg.Add(1)

	assert.NotNil(da.eventListeners)
	da.AddEventListener(EventError, func(writer Logger, ts TimeSource, eventFlag uint64, state ...interface{}) {
		defer wg.Done()
		assert.Equal(EventError, eventFlag)
		assert.NotEmpty(state)
		assert.Len(state, 2)
		assert.Equal("Hello", state[0])
		assert.Equal("World", state[1])
	})
	assert.True(da.CheckVerbosity(EventError))
	assert.True(da.CheckHasHandler(EventError))

	da.OnEvent(EventError, "Hello", "World")
	wg.Wait()
}

func TestDiagnosticsAgentOnEventMultipleListeners(t *testing.T) {
	assert := assert.New(t)

	buffer := bytes.NewBuffer([]byte{})
	da := NewDiagnosticsAgent(EventAll, NewLogWriter(buffer))
	defer da.Close()

	wg := sync.WaitGroup{}
	wg.Add(2)

	assert.NotNil(da.eventListeners)
	da.AddEventListener(EventError, func(writer Logger, ts TimeSource, eventFlag uint64, state ...interface{}) {
		defer wg.Done()
		assert.Equal(EventError, eventFlag)
		assert.NotEmpty(state)
		assert.Len(state, 2)
		assert.Equal("Hello", state[0])
		assert.Equal("World", state[1])
	})
	da.AddEventListener(EventError, func(writer Logger, ts TimeSource, eventFlag uint64, state ...interface{}) {
		defer wg.Done()
		assert.Equal(EventError, eventFlag)
		assert.NotEmpty(state)
		assert.Len(state, 2)
		assert.Equal("Hello", state[0])
		assert.Equal("World", state[1])
	})
	assert.True(da.CheckVerbosity(EventError))
	assert.True(da.CheckHasHandler(EventError))

	da.OnEvent(EventError, "Hello", "World")
	wg.Wait()
}

func TestDiagnosticsAgentOnEventUnhandled(t *testing.T) {
	assert := assert.New(t)

	buffer := bytes.NewBuffer([]byte{})
	da := NewDiagnosticsAgent(EventAll, NewLogWriter(buffer))
	defer da.Close()

	assert.NotNil(da.eventListeners)
	da.AddEventListener(EventError, func(writer Logger, ts TimeSource, eventFlag uint64, state ...interface{}) {
		assert.FailNow("The Error Handler shouldn't have fired")
	})
	assert.True(da.CheckVerbosity(EventError))
	assert.True(da.CheckVerbosity(EventFatalError))
	assert.True(da.CheckHasHandler(EventError))
	assert.False(da.CheckHasHandler(EventFatalError))

	da.OnEvent(EventFatalError, "Hello", "World")
}

func TestDiagnosticsAgentOnEventUnflagged(t *testing.T) {
	assert := assert.New(t)

	buffer := bytes.NewBuffer([]byte{})
	da := NewDiagnosticsAgent(EventFlagCombine(EventInfo, EventRequestComplete), NewLogWriter(buffer))
	defer da.Close()

	assert.NotNil(da.eventListeners)
	da.AddEventListener(EventError, func(writer Logger, ts TimeSource, eventFlag uint64, state ...interface{}) {
		assert.FailNow("The Error Handler shouldn't have fired")
	})
	assert.False(da.CheckVerbosity(EventError))
	assert.True(da.CheckHasHandler(EventError))

	da.OnEvent(EventError, "Hello", "World")
}

func TestDiagnosticsAgentEventf(t *testing.T) {
	assert := assert.New(t)

	buffer := bytes.NewBuffer([]byte{})
	da := NewDiagnosticsAgent(EventAll, NewLogWriter(buffer))
	defer da.Close()

	wg := sync.WaitGroup{}
	wg.Add(1)

	da.AddEventListener(EventInfo, func(writer Logger, ts TimeSource, eventFlag uint64, state ...interface{}) {
		defer wg.Done()
	})

	da.Eventf(EventInfo, "Informational", ColorWhite, "%s World", "Hello")
	wg.Wait()

	assert.NotZero(buffer.Len())
	assert.True(strings.HasSuffix(buffer.String(), "Hello World\n"), buffer.String())
}

func TestDiagnosticsAgentErrorf(t *testing.T) {
	assert := assert.New(t)

	stdout := bytes.NewBuffer([]byte{})
	stderr := bytes.NewBuffer([]byte{})
	da := NewDiagnosticsAgent(EventAll, NewLogWriter(stdout, stderr))
	defer da.Close()

	wg := sync.WaitGroup{}
	wg.Add(1)

	da.AddEventListener(EventError, func(writer Logger, ts TimeSource, eventFlag uint64, state ...interface{}) {
		defer wg.Done()
	})

	da.ErrorEventf(EventError, "Info Error", ColorWhite, "%s World", "Hello")
	wg.Wait()

	assert.Zero(stdout.Len())
	assert.NotZero(stderr.Len())
	assert.True(strings.HasSuffix(stderr.String(), "Hello World\n"), stderr.String())
}

func TestDiagnosticsAgentFireEvent(t *testing.T) {
	assert := assert.New(t)

	buffer := bytes.NewBuffer([]byte{})
	da := NewDiagnosticsAgent(EventAll, NewLogWriter(buffer))
	defer da.Close()
	da.writer.SetUseAnsiColors(false)

	ts := TimeInstance(time.Date(2016, 01, 02, 03, 04, 05, 06, time.UTC))
	da.AddEventListener(EventInfo, func(wr Logger, ts TimeSource, e uint64, state ...interface{}) {
		wr.WriteWithTimeSource(ts, []byte(fmt.Sprintf("Hello World")))
	})

	err := da.fireEvent(ts, EventInfo)
	assert.Nil(err)

	assert.True(strings.HasPrefix(buffer.String(), time.Time(ts).Format(DefaultTimeFormat)), buffer.String())
	assert.True(strings.HasSuffix(buffer.String(), "Hello World\n"))
}

func TestDiagnosticsAgentWriteEventMessageWithOutput(t *testing.T) {
	assert := assert.New(t)

	buffer := bytes.NewBuffer([]byte{})
	da := NewDiagnosticsAgent(EventAll, NewLogWriter(buffer))
	defer da.Close()

	da.writer.SetUseAnsiColors(false)

	ts := TimeInstance(time.Date(2016, 01, 02, 03, 04, 05, 06, time.UTC))
	err := da.writeEventMessageWithOutput(da.writer.PrintfWithTimeSource, ts, "test", ColorWhite, "%s World", "Hello")
	assert.Nil(err)
	assert.True(strings.HasPrefix(buffer.String(), time.Time(ts).Format(DefaultTimeFormat)))
	assert.True(strings.HasSuffix(buffer.String(), "Hello World\n"))
}
