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
