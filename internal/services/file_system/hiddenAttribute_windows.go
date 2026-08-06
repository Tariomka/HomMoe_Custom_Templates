//go:build windows

package file_system

import (
	"os"
	"syscall"
)

// hasHiddenAttribute reports whether a Windows file carries the hidden or
// system attribute, so such entries can be filtered out unless the caller opts
// to show hidden entries.
func hasHiddenAttribute(info os.FileInfo) bool {
	data, ok := info.Sys().(*syscall.Win32FileAttributeData)
	if !ok {
		return false
	}

	const hiddenOrSystem = syscall.FILE_ATTRIBUTE_HIDDEN | syscall.FILE_ATTRIBUTE_SYSTEM
	return data.FileAttributes&hiddenOrSystem != 0
}
