//go:build integration_test

// Package performance_test benchmarks the editor.Window UI through the shared
// integration_common.AppRunner. The SAME benchmark runs headless by default
// (CI-safe) or on screen with `go test ... -args headed`.
package performance_test

import (
	"testing"
	"time"

	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers/integration_common"
	"github.com/stretchr/testify/assert"
)

// TestMain defers to the shared runner so headed mode can start app.Main on the
// main goroutine before the testing framework parses flags.
func TestMain(m *testing.M) {
	integration_common.RunMain(m)
}

// BenchmarkEditorWindow_TabCycling starts the editor, clicks through every tab
// several times and gracefully shuts down. It runs identically headless or with
// a real on-screen window (the latter additionally renders the UI); select the
// mode with `go test ... -args headed`.
func BenchmarkEditorWindow_TabCycling(b *testing.B) {
	runner := integration_common.NewAppRunner()
	runner.Start()
	defer runner.Stop()

	points := integration_common.CalibrateTabPoints(b, runner)
	runner.SetRenderDelay(100 * time.Millisecond)

	// TODO: add benchmark assertion with
	// actual := testing.Benchmark()
	// assert.LessOrEqual(b, actual.AllocedBytesPerOp(), 155_000)

	for b.Loop() {
		for idx := range points {
			runner.ClickAt(points[idx])
			assert.Equal(b, idx, runner.SelectedTabIndex())
		}
	}
}
