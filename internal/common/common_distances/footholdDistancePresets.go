package common_distances

import "github.com/Tariomka/hommoe_custom_templates/internal/models"

type FootholdDistancePresets struct {
	Crossroads       models.DistancePreset
	NearMainCastle   models.DistancePreset
	NearSecondCastle models.DistancePreset
}

func GetFootholdDistancePresets() FootholdDistancePresets {
	return FootholdDistancePresets{
		Crossroads:       models.DistancePreset{Name: "Foothold Crossroads", Min: 0.2, Max: 0.3},
		NearMainCastle:   models.DistancePreset{Name: "Foothold Main Castle", Min: 0.2, Max: 0.4},
		NearSecondCastle: models.DistancePreset{Name: "Foothold Second Castle", Min: 0.5, Max: 0.5},
	}
}
