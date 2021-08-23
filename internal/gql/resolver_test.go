package gql

import (
	"testing"

	"github.com/tj/assert"
)

func TestResolvers(t *testing.T) {
	handler, err := New(Config{})
	assert.Nil(t, err)
	assert.NotNil(t, handler)
}
