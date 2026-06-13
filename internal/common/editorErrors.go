package common

import "errors"

var (
	ErrNoTemplateName           = errors.New("template name is required")
	ErrGeneratedTemplateInvalid = errors.New("generated template was not valid")
	ErrNothingToSave            = errors.New("nothing to save - generate a template first")
	ErrNoOutputPath             = errors.New("output path is required")
)
