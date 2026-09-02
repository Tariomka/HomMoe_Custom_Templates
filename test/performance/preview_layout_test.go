package performance_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/mappers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/preview"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/preview_service"
	zone_services "github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/require"
)

const previewCanvasSide = 600.0

// BenchmarkPreviewLayoutService_BuildPreviewLayout measures one preview layout
// build, which the preview canvas performs whenever its inputs change. Random
// is the shipped default topology and its relaxation passes are O(zones²), so
// the zone-count cases guard against super-linear regressions. It needs no
// window, so it runs identically headless or headed.
func BenchmarkPreviewLayoutService_BuildPreviewLayout(b *testing.B) {
	benchmarkCases := []struct {
		name             string
		topology         config.MapTopology
		playerCount      int
		neutralZoneCount int
	}{
		{name: "RandomDefault", topology: config.TopologyRandom, playerCount: 2, neutralZoneCount: 0},
		{name: "RandomMedium", topology: config.TopologyRandom, playerCount: 4, neutralZoneCount: 8},
		{name: "RandomLarge", topology: config.TopologyRandom, playerCount: 8, neutralZoneCount: 16},
		{name: "RingLarge", topology: config.TopologyRing, playerCount: 8, neutralZoneCount: 16},
		{name: "CirclesLarge", topology: config.TopologyCircles, playerCount: 8, neutralZoneCount: 16},
	}

	for _, benchmarkCase := range benchmarkCases {
		b.Run(benchmarkCase.name, func(b *testing.B) {
			configuration := config.NewGeneratorConfig()
			configuration.Topology = benchmarkCase.topology
			configuration.PlayerCount = benchmarkCase.playerCount
			configuration.ZoneConfiguration.NeutralZoneCount = benchmarkCase.neutralZoneCount
			generated, _ := test_helpers.NewTemplateGenerator(configuration).Generate()
			template := mappers.NewTemplateMapper().ToEntity(*generated)
			service := preview_service.NewPreviewLayoutService(zone_services.NewZoneTierService())

			var layout preview.Layout

			b.ReportAllocs()
			for b.Loop() {
				layout = service.BuildPreviewLayout(&template, benchmarkCase.topology, previewCanvasSide)
			}

			require.NotEmpty(b, layout.Positions)
		})
	}
}
