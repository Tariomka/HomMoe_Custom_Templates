package zones

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
)

type roadPolicyScope struct {
	eligibleNames   map[string]map[string]bool
	eligibleOrder   map[string][]string
	contentGroups   map[string]map[string]bool
	validateContent bool
	footholdCount   int
}

func newRoadPolicyScope(
	request models.RoadReconciliationRequest,
	validateContent bool) roadPolicyScope {
	scope := roadPolicyScope{
		eligibleNames:   make(map[string]map[string]bool),
		eligibleOrder:   make(map[string][]string),
		contentGroups:   make(map[string]map[string]bool, len(request.MandatoryContent)),
		validateContent: validateContent,
	}
	if request.SpawnRemoteFootholds && request.RemoteFootholdCount > 0 {
		scope.footholdCount = request.RemoteFootholdCount
	}

	for _, connection := range request.Connections {
		if connection.Name == "" || !isRoadEligible(connection, request.GenerateRoads) {
			continue
		}

		scope.addEligible(connection.From, connection.Name)
		if connection.To != connection.From {
			scope.addEligible(connection.To, connection.Name)
		}
	}

	for _, group := range request.MandatoryContent {
		items := make(map[string]bool, len(group.Content))
		for _, item := range group.Content {
			if item.Name != "" {
				items[item.Name] = true
			}
		}
		scope.contentGroups[group.Name] = items
	}
	return scope
}

func (this roadPolicyScope) addEligible(zoneName string, connectionName string) {
	if zoneName == "" {
		return
	}

	names := this.eligibleNames[zoneName]
	if names == nil {
		names = make(map[string]bool)
		this.eligibleNames[zoneName] = names
	}
	if names[connectionName] {
		return
	}
	names[connectionName] = true
	this.eligibleOrder[zoneName] = append(this.eligibleOrder[zoneName], connectionName)
}

// allowsConnection reports whether the zone may hold a road targeting the named
// connection: it must exist, be incident to that zone and have roads enabled.
func (this roadPolicyScope) allowsConnection(zoneName string, connectionName string) bool {
	return this.eligibleNames[zoneName][connectionName]
}

func (this roadPolicyScope) hasContentItem(zone template_model.Zone, itemName string) bool {
	for _, groupName := range zone.MandatoryContent {
		if this.contentGroups[groupName][itemName] {
			return true
		}
	}
	return false
}
