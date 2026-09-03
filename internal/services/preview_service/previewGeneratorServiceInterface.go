package preview_service

import (
	"image"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
)

// IPreviewGeneratorService is the contract for rendering the preview image that
// is embedded next to a saved template.
type IPreviewGeneratorService interface {
	// CreatePreviewImage renders the template's preview, or nil when no image
	// can be produced.
	CreatePreviewImage(template *template_model.Template, topology config.MapTopology) *image.RGBA
}
