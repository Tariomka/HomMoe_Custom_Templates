package models

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/preview"
)

type previewLayoutCacheKey struct {
	templateRevision uint64
	topology         config.MapTopology
	canvasSide       float64
}

// PreviewLayoutCache memoizes the last preview layout so a redrawing canvas
// does not recompute a deterministic result from unchanged inputs. It holds a
// single entry: the inputs only change on generation, topology switches and
// resizes, so nothing is gained by remembering older layouts.
type PreviewLayoutCache struct {
	key      previewLayoutCacheKey
	layout   preview.Layout
	hasEntry bool
}

func NewPreviewLayoutCache() *PreviewLayoutCache {
	return &PreviewLayoutCache{}
}

// Get returns the layout for the given inputs, calling build only when they
// differ from the cached ones.
func (this *PreviewLayoutCache) Get(
	templateRevision uint64,
	topology config.MapTopology,
	canvasSide float64,
	build func() preview.Layout) preview.Layout {
	key := previewLayoutCacheKey{templateRevision, topology, canvasSide}
	if this.hasEntry && this.key == key {
		return this.layout
	}

	this.key = key
	this.layout = build()
	this.hasEntry = true
	return this.layout
}
