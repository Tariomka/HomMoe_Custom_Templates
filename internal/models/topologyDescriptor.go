package models

import "github.com/Tariomka/hommoe_custom_templates/internal/models/generator"

type TopologyDescriptor struct {
	Type        generator.MapTopology
	Label       string
	Description string
}
