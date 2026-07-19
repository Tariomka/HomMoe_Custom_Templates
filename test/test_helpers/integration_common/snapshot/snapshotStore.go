package snapshot

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
)

const (
	goldenExtension  = ".golden"
	failureExtension = ".failure"
	snapshotsDirName = "__snapshots__"
)

// SnapshotStore resolves golden/failure snapshot paths and persists snapshots
// as PNG files. The default root is the __snapshots__ directory next to this
// source file, so paths are stable no matter which test package is running.
type SnapshotStore struct {
	root string
}

// NewSnapshotStore builds a store rooted at the package's __snapshots__ folder.
func NewSnapshotStore() SnapshotStore {
	_, sourceFile, _, _ := runtime.Caller(0)
	return SnapshotStore{root: filepath.Join(filepath.Dir(sourceFile), snapshotsDirName)}
}

// NewSnapshotStoreWithRoot builds a store rooted at an explicit directory
// (used by unit tests).
func NewSnapshotStoreWithRoot(root string) SnapshotStore {
	return SnapshotStore{root: root}
}

// GoldenPath returns the golden snapshot path for the given test file, test
// name and 1-based action number.
func (this SnapshotStore) GoldenPath(testFileName, testName string, actionNumber int) string {
	return this.snapshotPath(testFileName, testName, actionNumber, goldenExtension)
}

// FailurePath returns the failure snapshot path (same location and name as the
// golden, with the .failure extension).
func (this SnapshotStore) FailurePath(testFileName, testName string, actionNumber int) string {
	return this.snapshotPath(testFileName, testName, actionNumber, failureExtension)
}

// SaveGolden writes the screenshot as a PNG golden snapshot, creating parent
// directories as needed, and removes any stale failure snapshot beside it.
func (this SnapshotStore) SaveGolden(path string, screenshot image.Image) error {
	if err := this.savePNG(path, screenshot); err != nil {
		return err
	}
	return this.DeleteFailure(failurePathForGolden(path))
}

// LoadGolden decodes a previously saved golden snapshot.
func (this SnapshotStore) LoadGolden(path string) (image.Image, error) {
	goldenFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer goldenFile.Close()
	return png.Decode(goldenFile)
}

// SaveFailure writes the screenshot as a PNG failure snapshot beside the golden.
func (this SnapshotStore) SaveFailure(path string, screenshot image.Image) error {
	return this.savePNG(path, screenshot)
}

// DeleteFailure removes a failure snapshot if present; a missing file is not an
// error (used to clear stale failures after a passing validation).
func (this SnapshotStore) DeleteFailure(path string) error {
	err := os.Remove(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func (this SnapshotStore) snapshotPath(testFileName, testName string, actionNumber int, extension string) string {
	fileName := fmt.Sprintf("%s_%d%s", sanitizeSnapshotName(testName), actionNumber, extension)
	return filepath.Join(this.root, sanitizeSnapshotName(testFileName), fileName)
}

func (this SnapshotStore) savePNG(path string, screenshot image.Image) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		return err
	}
	snapshotFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer snapshotFile.Close()
	return png.Encode(snapshotFile, screenshot)
}

// failurePathForGolden swaps the .golden extension for .failure.
func failurePathForGolden(goldenPath string) string {
	return goldenPath[:len(goldenPath)-len(goldenExtension)] + failureExtension
}

// sanitizeSnapshotName makes a test or file name safe as a path segment on
// every supported OS by collapsing path-hostile characters (subtest slashes,
// spaces, Windows-reserved punctuation) into underscores.
func sanitizeSnapshotName(name string) string {
	sanitizer := regexp.MustCompile(`[^A-Za-z0-9_.-]+`)
	return sanitizer.ReplaceAllString(name, "_")
}
