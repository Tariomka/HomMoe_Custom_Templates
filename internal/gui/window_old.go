package gui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/Tariomka/hommoe_custom_templates/internal/constants"
	"github.com/Tariomka/hommoe_custom_templates/internal/gui/components/widgets"
	"github.com/Tariomka/hommoe_custom_templates/internal/gui/utils"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/generator"
	"github.com/Tariomka/hommoe_custom_templates/internal/services"
)

var (
	gameModes      = []string{"Classic", "SingleHero"}
	mapSizes       = []int{64, 80, 96, 112, 128, 144, 160, 176, 192, 208, 240}
	expMapSizes    = []int{256, 272, 288, 304, 320, 336, 352, 368, 384, 400, 416, 432, 448, 464, 480, 496, 512}
	topologyLabels = []string{"Random", "Ring", "Hub", "Chain", "Shared Web"}
	topologyValues = []models.MapTopology{
		generator.TopologyRandom,
		generator.TopologyDefault,
		generator.TopologyHubAndSpoke,
		generator.TopologyChain,
		generator.TopologySharedWeb,
	}
	victoryLabels = []string{"Standard", "Lost Starting City", "Hold City", "Tournament"}
	victoryIDs    = []string{"win_condition_1", "win_condition_3", "win_condition_5", "win_condition_6"}
	roadDistances = []string{"Any", "Next To", "Near", "Medium", "Far", "Very Far"}

	mainTabLabels = []string{
		"Map Setup",
		"Generation Options",
		"Game Rules",
		"Zone Content (EXP)",
	}
)

// WindowOld holds all interactive widget state and the current SettingsFile.
type WindowOld struct {
	// File state.
	currentPath string
	dirty       bool

	// Top-level layout.
	tabs    *tabs
	scrolls [4]widget.List

	// Toolbar buttons.
	btnNew       widget.Clickable
	btnOpen      widget.Clickable
	btnSave      widget.Clickable
	btnSaveAs    widget.Clickable
	btnTemplates widget.Clickable

	// Footer buttons.
	btnGenerate     widget.Clickable
	btnSaveTemplate widget.Clickable
	btnPickOutput   widget.Clickable
	btnRevealOutput widget.Clickable

	// — Tab 1: Map Setup —
	templateName widget.Editor
	gameMode     *segmentGroup
	playerCnt    widget.Float
	mapSizeSld   widget.Float
	chkExpSizes  widget.Bool
	topology     *comboBox

	// — Tab 2: Generation Options —
	chkRoads               widget.Bool
	chkPortals             widget.Bool
	sldMaxPortals          widget.Float
	chkFootholds           widget.Bool
	chkBalancedZones       widget.Bool
	chkPlayerIsolation     widget.Bool
	chkMatchPlayerFactions widget.Bool
	sldMinNeutralBetween   widget.Float

	chkAdvancedZones       widget.Bool
	sldNeutralLowNoCastle  widget.Float
	sldNeutralLowCastle    widget.Float
	sldNeutralMedNoCastle  widget.Float
	sldNeutralMedCastle    widget.Float
	sldNeutralHighNoCastle widget.Float
	sldNeutralHighCastle   widget.Float
	sldNeutralCount        widget.Float
	sldPlayerCastles       widget.Float
	sldNeutralCastles      widget.Float
	sldHubSize             widget.Float
	sldHubCastles          widget.Float
	sldPlayerZoneSize      widget.Float
	sldNeutralZoneSize     widget.Float
	sldGuardRandom         widget.Float
	sldResourceDensity     widget.Float
	sldStructureDensity    widget.Float
	sldNeutralStack        widget.Float
	sldBorderGuard         widget.Float

	// — Tab 3: Game Rules —
	victory               *comboBox
	chkLostStartCity      widget.Bool
	sldLostCityDay        widget.Float
	chkLostStartHero      widget.Bool
	chkCityHold           widget.Bool
	sldCityHoldDays       widget.Float
	chkGladiatorArena     widget.Bool
	sldGladiatorDelay     widget.Float
	sldGladiatorCountDay  widget.Float
	chkTournament         widget.Bool
	sldTournamentDay      widget.Float
	sldTournamentInterval widget.Float
	sldTournamentPoints   widget.Float
	chkTournamentSaveArmy widget.Bool
	sldHeroMin            widget.Float
	sldHeroMax            widget.Float
	sldHeroIncr           widget.Float
	sldFactionLawsExp     widget.Float
	sldAstrologyExp       widget.Float

	// — Tab 4: Zone Content —
	zcMines      *zoneContentSection
	zcTreasures  *zoneContentSection
	zcHires      *zoneContentSection
	zcBanks      *zoneContentSection
	btnZoneReset widget.Clickable

	// Output / status.
	outputPath   widget.Editor
	lastTemplate *models.RmgTemplate
	statusMsg    string
	statusErr    bool

	// Preview panel.
	preview        previewState
	btnSavePreview widget.Clickable

	// Persistent settings file model. Updated continuously from widgets.
	settingsFile *models.SettingsFile
}

func NewState() *WindowOld {
	state := &WindowOld{
		tabs:         newTabs(mainTabLabels),
		gameMode:     newSegmentGroup(gameModes),
		topology:     newComboBox(topologyLabels),
		victory:      newComboBox(victoryLabels),
		settingsFile: models.NewSettingsFile(),
		zcMines:      newZoneContentSection("Mines", constants.ContentItemGroup.Mines, 3, true),
		zcTreasures:  newZoneContentSection("Treasures", constants.ContentItemGroup.Treasures, 10, false),
		zcHires:      newZoneContentSection("Random Hires", constants.ContentItemGroup.HireBuildings, 10, false),
		zcBanks:      newZoneContentSection("Resource Banks", constants.ContentItemGroup.ResourceBanks, 10, false),
	}
	for i := range state.scrolls {
		state.scrolls[i].Axis = layout.Vertical
	}
	state.templateName.SingleLine = true
	state.outputPath.SingleLine = true
	if workingDir, err := os.Getwd(); err == nil {
		state.outputPath.SetText(workingDir)
	}

	// Seed default in-memory zone content (mirrors C# InitializeDefaultPlayerZoneContents).
	state.seedDefaultPlayerZoneContent()
	state.applyFromSettingsFile()
	return state
}

// applyFromSettingsFile pushes values from this.settingsFile into all widget states.
func (this *WindowOld) applyFromSettingsFile() {
	settingsFile := this.settingsFile
	this.templateName.SetText(settingsFile.TemplateName)
	this.gameMode.Selected = 0
	this.topology.SelectByName(topologyLabelFor(settingsFile.Topology))

	this.chkExpSizes.Value = settingsFile.ExperimentalMapSizes
	this.mapSizeSld.Value = mapSizeToSlider(settingsFile.MapSize, this.chkExpSizes.Value)
	this.playerCnt.Value = utils.Normalize(float32(settingsFile.PlayerCount), 2, 8)

	this.chkRoads.Value = settingsFile.GenerateRoads
	this.chkPortals.Value = settingsFile.RandomPortals
	this.sldMaxPortals.Value = utils.Normalize(float32(settingsFile.MaxPortalConnections), 1, 32)
	this.chkFootholds.Value = settingsFile.SpawnRemoteFootholds
	this.chkBalancedZones.Value = settingsFile.ExperimentalBalancedZonePlacement
	this.chkPlayerIsolation.Value = settingsFile.NoDirectPlayerConn
	this.chkMatchPlayerFactions.Value = settingsFile.MatchPlayerCastleFactions
	this.sldMinNeutralBetween.Value = utils.Normalize(float32(settingsFile.MinNeutralZonesBetweenPlayers), 0, 8)

	this.chkAdvancedZones.Value = settingsFile.AdvancedMode
	this.sldNeutralCount.Value = utils.Normalize(float32(settingsFile.NeutralZoneCount), 0, 16)
	this.sldPlayerCastles.Value = utils.Normalize(float32(settingsFile.PlayerZoneCastles), 0, 4)
	this.sldNeutralCastles.Value = utils.Normalize(float32(settingsFile.NeutralZoneCastles), 0, 4)
	this.sldNeutralLowNoCastle.Value = utils.Normalize(float32(settingsFile.NeutralLowNoCastleCount), 0, 8)
	this.sldNeutralLowCastle.Value = utils.Normalize(float32(settingsFile.NeutralLowCastleCount), 0, 8)
	this.sldNeutralMedNoCastle.Value = utils.Normalize(float32(settingsFile.NeutralMediumNoCastleCount), 0, 8)
	this.sldNeutralMedCastle.Value = utils.Normalize(float32(settingsFile.NeutralMediumCastleCount), 0, 8)
	this.sldNeutralHighNoCastle.Value = utils.Normalize(float32(settingsFile.NeutralHighNoCastleCount), 0, 8)
	this.sldNeutralHighCastle.Value = utils.Normalize(float32(settingsFile.NeutralHighCastleCount), 0, 8)
	this.sldHubSize.Value = float32((settingsFile.HubZoneSize - 0.5) / 1.5)
	this.sldHubCastles.Value = utils.Normalize(float32(settingsFile.HubZoneCastles), 0, 4)
	this.sldPlayerZoneSize.Value = float32((settingsFile.PlayerZoneSize - 0.5) / 1.5)
	this.sldNeutralZoneSize.Value = float32((settingsFile.NeutralZoneSize - 0.5) / 1.5)
	this.sldGuardRandom.Value = utils.Normalize(float32(settingsFile.GuardRandomization), 0, 0.5)
	this.sldResourceDensity.Value = utils.Normalize(float32(settingsFile.EffectiveResourceDensity()), 25, 200)
	this.sldStructureDensity.Value = utils.Normalize(float32(settingsFile.EffectiveStructureDensity()), 25, 200)
	this.sldNeutralStack.Value = utils.Normalize(float32(settingsFile.NeutralStackStrengthPercent), 25, 200)
	this.sldBorderGuard.Value = utils.Normalize(float32(settingsFile.BorderGuardStrengthPercent), 25, 200)

	// Game rules.
	this.victory.Selected = victoryIndex(settingsFile.VictoryCondition)
	this.chkLostStartCity.Value = settingsFile.LostStartCity
	this.sldLostCityDay.Value = utils.Normalize(float32(settingsFile.LostStartCityDay), 1, 30)
	this.chkLostStartHero.Value = settingsFile.LostStartHero
	this.chkCityHold.Value = settingsFile.CityHold
	this.sldCityHoldDays.Value = utils.Normalize(float32(settingsFile.CityHoldDays), 1, 30)
	this.chkGladiatorArena.Value = settingsFile.GladiatorArena
	this.sldGladiatorDelay.Value = utils.Normalize(float32(settingsFile.GladiatorArenaDaysDelayStart), 1, 90)
	this.sldGladiatorCountDay.Value = utils.Normalize(float32(settingsFile.GladiatorArenaCountDay), 1, 14)
	this.chkTournament.Value = settingsFile.Tournament
	this.sldTournamentDay.Value = utils.Normalize(float32(settingsFile.TournamentFirstTournamentDay), 1, 60)
	this.sldTournamentInterval.Value = utils.Normalize(float32(settingsFile.TournamentInterval), 1, 30)
	this.sldTournamentPoints.Value = utils.Normalize(float32(settingsFile.TournamentPointsToWin), 1, 10)
	this.chkTournamentSaveArmy.Value = settingsFile.TournamentSaveArmy
	this.sldHeroMin.Value = utils.Normalize(float32(settingsFile.HeroCountMin), 1, 16)
	this.sldHeroMax.Value = utils.Normalize(float32(settingsFile.HeroCountMax), 1, 16)
	this.sldHeroIncr.Value = utils.Normalize(float32(settingsFile.HeroCountIncrement), 1, 5)
	this.sldFactionLawsExp.Value = utils.Normalize(float32(settingsFile.FactionLawsExpPercent), 25, 200)
	this.sldAstrologyExp.Value = utils.Normalize(float32(settingsFile.AstrologyExpPercent), 25, 200)

	// Zone content.
	if len(settingsFile.PlayerZoneMandatoryContent) > 0 {
		this.applyZoneContentItems(settingsFile.PlayerZoneMandatoryContent)
	}
}

// captureToSettingsFile pulls live widget state back into this.settingsFile.
func (this *WindowOld) captureToSettingsFile() *models.SettingsFile {
	settingsFile := this.settingsFile
	settingsFile.TemplateName = strings.TrimSpace(this.templateName.Text())
	settingsFile.PlayerCount = int(roundHalfAway(float64(utils.Denormalize(this.playerCnt.Value, 2, 8))))
	settingsFile.MapSize = sliderToMapSize(this.mapSizeSld.Value, this.chkExpSizes.Value)
	settingsFile.ExperimentalMapSizes = this.chkExpSizes.Value
	settingsFile.Topology = topologyValues[this.topology.Selected]

	settingsFile.GenerateRoads = this.chkRoads.Value
	settingsFile.RandomPortals = this.chkPortals.Value
	settingsFile.MaxPortalConnections = roundedRange(this.sldMaxPortals.Value, 1, 32)
	settingsFile.SpawnRemoteFootholds = this.chkFootholds.Value
	settingsFile.ExperimentalBalancedZonePlacement = this.chkBalancedZones.Value
	settingsFile.NoDirectPlayerConn = this.chkPlayerIsolation.Value
	settingsFile.MatchPlayerCastleFactions = this.chkMatchPlayerFactions.Value
	settingsFile.MinNeutralZonesBetweenPlayers = roundedRange(this.sldMinNeutralBetween.Value, 0, 8)

	settingsFile.AdvancedMode = this.chkAdvancedZones.Value
	settingsFile.NeutralZoneCount = roundedRange(this.sldNeutralCount.Value, 0, 16)
	settingsFile.PlayerZoneCastles = roundedRange(this.sldPlayerCastles.Value, 0, 4)
	settingsFile.NeutralZoneCastles = roundedRange(this.sldNeutralCastles.Value, 0, 4)
	settingsFile.NeutralLowNoCastleCount = roundedRange(this.sldNeutralLowNoCastle.Value, 0, 8)
	settingsFile.NeutralLowCastleCount = roundedRange(this.sldNeutralLowCastle.Value, 0, 8)
	settingsFile.NeutralMediumNoCastleCount = roundedRange(this.sldNeutralMedNoCastle.Value, 0, 8)
	settingsFile.NeutralMediumCastleCount = roundedRange(this.sldNeutralMedCastle.Value, 0, 8)
	settingsFile.NeutralHighNoCastleCount = roundedRange(this.sldNeutralHighNoCastle.Value, 0, 8)
	settingsFile.NeutralHighCastleCount = roundedRange(this.sldNeutralHighCastle.Value, 0, 8)
	settingsFile.HubZoneSize = float64(0.5 + this.sldHubSize.Value*1.5)
	settingsFile.HubZoneCastles = roundedRange(this.sldHubCastles.Value, 0, 4)
	settingsFile.PlayerZoneSize = float64(0.5 + this.sldPlayerZoneSize.Value*1.5)
	settingsFile.NeutralZoneSize = float64(0.5 + this.sldNeutralZoneSize.Value*1.5)
	settingsFile.GuardRandomization = float64(utils.Denormalize(this.sldGuardRandom.Value, 0, 0.5))
	rd := roundedRange(this.sldResourceDensity.Value, 25, 200)
	sd := roundedRange(this.sldStructureDensity.Value, 25, 200)
	settingsFile.ResourceDensityPercent = &rd
	settingsFile.StructureDensityPercent = &sd
	settingsFile.NeutralStackStrengthPercent = roundedRange(this.sldNeutralStack.Value, 25, 200)
	settingsFile.BorderGuardStrengthPercent = roundedRange(this.sldBorderGuard.Value, 25, 200)

	settingsFile.VictoryCondition = victoryIDs[this.victory.Selected]
	settingsFile.LostStartCity = this.chkLostStartCity.Value
	settingsFile.LostStartCityDay = roundedRange(this.sldLostCityDay.Value, 1, 30)
	settingsFile.LostStartHero = this.chkLostStartHero.Value
	settingsFile.CityHold = this.chkCityHold.Value || this.victory.Selected == 2
	settingsFile.CityHoldDays = roundedRange(this.sldCityHoldDays.Value, 1, 30)
	settingsFile.GladiatorArena = this.chkGladiatorArena.Value
	settingsFile.GladiatorArenaDaysDelayStart = roundedRange(this.sldGladiatorDelay.Value, 1, 90)
	settingsFile.GladiatorArenaCountDay = roundedRange(this.sldGladiatorCountDay.Value, 1, 14)
	settingsFile.Tournament = this.chkTournament.Value || this.victory.Selected == 3
	settingsFile.TournamentFirstTournamentDay = roundedRange(this.sldTournamentDay.Value, 1, 60)
	settingsFile.TournamentInterval = roundedRange(this.sldTournamentInterval.Value, 1, 30)
	settingsFile.TournamentPointsToWin = roundedRange(this.sldTournamentPoints.Value, 1, 10)
	settingsFile.TournamentSaveArmy = this.chkTournamentSaveArmy.Value
	settingsFile.HeroCountMin = roundedRange(this.sldHeroMin.Value, 1, 16)
	settingsFile.HeroCountMax = max(roundedRange(this.sldHeroMax.Value, 1, 16), settingsFile.HeroCountMin)
	settingsFile.HeroCountIncrement = roundedRange(this.sldHeroIncr.Value, 1, 5)
	settingsFile.FactionLawsExpPercent = roundedRange(this.sldFactionLawsExp.Value, 25, 200)
	settingsFile.AstrologyExpPercent = roundedRange(this.sldAstrologyExp.Value, 25, 200)

	settingsFile.PlayerZoneMandatoryContent = this.collectZoneContentItems()
	return settingsFile
}

// generate runs the template generator and stores the result.
func (this *WindowOld) generate() {
	captured := this.captureToSettingsFile()
	generatorSettings := services.SettingsToGenerator(captured)
	if generatorSettings.TemplateName == "" {
		this.setStatus("Template name is required.", true)
		return
	}
	template, err := services.Generate(generatorSettings)
	if err != nil {
		this.setStatus(fmt.Sprintf("Generation failed: %value", err), true)
		this.lastTemplate = nil
		return
	}
	this.lastTemplate = template
	zoneCount := 0
	connectionCount := 0
	if len(template.Variants) > 0 {
		zoneCount = len(template.Variants[0].Zones)
		connectionCount = len(template.Variants[0].Connections)
	}
	this.setStatus(fmt.Sprintf("Generated '%s' — %d zones, %d connections.", template.Name, zoneCount, connectionCount), false)
}

// saveTemplate writes the most recently generated template as .rmg.json.
func (this *WindowOld) saveTemplate() {
	if this.lastTemplate == nil {
		this.setStatus("Nothing to save — generate a template first.", true)
		return
	}
	dir := strings.TrimSpace(this.outputPath.Text())
	if dir == "" {
		this.setStatus("Output directory is empty.", true)
		return
	}
	out, err := services.WriteTemplate(dir, this.lastTemplate)
	if err != nil {
		this.setStatus(fmt.Sprintf("Save failed: %value", err), true)
		return
	}
	this.setStatus("Saved template to "+out, false)
}

func (this *WindowOld) setStatus(msg string, isErr bool) {
	this.statusMsg = msg
	this.statusErr = isErr
}

// fileNew clears the in-memory model.
func (this *WindowOld) fileNew() {
	this.settingsFile = models.NewSettingsFile()
	this.currentPath = ""
	this.dirty = false
	this.seedDefaultPlayerZoneContent()
	this.applyFromSettingsFile()
	this.setStatus("New settings file.", false)
}

// fileOpen presents a dialog and loads the chosen .gen.json file.
func (this *WindowOld) fileOpen() {
	path, err := utils.PickOpenFile("Open settings", "Settings (*.gen.json)|*.gen.json|All files|*.*", this.suggestDir())
	if err != nil {
		this.setStatus("Open dialog failed: "+err.Error(), true)
		return
	}
	if path == "" {
		return
	}
	loaded, err := services.LoadSettingsFile(path)
	if err != nil {
		this.setStatus("Load failed: "+err.Error(), true)
		return
	}
	this.settingsFile = loaded
	this.currentPath = path
	this.dirty = false
	this.applyFromSettingsFile()
	this.setStatus("Loaded "+path, false)
}

// fileSave writes to the current path or prompts via Save As if none.
func (this *WindowOld) fileSave() {
	if this.currentPath == "" {
		this.fileSaveAs()
		return
	}
	if err := services.SaveSettingsFile(this.currentPath, this.captureToSettingsFile()); err != nil {
		this.setStatus("Save failed: "+err.Error(), true)
		return
	}
	this.dirty = false
	this.setStatus("Saved "+this.currentPath, false)
}

// fileSaveAs prompts for a destination path then writes the settings file.
func (this *WindowOld) fileSaveAs() {
	defaultName := services.SanitizeFilename(strings.TrimSpace(this.templateName.Text())) + ".gen.json"
	path, err := utils.PickSaveFile("Save settings as", "Settings (*.gen.json)|*.gen.json", this.suggestDir(), defaultName)
	if err != nil {
		this.setStatus("Save dialog failed: "+err.Error(), true)
		return
	}
	if path == "" {
		return
	}
	if err := services.SaveSettingsFile(path, this.captureToSettingsFile()); err != nil {
		this.setStatus("Save failed: "+err.Error(), true)
		return
	}
	this.currentPath = path
	this.dirty = false
	this.setStatus("Saved "+path, false)
}

func (this *WindowOld) suggestDir() string {
	if this.currentPath != "" {
		return filepath.Dir(this.currentPath)
	}
	if outputDir := strings.TrimSpace(this.outputPath.Text()); outputDir != "" {
		return outputDir
	}
	workingDir, _ := os.Getwd()
	return workingDir
}

// openTemplatesFolder opens the official Steam templates directory in Explorer.
func (this *WindowOld) openTemplatesFolder() {
	dir := utils.FindOldenEraTemplatesDir()
	if dir == "" {
		this.setStatus("Heroes Olden Era templates folder not found.", true)
		return
	}
	if err := utils.RevealInExplorer(dir); err != nil {
		this.setStatus("Open folder failed: "+err.Error(), true)
		return
	}
	this.setStatus("Opened "+dir, false)
}

// pickOutputDir presents a folder picker for the template output directory.
func (this *WindowOld) pickOutputDir() {
	cur := strings.TrimSpace(this.outputPath.Text())
	dir, err := utils.PickFolder("Select output directory", cur)
	if err != nil {
		this.setStatus("Folder dialog failed: "+err.Error(), true)
		return
	}
	if dir == "" {
		return
	}
	this.outputPath.SetText(dir)
}

// — Layout —

func (this *WindowOld) Layout(gtx layout.Context, theme *material.Theme) layout.Dimensions {
	// Process button clicks first.
	if this.btnGenerate.Clicked(gtx) {
		this.generate()
	}
	if this.btnSaveTemplate.Clicked(gtx) {
		this.saveTemplate()
	}
	if this.btnNew.Clicked(gtx) {
		this.fileNew()
	}
	if this.btnOpen.Clicked(gtx) {
		this.fileOpen()
	}
	if this.btnSave.Clicked(gtx) {
		this.fileSave()
	}
	if this.btnSaveAs.Clicked(gtx) {
		this.fileSaveAs()
	}
	if this.btnTemplates.Clicked(gtx) {
		this.openTemplatesFolder()
	}
	if this.btnPickOutput.Clicked(gtx) {
		this.pickOutputDir()
	}
	if this.btnRevealOutput.Clicked(gtx) {
		_ = utils.RevealInExplorer(strings.TrimSpace(this.outputPath.Text()))
	}
	if this.btnZoneReset.Clicked(gtx) {
		this.seedDefaultPlayerZoneContent()
	}

	fillBackground(gtx, colBackground)
	return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(widgets.NewTitleBarWidget(theme, "⚔  Olden Era — Template Generator")),
			layout.Rigid(layout.Spacer{Height: unit.Dp(6)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions { return this.layoutToolbar(gtx, theme) }),
			layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions { return this.tabs.Layout(gtx, theme) }),
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Flexed(1, widgets.NewPanelWidget(unit.Dp(0), func(gtx layout.Context) layout.Dimensions {
						return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return this.layoutActiveTab(gtx, theme)
						})
					})),
					layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Min.X = gtx.Dp(unit.Dp(380))
						gtx.Constraints.Max.X = gtx.Dp(unit.Dp(440))
						return this.layoutPreviewPanel(gtx, theme)
					}),
				)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
			layout.Rigid(this.layoutFooterWidget(theme)),
		)
	})
}

func (this *WindowOld) layoutToolbar(gtx layout.Context, theme *material.Theme) layout.Dimensions {
	row := layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}
	return row.Layout(gtx,
		layout.Rigid(widgets.NewButtonWidget(theme, "🗎 New", &this.btnNew, false)),
		layout.Rigid(widgets.NewHorizontalSpacerWidget(6)),
		layout.Rigid(widgets.NewButtonWidget(theme, "🗀 Open…", &this.btnOpen, false)),
		layout.Rigid(widgets.NewHorizontalSpacerWidget(6)),
		layout.Rigid(widgets.NewButtonWidget(theme, "🖫 Save", &this.btnSave, false)),
		layout.Rigid(widgets.NewHorizontalSpacerWidget(6)),
		layout.Rigid(widgets.NewButtonWidget(theme, "🖫 Save As…", &this.btnSaveAs, false)),
		layout.Rigid(layout.Spacer{Width: unit.Dp(12)}.Layout),
		layout.Rigid(widgets.NewButtonWidget(theme, "🗀 Open templates folder", &this.btnTemplates, false)),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Left: unit.Dp(12)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				path := this.currentPath
				if path == "" {
					path = "(unsaved)"
				}
				if this.dirty {
					path += " *"
				}
				label := material.Body2(theme, "File: "+path)
				label.Color = colTextDim
				label.TextSize = unit.Sp(11)
				label.MaxLines = 1
				label.Truncator = "…"
				label.Alignment = text.End
				return label.Layout(gtx)
			})
		}),
	)
}

func (this *WindowOld) layoutActiveTab(gtx layout.Context, theme *material.Theme) layout.Dimensions {
	idx := this.tabs.Selected
	if idx < 0 || idx >= len(this.scrolls) {
		idx = 0
	}
	list := material.List(theme, &this.scrolls[idx])
	var sections []layout.Widget
	switch idx {
	case 0:
		sections = this.tabMapSetup(theme)
	case 1:
		sections = this.tabGenerationOptions(theme)
	case 2:
		sections = this.tabGameRules(theme)
	case 3:
		sections = this.tabZoneContent(theme)
	}
	return list.Layout(gtx, len(sections), func(gtx layout.Context, i int) layout.Dimensions {
		return sections[i](gtx)
	})
}

// — Footer (Generate + Save Template + Output) —

func (this *WindowOld) layoutFooterWidget(theme *material.Theme) layout.Widget {
	return widgets.NewPanelWidget(unit.Dp(10), func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Min.X = gtx.Dp(120)
						label := material.Body1(theme, "Output dir")
						label.Color = colText
						label.TextSize = unit.Sp(13)
						return label.Layout(gtx)
					}),
					layout.Flexed(1, widgets.NewTextboxWidget(theme, &this.outputPath, "Choose folder")),
					layout.Rigid(layout.Spacer{Width: unit.Dp(6)}.Layout),
					layout.Rigid(widgets.NewButtonWidget(theme, "Browse…", &this.btnPickOutput, false)),
					layout.Rigid(layout.Spacer{Width: unit.Dp(6)}.Layout),
					layout.Rigid(widgets.NewButtonWidget(theme, "Reveal", &this.btnRevealOutput, false)),
				)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						msg := this.statusMsg
						if msg == "" {
							msg = "Ready."
						}
						col := colTextDim
						if this.statusErr {
							col = colError
						} else if this.lastTemplate != nil {
							col = colGoldBright
						}
						label := material.Body2(theme, msg)
						label.Color = col
						label.TextSize = unit.Sp(12)
						label.MaxLines = 2
						return label.Layout(gtx)
					}),
					layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Min.X = gtx.Dp(190)
						return goldButton{Text: "⚔  Generate Template", Click: &this.btnGenerate}.Layout(gtx, theme)
					}),
					layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Min.X = gtx.Dp(180)
						return goldButton{Text: "💾  Save Template", Click: &this.btnSaveTemplate, Disabled: this.lastTemplate == nil}.Layout(gtx, theme)
					}),
				)
			}),
		)
	})
}
