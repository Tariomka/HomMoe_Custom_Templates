package editor_state_v1

// ZoneContentRow is the frozen v1 snapshot - see editorState.go. Never edit.
type ZoneContentRow struct {
	Sid     string `json:"sid"`
	Count   int    `json:"count"`
	IsGroup bool   `json:"isGroup"`
	IsMine  bool   `json:"isMine,omitempty"`

	Rules []ContentRuleRow `json:"rules,omitempty"`
}
