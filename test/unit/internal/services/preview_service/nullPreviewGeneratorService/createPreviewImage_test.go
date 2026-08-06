package nullPreviewGeneratorService_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/preview_service"
	"github.com/stretchr/testify/assert"
)

func TestWhenTemplateIsProvided_ReturnsNoImage(t *testing.T) {
	t.Parallel()
	// Arrange
	generator := preview_service.NewNullPreviewGenerator()
	template := &entities.RmgTemplate{Variants: []entities.Variant{{}}}

	// Act
	image := generator.CreatePreviewImage(template, config.TopologyRing)

	// Assert
	assert.Nil(t, image)
}

func TestWhenTemplateIsNil_ReturnsNoImage(t *testing.T) {
	t.Parallel()
	// Arrange
	generator := preview_service.NewNullPreviewGenerator()

	// Act
	image := generator.CreatePreviewImage(nil, config.TopologyRing)

	// Assert
	assert.Nil(t, image)
}
