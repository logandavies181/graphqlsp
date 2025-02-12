package state

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDef(t *testing.T) {
	s, err := NewFromFile("../test/schema.graphql")
	assert.NoError(t, err)

	pos := s.GetDefinitionOf(55, 9)

	assert.NotNil(t, pos)
	assert.Equal(t, &Position{1, 6, 5, false}, pos)
}

func TestPreludeState(t *testing.T) {
	s := PreludeState()
	assert.NotNil(t, s)
	assert.NotNil(t, s.schema.Types["Int"])
	assert.Equal(t, 4, s.schema.Types["Int"].Position.Line)
	assert.Equal(t, 8, s.schema.Types["Int"].Position.Column)
}
