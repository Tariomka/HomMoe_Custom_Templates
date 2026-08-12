package stateHandler_test

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/validators"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/mock"
)

// newPassingValidator returns a validator mock that reports no issues.
func newPassingValidator() *test_helpers.EditorStateValidatorMock {
	validator := &test_helpers.EditorStateValidatorMock{}
	validator.On("Validate", mock.Anything).Return([]validators.ValidationIssue{})
	return validator
}

// newValidatorReporting returns a validator mock that reports one issue per
// message. The issues carry no fix, so they may only be used with
// fixIssues=false.
func newValidatorReporting(messages ...string) *test_helpers.EditorStateValidatorMock {
	issues := make([]validators.ValidationIssue, 0, len(messages))
	for _, message := range messages {
		issues = append(issues, validators.ValidationIssue{Message: message})
	}

	validator := &test_helpers.EditorStateValidatorMock{}
	validator.On("Validate", mock.Anything).Return(issues)
	return validator
}
