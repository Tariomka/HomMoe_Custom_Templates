package common_distances

import (
	"strings"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
)

func GetContentDistancePresets() []models.DistancePreset {
	return distancePresets(contentNearDistance())
}

func GetPortalPlacementDistancePresets() []models.DistancePreset {
	return distancePresets(portalNearDistance())
}

func GetContentDistancePreset(name string) (models.DistancePreset, bool) {
	return getDistancePreset(GetContentDistancePresets(), name)
}

func GetPortalPlacementDistancePreset(name string) (models.DistancePreset, bool) {
	return getDistancePreset(GetPortalPlacementDistancePresets(), name)
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

func distancePresets(near models.DistancePreset) []models.DistancePreset {
	return []models.DistancePreset{
		{Name: "Next To", Min: 0.05, Max: 0.1},
		near,
		{Name: "Medium", Min: 0.25, Max: 0.5},
		{Name: "Far", Min: 0.5, Max: 0.75},
		{Name: "Very Far", Min: 0.75, Max: 0.9},
	}
}

func contentNearDistance() models.DistancePreset {
	return models.DistancePreset{Name: "Near", Min: 0.1, Max: 0.25}
}

func portalNearDistance() models.DistancePreset {
	return models.DistancePreset{Name: "Near", Min: 0.075, Max: 0.35}
}
