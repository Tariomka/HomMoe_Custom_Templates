package test_helpers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/stretchr/testify/mock"
)

// TemplateHandlerMock is a testify mock of interfaces.ITemplateHandler, used
// to unit-test app/gui/drivers.State without the real generator stack.
type TemplateHandlerMock struct {
	mock.Mock
}

func (this *TemplateHandlerMock) GenerateTemplate(stateDto dtos.EditorStateDto) (dtos.TemplateLoadDto, error) {
	arguments := this.Called(stateDto)
	template, _ := arguments.Get(0).(dtos.TemplateLoadDto)
	return template, arguments.Error(1)
}

func (this *TemplateHandlerMock) UpdateTemplate(templateDto dtos.TemplateUpdateDto) (dtos.TemplateLoadDto, error) {
	arguments := this.Called(templateDto)
	template, _ := arguments.Get(0).(dtos.TemplateLoadDto)
	return template, arguments.Error(1)
}

func (this *TemplateHandlerMock) SaveTemplate(templateDto dtos.TemplateSaveDto) (string, error) {
	arguments := this.Called(templateDto)
	return arguments.String(0), arguments.Error(1)
}

func (this *TemplateHandlerMock) LoadState(path string, fixIssues bool) (*dtos.EditorStateDto, []string, error) {
	arguments := this.Called(path, fixIssues)
	state, _ := arguments.Get(0).(*dtos.EditorStateDto)
	warnings, _ := arguments.Get(1).([]string)
	return state, warnings, arguments.Error(2)
}

func (this *TemplateHandlerMock) SaveState(stateDto dtos.EditorStateSaveDto) (string, error) {
	arguments := this.Called(stateDto)
	return arguments.String(0), arguments.Error(1)
}
