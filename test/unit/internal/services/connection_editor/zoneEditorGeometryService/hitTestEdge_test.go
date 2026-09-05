package zoneEditorGeometryService_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestWhenAPointSitsOnACurve_TheEdgeHitTestReturnsThatEdge(t *testing.T) {
	t.Parallel()
	// Arrange
	service, _ := newGeometryFixture(nil)
	edges := []models.ZoneEditorEdge{straightEdge(0, 0)}

	// Act
	index := service.HitTestEdge(data.NewVec2(50.0, 5.0), edges)

	// Assert
	assert.Equal(t, 0, index)
}

func TestWhenAPointIsFarFromEveryCurve_TheEdgeHitTestReturnsNoIndex(t *testing.T) {
	t.Parallel()
	// Arrange
	service, _ := newGeometryFixture(nil)
	edges := []models.ZoneEditorEdge{straightEdge(0, 0)}

	// Act
	index := service.HitTestEdge(data.NewVec2(50.0, 50.0), edges)

	// Assert
	assert.Equal(t, -1, index)
}

func TestWhenTwoCurvesAreInReach_TheNearestOneIsReturned(t *testing.T) {
	t.Parallel()
	// Arrange
	service, _ := newGeometryFixture(nil)
	edges := []models.ZoneEditorEdge{straightEdge(0, 0), straightEdge(1, 20)}

	// Act
	index := service.HitTestEdge(data.NewVec2(50.0, 18.0), edges)

	// Assert
	assert.Equal(t, 1, index)
}

func TestWhenThereAreNoEdges_TheEdgeHitTestReturnsNoIndex(t *testing.T) {
	t.Parallel()
	// Arrange
	service, _ := newGeometryFixture(nil)

	// Act
	index := service.HitTestEdge(data.NewVec2(50.0, 0.0), nil)

	// Assert
	assert.Equal(t, -1, index)
}

// straightEdge is a degenerate quadratic curve: a horizontal line at the given
// height, which keeps the expected distances in these tests obvious.
func straightEdge(connectionIndex int, height float64) models.ZoneEditorEdge {
	return models.ZoneEditorEdge{
		ConnectionIndex: connectionIndex,
		StartPoint:      data.NewVec2(0.0, height),
		ControlPoint:    data.NewVec2(50.0, height),
		EndPoint:        data.NewVec2(100.0, height),
	}
}
