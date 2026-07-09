package string_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
	"github.com/stretchr/testify/assert"
)

func TestWhenNameIsAlreadyClean_ReturnsItUnchanged(t *testing.T) {
	// Arrange
	name := "clean_name-1"

	// Act
	sanitized := helpers.SanitizeFilename(name)

	// Assert
	assert.Equal(t, "clean_name-1", sanitized)
}

func TestWhenNameHasSurroundingWhitespace_TrimsIt(t *testing.T) {
	// Arrange
	name := "   spaced   "

	// Act
	sanitized := helpers.SanitizeFilename(name)

	// Assert
	assert.Equal(t, "spaced", sanitized)
}

func TestWhenNameContainsUnsafeRune_ReplacesItWithUnderscore(t *testing.T) {
	testCases := []struct {
		subtestName string
		input       string
	}{
		{"WhenNameContainsForwardSlash_ReplacesItWithUnderscore", "a/b"},
		{"WhenNameContainsBackslash_ReplacesItWithUnderscore", `a\b`},
		{"WhenNameContainsColon_ReplacesItWithUnderscore", "a:b"},
		{"WhenNameContainsAsterisk_ReplacesItWithUnderscore", "a*b"},
		{"WhenNameContainsQuestionMark_ReplacesItWithUnderscore", "a?b"},
		{"WhenNameContainsDoubleQuote_ReplacesItWithUnderscore", `a"b`},
		{"WhenNameContainsLessThan_ReplacesItWithUnderscore", "a<b"},
		{"WhenNameContainsGreaterThan_ReplacesItWithUnderscore", "a>b"},
		{"WhenNameContainsPipe_ReplacesItWithUnderscore", "a|b"},
	}
	for _, testCase := range testCases {
		t.Run(testCase.subtestName, func(t *testing.T) {
			// Arrange
			name := testCase.input

			// Act
			sanitized := helpers.SanitizeFilename(name)

			// Assert
			assert.Equal(t, "a_b", sanitized)
		})
	}
}

func TestWhenNameContainsEveryUnsafeRune_ReplacesAllOfThem(t *testing.T) {
	// Arrange
	name := `bad/\*?":<>|name`

	// Act
	sanitized := helpers.SanitizeFilename(name)

	// Assert
	assert.Equal(t, "bad_________name", sanitized)
}

func TestWhenNameIsEmpty_ReturnsEmptyString(t *testing.T) {
	// Arrange
	name := ""

	// Act
	sanitized := helpers.SanitizeFilename(name)

	// Assert
	assert.Empty(t, sanitized)
}

func TestWhenNameIsOnlyWhitespace_ReturnsEmptyString(t *testing.T) {
	// Arrange
	name := "   "

	// Act
	sanitized := helpers.SanitizeFilename(name)

	// Assert
	assert.Empty(t, sanitized)
}
