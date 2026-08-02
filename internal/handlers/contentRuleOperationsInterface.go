package handlers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
)

type IContentRuleOperations interface {
	GetContentRuleEditorOptions(content models.SidMapping) dtos.ContentRuleEditorOptionsDto
	DescribeContentRule(
		content models.SidMapping,
		savedRule models.ContentRuleRowSave,
	) dtos.ContentRuleDescriptionDto
}
