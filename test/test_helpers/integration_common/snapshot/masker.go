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
