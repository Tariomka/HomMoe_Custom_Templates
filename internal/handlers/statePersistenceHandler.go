package handlers

import (
	"strings"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_errors"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/file_service"
)

type statePersistenceHandler struct {
	fileService     *file_service.FileService
	stateValidation *stateValidationHandler
}

func newStatePersistenceHandler(
	fileService *file_service.FileService,
	stateValidation *stateValidationHandler,
) *statePersistenceHandler {
	return &statePersistenceHandler{
		fileService:     fileService,
		stateValidation: stateValidation,
	}
}

func (this *statePersistenceHandler) LoadState(
	path string,
	fixIssues bool,
) (*dtos.EditorStateDto, []string, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return nil, nil, common_errors.ErrNoOutputPath
	}

	loaded, err := this.fileService.LoadSettingsFile(path)
	if err != nil {
		return nil, nil, err
	}

	validation := this.stateValidation.validateEditorState(*loaded, fixIssues)
	return &validation.State, validation.Warnings, nil
}

func (this *statePersistenceHandler) SaveState(stateDto dtos.EditorStateSaveDto) (string, error) {
	if stateDto.State == nil {
		return "", common_errors.ErrNothingToSave
	}

	outputPath := strings.TrimSpace(stateDto.OutputPath)
	if outputPath == "" {
		return "", common_errors.ErrNoOutputPath
	}

	err := this.fileService.SaveSettings(outputPath, stateDto.State)
	if err != nil {
		return "", err
	}

	return outputPath, nil
}
