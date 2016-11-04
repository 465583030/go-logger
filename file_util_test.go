package logger

import (
	"os"
	"testing"

	assert "github.com/blendlabs/go-assert"
)

func TestFileCreateOrOpen(t *testing.T) {
	assert := assert.New(t)

	tempFilePath := os.TempDir() + UUIDv4()
	f, err := File.CreateOrOpen(tempFilePath)
	assert.Nil(err)
	assert.NotNil(f)
	defer f.Close()
	defer func() {
		os.Remove(tempFilePath)
	}()
	_, err = f.Stat()
	assert.Nil(err)
}

func TestFileParseSize(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(2*Gigabyte, File.ParseSize("2gb", 1))
	assert.Equal(3*Megabyte, File.ParseSize("3mb", 1))
	assert.Equal(123*Kilobyte, File.ParseSize("123kb", 1))
	assert.Equal(12345, File.ParseSize("12345", 1))
	assert.Equal(1, File.ParseSize("", 1))
}
