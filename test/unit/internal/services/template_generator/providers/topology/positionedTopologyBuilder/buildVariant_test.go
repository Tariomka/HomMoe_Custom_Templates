package positionedTopologyBuilder_test

import (
	"strings"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
)

func newPositionedTopologyBuilder() *topology.PositionedTopologyBuilder {
	return topology.NewPositionedTopologyBuilder(test_helpers.NewZoneFactories())
}

func TestWhenLayoutIsBuilt_StampsGeneratorPosition(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	tuning := test_helpers.NewGenerationTuning(configuration, 1)
	builder := newPositionedTopologyBuilder()
	expectedPosition := [2]float64{0.25, 0.75}
	layoutBuilder := func([]string, neutral_zone.Plans) ([]string, models.Positions, []models.ConnectionIndexes) {
		return []string{"A"}, models.Positions{data.NewVec2(expectedPosition[0], expectedPosition[1])}, nil
	}

	// Act
	variant := builder.BuildVariant(*configuration, []string{"A"}, nil, tuning, "", layoutBuilder, nil)

	// Assert
	assert.Equal(t, &expectedPosition, variant.Zones[0].GeneratorPosition)
}

func TestWhenLayoutContainsPair_CreatesDirectConnection(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	neutralZones := neutral_zone.Plans{}
	neutralZones.AddPlan("N", neutral_zone.QualityMedium, 1)
	tuning := test_helpers.NewGenerationTuning(configuration, 2)
	builder := newPositionedTopologyBuilder()
	layoutBuilder := func([]string, neutral_zone.Plans) ([]string, models.Positions, []models.ConnectionIndexes) {
		return []string{"A", "N"}, models.Positions{data.NewVec2(0.0, 0.0), data.NewVec2(1.0, 0.0)},
			[]models.ConnectionIndexes{{X: 0, Y: 1}}
	}

	// Act
	variant := builder.BuildVariant(*configuration, []string{"A"}, neutralZones, tuning, "", layoutBuilder, nil)

	// Assert
	assert.Equal(t, []string{"Rnd-A-N"}, connectionNames(variant.Connections))
}

func TestWhenDirectPlayerConnectionsAreDisabled_SkipsLayoutPlayerPair(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.NoDirectPlayerConnections = true
	tuning := test_helpers.NewGenerationTuning(configuration, 2)
	builder := newPositionedTopologyBuilder()
	layoutBuilder := func([]string, neutral_zone.Plans) ([]string, models.Positions, []models.ConnectionIndexes) {
		return []string{"A", "B"}, models.Positions{data.NewVec2(0.0, 0.0), data.NewVec2(1.0, 0.0)},
			[]models.ConnectionIndexes{{X: 0, Y: 1}}
	}

	// Act
	variant := builder.BuildVariant(*configuration, []string{"A", "B"}, nil, tuning, "", layoutBuilder, nil)

	// Assert
	assert.NotContains(t, connectionNames(variant.Connections), "Rnd-A-B")
}

func TestWhenRandomPortalsAreEnabled_AddsPortalConnection(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.RandomPortals = true
	configuration.MaxPortalConnections = 1
	neutralZones := neutral_zone.Plans{}
	neutralZones.AddPlan("N1", neutral_zone.QualityLow, 0)
	neutralZones.AddPlan("N2", neutral_zone.QualityMedium, 1)
	tuning := test_helpers.NewGenerationTuning(configuration, 3)
	builder := newPositionedTopologyBuilder()
	layoutBuilder := func([]string, neutral_zone.Plans) ([]string, models.Positions, []models.ConnectionIndexes) {
		return []string{"A", "N1", "N2"}, models.Positions{
			data.NewVec2(0.0, 0.0), data.NewVec2(0.5, 0.0), data.NewVec2(1.0, 0.0)}, nil
	}

	// Act
	variant := builder.BuildVariant(*configuration, []string{"A"}, neutralZones, tuning, "", layoutBuilder, nil)

	// Assert
	assert.Equal(t, 1, countConnectionsWithPrefix(variant.Connections, "Portal-"))
}

func TestWhenLayoutIsDisconnected_AddsBridgeConnection(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	neutralZones := neutral_zone.Plans{}
	neutralZones.AddPlan("N1", neutral_zone.QualityLow, 0)
	neutralZones.AddPlan("N2", neutral_zone.QualityMedium, 1)
	tuning := test_helpers.NewGenerationTuning(configuration, 3)
	builder := newPositionedTopologyBuilder()
	layoutBuilder := func([]string, neutral_zone.Plans) ([]string, models.Positions, []models.ConnectionIndexes) {
		return []string{"A", "N1", "N2"}, models.Positions{
				data.NewVec2(0.0, 0.0), data.NewVec2(0.5, 0.0), data.NewVec2(1.0, 0.0)},
			[]models.ConnectionIndexes{{X: 0, Y: 1}}
	}

	// Act
	variant := builder.BuildVariant(*configuration, []string{"A"}, neutralZones, tuning, "", layoutBuilder, nil)

	// Assert
	assert.Equal(t, 1, countConnectionsWithPrefix(variant.Connections, "Bridge-"))
}

func TestWhenZoneDecoratorProvided_AppliesItToBuiltZones(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	tuning := test_helpers.NewGenerationTuning(configuration, 1)
	builder := newPositionedTopologyBuilder()
	layoutBuilder := func([]string, neutral_zone.Plans) ([]string, models.Positions, []models.ConnectionIndexes) {
		return []string{"A"}, models.Positions{data.NewVec2(0.0, 0.0)}, nil
	}
	expectedRing := 3
	decorator := func(zones []entities.Zone, _ []string, _ []string, _ neutral_zone.Plans) {
		zones[0].GeneratorRing = &expectedRing
	}

	// Act
	variant := builder.BuildVariant(*configuration, []string{"A"}, nil, tuning, "", layoutBuilder, decorator)

	// Assert
	assert.Equal(t, &expectedRing, variant.Zones[0].GeneratorRing)
}

func connectionNames(connections []entities.Connection) []string {
	names := make([]string, len(connections))
	for index, connection := range connections {
		names[index] = connection.Name
	}
	return names
}

func countConnectionsWithPrefix(connections []entities.Connection, prefix string) int {
	count := 0
	for _, connection := range connections {
		if strings.HasPrefix(connection.Name, prefix) {
			count++
		}
	}
	return count
}
