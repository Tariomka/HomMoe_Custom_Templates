package guiHandler_test

import (
	"slices"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_errors"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenTemplateIsNil_ReturnsProvidedTemplateInvalidError(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	templateDto := dtos.TemplateUpdateDto{Template: nil}

	// Act
	_, err := handler.UpdateTemplate(templateDto)

	// Assert
	assert.ErrorIs(t, err, common_errors.ErrProvidedTemplateInvalid)
}

func TestWhenTemplateHasNoVariants_ReturnsProvidedTemplateInvalidError(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
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
	handler := newProductionGuiHandler()
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
	handler := newProductionGuiHandler()
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
	handler := newProductionGuiHandler()
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
	handler := newProductionGuiHandler()
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

func TestWhenEditorStateIsProvided_MandatoryContentMatchesMappedConfiguration(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	template := generateDefaultTemplate(t, handler)
	template.MandatoryContent = nil
	editorState := editor_state_model.NewDefaultEditorStateModel()
	configuration := test_helpers.NewConfigMapper().FromEditorState(editorState)
	expectedContent := newMandatoryContentProvider().CreateContentsForZones(
		*configuration,
		template.Variants[0].Zones,
	)
	templateDto := dtos.TemplateUpdateDto{
		Template:    template,
		Zones:       template.Variants[0].Zones,
		Connections: template.Variants[0].Connections,
		EditorState: &editorState,
	}

	// Act
	loadDto, err := handler.UpdateTemplate(templateDto)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedContent, loadDto.Template.MandatoryContent)
}

func TestWhenZoneWasPromotedToHighTier_UsesHighTierEditorRows(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	zones := []entities.Zone{{
		Name:               "Neutral-G",
		Layout:             registry.GetLayoutValues().TreasureZone,
		GuardedContentPool: []string{"classic_template_pool_random_t4_item"},
		MainObjects:        []entities.MainObject{{Type: "City"}},
	}}
	template := &entities.RmgTemplate{Variants: []entities.Variant{{Zones: zones}}}
	editorState := editor_state_model.NewDefaultEditorStateModel()
	editorState.SpawnRemoteFootholds = false
	editorState.MediumNeutralContentRows = []models.ZoneContentRow{{Sid: "medium_only", Count: 1}}
	editorState.HighNeutralContentRows = []models.ZoneContentRow{{Sid: "high_only", Count: 1}}
	templateDto := dtos.TemplateUpdateDto{
		Template:    template,
		Zones:       zones,
		Connections: nil,
		EditorState: &editorState,
	}

	// Act
	loadDto, err := handler.UpdateTemplate(templateDto)

	// Assert
	require.NoError(t, err)
	require.Len(t, loadDto.Template.MandatoryContent, 1)
	require.Len(t, loadDto.Template.MandatoryContent[0].Content, 1)
	assert.Equal(t, "high_only", loadDto.Template.MandatoryContent[0].Content[0].SID)
}

func TestWhenEditorStateIsNil_MandatoryContentIsLeftUntouched(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
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
	handler := newProductionGuiHandler()
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
	handler := newProductionGuiHandler()
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
