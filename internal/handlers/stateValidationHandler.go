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
	if fixIssues {
		normalizeInactiveNeutralCounts(&stateDto)
	}

	return dtos.EditorStateValidationDto{State: stateDto, Warnings: warnings}
}

func normalizeInactiveNeutralCounts(stateDto *dtos.EditorStateDto) {
	if stateDto.AdvancedMode {
		stateDto.NeutralZoneCount = 0
		return
	}

	stateDto.NeutralLowestNoCastleCount = 0
	stateDto.NeutralLowestCastleCount = 0
	stateDto.NeutralLowNoCastleCount = 0
	stateDto.NeutralLowCastleCount = 0
	stateDto.NeutralMediumNoCastleCount = 0
	stateDto.NeutralMediumCastleCount = 0
	stateDto.NeutralHighNoCastleCount = 0
	stateDto.NeutralHighCastleCount = 0
}
