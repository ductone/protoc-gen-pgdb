package pgtest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPGConfig(t *testing.T) {
	assert := assert.New(t)

	config := New().From("/usr/bin").DataDir("/tmp/data").Persistent()

	assert.True(config.IsPersistent)
	assert.EqualValues("/tmp/data", config.Dir)
	assert.EqualValues("/usr/bin", config.BinDir)
}
