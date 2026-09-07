package preview_service

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/preview"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
)

// IPreviewLayoutService is the contract for computing the preview geometry of a
// generated template.
type IPreviewLayoutService interface {
	// BuildPreviewLayout computes zone positions, radius and connections for a
	// preview canvas of the given side length.
	BuildPreviewLayout(template *template_model.Template, topology config.MapTopology, side float64) preview.Layout
}
