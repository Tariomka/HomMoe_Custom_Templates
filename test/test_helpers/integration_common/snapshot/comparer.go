package snapshot

import (
	"fmt"
	"image"
)

// The defaults below were derived from a CI failure artifact captured on
// ubuntu-latest under Xvfb + Mesa llvmpipe and compared against goldens rendered
// on a real GPU. Two effects were measured on that frame pair:
//
//   - Software rasterization produces about 0.75x the anti-aliasing coverage of
//     a real GPU inside rounded-rect clips, so the toolbar buttons and the tab
//     strip come out slightly dimmer. Its largest single-channel delta is ~40,
//     and it covers ~0.0005% of the frame once measured at a tolerance of 64.
//   - A layout shift of a single pixel per section header (a glyph resolved
//     through OS font fallback) moved the whole options panel. It covered 3.42%
//     of the frame at that same tolerance, yet its mean distance was only 0.66%
//   - which is how it slipped past the previous 2% mean-only gate.
//
// The gap between 0.0005% and 3.42% is what the pixel-fraction gate exists to
// exploit. The mean gate is kept as a coarse safety net for a change that is
// spread too thinly to trip the fraction gate.
const (
	// DefaultMeanThreshold is the maximum allowed mean color distance (0.25%).
	DefaultMeanThreshold = 0.0025
	// DefaultPixelTolerance is how far a single channel may drift before the
	// pixel counts as changed. 64 of 255 is well above software-rasterizer
	// anti-aliasing noise and well below any real color change.
	DefaultPixelTolerance = 64
	// DefaultChangedPixelThreshold is the maximum allowed share of changed
	// pixels (0.05%). A solid 25x25 widget covers more of the frame than this.
	DefaultChangedPixelThreshold = 0.0005
)

// Comparer measures how different two screenshots are. Alpha is ignored:
// screenshots are opaque.
type Comparer struct {
	// MeanThreshold bounds Difference.MeanDistance.
	MeanThreshold float64
	// PixelTolerance is the per-channel delta above which a pixel counts as
	// changed, in 0-255 units.
	PixelTolerance uint8
	// ChangedPixelThreshold bounds Difference.ChangedPixelFraction.
	ChangedPixelThreshold float64
}

// NewComparer builds a comparer with the measured defaults documented above.
func NewComparer() Comparer {
	return Comparer{
		MeanThreshold:         DefaultMeanThreshold,
		PixelTolerance:        DefaultPixelTolerance,
		ChangedPixelThreshold: DefaultChangedPixelThreshold,
	}
}

// Compare measures the two images. It errors when the image dimensions differ
// (callers treat that as a validation failure).
func (this Comparer) Compare(golden, actual image.Image) (Difference, error) {
	goldenBounds := golden.Bounds()
	actualBounds := actual.Bounds()
	if goldenBounds.Dx() != actualBounds.Dx() || goldenBounds.Dy() != actualBounds.Dy() {
		return Difference{}, fmt.Errorf(
			"snapshot dimensions differ: golden %dx%d, actual %dx%d",
			goldenBounds.Dx(), goldenBounds.Dy(), actualBounds.Dx(), actualBounds.Dy())
	}

	width := goldenBounds.Dx()
	height := goldenBounds.Dy()
	if width == 0 || height == 0 {
		return Difference{}, nil
	}

	tolerance := uint64(this.PixelTolerance)

	var totalDistance uint64
	var changedPixels uint64
	for row := range height {
		for column := range width {
			goldenRed, goldenGreen, goldenBlue, _ := golden.At(goldenBounds.Min.X+column, goldenBounds.Min.Y+row).RGBA()
			actualRed, actualGreen, actualBlue, _ := actual.At(actualBounds.Min.X+column, actualBounds.Min.Y+row).RGBA()

			redDistance := channelDistance(goldenRed, actualRed)
			greenDistance := channelDistance(goldenGreen, actualGreen)
			blueDistance := channelDistance(goldenBlue, actualBlue)

			totalDistance += redDistance + greenDistance + blueDistance
			if max(redDistance, greenDistance, blueDistance) > tolerance {
				changedPixels++
			}
		}
	}

	pixelCount := float64(width) * float64(height)
	return Difference{
		MeanDistance:         float64(totalDistance) / (pixelCount * 3 * 255),
		ChangedPixelFraction: float64(changedPixels) / pixelCount,
	}, nil
}

// Matches reports whether a difference returned by Compare is within both
// allowed thresholds.
func (this Comparer) Matches(difference Difference) bool {
	return difference.MeanDistance < this.MeanThreshold &&
		difference.ChangedPixelFraction < this.ChangedPixelThreshold
}

// Describe renders a difference next to the thresholds it was judged against,
// so a failing snapshot says which gate tripped and by how much.
func (this Comparer) Describe(difference Difference) string {
	return fmt.Sprintf(
		"mean %.4f%% (allowed < %.4f%%), changed pixels %.4f%% (allowed < %.4f%%, tolerance %d/255)",
		difference.MeanDistance*100, this.MeanThreshold*100,
		difference.ChangedPixelFraction*100, this.ChangedPixelThreshold*100,
		this.PixelTolerance)
}

// channelDistance converts the 16-bit color values returned by image.Color.RGBA
// to 8-bit and returns their absolute difference.
func channelDistance(golden, actual uint32) uint64 {
	goldenByte := golden >> 8
	actualByte := actual >> 8
	if goldenByte > actualByte {
		return uint64(goldenByte - actualByte)
	}
	return uint64(actualByte - goldenByte)
}
