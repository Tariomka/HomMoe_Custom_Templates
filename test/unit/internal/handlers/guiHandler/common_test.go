package guiHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/composition"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers/handler_interfaces"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/connection_editor"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/generation_tuning"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers"
	zone_services "github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
	"github.com/stretchr/testify/require"
)

// generateDefaultTemplate produces a real generated template from the default
// editor state so update tests operate on realistic zones and connections.
func generateDefaultTemplate(t *testing.T, handler handler_interfaces.ITemplateHandler) *entities.RmgTemplate {
	t.Helper()

	loadDto, err := handler.GenerateTemplate(dtos.NewDefaultEditorStateDto())
	require.NoError(t, err)
	require.NotNil(t, loadDto.Template)
	require.NotEmpty(t, loadDto.Template.Variants)

	return loadDto.Template
}

// newProductionGuiHandler builds the same handler graph the application uses.
func newProductionGuiHandler() handler_interfaces.IGuiHandler {
	return composition.InitializeGuiHandler()
}

// newManualReapplyService builds the castle re-apply service with the same
// collaborators the handler graph wires.
func newManualReapplyService() *connection_editor.ManualReapplyService {
	return connection_editor.NewManualReapplyService(
		connection_editor.NewDefaultZoneEditorService(),
		zone_services.NewZoneClassifier(),
		generation_tuning.NewGenerationTuningFactory(),
	)
}

// newMandatoryContentProvider builds the mandatory-content provider with the
// same collaborators the handler graph wires.
func newMandatoryContentProvider() *providers.MandatoryContentProvider {
	return providers.NewMandatoryContentProvider(
		zone_services.NewZoneClassifier(),
		connection_editor.NewDefaultZoneEditorService(),
	)
}
