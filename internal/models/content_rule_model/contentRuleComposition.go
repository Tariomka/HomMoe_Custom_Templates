package content_rule_model

// ContentRuleComposition is the manage-rules editor state a content rule is
// composed from: the selected rule type plus the value its editor control holds.
type ContentRuleComposition struct {
	Key             ContentRuleKey
	Name            string
	DistanceNames   []string
	DistanceIndex   int
	IsGuarded       bool
	IsSoloEncounter bool
	VariantIDs      []int
	VariantIndex    int
}
