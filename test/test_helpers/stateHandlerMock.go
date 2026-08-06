package test_helpers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/stretchr/testify/mock"
)

// StateHandlerMock is a testify mock of
// handler_interfaces.IStateHandler, used to unit-test collaborators without
// the real validation and persistence stack.
type StateHandlerMock struct {
	mock.Mock
}

func (this *StateHandlerMock) LoadState(path string, fixIssues bool) (*dtos.EditorStateDto, []string, error) {
	arguments := this.Called(path, fixIssues)
	state, _ := arguments.Get(0).(*dtos.EditorStateDto)
	warnings, _ := arguments.Get(1).([]string)
	return state, warnings, arguments.Error(2)
}

func (this *StateHandlerMock) SaveState(stateDto dtos.EditorStateSaveDto) (string, error) {
	arguments := this.Called(stateDto)
	return arguments.String(0), arguments.Error(1)
}

func (this *StateHandlerMock) ValidateEditorState(
	state dtos.EditorStateDto,
	fixIssues bool) dtos.EditorStateValidationDto {
	arguments := this.Called(state, fixIssues)
	validation, _ := arguments.Get(0).(dtos.EditorStateValidationDto)
	return validation
}
