package template_content_model

import "maps"

type ContentPool map[string]any

func (this ContentPool) Clone() ContentPool { return maps.Clone(this) }
