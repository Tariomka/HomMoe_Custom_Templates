package previewRenderer_test

import (
	"fmt"
	"image"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/services"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func simpleTemplate(name string) *entities.RmgTemplate {
	return &entities.RmgTemplate{
		Name: name,
		Variants: []entities.Variant{{
			Zones: []entities.Zone{
				{Name: "Spawn-A", MainObjects: []entities.MainObject{{Type: "Spawn", Spawn: "Player1"}}},
				{Name: "Spawn-B", MainObjects: []entities.MainObject{{Type: "Spawn", Spawn: "Player2"}}},
				{Name: "Neutral-C"},
			},
			Connections: []entities.Connection{
				{From: "Spawn-A", To: "Neutral-C", ConnectionType: "Direct"},
				{From: "Neutral-C", To: "Spawn-B", ConnectionType: "Direct"},
			},
		}},
	}
}

func TestWhenAnySideIsRequested_ReturnsFixed700SquareImage(t *testing.T) {
	// Arrange
	rmgTemplate := simpleTemplate("T")

	// Act
	img := services.RenderPreviewImage(rmgTemplate, config.TopologyRing, 400)

	// Assert
	require.NotNil(t, img)
	assert.Equal(t, image.Rect(0, 0, 700, 700), img.Bounds())
}

func TestWhenTemplateIsEmpty_StillReturnsBackgroundCanvas(t *testing.T) {
	// Arrange
	rmgTemplate := &entities.RmgTemplate{}

	// Act
	img := services.RenderPreviewImage(rmgTemplate, config.TopologyRing, 100)

	// Assert
	assert.NotNil(t, img)
}

func TestWhenPortalConnectionExists_RendersWithoutPanic(t *testing.T) {
	// Arrange
	rmgTemplate := simpleTemplate("T")
	rmgTemplate.Variants[0].Connections = append(rmgTemplate.Variants[0].Connections,
		entities.Connection{From: "Spawn-A", To: "Spawn-B", ConnectionType: "Portal"},
	)

	// Act
	img := services.RenderPreviewImage(rmgTemplate, config.TopologyRing, 300)

	// Assert
	assert.NotNil(t, img)
}

func TestWhenPlayerZoneHasExtraCityCastle_RendersWithoutPanic(t *testing.T) {
	// Arrange
	rmgTemplate := simpleTemplate("T")
	rmgTemplate.Variants[0].Zones[0].MainObjects = append(rmgTemplate.Variants[0].Zones[0].MainObjects,
		entities.MainObject{Type: "City"})

	// Act
	img := services.RenderPreviewImage(rmgTemplate, config.TopologyRing, 300)

	// Assert
	assert.NotNil(t, img)
}

func TestWhenTemplateContainsHubZone_RendersWithoutPanic(t *testing.T) {
	// Arrange
	rmgTemplate := &entities.RmgTemplate{
		Variants: []entities.Variant{{
			Zones: []entities.Zone{
				{Name: "Hub"},
				{Name: "Spawn-A", MainObjects: []entities.MainObject{{Type: "Spawn", Spawn: "Player1"}}},
			},
			Connections: []entities.Connection{{From: "Hub", To: "Spawn-A", ConnectionType: "Direct"}},
		}},
	}

	// Act
	img := services.RenderPreviewImage(rmgTemplate, config.TopologyRing, 300)

	// Assert
	assert.NotNil(t, img)
}

func manualPositionTemplate(positionA, positionB [2]float64, connectionType string) *entities.RmgTemplate {
	return &entities.RmgTemplate{
		Name: "T",
		Variants: []entities.Variant{{
			Zones: []entities.Zone{
				{Name: "Spawn-A", ManualPosition: &positionA},
				{Name: "Spawn-B", ManualPosition: &positionB},
			},
			Connections: []entities.Connection{
				{From: "Spawn-A", To: "Spawn-B", ConnectionType: connectionType},
			},
		}},
	}
}

func TestWhenConnectedZonesShareOnePixel_SkipsDegenerateLineWithoutPanic(t *testing.T) {
	// Arrange
	rmgTemplate := manualPositionTemplate([2]float64{0.5, 0.5}, [2]float64{0.5, 0.5}, "Direct")

	// Act
	img := services.RenderPreviewImage(rmgTemplate, config.TopologyRing, 700)

	// Assert
	assert.NotNil(t, img)
}

func TestWhenConnectedZonesAreNearlyTouching_RendersShortLineWithoutPanic(t *testing.T) {
	// Arrange
	rmgTemplate := manualPositionTemplate([2]float64{0.5, 0.5}, [2]float64{0.505, 0.5}, "Direct")

	// Act
	img := services.RenderPreviewImage(rmgTemplate, config.TopologyRing, 700)

	// Assert
	assert.NotNil(t, img)
}

func TestWhenPortalConnectsNearbyZones_RendersDashedLineWithoutPanic(t *testing.T) {
	// Arrange
	rmgTemplate := manualPositionTemplate([2]float64{0.5, 0.5}, [2]float64{0.51, 0.5}, "Portal")

	// Act
	img := services.RenderPreviewImage(rmgTemplate, config.TopologyRing, 700)

	// Assert
	assert.NotNil(t, img)
}

func TestWhenParallelEdgesFanBeyondCanvasEdge_ClipsLinesWithoutPanic(t *testing.T) {
	// Arrange - two zones hug the top border and fourteen parallel edges fan
	// their control points far off-canvas, forcing the line rasteriser to clip.
	rmgTemplate := &entities.RmgTemplate{
		Name: "T",
		Variants: []entities.Variant{{
			Zones: []entities.Zone{
				{Name: "Neutral-A", ManualPosition: &[2]float64{0.35, 0.01}},
				{Name: "Neutral-B", ManualPosition: &[2]float64{0.65, 0.01}},
			},
		}},
	}
	for index := 0; index < 14; index++ {
		rmgTemplate.Variants[0].Connections = append(rmgTemplate.Variants[0].Connections,
			entities.Connection{From: "Neutral-A", To: "Neutral-B", ConnectionType: "Direct"})
	}

	// Act
	img := services.RenderPreviewImage(rmgTemplate, config.TopologyRing, 700)

	// Assert
	assert.NotNil(t, img)
}

func TestWhenGeneratedTemplatesAreRendered_EveryTopologyProducesFullSizeCanvas(t *testing.T) {
	topologies := []config.MapTopology{
		config.TopologyRing,
		config.TopologyChain,
		config.TopologyHubAndSpoke,
		config.TopologySharedWeb,
		config.TopologyRandom,
		config.TopologyCircles,
		config.TopologySquare,
		config.TopologyGeometric,
		config.TopologyCross,
		config.TopologyFractal,
	}
	playerCounts := []int{2, 3, 4, 6, 8}
	for _, topology := range topologies {
		for _, playerCount := range playerCounts {
			subtestName := fmt.Sprintf("WhenTopologyIs%sWith%dPlayers_ReturnsFullSizeCanvas", topology, playerCount)
			t.Run(subtestName, func(t *testing.T) {
				// Arrange
				configuration := config.NewGeneratorConfig()
				configuration.Topology = topology
				configuration.PlayerCount = playerCount
				configuration.ZoneConfiguration.NeutralZoneCount = playerCount
				rmgTemplate := template_generator.NewTemplateGenerator(configuration).Generate()

				// Act
				img := services.RenderPreviewImage(rmgTemplate, topology, 600)

				// Assert
				require.NotNil(t, img)
				assert.Equal(t, image.Rect(0, 0, 700, 700), img.Bounds())
			})
		}
	}
}

func TestWhenGeneratedTournamentTemplatesAreRendered_EveryTopologyProducesImage(t *testing.T) {
	topologies := []config.MapTopology{
		config.TopologyRing,
		config.TopologyChain,
		config.TopologyHubAndSpoke,
		config.TopologySharedWeb,
		config.TopologyRandom,
		config.TopologyCircles,
		config.TopologySquare,
		config.TopologyGeometric,
		config.TopologyCross,
		config.TopologyFractal,
	}
	playerCounts := []int{2, 3, 4, 6, 8}
	for _, topology := range topologies {
		for _, playerCount := range playerCounts {
			subtestName := fmt.Sprintf("WhenTournamentTopologyIs%sWith%dPlayers_ReturnsImage", topology, playerCount)
			t.Run(subtestName, func(t *testing.T) {
				// Arrange
				configuration := config.NewGeneratorConfig()
				configuration.Topology = topology
				configuration.PlayerCount = playerCount
				configuration.ZoneConfiguration.NeutralZoneCount = playerCount
				configuration.TournamentRules = &config.TournamentRules{
					Enabled:            true,
					FirstTournamentDay: 14,
					Interval:           7,
					PointsToWin:        2,
				}
				rmgTemplate := template_generator.NewTemplateGenerator(configuration).Generate()

				// Act
				img := services.RenderPreviewImage(rmgTemplate, topology, 600)

				// Assert
				assert.NotNil(t, img)
			})
		}
	}
}
