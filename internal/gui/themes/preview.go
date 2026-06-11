package themes

import "image/color"

// Map-preview palette. Zone tiers keep their meaningful metallic identities
// (bronze / silver / gold / spawn-green / hub-blue) so the legend stays
// readable, while the chrome around them — background, frame, connection
// lines, labels — follows the Crimson Night theme.
var (
	ColorPreviewBg    = color.NRGBA{R: 0x19, G: 0x16, B: 0x1A, A: 0xFF} // #19161A canvas
	ColorPreviewFrame = color.NRGBA{R: 0x9C, G: 0x2E, B: 0x42, A: 0xFF} // #9C2E42 crimson frame

	ColorPreviewBronzeFill = color.NRGBA{R: 0x5E, G: 0x40, B: 0x26, A: 0xFF}
	ColorPreviewBronzeEdge = color.NRGBA{R: 0xCD, G: 0x7F, B: 0x32, A: 0xFF}
	ColorPreviewSilverFill = color.NRGBA{R: 0x46, G: 0x4A, B: 0x52, A: 0xFF}
	ColorPreviewSilverEdge = color.NRGBA{R: 0xC4, G: 0xC6, B: 0xCC, A: 0xFF}
	ColorPreviewGoldFill   = color.NRGBA{R: 0x75, G: 0x58, B: 0x16, A: 0xFF}
	ColorPreviewGoldEdge   = color.NRGBA{R: 0xFF, G: 0xD2, B: 0x32, A: 0xFF}
	ColorPreviewSpawnFill  = color.NRGBA{R: 0x28, G: 0x58, B: 0x34, A: 0xFF}
	ColorPreviewSpawnEdge  = color.NRGBA{R: 0x64, G: 0xC8, B: 0x78, A: 0xFF}
	ColorPreviewHubFill    = color.NRGBA{R: 0x35, G: 0x4E, B: 0x60, A: 0xFF}
	ColorPreviewHubEdge    = color.NRGBA{R: 0x82, G: 0xB4, B: 0xC8, A: 0xFF}

	ColorPreviewDirectLine = color.NRGBA{R: 0xC4, G: 0x44, B: 0x58, A: 0xFF} // crimson roads
	ColorPreviewPortalLine = color.NRGBA{R: 0x5A, G: 0xAA, B: 0xD2, A: 0xB4} // portals stay cool blue

	ColorPreviewZoneLabel   = color.NRGBA{R: 0xF2, G: 0xEE, B: 0xF2, A: 0xFF} // zone value text
	ColorPreviewCastleBadge = color.NRGBA{R: 0xFF, G: 0xD9, B: 0x9C, A: 0xFF} // castle-count badge
)

// Zone-editor canvas accents.
var (
	ColorEditorEdgeSelected = color.NRGBA{R: 0xFF, G: 0x9E, B: 0x2E, A: 0xFF} // selected connection — amber pops against crimson roads
	ColorEditorUserAddedDot = color.NRGBA{R: 0xFF, G: 0xB3, B: 0xC0, A: 0xFF} // marker on user-added zones
	ColorEditorGuardLabel   = color.NRGBA{R: 0xEF, G: 0xE8, B: 0xD6, A: 0xFF} // guard-value text on edges
)
