package template_content_model

import "github.com/Tariomka/hommoe_custom_templates/internal/helpers"

type ContentCountLimit struct {
	Name   string
	Limits []ContentLimit
}

func (this ContentCountLimit) Clone() ContentCountLimit {
	clone := this
	clone.Limits = helpers.MapSlice(this.Limits, ContentLimit.Clone)
	return clone
}
