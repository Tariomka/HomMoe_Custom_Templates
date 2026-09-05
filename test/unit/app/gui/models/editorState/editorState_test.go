package editorState_test

import (
	"github.com/Tariomka/hommoe_custom_templates/app/gui/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
)

type stateValidationFunc func(editor_state_model.EditorState, bool) editor_state_dto.EditorStateValidationDto

func (this stateValidationFunc) ValidateEditorState(
	stateDto editor_state_model.EditorState,
	fixIssues bool,
) editor_state_dto.EditorStateValidationDto {
	return this(stateDto, fixIssues)
}

func newEditorState() *models.EditorState {
	return newEditorStateWithValidation(
		func(stateDto editor_state_model.EditorState, _ bool) editor_state_dto.EditorStateValidationDto {
			return editor_state_dto.EditorStateValidationDto{State: stateDto}
		},
	)
}

func newEditorStateWithValidation(validation stateValidationFunc) *models.EditorState {
	return models.NewEditorState(validation)
}
