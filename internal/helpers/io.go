package helpers

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/andygrunwald/vdf"
)

var (
	libraryFoldersVDFRelativePath = filepath.Join("steamapps", "libraryfolders.vdf")
	windowsSteamPath              = filepath.Join("C:", "Program Files (x86)", "Steam")
	unixSteamPath                 = filepath.Join("~", ".local", "share", "Steam")
	unixSteamAltPath              = filepath.Join("~", ".steam", "steam")

	oldenEraID                   = "3105440"
	oldenEraTemplateRelativePath = filepath.Join("steamapps", "common", "Heroes of Might and Magic Olden Era", "HeroesOldenEra_Data", "StreamingAssets", "map_templates")

	windowsUserPath            = filepath.Join("C:", "Users", os.Getenv("USERNAME"))
	unixCompatUserRelativePath = filepath.Join("compatdata", oldenEraID, "pfx", "drive_c", "users", "steamuser")
	customTemplateRelativeGlob = filepath.Join("AppData", "LocalLow", "Unfrozen", "HeroesOldenEra", "users", "*", "my_map_templates")
)

// FindOldenEraTemplatesDir tries to locate a template folder in common directories.
// Depending on the useInstallDir flag, it either looks for the official Steam install directory
// or tries to find the user directory and tries to resolve the custom template path from there.
// Returns "" if it cannot be located
func FindOldenEraTemplatesDir(useInstallDir bool) string {
	if !useInstallDir && runtime.GOOS == "windows" {
		templatePathPattern := filepath.Join(windowsUserPath, customTemplateRelativeGlob)
		return tryResolveGlob(templatePathPattern)
	}

	content, err := getVDFContent()
	if err != nil {
		return ""
	}

	directory := getBasePath(content)
	if directory == "" {
		return ""
	}

	if !useInstallDir /*&& runtime.GOOS != "windows" is redundant here*/ {
		templatePathPattern := filepath.Join(directory, unixCompatUserRelativePath, customTemplateRelativeGlob)
		return tryResolveGlob(templatePathPattern)
	}

	directory = filepath.Join(directory, oldenEraTemplateRelativePath)
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		return ""
	}

	return directory
}

func getVDFContent() (map[string]any, error) {
	vdfPath, err := getVDFFilePath()
	if err != nil {
		return nil, err
	}

	vdfFile, err := os.Open(vdfPath)
	if err != nil {
		return nil, err
	}
	defer vdfFile.Close()

	parser := vdf.NewParser(vdfFile)
	content, err := parser.Parse()
	return content, err
}

func getVDFFilePath() (path string, err error) {
	steamPath := ""
	if runtime.GOOS == "windows" {
		steamPath = windowsSteamPath
	} else {
		steamPath = unixSteamPath
		if _, err := os.Stat(steamPath); os.IsNotExist(err) {
			steamPath = unixSteamAltPath
		}
	}
	vdfPath := filepath.Join(steamPath, libraryFoldersVDFRelativePath)
	if _, err := os.Stat(vdfPath); os.IsNotExist(err) {
		return "", err
	}
	return vdfPath, nil
}

func getBasePath(vdfContent map[string]any) string {
	directory := ""

	for _, data := range vdfContent["libraryfolders"].(map[string]any) {
		if directory != "" {
			break
		}
		library, ok := data.(map[string]any)
		if !ok {
			continue
		}
		apps, ok := library["apps"].(map[string]any)
		if !ok {
			continue
		}
		for appID := range apps {
			if appID == oldenEraID {
				path, ok := library["path"].(string)
				if !ok {
					break
				}
				directory = path
				break
			}
		}
	}

	return directory
}

func tryResolveGlob(pattern string) string {
	matches, err := filepath.Glob(pattern)
	if err != nil || len(matches) == 0 {
		return ""
	}

	for _, path := range matches {
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			return path
		}
	}

	return ""
}
