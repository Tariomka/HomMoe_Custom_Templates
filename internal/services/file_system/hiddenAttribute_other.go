//go:build !windows

package file_system

import "os"

// hasHiddenAttribute is Windows-only; elsewhere the dot-prefix convention is
// the only notion of a hidden entry.
func hasHiddenAttribute(_ os.FileInfo) bool { return false }
