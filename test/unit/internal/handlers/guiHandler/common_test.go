package guiHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/composition"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers/handler_interfaces"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/connection_editor"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/generation_tuning"
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

// newManualReapplyService builds the castle re-apply service with the same
// collaborators the handler graph wires.
func newManualReapplyService() connection_editor.IManualReapplyService {
	return connection_editor.NewManualReapplyService(
		test_helpers.NewZoneEditorService(),
		zone_services.NewCastleFactory(),
		zone_services.NewZoneTierService(),
		generation_tuning.NewGenerationTuningFactory(),
	)
}

// newMandatoryContentProvider builds the mandatory-content provider with the
// same collaborators the handler graph wires.
func newMandatoryContentProvider() provider_interfaces.IMandatoryContentProvider {
	return providers.NewMandatoryContentProvider(
		zone_services.NewZoneTierService(),
		test_helpers.NewZoneEditorService(),
	)
}
