package logger

import (
	"os"
	"testing"

	"github.com/blendlabs/go-assert"
)

func TestEventFlagSetEnable(t *testing.T) {
	assert := assert.New(t)

	set := NewEventFlagSet()
	set.Enable("TEST")
	assert.True(set.IsEnabled("TEST"))
	assert.False(set.IsEnabled("NOT_TEST"))
}

func TestEventFlagSetDisable(t *testing.T) {
	assert := assert.New(t)

	set := NewEventFlagSet()
	set.Enable("TEST")
	assert.True(set.IsEnabled("TEST"))
	set.Disable("TEST")
	assert.False(set.IsEnabled("TEST"))
}

func TestEventFlagSetEnableAll(t *testing.T) {
	assert := assert.New(t)

	set := NewEventFlagSet()
	set.EnableAll()
	assert.True(set.IsEnabled("TEST"))
	assert.True(set.IsEnabled("NOT_TEST"))
	assert.True(set.IsEnabled("NOT_TEST"))
	set.Disable("TEST")
	assert.True(set.IsEnabled("NOT_TEST"))
	assert.False(set.IsEnabled("TEST"))
}

func TestEventFlagSetFromEnvironment(t *testing.T) {
	assert := assert.New(t)

	oldLogVerbosity := os.Getenv(EnvironmentVariableLogEvents)
	defer func() {
		os.Setenv(EnvironmentVariableLogEvents, oldLogVerbosity)
	}()
	os.Setenv(EnvironmentVariableLogEvents, "error,info,web.request")

	set := NewEventFlagSetFromEnvironment()
	assert.True(set.IsEnabled(EventError))
	assert.True(set.IsEnabled(EventInfo))
	assert.True(set.IsEnabled(EventWebRequest))
	assert.False(set.IsEnabled(EventFatalError))
}

func TestEventFlagSetFromEnvironmentWithDisabled(t *testing.T) {
	assert := assert.New(t)

	oldLogVerbosity := os.Getenv(EnvironmentVariableLogEvents)
	defer func() {
		os.Setenv(EnvironmentVariableLogEvents, oldLogVerbosity)
	}()
	os.Setenv(EnvironmentVariableLogEvents, "all,-debug")

	set := NewEventFlagSetFromEnvironment()
	assert.True(set.IsEnabled(EventError))
	assert.False(set.IsEnabled(EventDebug))
}

func TestEventFlagSetFromEnvironmentAll(t *testing.T) {
	assert := assert.New(t)

	oldLogVerbosity := os.Getenv(EnvironmentVariableLogEvents)
	defer func() {
		os.Setenv(EnvironmentVariableLogEvents, oldLogVerbosity)
	}()
	os.Setenv(EnvironmentVariableLogEvents, "all")

	set := NewEventFlagSetFromEnvironment()
	assert.True(set.IsAllEnabled())
	assert.False(set.IsNoneEnabled())
	assert.True(set.IsEnabled(EventError))
}

func TestEventFlagSetFromEnvironmentNone(t *testing.T) {
	assert := assert.New(t)

	oldLogVerbosity := os.Getenv(EnvironmentVariableLogEvents)
	defer func() {
		os.Setenv(EnvironmentVariableLogEvents, oldLogVerbosity)
	}()
	os.Setenv(EnvironmentVariableLogEvents, "none")

	set := NewEventFlagSetFromEnvironment()
	assert.False(set.IsAllEnabled())
	assert.True(set.IsNoneEnabled())
	assert.False(set.IsEnabled(EventError))
}

func TestEventFlagNoneEnableEvents(t *testing.T) {
	assert := assert.New(t)

	flags := NewEventFlagSetNone()
	assert.False(flags.IsEnabled("test_flag"))

	flags.Enable("test_flag")
	assert.True(flags.IsEnabled("test_flag"))
}
