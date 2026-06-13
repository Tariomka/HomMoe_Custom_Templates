package template_variant

// TypedRef is a generic polymorphic reference shape used across the schema for
// biome selectors, factions, road endpoints, placement rules and similar lookups.
// Observed `type` values include: "FromList", "MatchMainObject", "MatchZone",
// "Match", "MainObject", "Connection", "Crossroads", "MandatoryContent", "Road".
type TypedRef struct {
	Type string   `json:"type"`
	Args []string `json:"args"`
}
