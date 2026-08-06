package zone_interfaces

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
)

type IRoadFactory interface {
	CreateConnectorZoneRoads(
		connectionNames []string,
		generateRoads bool) []entities.Road

	CreateOuterZoneRoads(
		connectionNames []string,
		mainObjectCount int,
		footholdCount int,
		generateRoads bool) []entities.Road
}
