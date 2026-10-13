package editor_state_v2

// TemplateIdentity names the template and the game mode it targets.
type TemplateIdentity struct {
	TemplateName string `json:"templateName"`
	GameMode     string `json:"gameMode"`
}
