package constants

import "gioui.org/unit"

const (
	DefaultLabelWidth      = 150
	DefaultLabelWidthLong  = 220 // +70 to the Default Width
	DefaultLabelWidthShort = 80  // -70 from the Default Width

	DefaultPadding      = unit.Dp(8)
	DefaultPaddingLarge = unit.Dp(10) // +2 to the Default Padding
	DefaultPaddingSmall = unit.Dp(6)  // -2 from the Default Padding

	DefaultRoundness             = unit.Dp(4)
	DefaultRoundnessLarge        = unit.Dp(6)
	DefaultRoundnessOverlineText = unit.Dp(10)

	DefaultPreviewWidthMinimum = 380
	DefaultPreviewWidthMaximum = 440
)
