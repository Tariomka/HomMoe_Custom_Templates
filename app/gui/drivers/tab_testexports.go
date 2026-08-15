//go:build integration_test

package drivers

import "github.com/Tariomka/hommoe_custom_templates/app/gui/interfaces"

// GetPanel ONLY FOR INTEGRATION TEST USE
func (this *Tab) GetPanel() interfaces.IPanel {
	return this.panel
}
