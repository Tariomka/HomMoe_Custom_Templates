package handlers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers/handler_interfaces"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/preview_service"
)

type previewHandler struct {
	previewLayout *preview_service.PreviewLayoutService
}

func NewPreviewHandler(
	previewLayout *preview_service.PreviewLayoutService) handler_interfaces.IPreviewHandler {
	return &previewHandler{previewLayout: previewLayout}
}

func (this *previewHandler) BuildPreviewLayout(request dtos.PreviewLayoutRequestDto) (dtos.PreviewLayoutDto, error) {
	template := request.Template
	if template == nil && (request.Zones != nil || request.Connections != nil) {
		template = &entities.RmgTemplate{
			Variants: []entities.Variant{variant_content.NewVariantBuilder().
				WithZones(request.Zones...).
				WithConnections(request.Connections...).
				Build()},
		}
	}
	layout := this.previewLayout.BuildPreviewLayout(template, request.Topology, request.CanvasSide)
	return dtos.PreviewLayoutDto{Layout: layout}, nil
}
