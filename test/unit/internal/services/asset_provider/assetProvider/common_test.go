package assetProvider_test

import (
	"image"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/preview"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/asset_provider"
	"github.com/stretchr/testify/require"
)

const (
	canvasSide   = 64
	canvasCenter = canvasSide / 2
	spriteScale  = 0.5
)

// newCanvas returns a blank transparent canvas the sprite fits into.
func newCanvas() *image.RGBA {
	return image.NewRGBA(image.Rect(0, 0, canvasSide, canvasSide))
}

// mustNewProvider fails the test immediately when the embedded assets cannot load.
func mustNewProvider(t *testing.T) *asset_provider.AssetProvider {
	t.Helper()
	provider, err := asset_provider.NewAssetProvider()
	require.NoError(t, err)
	return provider
}

// renderPlayer draws the given player zone centered on a fresh canvas and returns it.
func renderPlayer(t *testing.T, zone preview.Zone) *image.RGBA {
	t.Helper()
	canvas := newCanvas()
	mustNewProvider(t).DrawPlayerZone(canvas, zone, image.Pt(canvasCenter, canvasCenter), spriteScale)
	return canvas
}

// renderNeutral draws the given neutral zone centered on a fresh canvas and returns it.
func renderNeutral(t *testing.T, zone preview.Zone) *image.RGBA {
	t.Helper()
	canvas := newCanvas()
	mustNewProvider(t).DrawNeutralZone(canvas, zone, image.Pt(canvasCenter, canvasCenter), spriteScale)
	return canvas
}

// renderArenaMarker draws the swords marker centered on a fresh canvas and returns it.
func renderArenaMarker(t *testing.T) *image.RGBA {
	t.Helper()
	canvas := newCanvas()
	mustNewProvider(t).DrawArenaMarker(canvas, image.Pt(canvasCenter, canvasCenter), spriteScale)
	return canvas
}
