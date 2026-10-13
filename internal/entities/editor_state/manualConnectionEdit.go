package editor_state

type ManualConnectionEdit struct {
	Name           string `json:"name"`
	From           string `json:"from"`
	To             string `json:"to"`
	ConnectionType string `json:"type"`
	IsUserAdded    bool   `json:"isUserAdded,omitempty"`

	FieldOverrides ConnectionFieldOverrides `json:"fieldOverrides,omitzero"`
	Unset          []string                 `json:"unset,omitempty"`
}
