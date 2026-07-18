package dialogs

import (
	"math"
	"strconv"
	"strings"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/widget/material"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/themes"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/widgets"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutralZone"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/connection_editor"
)

func (this *ZoneEditorDialog) zonePropertyRows(theme *material.Theme, zone *entities.Zone) []layout.Widget {
	isNeutral := strings.HasPrefix(zone.Name, "Neutral-")
	isSpawn := this.playerZones[zone.Name]
	rows := []layout.Widget{
		func(gtx layout.Context) layout.Dimensions {
			label := material.Body1(theme, zone.Name)
			label.Color = themes.ColorsBase.AccentBright
			label.Font = font.Font{Weight: font.SemiBold}
			return label.Layout(gtx)
		},
	}
	if isSpawn {
		rows = append(rows,
			widgets.NewDimmedLabelWidget(theme, "Player spawn zone - content is managed by the generator."))
	} else if !isNeutral {
		rows = append(rows, widgets.NewDimmedLabelWidget(theme, "Quality presets apply to neutral zones only."))
	}
	rows = append(rows,
		widgets.NewVerticalSpacerWidget(6),
		widgets.NewLabeledRowWidget(theme, "Size", 110,
			widgets.NewTextboxWidget(theme, &this.zoneSizeEdit, "0.1 - 2.0", false)),
		widgets.NewVerticalSpacerWidget(4),
		widgets.NewLabeledRowWidget(theme, "Guard x", 110,
			widgets.NewTextboxWidget(theme, &this.zoneGuardEdit, "guard multiplier", false)),
		widgets.NewLabeledRowWidget(theme, "Weekly +", 110,
			widgets.NewTextboxWidget(theme, &this.zoneWeeklyEdit, "0.15", false)))
	if isNeutral {
		rows = append(rows,
			widgets.NewVerticalSpacerWidget(4),
			widgets.NewLabeledRowWidget(theme, "Quality", 110, this.qualityDropdown.GetWidget(theme)),
			widgets.NewLabeledRowWidget(theme, "Castles", 110, this.castleDropdown.GetWidget(theme)),
			widgets.NewDimmedLabelWidget(theme, "Changing quality or castles regenerates the zone's content."))
	}
	rows = append(rows,
		widgets.NewVerticalSpacerWidget(10),
		widgets.NewButtonWidget(theme, "Delete this zone", &this.sideZoneDelete, isSpawn))
	return rows
}

// syncZoneProps loads the zone property widgets from the selected zone.
// Called once whenever the zone selection changes.
func (this *ZoneEditorDialog) syncZoneProps(zone *entities.Zone) {
	quality := neutralZone.GetQualityFrom(*zone)
	this.qualityDropdown.SelectByName(connection_editor.QualityLabels[quality.GetIndex()])
	castles := min(connection_editor.CountZoneCastles(*zone), 4)
	this.castleDropdown.SelectByName(strconv.Itoa(castles))
	this.zoneSizeEdit.SetText(strconv.FormatFloat(zone.Size, 'f', -1, 64))
	this.zoneGuardEdit.SetText(strconv.FormatFloat(zone.GuardMultiplier, 'f', -1, 64))
	this.zoneWeeklyEdit.SetText(formatIncrement(zone.GuardWeeklyIncrement))
}

// writebackZoneProps copies the zone widget state back into the selected zone
// after the panel has been laid out for this frame.
func (this *ZoneEditorDialog) writebackZoneProps(zone *entities.Zone) {
	if value, err := strconv.ParseFloat(strings.TrimSpace(this.zoneSizeEdit.Text()), 64); err == nil {
		zone.Size = math.Round(math.Min(math.Max(value, 0.1), 2.0)*100) / 100
	}
	if value, err := strconv.ParseFloat(strings.TrimSpace(this.zoneGuardEdit.Text()), 64); err == nil {
		zone.GuardMultiplier = value
	}
	if value, err := strconv.ParseFloat(strings.TrimSpace(this.zoneWeeklyEdit.Text()), 64); err == nil {
		zone.GuardWeeklyIncrement = value
	}
	if strings.HasPrefix(zone.Name, "Neutral-") && (this.qualityDropdown.WasUpdated || this.castleDropdown.WasUpdated) {
		quality := neutralZone.Quality(this.qualityDropdown.GetSelectedIndex())
		castles := this.castleDropdown.GetSelectedIndex()
		connection_editor.ApplyNeutralZoneQuality(zone, quality, castles, this.tuning)
		this.geometryDirty = true // tier color / castle glyph live in previewZones
		this.syncedZoneFor = ""   // re-sync dependent fields next frame
	}
}
