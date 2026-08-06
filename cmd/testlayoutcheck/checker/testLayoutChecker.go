package checker

import (
	"fmt"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

const (
	integrationTestTag = "integration_test"
	guiTag             = "gui"

	testFileSuffix    = "_test.go"
	testExportsSuffix = "_testexports.go"

	unitTestDirectory       = "test/unit/"
	guiIntegrationDirectory = "test/integration/gui/"
	testDirectory           = "test/"

	// RuleMissingIntegrationTag covers AGENTS.md 4.6.1: consuming a test-only export requires the tag.
	RuleMissingIntegrationTag = "missing-integration-tag"
	// RuleTaggedUnitTest covers AGENTS.md 4.6.1: unit tests must never see test-only exports.
	RuleTaggedUnitTest = "tagged-unit-test"
	// RuleMissingGuiTag covers AGENTS.md 4.6.2: GPU-dependent tests must be opt-in.
	RuleMissingGuiTag = "missing-gui-tag"
	// RuleTaggedProductionFile covers AGENTS.md 4.6.1: only *_testexports.go may carry the tag.
	RuleTaggedProductionFile = "tagged-production-file"
)

// TestLayoutChecker enforces the build-tag placement rules of AGENTS.md 4.6.1 and 4.6.2.
type TestLayoutChecker struct {
	fileSet *token.FileSet
}

// NewTestLayoutChecker creates a checker with its own file set.
func NewTestLayoutChecker() *TestLayoutChecker {
	return &TestLayoutChecker{fileSet: token.NewFileSet()}
}

// Check walks root and reports every misplaced or missing integration_test / gui build constraint.
func (this *TestLayoutChecker) Check(root string) ([]Violation, error) {
	files, err := this.parseTree(root)
	if err != nil {
		return nil, err
	}

	accessors := collectAccessorNames(files)

	violations := make([]Violation, 0)
	for _, file := range files {
		violations = append(violations, inspectFile(file, accessors)...)
	}

	return violations, nil
}

func (this *TestLayoutChecker) parseTree(root string) ([]*goFile, error) {
	files := make([]*goFile, 0)
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		if entry.IsDir() {
			if isSkippedDirectory(path, root, entry.Name()) {
				return filepath.SkipDir
			}

			return nil
		}

		if !strings.HasSuffix(entry.Name(), ".go") {
			return nil
		}

		file, parseErr := this.parseFile(root, path)
		if parseErr != nil {
			return parseErr
		}

		files = append(files, file)

		return nil
	})
	if err != nil {
		return nil, err
	}

	return files, nil
}

func (this *TestLayoutChecker) parseFile(root string, path string) (*goFile, error) {
	source, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	syntax, err := parser.ParseFile(this.fileSet, path, source, parser.ParseComments|parser.SkipObjectResolution)
	if err != nil {
		return nil, fmt.Errorf("parsing %s: %w", path, err)
	}

	relative, err := filepath.Rel(root, path)
	if err != nil {
		return nil, err
	}

	return &goFile{
		relativePath: filepath.ToSlash(relative),
		name:         filepath.Base(path),
		constraints:  parseConstraints(string(source)),
		syntax:       syntax,
	}, nil
}

// collectAccessorNames gathers every function and method declared in a *_testexports.go file.
func collectAccessorNames(files []*goFile) map[string]struct{} {
	accessors := make(map[string]struct{})
	for _, file := range files {
		if !file.isTestExportsFile() {
			continue
		}

		for _, name := range file.declaredFunctionNames() {
			accessors[name] = struct{}{}
		}
	}

	return accessors
}

func inspectFile(file *goFile, accessors map[string]struct{}) []Violation {
	violations := make([]Violation, 0)
	requiresIntegrationTag := file.requiresTag(integrationTestTag)

	if requiresIntegrationTag && file.isUnder(unitTestDirectory) {
		violations = append(violations, newViolation(file, RuleTaggedUnitTest,
			"unit tests assert production code and must never carry //go:build "+integrationTestTag))
	}

	if requiresIntegrationTag && !file.isUnder(testDirectory) && !file.isTestExportsFile() {
		violations = append(violations, newViolation(file, RuleTaggedProductionFile,
			"only *_testexports.go may carry //go:build "+integrationTestTag+" outside test/"))
	}

	if file.isUnder(guiIntegrationDirectory) && !file.requiresTag(guiTag) {
		violations = append(violations, newViolation(file, RuleMissingGuiTag,
			"GPU-dependent tests must require the "+guiTag+" tag so CI can exclude them"))
	}

	if accessor := file.firstReferencedName(accessors); accessor != "" &&
		file.isTestFile() && !requiresIntegrationTag {
		violations = append(violations, newViolation(file, RuleMissingIntegrationTag,
			"references the test-only export "+accessor+", so it must carry //go:build "+integrationTestTag))
	}

	return violations
}

func newViolation(file *goFile, rule string, detail string) Violation {
	return Violation{Path: file.relativePath, Rule: rule, Detail: detail}
}

func isSkippedDirectory(path string, root string, name string) bool {
	if path == root {
		return false
	}

	return strings.HasPrefix(name, ".") || name == "data" || name == "output" || name == "tmp"
}
