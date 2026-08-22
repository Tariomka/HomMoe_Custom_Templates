package test_helpers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/validators"
	"github.com/stretchr/testify/mock"
)

// EditorStateValidatorMock is a testify mock of
// validators.IEditorStateValidator, used to unit-test collaborators with a
// controlled set of validation issues.
type EditorStateValidatorMock struct {
	mock.Mock
}

func (this *EditorStateValidatorMock) Validate(
	state *editor_state_model.EditorState,
) []validators.ValidationIssue {
	arguments := this.Called(state)
	issues, _ := arguments.Get(0).([]validators.ValidationIssue)
	return issues
}
