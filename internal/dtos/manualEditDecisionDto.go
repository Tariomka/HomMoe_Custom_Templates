package dtos

import "github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"

// ManualEditDecisionDto says whether a freshly generated template should have
// the manual zone/connection edits reapplied, and which castle options moved
// since the generation those edits were made against.
//
// A nil ReapplyWithCastleChanges means the edits must be dropped instead.
type ManualEditDecisionDto struct {
	ReapplyWithCastleChanges *editor_state_model.CastleSettingChanges
}
