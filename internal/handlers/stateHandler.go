package handlers

import (
	"strings"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_errors"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers/handler_interfaces"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/file_service"
	"github.com/Tariomka/hommoe_custom_templates/internal/validators"
)

type stateHandler struct {
	fileService     *file_service.FileService
	editorValidator *validators.EditorStateValidator
}

func NewStateHandler(
	fileService *file_service.FileService,
	editorValidator *validators.EditorStateValidator) handler_interfaces.IStateHandler {
	return &stateHandler{
		fileService:     fileService,
		editorValidator: editorValidator,
	}
}

func (this *stateHandler) LoadState(path string, fixIssues bool) (*dtos.EditorStateDto, []string, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return nil, nil, common_errors.ErrNoOutputPath
	}

	loaded, err := this.fileService.LoadSettingsFile(path)
	if err != nil {
		return nil, nil, err
	}

	validation := this.ValidateEditorState(*loaded, fixIssues)
	return &validation.State, validation.Warnings, nil
}

func (this *stateHandler) SaveState(stateDto dtos.EditorStateSaveDto) (string, error) {
	if stateDto.State == nil {
		return "", common_errors.ErrNothingToSave
	}

	outputPath := strings.TrimSpace(stateDto.OutputPath)
	if outputPath == "" {
		return "", common_errors.ErrNoOutputPath
	}

	return this.fileService.SaveSettings(outputPath, stateDto.State)
}

func (this *stateHandler) ValidateEditorState(
	stateDto dtos.EditorStateDto,
	fixIssues bool) dtos.EditorStateValidationDto {
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
