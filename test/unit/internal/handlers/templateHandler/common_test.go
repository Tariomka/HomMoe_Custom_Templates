package templateHandler_test

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers/handler_interfaces"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/mock"
)

// templateHandlerFixture holds the handler together with every mock it was
// built from, so a test can arrange one collaborator and assert on another.
type templateHandlerFixture struct {
	handler           handler_interfaces.ITemplateHandler
	templateGenerator *test_helpers.TemplateGeneratorMock
	mapper            *test_helpers.GeneratorConfigMapperMock
	contentProvider   *test_helpers.MandatoryContentProviderMock
	connectionEditor  *test_helpers.ConnectionEditorServiceMock
	zoneEditor        *test_helpers.ZoneEditorServiceMock
	roadPolicy        *test_helpers.RoadPolicyServiceMock
	manualEditDelta   *test_helpers.ManualEditDeltaServiceMock
	fileService       *test_helpers.FileServiceMock
	previewGenerator  *test_helpers.PreviewGeneratorServiceMock
	stateHandler      *test_helpers.StateHandlerMock
}

func newTemplateHandlerFixture() *templateHandlerFixture {
	fixture := &templateHandlerFixture{
		templateGenerator: &test_helpers.TemplateGeneratorMock{},
		mapper:            &test_helpers.GeneratorConfigMapperMock{},
		contentProvider:   &test_helpers.MandatoryContentProviderMock{},
		connectionEditor:  &test_helpers.ConnectionEditorServiceMock{},
		zoneEditor:        &test_helpers.ZoneEditorServiceMock{},
		roadPolicy:        &test_helpers.RoadPolicyServiceMock{},
		manualEditDelta:   &test_helpers.ManualEditDeltaServiceMock{},
		fileService:       &test_helpers.FileServiceMock{},
		previewGenerator:  &test_helpers.PreviewGeneratorServiceMock{},
		stateHandler:      &test_helpers.StateHandlerMock{},
	}

	fixture.handler = handlers.NewTemplateHandler(
		fixture.templateGenerator,
		fixture.mapper,
		fixture.contentProvider,
		fixture.connectionEditor,
		fixture.zoneEditor,
		fixture.roadPolicy,
		fixture.manualEditDelta,
		fixture.fileService,
		fixture.previewGenerator,
		fixture.stateHandler,
	)

	return fixture
}

func toDto(state editor_state_model.EditorState) editor_state_dto.EditorStateDto {
	return editor_state_dto.EditorStateDto{EditorState: state}
}

// arrangeGeneration makes generating from state succeed with generated, the
// state passing validation unchanged and mapping to configuration.
func (this *templateHandlerFixture) arrangeGeneration(
	state editor_state_model.EditorState,
	configuration *config.GeneratorConfig,
	generated *template_model.Template) {
	this.stateHandler.On("ValidateEditorState", state, true).
		Return(editor_state_dto.EditorStateValidationDto{State: state})
	this.mapper.On("FromEditorState", state).Return(configuration)
	this.templateGenerator.On("SetConfiguration", mock.Anything).Return()
	this.templateGenerator.On("Generate").Return(generated, nil)
}

// newTemplateWithLayout builds a template whose single variant holds the given layout.
func newTemplateWithLayout(
	zones []template_model.Zone,
	connections []template_model.Connection) *template_model.Template {
	return &template_model.Template{
		Variants: []template_model.Variant{{Zones: zones, Connections: connections}},
	}
}
