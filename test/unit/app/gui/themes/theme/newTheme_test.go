package theme_test

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"gioui.org/font/gofont"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestWhenAppSourcesContainNonAsciiRunes_AllAreCoveredByTheBundledFonts guards the
// decision documented on themes.NewTheme: the shaper is built with
// text.NoSystemFonts, so a rune the Go collection does not carry renders as
// .notdef instead of silently borrowing an OS face whose metrics shift every row
// below it. The scan is static because the offending literal is usually added far
// from this package.
func TestWhenAppSourcesContainNonAsciiRunes_AllAreCoveredByTheBundledFonts(t *testing.T) {
	t.Parallel()

	// Arrange
	appDirectory := filepath.Join(repositoryRoot(t), "app")
	literalRunes, err := collectNonASCIIRunes(appDirectory)
	require.NoError(t, err)
	require.NotEmpty(t, literalRunes, "scan found no non-ASCII runes at all, which means it is not working")

	// Act
	uncovered := make([]string, 0)
	for _, occurrence := range literalRunes {
		if !isCoveredByBundledFonts(occurrence.value) {
			uncovered = append(uncovered, occurrence.String())
		}
	}

	// Assert
	assert.Empty(t, uncovered, "these runes are missing from gofont.Collection() and will render as .notdef")
}

type runeOccurrence struct {
	value    rune
	position string
}

func (this runeOccurrence) String() string {
	return fmt.Sprintf("U+%04X %q at %s", this.value, this.value, this.position)
}

func repositoryRoot(t *testing.T) string {
	t.Helper()

	_, thisFile, _, ok := runtime.Caller(0)
	require.True(t, ok, "unable to locate the test file on disk")

	// <root>/test/unit/app/gui/themes/theme/newTheme_test.go
	return filepath.Join(filepath.Dir(thisFile), "..", "..", "..", "..", "..", "..")
}

func collectNonASCIIRunes(root string) ([]runeOccurrence, error) {
	fileSet := token.NewFileSet()
	occurrences := make([]runeOccurrence, 0)

	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".go") {
			return nil
		}

		syntax, parseErr := parser.ParseFile(fileSet, path, nil, parser.SkipObjectResolution)
		if parseErr != nil {
			return fmt.Errorf("parsing %s: %w", path, parseErr)
		}

		ast.Inspect(syntax, func(node ast.Node) bool {
			literal, isLiteral := node.(*ast.BasicLit)
			if !isLiteral || (literal.Kind != token.STRING && literal.Kind != token.CHAR) {
				return true
			}

			text, unquoteErr := strconv.Unquote(literal.Value)
			if unquoteErr != nil {
				return true
			}

			for _, value := range text {
				if value < 0x80 {
					continue
				}

				occurrences = append(occurrences, runeOccurrence{
					value:    value,
					position: fileSet.Position(literal.Pos()).String(),
				})
			}

			return true
		})

		return nil
	})
	if err != nil {
		return nil, err
	}

	return occurrences, nil
}

// isCoveredByBundledFonts asks every face of the collection themes.NewTheme hands
// to the shaper whether it can map the rune to a real glyph.
func isCoveredByBundledFonts(value rune) bool {
	for _, collectionFace := range gofont.Collection() {
		if _, ok := collectionFace.Face.Face().NominalGlyph(value); ok {
			return true
		}
	}

	return false
}
