//go:build integration_test && gui

package gui_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers/integration_common"
)

// TestMain defers to the shared runner so headed mode can start app.Main on the
// main goroutine before the testing framework parses flags. Headless runs (the
// default) never start app.Main and need no display.
func TestMain(m *testing.M) {
	integration_common.RunMain(m)
}
