package templateHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWhenExpandIsRequested_PassesTheFreshLayoutAndStoredEditsToTheDeltaService(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newTemplateHandlerFixture()
	state := editor_state_model.NewDefaultEditorStateModel()
	state.ManualEditSettings = editor_state_model.ManualEditSettingsFromSnapshot(
		[]template_model.Zone{{Name: gofakeit.Word()}}, []template_model.Connection{{Name: gofakeit.Word()}})
	configuration := namedConfiguration()
	configuration.Topology = config.TopologySquare
	freshZones := []template_model.Zone{{Name: "Fresh"}}
	freshConnections := []template_model.Connection{{Name: "Fresh-Link"}}
	castleChanges := editor_state_model.CastleSettingChanges{PlayerCastles: true}
	arrangeExpandConfiguration(fixture, state, configuration)
	fixture.manualEditDelta.On("Expand", mock.Anything).Return(models.ManualEditExpansion{})

	// Act
	fixture.handler.ExpandManualEdits(dtos.ManualEditExpandRequestDto{
		Template:      newTemplateWithLayout(freshZones, freshConnections),
		State:         toDto(state),
		CastleChanges: castleChanges,
	})

	// Assert
	fixture.manualEditDelta.AssertCalled(t, "Expand", models.ManualEditExpandRequest{
		FreshZones:       freshZones,
		FreshConnections: freshConnections,
		ManualEdits:      state.ManualEditSettings,
		CastleChanges:    castleChanges,
		Configuration:    configuration,
		NonDeterministic: false,
	})
}

func TestWhenExpandConfigurationHasARandomTopology_FlagsTheLayoutAsNonDeterministic(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newTemplateHandlerFixture()
	state := editor_state_model.NewDefaultEditorStateModel()
	configuration := namedConfiguration()
	configuration.Topology = config.TopologyRandom
	arrangeExpandConfiguration(fixture, state, configuration)
	fixture.manualEditDelta.On("Expand", mock.Anything).Return(models.ManualEditExpansion{})

	// Act
	fixture.handler.ExpandManualEdits(dtos.ManualEditExpandRequestDto{
		Template: newTemplateWithLayout(nil, nil),
		State:    toDto(state),
	})

	// Assert
	fixture.manualEditDelta.AssertCalled(t, "Expand", models.ManualEditExpandRequest{
		Configuration:    configuration,
		NonDeterministic: true,
	})
}

func TestWhenExpandIsRequestedWithoutATemplate_PassesNoFreshLayout(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newTemplateHandlerFixture()
	state := editor_state_model.NewDefaultEditorStateModel()
	configuration := namedConfiguration()
	arrangeExpandConfiguration(fixture, state, configuration)
	fixture.manualEditDelta.On("Expand", mock.Anything).Return(models.ManualEditExpansion{})

	// Act
	fixture.handler.ExpandManualEdits(dtos.ManualEditExpandRequestDto{State: toDto(state)})

	// Assert
	fixture.manualEditDelta.AssertCalled(t, "Expand", models.ManualEditExpandRequest{
		Configuration:    configuration,
		NonDeterministic: false,
	})
}

func TestWhenExpandIsRequestedWithAVariantlessTemplate_PassesNoFreshLayout(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newTemplateHandlerFixture()
	state := editor_state_model.NewDefaultEditorStateModel()
	configuration := namedConfiguration()
	arrangeExpandConfiguration(fixture, state, configuration)
	fixture.manualEditDelta.On("Expand", mock.Anything).Return(models.ManualEditExpansion{})

	// Act
	fixture.handler.ExpandManualEdits(dtos.ManualEditExpandRequestDto{
		Template: &template_model.Template{},
		State:    toDto(state),
	})

	// Assert
	fixture.manualEditDelta.AssertCalled(t, "Expand", models.ManualEditExpandRequest{
		Configuration:    configuration,
		NonDeterministic: false,
	})
}

func TestWhenExpandSucceeds_ReturnsTheDeltaServiceExpansion(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newTemplateHandlerFixture()
	state := editor_state_model.NewDefaultEditorStateModel()
	expansion := models.ManualEditExpansion{
		Zones:       []template_model.Zone{{Name: gofakeit.Word()}},
		Connections: []template_model.Connection{{Name: gofakeit.Word()}},
		ManualEdits: editor_state_model.ManualEditSettingsFromSnapshot([]template_model.Zone{{Name: "Z"}}, nil),
		Changed:     true,
		Warnings:    []string{gofakeit.Sentence(3)},
	}
	arrangeExpandConfiguration(fixture, state, namedConfiguration())
	fixture.manualEditDelta.On("Expand", mock.Anything).Return(expansion)

	// Act
	result := fixture.handler.ExpandManualEdits(dtos.ManualEditExpandRequestDto{
		Template: newTemplateWithLayout(nil, nil),
		State:    toDto(state),
	})

	// Assert
	assert.Equal(t, dtos.ManualEditExpansionDto(expansion), result)
}

func TestWhenExpandIsRequested_MapsTheValidatedSettings(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newTemplateHandlerFixture()
	state := editor_state_model.NewDefaultEditorStateModel()
	validated := state.Clone()
	validated.ResourceDensityPercent = state.ResourceDensityPercent + 1
	fixture.stateHandler.On("ValidateEditorState", state, true).
		Return(editor_state_dto.EditorStateValidationDto{State: validated})
	fixture.mapper.On("FromEditorState", validated).Return(namedConfiguration())
	fixture.manualEditDelta.On("Expand", mock.Anything).Return(models.ManualEditExpansion{})

	// Act
	fixture.handler.ExpandManualEdits(dtos.ManualEditExpandRequestDto{State: toDto(state)})

	// Assert
	fixture.mapper.AssertCalled(t, "FromEditorState", validated)
}

// arrangeExpandConfiguration makes validation keep the state as it is and
// maps it to the given configuration.
func arrangeExpandConfiguration(
	fixture *templateHandlerFixture,
	state editor_state_model.EditorState,
	configuration *config.GeneratorConfig) {
	fixture.stateHandler.On("ValidateEditorState", state, true).
		Return(editor_state_dto.EditorStateValidationDto{State: state})
	fixture.mapper.On("FromEditorState", state).Return(configuration)
}
