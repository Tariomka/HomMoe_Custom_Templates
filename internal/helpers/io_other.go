//go:build !windows

package helpers

// getSteamPathFromRegistry is a no-op on non-Windows platforms; the Windows
// implementation reads the Steam install location from the registry.
func getSteamPathFromRegistry() string {
	return ""
}
