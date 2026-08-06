package preview_service

import (
	"image"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
)

// IPreviewGeneratorService is the contract for rendering the preview image that
// is embedded next to a saved template.
type IPreviewGeneratorService interface {
	// CreatePreviewImage renders the template's preview, or nil when no image
	// can be produced.
	CreatePreviewImage(template *entities.RmgTemplate, topology config.MapTopology) *image.RGBA
}
