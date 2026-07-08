package editorState_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config/config_inner"
	"github.com/stretchr/testify/assert"
)

func TestWhenTopologyChangedSinceSnapshot_ReportsLayoutChanged(t *testing.T) {
	// Arrange
	state := models.NewEditorState()
	state.SnapshotCurrentState()
	state.UpdateCurrentState(func(dto *dtos.EditorStateDto) { dto.Topology = config_inner.TopologyChain })

	// Act
	layoutChanged := state.WasLayoutChanged()

	// Assert
	assert.True(t, layoutChanged)
}

func TestWhenOnlyDensityChangedSinceSnapshot_ReportsLayoutNotChanged(t *testing.T) {
	// Arrange
	state := models.NewEditorState()
	state.SnapshotCurrentState()
	state.UpdateCurrentState(func(dto *dtos.EditorStateDto) { dto.ResourceDensityPercent = 50 })

	// Act
	layoutChanged := state.WasLayoutChanged()

	// Assert
	assert.False(t, layoutChanged)
}
