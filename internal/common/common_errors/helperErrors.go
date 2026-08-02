package common_errors

import "errors"

var (
	ErrGameInVDFNotFound    = errors.New("could not find game path in Steam VDF file")
	ErrTemplatesDirNotFound = errors.New("game templates directory not found")
)
