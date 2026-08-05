package models

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/preview"
)

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
// differ from the cached ones. A failed build is not cached, so the next call
// retries it.
func (this *PreviewLayoutCache) Get(
	templateRevision uint64,
	topology config.MapTopology,
	canvasSide float64,
	build func() (preview.Layout, error)) (preview.Layout, error) {
	key := previewLayoutCacheKey{templateRevision, topology, canvasSide}
	if this.hasEntry && this.key == key {
		return this.layout, nil
	}

	layout, err := build()
	if err != nil {
		return preview.Layout{}, err
	}

	this.key = key
	this.layout = layout
	this.hasEntry = true
	return layout, nil
}
