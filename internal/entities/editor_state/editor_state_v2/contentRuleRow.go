package editor_state_v2

type ContentRuleRow struct {
	Name            string `json:"name,omitempty"`
	DistanceName    string `json:"distanceName,omitempty"`
	IsGuarded       *bool  `json:"isGuarded,omitempty"`
	IsSoloEncounter *bool  `json:"isSoloEncounter,omitempty"`
	VariantID       *int   `json:"variantId,omitempty"`
}
