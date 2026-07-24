package content_rules

import "strings"

type DistanceCatalog struct {
	variations []DistanceVariation
}

func NewDistanceCatalog() *DistanceCatalog {
	return &DistanceCatalog{variations: []DistanceVariation{
		{Name: "Next To", Min: 0.05, Max: 0.1},
		{Name: "Near", Min: 0.1, Max: 0.25},
		{Name: "Medium", Min: 0.25, Max: 0.5},
		{Name: "Far", Min: 0.5, Max: 0.75},
		{Name: "Very Far", Min: 0.75, Max: 0.9},
	}}
}

func (this *DistanceCatalog) GetDisplayNames() []string {
	names := make([]string, len(this.variations))
	for index, variation := range this.variations {
		names[index] = variation.Name
	}
	return names
}

func (this *DistanceCatalog) GetByName(name string) (DistanceVariation, bool) {
	trimmedName := strings.TrimSpace(name)
	for _, variation := range this.variations {
		if strings.EqualFold(variation.Name, trimmedName) {
			return variation, true
		}
	}
	return DistanceVariation{}, false
}
