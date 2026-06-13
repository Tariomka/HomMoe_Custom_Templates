package utils

import "github.com/Tariomka/hommoe_custom_templates/internal/models"

func CloneRuleRows(rules []models.ContentRuleRowSave) []models.ContentRuleRowSave {
	if len(rules) == 0 {
		return nil
	}
	return append([]models.ContentRuleRowSave(nil), rules...)
}
