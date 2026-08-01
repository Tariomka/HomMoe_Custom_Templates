package panels

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/components"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/drivers"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/interfaces"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/utils"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/widgets"
	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_topologies"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
)

type LayoutPanel struct {
	topology *components.DropdownSelector

	chkRoads               widget.Bool
	chkPortals             widget.Bool
	sldMaxPortals          widget.Float
	chkFootholds           widget.Bool
	sldRemoteFootholds     widget.Float
	chkPlayerIsolation     widget.Bool
	chkMatchPlayerFactions widget.Bool
	chkAbandonedOutposts   widget.Bool
	sldAbandonedOutposts   widget.Float

	chkAdvancedZones         widget.Bool
	sldNeutralLowestNoCastle widget.Float
	sldNeutralLowestCastle   widget.Float
	sldNeutralLowNoCastle    widget.Float
	sldNeutralLowCastle      widget.Float
	sldNeutralMedNoCastle    widget.Float
	sldNeutralMedCastle      widget.Float
	sldNeutralHighNoCastle   widget.Float
	sldNeutralHighCastle     widget.Float

	sldNeutralLowestCastlesPerZone widget.Float
	sldNeutralLowCastlesPerZone    widget.Float
	sldNeutralMedCastlesPerZone    widget.Float
	sldNeutralHighCastlesPerZone   widget.Float

	sldNeutralCount       widget.Float
	sldPlayerOwnedCastles widget.Float
	sldPlayerCastles      widget.Float
	sldNeutralCastles     widget.Float
	sldHubSize            widget.Float
	sldHubCastles         widget.Float
	sldPlayerZoneSize     widget.Float
	sldNeutralZoneSize    widget.Float
	sldGuardRandom        widget.Float
	sldResourceDensity    widget.Float
	sldStructureDensity   widget.Float
	sldNeutralStack       widget.Float
	sldBorderGuard        widget.Float

	editConnectionsBtn widget.Clickable
	btnPlayerContent   widget.Clickable
	btnLowestContent   widget.Clickable
	btnLowContent      widget.Clickable
	btnMedContent      widget.Clickable
	btnHighContent     widget.Clickable
	btnHubContent      widget.Clickable

	scroll widget.List

	state              *drivers.State
	previewHandler     interfaces.IPreviewHandler
	contentRuleHandler interfaces.IContentRuleHandler
	zoneEditorHandler  interfaces.IZoneEditorHandler
}

func NewLayoutPanel(
	state *drivers.State,
	previewHandler interfaces.IPreviewHandler,
	contentRuleHandler interfaces.IContentRuleHandler,
	zoneEditorHandler interfaces.IZoneEditorHandler,
) *LayoutPanel {
	panel := &LayoutPanel{
		topology: components.NewDropdownSelector(func() []string {
			labels := make([]string, 0)
			for topology := range common_topologies.GetTopologyDescriptorSeq() {
				labels = append(labels, topology.Label)
			}
			return labels
		}()),
		state:              state,
		previewHandler:     previewHandler,
		contentRuleHandler: contentRuleHandler,
		zoneEditorHandler:  zoneEditorHandler,
	}
	panel.scroll.Axis = layout.Vertical
	panel.LoadFromState()
	return panel
}

func (this *LayoutPanel) GetPanelWidget(theme *material.Theme) layout.Widget {
	widgetsList := []layout.Widget{
		widgets.NewHorizontallySplitWidget(theme,
			func(theme *material.Theme) layout.Widget {
				return func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(this.getTopologySectionWidget(theme)),
						layout.Rigid(this.getConnectivityWidget(theme)),
						layout.Rigid(this.getZoneSizesWidget(theme)),
						layout.Rigid(this.getDifficultyAndDensityWidget(theme)),
					)
				}
			},
			func(theme *material.Theme) layout.Widget {
				return func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(this.getManualZoneEditWidget(theme)),
						layout.Rigid(this.getZonesWidget(theme)),
					)
				}
			}),
	}

	return func(gtx layout.Context) layout.Dimensions {
		this.handleConnectionEditorClick(gtx)
		this.handleZoneContentDialogClicks(gtx)
		return material.List(theme, &this.scroll).Layout(gtx, len(widgetsList),
			func(gtx layout.Context, index int) layout.Dimensions { return widgetsList[index](gtx) })
	}
}

func (this *LayoutPanel) LoadFromState() {
	settings := this.state.GetStateData()

	this.topology.SelectByName(common_topologies.GetTopologyDescriptorFromType(settings.Topology).Label)

	this.chkRoads.Value = settings.GenerateRoads
	this.chkPortals.Value = settings.RandomPortals
	this.sldMaxPortals.Value = utils.Normalize(float32(settings.MaxPortalConnections), 1, 32)
	this.chkFootholds.Value = settings.SpawnRemoteFootholds
	this.sldRemoteFootholds.Value = utils.Normalize(float32(settings.RemoteFootholdCount), 0, 4)
	this.chkPlayerIsolation.Value = settings.NoDirectPlayerConn
	this.chkMatchPlayerFactions.Value = settings.MatchPlayerCastleFactions
	this.chkAbandonedOutposts.Value = settings.SpawnAbandonedOutposts
	this.sldAbandonedOutposts.Value = utils.Normalize(float32(settings.AbandonedOutpostCount), 0, 4)

	this.chkAdvancedZones.Value = settings.AdvancedMode
	this.sldNeutralCount.Value = utils.Normalize(float32(settings.NeutralZoneCount), 0, 16)
	this.sldPlayerOwnedCastles.Value = utils.Normalize(float32(settings.PlayerOwnedCastles), 0, 4)
	this.sldPlayerCastles.Value = utils.Normalize(float32(settings.PlayerZoneCastles), 0, 4)
	this.sldNeutralCastles.Value = utils.Normalize(float32(settings.NeutralZoneCastles), 0, 4)
	this.sldNeutralLowestNoCastle.Value = utils.Normalize(float32(settings.NeutralLowestNoCastleCount), 0, 8)
	this.sldNeutralLowestCastle.Value = utils.Normalize(float32(settings.NeutralLowestCastleCount), 0, 8)
	this.sldNeutralLowNoCastle.Value = utils.Normalize(float32(settings.NeutralLowNoCastleCount), 0, 8)
	this.sldNeutralLowCastle.Value = utils.Normalize(float32(settings.NeutralLowCastleCount), 0, 8)
	this.sldNeutralMedNoCastle.Value = utils.Normalize(float32(settings.NeutralMediumNoCastleCount), 0, 8)
	this.sldNeutralMedCastle.Value = utils.Normalize(float32(settings.NeutralMediumCastleCount), 0, 8)
	this.sldNeutralHighNoCastle.Value = utils.Normalize(float32(settings.NeutralHighNoCastleCount), 0, 8)
	this.sldNeutralHighCastle.Value = utils.Normalize(float32(settings.NeutralHighCastleCount), 0, 8)
	this.sldNeutralLowestCastlesPerZone.Value = utils.Normalize(float32(settings.NeutralLowestCastlesPerZone), 1, 4)
	this.sldNeutralLowCastlesPerZone.Value = utils.Normalize(float32(settings.NeutralLowCastlesPerZone), 1, 4)
	this.sldNeutralMedCastlesPerZone.Value = utils.Normalize(float32(settings.NeutralMediumCastlesPerZone), 1, 4)
	this.sldNeutralHighCastlesPerZone.Value = utils.Normalize(float32(settings.NeutralHighCastlesPerZone), 1, 4)
	this.sldHubSize.Value = float32((settings.HubZoneSize - 0.5) / 1.5)
	this.sldHubCastles.Value = utils.Normalize(float32(settings.HubZoneCastles), 0, 4)
	this.sldPlayerZoneSize.Value = float32((settings.PlayerZoneSize - 0.5) / 1.5)
	this.sldNeutralZoneSize.Value = float32((settings.NeutralZoneSize - 0.5) / 1.5)
	this.sldGuardRandom.Value = utils.Normalize(float32(settings.GuardRandomization), 0, 0.5)
	this.sldResourceDensity.Value = utils.Normalize(float32(settings.ResourceDensityPercent), 25, 200)
	this.sldStructureDensity.Value = utils.Normalize(float32(settings.StructureDensityPercent), 25, 200)
	this.sldNeutralStack.Value = utils.Normalize(float32(settings.NeutralStackStrengthPercent), 25, 200)
	this.sldBorderGuard.Value = utils.Normalize(float32(settings.BorderGuardStrengthPercent), 25, 200)
}

func (this *LayoutPanel) SaveToState() {
	this.state.UpdateState(func(settings *dtos.EditorStateDto) {
		settings.Topology = this.getCurrentTopology().Type

		settings.GenerateRoads = this.chkRoads.Value
		settings.RandomPortals = this.chkPortals.Value
		settings.MaxPortalConnections = utils.RoundedRange(this.sldMaxPortals.Value, 1, 32)
		settings.SpawnRemoteFootholds = this.chkFootholds.Value
		settings.RemoteFootholdCount = utils.RoundedRange(this.sldRemoteFootholds.Value, 0, 4)
		settings.NoDirectPlayerConn = this.chkPlayerIsolation.Value
		settings.MatchPlayerCastleFactions = this.chkMatchPlayerFactions.Value
		settings.SpawnAbandonedOutposts = this.chkAbandonedOutposts.Value
		settings.AbandonedOutpostCount = utils.RoundedRange(this.sldAbandonedOutposts.Value, 0, 4)

		settings.AdvancedMode = this.chkAdvancedZones.Value
		settings.NeutralZoneCount = utils.RoundedRange(this.sldNeutralCount.Value, 0, 16)
		settings.PlayerOwnedCastles = utils.RoundedRange(this.sldPlayerOwnedCastles.Value, 0, 4)
		settings.PlayerZoneCastles = utils.RoundedRange(this.sldPlayerCastles.Value, 0, 4)
		settings.NeutralZoneCastles = utils.RoundedRange(this.sldNeutralCastles.Value, 0, 4)
		settings.NeutralLowestNoCastleCount = utils.RoundedRange(this.sldNeutralLowestNoCastle.Value, 0, 8)
		settings.NeutralLowestCastleCount = utils.RoundedRange(this.sldNeutralLowestCastle.Value, 0, 8)
		settings.NeutralLowNoCastleCount = utils.RoundedRange(this.sldNeutralLowNoCastle.Value, 0, 8)
		settings.NeutralLowCastleCount = utils.RoundedRange(this.sldNeutralLowCastle.Value, 0, 8)
		settings.NeutralMediumNoCastleCount = utils.RoundedRange(this.sldNeutralMedNoCastle.Value, 0, 8)
		settings.NeutralMediumCastleCount = utils.RoundedRange(this.sldNeutralMedCastle.Value, 0, 8)
		settings.NeutralHighNoCastleCount = utils.RoundedRange(this.sldNeutralHighNoCastle.Value, 0, 8)
		settings.NeutralHighCastleCount = utils.RoundedRange(this.sldNeutralHighCastle.Value, 0, 8)
		settings.NeutralLowestCastlesPerZone = utils.RoundedRange(this.sldNeutralLowestCastlesPerZone.Value, 1, 4)
		settings.NeutralLowCastlesPerZone = utils.RoundedRange(this.sldNeutralLowCastlesPerZone.Value, 1, 4)
		settings.NeutralMediumCastlesPerZone = utils.RoundedRange(this.sldNeutralMedCastlesPerZone.Value, 1, 4)
		settings.NeutralHighCastlesPerZone = utils.RoundedRange(this.sldNeutralHighCastlesPerZone.Value, 1, 4)
		settings.HubZoneSize = float64(0.5 + this.sldHubSize.Value*1.5)
		settings.HubZoneCastles = utils.RoundedRange(this.sldHubCastles.Value, 0, 4)
		settings.PlayerZoneSize = float64(0.5 + this.sldPlayerZoneSize.Value*1.5)
		settings.NeutralZoneSize = float64(0.5 + this.sldNeutralZoneSize.Value*1.5)
		settings.GuardRandomization = float64(utils.Denormalize(this.sldGuardRandom.Value, 0, 0.5))
		settings.ResourceDensityPercent = utils.RoundedRange(this.sldResourceDensity.Value, 25, 200)
		settings.StructureDensityPercent = utils.RoundedRange(this.sldStructureDensity.Value, 25, 200)
		settings.NeutralStackStrengthPercent = utils.RoundedRange(this.sldNeutralStack.Value, 25, 200)
		settings.BorderGuardStrengthPercent = utils.RoundedRange(this.sldBorderGuard.Value, 25, 200)
	})
}

func (this *LayoutPanel) getCurrentTopology() models.TopologyDescriptor {
	return common_topologies.GetTopologyDescriptorFromIndex(this.topology.GetSelectedIndex())
}
