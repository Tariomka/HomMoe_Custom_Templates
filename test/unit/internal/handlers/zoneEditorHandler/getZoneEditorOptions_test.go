package zoneEditorHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWhenEditorOptionsAreRequested_ReturnsTheStatesTopology(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newZoneEditorHandlerFixture()
	state := editor_state_model.NewDefaultEditorStateModel()
	state.Topology = config.TopologyChain
	fixture.mapper.On("FromEditorState", state).Return(config.NewGeneratorConfig())
	fixture.tuningFactory.On("Create", mock.Anything, mock.Anything).Return(models.GenerationTuning{})

	// Act
	options := fixture.handler.GetZoneEditorOptions(toDto(state), gofakeit.IntRange(1, 20))

	// Assert
	assert.Equal(t, config.TopologyChain, options.Topology)
}

func TestWhenEditorOptionsAreRequested_ReturnsTheStatesRoadFlag(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newZoneEditorHandlerFixture()
	state := editor_state_model.NewDefaultEditorStateModel()
	state.GenerateRoads = true
	fixture.mapper.On("FromEditorState", state).Return(config.NewGeneratorConfig())
	fixture.tuningFactory.On("Create", mock.Anything, mock.Anything).Return(models.GenerationTuning{})

	// Act
	options := fixture.handler.GetZoneEditorOptions(toDto(state), gofakeit.IntRange(1, 20))

	// Assert
	assert.True(t, options.GenerateRoads)
}

func TestWhenEditorOptionsAreRequested_ReturnsTheTuningForTheZoneCount(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newZoneEditorHandlerFixture()
	state := editor_state_model.NewDefaultEditorStateModel()
	totalZoneCount := gofakeit.IntRange(1, 20)
	configuration := config.NewGeneratorConfig()
	expected := models.GenerationTuning{ContentScale: gofakeit.Float64Range(0.5, 2)}
	fixture.mapper.On("FromEditorState", state).Return(configuration)
	fixture.tuningFactory.On("Create", configuration, totalZoneCount).Return(expected)

	// Act
	options := fixture.handler.GetZoneEditorOptions(toDto(state), totalZoneCount)

	// Assert
	assert.Equal(t, expected, options.Tuning)
}
