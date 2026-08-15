//go:build integration_test

package panels

// ScrollPosition ONLY FOR INTEGRATION TEST USE
// Returns the list's first visible child and its pixel offset, so a test can
// prove a synthetic wheel event actually moved the panel. Reading the real
// widget.List matters because layout.List clamps at both ends, so accumulating
// the injected deltas would diverge from the truth at the first clamp.
func (this *LayoutPanel) ScrollPosition() (int, int) {
	return this.scroll.Position.First, this.scroll.Position.Offset
}
