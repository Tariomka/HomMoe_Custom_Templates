package models

// ContentRuleRowSave is the lightweight, JSON-friendly representation of a
// single content rule attached to a ZoneContentRowSave. It mirrors the C#
// `ContentRuleRowSave` record: every field is optional so that each concrete
// rule only serialises the data it actually needs.
type ContentRuleRowSave struct {
	// Name identifies which rule this row represents (e.g. "Guarded",
	// "Distance to road", "Distance to town", "Solo Encounter", "Variant").
	Name string `json:"name,omitempty"`
	// DistanceName is the distance-variation label used by the distance rules
	// ("Next To", "Near", "Medium", "Far", "Very Far").
	DistanceName string `json:"distanceName,omitempty"`
	// IsGuarded is set by the Guarded rule.
	IsGuarded *bool `json:"isGuarded,omitempty"`
	// IsSoloEncounter is set by the Solo Encounter rule.
	IsSoloEncounter *bool `json:"isSoloEncounter,omitempty"`
	// VariantId is set by the Variant rule.
	VariantId *int `json:"variantId,omitempty"`
}
