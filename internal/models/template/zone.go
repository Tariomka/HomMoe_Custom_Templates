package template

// Zone represents a single zone in the map
type Zone struct {
	Name              string          `json:"name"`
	Type              string          `json:"type"` // "player" or "neutral"
	Letter            string          `json:"letter"`
	Owner             int             `json:"owner,omitempty"`
	DefenseValue      int             `json:"defenseValue,omitempty"`
	Layout            ZoneLayout      `json:"layout"`
	GuardSettings     GuardSettings   `json:"guardSettings"`
	BiomeSelectors    []BiomeSelector `json:"biomeSelectors,omitempty"`
	ContentPools      ContentPools    `json:"contentPools"`
	MandatoryContents []string        `json:"mandatoryContents,omitempty"`
	MainObjects       []MainObject    `json:"mainObjects,omitempty"`
	Roads             []string        `json:"roads,omitempty"`
	WaterSlots        []WaterSlot     `json:"waterSlots,omitempty"`
	ConnectionRules   []string        `json:"connectionRules,omitempty"`
}
