package template_content_model

import "github.com/Tariomka/hommoe_custom_templates/internal/helpers"

type MandatoryContent struct {
	Name    string
	Content []MandatoryContentItem
}

func (this MandatoryContent) Clone() MandatoryContent {
	clone := this
	clone.Content = helpers.MapSlice(this.Content, MandatoryContentItem.Clone)
	return clone
}
