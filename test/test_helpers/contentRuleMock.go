package test_helpers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/stretchr/testify/mock"
)

// ContentRuleMock is a testify mock of content_rules.IContentRule, used to
// unit-test collaborators without a concrete rule implementation.
type ContentRuleMock struct {
	mock.Mock
}

func (this *ContentRuleMock) Name() string {
	arguments := this.Called()
	return arguments.String(0)
}

func (this *ContentRuleMock) Description() string {
	arguments := this.Called()
	return arguments.String(0)
}

func (this *ContentRuleMock) Marker() string {
	arguments := this.Called()
	return arguments.String(0)
}

func (this *ContentRuleMock) DisplayText() string {
	arguments := this.Called()
	return arguments.String(0)
}

func (this *ContentRuleMock) SerializeToRowSave() editor_state_model.ContentRuleRow {
	arguments := this.Called()
	row, _ := arguments.Get(0).(editor_state_model.ContentRuleRow)
	return row
}

func (this *ContentRuleMock) Apply(item *entities.MandatoryContentItem) {
	this.Called(item)
}
