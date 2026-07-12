package helpers

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/Tariomka/hommoe_custom_templates/internal/common"
	"github.com/andygrunwald/vdf"
)

const oldenEraID = "3105440"

// FindOldenEraTemplatesDir tries to locate a template folder in common directories.
// Depending on the useInstallDir flag, it either looks for the official Steam install directory
// or tries to find the user directory and tries to resolve the custom template path from there.
func FindOldenEraTemplatesDir(useInstallDir bool) (string, error) {
	customTemplateRelativeGlob := filepath.Join(
		"AppData", "LocalLow", "Unfrozen", "HeroesOldenEra", "users", "*", "my_map_templates")

	if !useInstallDir && runtime.GOOS == "windows" {
		userPath, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}

		templatePathPattern := filepath.Join(userPath, customTemplateRelativeGlob)
		return resolveGlob(templatePathPattern)
	}

	content, err := getVDFContent()
	if err != nil {
		return "", err
	}

	directory := getBasePath(content)
	if directory == "" {
		return "", common.ErrGameInVDFNotFound
	}

	if !useInstallDir /*&& runtime.GOOS != "windows" is redundant here*/ {
		templatePathPattern := filepath.Join(
			directory,
			"steamapps", "compatdata", oldenEraID, "pfx", "drive_c", "users", "steamuser",
			customTemplateRelativeGlob)
		return resolveGlob(templatePathPattern)
	}

	directory = filepath.Join(directory,
		"steamapps",
		"common",
		"Heroes of Might and Magic Olden Era",
		"HeroesOldenEra_Data",
		"StreamingAssets",
		"map_templates")
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		return "", err
	}

	return directory, nil
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
	steamPath := getSteamPath()
	vdfPath := filepath.Join(steamPath, "steamapps", "libraryfolders.vdf")
	if _, err := os.Stat(vdfPath); os.IsNotExist(err) {
		return "", err
	}
	return vdfPath, nil
}

func getSteamPath() string {
	switch runtime.GOOS {
	case "windows":
		if registryPath := getSteamPathFromRegistry(); registryPath != "" {
			return registryPath
		}

		if programFilesPath := os.Getenv("ProgramFiles(x86)"); programFilesPath != "" {
			return filepath.Join(programFilesPath, "Steam")
		}

		return `C:\Program Files (x86)\Steam`

	default:
		userPath, err := os.UserHomeDir()
		if err != nil {
			return ""
		}

		steamPath := filepath.Join(userPath, ".local", "share", "Steam")
		if _, err := os.Stat(steamPath); os.IsNotExist(err) {
			steamPath = filepath.Join(userPath, ".steam", "steam")
		}
		return steamPath
	}
}

func getBasePath(vdfContent map[string]any) string {
	libraryFolders, ok := vdfContent["libraryfolders"].(map[string]any)
	if !ok {
		return ""
	}

	for _, data := range libraryFolders {
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

				return path
			}
		}
	}

	return ""
}

func resolveGlob(pattern string) (string, error) {
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return "", fmt.Errorf("resolve templates glob %q: %w", pattern, err)
	}

	for _, path := range matches {
		if _, statErr := os.Stat(path); statErr == nil {
			return path, nil
		}
	}

	return "", common.ErrTemplatesDirNotFound
}
