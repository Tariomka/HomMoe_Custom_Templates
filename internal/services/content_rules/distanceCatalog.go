package content_rules

import (
	"strings"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_distances"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
)

type DistanceCatalog struct {
	variations []models.DistancePreset
}

func NewDistanceCatalog() *DistanceCatalog {
	return &DistanceCatalog{variations: common_distances.GetContentDistancePresets()}
}

func (this *DistanceCatalog) GetDisplayNames() []string {
	names := make([]string, len(this.variations))
	for index, variation := range this.variations {
		names[index] = variation.Name
	}
	return names
}

func (this *DistanceCatalog) GetByName(name string) (models.DistancePreset, bool) {
	trimmedName := strings.TrimSpace(name)
	for _, variation := range this.variations {
		if strings.EqualFold(variation.Name, trimmedName) {
			return variation, true
		}
	}
	return models.DistancePreset{}, false
}

func defaultDistancePreset() models.DistancePreset {
	preset, _ := common_distances.GetContentDistancePreset("Medium")
	return preset
}
