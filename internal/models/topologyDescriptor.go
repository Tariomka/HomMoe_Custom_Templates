package models

import "github.com/Tariomka/hommoe_custom_templates/internal/models/config"

type TopologyDescriptor struct {
	Type         config.MapTopology
	Label        string
	Description  string
	Capabilities TopologyCapabilities
}
