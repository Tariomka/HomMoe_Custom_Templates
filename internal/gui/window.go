package gui

import "github.com/Tariomka/hommoe_custom_templates/internal/gui/components"

type Window struct {
	state *components.State

	tabs        []*components.Tab
	selectedTab int

	toolbar      *components.Toolbar
	previewPanel *components.PreviewPanel
}

func NewWindow() *Window {
	window := Window{
		state: components.NewUiState(),
		// tabs: []*components.Tab{
		// 	components.NewTab("General", components.NewGeneralTab),
		// 	components.NewTab("Appearance", components.NewAppearanceTab),
		// },
	}

	return &window
}
