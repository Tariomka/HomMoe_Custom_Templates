package editor_state_v1

// TemplateIdentity is the frozen v1 snapshot - see editorState.go. Never edit.
type TemplateIdentity struct {
	TemplateName string `json:"templateName"`
	GameMode     string `json:"gameMode"`
}
