package editor_state_dto

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
)

type GenerationSettingsDto struct {
	PlayerZoneSize              float64
	NeutralZoneSize             float64
	HubZoneSize                 float64
	GuardRandomization          float64
	Topology                    config.MapTopology
	RandomPortals               bool
	MaxPortalConnections        int
	SpawnRemoteFootholds        bool
	RemoteFootholdCount         int
	GenerateRoads               bool
	NoDirectPlayerConn          bool
	ResourceDensityPercent      int
	StructureDensityPercent     int
	NeutralStackStrengthPercent int
	BorderGuardStrengthPercent  int
}
