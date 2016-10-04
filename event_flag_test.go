package logger

import (
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
	assert.True(set.IsEnabled("TEST"))
}
