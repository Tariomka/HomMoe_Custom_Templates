package test_helpers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/stretchr/testify/mock"
)

// StateHandlerMock is a testify mock of
// handler_interfaces.IStateHandler, used to unit-test collaborators without
// the real validation and persistence stack.
type StateHandlerMock struct {
	mock.Mock
}

func (this *StateHandlerMock) LoadState(
	path string,
	fixIssues bool,
) (*editor_state_model.EditorState, []string, error) {
	arguments := this.Called(path, fixIssues)
	state, _ := arguments.Get(0).(*editor_state_model.EditorState)
	warnings, _ := arguments.Get(1).([]string)
	return state, warnings, arguments.Error(2)
}

func (this *StateHandlerMock) SaveState(stateDto editor_state_dto.EditorStateSaveDto) (string, error) {
	arguments := this.Called(stateDto)
	return arguments.String(0), arguments.Error(1)
}

func (this *StateHandlerMock) ValidateEditorState(
	state editor_state_model.EditorState,
	fixIssues bool) editor_state_dto.EditorStateValidationDto {
	arguments := this.Called(state, fixIssues)
	validation, _ := arguments.Get(0).(editor_state_dto.EditorStateValidationDto)
	return validation
}
