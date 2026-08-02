package common_errors

import "errors"

var (
	ErrNoVariantProvided = errors.New("provided variant mapping has no variants")
)
