package template_content_model

import "maps"

type ContentList map[string]any

func (this ContentList) Clone() ContentList { return maps.Clone(this) }
