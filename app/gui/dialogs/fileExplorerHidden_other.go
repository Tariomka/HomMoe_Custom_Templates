//go:build !windows

package dialogs

import "os"

// hasHiddenAttr has no meaning outside Windows; dotfile filtering already covers
// the Unix hidden-file convention.
func hasHiddenAttr(info os.FileInfo) bool {
	return false
}
