package dtos

import "github.com/Tariomka/hommoe_custom_templates/internal/models"

type ContentRuleDescriptionDto struct {
	Key          ContentRuleKey
	DisplayText  string
	Marker       string
	VariantLabel string
	Valid        bool
	SavedRule    models.ContentRuleRowSave
}
