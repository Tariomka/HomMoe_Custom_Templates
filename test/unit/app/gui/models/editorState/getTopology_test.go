package editorState_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/stretchr/testify/assert"
)

func TestWhenTopologyWasUpdated_GetTopologyReturnsIt(t *testing.T) {
	t.Parallel()
	// Arrange
	state := newEditorState()
	state.UpdateCurrentState(func(dto *editor_state_dto.EditorStateDto) { dto.Topology = config.TopologyHubAndSpoke })

	// Act
	actual := state.GetTopology()

	// Assert
	assert.Equal(t, config.TopologyHubAndSpoke, actual)
}
