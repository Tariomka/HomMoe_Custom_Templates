package snapshot

import (
	"fmt"
	"image"
)

// DefaultSnapshotThreshold is the maximum allowed normalized mean color
// distance between a golden snapshot and an actual screenshot (2%).
// Pipeline has discrepancies, I don't want to investigate them right now.
const DefaultSnapshotThreshold = 0.005

// Comparer measures how different two screenshots are using a
// normalized mean color distance over the RGB channels (alpha is ignored:
// screenshots are opaque). 0 means identical, 1 means maximally different.
type Comparer struct {
	Threshold float64
}

// NewComparer builds a comparer with the default 2% threshold.
func NewComparer() Comparer {
	return Comparer{Threshold: DefaultSnapshotThreshold}
}

// Compare returns the normalized mean color distance between the two images:
// sum(|deltaR|+|deltaG|+|deltaB|) / (width*height*3*255). It errors when the
// image dimensions differ (callers treat that as a validation failure).
func (this Comparer) Compare(golden, actual image.Image) (float64, error) {
	goldenBounds := golden.Bounds()
	actualBounds := actual.Bounds()
	if goldenBounds.Dx() != actualBounds.Dx() || goldenBounds.Dy() != actualBounds.Dy() {
		return 0, fmt.Errorf(
			"snapshot dimensions differ: golden %dx%d, actual %dx%d",
			goldenBounds.Dx(), goldenBounds.Dy(), actualBounds.Dx(), actualBounds.Dy())
	}

	width := goldenBounds.Dx()
	height := goldenBounds.Dy()
	if width == 0 || height == 0 {
		return 0, nil
	}

	var totalDistance uint64
	for row := range height {
		for column := range width {
			goldenRed, goldenGreen, goldenBlue, _ := golden.At(goldenBounds.Min.X+column, goldenBounds.Min.Y+row).RGBA()
			actualRed, actualGreen, actualBlue, _ := actual.At(actualBounds.Min.X+column, actualBounds.Min.Y+row).RGBA()
			totalDistance += channelDistance(goldenRed, actualRed)
			totalDistance += channelDistance(goldenGreen, actualGreen)
			totalDistance += channelDistance(goldenBlue, actualBlue)
		}
	}

	maxDistance := float64(width) * float64(height) * 3 * 255
	return float64(totalDistance) / maxDistance, nil
}

// Matches reports whether a difference returned by Compare is within the
// allowed threshold.
func (this Comparer) Matches(difference float64) bool {
	return difference < this.Threshold
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
