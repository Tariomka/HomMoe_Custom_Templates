package handlers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/validators"
)

type stateValidationHandler struct {
	editorValidator *validators.EditorStateValidator
}

func newStateValidationHandler(
	editorValidator *validators.EditorStateValidator,
) *stateValidationHandler {
	return &stateValidationHandler{editorValidator: editorValidator}
}

func (this *stateValidationHandler) validateEditorState(
	stateDto dtos.EditorStateDto,
	fixIssues bool,
) dtos.EditorStateValidationDto {
	issues := this.editorValidator.Validate(&stateDto)
	warnings := make([]string, 0, len(issues))
	for _, issue := range issues {
		if fixIssues {
			issue.Fix(&stateDto)
		}
		warnings = append(warnings, issue.Message)
	}

	return dtos.EditorStateValidationDto{State: stateDto, Warnings: warnings}
}
