package positionedTopologyBuilder_test

import (
	"strings"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
)

func TestWhenLayoutIsBuilt_StampsGeneratorPosition(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	tuning := test_helpers.NewGenerationTuning(configuration, 1)
	builder := newPositionedTopologyBuilder()
	expectedPosition := data.NewVec2(0.25, 0.75)
	layoutBuilder := func([]string, neutral_zone.Plans) ([]string, models.Positions, []models.ConnectionIndexes) {
		return []string{"A"}, models.Positions{expectedPosition}, nil
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

func TestWhenRoadlessIsolatedPlayersAreLinkedThroughNeutral_NoFallbackIsCreated(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.GenerateRoads = false
	configuration.NoDirectPlayerConnections = true
	neutralZones := neutral_zone.Plans{}
	neutralZones.AddPlan("N", neutral_zone.QualityMedium, 1)
	tuning := test_helpers.NewGenerationTuning(configuration, 3)
	builder := newPositionedTopologyBuilder()
	layoutBuilder := func([]string, neutral_zone.Plans) ([]string, models.Positions, []models.ConnectionIndexes) {
		return []string{"A", "N", "B"}, models.Positions{
			data.NewVec2(0.0, 0.0), data.NewVec2(0.5, 0.0), data.NewVec2(1.0, 0.0),
		}, []models.ConnectionIndexes{{X: 0, Y: 1}, {X: 1, Y: 2}}
	}

	// Act
	variant := builder.BuildVariant(*configuration, []string{"A", "B"}, neutralZones, tuning, "", layoutBuilder, nil)

	// Assert
	assert.Equal(t, []string{"Rnd-A-N", "Rnd-N-B"}, connectionNames(variant.Connections))
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

func TestWhenRoadlessLayoutIsDisconnected_PreservesBridgeWithoutConnectionRoads(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.GenerateRoads = false
	neutralZones := neutral_zone.Plans{}
	neutralZones.AddPlan("N1", neutral_zone.QualityLow, 0)
	neutralZones.AddPlan("N2", neutral_zone.QualityMedium, 1)
	tuning := test_helpers.NewGenerationTuning(configuration, 3)
	builder := newPositionedTopologyBuilder()
	layoutBuilder := func([]string, neutral_zone.Plans) ([]string, models.Positions, []models.ConnectionIndexes) {
		return []string{"A", "N1", "N2"}, models.Positions{
			data.NewVec2(0.0, 0.0), data.NewVec2(0.5, 0.0), data.NewVec2(1.0, 0.0),
		}, []models.ConnectionIndexes{{X: 0, Y: 1}}
	}

	// Act
	variant := builder.BuildVariant(*configuration, []string{"A"}, neutralZones, tuning, "", layoutBuilder, nil)

	// Assert
	assert.Equal(t, []string{"Rnd-A-N1", "Bridge-N1-N2"}, connectionNames(variant.Connections))
}

func TestWhenRoadlessLayoutIsDisconnected_DoesNotMaterializeBridgeConnectionRoads(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.GenerateRoads = false
	neutralZones := neutral_zone.Plans{}
	neutralZones.AddPlan("N1", neutral_zone.QualityLow, 0)
	neutralZones.AddPlan("N2", neutral_zone.QualityMedium, 1)
	tuning := test_helpers.NewGenerationTuning(configuration, 3)
	builder := newPositionedTopologyBuilder()
	layoutBuilder := func([]string, neutral_zone.Plans) ([]string, models.Positions, []models.ConnectionIndexes) {
		return []string{"A", "N1", "N2"}, models.Positions{
			data.NewVec2(0.0, 0.0), data.NewVec2(0.5, 0.0), data.NewVec2(1.0, 0.0),
		}, []models.ConnectionIndexes{{X: 0, Y: 1}}
	}

	// Act
	variant := builder.BuildVariant(*configuration, []string{"A"}, neutralZones, tuning, "", layoutBuilder, nil)

	// Assert
	assert.Empty(t, connectionRoadTargets(variant.Zones))
}

func TestWhenRoadlessIsolationHasNoNeutral_FallbackIsKeptWithoutConnectionRoads(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.GenerateRoads = false
	configuration.NoDirectPlayerConnections = true
	tuning := test_helpers.NewGenerationTuning(configuration, 2)
	builder := newPositionedTopologyBuilder()
	layoutBuilder := func([]string, neutral_zone.Plans) ([]string, models.Positions, []models.ConnectionIndexes) {
		return []string{"A", "B"}, models.Positions{data.NewVec2(0.0, 0.0), data.NewVec2(1.0, 0.0)}, nil
	}

	// Act
	variant := builder.BuildVariant(*configuration, []string{"A", "B"}, nil, tuning, "", layoutBuilder, nil)

	// Assert
	assert.Equal(t, []string{"Fallback-A-B"}, connectionNames(variant.Connections))
}

func TestWhenRoadlessIsolationHasNoNeutral_DoesNotMaterializeFallbackConnectionRoads(t *testing.T) {
	t.Parallel()
	// Arrange
	configuration := config.NewGeneratorConfig()
	configuration.GenerateRoads = false
	configuration.NoDirectPlayerConnections = true
	tuning := test_helpers.NewGenerationTuning(configuration, 2)
	builder := newPositionedTopologyBuilder()
	layoutBuilder := func([]string, neutral_zone.Plans) ([]string, models.Positions, []models.ConnectionIndexes) {
		return []string{"A", "B"}, models.Positions{data.NewVec2(0.0, 0.0), data.NewVec2(1.0, 0.0)}, nil
	}

	// Act
	variant := builder.BuildVariant(*configuration, []string{"A", "B"}, nil, tuning, "", layoutBuilder, nil)

	// Assert
	assert.Empty(t, connectionRoadTargets(variant.Zones))
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
	decorator := func(zones []template_model.Zone, _ []string, _ []string, _ neutral_zone.Plans) {
		zones[0].GeneratorRing = &expectedRing
	}

	// Act
	variant := builder.BuildVariant(*configuration, []string{"A"}, nil, tuning, "", layoutBuilder, decorator)

	// Assert
	assert.Equal(t, &expectedRing, variant.Zones[0].GeneratorRing)
}

func connectionNames(connections []template_model.Connection) []string {
	names := make([]string, len(connections))
	for index, connection := range connections {
		names[index] = connection.Name
	}
	return names
}

func countConnectionsWithPrefix(connections []template_model.Connection, prefix string) int {
	count := 0
	for _, connection := range connections {
		if strings.HasPrefix(connection.Name, prefix) {
			count++
		}
	}
	return count
}

func connectionRoadTargets(zones []template_model.Zone) []string {
	var connectionNames []string
	for _, zone := range zones {
		for _, road := range zone.Roads {
			if road.From.Type == "Connection" && len(road.From.Args) > 0 {
				connectionNames = append(connectionNames, road.From.Args[0])
			}
			if road.To.Type == "Connection" && len(road.To.Args) > 0 {
				connectionNames = append(connectionNames, road.To.Args[0])
			}
		}
	}
	return connectionNames
}

func newPositionedTopologyBuilder() *topology.PositionedTopologyBuilder {
	return topology.NewPositionedTopologyBuilder(test_helpers.NewZoneFactories())
}
