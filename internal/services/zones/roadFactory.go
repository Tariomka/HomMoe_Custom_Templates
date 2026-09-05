package zones

import (
	"fmt"
	"strconv"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones/zone_interfaces"
)

type RoadFactory struct{}

func NewRoadFactory() zone_interfaces.IRoadFactory {
	return &RoadFactory{}
}

func (this *RoadFactory) CreateConnectorZoneRoads(
	connectionNames []string,
	generateRoads bool) []template_model.Road {
	if !generateRoads {
		return nil
	}

	distinctNames := helpers.GetUniqueElements(connectionNames)
	if len(distinctNames) == 0 {
		return nil
	}

	if len(distinctNames) == 1 {
		return []template_model.Road{
			variant_content.NewRoadBuilder().
				WithFrom(variant_content.NewRefBuilder().BuildConnectionType(distinctNames[0])).
				WithTo(variant_content.NewRefBuilder().BuildConnectionType(distinctNames[0])).
				Build()}
	}
	var roads []template_model.Road
	for _, name := range distinctNames[1:] {
		roads = append(roads,
			variant_content.NewRoadBuilder().
				WithFrom(variant_content.NewRefBuilder().BuildConnectionType(distinctNames[0])).
				WithTo(variant_content.NewRefBuilder().BuildConnectionType(name)).
				Build())
	}
	return roads
}

func (this *RoadFactory) CreateOuterZoneRoads(
	connectionNames []string,
	mainObjectCount int,
	footholdCount int,
	generateRoads bool) []template_model.Road {
	if !generateRoads {
		return nil
	}

	if mainObjectCount == 0 {
		return this.CreateConnectorZoneRoads(connectionNames, generateRoads)
	}

	var roads []template_model.Road
	for index := range mainObjectCount - 1 {
		roads = append(roads,
			variant_content.NewRoadBuilder().
				WithStoneType().
				WithFrom(variant_content.NewRefBuilder().BuildMainObjectType("0")).
				WithTo(variant_content.NewRefBuilder().BuildMainObjectType(strconv.Itoa(index+1))).
				Build())
	}
	for index := range footholdCount {
		roads = append(roads,
			variant_content.NewRoadBuilder().
				WithFrom(variant_content.NewRefBuilder().BuildMainObjectType("0")).
				WithTo(variant_content.NewRefBuilder().
					BuildMandatoryContentType(fmt.Sprintf("name_remote_foothold_%d", index+1))).
				Build())
	}
	for _, name := range connectionNames {
		roads = append(roads,
			variant_content.NewRoadBuilder().
				WithFrom(variant_content.NewRefBuilder().BuildMainObjectType("0")).
				WithTo(variant_content.NewRefBuilder().BuildConnectionType(name)).
				Build())
	}
	return roads
}
