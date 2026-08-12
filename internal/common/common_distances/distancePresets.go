package common_distances

import (
	"strings"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
)

func GetContentDistancePresets() []models.DistancePreset {
	return []models.DistancePreset{
		{Name: "Next To", Min: 0.05, Max: 0.1},
		{Name: "Near", Min: 0.1, Max: 0.25},
		{Name: "Medium", Min: 0.25, Max: 0.5},
		{Name: "Far", Min: 0.5, Max: 0.75},
		{Name: "Very Far", Min: 0.75, Max: 0.9},
	}
}

func GetContentDistancePreset(name string) (models.DistancePreset, bool) {
	return getDistancePreset(GetContentDistancePresets(), name)
}

func getDistancePreset(presets []models.DistancePreset, name string) (models.DistancePreset, bool) {
	trimmedName := strings.TrimSpace(name)
	for _, preset := range presets {
		if strings.EqualFold(preset.Name, trimmedName) {
			return preset, true
		}
	}

	return models.DistancePreset{}, false
}
