// Package themes is the single source of truth for the application's visual
// identity. Every color used anywhere in the UI lives in this package.
//
// The active palette is "Crimson Night": a layered dark theme built from
// cool charcoal surfaces with a subtle violet undertone, crimson accents for
// everything interactive, and soft near-white text. Derived from the design
// reference in Display.pdf.
package themes

import "image/color"

// ─── Raw palette ────────────────────────────────────────────────────────────
// The primitive swatches of Crimson Night. Prefer the semantic names below;
// these exist so every tone is defined exactly once.

var (
	// Surfaces - darkest to lightest. Charcoal with a faint warm-violet cast
	// so the UI reads "dark" rather than "black".
	swatchNight    = color.NRGBA{R: 0x16, G: 0x14, B: 0x18, A: 0xFF} // #161418 window backdrop
	swatchCoal     = color.NRGBA{R: 0x1E, G: 0x1B, B: 0x20, A: 0xFF} // #1E1B20 raised panels
	swatchGraphite = color.NRGBA{R: 0x28, G: 0x24, B: 0x2A, A: 0xFF} // #28242A inputs, wells
	swatchSlate    = color.NRGBA{R: 0x2D, G: 0x2A, B: 0x31, A: 0xFF} // #2D2A31 buttons, chrome

	// Crimson ramp - the accent family.
	swatchCrimsonDeep   = color.NRGBA{R: 0x77, G: 0x1A, B: 0x2C, A: 0xFF} // #771A2C filled primary surfaces
	swatchCrimsonDim    = color.NRGBA{R: 0x8E, G: 0x22, B: 0x37, A: 0xFF} // #8E2237 muted accent
	swatchCrimson       = color.NRGBA{R: 0xD9, G: 0x30, B: 0x4C, A: 0xFF} // #D9304C the signature accent
	swatchCrimsonBright = color.NRGBA{R: 0xFF, G: 0x5E, B: 0x76, A: 0xFF} // #FF5E76 highlights, active text
	swatchWine          = color.NRGBA{R: 0x5C, G: 0x24, B: 0x33, A: 0xFF} // #5C2433 selection fills
	swatchMulberry      = color.NRGBA{R: 0x4A, G: 0x32, B: 0x3A, A: 0xFF} // #4A323A smoky crimson borders

	// Ink - text tones.
	swatchInk    = color.NRGBA{R: 0xEE, G: 0xEA, B: 0xEE, A: 0xFF} // #EEEAEE primary text
	swatchInkDim = color.NRGBA{R: 0x9E, G: 0x96, B: 0xA2, A: 0xFF} // #9E96A2 secondary text

	// Signals - kept away from crimson so they stay legible next to accents.
	swatchEmber      = color.NRGBA{R: 0xFF, G: 0x96, B: 0x4B, A: 0xFF} // #FF964B errors
	swatchAmber      = color.NRGBA{R: 0xE0, G: 0xAC, B: 0x4E, A: 0xFF} // #E0AC4E warnings
	swatchAmberDark  = color.NRGBA{R: 0x6A, G: 0x50, B: 0x24, A: 0xFF} // #6A5024 warning borders
	swatchAmberShade = color.NRGBA{R: 0x2A, G: 0x21, B: 0x1A, A: 0xFF} // #2A211A warning fills
)

// ─── Semantic theme colors ──────────────────────────────────────────────────
// Use these throughout the UI; they map intent onto the raw palette.

var (
	// Surfaces.
	ColorBackground = swatchNight    // window background
	ColorPanel      = swatchCoal     // panels, cards, dialogs
	ColorInput      = swatchGraphite // text inputs, dropdown bodies, wells
	ColorButton     = swatchSlate    // standard button fill
	ColorBorder     = swatchMulberry // default border / outline

	// Accent - crimson everywhere something is interactive or highlighted.
	ColorAccent       = swatchCrimson       // headers, focus, active borders
	ColorAccentDim    = swatchCrimsonDim    // de-emphasised accent
	ColorAccentBright = swatchCrimsonBright // active labels, hot highlights

	// Text.
	ColorText    = swatchInk
	ColorTextDim = swatchInkDim

	// Status.
	ColorError          = swatchEmber
	ColorWarnBackground = swatchAmberShade
	ColorWarnBorder     = swatchAmberDark
	ColorWarnText       = swatchAmber

	// Controls.
	ColorPrimaryButton  = swatchCrimsonDeep                               // generate / confirm fill
	ColorButtonDisabled = color.NRGBA{R: 0x2A, G: 0x23, B: 0x26, A: 0xFF} // #2A2326 disabled fill
	ColorBorderDisabled = color.NRGBA{R: 0x41, G: 0x37, B: 0x3B, A: 0xFF} // #41373B disabled outline
	ColorSelection      = swatchWine                                      // checked picker rows, dropdown selection
	ColorHover          = swatchCrimsonDim                                // hovered control border highlight
	ColorScrim          = color.NRGBA{A: 0xB0}                            // modal backdrop dimmer
)
