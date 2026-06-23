package utils

import (
	"os/exec"
	"runtime"
	"strings"
)

// RevealInExplorer opens the given directory in a file manager window.
func RevealInExplorer(path string) error {
	switch runtime.GOOS {
	case "windows":
		return exec.Command("explorer.exe", path).Start()
	case "darwin":
		return exec.Command("open", path).Start()
	default:
		return exec.Command("xdg-open", path).Start()
	}
}

// PickFolder shows a native folder picker and returns the chosen path.
func PickFolder(title, initialDir string) (string, error) {
	if runtime.GOOS != "windows" {
		return "", nil
	}
	script := `
Add-Type -AssemblyName System.Windows.Forms | Out-Null
$d = New-Object System.Windows.Forms.FolderBrowserDialog
$d.Description = '` + escapePS(title) + `'
$d.SelectedPath = '` + escapePS(initialDir) + `'
if ($d.ShowDialog() -eq [System.Windows.Forms.DialogResult]::OK) { Write-Output $d.SelectedPath }
`
	return runPowerShell(script)
}

func runPowerShell(script string) (string, error) {
	command := exec.Command("powershell.exe", "-NoProfile", "-NonInteractive", "-Command", script)
	out, err := command.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// escapePS escapes a value for use inside single-quoted PowerShell strings.
func escapePS(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}
