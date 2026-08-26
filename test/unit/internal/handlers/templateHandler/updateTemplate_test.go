package templateHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_errors"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWhenUpdatedTemplateIsMissing_ReturnsProvidedTemplateInvalidError(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newTemplateHandlerFixture()

	// Act
	_, err := fixture.handler.UpdateTemplate(dtos.TemplateUpdateDto{})

	// Assert
	assert.ErrorIs(t, err, common_errors.ErrProvidedTemplateInvalid)
}

func TestWhenUpdatedTemplateHasNoVariants_ReturnsProvidedTemplateInvalidError(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newTemplateHandlerFixture()

	// Act
	_, err := fixture.handler.UpdateTemplate(dtos.TemplateUpdateDto{Template: &entities.RmgTemplate{}})

	// Assert
	assert.ErrorIs(t, err, common_errors.ErrProvidedTemplateInvalid)
}

func TestWhenTemplateIsUpdated_ReplacesTheFirstVariantsZones(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newTemplateHandlerFixture()
	zones := []entities.Zone{{Name: gofakeit.Word()}}
	arrangeUpdateCollaborators(fixture, false)

	// Act
	loadDto, _ := fixture.handler.UpdateTemplate(dtos.TemplateUpdateDto{
		Template: singleVariantTemplate(),
		Zones:    zones,
	})

	// Assert
	assert.Equal(t, zones, loadDto.Template.Variants[0].Zones)
}

func TestWhenTemplateIsUpdated_ReplacesTheFirstVariantsConnections(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newTemplateHandlerFixture()
	connections := []entities.Connection{{Name: gofakeit.Word()}}
	arrangeUpdateCollaborators(fixture, false)

	// Act
	loadDto, _ := fixture.handler.UpdateTemplate(dtos.TemplateUpdateDto{
		Template:    singleVariantTemplate(),
		Connections: connections,
	})

	// Assert
	assert.Equal(t, connections, loadDto.Template.Variants[0].Connections)
}

func TestWhenTemplateIsUpdated_LeavesTheSourceTemplateUntouched(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newTemplateHandlerFixture()
	source := singleVariantTemplate()
	originalZones := source.Variants[0].Zones
	arrangeUpdateCollaborators(fixture, false)

	// Act
	_, _ = fixture.handler.UpdateTemplate(dtos.TemplateUpdateDto{
		Template: source,
		Zones:    []entities.Zone{{Name: gofakeit.Word()}},
	})

	// Assert
	assert.Equal(t, originalZones, source.Variants[0].Zones)
}

func TestWhenTemplateIsUpdated_RebuildsTheZoneConnectionRoads(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newTemplateHandlerFixture()
	zones := []entities.Zone{{Name: gofakeit.Word()}}
	connections := []entities.Connection{{Name: gofakeit.Word()}}
	arrangeUpdateCollaborators(fixture, false)

	// Act
	_, _ = fixture.handler.UpdateTemplate(dtos.TemplateUpdateDto{
		Template:    singleVariantTemplate(),
		Zones:       zones,
		Connections: connections,
	})

	// Assert
	fixture.zoneEditor.AssertCalled(t, "RebuildZoneConnectionRoads", zones, connections)
}

func TestWhenNoEditorStateIsSupplied_KeepsTheExistingMandatoryContent(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newTemplateHandlerFixture()
	arrangeUpdateCollaborators(fixture, false)

	// Act
	_, _ = fixture.handler.UpdateTemplate(dtos.TemplateUpdateDto{Template: singleVariantTemplate()})

	// Assert
	fixture.contentProvider.AssertNotCalled(t, "CreateContentsForZones", mock.Anything, mock.Anything)
}

func TestWhenEditorStateIsSupplied_RebuildsTheMandatoryContentFromTheFinalZones(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newTemplateHandlerFixture()
	state := editor_state_model.NewDefaultEditorStateModel()
	zones := []entities.Zone{{Name: gofakeit.Word()}}
	expected := []entities.MandatoryContent{{Name: gofakeit.Word()}}
	configuration := namedConfiguration()
	arrangeUpdateCollaborators(fixture, false)
	fixture.mapper.On("FromEditorState", state).Return(configuration)
	fixture.contentProvider.On("CreateContentsForZones", *configuration, zones).Return(expected)

	// Act
	loadDto, _ := fixture.handler.UpdateTemplate(dtos.TemplateUpdateDto{
		Template:    singleVariantTemplate(),
		Zones:       zones,
		EditorState: &editor_state_dto.EditorStateDto{EditorState: state},
	})

	// Assert
	assert.Equal(t, expected, loadDto.Template.MandatoryContent)
}

func TestWhenTheUpdatedGraphHasErrors_ReturnsZonesMissingError(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newTemplateHandlerFixture()
	arrangeUpdateCollaborators(fixture, true)

	// Act
	_, err := fixture.handler.UpdateTemplate(dtos.TemplateUpdateDto{Template: singleVariantTemplate()})

	// Assert
	assert.ErrorIs(t, err, common_errors.ErrZonesMissing)
}

func TestWhenTheUpdatedGraphIsSound_ReturnsNoError(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newTemplateHandlerFixture()
	arrangeUpdateCollaborators(fixture, false)

	// Act
	_, err := fixture.handler.UpdateTemplate(dtos.TemplateUpdateDto{Template: singleVariantTemplate()})

	// Assert
	assert.NoError(t, err)
}

// arrangeUpdateCollaborators stubs the two collaborators every UpdateTemplate
// call reaches, with the requested graph-error verdict.
func arrangeUpdateCollaborators(fixture *templateHandlerFixture, hasErrors bool) {
	fixture.zoneEditor.On("RebuildZoneConnectionRoads", mock.Anything, mock.Anything).Return()
	fixture.connectionEditor.On("ComputeHasErrors", mock.Anything, mock.Anything).Return(hasErrors)
}

func singleVariantTemplate() *entities.RmgTemplate {
	return &entities.RmgTemplate{
		Name:     gofakeit.Word(),
		Variants: []entities.Variant{{Zones: []entities.Zone{{Name: gofakeit.Word()}}}},
	}
}
