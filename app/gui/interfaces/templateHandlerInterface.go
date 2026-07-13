package interfaces

import "github.com/Tariomka/hommoe_custom_templates/internal/dtos"

// ITemplateHandler is the boundary between the GUI drivers and the template
// generation/persistence backend, allowing State to be tested with a mock
// instead of the full generator and preview stack.
type ITemplateHandler interface {
	GenerateTemplate(stateDto dtos.EditorStateDto) (dtos.TemplateLoadDto, error)
	UpdateTemplate(templateDto dtos.TemplateUpdateDto) (dtos.TemplateLoadDto, error)
	SaveTemplate(templateDto dtos.TemplateSaveDto) (string, error)
	LoadState(path string, fixIssues bool) (*dtos.EditorStateDto, []string, error)
	SaveState(stateDto dtos.EditorStateSaveDto) (string, error)
}
