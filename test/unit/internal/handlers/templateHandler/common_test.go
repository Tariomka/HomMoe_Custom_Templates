package templateHandler_test

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers/handler_interfaces"
	"github.com/Tariomka/hommoe_custom_templates/internal/mappers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
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
	manualReapply     *test_helpers.ManualReapplyServiceMock
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
		manualReapply:     &test_helpers.ManualReapplyServiceMock{},
		fileService:       &test_helpers.FileServiceMock{},
		previewGenerator:  &test_helpers.PreviewGeneratorServiceMock{},
		stateHandler:      &test_helpers.StateHandlerMock{},
	}

	fixture.handler = handlers.NewTemplateHandler(
		fixture.templateGenerator,
		fixture.mapper,
		mappers.NewTemplateMapper(),
		fixture.contentProvider,
		fixture.connectionEditor,
		fixture.zoneEditor,
		fixture.manualReapply,
		fixture.fileService,
		fixture.previewGenerator,
		fixture.stateHandler,
	)

	return fixture
}

func toDto(state editor_state_model.EditorState) editor_state_dto.EditorStateDto {
	return editor_state_dto.EditorStateDto{EditorState: state}
}
