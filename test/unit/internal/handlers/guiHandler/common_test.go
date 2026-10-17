package guiHandler_test

import (
	"reflect"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/composition"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers/handler_interfaces"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/provider_interfaces"
	zone_services "github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/require"
)

// generateDefaultTemplate produces a real generated template from the default
// editor state so update tests operate on realistic zones and connections.
func generateDefaultTemplate(t *testing.T, handler handler_interfaces.ITemplateHandler) *template_model.Template {
	t.Helper()

	loadDto, err := handler.GenerateTemplate(toDto(editor_state_model.NewDefaultEditorStateModel()))
	require.NoError(t, err)
	require.NotNil(t, loadDto.Template)
	require.NotEmpty(t, loadDto.Template.Variants)

	return loadDto.Template
}

func toDto(state editor_state_model.EditorState) editor_state_dto.EditorStateDto {
	return editor_state_dto.EditorStateDto{EditorState: state}
}

func toDtoPointer(state *editor_state_model.EditorState) *editor_state_dto.EditorStateDto {
	return &editor_state_dto.EditorStateDto{EditorState: *state}
}

// newProductionGuiHandler builds the same handler graph the application uses.
func newProductionGuiHandler() handler_interfaces.IGuiHandler {
	return composition.InitializeGuiHandler()
}

// newDeterministicState returns the default state with a topology whose layout
// two generations reproduce exactly (the default Random topology does not).
func newDeterministicState() editor_state_model.EditorState {
	state := editor_state_model.NewDefaultEditorStateModel()
	state.Topology = config.TopologyCircles
	return state
}

// generateTemplateFor produces a real generated template from state.
func generateTemplateFor(
	t *testing.T,
	handler handler_interfaces.ITemplateHandler,
	state editor_state_model.EditorState) *template_model.Template {
	t.Helper()

	loadDto, err := handler.GenerateTemplate(toDto(state))
	require.NoError(t, err)
	require.NotNil(t, loadDto.Template)
	require.NotEmpty(t, loadDto.Template.Variants)

	return loadDto.Template
}

// zoneEditsWithOverrides returns the zone edits that carry at least one field override.
func zoneEditsWithOverrides(settings editor_state_model.ManualEditSettings) []editor_state_model.ManualZoneEdit {
	var overridden []editor_state_model.ManualZoneEdit
	for _, edit := range settings.ManualZones {
		if !reflect.ValueOf(edit.FieldOverrides).IsZero() || len(edit.Unset) > 0 {
			overridden = append(overridden, edit)
		}
	}
	return overridden
}

// connectionEditsWithOverrides returns the connection edits that carry at least one field override.
func connectionEditsWithOverrides(
	settings editor_state_model.ManualEditSettings) []editor_state_model.ManualConnectionEdit {
	var overridden []editor_state_model.ManualConnectionEdit
	for _, edit := range settings.ManualConnections {
		if !reflect.ValueOf(edit.FieldOverrides).IsZero() || len(edit.Unset) > 0 {
			overridden = append(overridden, edit)
		}
	}
	return overridden
}

// newMandatoryContentProvider builds the mandatory-content provider with the
// same collaborators the handler graph wires.
func newMandatoryContentProvider() provider_interfaces.IMandatoryContentProvider {
	return providers.NewMandatoryContentProvider(
		zone_services.NewZoneTierService(),
		test_helpers.NewZoneEditorService(),
	)
}
