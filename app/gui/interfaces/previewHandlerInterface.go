package interfaces

import "github.com/Tariomka/hommoe_custom_templates/internal/dtos"

type IPreviewHandler interface {
	BuildPreviewLayout(request dtos.PreviewLayoutRequestDto) (dtos.PreviewLayoutDto, error)
}
