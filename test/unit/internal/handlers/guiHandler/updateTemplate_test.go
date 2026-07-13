package guiHandler_test

import (
	"slices"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_errors"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers"
	"github.com/Tariomka/hommoe_custom_templates/internal/mappers"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenTemplateIsNil_ReturnsProvidedTemplateInvalidError(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := handlers.NewGuiHandler()
	templateDto := dtos.TemplateUpdateDto{Template: nil}

	// Act
	_, err := handler.UpdateTemplate(templateDto)

	// Assert
	assert.ErrorIs(t, err, common_errors.ErrProvidedTemplateInvalid)
}

func TestWhenTemplateHasNoVariants_ReturnsProvidedTemplateInvalidError(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := handlers.NewGuiHandler()
	templateDto := dtos.TemplateUpdateDto{
		Template: &entities.RmgTemplate{Name: gofakeit.ProductName()},
	}

	// Act
	_, err := handler.UpdateTemplate(templateDto)

	// Assert
	assert.ErrorIs(t, err, common_errors.ErrProvidedTemplateInvalid)
}

func TestWhenGeneratedZonesAndConnectionsAreReapplied_ReturnsNoError(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
	// Arrange
	handler := handlers.NewGuiHandler()
	template := generateDefaultTemplate(t, handler)
	brokenConnections := slices.Clone(template.Variants[0].Connections)
	brokenConnections = append(brokenConnections, entities.Connection{
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
	assert.ErrorIs(t, err, common_errors.ErrZonesMissing)
}

func TestWhenUpdateSucceeds_ReturnedTemplateIsProvidedTemplateInstance(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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

func TestWhenUpdateReplacesZones_CallersTemplateZonesAreNotMutated(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := handlers.NewGuiHandler()
	template := generateDefaultTemplate(t, handler)
	originalZones := template.Variants[0].Zones
	require.NotEmpty(t, originalZones)
	templateDto := dtos.TemplateUpdateDto{
		Template:    template,
		Zones:       slices.Clone(originalZones)[:1],
		Connections: nil,
	}

	// Act
	_, err := handler.UpdateTemplate(templateDto)

	// Assert
	require.NoError(t, err)
	assert.Len(t, template.Variants[0].Zones, len(originalZones))
}

func TestWhenUpdateReplacesConnections_CallersTemplateConnectionsAreNotMutated(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := handlers.NewGuiHandler()
	template := generateDefaultTemplate(t, handler)
	originalConnections := template.Variants[0].Connections
	require.NotEmpty(t, originalConnections)
	templateDto := dtos.TemplateUpdateDto{
		Template:    template,
		Zones:       template.Variants[0].Zones,
		Connections: nil,
	}

	// Act
	_, err := handler.UpdateTemplate(templateDto)

	// Assert
	require.NoError(t, err)
	assert.Len(t, template.Variants[0].Connections, len(originalConnections))
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
