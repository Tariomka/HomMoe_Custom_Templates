package zoneEditorHandler_test

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers/handler_interfaces"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
)

// zoneEditorHandlerFixture holds the handler together with every mock it was
// built from, so a test can arrange one collaborator and assert on another.
type zoneEditorHandlerFixture struct {
	handler          handler_interfaces.IZoneEditorHandler
	mapper           *test_helpers.GeneratorConfigMapperMock
	tierService      *test_helpers.ZoneTierServiceMock
	connectionEditor *test_helpers.ConnectionEditorServiceMock
	zoneEditor       *test_helpers.ZoneEditorServiceMock
	geometry         *test_helpers.ZoneEditorGeometryServiceMock
	tuningFactory    *test_helpers.GenerationTuningFactoryMock
}

func newZoneEditorHandlerFixture() *zoneEditorHandlerFixture {
	fixture := &zoneEditorHandlerFixture{
		mapper:           &test_helpers.GeneratorConfigMapperMock{},
		tierService:      &test_helpers.ZoneTierServiceMock{},
		connectionEditor: &test_helpers.ConnectionEditorServiceMock{},
		zoneEditor:       &test_helpers.ZoneEditorServiceMock{},
		geometry:         &test_helpers.ZoneEditorGeometryServiceMock{},
		tuningFactory:    &test_helpers.GenerationTuningFactoryMock{},
	}

	fixture.handler = handlers.NewZoneEditorHandler(
		fixture.mapper,
		fixture.tierService,
		fixture.connectionEditor,
		fixture.zoneEditor,
		fixture.geometry,
		fixture.tuningFactory,
	)

	return fixture
}

func toDto(state editor_state_model.EditorState) editor_state_dto.EditorStateDto {
	return editor_state_dto.EditorStateDto{EditorState: state}
}
