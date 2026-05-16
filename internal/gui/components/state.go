package components

import (
	"os"
	"path/filepath"
	"strings"

	"gioui.org/widget"
	"github.com/Tariomka/hommoe_custom_templates/internal/gui/utils"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/services"
)

type State struct {
	// Persistent settings file model. Updated continuously from widgets.
	settingsFile *models.SettingsFile

	// File state
	currentPath string
	dirty       bool

	// Output / status
	outputPath   widget.Editor
	lastTemplate *models.RmgTemplate
	statusMsg    string
	statusErr    bool
}

func NewUiState() *State {
	state := &State{
		settingsFile: models.NewSettingsFile(),
	}
	state.outputPath.SingleLine = true
	if workingDir, err := os.Getwd(); err == nil {
		state.outputPath.SetText(workingDir)
	}
	return state
}

func (this *State) GetStatus() (msg string, isErr bool) {
	return this.statusMsg, this.statusErr
}

func (this *State) GetSettingsFile() *models.SettingsFile {
	return this.settingsFile
}

func (this *State) GetCurrentPath() string {
	return this.currentPath
}

func (this *State) IsDirty() bool {
	return this.dirty
}

func (this *State) GetOutputPath() string {
	return this.outputPath.Text()
}

func (this *State) Reset() {
	this.settingsFile = models.NewSettingsFile()
	this.currentPath = ""
	this.dirty = false
	this.SetStatus("New settings file.", false)
}

func (this *State) SetStatus(msg string, isErr bool) {
	this.statusMsg = msg
	this.statusErr = isErr
}

func (this *State) SuggestDirectory() string {
	if this.currentPath != "" {
		return filepath.Dir(this.currentPath)
	}
	if outputDir := strings.TrimSpace(this.outputPath.Text()); outputDir != "" {
		return outputDir
	}
	workingDir, _ := os.Getwd()
	return workingDir
}

func (this *State) Load() {
	path, err := utils.PickOpenFile("Open settings", "Settings (*.gen.json)|*.gen.json|All files|*.*", this.SuggestDirectory())
	if err != nil {
		this.SetStatus("Open dialog failed: "+err.Error(), true)
		return
	}
	if path == "" {
		return
	}
	loaded, err := services.LoadSettingsFile(path)
	if err != nil {
		this.SetStatus("Load failed: "+err.Error(), true)
		return
	}
	this.settingsFile = loaded
	this.currentPath = path
	this.dirty = false
	// this.applyFromSettingsFile()
	this.SetStatus("Loaded "+path, false)
}

func (this *State) Save() {
	if this.currentPath == "" {
		this.SaveAs(this.settingsFile.TemplateName)
		return
	}
	// if err := services.SaveSettingsFile(this.currentPath, this.captureToSettingsFile()); err != nil {
	// 	this.SetStatus("Save failed: "+err.Error(), true)
	// 	return
	// }
	this.dirty = false
	this.SetStatus("Saved "+this.currentPath, false)
}

func (this *State) SaveAs(templateName string) {
	defaultName := services.SanitizeFilename(strings.TrimSpace(templateName)) + ".gen.json"
	path, err := utils.PickSaveFile("Save settings as", "Settings (*.gen.json)|*.gen.json", this.SuggestDirectory(), defaultName)
	if err != nil {
		this.SetStatus("Save dialog failed: "+err.Error(), true)
		return
	}
	if path == "" {
		return
	}
	// if err := services.SaveSettingsFile(path, this.captureToSettingsFile()); err != nil {
	// 	this.SetStatus("Save failed: "+err.Error(), true)
	// 	return
	// }
	this.currentPath = path
	this.dirty = false
	this.SetStatus("Saved "+path, false)
}

func (this *State) UpdateState(updateFunc func(*models.SettingsFile)) {
	updateFunc(this.settingsFile)
}

// // applyFromSettingsFile pushes values from this.settingsFile into all widget states.
// func (this *State) applyFromSettingsFile() {
// 	settingsFile := this.settingsFile
// 	this.templateName.SetText(settingsFile.TemplateName)
// 	this.gameMode.Selected = 0
// 	this.topology.SelectByName(topologyLabelFor(settingsFile.Topology))

// 	this.chkExpSizes.Value = settingsFile.ExperimentalMapSizes
// 	this.mapSizeSld.Value = mapSizeToSlider(settingsFile.MapSize, this.chkExpSizes.Value)
// 	this.playerCnt.Value = utils.Normalize(float32(settingsFile.PlayerCount), 2, 8)

// 	this.chkRoads.Value = settingsFile.GenerateRoads
// 	this.chkPortals.Value = settingsFile.RandomPortals
// 	this.sldMaxPortals.Value = utils.Normalize(float32(settingsFile.MaxPortalConnections), 1, 32)
// 	this.chkFootholds.Value = settingsFile.SpawnRemoteFootholds
// 	this.chkBalancedZones.Value = settingsFile.ExperimentalBalancedZonePlacement
// 	this.chkPlayerIsolation.Value = settingsFile.NoDirectPlayerConn
// 	this.chkMatchPlayerFactions.Value = settingsFile.MatchPlayerCastleFactions
// 	this.sldMinNeutralBetween.Value = utils.Normalize(float32(settingsFile.MinNeutralZonesBetweenPlayers), 0, 8)

// 	this.chkAdvancedZones.Value = settingsFile.AdvancedMode
// 	this.sldNeutralCount.Value = utils.Normalize(float32(settingsFile.NeutralZoneCount), 0, 16)
// 	this.sldPlayerCastles.Value = utils.Normalize(float32(settingsFile.PlayerZoneCastles), 0, 4)
// 	this.sldNeutralCastles.Value = utils.Normalize(float32(settingsFile.NeutralZoneCastles), 0, 4)
// 	this.sldNeutralLowNoCastle.Value = utils.Normalize(float32(settingsFile.NeutralLowNoCastleCount), 0, 8)
// 	this.sldNeutralLowCastle.Value = utils.Normalize(float32(settingsFile.NeutralLowCastleCount), 0, 8)
// 	this.sldNeutralMedNoCastle.Value = utils.Normalize(float32(settingsFile.NeutralMediumNoCastleCount), 0, 8)
// 	this.sldNeutralMedCastle.Value = utils.Normalize(float32(settingsFile.NeutralMediumCastleCount), 0, 8)
// 	this.sldNeutralHighNoCastle.Value = utils.Normalize(float32(settingsFile.NeutralHighNoCastleCount), 0, 8)
// 	this.sldNeutralHighCastle.Value = utils.Normalize(float32(settingsFile.NeutralHighCastleCount), 0, 8)
// 	this.sldHubSize.Value = float32((settingsFile.HubZoneSize - 0.5) / 1.5)
// 	this.sldHubCastles.Value = utils.Normalize(float32(settingsFile.HubZoneCastles), 0, 4)
// 	this.sldPlayerZoneSize.Value = float32((settingsFile.PlayerZoneSize - 0.5) / 1.5)
// 	this.sldNeutralZoneSize.Value = float32((settingsFile.NeutralZoneSize - 0.5) / 1.5)
// 	this.sldGuardRandom.Value = utils.Normalize(float32(settingsFile.GuardRandomization), 0, 0.5)
// 	this.sldResourceDensity.Value = utils.Normalize(float32(settingsFile.EffectiveResourceDensity()), 25, 200)
// 	this.sldStructureDensity.Value = utils.Normalize(float32(settingsFile.EffectiveStructureDensity()), 25, 200)
// 	this.sldNeutralStack.Value = utils.Normalize(float32(settingsFile.NeutralStackStrengthPercent), 25, 200)
// 	this.sldBorderGuard.Value = utils.Normalize(float32(settingsFile.BorderGuardStrengthPercent), 25, 200)

// 	// Game rules.
// 	this.victory.Selected = victoryIndex(settingsFile.VictoryCondition)
// 	this.chkLostStartCity.Value = settingsFile.LostStartCity
// 	this.sldLostCityDay.Value = utils.Normalize(float32(settingsFile.LostStartCityDay), 1, 30)
// 	this.chkLostStartHero.Value = settingsFile.LostStartHero
// 	this.chkCityHold.Value = settingsFile.CityHold
// 	this.sldCityHoldDays.Value = utils.Normalize(float32(settingsFile.CityHoldDays), 1, 30)
// 	this.chkGladiatorArena.Value = settingsFile.GladiatorArena
// 	this.sldGladiatorDelay.Value = utils.Normalize(float32(settingsFile.GladiatorArenaDaysDelayStart), 1, 90)
// 	this.sldGladiatorCountDay.Value = utils.Normalize(float32(settingsFile.GladiatorArenaCountDay), 1, 14)
// 	this.chkTournament.Value = settingsFile.Tournament
// 	this.sldTournamentDay.Value = utils.Normalize(float32(settingsFile.TournamentFirstTournamentDay), 1, 60)
// 	this.sldTournamentInterval.Value = utils.Normalize(float32(settingsFile.TournamentInterval), 1, 30)
// 	this.sldTournamentPoints.Value = utils.Normalize(float32(settingsFile.TournamentPointsToWin), 1, 10)
// 	this.chkTournamentSaveArmy.Value = settingsFile.TournamentSaveArmy
// 	this.sldHeroMin.Value = utils.Normalize(float32(settingsFile.HeroCountMin), 1, 16)
// 	this.sldHeroMax.Value = utils.Normalize(float32(settingsFile.HeroCountMax), 1, 16)
// 	this.sldHeroIncr.Value = utils.Normalize(float32(settingsFile.HeroCountIncrement), 1, 5)
// 	this.sldFactionLawsExp.Value = utils.Normalize(float32(settingsFile.FactionLawsExpPercent), 25, 200)
// 	this.sldAstrologyExp.Value = utils.Normalize(float32(settingsFile.AstrologyExpPercent), 25, 200)

// 	// Zone content.
// 	if len(settingsFile.PlayerZoneMandatoryContent) > 0 {
// 		this.applyZoneContentItems(settingsFile.PlayerZoneMandatoryContent)
// 	}
// }

// // captureToSettingsFile pulls live widget state back into this.settingsFile.
// func (this *State) captureToSettingsFile() *models.SettingsFile {
// 	settingsFile := this.settingsFile
// 	settingsFile.TemplateName = strings.TrimSpace(this.templateName.Text())
// 	settingsFile.PlayerCount = int(roundHalfAway(float64(utils.Denormalize(this.playerCnt.Value, 2, 8))))
// 	settingsFile.MapSize = sliderToMapSize(this.mapSizeSld.Value, this.chkExpSizes.Value)
// 	settingsFile.ExperimentalMapSizes = this.chkExpSizes.Value
// 	settingsFile.Topology = topologyValues[this.topology.Selected]

// 	settingsFile.GenerateRoads = this.chkRoads.Value
// 	settingsFile.RandomPortals = this.chkPortals.Value
// 	settingsFile.MaxPortalConnections = roundedRange(this.sldMaxPortals.Value, 1, 32)
// 	settingsFile.SpawnRemoteFootholds = this.chkFootholds.Value
// 	settingsFile.ExperimentalBalancedZonePlacement = this.chkBalancedZones.Value
// 	settingsFile.NoDirectPlayerConn = this.chkPlayerIsolation.Value
// 	settingsFile.MatchPlayerCastleFactions = this.chkMatchPlayerFactions.Value
// 	settingsFile.MinNeutralZonesBetweenPlayers = roundedRange(this.sldMinNeutralBetween.Value, 0, 8)

// 	settingsFile.AdvancedMode = this.chkAdvancedZones.Value
// 	settingsFile.NeutralZoneCount = roundedRange(this.sldNeutralCount.Value, 0, 16)
// 	settingsFile.PlayerZoneCastles = roundedRange(this.sldPlayerCastles.Value, 0, 4)
// 	settingsFile.NeutralZoneCastles = roundedRange(this.sldNeutralCastles.Value, 0, 4)
// 	settingsFile.NeutralLowNoCastleCount = roundedRange(this.sldNeutralLowNoCastle.Value, 0, 8)
// 	settingsFile.NeutralLowCastleCount = roundedRange(this.sldNeutralLowCastle.Value, 0, 8)
// 	settingsFile.NeutralMediumNoCastleCount = roundedRange(this.sldNeutralMedNoCastle.Value, 0, 8)
// 	settingsFile.NeutralMediumCastleCount = roundedRange(this.sldNeutralMedCastle.Value, 0, 8)
// 	settingsFile.NeutralHighNoCastleCount = roundedRange(this.sldNeutralHighNoCastle.Value, 0, 8)
// 	settingsFile.NeutralHighCastleCount = roundedRange(this.sldNeutralHighCastle.Value, 0, 8)
// 	settingsFile.HubZoneSize = float64(0.5 + this.sldHubSize.Value*1.5)
// 	settingsFile.HubZoneCastles = roundedRange(this.sldHubCastles.Value, 0, 4)
// 	settingsFile.PlayerZoneSize = float64(0.5 + this.sldPlayerZoneSize.Value*1.5)
// 	settingsFile.NeutralZoneSize = float64(0.5 + this.sldNeutralZoneSize.Value*1.5)
// 	settingsFile.GuardRandomization = float64(utils.Denormalize(this.sldGuardRandom.Value, 0, 0.5))
// 	rd := roundedRange(this.sldResourceDensity.Value, 25, 200)
// 	sd := roundedRange(this.sldStructureDensity.Value, 25, 200)
// 	settingsFile.ResourceDensityPercent = &rd
// 	settingsFile.StructureDensityPercent = &sd
// 	settingsFile.NeutralStackStrengthPercent = roundedRange(this.sldNeutralStack.Value, 25, 200)
// 	settingsFile.BorderGuardStrengthPercent = roundedRange(this.sldBorderGuard.Value, 25, 200)

// 	settingsFile.VictoryCondition = victoryIDs[this.victory.Selected]
// 	settingsFile.LostStartCity = this.chkLostStartCity.Value
// 	settingsFile.LostStartCityDay = roundedRange(this.sldLostCityDay.Value, 1, 30)
// 	settingsFile.LostStartHero = this.chkLostStartHero.Value
// 	settingsFile.CityHold = this.chkCityHold.Value || this.victory.Selected == 2
// 	settingsFile.CityHoldDays = roundedRange(this.sldCityHoldDays.Value, 1, 30)
// 	settingsFile.GladiatorArena = this.chkGladiatorArena.Value
// 	settingsFile.GladiatorArenaDaysDelayStart = roundedRange(this.sldGladiatorDelay.Value, 1, 90)
// 	settingsFile.GladiatorArenaCountDay = roundedRange(this.sldGladiatorCountDay.Value, 1, 14)
// 	settingsFile.Tournament = this.chkTournament.Value || this.victory.Selected == 3
// 	settingsFile.TournamentFirstTournamentDay = roundedRange(this.sldTournamentDay.Value, 1, 60)
// 	settingsFile.TournamentInterval = roundedRange(this.sldTournamentInterval.Value, 1, 30)
// 	settingsFile.TournamentPointsToWin = roundedRange(this.sldTournamentPoints.Value, 1, 10)
// 	settingsFile.TournamentSaveArmy = this.chkTournamentSaveArmy.Value
// 	settingsFile.HeroCountMin = roundedRange(this.sldHeroMin.Value, 1, 16)
// 	settingsFile.HeroCountMax = max(roundedRange(this.sldHeroMax.Value, 1, 16), settingsFile.HeroCountMin)
// 	settingsFile.HeroCountIncrement = roundedRange(this.sldHeroIncr.Value, 1, 5)
// 	settingsFile.FactionLawsExpPercent = roundedRange(this.sldFactionLawsExp.Value, 25, 200)
// 	settingsFile.AstrologyExpPercent = roundedRange(this.sldAstrologyExp.Value, 25, 200)

// 	settingsFile.PlayerZoneMandatoryContent = this.collectZoneContentItems()
// 	return settingsFile
// }
