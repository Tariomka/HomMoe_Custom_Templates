package composition

import (
	"log/slog"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/preview_service"
)

// providePreviewGenerator keeps the injector error-free: a missing or unreadable
// asset set degrades to the null generator ("no preview images") instead of
// failing construction, so no consumer has to nil-check the generator.
func providePreviewGenerator(
	layoutService preview_service.IPreviewLayoutService) preview_service.IPreviewGeneratorService {
	previewGenerator, err := preview_service.NewPreviewGenerator(layoutService)
	if err != nil {
		slog.Error(
			"Preview Generator failed to initialize, preview images will not be generated",
			slog.String("error", err.Error()))
		return preview_service.NewNullPreviewGenerator()
	}

	return previewGenerator
}
