// Package gui — platform helpers for opening URLs, picking files, and
// revealing folders. On Windows we shell out to PowerShell which uses the
// native System.Windows.Forms dialogs. On other platforms the helpers fall
// back to xdg-open / open. No cgo is used.
package gui

import (
	"os/exec"
	"runtime"
	"strings"
)

// OpenURL opens the given URL in the user's default web browser.
func OpenURL(url string) error {
	switch runtime.GOOS {
	case "windows":
		// rundll32 is the most reliable launcher for hyperlinks on Windows.
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		return exec.Command("open", url).Start()
	default:
		return exec.Command("xdg-open", url).Start()
	}
}

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
// pair list, e.g. "Settings (*.oetgs)|*.oetgs|All files|*.*".
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
	cmd := exec.Command("powershell.exe", "-NoProfile", "-NonInteractive", "-Command", script)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// escapePS escapes a value for use inside single-quoted PowerShell strings.
func escapePS(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}

// FindOldenEraTemplatesDir tries to locate the official Steam install folder
// for "Heroes of Might and Magic Olden Era" and returns its map_templates
// directory, or "" if it cannot be located. Mirrors GetSteamTemplatesDir() in
// MainWindow.xaml.cs.
func FindOldenEraTemplatesDir() string {
	if runtime.GOOS != "windows" {
		return ""
	}
	script := `
$paths = @()
$reg = @(
  'HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\Steam App 3105440',
  'HKLM:\SOFTWARE\WOW6432Node\Microsoft\Windows\CurrentVersion\Uninstall\Steam App 3105440'
)
foreach ($k in $reg) {
  try {
    $v = Get-ItemProperty -Path $k -ErrorAction Stop
    if ($v.InstallLocation) { $paths += (Join-Path $v.InstallLocation 'HeroesOldenEra_Data\StreamingAssets\map_templates') }
  } catch {}
}
$paths += @(
  (Join-Path $env:ProgramFiles      'Steam\steamapps\common\Heroes of Might and Magic Olden Era\HeroesOldenEra_Data\StreamingAssets\map_templates'),
  (Join-Path ${env:ProgramFiles(x86)} 'Steam\steamapps\common\Heroes of Might and Magic Olden Era\HeroesOldenEra_Data\StreamingAssets\map_templates')
)
foreach ($p in $paths) { if ($p -and (Test-Path -Path $p -PathType Container)) { Write-Output $p; break } }
`
	out, err := runPowerShell(script)
	if err != nil {
		return ""
	}
	return out
}
