package logger

import (
	"testing"

	"github.com/blendlabs/go-assert"
)

func TestEventFlagCombine(t *testing.T) {
	assert := assert.New(t)

	three := EventFlagCombine(1, 2)
	assert.Equal(3, three)
}

func TestEventFlagAny(t *testing.T) {
	assert := assert.New(t)

	one := uint64(1 << 0)
	two := uint64(1 << 1)
	four := uint64(1 << 2)
	eight := uint64(1 << 3)
	sixteen := uint64(1 << 4)
	invalid := uint64(1 << 5)

	masterFlag := EventFlagCombine(one, two, four, eight)
	checkFlag := EventFlagCombine(one, sixteen)
	assert.True(EventFlagAny(masterFlag, checkFlag))
	assert.False(EventFlagAny(masterFlag, invalid))
}

func TestEventFlagAll(t *testing.T) {
	assert := assert.New(t)

	one := uint64(1 << 0)
	two := uint64(1 << 1)
	four := uint64(1 << 2)
	eight := uint64(1 << 3)
	sixteen := uint64(1 << 4)

	masterFlag := EventFlagCombine(one, two, four, eight)
	checkValidFlag := EventFlagCombine(one, two)
	checkInvalidFlag := EventFlagCombine(one, sixteen)
	assert.True(EventFlagAll(masterFlag, checkValidFlag))
	assert.False(EventFlagAll(masterFlag, checkInvalidFlag))
}
