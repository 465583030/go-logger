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
	a := All(NewLogWriter(buffer))
	a.Writer().SetShowTimestamp(false)
	a.Writer().SetShowLabel(false)
	a.Writer().SetUseAnsiColors(false)
	a.AddEventListener(EventInfo, func(writer Logger, ts TimeSource, eventFlag EventFlag, state ...interface{}) {
		assert.Equal(EventInfo, eventFlag)
		if len(state) > 0 {
			format = state[0].(string)
		}
	})
	a.Sync().Infof("this is a %s", "test")
	assert.Equal("this is a %s", format)
	assert.Equal("info this is a test\n", buffer.String())
}

func TestSyncAgentDebugf(t *testing.T) {
	assert := assert.New(t)

	buffer := bytes.NewBuffer(nil)
	a := None(NewLogWriter(buffer))
	a.EnableEvent(EventDebug)
	a.Writer().SetShowTimestamp(false)
	a.Writer().SetShowLabel(false)
	a.Writer().SetUseAnsiColors(false)
	a.Sync().Infof("this is a %s", "test")
	a.Sync().Debugf("this is a %s", "test")
	assert.Equal("debug this is a test\n", buffer.String())
}
