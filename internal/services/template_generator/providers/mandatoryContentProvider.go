package providers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template"
)

type MandatoryContentProvider struct{}

func NewMandatoryContentProvider() *MandatoryContentProvider {
	return &MandatoryContentProvider{}
}

func (this *MandatoryContentProvider) CreateContents(
	configuration config.GeneratorConfig,
	playerLabels []string,
	neutralZones []models.NeutralZonePlan) []template.MandatoryContent {
	var groups []template.MandatoryContent
	for _, letter := range playerLabels {
		groups = append(groups, template.MandatoryContent{
			Name:    "mandatory_content_side_" + letter,
			Content: BuildPlayerZoneMandatoryContent(configuration),
		})
	}
	for _, nz := range neutralZones {
		var content []template.MandatoryContentItem
		switch nz.Quality {
		case models.QualityLow:
			content = BuildLowNeutralMandatoryContent(configuration, nz.CastleCount)
		case models.QualityHigh:
			content = BuildHighNeutralMandatoryContent(configuration, nz.CastleCount)
		default:
			content = BuildMediumNeutralMandatoryContent(configuration, nz.CastleCount)
		}
		groups = append(groups, template.MandatoryContent{
			Name:    "mandatory_content_neutral_" + nz.Letter,
			Content: content,
		})
	}
	return groups
}

// func (this *MandatoryContentProvider) RowsToMandatoryContent(rows []models.ZoneContentRowSave) []template.MandatoryContentItem {
// 	if len(rows) == 0 {
// 		return nil
// 	}
// 	var out []template.MandatoryContentItem
// 	for _, raw := range rows {
// 		row := raw.Normalised()
// 		if row.Sid == "" {
// 			continue
// 		}
// 		for i := 0; i < row.Count; i++ {
// 			out = append(out, this.rowToMandatoryItem(row))
// 		}
// 	}
// 	return out
// }

// func (this *MandatoryContentProvider) rowToMandatoryItem(row models.ZoneContentRowSave) template.MandatoryContentItem {
// 	item := template.MandatoryContentItem{
// 		IsGuarded: row.IsGuarded,
// 		IsMine:    row.IsMine,
// 	}
// 	if row.IsGroup {
// 		item.IncludeLists = []string{row.Sid}
// 	} else {
// 		item.SID = row.Sid
// 	}
// 	if row.NearCastle {
// 		item.Rules = append(item.Rules, ruleNearCastle(1))
// 	}
// 	if d, ok := distanceForLabel(row.RoadDistance); ok {
// 		item.Rules = append(item.Rules, ruleRoadDistance(d, 1))
// 	}
// 	return item
// }

func (this *MandatoryContentProvider) createContentWithFoothold(
	configuration config.GeneratorConfig,
	castleCount int,
	rows []template.MandatoryContentItem) []template.MandatoryContentItem {
	var content []template.MandatoryContentItem
	if configuration.SpawnRemoteFootholds {
		content = append(content, presetRemoteFoothold(castleCount))
	}
	content = append(content, rows...)
	return content
}
