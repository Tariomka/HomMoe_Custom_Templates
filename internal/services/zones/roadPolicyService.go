package zones

import (
	"slices"
	"strconv"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/road_helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones/zone_interfaces"
)

type RoadPolicyService struct {
	roadFactory zone_interfaces.IRoadFactory
}

func NewRoadPolicyService(roadFactory zone_interfaces.IRoadFactory) zone_interfaces.IRoadPolicyService {
	return &RoadPolicyService{roadFactory: roadFactory}
}

// Reconcile stamps the road setting onto every non-portal connection, then
// drops each road whose endpoint is confirmed invalid - a missing main-object
// index, a connection that is gone, not incident or has its roads off, or a
// named content item absent from the groups the zone references - and appends
// the still-missing approach and foothold routes.
//
// The request's content is authoritative: nil or empty means the zones own no
// content, so every road to a named content item goes.
//
// Castle roads are validated here but never created: the zones already carry
// the arena marker as a main object, so recreating every main-object route
// would road-link the arena as if it were a castle.
func (this *RoadPolicyService) Reconcile(request models.RoadReconciliationRequest) {
	this.reconcile(request, true)
}

// RebuildZoneConnectionRoads is the entry point for callers that hold no
// generation settings: roads are assumed on and the absent content set is
// unknown rather than empty, so content roads are left alone instead of being
// deleted on a guess. Castle roads are rebuilt first because this is also the
// path that follows castles added or removed by hand.
//
// Assuming roads are on also means every non-portal connection leaves here with
// its road flag forced to true, whatever it carried on the way in. Callers that
// hold the road setting must build a models.RoadReconciliationRequest and use
// Reconcile instead.
func (this *RoadPolicyService) RebuildZoneConnectionRoads(
	zones []template_model.Zone,
	connections []template_model.Connection) {
	for index := range zones {
		this.RebuildCastleRoads(&zones[index])
	}

	this.reconcile(models.RoadReconciliationRequest{
		Zones:         zones,
		Connections:   connections,
		GenerateRoads: true,
	}, false)
}

// RebuildCastleRoads reconciles the zone's castle<->castle roads with its
// current main objects: routes anchored on a main object that no longer exists
// are dropped, every other route - including hand-authored ones between other
// castles - is preserved with its attributes, and the missing primary routes
// are appended. The road setting does not gate this: a zone's internal roads
// exist even when zone-to-zone roads are off.
//
// The gladiator-arena marker never gains a route: it is a win-condition prop
// that happens to be filed among the main objects, not a settlement.
func (this *RoadPolicyService) RebuildCastleRoads(zone *template_model.Zone) {
	mainObjectCount := len(zone.MainObjects)
	anchors := roadAnchorIndices(zone.MainObjects)

	// The caller's array can be shared with a preview or revert baseline, and
	// DeleteFunc compacts and zeroes in place, so work on a copy. Clone keeps a
	// nil list nil and an empty list empty.
	roads := slices.DeleteFunc(slices.Clone(zone.Roads), func(road template_model.Road) bool {
		return road_helpers.IsRoadTypeCastle(road) && !isValidCastleRoad(road, mainObjectCount)
	})

	candidates := rebaseMainObjectAnchors(
		this.roadFactory.CreateOuterZoneRoads(nil, len(anchors), 0, false), anchors)
	for _, candidate := range candidates {
		if hasUndirectedRoad(roads, candidate) {
			continue
		}
		roads = append(roads, candidate)
	}
	zone.Roads = roads
}

func (this *RoadPolicyService) reconcile(
	request models.RoadReconciliationRequest,
	validateContent bool) {
	this.stampConnectionRoads(request.Connections, request.GenerateRoads)

	scope := newRoadPolicyScope(request, validateContent)
	for index := range request.Zones {
		zone := &request.Zones[index]
		zone.Roads = this.keepValidRoads(*zone, scope)
		this.appendMissingRoutes(zone, scope)
	}
}

func (this *RoadPolicyService) stampConnectionRoads(
	connections []template_model.Connection,
	generateRoads bool) {
	for index := range connections {
		if connections[index].IsExplicitPortal() {
			continue
		}

		road := generateRoads
		connections[index].Road = &road
	}
}

func (this *RoadPolicyService) keepValidRoads(
	zone template_model.Zone,
	scope roadPolicyScope) []template_model.Road {
	if len(zone.Roads) == 0 {
		return zone.Roads
	}

	var kept []template_model.Road
	for _, road := range zone.Roads {
		if isValidRef(zone, road.From, scope) && isValidRef(zone, road.To, scope) {
			kept = append(kept, road)
		}
	}
	return kept
}

func (this *RoadPolicyService) appendMissingRoutes(zone *template_model.Zone, scope roadPolicyScope) {
	names := scope.eligibleOrder[zone.Name]
	anchors := roadAnchorIndices(zone.MainObjects)
	if len(anchors) == 0 {
		this.appendMissingRoads(zone, this.roadFactory.CreateConnectorZoneRoads(names, true), scope)
		return
	}

	this.appendMissingRoads(
		zone,
		rebaseMainObjectAnchors(
			this.roadFactory.CreateOuterZoneRoads(names, 1, scope.footholdCount, true), anchors[:1]),
		scope)
}

func (this *RoadPolicyService) appendMissingRoads(
	zone *template_model.Zone,
	candidates []template_model.Road,
	scope roadPolicyScope) {
	contentType := registry.GetRoadConnectionTypeValues().MandatoryContent
	for _, candidate := range candidates {
		if candidate.To.Type == contentType && !scope.hasContentItem(*zone, refName(candidate.To)) {
			continue
		}

		if hasUndirectedRoad(zone.Roads, candidate) {
			continue
		}

		zone.Roads = append(zone.Roads, candidate)
	}
}

// roadAnchorIndices returns the main-object indices a road may be anchored on.
// The gladiator arena is filed among the main objects but is a win-condition
// marker, not a settlement, so a zone holding nothing else anchors nothing.
func roadAnchorIndices(mainObjects []template_model.MainObject) []int {
	arenaType := registry.GetMainObjectTypeValues().GladiatorArena
	var anchors []int
	for index, mainObject := range mainObjects {
		if mainObject.Type != arenaType {
			anchors = append(anchors, index)
		}
	}
	return anchors
}

// rebaseMainObjectAnchors maps the factory's dense 0..n-1 anchor numbering onto
// the zone's real main-object indices, so a route never lands on the slot an
// arena marker occupies - which an imported zone may well hold at index 0.
func rebaseMainObjectAnchors(roads []template_model.Road, anchors []int) []template_model.Road {
	if len(anchors) == 0 || anchors[len(anchors)-1] == len(anchors)-1 {
		return roads
	}

	for index := range roads {
		roads[index].From = rebaseAnchorRef(roads[index].From, anchors)
		roads[index].To = rebaseAnchorRef(roads[index].To, anchors)
	}
	return roads
}

func rebaseAnchorRef(ref template_model.TypedRef, anchors []int) template_model.TypedRef {
	if ref.Type != registry.GetRoadConnectionTypeValues().MainObject {
		return ref
	}

	position, err := strconv.Atoi(refName(ref))
	if err != nil || position < 0 || position >= len(anchors) {
		return ref
	}

	args := slices.Clone(ref.Args)
	args[0] = strconv.Itoa(anchors[position])
	ref.Args = args
	return ref
}

func isRoadEligible(connection template_model.Connection, generateRoads bool) bool {
	return connection.IsExplicitPortal() || generateRoads
}

func isValidRef(zone template_model.Zone, ref template_model.TypedRef, scope roadPolicyScope) bool {
	referenceTypes := registry.GetRoadConnectionTypeValues()
	name := refName(ref)
	switch {
	case ref.Type == referenceTypes.MainObject:
		return isValidMainObjectRef(ref, len(zone.MainObjects))
	case ref.Type == referenceTypes.Connection && name != "":
		return scope.allowsConnection(zone.Name, name)
	case ref.Type == referenceTypes.MandatoryContent && name != "" && scope.validateContent:
		return scope.hasContentItem(zone, name)
	default:
		return true
	}
}

func isValidCastleRoad(road template_model.Road, mainObjectCount int) bool {
	return isValidMainObjectRef(road.From, mainObjectCount) &&
		isValidMainObjectRef(road.To, mainObjectCount)
}

func isValidMainObjectRef(ref template_model.TypedRef, mainObjectCount int) bool {
	index, err := strconv.Atoi(refName(ref))
	if err != nil {
		return true
	}

	return index >= 0 && index < mainObjectCount
}

func hasUndirectedRoad(roads []template_model.Road, candidate template_model.Road) bool {
	for _, road := range roads {
		if refEquals(road.From, candidate.From) && refEquals(road.To, candidate.To) {
			return true
		}
		if refEquals(road.From, candidate.To) && refEquals(road.To, candidate.From) {
			return true
		}
	}
	return false
}

func refEquals(left template_model.TypedRef, right template_model.TypedRef) bool {
	return left.Type == right.Type && slices.Equal(left.Args, right.Args)
}

func refName(ref template_model.TypedRef) string {
	if len(ref.Args) == 0 {
		return ""
	}

	return ref.Args[0]
}
