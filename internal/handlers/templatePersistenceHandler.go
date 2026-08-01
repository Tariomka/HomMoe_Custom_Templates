package handlers

import (
	"strings"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_errors"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/file_service"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/preview_service"
)

type templatePersistenceHandler struct {
	fileService      *file_service.FileService
	previewGenerator *preview_service.PreviewGeneratorService
}

func newTemplatePersistenceHandler(
	fileService *file_service.FileService,
	previewGenerator *preview_service.PreviewGeneratorService,
) *templatePersistenceHandler {
	return &templatePersistenceHandler{
		fileService:      fileService,
		previewGenerator: previewGenerator,
	}
}

func (this *templatePersistenceHandler) SaveTemplate(templateDto dtos.TemplateSaveDto) (string, error) {
	if templateDto.Template == nil {
		return "", common_errors.ErrNothingToSave
	}

	outputPath := strings.TrimSpace(templateDto.OutputPath)
	if outputPath == "" {
		return "", common_errors.ErrNoOutputPath
	}

	out, err := this.fileService.SaveTemplate(outputPath, templateDto.Template)
	if err != nil {
		return "", err
	}

	if this.previewGenerator != nil {
		previewImage := this.previewGenerator.CreatePreviewImage(templateDto.Template, templateDto.Topology)
		_, err = this.fileService.SavePreviewImage(outputPath, previewImage, templateDto.Template.Name)
		if err != nil {
			return out, err
		}
	}

	return out, nil
}
