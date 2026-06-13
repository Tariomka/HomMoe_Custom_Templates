package template_content

// MandatoryContent is a named group of objects that must be placed somewhere on the map.
type MandatoryContent struct {
	Name    string                 `json:"name"`
	Content []MandatoryContentItem `json:"content"`
}
