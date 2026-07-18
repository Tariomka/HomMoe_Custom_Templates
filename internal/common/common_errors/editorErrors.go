package common_errors

import "errors"

var (
	ErrNoTemplateName           = errors.New("template name is required")
	ErrGeneratedTemplateInvalid = errors.New("generated template was not valid")
	ErrNothingToSave            = errors.New("nothing to save")
	ErrNoOutputPath             = errors.New("output path is required")
	ErrProvidedTemplateInvalid  = errors.New("provided template is invalid")
	ErrInvalidPlayerCount       = errors.New("player count must be between 2 and 8")
	ErrZonesMissing             = errors.New("some connections point to non existing zone(s)")
	ErrGameInVDFNotFound        = errors.New("could not find game path in Steam VDF file")
	ErrTemplatesDirNotFound     = errors.New("game templates directory not found")
)
