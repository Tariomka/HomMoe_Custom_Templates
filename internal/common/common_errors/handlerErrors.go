package common_errors

import "errors"

var (
	ErrNoTemplateName           = errors.New("template name is required")
	ErrGeneratedTemplateInvalid = errors.New("generated template was not valid")
	ErrNothingToSave            = errors.New("nothing to save")
	ErrNoOutputPath             = errors.New("output path is required")
	ErrProvidedTemplateInvalid  = errors.New("provided template is invalid")
	ErrZonesMissing             = errors.New("some connections point to non existing zone(s)")
)
