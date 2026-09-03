package template_model

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model/template_variant_model"
)

// The zone, connection and main-object converters are re-exported because the
// seams that still speak entities - the update handler, the zone editor, the
// arena provider - live outside this tree and may not name a *_model
// subpackage. Type aliases carry across a package boundary; functions do not.

func ToZoneModels(entities []template.Zone) []Zone {
	return template_variant_model.ToZoneModels(entities)
}

func ToZoneEntities(models []Zone) []template.Zone {
	return template_variant_model.ToZoneEntities(models)
}

func ToConnectionModels(entities []template.Connection) []Connection {
	return template_variant_model.ToConnectionModels(entities)
}

func ToConnectionEntities(models []Connection) []template.Connection {
	return template_variant_model.ToConnectionEntities(models)
}

func ToZoneEntity(model Zone) template.Zone {
	return template_variant_model.ToZoneEntity(model)
}

func ToZoneModel(entity template.Zone) Zone {
	return template_variant_model.ToZoneModel(entity)
}

func ToMainObjectModel(entity template.MainObject) MainObject {
	return template_variant_model.ToMainObjectModel(entity)
}

func ToMainObjectModels(entities []template.MainObject) []MainObject {
	return template_variant_model.ToMainObjectModels(entities)
}

func ToRoadModel(entity template.Road) Road {
	return template_variant_model.ToRoadModel(entity)
}

func ToRoadModels(entities []template.Road) []Road {
	return template_variant_model.ToRoadModels(entities)
}
