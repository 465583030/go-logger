package logger

import (
	"net/http"
	"testing"

	"os"

	assert "github.com/blendlabs/go-assert"
)

func TestOpenOrCreateFile(t *testing.T) {
	assert := assert.New(t)

	tempFilePath := os.TempDir() + "open_or_create_test.txt"
	f, err := OpenOrCreateFile(tempFilePath)
	assert.Nil(err)
	assert.NotNil(f)
	defer f.Close()
}

func TestGetIP(t *testing.T) {
	assert := assert.New(t)

	hdr := http.Header{}
	hdr.Set("X-Forwarded-For", "1")
	r := http.Request{
		Header: hdr,
	}
	assert.Equal("1", GetIP(&r))

	hdr = http.Header{}
	hdr.Set("X-FORWARDED-FOR", "1")
	r = http.Request{
		Header: hdr,
	}
	assert.Equal("1", GetIP(&r))

	hdr = http.Header{}
	hdr.Set("X-FORWARDED-FOR", "1,2,3")
	r = http.Request{
		Header: hdr,
	}
	assert.Equal("1", GetIP(&r))

	hdr = http.Header{}
	hdr.Set("X-Real-Ip", "1")
	r = http.Request{
		Header: hdr,
	}
	assert.Equal("1", GetIP(&r))

	hdr = http.Header{}
	hdr.Set("X-REAL-IP", "1")
	r = http.Request{
		Header: hdr,
	}
	assert.Equal("1", GetIP(&r))

	hdr = http.Header{}
	hdr.Set("X-REAL-IP", "1,2,3")
	r = http.Request{
		Header: hdr,
	}
	assert.Equal("1", GetIP(&r))

	r = http.Request{
		RemoteAddr: "1:1",
	}
	assert.Equal("1", GetIP(&r))
}
