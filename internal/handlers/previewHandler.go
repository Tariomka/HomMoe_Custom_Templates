package handlers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers/handler_interfaces"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/preview_service"
)

type previewHandler struct {
	previewLayout preview_service.IPreviewLayoutService
}

func NewPreviewHandler(previewLayout preview_service.IPreviewLayoutService) handler_interfaces.IPreviewHandler {
	return &previewHandler{previewLayout: previewLayout}
}

func (this *previewHandler) BuildPreviewLayout(request dtos.PreviewLayoutRequestDto) dtos.PreviewLayoutDto {
	layout := this.previewLayout.BuildPreviewLayout(request.Template, request.Topology, request.CanvasSide)
	return dtos.PreviewLayoutDto{Layout: layout}
}
