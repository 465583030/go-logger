package logger

import (
	"bytes"
	"testing"

	assert "github.com/blendlabs/go-assert"
)

func TestLogWriterPrintf(t *testing.T) {
	assert := assert.New(t)

	buffer := bytes.NewBuffer(nil)
	writer := NewWriterWithError(buffer, buffer)
	writer.showTimestamp = false
	writer.showLabel = false
	writer.useAnsiColors = false
	writer.Printf("test %s", "string")
	assert.Equal("test string\n", string(buffer.Bytes()))
}

func TestLogWriterPrintfWithLabel(t *testing.T) {
	assert := assert.New(t)

	buffer := bytes.NewBuffer(nil)
	writer := NewWriterWithError(buffer, buffer)
	writer.label = "unit-test"
	writer.showTimestamp = false
	writer.showLabel = true
	writer.useAnsiColors = false
	writer.Printf("test %s", "string")
	assert.Equal("unit-test test string\n", string(buffer.Bytes()))
}

func TestLogWriterPrintfWithLabelColorized(t *testing.T) {
	assert := assert.New(t)

	buffer := bytes.NewBuffer([]byte{})
	writer := NewWriterWithError(buffer, buffer)
	writer.label = "unit-test"
	writer.showTimestamp = false
	writer.showLabel = true
	writer.useAnsiColors = true
	writer.Printf("test %s", "string")
	assert.Equal(ColorBlue.Apply("unit-test")+" test string\n", string(buffer.Bytes()))
}

func TestWriterErrorOutputCoalesced(t *testing.T) {
	assert := assert.New(t)

	buffer := bytes.NewBuffer(nil)
	writer := NewWriter(buffer)
	writer.showTimestamp = false
	writer.useAnsiColors = false

	writer.Errorf("test %s", "string")
	assert.Equal("test string\n", string(buffer.Bytes()))
}

func TestWriterErrorOutput(t *testing.T) {
	assert := assert.New(t)

	stdout := bytes.NewBuffer(nil)
	stderr := bytes.NewBuffer(nil)
	writer := NewWriterWithError(stdout, stderr)
	writer.showTimestamp = false
	writer.useAnsiColors = false

	writer.Errorf("test %s", "string")
	assert.Equal(0, stdout.Len())
	assert.Equal("test string\n", string(stderr.Bytes()))
}
