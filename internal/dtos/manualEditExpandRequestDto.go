package dtos

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
)

// ManualEditExpandRequestDto carries the freshly generated template, the
// state holding the compact manual edits, and the castle options that moved
// since the edits were recorded.
type ManualEditExpandRequestDto struct {
	Template      *template_model.Template
	State         editor_state_dto.EditorStateDto
	CastleChanges editor_state_model.CastleSettingChanges
}
