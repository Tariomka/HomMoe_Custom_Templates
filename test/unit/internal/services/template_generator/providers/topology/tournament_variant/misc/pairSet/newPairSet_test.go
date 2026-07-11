package pairSet_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/tournament_variant/misc"
	"github.com/stretchr/testify/assert"
)

func TestWhenNewPairSetIsCreated_SetIsEmpty(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	set := misc.NewPairSet()

	// Assert
	assert.Empty(t, *set)
}
