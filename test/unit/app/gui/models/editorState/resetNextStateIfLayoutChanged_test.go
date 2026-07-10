package editorState_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config/config_inner"
	"github.com/stretchr/testify/assert"
)

func layoutChangedStateWithNext() *models.EditorState {
	state := models.NewEditorState()
	state.SnapshotCurrentState()
	state.SetNextState(state.GetCurrentState())
	state.UpdateCurrentState(func(dto *dtos.EditorStateDto) { dto.Topology = config_inner.TopologyChain })
	return state
}

func nonLayoutChangedStateWithNext() *models.EditorState {
	state := models.NewEditorState()
	state.SnapshotCurrentState()
	state.SetNextState(state.GetCurrentState())
	state.UpdateCurrentState(func(dto *dtos.EditorStateDto) { dto.ResourceDensityPercent = 50 })
	return state
}

func TestWhenLayoutOptionChangedSinceSnapshot_ReportsReset(t *testing.T) {
	t.Parallel()
	// Arrange
	state := layoutChangedStateWithNext()

	// Act
	wasReset := state.ResetNextStateIfLayoutChanged()

	// Assert
	assert.True(t, wasReset)
}

func TestWhenLayoutOptionChangedSinceSnapshot_NextStateIsCleared(t *testing.T) {
	t.Parallel()
	// Arrange
	state := layoutChangedStateWithNext()

	// Act
	state.ResetNextStateIfLayoutChanged()

	// Assert
	assert.False(t, state.HasNextState())
}

func TestWhenOnlyNonLayoutOptionChangedSinceSnapshot_ReportsNoReset(t *testing.T) {
	t.Parallel()
	// Arrange
	state := nonLayoutChangedStateWithNext()

	// Act
	wasReset := state.ResetNextStateIfLayoutChanged()

	// Assert
	assert.False(t, wasReset)
}

func TestWhenOnlyNonLayoutOptionChangedSinceSnapshot_NextStateIsKept(t *testing.T) {
	t.Parallel()
	// Arrange
	state := nonLayoutChangedStateWithNext()

	// Act
	state.ResetNextStateIfLayoutChanged()

	// Assert
	assert.True(t, state.HasNextState())
}
