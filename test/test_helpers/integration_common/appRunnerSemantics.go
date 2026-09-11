//go:build integration_test

package integration_common

import (
	"image"

	"gioui.org/f32"
	"gioui.org/io/semantic"
)

// ButtonBounds returns the window-pixel bounds of the button labelled label.
// Every button flavour in app/gui/widgets, and every file explorer listing row,
// publishes semantic.Button plus its label through utils.AddButtonSemantics, so
// a widget can be addressed by what it says rather than by a pinned coordinate
// that a layout change silently invalidates.
//
// It fails the test unless exactly one button matches: two matches mean the
// caller cannot know which of them a click would have hit. Use ButtonBoundsIn
// to disambiguate by region.
func (this *AppRunner) ButtonBounds(label string) image.Rectangle {
	this.tb.Helper()
	return this.ButtonBoundsIn(image.Rect(0, 0, WindowWidth, WindowHeight), label)
}

// ButtonBoundsIn is ButtonBounds restricted to the buttons centered inside area.
// An open modal needs this: the editor behind the scrim is still laid out, so
// the dialog's own "Save" and the toolbar's "Save" are both present in the frame
// and only the dialog's panel rectangle tells them apart.
func (this *AppRunner) ButtonBoundsIn(area image.Rectangle, label string) image.Rectangle {
	this.tb.Helper()
	this.mu.Lock()
	defer this.mu.Unlock()
	this.frameLocked()

	matches := make([]image.Rectangle, 0, 1)
	for _, node := range this.router.AppendSemantics(nil) {
		bounds := node.Desc.Bounds
		if node.Desc.Class != semantic.Button || node.Desc.Label != label {
			continue
		}
		if bounds.Min.Add(bounds.Max).Div(2).In(area) {
			matches = append(matches, bounds)
		}
	}

	switch len(matches) {
	case 1:
		return matches[0]
	case 0:
		this.tb.Fatalf("no button labelled %q inside %v", label, area)
	default:
		this.tb.Fatalf("%d buttons labelled %q inside %v, cannot pick one", len(matches), label, area)
	}

	return image.Rectangle{}
}

// ClickButton taps the center of the button labelled label.
func (this *AppRunner) ClickButton(label string) {
	this.tb.Helper()
	this.clickBounds(this.ButtonBounds(label))
}

// ButtonLabelsIn returns the label of every button the current frame publishes
// inside area, for the assertions that are about what the UI offers rather than
// about pressing any one of them.
func (this *AppRunner) ButtonLabelsIn(area image.Rectangle) []string {
	this.tb.Helper()
	this.mu.Lock()
	defer this.mu.Unlock()
	this.frameLocked()

	labels := make([]string, 0)
	for _, node := range this.router.AppendSemantics(nil) {
		bounds := node.Desc.Bounds
		if node.Desc.Class != semantic.Button || node.Desc.Label == "" {
			continue
		}
		if bounds.Min.Add(bounds.Max).Div(2).In(area) {
			labels = append(labels, node.Desc.Label)
		}
	}

	return labels
}

// ClickButtonIn taps the center of the button labelled label inside area.
func (this *AppRunner) ClickButtonIn(area image.Rectangle, label string) {
	this.tb.Helper()
	this.clickBounds(this.ButtonBoundsIn(area, label))
}

// ClickTopmostCheckboxIn taps the highest checkbox inside area. A checkbox is
// the one widget flavour that cannot be addressed by what it says: Gio's
// material.CheckBox publishes semantic.CheckBox but no label, so a caller picks
// one by where it sits and confirms the tap by asserting the state it flipped.
// The whole row is the clickable area, so the tap lands on its centre.
func (this *AppRunner) ClickTopmostCheckboxIn(area image.Rectangle) {
	this.tb.Helper()
	this.mu.Lock()
	this.frameLocked()
	topmost := image.Rectangle{}
	for _, node := range this.router.AppendSemantics(nil) {
		bounds := node.Desc.Bounds
		if node.Desc.Class != semantic.CheckBox || !bounds.Min.In(area) {
			continue
		}
		if topmost.Empty() || bounds.Min.Y < topmost.Min.Y {
			topmost = bounds
		}
	}
	this.mu.Unlock()

	if topmost.Empty() {
		this.tb.Fatalf("no checkbox inside %v", area)
		return
	}

	this.clickBounds(topmost)
}

func (this *AppRunner) clickBounds(bounds image.Rectangle) {
	this.tb.Helper()
	center := bounds.Min.Add(bounds.Max).Div(2)
	this.ClickAt(f32.Pt(float32(center.X), float32(center.Y)))
}
