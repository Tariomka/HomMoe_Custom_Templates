package editor_state

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
)

// GenerationSettings holds the layout and difficulty knobs the generator reads:
// zone sizes, topology, connectivity and the content/guard density scaling.
type GenerationSettings struct {
	PlayerZoneSize              float64            `json:"playerZoneSize"`
	NeutralZoneSize             float64            `json:"neutralZoneSize"`
	HubZoneSize                 float64            `json:"hubZoneSize"`
	GuardRandomization          float64            `json:"guardRandomization"`
	Topology                    config.MapTopology `json:"topology"`
	RandomPortals               bool               `json:"randomPortals"`
	MaxPortalConnections        int                `json:"maxPortalConns"`
	SpawnRemoteFootholds        bool               `json:"spawnFootholds"`
	RemoteFootholdCount         int                `json:"remoteFootholdCount"`
	GenerateRoads               bool               `json:"generateRoads"`
	NoDirectPlayerConn          bool               `json:"isolateplayers"`
	ResourceDensityPercent      int                `json:"resourceDensity"`
	StructureDensityPercent     int                `json:"structureDensity"`
	NeutralStackStrengthPercent int                `json:"neutralStackStrength"`
	BorderGuardStrengthPercent  int                `json:"borderGuardStrength"`
}
