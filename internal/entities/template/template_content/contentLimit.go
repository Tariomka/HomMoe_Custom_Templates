package template_content

// ContentLimit caps a single object SID at MaxCount instances, optionally per variant index.
// A few templates (e.g. "Blitz") allow caps to apply to whole content lists via `includeLists`
// in place of a single `sid`.
type ContentLimit struct {
	SID          string            `json:"sid,omitempty"`
	IncludeLists []string          `json:"includeLists,omitempty"`
	Content      []WeightedContent `json:"content,omitempty"`
	Variant      *int              `json:"variant,omitempty"`
	MaxCount     int               `json:"maxCount"`
}
