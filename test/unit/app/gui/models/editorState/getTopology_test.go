package editorState_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/stretchr/testify/assert"
)

func TestWhenTopologyWasUpdated_GetTopologyReturnsIt(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()
	state.UpdateCurrentState(func(dto *dtos.EditorStateDto) { dto.Topology = config.TopologyHubAndSpoke })

	// Act
	actual := state.GetTopology()

	// Assert
	assert.Equal(t, config.TopologyHubAndSpoke, actual)
}
