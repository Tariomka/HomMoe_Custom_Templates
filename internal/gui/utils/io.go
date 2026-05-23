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

// PickOpenFile shows a native file open dialog and returns the chosen path
// (empty string if cancelled). filter is a comma-separated "Description|*.ext"
// pair list, e.g. "Settings (*.gen.json)|*.gen.json|All files|*.*".
func PickOpenFile(title, filter, initialDir string) (string, error) {
	if runtime.GOOS != "windows" {
		return "", nil // unsupported on non-Windows builds
	}
	script := `
Add-Type -AssemblyName System.Windows.Forms | Out-Null
$d = New-Object System.Windows.Forms.OpenFileDialog
$d.Title = '` + escapePS(title) + `'
$d.Filter = '` + escapePS(filter) + `'
$d.InitialDirectory = '` + escapePS(initialDir) + `'
if ($d.ShowDialog() -eq [System.Windows.Forms.DialogResult]::OK) { Write-Output $d.FileName }
`
	return runPowerShell(script)
}

// PickSaveFile shows a native file save dialog and returns the chosen path.
func PickSaveFile(title, filter, initialDir, defaultName string) (string, error) {
	if runtime.GOOS != "windows" {
		return "", nil
	}
	script := `
Add-Type -AssemblyName System.Windows.Forms | Out-Null
$d = New-Object System.Windows.Forms.SaveFileDialog
$d.Title = '` + escapePS(title) + `'
$d.Filter = '` + escapePS(filter) + `'
$d.InitialDirectory = '` + escapePS(initialDir) + `'
$d.FileName = '` + escapePS(defaultName) + `'
if ($d.ShowDialog() -eq [System.Windows.Forms.DialogResult]::OK) { Write-Output $d.FileName }
`
	return runPowerShell(script)
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
