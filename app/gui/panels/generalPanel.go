package panels

import (
	"fmt"
	"slices"
	"strings"

	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/components"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/constants"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/drivers"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/utils"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/widgets"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/linq"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
)

type GeneralPanel struct {
	// Template section

	templateName      widget.Editor
	playerCount       widget.Float
	mapSizeSelector   *components.DropdownSelector
	checkMoreMapSizes widget.Bool

	// Hero section

	gameMode               *components.SegmentButtonGroup
	heroMinimumCount       widget.Float
	heroMaximumCount       widget.Float
	heroIncrementPerCastle widget.Float

	// Rules section

	factionLawXpMultiplier widget.Float
	astrologyXpMultiplier  widget.Float

	victorySelector *components.DropdownSelector

	checkLostStartCity widget.Bool
	lostCityDayCount   widget.Float

	checkLostStartHero widget.Bool

	checkCityHold     widget.Bool
	cityHoldDaysCount widget.Float

	checkGladiatorArena widget.Bool
	gladiatorDelayCount widget.Float
	gladiatorDayCount   widget.Float

	checkTournament         widget.Bool
	tournamentDayCount      widget.Float
	tournamentIntervalCount widget.Float
	tournamentPointsCount   widget.Float
	checkTournamentSaveArmy widget.Bool

	// The rest of the stuff

	scroll widget.List

	state *drivers.State
}

func NewGeneralPanel(state *drivers.State) *GeneralPanel {
	panel := &GeneralPanel{
		templateName: widget.Editor{SingleLine: true},
		gameMode:     components.NewSegmentButtonGroup(registry.GetGameModeList()),
		victorySelector: components.NewDropdownSelector(func() []string {
			labels := make([]string, 0)
			for _, victory := range constants.GetVictoryConditionList() {
				labels = append(labels, victory.Label)
			}
			return labels
		}()),
		mapSizeSelector: components.NewDropdownSelector(func() []string {
			return linq.FromSlice(constants.GetMapSizes(state.GetStateData().ExperimentalMapSizes)).
				SelectString(func(ms constants.MapSize) string { return ms.Label }).
				ToSlice()
		}()),
		state: state,
	}
	panel.scroll.Axis = layout.Vertical
	panel.LoadFromState()
	return panel
}

func (this *GeneralPanel) GetPanelWidget(theme *material.Theme) layout.Widget {
	widgetsList := []layout.Widget{
		widgets.NewHorizontallySplitWidget(theme, this.getTemplateSectionWidget, this.getMapSectionWidget),
		this.getRulesWidget(theme),
	}

	return func(gtx layout.Context) layout.Dimensions {
		return material.List(theme, &this.scroll).Layout(
			gtx, len(widgetsList),
			func(gtx layout.Context, index int) layout.Dimensions { return widgetsList[index](gtx) })
	}
}

func (this *GeneralPanel) LoadFromState() {
	settings := this.state.GetStateData()

	this.templateName.SetText(settings.TemplateName)
	this.playerCount.Value = utils.Normalize(float32(settings.PlayerCount), 2, 8)
	this.checkMoreMapSizes.Value = settings.ExperimentalMapSizes
	this.updateMapSizeSelectorItems()
	this.mapSizeSelector.SelectByName(constants.GetMapSize(settings.MapSize).Label)

	gameModeIndex := max(slices.Index(constants.GetGameModes(), settings.GameMode), 0)
	this.gameMode.SetSelectedIndex(gameModeIndex)
	this.heroMinimumCount.Value = utils.Normalize(float32(settings.HeroCountMin), 1, 12)
	this.heroMaximumCount.Value = utils.Normalize(float32(settings.HeroCountMax), 1, 12)
	this.heroIncrementPerCastle.Value = utils.Normalize(float32(settings.HeroCountIncrement), 1, 10)

	this.factionLawXpMultiplier.Value = utils.Normalize(float32(settings.FactionLawXpPercent), 25, 200)
	this.astrologyXpMultiplier.Value = utils.Normalize(float32(settings.AstrologyXpPercent), 25, 200)

	victory, victoryKnown := constants.GetVictoryCondition(settings.VictoryCondition)
	if !victoryKnown {
		victory = constants.GetVictoryConditionValues().Standard
		this.state.SetStatus(
			fmt.Sprintf("Unknown victory condition %q in file - reset to Standard.", settings.VictoryCondition),
			false)
	}
	this.victorySelector.SelectByName(victory.Label)

	this.checkLostStartCity.Value = settings.LostStartCity
	this.lostCityDayCount.Value = utils.Normalize(float32(settings.LostStartCityDay), 1, 30)

	this.checkLostStartHero.Value = settings.LostStartHero

	this.checkCityHold.Value = settings.CityHold
	this.cityHoldDaysCount.Value = utils.Normalize(float32(settings.CityHoldDays), 1, 30)

	this.checkGladiatorArena.Value = settings.GladiatorArena
	this.gladiatorDelayCount.Value = utils.Normalize(float32(settings.GladiatorArenaDaysDelayStart), 1, 60)
	this.gladiatorDayCount.Value = utils.Normalize(float32(settings.GladiatorArenaCountDay), 1, 30)

	this.checkTournament.Value = settings.Tournament
	this.tournamentDayCount.Value = utils.Normalize(float32(settings.TournamentFirstTournamentDay), 3, 30)
	this.tournamentIntervalCount.Value = utils.Normalize(float32(settings.TournamentInterval), 3, 30)
	this.tournamentPointsCount.Value = utils.Normalize(float32(settings.TournamentPointsToWin), 1, 10)
	this.checkTournamentSaveArmy.Value = settings.TournamentSaveArmy
}

func (this *GeneralPanel) SaveToState() {
	this.state.UpdateState(func(settings *dtos.EditorStateDto) {
		settings.TemplateName = strings.TrimSpace(this.templateName.Text())
		settings.PlayerCount = int(utils.RoundHalfAway(float64(utils.Denormalize(this.playerCount.Value, 2, 8))))
		settings.MapSize = this.getCurrentMapSize().Size
		settings.ExperimentalMapSizes = this.checkMoreMapSizes.Value

		settings.GameMode = constants.GetGameModes()[this.gameMode.GetSelectedIndex()]
		settings.HeroCountMin = utils.RoundedRange(this.heroMinimumCount.Value, 1, 12)
		settings.HeroCountMax = max(utils.RoundedRange(this.heroMaximumCount.Value, 1, 12), settings.HeroCountMin)
		settings.HeroCountIncrement = utils.RoundedRange(this.heroIncrementPerCastle.Value, 1, 10)

		settings.FactionLawXpPercent = utils.RoundedRange(this.factionLawXpMultiplier.Value, 25, 200)
		settings.AstrologyXpPercent = utils.RoundedRange(this.astrologyXpMultiplier.Value, 25, 200)

		settings.VictoryCondition = this.getCurrentVictoryCondition().ID

		settings.LostStartCity = this.checkLostStartCity.Value
		settings.LostStartCityDay = utils.RoundedRange(this.lostCityDayCount.Value, 1, 30)

		settings.LostStartHero = this.checkLostStartHero.Value

		settings.CityHold = this.checkCityHold.Value || this.isHoldCity()
		settings.CityHoldDays = utils.RoundedRange(this.cityHoldDaysCount.Value, 1, 30)

		settings.GladiatorArena = this.checkGladiatorArena.Value
		settings.GladiatorArenaDaysDelayStart = utils.RoundedRange(this.gladiatorDelayCount.Value, 1, 60)
		settings.GladiatorArenaCountDay = utils.RoundedRange(this.gladiatorDayCount.Value, 1, 30)

		settings.Tournament = this.checkTournament.Value || this.isTournament()
		settings.TournamentFirstTournamentDay = utils.RoundedRange(this.tournamentDayCount.Value, 3, 30)
		settings.TournamentInterval = utils.RoundedRange(this.tournamentIntervalCount.Value, 3, 30)
		settings.TournamentPointsToWin = utils.RoundedRange(this.tournamentPointsCount.Value, 1, 10)
		settings.TournamentSaveArmy = this.checkTournamentSaveArmy.Value
	})
}

func (this *GeneralPanel) getTemplateSectionWidget(theme *material.Theme) layout.Widget {
	return widgets.NewSectionWidget(theme, "Template", []layout.Widget{
		widgets.NewLabeledRowWidget(theme, "Template name", constants.DefaultLabelWidth,
			widgets.NewTextboxWidget(theme, &this.templateName, "Enter template name", false)),
		widgets.NewSliderRowWidget(theme, "Players", constants.DefaultLabelWidth,
			&this.playerCount, utils.RoundedRangeFormatter(2, 8)),
		widgets.NewLabeledRowWidget(theme, "Map size", constants.DefaultLabelWidth,
			func(gtx layout.Context) layout.Dimensions {
				return this.updateMapSizeSelector(gtx).GetWidget(theme)(gtx)
			}),
		widgets.NewLabeledCheckboxRowWidget(theme, &this.checkMoreMapSizes,
			"Allow non official larger map sizes (>240)"),
	})
}

func (this *GeneralPanel) getMapSectionWidget(theme *material.Theme) layout.Widget {
	widgetList := []layout.Widget{
		widgets.NewLabeledRowWidget(theme, "Game mode", constants.DefaultLabelWidth, this.gameMode.GetWidget(theme)),
	}

	if !this.isSingleHero() {
		widgetList = append(widgetList,
			widgets.NewSliderRowWidget(theme, "Hero count min", constants.DefaultLabelWidth,
				&this.heroMinimumCount, utils.RoundedRangeFormatter(1, 12)),
			widgets.NewSliderRowWidget(theme, "Hero count max", constants.DefaultLabelWidth,
				&this.heroMaximumCount, utils.RoundedRangeFormatter(1, 12)),
			widgets.NewSliderRowWidget(theme, "Increment", constants.DefaultLabelWidth,
				&this.heroIncrementPerCastle, utils.RoundedRangeFormatter(1, 10)))
	}

	return widgets.NewSectionWidget(theme, "Hero Restrictions", widgetList)
}

func (this *GeneralPanel) getRulesWidget(theme *material.Theme) layout.Widget {
	return widgets.NewSectionWidget(theme, "Rules", []layout.Widget{
		widgets.NewSliderRowWidget(
			theme, "Faction law experience", constants.DefaultLabelWidth,
			&this.factionLawXpMultiplier, utils.RoundedRangePercentFormatter(25, 200)),
		widgets.NewSliderRowWidget(
			theme, "Astrology experience", constants.DefaultLabelWidth,
			&this.astrologyXpMultiplier, utils.RoundedRangePercentFormatter(25, 200)),
		widgets.NewSectionWidget(theme, "Conditions", []layout.Widget{
			func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Flexed(0.5, widgets.NewLabeledRowWidget(
						theme, "Victory", constants.DefaultLabelWidth,
						this.victorySelector.GetWidget(theme))),
					widgets.NewDefaultWidgetSpacer(),
					layout.Flexed(0.5, this.getConditionOptionsWidget(theme)))
			},
		}),
	})
}

func (this *GeneralPanel) getConditionOptionsWidget(theme *material.Theme) layout.Widget {
	children := []layout.FlexChild{}
	if this.isTournament() {
		children = this.getTournamentOptionsRigidWidgets(theme)
	} else {
		children = this.getNonTournamentOptionsRigidWidgets(theme)
	}

	return func(gtx layout.Context) layout.Dimensions {
		this.updateConditionOptions()
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx, children...)
	}
}

func (this *GeneralPanel) getTournamentOptionsRigidWidgets(theme *material.Theme) []layout.FlexChild {
	return []layout.FlexChild{
		layout.Rigid(widgets.NewLabeledCheckboxRowWidget(theme, &this.checkTournament, "Enable tournament")),
		layout.Rigid(widgets.NewSliderRowWidget(theme, "First tournament day", constants.DefaultLabelWidth,
			&this.tournamentDayCount, utils.RoundedRangeFormatter(3, 30))),
		layout.Rigid(widgets.NewSliderRowWidget(theme, "Interval (days)", constants.DefaultLabelWidth,
			&this.tournamentIntervalCount, utils.RoundedRangeFormatter(3, 30))),
		layout.Rigid(widgets.NewSliderRowWidget(theme, "Points to win", constants.DefaultLabelWidth,
			&this.tournamentPointsCount, utils.RoundedRangeFormatter(1, 10))),
		layout.Rigid(widgets.NewLabeledCheckboxRowWidget(theme,
			&this.checkTournamentSaveArmy, "Save army between rounds")),
	}
}

func (this *GeneralPanel) getNonTournamentOptionsRigidWidgets(theme *material.Theme) []layout.FlexChild {
	return []layout.FlexChild{
		layout.Rigid(widgets.NewLabeledCheckboxRowWidget(theme,
			&this.checkLostStartCity, "Lose if start city is captured")),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if !this.checkLostStartCity.Value {
				return layout.Dimensions{}
			}

			return widgets.NewSliderRowWidget(theme, "Grace period (days)", constants.DefaultLabelWidth,
				&this.lostCityDayCount, utils.RoundedRangeFormatter(1, 30))(gtx)
		}),
		layout.Rigid(widgets.NewLabeledCheckboxRowWidget(theme,
			&this.checkLostStartHero, "Lose if start hero is killed")),
		layout.Rigid(widgets.NewLabeledCheckboxRowWidget(theme, &this.checkCityHold, "Win if you hold a target city")),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if !this.checkCityHold.Value && !this.isHoldCity() {
				return layout.Dimensions{}
			}

			return widgets.NewSliderRowWidget(theme, "Days to hold", constants.DefaultLabelWidth,
				&this.cityHoldDaysCount, utils.RoundedRangeFormatter(1, 30))(gtx)
		}),
		layout.Rigid(widgets.NewLabeledCheckboxRowWidget(theme, &this.checkGladiatorArena, "Enable gladiator arena")),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if !this.checkGladiatorArena.Value {
				return layout.Dimensions{}
			}

			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(widgets.NewSliderRowWidget(theme, "Days delay start", constants.DefaultLabelWidth,
					&this.gladiatorDelayCount, utils.RoundedRangeFormatter(1, 60))),
				layout.Rigid(widgets.NewSliderRowWidget(theme, "Count days", constants.DefaultLabelWidth,
					&this.gladiatorDayCount, utils.RoundedRangeFormatter(1, 30))))
		}),
	}
}

func (this *GeneralPanel) updateMapSizeSelector(gtx layout.Context) *components.DropdownSelector {
	if !this.checkMoreMapSizes.Update(gtx) {
		return this.mapSizeSelector
	}

	this.updateMapSizeSelectorItems()
	return this.mapSizeSelector
}

func (this *GeneralPanel) updateMapSizeSelectorItems() {
	labels := []string{}
	for _, mapSize := range constants.GetMapSizes(this.checkMoreMapSizes.Value) {
		labels = append(labels, mapSize.Label)
	}
	this.mapSizeSelector.SetItems(labels)
	this.mapSizeSelector.SelectByName(constants.GetMapSize(this.state.GetStateData().MapSize).Label)
}

func (this *GeneralPanel) updateConditionOptions() {
	if this.victorySelector.WasUpdated {
		this.checkTournament.Value = false
		this.checkCityHold.Value = false
		this.checkGladiatorArena.Value = false
		this.checkLostStartCity.Value = false
		this.checkLostStartHero.Value = false
	}

	victoryConditions := constants.GetVictoryConditionValues()
	switch this.getCurrentVictoryCondition() {
	case victoryConditions.Standard:
		break
	case victoryConditions.LostStartingCity:
		this.checkLostStartCity.Value = true
	case victoryConditions.GuardianArena:
		this.checkGladiatorArena.Value = true
		this.checkLostStartHero.Value = true
	case victoryConditions.HoldCity:
		this.checkCityHold.Value = true
	case victoryConditions.Tournament:
		this.checkTournament.Value = true
	}
}

func (this *GeneralPanel) getCurrentMapSize() constants.MapSize {
	return constants.GetMapSizes(true)[this.mapSizeSelector.GetSelectedIndex()]
}

func (this *GeneralPanel) getCurrentVictoryCondition() constants.Victory {
	return constants.GetVictoryConditionList()[this.victorySelector.GetSelectedIndex()]
}

func (this *GeneralPanel) isSingleHero() bool {
	return this.gameMode.GetSelectedIndex() == 1
}

func (this *GeneralPanel) isHoldCity() bool {
	return this.getCurrentVictoryCondition() == constants.GetVictoryConditionValues().HoldCity
}

func (this *GeneralPanel) isTournament() bool {
	return this.getCurrentVictoryCondition() == constants.GetVictoryConditionValues().Tournament
}
