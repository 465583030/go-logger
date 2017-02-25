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

func TestNewEventQueue(t *testing.T) {
	assert := assert.New(t)

	eq := newEventQueue()
	defer eq.Close()
	assert.Zero(eq.Len())
	assert.Equal(DefaultAgentQueueWorkers, eq.NumWorkers())
	assert.Equal(DefaultAgentQueueLength, eq.MaxWorkItems())
}

func TestNewAgent(t *testing.T) {
	assert := assert.New(t)

	buffer := bytes.NewBuffer([]byte{})
	da := New(NewEventFlagSetAll(), NewLogWriter(buffer))
	defer da.Close()

	assert.NotNil(da)
	assert.NotNil(da.Events())
	assert.True(da.Events().IsAllEnabled())
	assert.NotNil(da.eventListeners)
	assert.NotNil(da.eventQueue)
}

func TestNewAgentFromEnvironment(t *testing.T) {
	assert := assert.New(t)

	oldLogVerbosity := os.Getenv(EnvironmentVariableLogEvents)
	defer func() {
		os.Setenv(EnvironmentVariableLogEvents, oldLogVerbosity)
	}()
	os.Setenv(EnvironmentVariableLogEvents, "all")

	oldLogLabel := os.Getenv(EnvironmentVariableLogLabel)
	defer func() {
		os.Setenv(EnvironmentVariableLogLabel, oldLogLabel)
	}()
	os.Setenv(EnvironmentVariableLogLabel, "Testing Harness")

	da := NewFromEnvironment()
	defer da.Close()

	assert.NotNil(da.Events())
	assert.True(da.Writer().UseAnsiColors())
	assert.True(da.Writer().ShowTimestamp())
	assert.True(da.Writer().ShowLabel())
	assert.Equal("Testing Harness", da.Writer().Label())
}

func TestNewAgentFromEnvironmentCustomVerbosity(t *testing.T) {
	assert := assert.New(t)

	oldLogVerbosity := os.Getenv(EnvironmentVariableLogEvents)
	defer func() {
		os.Setenv(EnvironmentVariableLogEvents, oldLogVerbosity)
	}()
	os.Setenv(EnvironmentVariableLogEvents, "error,info,web.request")

	oldLogLabel := os.Getenv(EnvironmentVariableLogLabel)
	defer func() {
		os.Setenv(EnvironmentVariableLogLabel, oldLogLabel)
	}()
	os.Setenv(EnvironmentVariableLogLabel, "Testing Harness")

	da := NewFromEnvironment()
	defer da.Close()

	assert.True(da.IsEnabled(EventError))
	assert.True(da.IsEnabled(EventWebRequest))
	assert.True(da.IsEnabled(EventInfo))
	assert.False(da.IsEnabled(EventWarning))
	assert.False(da.IsEnabled(EventFatalError))
	assert.True(da.Writer().UseAnsiColors())
	assert.True(da.Writer().ShowTimestamp())
	assert.True(da.Writer().ShowLabel())
	assert.Equal("Testing Harness", da.Writer().Label())
}

func TestAgentEnableDisableEvent(t *testing.T) {
	assert := assert.New(t)

	da := New(NewEventFlagSet())
	da.EnableEvent("TEST")
	assert.True(da.IsEnabled("TEST"))
	da.EnableEvent("FOO")
	assert.True(da.IsEnabled("FOO"))

	da.DisableEvent("TEST")
	assert.False(da.IsEnabled("TEST"))
	assert.True(da.IsEnabled("FOO"))
}

func TestAgentVerbosity(t *testing.T) {
	assert := assert.New(t)

	da := New(NewEventFlagSetAll())
	da.SetVerbosity(NewEventFlagSetWithEvents(EventInfo))
	assert.True(da.IsEnabled(EventInfo))
	assert.False(da.IsEnabled(EventWebRequest))
}

func TestAgentAddEventListener(t *testing.T) {
	assert := assert.New(t)

	da := New(NewEventFlagSetAll())

	assert.NotNil(da.eventListeners)
	da.AddEventListener(EventError, func(writer Logger, ts TimeSource, eventFlag EventFlag, state ...interface{}) {})
	assert.True(da.IsEnabled(EventError))
	assert.True(da.HasListener(EventError))
}

func TestAgentOnEvent(t *testing.T) {
	assert := assert.New(t)

	buffer := bytes.NewBuffer([]byte{})
	da := New(NewEventFlagSetAll(), NewLogWriter(buffer))
	defer da.Close()

	wg := sync.WaitGroup{}
	wg.Add(1)

	assert.NotNil(da.eventListeners)
	da.AddEventListener(EventError, func(writer Logger, ts TimeSource, eventFlag EventFlag, state ...interface{}) {
		defer wg.Done()
		assert.Equal(EventError, eventFlag)
		assert.NotEmpty(state)
		assert.Len(state, 2)
		assert.Equal("Hello", state[0])
		assert.Equal("World", state[1])
	})
	assert.True(da.IsEnabled(EventError))
	assert.True(da.HasListener(EventError))

	da.OnEvent(EventError, "Hello", "World")
	wg.Wait()
}

func TestAgentOnEventMultipleListeners(t *testing.T) {
	assert := assert.New(t)

	buffer := bytes.NewBuffer([]byte{})
	da := New(NewEventFlagSetAll(), NewLogWriter(buffer))
	defer da.Close()

	wg := sync.WaitGroup{}
	wg.Add(2)

	assert.NotNil(da.eventListeners)
	da.AddEventListener(EventError, func(writer Logger, ts TimeSource, eventFlag EventFlag, state ...interface{}) {
		defer wg.Done()
		assert.Equal(EventError, eventFlag)
		assert.NotEmpty(state)
		assert.Len(state, 2)
		assert.Equal("Hello", state[0])
		assert.Equal("World", state[1])
	})
	da.AddEventListener(EventError, func(writer Logger, ts TimeSource, eventFlag EventFlag, state ...interface{}) {
		defer wg.Done()
		assert.Equal(EventError, eventFlag)
		assert.NotEmpty(state)
		assert.Len(state, 2)
		assert.Equal("Hello", state[0])
		assert.Equal("World", state[1])
	})
	assert.True(da.IsEnabled(EventError))
	assert.True(da.HasListener(EventError))

	da.OnEvent(EventError, "Hello", "World")
	wg.Wait()
}

func TestAgentOnEventUnhandled(t *testing.T) {
	assert := assert.New(t)

	buffer := bytes.NewBuffer([]byte{})
	da := New(NewEventFlagSetAll(), NewLogWriter(buffer))
	defer da.Close()

	assert.NotNil(da.eventListeners)
	da.AddEventListener(EventError, func(writer Logger, ts TimeSource, eventFlag EventFlag, state ...interface{}) {
		assert.FailNow("The Error Handler shouldn't have fired")
	})
	assert.True(da.IsEnabled(EventError))
	assert.True(da.IsEnabled(EventFatalError))
	assert.True(da.HasListener(EventError))
	assert.False(da.HasListener(EventFatalError))

	da.OnEvent(EventFatalError, "Hello", "World")
}

func TestAgentOnEventUnflagged(t *testing.T) {
	assert := assert.New(t)

	buffer := bytes.NewBuffer([]byte{})
	da := New(NewEventFlagSetWithEvents(EventInfo, EventWebRequest), NewLogWriter(buffer))
	defer da.Close()

	assert.NotNil(da.eventListeners)
	da.AddEventListener(EventError, func(writer Logger, ts TimeSource, eventFlag EventFlag, state ...interface{}) {
		assert.FailNow("The Error Handler shouldn't have fired")
	})
	assert.False(da.IsEnabled(EventError))
	assert.True(da.HasListener(EventError))

	da.OnEvent(EventError, "Hello", "World")
}

func TestAgentEventf(t *testing.T) {
	assert := assert.New(t)

	buffer := bytes.NewBuffer([]byte{})
	da := New(NewEventFlagSetAll(), NewLogWriter(buffer))
	defer da.Close()

	wg := sync.WaitGroup{}
	wg.Add(1)

	da.AddEventListener(EventInfo, func(writer Logger, ts TimeSource, eventFlag EventFlag, state ...interface{}) {
		defer wg.Done()
	})

	da.Eventf(EventInfo, ColorWhite, "%s World", "Hello")
	wg.Wait()

	assert.NotZero(buffer.Len())
	assert.True(strings.HasSuffix(buffer.String(), "Hello World\n"), buffer.String())
}

func TestAgentErrorf(t *testing.T) {
	assert := assert.New(t)

	stdout := bytes.NewBuffer([]byte{})
	stderr := bytes.NewBuffer([]byte{})
	da := New(NewEventFlagSetAll(), NewLogWriter(stdout, stderr))
	defer da.Close()

	wg := sync.WaitGroup{}
	wg.Add(1)

	da.AddEventListener(EventError, func(writer Logger, ts TimeSource, eventFlag EventFlag, state ...interface{}) {
		defer wg.Done()
	})

	da.Errorf("%s World", "Hello")
	wg.Wait()

	assert.Zero(stdout.Len())
	assert.NotZero(stderr.Len())
	assert.True(strings.HasSuffix(stderr.String(), "Hello World\n"), stderr.String())
}

func TestAgentFireEvent(t *testing.T) {
	assert := assert.New(t)

	buffer := bytes.NewBuffer([]byte{})
	da := New(NewEventFlagSetAll(), NewLogWriter(buffer))
	defer da.Close()
	da.writer.SetUseAnsiColors(false)

	ts := TimeInstance(time.Date(2016, 01, 02, 03, 04, 05, 06, time.UTC))
	da.AddEventListener(EventInfo, func(wr Logger, ts TimeSource, e EventFlag, state ...interface{}) {
		wr.WriteWithTimeSource(ts, []byte(fmt.Sprintf("Hello World")))
	})

	err := da.fireEvent(ts, EventInfo)
	assert.Nil(err)

	assert.True(strings.HasPrefix(buffer.String(), time.Time(ts).Format(DefaultTimeFormat)), buffer.String())
	assert.True(strings.HasSuffix(buffer.String(), "Hello World\n"))
}

func TestAgentWriteEventMessageWithOutput(t *testing.T) {
	assert := assert.New(t)

	buffer := bytes.NewBuffer([]byte{})
	da := New(NewEventFlagSetAll(), NewLogWriter(buffer))
	defer da.Close()

	da.writer.SetUseAnsiColors(false)

	ts := TimeInstance(time.Date(2016, 01, 02, 03, 04, 05, 06, time.UTC))
	err := da.writeEventMessageWithOutput(da.writer.PrintfWithTimeSource, ts, EventFlag("test"), ColorWhite, "%s World", "Hello")
	assert.Nil(err)
	assert.True(strings.HasPrefix(buffer.String(), time.Time(ts).Format(DefaultTimeFormat)))
	assert.True(strings.HasSuffix(buffer.String(), "Hello World\n"))
}

func TestAgentRemoveListeners(t *testing.T) {
	assert := assert.New(t)

	da := New(NewEventFlagSetAll())
	da.AddEventListener(EventError, func(writer Logger, ts TimeSource, eventFlag EventFlag, state ...interface{}) {})
	da.AddEventListener(EventInfo, func(writer Logger, ts TimeSource, eventFlag EventFlag, state ...interface{}) {})
	da.RemoveListeners(EventError)

	assert.False(da.HasListener(EventError))

}

func BenchmarkAgentIsEnabled(b *testing.B) {
	for iter := 0; iter < b.N; iter++ {
		for subIter := 0; subIter < 50; subIter++ {
			da := New(NewEventFlagSetWithEvents(EventFatalError, EventError, EventWebRequest, EventInfo))
			da.IsEnabled(EventFatalError)
			da.IsEnabled(EventWebUserError)
			da.IsEnabled(EventDebug)
			da.IsEnabled(EventInfo)
			da.IsEnabled(EventWebRequest)
		}
	}
}
