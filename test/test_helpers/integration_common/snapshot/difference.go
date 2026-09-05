package snapshot

import "fmt"

// Difference is what Comparer.Compare measured between a golden snapshot and an
// actual screenshot. The two numbers answer different questions: MeanDistance
// asks "how different is the frame overall", ChangedPixelFraction asks "how much
// of the frame changed noticeably".
type Difference struct {
	// MeanDistance is sum(|deltaR|+|deltaG|+|deltaB|) over every pixel divided
	// by width*height*3*255. 0 means identical, 1 means maximally different.
	MeanDistance float64
	// ChangedPixelFraction is the share of pixels whose largest single-channel
	// delta exceeded the comparer's PixelTolerance.
	ChangedPixelFraction float64
}

// String renders both measurements as percentages.
func (this Difference) String() string {
	return fmt.Sprintf(
		"mean %.4f%%, changed pixels %.4f%%",
		this.MeanDistance*100, this.ChangedPixelFraction*100)
}
