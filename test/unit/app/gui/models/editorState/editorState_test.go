package editorState_test

import (
	"github.com/Tariomka/hommoe_custom_templates/app/gui/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
)

type stateValidationFunc func(editor_state_dto.EditorStateDto, bool) editor_state_dto.EditorStateValidationDto

func (this stateValidationFunc) ValidateEditorState(
	stateDto editor_state_dto.EditorStateDto,
	fixIssues bool,
) editor_state_dto.EditorStateValidationDto {
	return this(stateDto, fixIssues)
}

func newEditorState() *models.EditorState {
	return newEditorStateWithValidation(
		func(stateDto editor_state_dto.EditorStateDto, _ bool) editor_state_dto.EditorStateValidationDto {
			return editor_state_dto.EditorStateValidationDto{State: stateDto}
		},
	)
}

func newEditorStateWithValidation(validation stateValidationFunc) *models.EditorState {
	return models.NewEditorState(validation)
}
