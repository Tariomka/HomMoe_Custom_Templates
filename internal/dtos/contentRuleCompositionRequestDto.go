package dtos

// ContentRuleCompositionRequestDto carries the manage-rules editor state that a
// content rule is composed from: the selected rule type plus the value the
// matching editor control currently holds.
type ContentRuleCompositionRequestDto struct {
	Option          ContentRuleOptionDto
	DistanceNames   []string
	DistanceIndex   int
	IsGuarded       bool
	IsSoloEncounter bool
	VariantIDs      []int
	VariantIndex    int
}
