package drivers

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/models"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/widgets"
	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_errors"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers/handler_interfaces"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
)

const (
	configFileExtension    = ".gen.json"
	chooseOutputFolderHint = "Choose the game templates folder using the output folder picker before exporting."
)

type State struct {
	handler      handler_interfaces.IGuiHandler
	fileSystem   handler_interfaces.IFileSystemHandler
	regeneration handler_interfaces.IRegenerationHandler

	innerState *models.EditorState

	currentPath string
	unsaved     bool

	outputPath   widget.Editor
	lastTemplate *template_model.Template
	// templateRevision counts every replacement of lastTemplate, letting the
	// preview cache detect a new template without comparing its contents.
	templateRevision uint64
	statusMsg        string
	statusErr        bool

	confirmExit bool
	onExit      func() // app closing callback for Gio app.DestroyEvent flow.

	applyNextStateAt time.Time               // for debouncing
	pendingBaseZones dtos.ZoneEditorZonesDto // uncommitted layout produced by PreviewBaseZones.

	dialogs *DialogHost
}

func NewUIState(
	handler handler_interfaces.IGuiHandler,
	fileSystem handler_interfaces.IFileSystemHandler,
	regeneration handler_interfaces.IRegenerationHandler,
	findTemplateDir bool) *State {
	state := &State{
		handler:      handler,
		fileSystem:   fileSystem,
		regeneration: regeneration,
		innerState:   models.NewEditorState(handler),
	}
	state.outputPath.SingleLine = true
	state.dialogs = &DialogHost{}
	if !findTemplateDir {
		return state
	}

	state.tryResolveTemplateDirectory()
	return state
}

func (this *State) GetStatus() (msg string, isErr bool) { return this.statusMsg, this.statusErr }

func (this *State) GetDialogHost() *DialogHost { return this.dialogs }

func (this *State) GetStateData() editor_state_model.EditorState {
	return this.innerState.GetCurrentState()
}

func (this *State) GetStateDto() editor_state_dto.EditorStateDto {
	return editor_state_dto.EditorStateDto{EditorState: this.GetStateData()}
}

func (this *State) GetTemplateName() string { return this.innerState.GetTemplateName() }

func (this *State) GetMapSize() int { return this.innerState.GetMapSize() }

func (this *State) GetTopology() config.MapTopology { return this.innerState.GetTopology() }

func (this *State) GetExperimentalMapSizes() bool { return this.innerState.GetExperimentalMapSizes() }

func (this *State) GetCurrentPath() string { return this.currentPath }

func (this *State) IsUnsaved() bool { return this.unsaved }

func (this *State) GetLastTemplate() *template_model.Template { return this.lastTemplate }

func (this *State) GetTemplateRevision() uint64 { return this.templateRevision }

func (this *State) GetOutputPath() string { return this.outputPath.Text() }

func (this *State) GetOutputPathWidget(theme *material.Theme) layout.Widget {
	return widgets.NewTextboxWidget(theme, &this.outputPath, "Choose folder", true)
}

func (this *State) Reset() {
	this.innerState.ResetState()
	this.currentPath = ""
	this.unsaved = false
	this.clearGeneratedState()
	this.SetStatus("New settings file.", false)
}

func (this *State) UpdateState(updateFunc func(*editor_state_model.EditorState)) {
	this.innerState.UpdateCurrentState(updateFunc)
	if this.innerState.WasStateChanged() {
		this.flagAsUnsaved()
	}
}

func (this *State) SetStatus(msg string, isErr bool) {
	this.statusMsg = msg
	this.statusErr = isErr
}

func (this *State) tryResolveTemplateDirectory() {
	templateDir, err := this.fileSystem.FindGameTemplateDirectory()
	templateDir = strings.TrimSpace(templateDir)
	if err == nil && templateDir != "" {
		this.outputPath.SetText(templateDir)
		return
	}

	if err == nil || errors.Is(err, common_errors.ErrTemplatesDirNotFound) {
		this.SetStatus("Game template directory not found. "+chooseOutputFolderHint, true)
		return
	}

	this.SetStatus(fmt.Sprintf("Failed to find game template directory: %v. %s", err, chooseOutputFolderHint), true)
}

func (this *State) hasTemplateVariants() bool {
	return this.lastTemplate != nil && len(this.lastTemplate.Variants) > 0
}

// setLastTemplate is the only writer of lastTemplate, so templateRevision
// cannot drift away from the template the preview is showing.
func (this *State) setLastTemplate(template *template_model.Template) {
	this.lastTemplate = template
	this.templateRevision++
}

func (this *State) getPreviousStateDto() *editor_state_dto.EditorStateDto {
	state := this.innerState.GetPreviousState()
	if state == nil {
		return nil
	}

	return &editor_state_dto.EditorStateDto{EditorState: *state}
}

func (this *State) getNextStateDto() *editor_state_dto.EditorStateDto {
	state := this.innerState.GetNextState()
	if state == nil {
		return nil
	}

	return &editor_state_dto.EditorStateDto{EditorState: *state}
}

func (this *State) flagAsUnsaved() {
	this.unsaved = true
	// New edits invalidate a pending exit confirmation.
	this.confirmExit = false
}
