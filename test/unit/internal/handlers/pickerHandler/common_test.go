package pickerHandler_test

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers/handler_interfaces"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
)

// pickerHandlerFixture holds the handler together with the picker service it
// was built from, so a test can arrange one and assert on the other.
type pickerHandlerFixture struct {
	handler handler_interfaces.IPickerHandler
	service *test_helpers.PickerEntryServiceMock
}

func newPickerHandlerFixture() *pickerHandlerFixture {
	fixture := &pickerHandlerFixture{service: &test_helpers.PickerEntryServiceMock{}}
	fixture.handler = handlers.NewPickerHandler(fixture.service)

	return fixture
}
