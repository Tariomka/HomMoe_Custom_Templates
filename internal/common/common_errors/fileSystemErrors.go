package common_errors

import "errors"

var (
	ErrDirectoryNameEmpty   = errors.New("directory name is empty")
	ErrDirectoryNameInvalid = errors.New("directory name is invalid")
)
