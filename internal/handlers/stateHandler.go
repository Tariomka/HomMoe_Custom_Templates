package handlers

import (
	"strings"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_errors"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers/handler_interfaces"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/file_service"
	"github.com/Tariomka/hommoe_custom_templates/internal/validators"
)

type stateHandler struct {
	fileService     file_service.IFileService
	editorValidator validators.IEditorStateValidator
}

func NewStateHandler(
	fileService file_service.IFileService,
	editorValidator validators.IEditorStateValidator) handler_interfaces.IStateHandler {
	return &stateHandler{
		fileService:     fileService,
		editorValidator: editorValidator,
	}
}

func (this *stateHandler) LoadState(path string, fixIssues bool) (*editor_state_dto.EditorStateDto, []string, error) {
	// TODO: should just return EditorStateValidationDto, instead of EditorStateDto + warnings
	path = strings.TrimSpace(path)
	if path == "" {
		return nil, nil, common_errors.ErrNoOutputPath
	}

	loaded, err := this.fileService.LoadSettingsFile(path)
	if err != nil {
		return nil, nil, err
	}

	validation := this.ValidateEditorState(*loaded, fixIssues)
	return &editor_state_dto.EditorStateDto{EditorState: validation.State}, validation.Warnings, nil
}

func (this *stateHandler) SaveState(stateDto editor_state_dto.EditorStateSaveDto) (string, error) {
	if stateDto.State == nil {
		return "", common_errors.ErrNothingToSave
	}

	outputPath := strings.TrimSpace(stateDto.OutputPath)
	if outputPath == "" {
		return "", common_errors.ErrNoOutputPath
	}

	return this.fileService.SaveSettings(outputPath, &stateDto.State.EditorState)
}

func (this *stateHandler) ValidateEditorState(
	state editor_state_model.EditorState,
	fixIssues bool) editor_state_dto.EditorStateValidationDto {
	// Cloned on entry so the fixes below never write through to the caller's
	// slices, and so the returned state does not alias them either.
	state = state.Clone()
	issues := this.editorValidator.Validate(&state)
	warnings := make([]string, 0, len(issues))
	for _, issue := range issues {
		if fixIssues {
			issue.Fix(&state)
		}
		warnings = append(warnings, issue.Message)
	}
	if fixIssues {
		normalizeInactiveNeutralCounts(&state)
	}

	return editor_state_dto.EditorStateValidationDto{State: state, Warnings: warnings}
}

func normalizeInactiveNeutralCounts(state *editor_state_model.EditorState) {
	if state.AdvancedMode {
		state.NeutralZoneCount = 0
		return
	}

	state.NeutralLowestNoCastleCount = 0
	state.NeutralLowestCastleCount = 0
	state.NeutralLowNoCastleCount = 0
	state.NeutralLowCastleCount = 0
	state.NeutralMediumNoCastleCount = 0
	state.NeutralMediumCastleCount = 0
	state.NeutralHighNoCastleCount = 0
	state.NeutralHighCastleCount = 0
}
