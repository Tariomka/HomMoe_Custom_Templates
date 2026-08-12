// Command testlayoutcheck verifies the build-tag placement rules of AGENTS.md 4.6.1 and 4.6.2.
//
// Usage: testlayoutcheck [root]   (root defaults to the working directory)
package main

import (
	"fmt"
	"os"

	"github.com/Tariomka/hommoe_custom_templates/cmd/testlayoutcheck/checker"
)

func main() {
	root := "."
	if len(os.Args) > 1 {
		root = os.Args[1]
	}

	violations, err := checker.NewTestLayoutChecker().Check(root)
	if err != nil {
		fmt.Fprintln(os.Stderr, "test-layout check failed:", err)
		os.Exit(2)
	}

	for _, violation := range violations {
		fmt.Fprintf(os.Stderr, "%s: %s: %s\n", violation.Path, violation.Rule, violation.Detail)
	}

	if len(violations) > 0 {
		fmt.Fprintf(os.Stderr, "%d test-layout violation(s) found\n", len(violations))
		os.Exit(1)
	}

	fmt.Fprintln(os.Stdout, "test-layout check passed")
}
