//go:build integration_test

package drivers

import "github.com/Tariomka/hommoe_custom_templates/app/gui/interfaces"

// GetTopDialog ONLY FOR INTEGRATION TEST USE
func (this *DialogHost) GetTopDialog() interfaces.IDialog {
	return this.getTopDialog()
}
