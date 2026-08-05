package models

import "github.com/Tariomka/hommoe_custom_templates/internal/models/config"

// previewLayoutCacheKey identifies the inputs a preview layout is derived from.
// The template is represented by its revision so unchanged templates compare
// equal without walking their contents.
type previewLayoutCacheKey struct {
	templateRevision uint64
	topology         config.MapTopology
	canvasSide       float64
}
