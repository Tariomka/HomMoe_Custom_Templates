package dtos

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
)

// ManualEditConvertRequestDto carries an all-overridden manual snapshot and
// the settings it was made with, against which it is recorded again.
type ManualEditConvertRequestDto struct {
	Baseline    editor_state_dto.EditorStateDto
	ManualEdits editor_state_model.ManualEditSettings
}
