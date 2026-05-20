package components

import (
	"fmt"

	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/Tariomka/hommoe_custom_templates/internal/constants"
	"github.com/Tariomka/hommoe_custom_templates/internal/gui/components/content"
	"github.com/Tariomka/hommoe_custom_templates/internal/gui/components/widgets"
	"github.com/Tariomka/hommoe_custom_templates/internal/gui/utils"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
)

type RulesPanel struct {
	victory               *content.DropdownSelector
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

	scroll widget.List

	state *State
}

func NewRulesPanel(state *State) *RulesPanel {
	panel := &RulesPanel{
		victory: content.NewDropdownSelector(func() []string {
			labels := make([]string, 0)
			for _, victory := range constants.VictoryConditions {
				labels = append(labels, victory.Label)
			}
			return labels
		}()),
		state: state,
	}
	panel.scroll.Axis = layout.Vertical
	panel.LoadFromState()
	return panel
}

func (this *RulesPanel) GetPanelWidget(theme *material.Theme) layout.Widget {
	widgetsList := []layout.Widget{
		widgets.NewSectionWidget(theme, "Victory Condition", []layout.Widget{
			widgets.NewLabeledRowWidget(theme, "Victory", 160, func(gtx layout.Context) layout.Dimensions {
				return this.victory.Layout(gtx, theme)
			}),
		}),
		widgets.NewSectionWidget(theme, "Loss Conditions", []layout.Widget{
			widgets.NewLabeledCheckboxRowWidget(theme, &this.chkLostStartCity, "Lose if start city is captured"),
			func(gtx layout.Context) layout.Dimensions {
				if !this.chkLostStartCity.Value {
					return layout.Dimensions{}
				}
				return widgets.NewLabeledRowWidget(theme, "Grace period (days)", 200, widgets.NewLabeledSlider(theme, &this.sldLostCityDay, fmt.Sprintf("%d", utils.RoundedRange(this.sldLostCityDay.Value, 1, 30))))(gtx)
			},
			widgets.NewLabeledCheckboxRowWidget(theme, &this.chkLostStartHero, "Lose if start hero is killed"),
		}),
		widgets.NewSectionWidget(theme, "City Hold", []layout.Widget{
			widgets.NewLabeledCheckboxRowWidget(theme, &this.chkCityHold, "Win if you hold a target city"),
			func(gtx layout.Context) layout.Dimensions {
				if !this.chkCityHold.Value && this.victory.GetSelectedIndex() != 2 {
					return layout.Dimensions{}
				}
				return widgets.NewLabeledRowWidget(theme, "Days to hold", 200, widgets.NewLabeledSlider(theme, &this.sldCityHoldDays, fmt.Sprintf("%d", utils.RoundedRange(this.sldCityHoldDays.Value, 1, 30))))(gtx)
			},
		}),
		widgets.NewSectionWidget(theme, "Gladiator Arena", []layout.Widget{
			widgets.NewLabeledCheckboxRowWidget(theme, &this.chkGladiatorArena, "Enable gladiator arena"),
			func(gtx layout.Context) layout.Dimensions {
				if !this.chkGladiatorArena.Value {
					return layout.Dimensions{}
				}
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(widgets.NewLabeledRowWidget(theme, "Days delay start", 200, widgets.NewLabeledSlider(theme, &this.sldGladiatorDelay, fmt.Sprintf("%d", utils.RoundedRange(this.sldGladiatorDelay.Value, 1, 90))))),
					layout.Rigid(widgets.NewLabeledRowWidget(theme, "Count days", 200, widgets.NewLabeledSlider(theme, &this.sldGladiatorCountDay, fmt.Sprintf("%d", utils.RoundedRange(this.sldGladiatorCountDay.Value, 1, 14))))),
				)
			},
		}),
		widgets.NewSectionWidget(theme, "Tournament", []layout.Widget{
			widgets.NewLabeledCheckboxRowWidget(theme, &this.chkTournament, "Enable tournament"),
			func(gtx layout.Context) layout.Dimensions {
				if !this.chkTournament.Value && this.victory.GetSelectedIndex() != 3 {
					return layout.Dimensions{}
				}
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(widgets.NewLabeledRowWidget(theme, "First tournament day", 200, widgets.NewLabeledSlider(theme, &this.sldTournamentDay, fmt.Sprintf("%d", utils.RoundedRange(this.sldTournamentDay.Value, 1, 60))))),
					layout.Rigid(widgets.NewLabeledRowWidget(theme, "Interval (days)", 200, widgets.NewLabeledSlider(theme, &this.sldTournamentInterval, fmt.Sprintf("%d", utils.RoundedRange(this.sldTournamentInterval.Value, 1, 30))))),
					layout.Rigid(widgets.NewLabeledRowWidget(theme, "Points to win", 200, widgets.NewLabeledSlider(theme, &this.sldTournamentPoints, fmt.Sprintf("%d", utils.RoundedRange(this.sldTournamentPoints.Value, 1, 10))))),
					layout.Rigid(widgets.NewLabeledCheckboxRowWidget(theme, &this.chkTournamentSaveArmy, "Save army between rounds")),
				)
			},
		}),
		widgets.NewSectionWidget(theme, "Heroes", []layout.Widget{
			widgets.NewLabeledRowWidget(theme, "Hero count min", 200, widgets.NewLabeledSlider(theme, &this.sldHeroMin, fmt.Sprintf("%d", utils.RoundedRange(this.sldHeroMin.Value, 1, 16)))),
			widgets.NewLabeledRowWidget(theme, "Hero count max", 200, widgets.NewLabeledSlider(theme, &this.sldHeroMax, fmt.Sprintf("%d", utils.RoundedRange(this.sldHeroMax.Value, 1, 16)))),
			widgets.NewLabeledRowWidget(theme, "Increment", 200, widgets.NewLabeledSlider(theme, &this.sldHeroIncr, fmt.Sprintf("%d", utils.RoundedRange(this.sldHeroIncr.Value, 1, 5)))),
		}),
		widgets.NewSectionWidget(theme, "Experience modifiers", []layout.Widget{
			widgets.NewLabeledRowWidget(theme, "Faction laws exp %", 200, widgets.NewLabeledSlider(theme, &this.sldFactionLawsExp, fmt.Sprintf("%d%%", utils.RoundedRange(this.sldFactionLawsExp.Value, 25, 200)))),
			widgets.NewLabeledRowWidget(theme, "Astrology exp %", 200, widgets.NewLabeledSlider(theme, &this.sldAstrologyExp, fmt.Sprintf("%d%%", utils.RoundedRange(this.sldAstrologyExp.Value, 25, 200)))),
		}),
	}
	return func(gtx layout.Context) layout.Dimensions {
		return material.List(theme, &this.scroll).Layout(gtx, len(widgetsList), func(gtx layout.Context, index int) layout.Dimensions {
			return widgetsList[index](gtx)
		})
	}
}

func (this *RulesPanel) LoadFromState() {
	settings := this.state.GetSettingsFile()
	this.victory.SelectByName(constants.GetVictoryCondition(settings.VictoryCondition).Label)
	this.chkLostStartCity.Value = settings.LostStartCity
	this.sldLostCityDay.Value = utils.Normalize(float32(settings.LostStartCityDay), 1, 30)
	this.chkLostStartHero.Value = settings.LostStartHero
	this.chkCityHold.Value = settings.CityHold
	this.sldCityHoldDays.Value = utils.Normalize(float32(settings.CityHoldDays), 1, 30)
	this.chkGladiatorArena.Value = settings.GladiatorArena
	this.sldGladiatorDelay.Value = utils.Normalize(float32(settings.GladiatorArenaDaysDelayStart), 1, 90)
	this.sldGladiatorCountDay.Value = utils.Normalize(float32(settings.GladiatorArenaCountDay), 1, 14)
	this.chkTournament.Value = settings.Tournament
	this.sldTournamentDay.Value = utils.Normalize(float32(settings.TournamentFirstTournamentDay), 1, 60)
	this.sldTournamentInterval.Value = utils.Normalize(float32(settings.TournamentInterval), 1, 30)
	this.sldTournamentPoints.Value = utils.Normalize(float32(settings.TournamentPointsToWin), 1, 10)
	this.chkTournamentSaveArmy.Value = settings.TournamentSaveArmy
	this.sldHeroMin.Value = utils.Normalize(float32(settings.HeroCountMin), 1, 16)
	this.sldHeroMax.Value = utils.Normalize(float32(settings.HeroCountMax), 1, 16)
	this.sldHeroIncr.Value = utils.Normalize(float32(settings.HeroCountIncrement), 1, 5)
	this.sldFactionLawsExp.Value = utils.Normalize(float32(settings.FactionLawsExpPercent), 25, 200)
	this.sldAstrologyExp.Value = utils.Normalize(float32(settings.AstrologyExpPercent), 25, 200)
}

func (this *RulesPanel) SaveToState() {
	// TODO: check `.Update(gtx)` and on true update the value
	this.state.UpdateState(func(settings *models.SettingsFile) {
		settings.VictoryCondition = this.getCurrentVictoryCondition().ID
		settings.LostStartCity = this.chkLostStartCity.Value
		settings.LostStartCityDay = utils.RoundedRange(this.sldLostCityDay.Value, 1, 30)
		settings.LostStartHero = this.chkLostStartHero.Value
		settings.CityHold = this.chkCityHold.Value || this.victory.GetSelectedIndex() == 2
		settings.CityHoldDays = utils.RoundedRange(this.sldCityHoldDays.Value, 1, 30)
		settings.GladiatorArena = this.chkGladiatorArena.Value
		settings.GladiatorArenaDaysDelayStart = utils.RoundedRange(this.sldGladiatorDelay.Value, 1, 90)
		settings.GladiatorArenaCountDay = utils.RoundedRange(this.sldGladiatorCountDay.Value, 1, 14)
		settings.Tournament = this.chkTournament.Value || this.victory.GetSelectedIndex() == 3
		settings.TournamentFirstTournamentDay = utils.RoundedRange(this.sldTournamentDay.Value, 1, 60)
		settings.TournamentInterval = utils.RoundedRange(this.sldTournamentInterval.Value, 1, 30)
		settings.TournamentPointsToWin = utils.RoundedRange(this.sldTournamentPoints.Value, 1, 10)
		settings.TournamentSaveArmy = this.chkTournamentSaveArmy.Value
		settings.HeroCountMin = utils.RoundedRange(this.sldHeroMin.Value, 1, 16)
		settings.HeroCountMax = max(utils.RoundedRange(this.sldHeroMax.Value, 1, 16), settings.HeroCountMin)
		settings.HeroCountIncrement = utils.RoundedRange(this.sldHeroIncr.Value, 1, 5)
		settings.FactionLawsExpPercent = utils.RoundedRange(this.sldFactionLawsExp.Value, 25, 200)
		settings.AstrologyExpPercent = utils.RoundedRange(this.sldAstrologyExp.Value, 25, 200)
	})
}

func (this *RulesPanel) getCurrentVictoryCondition() constants.Victory {
	return constants.VictoryConditions[this.victory.GetSelectedIndex()]
}
