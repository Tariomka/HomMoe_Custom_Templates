package zone_interfaces

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
)

type IRoadFactory interface {
	CreateConnectorZoneRoads(
		connectionNames []string,
		generateRoads bool) []template_model.Road

	CreateOuterZoneRoads(
		connectionNames []string,
		mainObjectCount int,
		footholdCount int,
		generateRoads bool) []template_model.Road
}
