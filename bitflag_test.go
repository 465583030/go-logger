package logger

import (
	"testing"

	"github.com/blendlabs/go-assert"
)

func TestBitFlagCombine(t *testing.T) {
	assert := assert.New(t)

	three := BitFlagCombine(1, 2)
	assert.Equal(3, three)
}

func TestBitFlagAny(t *testing.T) {
	assert := assert.New(t)

	one := uint64(1 << 0)
	two := uint64(1 << 1)
	four := uint64(1 << 2)
	eight := uint64(1 << 3)
	sixteen := uint64(1 << 4)
	invalid := uint64(1 << 5)

	masterFlag := BitFlagCombine(one, two, four, eight)
	checkFlag := BitFlagCombine(one, sixteen)
	assert.True(BitFlagAny(masterFlag, checkFlag))
	assert.False(BitFlagAny(masterFlag, invalid))
}

func TestBitFlagAll(t *testing.T) {
	assert := assert.New(t)

	one := uint64(1 << 0)
	two := uint64(1 << 1)
	four := uint64(1 << 2)
	eight := uint64(1 << 3)
	sixteen := uint64(1 << 4)

	masterFlag := BitFlagCombine(one, two, four, eight)
	checkValidFlag := BitFlagCombine(one, two)
	checkInvalidFlag := BitFlagCombine(one, sixteen)
	assert.True(BitFlagAll(masterFlag, checkValidFlag))
	assert.False(BitFlagAll(masterFlag, checkInvalidFlag))
}
