package logger

import (
	"bytes"
	"testing"

	assert "github.com/blendlabs/go-assert"
)

func TestLogWriterPrintf(t *testing.T) {
	assert := assert.New(t)

	buffer := bytes.NewBuffer([]byte{})
	writer := NewLogWriter(buffer, buffer)
	writer.ShowTimestamp = false
	writer.ShowLabel = false
	writer.UseAnsiColors = false
	writer.Printf("test %s", "string")
	assert.Equal("test string\n", string(buffer.Bytes()))
}

func TestLogWriterPrintfWithLabel(t *testing.T) {
	assert := assert.New(t)

	buffer := bytes.NewBuffer([]byte{})
	writer := NewLogWriter(buffer, buffer)
	writer.Label = "unit-test"
	writer.ShowTimestamp = false
	writer.ShowLabel = true
	writer.UseAnsiColors = false
	writer.Printf("test %s", "string")
	assert.Equal("unit-test test string\n", string(buffer.Bytes()))
}

func TestLogWriterPrintfWithLabelColorized(t *testing.T) {
	assert := assert.New(t)

	buffer := bytes.NewBuffer([]byte{})
	writer := NewLogWriter(buffer, buffer)
	writer.Label = "unit-test"
	writer.ShowTimestamp = false
	writer.ShowLabel = true
	writer.UseAnsiColors = true
	writer.Printf("test %s", "string")
	assert.Equal(ColorBlue.Apply("unit-test")+" test string\n", string(buffer.Bytes()))
}

func TestWriterErrorOutputCoalesced(t *testing.T) {
	assert := assert.New(t)

	buffer := bytes.NewBuffer([]byte{})
	writer := NewLogWriter(buffer)
	writer.ShowTimestamp = false
	writer.UseAnsiColors = false

	writer.Errorf("test %s", "string")
	assert.Equal("test string\n", string(buffer.Bytes()))
}

func TestWriterErrorOutput(t *testing.T) {
	assert := assert.New(t)

	stdout := bytes.NewBuffer([]byte{})
	stderr := bytes.NewBuffer([]byte{})
	writer := NewLogWriter(stdout, stderr)
	writer.ShowTimestamp = false
	writer.UseAnsiColors = false

	writer.Errorf("test %s", "string")
	assert.Equal(0, stdout.Len())
	assert.Equal("test string\n", string(stderr.Bytes()))
}
