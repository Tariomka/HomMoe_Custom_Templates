package checker

import (
	"go/ast"
	"go/build/constraint"
	"strings"
)

// goFile is one parsed source file together with the build constraint that gates it.
type goFile struct {
	relativePath string
	name         string
	constraints  constraint.Expr
	syntax       *ast.File
}

// requiresTag reports whether the file drops out of the build when tag is absent.
func (this *goFile) requiresTag(tag string) bool {
	if this.constraints == nil || !mentionsTag(this.constraints, tag) {
		return false
	}

	return !this.constraints.Eval(func(candidate string) bool { return candidate != tag })
}

func (this *goFile) isTestFile() bool {
	return strings.HasSuffix(this.name, testFileSuffix)
}

func (this *goFile) isTestExportsFile() bool {
	return strings.HasSuffix(this.name, testExportsSuffix)
}

func (this *goFile) isUnder(directory string) bool {
	return strings.HasPrefix(this.relativePath, directory)
}

// declaredFunctionNames lists every function and method declared in the file.
func (this *goFile) declaredFunctionNames() []string {
	names := make([]string, 0, len(this.syntax.Decls))
	for _, declaration := range this.syntax.Decls {
		if function, ok := declaration.(*ast.FuncDecl); ok {
			names = append(names, function.Name.Name)
		}
	}

	return names
}

// firstReferencedName returns the first selector in the file whose name is in names, or "".
func (this *goFile) firstReferencedName(names map[string]struct{}) string {
	found := ""
	ast.Inspect(this.syntax, func(node ast.Node) bool {
		if found != "" {
			return false
		}

		selector, ok := node.(*ast.SelectorExpr)
		if !ok {
			return true
		}

		if _, exists := names[selector.Sel.Name]; exists {
			found = selector.Sel.Name
		}

		return found == ""
	})

	return found
}

// parseConstraints extracts the //go:build expression preceding the package clause.
func parseConstraints(source string) constraint.Expr {
	for line := range strings.SplitSeq(source, "\n") {
		trimmed := strings.TrimSpace(strings.TrimSuffix(line, "\r"))
		if strings.HasPrefix(trimmed, "package ") {
			return nil
		}

		if !constraint.IsGoBuild(trimmed) {
			continue
		}

		if expression, err := constraint.Parse(trimmed); err == nil {
			return expression
		}
	}

	return nil
}

// mentionsTag reports whether tag appears anywhere in the constraint expression.
func mentionsTag(expression constraint.Expr, tag string) bool {
	switch typed := expression.(type) {
	case *constraint.TagExpr:
		return typed.Tag == tag
	case *constraint.NotExpr:
		return mentionsTag(typed.X, tag)
	case *constraint.AndExpr:
		return mentionsTag(typed.X, tag) || mentionsTag(typed.Y, tag)
	case *constraint.OrExpr:
		return mentionsTag(typed.X, tag) || mentionsTag(typed.Y, tag)
	default:
		return false
	}
}
