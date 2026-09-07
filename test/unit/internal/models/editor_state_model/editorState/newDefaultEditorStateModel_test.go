package editorState_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/stretchr/testify/assert"
)

func TestWhenDefaultStateIsCreated_UsesRandomTopology(t *testing.T) {
	t.Parallel()
	// Arrange

	// Act
	state := editor_state_model.NewDefaultEditorStateModel()

	// Assert
	assert.Equal(t, config.TopologyRandom, state.Topology)
}
