package test_helpers

import (
	"image"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/stretchr/testify/mock"
)

// PreviewGeneratorServiceMock is a testify mock of
// preview_service.IPreviewGeneratorService, used to unit-test collaborators
// without rasterizing a preview.
type PreviewGeneratorServiceMock struct {
	mock.Mock
}

func (this *PreviewGeneratorServiceMock) CreatePreviewImage(
	template *entities.RmgTemplate,
	topology config.MapTopology) *image.RGBA {
	arguments := this.Called(template, topology)
	previewImage, _ := arguments.Get(0).(*image.RGBA)
	return previewImage
}
