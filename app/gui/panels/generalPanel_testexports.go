//go:build integration_test

package panels

// ScrollPosition ONLY FOR INTEGRATION TEST USE
// See layoutPanel_testexports.go for why the real widget.List is read.
func (this *GeneralPanel) ScrollPosition() (int, int) {
	return this.scroll.Position.First, this.scroll.Position.Offset
}
