package performance_test

import (
	"image"
	"math"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/preview"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/preview_service"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/require"
)

const (
	rasterZoneRadius  = 21.0
	rasterAnchorZone  = "Anchor"
	rasterCanvasSide  = 700.0
	rasterEdgeSpread  = 300.0
	rasterAnchorPoint = rasterCanvasSide / 2
)

// BenchmarkPreviewGeneratorService_CreatePreviewImage measures one full preview
// raster, which the preview panel and the exported PNG both pay for. The layout is
// fixed and holds no zones, so the connector strokes are not hidden behind per-zone
// artwork costs, and the edge counts are large enough that a per-edge full-image
// buffer would show up in the reported allocations.
func BenchmarkPreviewGeneratorService_CreatePreviewImage(b *testing.B) {
	benchmarkCases := []struct {
		name           string
		edgeCount      int
		connectionType preview.ConnectionType
		hasRoad        bool
		stacked        bool
	}{
		{name: "NoConnections", edgeCount: 0},
		{name: "DirectRoaded_1", edgeCount: 1, connectionType: preview.ConnectionTypeDirect, hasRoad: true},
		{name: "DirectRoadless_1", edgeCount: 1, connectionType: preview.ConnectionTypeDirect},
		{name: "DirectRoaded_32", edgeCount: 32, connectionType: preview.ConnectionTypeDirect, hasRoad: true},
		{name: "DirectRoadless_32", edgeCount: 32, connectionType: preview.ConnectionTypeDirect},
		{name: "DirectRoaded_128", edgeCount: 128, connectionType: preview.ConnectionTypeDirect, hasRoad: true},
		{name: "DirectRoadless_128", edgeCount: 128, connectionType: preview.ConnectionTypeDirect},
		{name: "PortalRoaded_32", edgeCount: 32, connectionType: preview.ConnectionTypePortal, hasRoad: true},
		{name: "PortalRoadless_32", edgeCount: 32, connectionType: preview.ConnectionTypePortal},
		{name: "PortalRoadless_128", edgeCount: 128, connectionType: preview.ConnectionTypePortal},
		{name: "DirectRoadedStacked_128", edgeCount: 128,
			connectionType: preview.ConnectionTypeDirect, hasRoad: true, stacked: true},
		{name: "DirectRoadlessStacked_128", edgeCount: 128,
			connectionType: preview.ConnectionTypeDirect, stacked: true},
	}

	for _, benchmarkCase := range benchmarkCases {
		b.Run(benchmarkCase.name, func(b *testing.B) {
			fixture := test_helpers.PreviewRasterFixture{Layout: preview.Layout{
				Positions: map[string]data.Vec2[float64]{
					rasterAnchorZone: data.NewVec2(rasterAnchorPoint, rasterAnchorPoint)},
				Connections: newRasterConnections(benchmarkCase.edgeCount,
					benchmarkCase.connectionType, benchmarkCase.hasRoad, benchmarkCase.stacked),
				ZoneRadius: rasterZoneRadius,
			}}
			generator, err := preview_service.NewPreviewGenerator(fixture)
			require.NoError(b, err)

			var canvas *image.RGBA

			b.ReportAllocs()
			for b.Loop() {
				canvas = generator.CreatePreviewImage(nil, config.TopologyRandom)
			}

			require.NotNil(b, canvas)
		})
	}
}

// newRasterConnections lays the edges out as chords through the canvas center, either
// fanned around it or stacked on one geometry. Stacking is the worst case for a
// per-edge composite, because every edge covers the same pixels.
func newRasterConnections(
	edgeCount int, connectionType preview.ConnectionType, hasRoad, stacked bool) []preview.Connection {
	center := data.NewVec2(rasterAnchorPoint, rasterAnchorPoint)
	connections := make([]preview.Connection, 0, edgeCount)
	for index := range edgeCount {
		angle := math.Pi * float64(index) / math.Max(float64(edgeCount), 1)
		if stacked {
			angle = 0
		}

		offset := data.NewVec2(rasterEdgeSpread*math.Cos(angle), rasterEdgeSpread*math.Sin(angle))
		connections = append(connections, preview.Connection{
			Start:   center.Add(offset),
			Ctrl:    center,
			End:     center.Subtract(offset),
			Type:    connectionType,
			HasRoad: hasRoad,
		})
	}

	return connections
}
