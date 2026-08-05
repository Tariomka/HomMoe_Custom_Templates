package performance_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/require"
)

// BenchmarkTemplateGenerator_Generate measures a full template generation per
// topology, covering the topology-lookup hot path. It needs no window, so it
// runs identically headless or headed.
func BenchmarkTemplateGenerator_Generate(b *testing.B) {
	benchmarkCases := []struct {
		name           string
		topology       config.MapTopology
		playerCount    int
		tournamentMode bool
	}{
		{name: "Ring", topology: config.TopologyRing, playerCount: 8},
		{name: "HubAndSpoke", topology: config.TopologyHubAndSpoke, playerCount: 8},
		{name: "GeometricHub", topology: config.TopologyGeometricHub, playerCount: 8},
		{name: "Fractal", topology: config.TopologyFractal, playerCount: 8},
		{name: "Tournament", topology: config.TopologyChain, playerCount: 2, tournamentMode: true},
	}

	for _, benchmarkCase := range benchmarkCases {
		b.Run(benchmarkCase.name, func(b *testing.B) {
			configuration := config.NewGeneratorConfig()
			configuration.Topology = benchmarkCase.topology
			configuration.PlayerCount = benchmarkCase.playerCount
			if benchmarkCase.tournamentMode {
				configuration.TournamentRules.Enabled = true
			}
			generator := test_helpers.NewTemplateGenerator(configuration)

			var template *entities.RmgTemplate

			b.ReportAllocs()
			for b.Loop() {
				template = generator.Generate()
			}

			require.NotEmpty(b, template.Variants)
		})
	}
}
