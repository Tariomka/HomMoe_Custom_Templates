package snapshot

import (
	"image"
	"image/color"
	"image/draw"
)

// Masker hides nondeterministic screen regions (e.g. the randomly
// generated map preview) by painting them a solid color before a screenshot is
// saved as a golden or compared against one. The same mask is applied to both
// sides, so masked regions can never cause a difference.
type Masker struct {
	masks []image.Rectangle
}

// AddRect registers a rectangle (in window pixel coordinates) to be masked.
func (this *Masker) AddRect(rect image.Rectangle) {
	this.masks = append(this.masks, rect)
}

// RemoveRect drops one previously registered rectangle and reports whether it
// was registered. It exists for the regions that are only nondeterministic while
// something transient is on screen, such as an open dialog: leaving their mask
// in place would blank a part of the editor that every later snapshot could
// otherwise compare.
func (this *Masker) RemoveRect(rect image.Rectangle) bool {
	for index, mask := range this.masks {
		if mask == rect {
			this.masks = append(this.masks[:index], this.masks[index+1:]...)
			return true
		}
	}

	return false
}

// Apply paints every registered rectangle onto the screenshot in place,
// clamping each rectangle to the screenshot bounds.
func (this *Masker) Apply(screenshot *image.RGBA) {
	maskFill := image.NewUniform(color.NRGBA{A: 255})
	for _, rect := range this.masks {
		clamped := rect.Intersect(screenshot.Bounds())
		if clamped.Empty() {
			continue
		}
		draw.Draw(screenshot, clamped, maskFill, image.Point{}, draw.Src)
	}
}
