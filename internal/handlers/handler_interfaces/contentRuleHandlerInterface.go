package handler_interfaces

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
)

type IContentRuleHandler interface {
	GetContentRuleEditorOptions(content models.SidMapping) dtos.ContentRuleEditorOptionsDto
	DescribeContentRule(
		content models.SidMapping,
		savedRule editor_state_model.ContentRuleRow) dtos.ContentRuleDescriptionDto
}
