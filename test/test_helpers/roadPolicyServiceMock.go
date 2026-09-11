package test_helpers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/stretchr/testify/mock"
)

// RoadPolicyServiceMock is a testify mock of
// zone_interfaces.IRoadPolicyService, used to prove what its callers hand it
// without running the real reconciliation.
type RoadPolicyServiceMock struct {
	mock.Mock
}

func (this *RoadPolicyServiceMock) Reconcile(request models.RoadReconciliationRequest) {
	this.Called(request)
}

func (this *RoadPolicyServiceMock) RebuildZoneConnectionRoads(
	zones []template_model.Zone,
	connections []template_model.Connection) {
	this.Called(zones, connections)
}

func (this *RoadPolicyServiceMock) RebuildCastleRoads(zone *template_model.Zone) {
	this.Called(zone)
}

func (this *RoadPolicyServiceMock) StampConnectionRoad(
	connection *template_model.Connection,
	generateRoads bool) {
	this.Called(connection, generateRoads)
}
