package zoneEditorHandler_test

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers/handler_interfaces"
	"github.com/Tariomka/hommoe_custom_templates/internal/mappers"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
)

// zoneEditorHandlerFixture holds the handler together with every mock it was
// built from, so a test can arrange one collaborator and assert on another.
type zoneEditorHandlerFixture struct {
	handler           handler_interfaces.IZoneEditorHandler
	mapper            *test_helpers.GeneratorConfigMapperMock
	editorStateMapper mappers.IEditorStateMapper
	zoneClassifier    *test_helpers.ZoneClassifierMock
	connectionEditor  *test_helpers.ConnectionEditorServiceMock
	zoneEditor        *test_helpers.ZoneEditorServiceMock
	geometry          *test_helpers.ZoneEditorGeometryServiceMock
	tuningFactory     *test_helpers.GenerationTuningFactoryMock
}

func newZoneEditorHandlerFixture() *zoneEditorHandlerFixture {
	fixture := &zoneEditorHandlerFixture{
		mapper:            &test_helpers.GeneratorConfigMapperMock{},
		editorStateMapper: mappers.NewEditorStateMapper(),
		zoneClassifier:    &test_helpers.ZoneClassifierMock{},
		connectionEditor:  &test_helpers.ConnectionEditorServiceMock{},
		zoneEditor:        &test_helpers.ZoneEditorServiceMock{},
		geometry:          &test_helpers.ZoneEditorGeometryServiceMock{},
		tuningFactory:     &test_helpers.GenerationTuningFactoryMock{},
	}

	fixture.handler = handlers.NewZoneEditorHandler(
		fixture.mapper,
		fixture.editorStateMapper,
		fixture.zoneClassifier,
		fixture.connectionEditor,
		fixture.zoneEditor,
		fixture.geometry,
		fixture.tuningFactory,
	)

	return fixture
}
