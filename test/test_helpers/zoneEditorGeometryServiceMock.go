package test_helpers

import (
	"image"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/stretchr/testify/mock"
)

// ZoneEditorGeometryServiceMock is a testify mock of
// connection_editor.IZoneEditorGeometryService, used to unit-test collaborators
// without computing a real canvas layout.
type ZoneEditorGeometryServiceMock struct {
	mock.Mock
}

func (this *ZoneEditorGeometryServiceMock) BuildGeometry(
	zones []entities.Zone,
	connections []entities.Connection,
	topology config.MapTopology,
	canvasSide int) models.ZoneEditorGeometry {
	arguments := this.Called(zones, connections, topology, canvasSide)
	geometry, _ := arguments.Get(0).(models.ZoneEditorGeometry)
	return geometry
}

func (this *ZoneEditorGeometryServiceMock) HitTestNode(
	position image.Point,
	positions map[string]image.Point,
	zoneRadius int) string {
	arguments := this.Called(position, positions, zoneRadius)
	return arguments.String(0)
}

func (this *ZoneEditorGeometryServiceMock) HitTestEdge(
	position image.Point,
	edges []models.ZoneEditorEdge) int {
	arguments := this.Called(position, edges)
	return arguments.Int(0)
}

func (this *ZoneEditorGeometryServiceMock) GridStep(zoneRadius int) float64 {
	arguments := this.Called(zoneRadius)
	step, _ := arguments.Get(0).(float64)
	return step
}

func (this *ZoneEditorGeometryServiceMock) SnapPosition(
	position image.Point,
	positions map[string]image.Point,
	zoneRadius int,
	draggedZone string) models.ZoneEditorSnapResult {
	arguments := this.Called(position, positions, zoneRadius, draggedZone)
	result, _ := arguments.Get(0).(models.ZoneEditorSnapResult)
	return result
}
