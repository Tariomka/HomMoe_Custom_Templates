package test_helpers

import (
	"strings"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_errors"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers/handler_interfaces"
)

type FailedLookupFileSystemHandler struct {
	handler_interfaces.IFileSystemHandler

	browseStart string
}

func NewFailedLookupFileSystemHandler(browseStart string) handler_interfaces.IFileSystemHandler {
	return &FailedLookupFileSystemHandler{
		IFileSystemHandler: NewFileSystemHandler(),
		browseStart:        browseStart,
	}
}

func (this *FailedLookupFileSystemHandler) FindGameTemplateDirectory() (string, error) {
	return "", common_errors.ErrTemplatesDirNotFound
}

func (this *FailedLookupFileSystemHandler) ResolveStartDirectory(preferred string) string {
	if strings.TrimSpace(preferred) == "" {
		return this.IFileSystemHandler.ResolveStartDirectory(this.browseStart)
	}

	return this.IFileSystemHandler.ResolveStartDirectory(preferred)
}
