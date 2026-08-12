package zoneEditorGeometryService_test

import (
	"image"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWhenAPointIsInsideAZone_TheHitTestNamesThatZone(t *testing.T) {
	t.Parallel()
	// Arrange
	service, _ := newGeometryFixture(nil)
	positions := chordPositions()

	// Act
	name := service.HitTestNode(image.Pt(140+fixtureZoneRadius, 350), positions, fixtureZoneRadius)

	// Assert
	assert.Equal(t, "A", name)
}

func TestWhenAPointIsJustOutsideEveryZone_TheHitTestNamesNone(t *testing.T) {
	t.Parallel()
	// Arrange
	service, _ := newGeometryFixture(nil)
	positions := chordPositions()

	// Act
	name := service.HitTestNode(image.Pt(140+fixtureZoneRadius+1, 350), positions, fixtureZoneRadius)

	// Assert
	assert.Empty(t, name)
}

func TestWhenTwoZonesCoverAPoint_TheNearestOneIsNamed(t *testing.T) {
	t.Parallel()
	// Arrange
	service, _ := newGeometryFixture(nil)
	positions := map[string]image.Point{
		"A": image.Pt(100, 100),
		"B": image.Pt(110, 100),
	}

	// Act
	name := service.HitTestNode(image.Pt(108, 100), positions, fixtureZoneRadius)

	// Assert
	assert.Equal(t, "B", name)
}

func TestWhenThereAreNoZones_TheHitTestNamesNone(t *testing.T) {
	t.Parallel()
	// Arrange
	service, _ := newGeometryFixture(nil)

	// Act
	name := service.HitTestNode(image.Pt(140, 350), nil, fixtureZoneRadius)

	// Assert
	assert.Empty(t, name)
}
