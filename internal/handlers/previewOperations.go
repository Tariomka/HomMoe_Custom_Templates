package handlers

import "github.com/Tariomka/hommoe_custom_templates/internal/dtos"

type PreviewOperations interface {
	BuildPreviewLayout(request dtos.PreviewLayoutRequestDto) (dtos.PreviewLayoutDto, error)
}
