package models

import "github.com/Tariomka/hommoe_custom_templates/internal/models/template"

// RmgTemplate is the on-disk .rmg.json template. It is a type alias to the
// authoritative model defined in the `template` package, so callers (CLI, GUI,
// generator) can keep using `models.RmgTemplate` while we round-trip through
// the single, correct schema.
type RmgTemplate = template.RmgTemplateModel

// RmgZone aliases the inner zone type for use in GUI code.
type RmgZone = template.Zone
