package zoneContentHandler_test

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers/handler_interfaces"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
)

// zoneContentHandlerFixture holds the handler together with the collaborators it
// was built from, so a test can arrange one and assert on the other.
type zoneContentHandlerFixture struct {
	handler       handler_interfaces.IZoneContentHandler
	contentRules  *test_helpers.TemplateHandlerMock
	contentEditor *test_helpers.ZoneContentEditorServiceMock
}

func newZoneContentHandlerFixture() *zoneContentHandlerFixture {
	fixture := &zoneContentHandlerFixture{
		contentRules:  &test_helpers.TemplateHandlerMock{},
		contentEditor: &test_helpers.ZoneContentEditorServiceMock{},
	}
	fixture.handler = handlers.NewZoneContentHandler(fixture.contentRules, fixture.contentEditor)

	return fixture
}
