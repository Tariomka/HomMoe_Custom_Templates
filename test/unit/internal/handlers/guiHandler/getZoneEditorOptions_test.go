package guiHandler_test

import (
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers"
	"github.com/Tariomka/hommoe_custom_templates/internal/mappers"
	"github.com/stretchr/testify/assert"
)

func TestWhenFullVariantZoneCountIsProvided_ReturnsMappedZoneEditorOptions(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := handlers.NewGuiHandler()
	state := dtos.NewDefaultEditorStateDto()
	fullVariantZoneCount := state.PlayerCount + state.NeutralZoneCount
	configuration := mappers.NewConfigMapper().FromEditorState(state)
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
