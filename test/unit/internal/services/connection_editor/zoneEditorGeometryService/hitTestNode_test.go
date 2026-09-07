package zoneEditorGeometryService_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestWhenAPointIsInsideAZone_TheHitTestNamesThatZone(t *testing.T) {
	t.Parallel()
	// Arrange
	service, _ := newGeometryFixture(nil)
	positions := chordPositions()

	// Act
	name := service.HitTestNode(data.NewVec2(140+fixtureZoneRadius, 350.0), positions, fixtureZoneRadius)

	// Assert
	assert.Equal(t, "A", name)
}

func TestWhenAPointIsJustOutsideEveryZone_TheHitTestNamesNone(t *testing.T) {
	t.Parallel()
	// Arrange
	service, _ := newGeometryFixture(nil)
	positions := chordPositions()

	// Act
	name := service.HitTestNode(data.NewVec2(140+fixtureZoneRadius+1, 350.0), positions, fixtureZoneRadius)

	// Assert
	assert.Empty(t, name)
}

func TestWhenTwoZonesCoverAPoint_TheNearestOneIsNamed(t *testing.T) {
	t.Parallel()
	// Arrange
	service, _ := newGeometryFixture(nil)
	positions := map[string]models.Position{
		"A": data.NewVec2(100.0, 100.0),
		"B": data.NewVec2(110.0, 100.0),
	}

	// Act
	name := service.HitTestNode(data.NewVec2(108.0, 100.0), positions, fixtureZoneRadius)

	// Assert
	assert.Equal(t, "B", name)
}

func TestWhenThereAreNoZones_TheHitTestNamesNone(t *testing.T) {
	t.Parallel()
	// Arrange
	service, _ := newGeometryFixture(nil)

	// Act
	name := service.HitTestNode(data.NewVec2(140.0, 350.0), nil, fixtureZoneRadius)

	// Assert
	assert.Empty(t, name)
}
