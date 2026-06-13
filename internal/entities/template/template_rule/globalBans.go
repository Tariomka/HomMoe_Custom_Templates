package template_rule

// GlobalBans declares globally banned content (items, magics, heroes) at the template level.
type GlobalBans struct {
	Items  []string `json:"items,omitempty"`
	Magics []string `json:"magics,omitempty"`
	Heroes []string `json:"heroes,omitempty"`
}
