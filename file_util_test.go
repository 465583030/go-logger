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
