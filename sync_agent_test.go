package logger

import (
	"bytes"
	"testing"

	assert "github.com/blendlabs/go-assert"
)

func TestSyncAgentInfof(t *testing.T) {
	assert := assert.New(t)

	buffer := bytes.NewBuffer(nil)
	var format string
	a := All(NewWriter(buffer))
	a.Writer().SetShowTimestamp(false)
	a.Writer().SetShowLabel(false)
	a.Writer().SetUseAnsiColors(false)
	a.AddEventListener(EventInfo, func(writer *Writer, ts TimeSource, eventFlag EventFlag, state ...interface{}) {
		assert.Equal(EventInfo, eventFlag)
		if len(state) > 0 {
			format = state[0].(string)
		}
	})
	a.Sync().Infof("this is a %s", "test")
	assert.Equal("this is a %s", format)
	assert.Equal("[info] this is a test\n", buffer.String())
}

func TestSyncAgentDebugf(t *testing.T) {
	assert := assert.New(t)

	buffer := bytes.NewBuffer(nil)
	a := None(NewWriter(buffer))
	a.EnableEvent(EventDebug)
	a.Writer().SetShowTimestamp(false)
	a.Writer().SetShowLabel(false)
	a.Writer().SetUseAnsiColors(false)
	a.Sync().Infof("this is a %s", "test")
	a.Sync().Debugf("this is a %s", "test")
	assert.Equal("[debug] this is a test\n", buffer.String())
}

func TestSyncAgentWarningf(t *testing.T) {
	assert := assert.New(t)

	buffer := bytes.NewBuffer(nil)
	a := None(NewWriter(buffer))
	a.EnableEvent(EventWarning)
	a.Writer().SetShowTimestamp(false)
	a.Writer().SetShowLabel(false)
	a.Writer().SetUseAnsiColors(false)
	a.Sync().Infof("this is a %s", "test")
	a.Sync().Warningf("this is a %s", "test")
	assert.Equal("[warning] this is a test\n", buffer.String())
}

func TestSyncAgentErrorf(t *testing.T) {
	assert := assert.New(t)

	buffer := bytes.NewBuffer(nil)
	a := None(NewWriter(buffer))
	a.EnableEvent(EventError)
	a.Writer().SetShowTimestamp(false)
	a.Writer().SetShowLabel(false)
	a.Writer().SetUseAnsiColors(false)
	a.Sync().Infof("this is a %s", "test")
	a.Sync().Errorf("this is a %s", "test")
	assert.Equal("[error] this is a test\n", buffer.String())
}

func TestSyncAgentOnEvent(t *testing.T) {
	assert := assert.New(t)

	a := None()
	a.EnableEvent(EventFlag("foo"))
	a.AddEventListener(EventFlag("foo"), func(writer *Writer, ts TimeSource, eventFlag EventFlag, state ...interface{}) {
		assert.NotEmpty(state)
		assert.Equal("bar", state[0])
	})
	a.OnEvent("foo", "bar")
}
