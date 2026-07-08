package guiHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/common"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers"
	"github.com/Tariomka/hommoe_custom_templates/internal/mappers"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenTemplateIsNil_ReturnsProvidedTemplateInvalidError(t *testing.T) {
	// Arrange
	handler := handlers.NewGuiHandler()
	templateDto := dtos.TemplateUpdateDto{Template: nil}

	// Act
	_, err := handler.UpdateTemplate(templateDto)

	// Assert
	assert.ErrorIs(t, err, common.ErrProvidedTemplateInvalid)
}

func TestWhenTemplateHasNoVariants_ReturnsProvidedTemplateInvalidError(t *testing.T) {
	// Arrange
	handler := handlers.NewGuiHandler()
	templateDto := dtos.TemplateUpdateDto{
		Template: &entities.RmgTemplate{Name: gofakeit.ProductName()},
	}

	// Act
	_, err := handler.UpdateTemplate(templateDto)

	// Assert
	assert.ErrorIs(t, err, common.ErrProvidedTemplateInvalid)
}

func TestWhenGeneratedZonesAndConnectionsAreReapplied_ReturnsNoError(t *testing.T) {
	// Arrange
	handler := handlers.NewGuiHandler()
	template := generateDefaultTemplate(t, handler)
	templateDto := dtos.TemplateUpdateDto{
		Template:    template,
		Zones:       template.Variants[0].Zones,
		Connections: template.Variants[0].Connections,
	}

	// Act
	_, err := handler.UpdateTemplate(templateDto)

	// Assert
	assert.NoError(t, err)
}

func TestWhenConnectionReferencesUnknownZone_ReturnsZonesMissingError(t *testing.T) {
	// Arrange
	handler := handlers.NewGuiHandler()
	template := generateDefaultTemplate(t, handler)
	brokenConnections := append(template.Variants[0].Connections, entities.Connection{
		Name: gofakeit.ProductName(),
		From: "No-Such-Zone",
		To:   "Another-Missing-Zone",
	})
	templateDto := dtos.TemplateUpdateDto{
		Template:    template,
		Zones:       template.Variants[0].Zones,
		Connections: brokenConnections,
	}

	// Act
	_, err := handler.UpdateTemplate(templateDto)

	// Assert
	assert.ErrorIs(t, err, common.ErrZonesMissing)
}

func TestWhenUpdateSucceeds_ReturnedTemplateIsProvidedTemplateInstance(t *testing.T) {
	// Arrange
	handler := handlers.NewGuiHandler()
	template := generateDefaultTemplate(t, handler)
	templateDto := dtos.TemplateUpdateDto{
		Template:    template,
		Zones:       template.Variants[0].Zones,
		Connections: template.Variants[0].Connections,
	}

	// Act
	loadDto, err := handler.UpdateTemplate(templateDto)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, *template, *loadDto.Template)
}

func TestWhenOnlySubsetOfZonesIsProvided_VariantZonesAreReplaced(t *testing.T) {
	// Arrange
	handler := handlers.NewGuiHandler()
	template := generateDefaultTemplate(t, handler)
	require.NotEmpty(t, template.Variants[0].Zones)
	templateDto := dtos.TemplateUpdateDto{
		Template:    template,
		Zones:       template.Variants[0].Zones[:1],
		Connections: nil,
	}

	// Act
	loadDto, err := handler.UpdateTemplate(templateDto)

	// Assert
	require.NoError(t, err)
	assert.Len(t, loadDto.Template.Variants[0].Zones, 1)
}

func TestWhenConfigIsProvided_MandatoryContentIsRebuiltFromZones(t *testing.T) {
	// Arrange
	handler := handlers.NewGuiHandler()
	template := generateDefaultTemplate(t, handler)
	template.MandatoryContent = nil
	configuration := mappers.NewConfigMapper().FromEditorState(dtos.NewDefaultEditorStateDto())
	templateDto := dtos.TemplateUpdateDto{
		Template:    template,
		Zones:       template.Variants[0].Zones,
		Connections: template.Variants[0].Connections,
		Config:      configuration,
	}

	// Act
	loadDto, err := handler.UpdateTemplate(templateDto)

	// Assert
	require.NoError(t, err)
	assert.NotEmpty(t, loadDto.Template.MandatoryContent)
}

func TestWhenConfigIsNil_MandatoryContentIsLeftUntouched(t *testing.T) {
	// Arrange
	handler := handlers.NewGuiHandler()
	template := generateDefaultTemplate(t, handler)
	template.MandatoryContent = nil
	templateDto := dtos.TemplateUpdateDto{
		Template:    template,
		Zones:       template.Variants[0].Zones,
		Connections: template.Variants[0].Connections,
	}

	// Act
	loadDto, err := handler.UpdateTemplate(templateDto)

	// Assert
	require.NoError(t, err)
	assert.Nil(t, loadDto.Template.MandatoryContent)
}

// generateDefaultTemplate produces a real generated template from the default
// editor state so update tests operate on realistic zones and connections.
func generateDefaultTemplate(t *testing.T, handler *handlers.GUIHandler) *entities.RmgTemplate {
	t.Helper()

	loadDto, err := handler.GenerateTemplate(dtos.NewDefaultEditorStateDto())
	require.NoError(t, err)
	require.NotNil(t, loadDto.Template)
	require.NotEmpty(t, loadDto.Template.Variants)

	return loadDto.Template
}
