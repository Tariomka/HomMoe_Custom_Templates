//go:build integration_test && gui

package performance_test

import (
	"testing"

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
	runner := integration_common.NewAppRunner(b)

	handler := integration_common.NewHandler(runner)

	for b.Loop() {
		handler.ClickGeneralTab()
		assert.Equal(b, 0, runner.SelectedTabIndex())
		handler.ClickLayoutAndZonesTab()
		assert.Equal(b, 1, runner.SelectedTabIndex())
		handler.ClickBonusesAndBansTab()
		assert.Equal(b, 2, runner.SelectedTabIndex())
		handler.ClickGeneralTab()
		assert.Equal(b, 0, runner.SelectedTabIndex())
	}
}
