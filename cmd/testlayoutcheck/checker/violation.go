package checker

// Violation is a single build-tag placement that contradicts AGENTS.md 4.6.1 or 4.6.2.
type Violation struct {
	// Path is the offending file, relative to the checked root, with forward slashes.
	Path string
	// Rule identifies which of the four enforced rules was broken.
	Rule string
	// Detail explains the breach and how to fix it.
	Detail string
}
