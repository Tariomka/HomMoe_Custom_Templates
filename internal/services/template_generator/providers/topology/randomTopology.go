package topology

import (
	"math/rand/v2"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/geometry"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	zone_services "github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
)

type RandomTopologyService struct {
	PositionedTopologyBuilder
}

func NewRandomTopologyService() *RandomTopologyService {
	return NewRandomTopologyServiceWithCreationServices(zone_services.NewCreationServices(nil, nil))
}

func NewRandomTopologyServiceWithCreationServices(
	creationServices *zone_services.CreationServices,
) *RandomTopologyService {
	return &RandomTopologyService{
		PositionedTopologyBuilder: *NewPositionedTopologyBuilder(creationServices),
	}
}

func (this *RandomTopologyService) CreateTopologyVariant(
	configuration config.GeneratorConfig,
	playerLabels []string,
	neutralZones neutral_zone.Plans,
	tuning models.GenerationTuning,
	holdCityNeutralLabel string) entities.Variant {
	return this.BuildVariant(
		configuration, playerLabels, neutralZones, tuning, holdCityNeutralLabel, this.createRandomLayout, nil)
}

// createRandomLayout scatters all zones uniformly over the map and connects
// them through a Delaunay triangulation of the random positions.
func (this *RandomTopologyService) createRandomLayout(
	playerLabels []string,
	neutralZones neutral_zone.Plans) ([]string, models.Positions, []models.ConnectionIndexes) {
	neutralLabels := make([]string, len(neutralZones))
	for i, nz := range neutralZones {
		neutralLabels[i] = nz.Label
	}
	allLabels := append(append([]string{}, playerLabels...), neutralLabels...)
	labelCount := len(allLabels)
	rand.Shuffle(labelCount, func(i, j int) { allLabels[i], allLabels[j] = allLabels[j], allLabels[i] })
	var positions models.Positions
	for range labelCount {
		positions.Add(data.NewVec2(rand.Float64()*0.9+0.05, rand.Float64()*0.9+0.05))
	}
	pairs := geometry.CreateDelaunayTriangulation(positions)
	return allLabels, positions, pairs
}
