package preview_service

import (
	"image"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
)

// NullPreviewGeneratorService is the no-op preview generator used when the
// asset set cannot be loaded. It keeps every consumer free of nil checks: a
// template still saves, just without a preview image.
type NullPreviewGeneratorService struct{}

func NewNullPreviewGenerator() IPreviewGeneratorService {
	return &NullPreviewGeneratorService{}
}

// CreatePreviewImage always returns nil, which callers treat as "no preview".
func (this *NullPreviewGeneratorService) CreatePreviewImage(
	_ *template_model.Template,
	_ config.MapTopology) *image.RGBA {
	return nil
}
