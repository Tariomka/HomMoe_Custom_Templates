package handlers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers/handler_interfaces"
	"github.com/Tariomka/hommoe_custom_templates/internal/mappers"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/preview_service"
)

type previewHandler struct {
	previewLayout  preview_service.IPreviewLayoutService
	templateMapper mappers.ITemplateMapper
}

func NewPreviewHandler(
	previewLayout preview_service.IPreviewLayoutService,
	templateMapper mappers.ITemplateMapper) handler_interfaces.IPreviewHandler {
	return &previewHandler{previewLayout: previewLayout, templateMapper: templateMapper}
}

func (this *previewHandler) BuildPreviewLayout(request dtos.PreviewLayoutRequestDto) (dtos.PreviewLayoutDto, error) {
	// The layout service still reads the entity shape; phase 4 moves it over.
	var template *entities.RmgTemplate
	switch {
	case request.Template != nil:
		template = new(this.templateMapper.ToEntity(*request.Template))
	case request.Zones != nil || request.Connections != nil:
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
