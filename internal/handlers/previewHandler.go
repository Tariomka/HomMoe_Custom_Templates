package handlers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/preview_service"
)

type previewHandler struct {
	previewLayout *preview_service.PreviewLayoutService
}

func newPreviewHandler(previewLayout *preview_service.PreviewLayoutService) *previewHandler {
	return &previewHandler{previewLayout: previewLayout}
}

func (this *previewHandler) BuildPreviewLayout(
	request dtos.PreviewLayoutRequestDto,
) (dtos.PreviewLayoutDto, error) {
	layout := this.previewLayout.BuildPreviewLayout(request.Template, request.Topology, request.CanvasSide)
	return dtos.PreviewLayoutDto{Layout: layout}, nil
}
