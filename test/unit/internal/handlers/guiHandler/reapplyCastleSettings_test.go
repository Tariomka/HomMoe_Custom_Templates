package guiHandler_test

import (
	"encoding/json"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenCastleSettingsChange_ReturnsServiceEquivalentZones(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	template := test_helpers.GetDefaultTemplateModel()
	expectedZones := cloneZones(t, template.Variants[0].Zones)
	actualZones := cloneZones(t, template.Variants[0].Zones)
	editorState := editor_state_model.NewDefaultEditorStateModel()
	editorState.NeutralZoneCastles = 2
	changes := editor_state_model.CastleSettingChanges{NeutralSimple: true}
	configuration := test_helpers.NewConfigMapper().FromEditorState(editorState)
	newManualReapplyService().ApplyCastleSettingChanges(expectedZones, changes, configuration)

	// Act
	result := handler.ReapplyCastleSettings(dtos.CastleSettingsReapplyRequestDto{
		Zones:       actualZones,
		Changes:     changes,
		EditorState: toDto(editorState),
	})

	// Assert
	assert.Equal(t, expectedZones, result)
}

func cloneZones(t *testing.T, zones []template_model.Zone) []template_model.Zone {
	t.Helper()

	encoded, err := json.Marshal(zones)
	require.NoError(t, err)

	var cloned []template_model.Zone
	require.NoError(t, json.Unmarshal(encoded, &cloned))
	return cloned
}
