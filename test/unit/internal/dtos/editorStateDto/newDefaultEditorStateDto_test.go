package editorStateDto_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/stretchr/testify/assert"
)

func TestWhenDefaultStateIsCreated_UsesRandomTopology(t *testing.T) {
	// Arrange - defaults require no setup.

	// Act
	state := dtos.NewDefaultEditorStateDto()

	// Assert
	assert.Equal(t, config.TopologyRandom, state.Topology)
}
