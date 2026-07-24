package dtos

type ContentRuleKey string

const (
	ContentRuleKeyDistanceToRoad ContentRuleKey = "distance-to-road"
	ContentRuleKeyDistanceToTown ContentRuleKey = "distance-to-town"
	ContentRuleKeyGuarded        ContentRuleKey = "guarded"
	ContentRuleKeySoloEncounter  ContentRuleKey = "solo-encounter"
	ContentRuleKeyVariant        ContentRuleKey = "variant"
)
