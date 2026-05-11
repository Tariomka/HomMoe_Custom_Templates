package gui

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/Tariomka/hommoe_custom_templates/internal/generator"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
)

// Domain knowledge ported from KnownValues.cs and TopologyOptions in MainWindow.xaml.cs.
var (
	gameModes      = []string{"Classic", "SingleHero"}
	mapSizes       = []int{64, 80, 96, 112, 128, 144, 160, 176, 192, 208, 240}
	expMapSizes    = []int{256, 272, 288, 304, 320, 336, 352, 368, 384, 400, 416, 432, 448, 464, 480, 496, 512}
	topologyLabels = []string{"Random", "Ring", "Hub", "Chain", "Shared Web"}
	topologyValues = []models.MapTopology{
		models.TopologyRandom,
		models.TopologyDefault,
		models.TopologyHubAndSpoke,
		models.TopologyChain,
		models.TopologySharedWeb,
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

// State holds all interactive widget state and the current SettingsFile.
type State struct {
	// File state.
	currentPath string
	dirty       bool

	// Top-level layout.
	tabs    *tabs
	scrolls [4]widget.List

	// Toolbar buttons.
	btnNew        widget.Clickable
	btnOpen       widget.Clickable
	btnSave       widget.Clickable
	btnSaveAs     widget.Clickable
	btnTemplates  widget.Clickable
	btnDiscord    widget.Clickable
	btnGitHub     widget.Clickable
	btnPatchNotes widget.Clickable

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
	sf *models.SettingsFile
}

func newState() *State {
	s := &State{
		tabs:        newTabs(mainTabLabels),
		gameMode:    newSegmentGroup(gameModes),
		topology:    newComboBox(topologyLabels),
		victory:     newComboBox(victoryLabels),
		sf:          models.NewSettingsFile(),
		zcMines:     newZoneContentSection("Mines", models.ContentItemGroup.Mines, 3, true),
		zcTreasures: newZoneContentSection("Treasures", models.ContentItemGroup.Treasures, 10, false),
		zcHires:     newZoneContentSection("Random Hires", models.ContentItemGroup.HireBuildings, 10, false),
		zcBanks:     newZoneContentSection("Resource Banks", models.ContentItemGroup.ResourceBanks, 10, false),
	}
	for i := range s.scrolls {
		s.scrolls[i].Axis = layout.Vertical
	}
	s.templateName.SingleLine = true
	s.outputPath.SingleLine = true
	if wd, err := os.Getwd(); err == nil {
		s.outputPath.SetText(wd)
	}

	// Seed default in-memory zone content (mirrors C# InitializeDefaultPlayerZoneContents).
	s.seedDefaultPlayerZoneContent()
	s.applyFromSettingsFile()
	return s
}

// applyFromSettingsFile pushes values from s.sf into all widget states.
func (s *State) applyFromSettingsFile() {
	sf := s.sf
	s.templateName.SetText(sf.TemplateName)
	s.gameMode.Selected = 0
	s.topology.SelectByName(topologyLabelFor(sf.Topology))

	s.chkExpSizes.Value = sf.ExperimentalMapSizes
	s.mapSizeSld.Value = mapSizeToSlider(sf.MapSize, s.chkExpSizes.Value)
	s.playerCnt.Value = mapRangeInv(float32(sf.PlayerCount), 2, 8)

	s.chkRoads.Value = sf.GenerateRoads
	s.chkPortals.Value = sf.RandomPortals
	s.sldMaxPortals.Value = mapRangeInv(float32(sf.MaxPortalConnections), 1, 32)
	s.chkFootholds.Value = sf.SpawnRemoteFootholds
	s.chkBalancedZones.Value = sf.ExperimentalBalancedZonePlacement
	s.chkPlayerIsolation.Value = sf.NoDirectPlayerConn
	s.chkMatchPlayerFactions.Value = sf.MatchPlayerCastleFactions
	s.sldMinNeutralBetween.Value = mapRangeInv(float32(sf.MinNeutralZonesBetweenPlayers), 0, 8)

	s.chkAdvancedZones.Value = sf.AdvancedMode
	s.sldNeutralCount.Value = mapRangeInv(float32(sf.NeutralZoneCount), 0, 16)
	s.sldPlayerCastles.Value = mapRangeInv(float32(sf.PlayerZoneCastles), 0, 4)
	s.sldNeutralCastles.Value = mapRangeInv(float32(sf.NeutralZoneCastles), 0, 4)
	s.sldNeutralLowNoCastle.Value = mapRangeInv(float32(sf.NeutralLowNoCastleCount), 0, 8)
	s.sldNeutralLowCastle.Value = mapRangeInv(float32(sf.NeutralLowCastleCount), 0, 8)
	s.sldNeutralMedNoCastle.Value = mapRangeInv(float32(sf.NeutralMediumNoCastleCount), 0, 8)
	s.sldNeutralMedCastle.Value = mapRangeInv(float32(sf.NeutralMediumCastleCount), 0, 8)
	s.sldNeutralHighNoCastle.Value = mapRangeInv(float32(sf.NeutralHighNoCastleCount), 0, 8)
	s.sldNeutralHighCastle.Value = mapRangeInv(float32(sf.NeutralHighCastleCount), 0, 8)
	s.sldHubSize.Value = float32((sf.HubZoneSize - 0.5) / 1.5)
	s.sldHubCastles.Value = mapRangeInv(float32(sf.HubZoneCastles), 0, 4)
	s.sldPlayerZoneSize.Value = float32((sf.PlayerZoneSize - 0.5) / 1.5)
	s.sldNeutralZoneSize.Value = float32((sf.NeutralZoneSize - 0.5) / 1.5)
	s.sldGuardRandom.Value = mapRangeInv(float32(sf.GuardRandomization), 0, 0.5)
	s.sldResourceDensity.Value = mapRangeInv(float32(sf.EffectiveResourceDensity()), 25, 200)
	s.sldStructureDensity.Value = mapRangeInv(float32(sf.EffectiveStructureDensity()), 25, 200)
	s.sldNeutralStack.Value = mapRangeInv(float32(sf.NeutralStackStrengthPercent), 25, 200)
	s.sldBorderGuard.Value = mapRangeInv(float32(sf.BorderGuardStrengthPercent), 25, 200)

	// Game rules.
	s.victory.Selected = victoryIndex(sf.VictoryCondition)
	s.chkLostStartCity.Value = sf.LostStartCity
	s.sldLostCityDay.Value = mapRangeInv(float32(sf.LostStartCityDay), 1, 30)
	s.chkLostStartHero.Value = sf.LostStartHero
	s.chkCityHold.Value = sf.CityHold
	s.sldCityHoldDays.Value = mapRangeInv(float32(sf.CityHoldDays), 1, 30)
	s.chkGladiatorArena.Value = sf.GladiatorArena
	s.sldGladiatorDelay.Value = mapRangeInv(float32(sf.GladiatorArenaDaysDelayStart), 1, 90)
	s.sldGladiatorCountDay.Value = mapRangeInv(float32(sf.GladiatorArenaCountDay), 1, 14)
	s.chkTournament.Value = sf.Tournament
	s.sldTournamentDay.Value = mapRangeInv(float32(sf.TournamentFirstTournamentDay), 1, 60)
	s.sldTournamentInterval.Value = mapRangeInv(float32(sf.TournamentInterval), 1, 30)
	s.sldTournamentPoints.Value = mapRangeInv(float32(sf.TournamentPointsToWin), 1, 10)
	s.chkTournamentSaveArmy.Value = sf.TournamentSaveArmy
	s.sldHeroMin.Value = mapRangeInv(float32(sf.HeroCountMin), 1, 16)
	s.sldHeroMax.Value = mapRangeInv(float32(sf.HeroCountMax), 1, 16)
	s.sldHeroIncr.Value = mapRangeInv(float32(sf.HeroCountIncrement), 1, 5)
	s.sldFactionLawsExp.Value = mapRangeInv(float32(sf.FactionLawsExpPercent), 25, 200)
	s.sldAstrologyExp.Value = mapRangeInv(float32(sf.AstrologyExpPercent), 25, 200)

	// Zone content.
	if len(sf.PlayerZoneMandatoryContent) > 0 {
		s.applyZoneContentItems(sf.PlayerZoneMandatoryContent)
	}
}

// captureToSettingsFile pulls live widget state back into s.sf.
func (s *State) captureToSettingsFile() *models.SettingsFile {
	sf := s.sf
	sf.TemplateName = strings.TrimSpace(s.templateName.Text())
	sf.PlayerCount = int(roundHalfAway(float64(mapRange(s.playerCnt.Value, 2, 8))))
	sf.MapSize = sliderToMapSize(s.mapSizeSld.Value, s.chkExpSizes.Value)
	sf.ExperimentalMapSizes = s.chkExpSizes.Value
	sf.Topology = topologyValues[s.topology.Selected]

	sf.GenerateRoads = s.chkRoads.Value
	sf.RandomPortals = s.chkPortals.Value
	sf.MaxPortalConnections = roundedRange(s.sldMaxPortals.Value, 1, 32)
	sf.SpawnRemoteFootholds = s.chkFootholds.Value
	sf.ExperimentalBalancedZonePlacement = s.chkBalancedZones.Value
	sf.NoDirectPlayerConn = s.chkPlayerIsolation.Value
	sf.MatchPlayerCastleFactions = s.chkMatchPlayerFactions.Value
	sf.MinNeutralZonesBetweenPlayers = roundedRange(s.sldMinNeutralBetween.Value, 0, 8)

	sf.AdvancedMode = s.chkAdvancedZones.Value
	sf.NeutralZoneCount = roundedRange(s.sldNeutralCount.Value, 0, 16)
	sf.PlayerZoneCastles = roundedRange(s.sldPlayerCastles.Value, 0, 4)
	sf.NeutralZoneCastles = roundedRange(s.sldNeutralCastles.Value, 0, 4)
	sf.NeutralLowNoCastleCount = roundedRange(s.sldNeutralLowNoCastle.Value, 0, 8)
	sf.NeutralLowCastleCount = roundedRange(s.sldNeutralLowCastle.Value, 0, 8)
	sf.NeutralMediumNoCastleCount = roundedRange(s.sldNeutralMedNoCastle.Value, 0, 8)
	sf.NeutralMediumCastleCount = roundedRange(s.sldNeutralMedCastle.Value, 0, 8)
	sf.NeutralHighNoCastleCount = roundedRange(s.sldNeutralHighNoCastle.Value, 0, 8)
	sf.NeutralHighCastleCount = roundedRange(s.sldNeutralHighCastle.Value, 0, 8)
	sf.HubZoneSize = float64(0.5 + s.sldHubSize.Value*1.5)
	sf.HubZoneCastles = roundedRange(s.sldHubCastles.Value, 0, 4)
	sf.PlayerZoneSize = float64(0.5 + s.sldPlayerZoneSize.Value*1.5)
	sf.NeutralZoneSize = float64(0.5 + s.sldNeutralZoneSize.Value*1.5)
	sf.GuardRandomization = float64(mapRange(s.sldGuardRandom.Value, 0, 0.5))
	rd := roundedRange(s.sldResourceDensity.Value, 25, 200)
	sd := roundedRange(s.sldStructureDensity.Value, 25, 200)
	sf.ResourceDensityPercent = &rd
	sf.StructureDensityPercent = &sd
	sf.NeutralStackStrengthPercent = roundedRange(s.sldNeutralStack.Value, 25, 200)
	sf.BorderGuardStrengthPercent = roundedRange(s.sldBorderGuard.Value, 25, 200)

	sf.VictoryCondition = victoryIDs[s.victory.Selected]
	sf.LostStartCity = s.chkLostStartCity.Value
	sf.LostStartCityDay = roundedRange(s.sldLostCityDay.Value, 1, 30)
	sf.LostStartHero = s.chkLostStartHero.Value
	sf.CityHold = s.chkCityHold.Value || s.victory.Selected == 2
	sf.CityHoldDays = roundedRange(s.sldCityHoldDays.Value, 1, 30)
	sf.GladiatorArena = s.chkGladiatorArena.Value
	sf.GladiatorArenaDaysDelayStart = roundedRange(s.sldGladiatorDelay.Value, 1, 90)
	sf.GladiatorArenaCountDay = roundedRange(s.sldGladiatorCountDay.Value, 1, 14)
	sf.Tournament = s.chkTournament.Value || s.victory.Selected == 3
	sf.TournamentFirstTournamentDay = roundedRange(s.sldTournamentDay.Value, 1, 60)
	sf.TournamentInterval = roundedRange(s.sldTournamentInterval.Value, 1, 30)
	sf.TournamentPointsToWin = roundedRange(s.sldTournamentPoints.Value, 1, 10)
	sf.TournamentSaveArmy = s.chkTournamentSaveArmy.Value
	sf.HeroCountMin = roundedRange(s.sldHeroMin.Value, 1, 16)
	sf.HeroCountMax = roundedRange(s.sldHeroMax.Value, 1, 16)
	if sf.HeroCountMax < sf.HeroCountMin {
		sf.HeroCountMax = sf.HeroCountMin
	}
	sf.HeroCountIncrement = roundedRange(s.sldHeroIncr.Value, 1, 5)
	sf.FactionLawsExpPercent = roundedRange(s.sldFactionLawsExp.Value, 25, 200)
	sf.AstrologyExpPercent = roundedRange(s.sldAstrologyExp.Value, 25, 200)

	sf.PlayerZoneMandatoryContent = s.collectZoneContentItems()
	return sf
}

// roundedRange snaps a [0,1] slider value to the nearest integer in [lo, hi].
func roundedRange(v float32, lo, hi int) int {
	r := int(roundHalfAway(float64(mapRange(v, float32(lo), float32(hi)))))
	if r < lo {
		r = lo
	}
	if r > hi {
		r = hi
	}
	return r
}

// mapSizeToSlider returns the [0,1] slider position for a map size value.
func mapSizeToSlider(size int, includeExp bool) float32 {
	all := mapSizes
	if includeExp {
		all = append(append([]int{}, mapSizes...), expMapSizes...)
	}
	for i, v := range all {
		if v == size {
			if len(all) <= 1 {
				return 0
			}
			return float32(i) / float32(len(all)-1)
		}
	}
	// Closest match.
	closest := 0
	best := 1 << 31
	for i, v := range all {
		d := v - size
		if d < 0 {
			d = -d
		}
		if d < best {
			best = d
			closest = i
		}
	}
	if len(all) <= 1 {
		return 0
	}
	return float32(closest) / float32(len(all)-1)
}

func sliderToMapSize(v float32, includeExp bool) int {
	all := mapSizes
	if includeExp {
		all = append(append([]int{}, mapSizes...), expMapSizes...)
	}
	if len(all) == 1 {
		return all[0]
	}
	idx := int(math.Round(float64(v) * float64(len(all)-1)))
	if idx < 0 {
		idx = 0
	}
	if idx >= len(all) {
		idx = len(all) - 1
	}
	return all[idx]
}

func topologyLabelFor(t models.MapTopology) string {
	for i, v := range topologyValues {
		if v == t {
			return topologyLabels[i]
		}
	}
	return topologyLabels[0]
}

func victoryIndex(id string) int {
	for i, v := range victoryIDs {
		if v == id {
			return i
		}
	}
	return 0
}

// generate runs the template generator and stores the result.
func (s *State) generate() {
	sf := s.captureToSettingsFile()
	settings := settingsFileToGenerator(sf)
	if settings.TemplateName == "" {
		s.setStatus("Template name is required.", true)
		return
	}
	tmpl, err := generator.Generate(settings)
	if err != nil {
		s.setStatus(fmt.Sprintf("Generation failed: %v", err), true)
		s.lastTemplate = nil
		return
	}
	s.lastTemplate = tmpl
	zones := 0
	conns := 0
	if len(tmpl.Variants) > 0 {
		zones = len(tmpl.Variants[0].Zones)
		conns = len(tmpl.Variants[0].Connections)
	}
	s.setStatus(fmt.Sprintf("Generated '%s' — %d zones, %d connections.", tmpl.Name, zones, conns), false)
}

// saveTemplate writes the most recently generated template as .rmg.json.
func (s *State) saveTemplate() {
	if s.lastTemplate == nil {
		s.setStatus("Nothing to save — generate a template first.", true)
		return
	}
	dir := strings.TrimSpace(s.outputPath.Text())
	if dir == "" {
		s.setStatus("Output directory is empty.", true)
		return
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		s.setStatus(fmt.Sprintf("Cannot create directory: %v", err), true)
		return
	}
	safeName := sanitizeFilename(s.lastTemplate.Name)
	if safeName == "" {
		safeName = "Generated_Template"
	}
	out := filepath.Join(dir, safeName+".rmg.json")
	data, err := json.MarshalIndent(s.lastTemplate, "", "  ")
	if err != nil {
		s.setStatus(fmt.Sprintf("Marshal failed: %v", err), true)
		return
	}
	if err := os.WriteFile(out, data, 0o644); err != nil {
		s.setStatus(fmt.Sprintf("Write failed: %v", err), true)
		return
	}
	s.setStatus("Saved template to "+out, false)
}

func sanitizeFilename(name string) string {
	bad := []rune{'/', '\\', ':', '*', '?', '"', '<', '>', '|'}
	out := []rune(strings.TrimSpace(name))
	for i, r := range out {
		for _, b := range bad {
			if r == b {
				out[i] = '_'
			}
		}
	}
	return string(out)
}

func (s *State) setStatus(msg string, isErr bool) {
	s.statusMsg = msg
	s.statusErr = isErr
}

// handleToolbar processes File toolbar clicks.
func (s *State) handleToolbar() {
	if s.btnNew.Clicked(layoutCtxNop()) {
		// no-op handled in Layout via direct check; this guard allows reuse
	}
}

func layoutCtxNop() layout.Context { return layout.Context{} }

// fileNew clears the in-memory model.
func (s *State) fileNew() {
	s.sf = models.NewSettingsFile()
	s.currentPath = ""
	s.dirty = false
	s.seedDefaultPlayerZoneContent()
	s.applyFromSettingsFile()
	s.setStatus("New settings file.", false)
}

// fileOpen presents a dialog and loads the chosen .oetgs file.
func (s *State) fileOpen() {
	path, err := PickOpenFile("Open settings", "Settings (*.oetgs)|*.oetgs|All files|*.*", s.suggestDir())
	if err != nil {
		s.setStatus("Open dialog failed: "+err.Error(), true)
		return
	}
	if path == "" {
		return
	}
	sf, err := LoadSettingsFile(path)
	if err != nil {
		s.setStatus("Load failed: "+err.Error(), true)
		return
	}
	s.sf = sf
	s.currentPath = path
	s.dirty = false
	s.applyFromSettingsFile()
	s.setStatus("Loaded "+path, false)
}

// fileSave writes to the current path or prompts via Save As if none.
func (s *State) fileSave() {
	if s.currentPath == "" {
		s.fileSaveAs()
		return
	}
	if err := SaveSettingsFile(s.currentPath, s.captureToSettingsFile()); err != nil {
		s.setStatus("Save failed: "+err.Error(), true)
		return
	}
	s.dirty = false
	s.setStatus("Saved "+s.currentPath, false)
}

// fileSaveAs prompts for a destination path then writes the settings file.
func (s *State) fileSaveAs() {
	defaultName := sanitizeFilename(strings.TrimSpace(s.templateName.Text())) + ".oetgs"
	path, err := PickSaveFile("Save settings as", "Settings (*.oetgs)|*.oetgs", s.suggestDir(), defaultName)
	if err != nil {
		s.setStatus("Save dialog failed: "+err.Error(), true)
		return
	}
	if path == "" {
		return
	}
	if err := SaveSettingsFile(path, s.captureToSettingsFile()); err != nil {
		s.setStatus("Save failed: "+err.Error(), true)
		return
	}
	s.currentPath = path
	s.dirty = false
	s.setStatus("Saved "+path, false)
}

func (s *State) suggestDir() string {
	if s.currentPath != "" {
		return filepath.Dir(s.currentPath)
	}
	if d := strings.TrimSpace(s.outputPath.Text()); d != "" {
		return d
	}
	wd, _ := os.Getwd()
	return wd
}

// openTemplatesFolder opens the official Steam templates directory in Explorer.
func (s *State) openTemplatesFolder() {
	dir := FindOldenEraTemplatesDir()
	if dir == "" {
		s.setStatus("Heroes Olden Era templates folder not found.", true)
		return
	}
	if err := RevealInExplorer(dir); err != nil {
		s.setStatus("Open folder failed: "+err.Error(), true)
		return
	}
	s.setStatus("Opened "+dir, false)
}

// pickOutputDir presents a folder picker for the template output directory.
func (s *State) pickOutputDir() {
	cur := strings.TrimSpace(s.outputPath.Text())
	dir, err := PickFolder("Select output directory", cur)
	if err != nil {
		s.setStatus("Folder dialog failed: "+err.Error(), true)
		return
	}
	if dir == "" {
		return
	}
	s.outputPath.SetText(dir)
}

// — Layout —

func (s *State) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	// Process button clicks first.
	if s.btnGenerate.Clicked(gtx) {
		s.generate()
	}
	if s.btnSaveTemplate.Clicked(gtx) {
		s.saveTemplate()
	}
	if s.btnNew.Clicked(gtx) {
		s.fileNew()
	}
	if s.btnOpen.Clicked(gtx) {
		s.fileOpen()
	}
	if s.btnSave.Clicked(gtx) {
		s.fileSave()
	}
	if s.btnSaveAs.Clicked(gtx) {
		s.fileSaveAs()
	}
	if s.btnTemplates.Clicked(gtx) {
		s.openTemplatesFolder()
	}
	if s.btnDiscord.Clicked(gtx) {
		_ = OpenURL("https://discord.gg/UqT8KshsxW")
	}
	if s.btnGitHub.Clicked(gtx) {
		_ = OpenURL("https://github.com/KhanDevelopsGames/Olden-Era---Template-Generator")
	}
	if s.btnPatchNotes.Clicked(gtx) {
		_ = OpenURL("https://github.com/KhanDevelopsGames/Olden-Era---Template-Generator/releases")
	}
	if s.btnPickOutput.Clicked(gtx) {
		s.pickOutputDir()
	}
	if s.btnRevealOutput.Clicked(gtx) {
		_ = RevealInExplorer(strings.TrimSpace(s.outputPath.Text()))
	}
	if s.btnZoneReset.Clicked(gtx) {
		s.seedDefaultPlayerZoneContent()
	}

	fillBackground(gtx, colBackground)
	return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions { return s.layoutTitleBar(gtx, th) }),
			layout.Rigid(layout.Spacer{Height: unit.Dp(6)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions { return s.layoutToolbar(gtx, th) }),
			layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions { return s.tabs.Layout(gtx, th) }),
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return borderedPanel(gtx, unit.Dp(0), func(gtx layout.Context) layout.Dimensions {
							return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								return s.layoutActiveTab(gtx, th)
							})
						})
					}),
					layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Min.X = gtx.Dp(unit.Dp(380))
						gtx.Constraints.Max.X = gtx.Dp(unit.Dp(440))
						return s.layoutPreviewPanel(gtx, th)
					}),
				)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions { return s.layoutFooter(gtx, th) }),
		)
	})
}

func (s *State) layoutTitleBar(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			lbl := material.H6(th, "⚔  Olden Era — Template Generator")
			lbl.Color = colGold
			lbl.Font = font.Font{Weight: font.SemiBold}
			return lbl.Layout(gtx)
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions { return layout.Dimensions{Size: gtx.Constraints.Min} }),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return toolbarButton{Text: "Discord", Click: &s.btnDiscord}.Layout(gtx, th)
		}),
		layout.Rigid(layout.Spacer{Width: unit.Dp(6)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return toolbarButton{Text: "GitHub", Click: &s.btnGitHub}.Layout(gtx, th)
		}),
		layout.Rigid(layout.Spacer{Width: unit.Dp(6)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return toolbarButton{Text: "Patch notes", Click: &s.btnPatchNotes}.Layout(gtx, th)
		}),
	)
}

func (s *State) layoutToolbar(gtx layout.Context, th *material.Theme) layout.Dimensions {
	row := layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}
	return row.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return toolbarButton{Text: "📄 New", Click: &s.btnNew}.Layout(gtx, th)
		}),
		layout.Rigid(layout.Spacer{Width: unit.Dp(6)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return toolbarButton{Text: "📂 Open…", Click: &s.btnOpen}.Layout(gtx, th)
		}),
		layout.Rigid(layout.Spacer{Width: unit.Dp(6)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return toolbarButton{Text: "💾 Save", Click: &s.btnSave}.Layout(gtx, th)
		}),
		layout.Rigid(layout.Spacer{Width: unit.Dp(6)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return toolbarButton{Text: "💾 Save As…", Click: &s.btnSaveAs}.Layout(gtx, th)
		}),
		layout.Rigid(layout.Spacer{Width: unit.Dp(12)}.Layout),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return toolbarButton{Text: "🗀 Open templates folder", Click: &s.btnTemplates}.Layout(gtx, th)
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{Left: unit.Dp(12)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				path := s.currentPath
				if path == "" {
					path = "(unsaved)"
				}
				if s.dirty {
					path += " *"
				}
				lbl := material.Body2(th, "File: "+path)
				lbl.Color = colTextDim
				lbl.TextSize = unit.Sp(11)
				lbl.MaxLines = 1
				lbl.Truncator = "…"
				lbl.Alignment = text.End
				return lbl.Layout(gtx)
			})
		}),
	)
}

func (s *State) layoutActiveTab(gtx layout.Context, th *material.Theme) layout.Dimensions {
	idx := s.tabs.Selected
	if idx < 0 || idx >= len(s.scrolls) {
		idx = 0
	}
	list := material.List(th, &s.scrolls[idx])
	var sections []layout.Widget
	switch idx {
	case 0:
		sections = s.tabMapSetup(th)
	case 1:
		sections = s.tabGenerationOptions(th)
	case 2:
		sections = s.tabGameRules(th)
	case 3:
		sections = s.tabZoneContent(th)
	}
	return list.Layout(gtx, len(sections), func(gtx layout.Context, i int) layout.Dimensions {
		return sections[i](gtx)
	})
}

// — Footer (Generate + Save Template + Output) —

func (s *State) layoutFooter(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return borderedPanel(gtx, unit.Dp(10), func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Min.X = gtx.Dp(120)
						lbl := material.Body1(th, "Output dir")
						lbl.Color = colText
						lbl.TextSize = unit.Sp(13)
						return lbl.Layout(gtx)
					}),
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						return drawEditor(gtx, th, &s.outputPath, "Choose folder")
					}),
					layout.Rigid(layout.Spacer{Width: unit.Dp(6)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return toolbarButton{Text: "Browse…", Click: &s.btnPickOutput}.Layout(gtx, th)
					}),
					layout.Rigid(layout.Spacer{Width: unit.Dp(6)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return toolbarButton{Text: "Reveal", Click: &s.btnRevealOutput}.Layout(gtx, th)
					}),
				)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						msg := s.statusMsg
						if msg == "" {
							msg = "Ready."
						}
						col := colTextDim
						if s.statusErr {
							col = colError
						} else if s.lastTemplate != nil {
							col = colGoldBright
						}
						lbl := material.Body2(th, msg)
						lbl.Color = col
						lbl.TextSize = unit.Sp(12)
						lbl.MaxLines = 2
						return lbl.Layout(gtx)
					}),
					layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Min.X = gtx.Dp(190)
						return goldButton{Text: "⚔  Generate Template", Click: &s.btnGenerate}.Layout(gtx, th)
					}),
					layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Min.X = gtx.Dp(180)
						return goldButton{Text: "💾  Save Template", Click: &s.btnSaveTemplate, Disabled: s.lastTemplate == nil}.Layout(gtx, th)
					}),
				)
			}),
		)
	})
}
