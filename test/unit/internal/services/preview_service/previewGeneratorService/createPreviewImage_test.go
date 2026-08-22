package previewGeneratorService_test

import (
	"image"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/preview_service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenTemplateIsRendered_ReturnsFullSizeCanvas(t *testing.T) {
	t.Parallel()
	// Arrange
	generator := mustNewGenerator(t)

	// Act
	canvas := generator.CreatePreviewImage(ringTemplate(), config.TopologyRing)

	// Assert
	assert.Equal(t, image.Rect(0, 0, 700, 700), canvas.Bounds())
}

func TestWhenTemplateIsNil_ReturnsBackgroundOnlyCanvas(t *testing.T) {
	t.Parallel()
	// Arrange
	generator := mustNewGenerator(t)
	backgroundOnly := generator.CreatePreviewImage(&entities.RmgTemplate{}, config.TopologyRing)

	// Act
	canvas := generator.CreatePreviewImage(nil, config.TopologyRing)

	// Assert
	assert.Equal(t, backgroundOnly.Pix, canvas.Pix)
}

func TestWhenTemplateHasZones_DrawsThemOverTheBackground(t *testing.T) {
	t.Parallel()
	// Arrange
	generator := mustNewGenerator(t)
	backgroundOnly := generator.CreatePreviewImage(&entities.RmgTemplate{}, config.TopologyRing)

	// Act
	canvas := generator.CreatePreviewImage(ringTemplate(), config.TopologyRing)

	// Assert
	assert.NotEqual(t, backgroundOnly.Pix, canvas.Pix)
}

func TestWhenTemplateHasConnections_DrawsLinesBetweenZones(t *testing.T) {
	t.Parallel()
	// Arrange
	generator := mustNewGenerator(t)
	disconnectedTemplate := ringTemplate()
	disconnectedTemplate.Variants[0].Connections = nil
	withoutConnections := generator.CreatePreviewImage(disconnectedTemplate, config.TopologyRing)

	// Act
	canvas := generator.CreatePreviewImage(ringTemplate(), config.TopologyRing)

	// Assert
	assert.NotEqual(t, withoutConnections.Pix, canvas.Pix)
}

func TestWhenConnectionIsPortal_DrawsDashedLineDifferentFromSolid(t *testing.T) {
	t.Parallel()
	// Arrange
	generator := mustNewGenerator(t)
	solidTemplate := ringTemplate()
	solidRender := generator.CreatePreviewImage(solidTemplate, config.TopologyRing)
	portalTemplate := ringTemplate()
	for index := range portalTemplate.Variants[0].Connections {
		portalTemplate.Variants[0].Connections[index].ConnectionType = "Portal"
	}

	// Act
	canvas := generator.CreatePreviewImage(portalTemplate, config.TopologyRing)

	// Assert
	assert.NotEqual(t, solidRender.Pix, canvas.Pix)
}

func TestWhenConnectionIsGladiatorArena_DrawsArenaMarkerOverTheSolidLine(t *testing.T) {
	t.Parallel()
	// Arrange
	generator := mustNewGenerator(t)
	solidRender := generator.CreatePreviewImage(ringTemplate(), config.TopologyRing)
	arenaTemplate := ringTemplate()
	arenaTemplate.Variants[0].Connections[0].ConnectionType = "GladiatorArena"

	// Act
	canvas := generator.CreatePreviewImage(arenaTemplate, config.TopologyRing)

	// Assert
	assert.NotEqual(t, solidRender.Pix, canvas.Pix)
}

func TestWhenZoneHostsTheArena_DrawsArenaBubbleInsteadOfThePlainOne(t *testing.T) {
	t.Parallel()
	// Arrange
	generator := mustNewGenerator(t)
	plainRender := generator.CreatePreviewImage(ringTemplate(), config.TopologyRing)
	arenaTemplate := ringTemplate()
	arenaTemplate.Variants[0].Zones[1].MainObjects = append(
		arenaTemplate.Variants[0].Zones[1].MainObjects,
		entities.MainObject{Type: "GladiatorArena"})

	// Act
	canvas := generator.CreatePreviewImage(arenaTemplate, config.TopologyRing)

	// Assert
	assert.NotEqual(t, plainRender.Pix, canvas.Pix)
}

func TestWhenSameTemplateIsRenderedTwice_ProducesIdenticalImages(t *testing.T) {
	t.Parallel()
	// Arrange
	generator := mustNewGenerator(t)
	firstRender := generator.CreatePreviewImage(ringTemplate(), config.TopologyRing)

	// Act
	secondRender := generator.CreatePreviewImage(ringTemplate(), config.TopologyRing)

	// Assert
	assert.Equal(t, firstRender.Pix, secondRender.Pix)
}

// mustNewGenerator fails the test immediately when the embedded assets cannot load.
func mustNewGenerator(t *testing.T) preview_service.IPreviewGeneratorService {
	t.Helper()
	generator, err := preview_service.NewPreviewGenerator(preview_service.NewPreviewLayoutService())
	require.NoError(t, err)
	return generator
}

// ringTemplate builds a small two-player ring template with plain connections.
func ringTemplate() *entities.RmgTemplate {
	return &entities.RmgTemplate{
		Variants: []entities.Variant{{
			Zones: []entities.Zone{
				{Name: "Spawn-A"}, {Name: "Neutral-B"},
				{Name: "Spawn-C"}, {Name: "Neutral-D"},
			},
			Connections: []entities.Connection{
				{From: "Spawn-A", To: "Neutral-B", ConnectionType: "Direct"},
				{From: "Neutral-B", To: "Spawn-C", ConnectionType: "Direct"},
				{From: "Spawn-C", To: "Neutral-D", ConnectionType: "Direct"},
				{From: "Neutral-D", To: "Spawn-A", ConnectionType: "Direct"},
			},
		}},
	}
}
