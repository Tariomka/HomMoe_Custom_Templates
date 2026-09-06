package editor_state_v1

import "github.com/Tariomka/hommoe_custom_templates/internal/entities/topology"

// GenerationSettings is the frozen v1 snapshot - see editorState.go. Never edit.
type GenerationSettings struct {
	PlayerZoneSize              float64              `json:"playerZoneSize"`
	NeutralZoneSize             float64              `json:"neutralZoneSize"`
	HubZoneSize                 float64              `json:"hubZoneSize"`
	GuardRandomization          float64              `json:"guardRandomization"`
	Topology                    topology.MapTopology `json:"topology"`
	RandomPortals               bool                 `json:"randomPortals"`
	MaxPortalConnections        int                  `json:"maxPortalConns"`
	SpawnRemoteFootholds        bool                 `json:"spawnFootholds"`
	RemoteFootholdCount         int                  `json:"remoteFootholdCount"`
	GenerateRoads               bool                 `json:"generateRoads"`
	NoDirectPlayerConn          bool                 `json:"isolateplayers"`
	ResourceDensityPercent      int                  `json:"resourceDensity"`
	StructureDensityPercent     int                  `json:"structureDensity"`
	NeutralStackStrengthPercent int                  `json:"neutralStackStrength"`
	BorderGuardStrengthPercent  int                  `json:"borderGuardStrength"`
}
