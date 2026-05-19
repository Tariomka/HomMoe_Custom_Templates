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
	window := Window{state: components.NewUiState()}
	window.toolbar = components.NewToolbar(window.state)
	window.tabs = []*components.Tab{
		components.NewTab("Map Setup", components.NewBasicSetupPanel(window.state)),
		components.NewTab("Generation Options", components.NewGenerationPanel(window.state)),
		components.NewTab("Game Rules", components.NewRulesPanel(window.state)),
		components.NewTab("Zone Content", components.NewZoneContentPanel(window.state)),
	}
	window.previewPanel = components.NewPreviewPanel(window.state)
	return &window
}

func (this *Window) save() {
	for _, tab := range this.tabs {
		tab.SaveToState()
	}
}

func (this *Window) load() {
	for _, tab := range this.tabs {
		tab.LoadFromState()
	}
}
