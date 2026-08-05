package guiHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
)

func TestWhenFullVariantZoneCountIsProvided_ReturnsMappedZoneEditorOptions(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	state := dtos.NewDefaultEditorStateDto()
	fullVariantZoneCount := state.PlayerCount + state.NeutralZoneCount
	configuration := test_helpers.NewConfigMapper().FromEditorState(state)
	expected := dtos.ZoneEditorOptionsDto{
		Topology:      state.Topology,
		Tuning:        test_helpers.NewGenerationTuning(configuration, fullVariantZoneCount),
		GenerateRoads: state.GenerateRoads,
	}

	// Act
	result := handler.GetZoneEditorOptions(state, fullVariantZoneCount)

	// Assert
	assert.Equal(t, expected, result)
}
