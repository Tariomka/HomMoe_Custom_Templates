package editor_state

type ManualEditSettings struct {
	Zones       []ManualZoneEdit       `json:"zones,omitempty"`
	Connections []ManualConnectionEdit `json:"connections,omitempty"`
}
