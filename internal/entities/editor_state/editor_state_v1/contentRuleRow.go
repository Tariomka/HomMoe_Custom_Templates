package editor_state_v1

// ContentRuleRow is the frozen v1 snapshot - see editorState.go. Never edit.
type ContentRuleRow struct {
	Name            string `json:"name,omitempty"`
	DistanceName    string `json:"distanceName,omitempty"`
	IsGuarded       *bool  `json:"isGuarded,omitempty"`
	IsSoloEncounter *bool  `json:"isSoloEncounter,omitempty"`
	VariantID       *int   `json:"variantId,omitempty"`
}
