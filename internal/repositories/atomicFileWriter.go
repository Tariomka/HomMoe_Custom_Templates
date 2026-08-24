package repositories

import (
	jsonV1 "encoding/json"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/constants"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
)

const (
	temporaryFilePrefix = "TEMP-"
	defaultName         = "Generated_Template"
	renameAttempts      = 5
	renameRetryDelay    = 20 * time.Millisecond
	jsonIndent          = "  "
)

type atomicFileWriter struct{}

func newAtomicFileWriter() *atomicFileWriter {
	return &atomicFileWriter{}
}

// Write encodes into "{directory}/TEMP-{fileName}{extension}" and renames it
// onto "{directory}/{fileName}{extension}", returning the destination path.
func (this *atomicFileWriter) Write(
	directory, fileName, extension string,
	encode func(file *os.File) error) (string, error) {
	if err := os.MkdirAll(directory, constants.FolderPermission); err != nil {
		return "", err
	}

	fileName = this.resolveFileName(fileName)
	destinationPath := filepath.Join(directory, fileName+extension)
	temporaryPath := filepath.Join(directory, temporaryFilePrefix+fileName+extension)

	if err := this.encodeToTemporaryFile(temporaryPath, encode); err != nil {
		this.discard(temporaryPath)
		return "", err
	}

	if err := this.commit(temporaryPath, destinationPath); err != nil {
		this.discard(temporaryPath)
		return "", err
	}

	return destinationPath, nil
}

// WriteJSON marshals value as indented JSON through write.
func (this *atomicFileWriter) WriteJSON(
	directory, fileName, extension string,
	value any) (string, error) {
	return this.Write(directory, fileName, extension, func(file *os.File) error {
		return json.MarshalWrite(
			file, value,
			jsontext.WithIndent(jsonIndent),
			jsonV1.OmitEmptyWithLegacySemantics(true),
			json.FormatNilSliceAsNull(true),
			json.FormatNilMapAsNull(true))
	})
}

func (this *atomicFileWriter) encodeToTemporaryFile(
	temporaryPath string,
	encode func(file *os.File) error) (err error) {
	file, err := os.OpenFile(temporaryPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, constants.FilePermission)
	if err != nil {
		return err
	}

	// A buffered write only reports an out-of-space condition on Close, so the
	// close error has to win when the encode itself succeeded.
	defer func() {
		if closeErr := file.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	if err = encode(file); err != nil {
		return err
	}

	return file.Sync()
}

// commit retries the rename because on Windows it fails while another process
// (the game, an editor) still holds the destination open.
func (this *atomicFileWriter) commit(temporaryPath, destinationPath string) error {
	var err error
	for attempt := range renameAttempts {
		if err = os.Rename(temporaryPath, destinationPath); err == nil {
			return nil
		}

		if attempt < renameAttempts-1 {
			time.Sleep(renameRetryDelay)
		}
	}
	return fmt.Errorf("could not replace %s: %w", destinationPath, err)
}

// discard drops the temporary file; the destination still holds the previous
// valid contents, so losing the half-written copy is the correct outcome.
func (this *atomicFileWriter) discard(temporaryPath string) { _ = os.Remove(temporaryPath) }

func (this *atomicFileWriter) resolveFileName(fileName string) string {
	safeName := helpers.SanitizeFilename(fileName)
	if safeName == "" {
		return defaultName
	}

	return safeName
}
