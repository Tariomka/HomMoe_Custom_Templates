package bonusHandler_test

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers/handler_interfaces"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
)

// bonusHandlerFixture holds the handler together with the service mock it was
// built from, so a test can arrange the collaborator and assert on the handler.
type bonusHandlerFixture struct {
	handler      handler_interfaces.IBonusHandler
	bonusService *test_helpers.BonusEntryServiceMock
}

func newBonusHandlerFixture() *bonusHandlerFixture {
	fixture := &bonusHandlerFixture{bonusService: &test_helpers.BonusEntryServiceMock{}}
	fixture.handler = handlers.NewBonusHandler(fixture.bonusService)

	return fixture
}
