package handler_interfaces

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
)

type IContentRuleHandler interface {
	GetContentRuleEditorOptions(content models.SidMapping) dtos.ContentRuleEditorOptionsDto
	DescribeContentRule(content models.SidMapping, savedRule models.ContentRuleRow) dtos.ContentRuleDescriptionDto
}
