package logger

import (
	"bytes"
	"errors"
	"sync"
	"testing"

	assert "github.com/blendlabs/go-assert"
)

func TestAgentDefault(t *testing.T) {
	assert := assert.New(t)
	defer SetDefault(nil)

	nilAgent := Default()
	assert.Nil(nilAgent.Output)
	assert.Nil(nilAgent.ErrorOutput)

	buffer := bytes.NewBuffer([]byte{})
	SetDefault(NewAgent(EventAll, buffer, buffer))
	notNilAgent := Default()
	assert.NotNil(notNilAgent.Output)
	assert.NotNil(notNilAgent.ErrorOutput)

	assert.NotNil(notNilAgent.Output.(*SyncWriter))
	assert.NotNil(notNilAgent.ErrorOutput.(*SyncWriter))
}

func TestAgentEventListeners(t *testing.T) {
	assert := assert.New(t)

	buffer := bytes.NewBuffer([]byte{})
	agent := NewAgent(EventAll, buffer, buffer)

	wg := sync.WaitGroup{}
	wg.Add(2)
	agent.AddEventListener(4, func(flag uint64, state ...interface{}) {
		defer wg.Done()
		assert.Equal(4, flag)
		assert.Len(state, 1)
		assert.Equal("not test state", state[0])
	})

	agent.AddEventListener(2, func(flag uint64, state ...interface{}) {
		defer wg.Done()
		assert.Equal(2, flag)
		assert.Len(state, 1)
		assert.Equal("test state", state[0])
	})
	agent.OnEvent(2, "test state")
	agent.OnEvent(4, "not test state")
	wg.Wait()
}

func TestAgentErrorListeners(t *testing.T) {
	assert := assert.New(t)

	buffer := bytes.NewBuffer([]byte{})
	agent := NewAgent(EventAll, buffer, buffer)
	wg := sync.WaitGroup{}
	wg.Add(1)
	agent.AddErrorListener(func(err error) {
		defer wg.Done()
		assert.Equal("test error", err.Error())
	})
	agent.OnError(errors.New("test error"))
	wg.Wait()
}

func TestAgentPrintf(t *testing.T) {
	assert := assert.New(t)

	buffer := bytes.NewBuffer([]byte{})
	agent := NewAgent(EventAll, buffer, buffer)
	agent.ShowTimestamp = false
	agent.ShowLabel = false
	agent.UseAnsiColors = false
	agent.Printf("test %s", "string")
	assert.Equal("test string\n", string(buffer.Bytes()))
}

func TestAgentPrintfWithLabel(t *testing.T) {
	assert := assert.New(t)

	buffer := bytes.NewBuffer([]byte{})
	agent := NewAgent(EventAll, buffer, buffer)
	agent.Label = "unit-test"
	agent.ShowTimestamp = false
	agent.ShowLabel = true
	agent.UseAnsiColors = false
	agent.Printf("test %s", "string")
	assert.Equal("unit-test test string\n", string(buffer.Bytes()))
}

func TestAgentPrintfWithLabelColorized(t *testing.T) {
	assert := assert.New(t)

	buffer := bytes.NewBuffer([]byte{})
	agent := NewAgent(EventAll, buffer, buffer)
	agent.Label = "unit-test"
	agent.ShowTimestamp = false
	agent.ShowLabel = true
	agent.UseAnsiColors = true
	agent.Printf("test %s", "string")
	assert.Equal(ColorBlue.Apply("unit-test")+" test string\n", string(buffer.Bytes()))
}

func TestAgentErrorOutputCoalesced(t *testing.T) {
	assert := assert.New(t)

	buffer := bytes.NewBuffer([]byte{})
	agent := NewAgent(EventAll, buffer)
	agent.ShowTimestamp = false
	agent.UseAnsiColors = false

	agent.Errorf("test %s", "string")
	assert.Equal("test string\n", string(buffer.Bytes()))
}

func TestAgentErrorOutput(t *testing.T) {
	assert := assert.New(t)

	stdout := bytes.NewBuffer([]byte{})
	stderr := bytes.NewBuffer([]byte{})
	agent := NewAgent(EventAll, stdout, stderr)
	agent.ShowTimestamp = false
	agent.UseAnsiColors = false

	agent.Errorf("test %s", "string")
	assert.Equal(0, stdout.Len())
	assert.Equal("test string\n", string(stderr.Bytes()))
}

func TestAgentCheckVerbosity(t *testing.T) {
	assert := assert.New(t)
	buffer := bytes.NewBuffer([]byte{})
	agent := NewAgent(EventInfo, buffer)
	assert.True(agent.CheckVerbosity(EventInfo))
	assert.False(agent.CheckVerbosity(EventDebug))
}
