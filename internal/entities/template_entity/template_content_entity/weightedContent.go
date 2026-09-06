package template_content_entity

// WeightedContent is a sid + weight pair used inside MandatoryContentItem.Content rosters.
type WeightedContent struct {
	SID    string `json:"sid"`
	Weight int    `json:"weight"`
}
