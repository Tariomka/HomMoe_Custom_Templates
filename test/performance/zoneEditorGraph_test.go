package performance_test

import (
	"image"
	"math"
	"testing"
	"time"

	"gioui.org/io/input"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/dialogs"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/themes"
	"github.com/Tariomka/hommoe_custom_templates/internal/common/constants"
	"github.com/Tariomka/hommoe_custom_templates/internal/composition"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/connection_editor"
	zone_services "github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
	"github.com/stretchr/testify/require"
)

// BenchmarkZoneEditorHandler_DescribeZoneEditorGraph measures the graph
// diagnostics the zone editor's status line asks for.
func BenchmarkZoneEditorHandler_DescribeZoneEditorGraph(b *testing.B) {
	handler := composition.InitializeGuiHandler()

	for _, benchmarkCase := range zoneEditorGraphCases() {
		b.Run(benchmarkCase.name, func(b *testing.B) {
			zones, connections := newZoneEditorGraph(
				b, benchmarkCase.playerCount, benchmarkCase.neutralCount, benchmarkCase.expectedConnections)

			var graph dtos.ZoneEditorGraphDto

			b.ReportAllocs()
			for b.Loop() {
				graph = handler.DescribeZoneEditorGraph(zones, connections)
			}

			require.False(b, graph.HasErrors)
		})
	}
}

// BenchmarkConnectionEditorService_FindIsolatedZones measures isolation
// detection alone, the part whose cost grows with zones × connections.
func BenchmarkConnectionEditorService_FindIsolatedZones(b *testing.B) {
	service := connection_editor.NewConnectionEditorService(zone_services.NewZoneTierService())

	for _, benchmarkCase := range zoneEditorGraphCases() {
		b.Run(benchmarkCase.name, func(b *testing.B) {
			zones, connections := newZoneEditorGraph(
				b, benchmarkCase.playerCount, benchmarkCase.neutralCount, benchmarkCase.expectedConnections)

			var isolated []string

			b.ReportAllocs()
			for b.Loop() {
				isolated = service.FindIsolatedZones(zones, connections)
			}

			require.Empty(b, isolated)
		})
	}
}

// BenchmarkZoneEditorDialog_IdleFrame measures laying out the zone editor on a
// frame with no input. It only records layout ops, so it needs no window or GPU.
func BenchmarkZoneEditorDialog_IdleFrame(b *testing.B) {
	handler := composition.InitializeGuiHandler()
	state := editor_state_dto.EditorStateDto{EditorState: editor_state_model.NewDefaultEditorStateModel()}
	theme := themes.NewTheme()

	for _, benchmarkCase := range zoneEditorGraphCases() {
		b.Run(benchmarkCase.name, func(b *testing.B) {
			zones, connections := newZoneEditorGraph(
				b, benchmarkCase.playerCount, benchmarkCase.neutralCount, benchmarkCase.expectedConnections)
			options := handler.GetZoneEditorOptions(state, len(zones))
			dialog := dialogs.NewZoneEditorDialog(
				zones, connections, options.Topology, options.Tuning, options.GenerateRoads, handler, nil, nil)

			var operations op.Ops
			router := new(input.Router)
			gtx := layout.Context{
				Ops:         &operations,
				Constraints: layout.Exact(image.Pt(1000, 720)),
				Metric:      unit.Metric{PxPerDp: 1, PxPerSp: 1},
				Source:      router.Source(),
				Now:         time.Now(),
			}
			frame := func() bool {
				operations.Reset()
				_, closed := dialog.Body(gtx, theme)
				router.Frame(&operations)

				return closed
			}
			frame()
			frame()

			closed := false

			b.ReportAllocs()
			for b.Loop() {
				closed = frame()
			}

			require.False(b, closed)
		})
	}
}

// zoneEditorGraphCase sizes one fixed zone-editor graph.
type zoneEditorGraphCase struct {
	name                string
	playerCount         int
	neutralCount        int
	expectedConnections int
}

// zoneEditorGraphCases are fixed graphs at the realistic sizes: the label pool
// caps a template at about 40 zones. They are deterministic so that separate
// before/after runs measure identical inputs.
func zoneEditorGraphCases() []zoneEditorGraphCase {
	return []zoneEditorGraphCase{
		{name: "Zones4", playerCount: 2, neutralCount: 2, expectedConnections: 6},
		{name: "Zones12", playerCount: 4, neutralCount: 8, expectedConnections: 18},
		{name: "Zones24", playerCount: 8, neutralCount: 16, expectedConnections: 36},
		{name: "Zones40", playerCount: 8, neutralCount: 32, expectedConnections: 60},
	}
}

// newZoneEditorGraph builds spawns then neutrals pinned on a circle, joined by
// a ring plus a chord from each zone in the first half to its opposite zone.
func newZoneEditorGraph(
	b *testing.B,
	playerCount, neutralCount, expectedConnections int,
) ([]template_model.Zone, []template_model.Connection) {
	b.Helper()
	labels := constants.GetZoneLabels()
	zones := make([]template_model.Zone, 0, playerCount+neutralCount)
	for index := range playerCount {
		zones = append(zones, template_model.Zone{Name: constants.GetPlayerZoneNameFor(labels[index])})
	}
	for index := range neutralCount {
		zones = append(zones, template_model.Zone{Name: constants.GetNeutralZoneNameFor(labels[index])})
	}

	zoneCount := len(zones)
	for index := range zones {
		angle := 2 * math.Pi * float64(index) / float64(zoneCount)
		zones[index].ManualPosition = new(data.NewVec2(0.5+0.4*math.Cos(angle), 0.5+0.4*math.Sin(angle)))
	}

	connections := make([]template_model.Connection, 0, expectedConnections)
	for index := range zoneCount {
		connections = append(connections, newZoneEditorGraphConnection(zones[index], zones[(index+1)%zoneCount]))
	}
	for index := range zoneCount / 2 {
		connections = append(connections, newZoneEditorGraphConnection(zones[index], zones[index+zoneCount/2]))
	}

	require.Len(b, zones, playerCount+neutralCount)
	require.Len(b, connections, expectedConnections)

	return zones, connections
}

func newZoneEditorGraphConnection(from, to template_model.Zone) template_model.Connection {
	return template_model.Connection{
		Name:           from.Name + "_" + to.Name,
		From:           from.Name,
		To:             to.Name,
		ConnectionType: "Direct",
	}
}
