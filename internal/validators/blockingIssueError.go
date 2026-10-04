package validators

// blockingIssueError carries a blocking issue's message verbatim while still
// matching its cause through [errors.Is].
type blockingIssueError struct {
	message string
	cause   error
}

func newBlockingIssueError(message string, cause error) error {
	return &blockingIssueError{message: message, cause: cause}
}

func (this *blockingIssueError) Error() string {
	return this.message
}

func (this *blockingIssueError) Unwrap() error {
	return this.cause
}
