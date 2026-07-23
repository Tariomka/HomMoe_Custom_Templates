//go:build windows

package helpers

import (
	"path/filepath"

	"golang.org/x/sys/windows/registry"
)

// getSteamPathFromRegistry reads the Steam install location from the current
// user's registry hive (HKCU\Software\Valve\Steam, value "SteamPath").
// It returns an empty string when the key or value is unavailable so callers
// can fall back to the conventional install paths.
func getSteamPathFromRegistry() string {
	key, err := registry.OpenKey(registry.CURRENT_USER, `Software\Valve\Steam`, registry.QUERY_VALUE)
	if err != nil {
		return ""
	}
	defer key.Close()

	steamPath, _, err := key.GetStringValue("SteamPath")
	if err != nil || steamPath == "" {
		return ""
	}

	// The registry stores the path with forward slashes; normalize it.
	return filepath.Clean(steamPath)
}
