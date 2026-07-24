package guiHandler_test

import (
	"encoding/json"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers"
	"github.com/Tariomka/hommoe_custom_templates/internal/mappers"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/connection_editor"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenCastleSettingsChange_ReturnsServiceEquivalentZones(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := handlers.NewGuiHandler()
	template := test_helpers.GetDefaultTemplate()
	expectedZones := cloneZones(t, template.Variants[0].Zones)
	actualZones := cloneZones(t, template.Variants[0].Zones)
	editorState := dtos.NewDefaultEditorStateDto()
	editorState.NeutralZoneCastles = 2
	changes := editor_state_dto.CastleSettingChanges{NeutralSimple: true}
	configuration := mappers.NewConfigMapper().FromEditorState(editorState)
	connection_editor.ApplyCastleSettingChanges(expectedZones, changes, configuration)

	// Act
	result := handler.ReapplyCastleSettings(dtos.CastleSettingsReapplyRequestDto{
		Zones:       actualZones,
		Changes:     changes,
		EditorState: editorState,
	})

	// Assert
	assert.Equal(t, expectedZones, result)
}

func cloneZones(t *testing.T, zones []entities.Zone) []entities.Zone {
	t.Helper()

	encoded, err := json.Marshal(zones)
	require.NoError(t, err)

	var cloned []entities.Zone
	require.NoError(t, json.Unmarshal(encoded, &cloned))
	return cloned
}
