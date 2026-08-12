package editorState_test

import (
	"github.com/Tariomka/hommoe_custom_templates/app/gui/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
)

type stateValidationFunc func(dtos.EditorStateDto, bool) dtos.EditorStateValidationDto

func (this stateValidationFunc) ValidateEditorState(
	stateDto dtos.EditorStateDto,
	fixIssues bool,
) dtos.EditorStateValidationDto {
	return this(stateDto, fixIssues)
}

func newEditorState() *models.EditorState {
	return newEditorStateWithValidation(func(stateDto dtos.EditorStateDto, _ bool) dtos.EditorStateValidationDto {
		return dtos.EditorStateValidationDto{State: stateDto}
	})
}

func newEditorStateWithValidation(validation stateValidationFunc) *models.EditorState {
	return models.NewEditorState(validation)
}
