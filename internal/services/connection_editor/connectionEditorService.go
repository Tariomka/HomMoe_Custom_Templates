// Package connection_editor contains the model- and logic-layer behaviour of the
// Zone Connection Editor. The visual canvas itself remains in the GUI layer.
package connection_editor

import (
	"strings"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_connections"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/linq"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones/zone_interfaces"
)

type ConnectionEditorService struct {
	zoneClassifier zone_interfaces.IZoneClassifier
}

func NewConnectionEditorService(
	zoneClassifier zone_interfaces.IZoneClassifier) IConnectionEditorService {
	return &ConnectionEditorService{zoneClassifier: zoneClassifier}
}

func (this *ConnectionEditorService) NewDefaultConnection(
	from string,
	to string,
	zones []entities.Zone,
	playerZoneNames map[string]bool) entities.Connection {
	quality := this.zoneClassifier.GetConnectionGuardQuality(
		from, to, zones, linq.FromMap(playerZoneNames).SelectKeys().ToSlice())
	return variant_content.NewConnectionBuilder().
		WithFrom(from).
		WithTo(to).
		WithConnectionTypeDirect().
		WithGuardValue(common_connections.GetGuardStrengthForQuality(quality).Default).
		WithGuardZone(from).
		WithGuardMatchGroup("rnd_guard_" + helpers.GetZoneLabel(from) + "_" + helpers.GetZoneLabel(to)).
		WithGuardWeeklyIncrement(common_connections.GetGuardWeeklyIncrements().Standard).
		WithIsUserAdded().
		Build()
}

func (this *ConnectionEditorService) FindIsolatedZones(
	zones []entities.Zone,
	connections []entities.Connection) []string {
	var isolated []string

	for _, zone := range zones {
		referenced := false
		for _, connection := range connections {
			if connection.From == zone.Name || connection.To == zone.Name {
				referenced = true
				break
			}
		}
		if !referenced {
			isolated = append(isolated, zone.Name)
		}
	}
	return isolated
}

func (this *ConnectionEditorService) ComputeHasErrors(zones []entities.Zone, connections []entities.Connection) bool {
	zoneNames := make(map[string]bool, len(zones))
	for _, zone := range zones {
		zoneNames[zone.Name] = true
	}
	for _, connection := range connections {
		if !zoneNames[connection.From] || !zoneNames[connection.To] {
			return true
		}
	}
	return false
}

func (this *ConnectionEditorService) HasDuplicateName(
	connections []entities.Connection,
	current *entities.Connection) bool {
	if current == nil || len(current.Name) == 0 {
		return false
	}

	for index := range connections {
		if &connections[index] == current {
			continue
		}
		if strings.EqualFold(connections[index].Name, current.Name) {
			return true
		}
	}

	return false
}
