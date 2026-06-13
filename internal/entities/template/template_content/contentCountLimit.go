package template_content

// ContentCountLimit is a named cap on how many of certain SIDs may appear.
type ContentCountLimit struct {
	Name   string         `json:"name"`
	Limits []ContentLimit `json:"limits"`
}
