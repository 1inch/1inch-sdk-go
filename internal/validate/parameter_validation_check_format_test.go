package validate

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// This is a project validation test to ensure all parameter validation functions in the SDK properly pair each parameter name with the correct string

func TestParameterConsistency(t *testing.T) {
	// Find the project root by looking for go.mod
	projectRoot, err := findProjectRoot()
	if err != nil {
		t.Fatalf("Failed to find project root: %v", err)
	}

	// Walk through sdk-clients directory looking for validation files
	sdkClientsPath := filepath.Join(projectRoot, "sdk-clients")
	err = filepath.Walk(sdkClientsPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Check validation.go and validate.go files
		if strings.HasSuffix(path, "validation.go") || strings.HasSuffix(path, "validate.go") {
			if err := checkFile(t, path, projectRoot); err != nil {
				t.Errorf("Failed to check file %v: %v", trimPath(path, projectRoot), err)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("Error walking through files: %v", err)
	}
}

func findProjectRoot() (string, error) {
	// Start from current directory and walk up until we find go.mod
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("could not find go.mod in any parent directory")
		}
		dir = parent
	}
}

func checkFile(t *testing.T, filePath string, projectRoot string) error {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("error parsing file: %w", err)
	}

	ast.Inspect(node, func(n ast.Node) bool {
		callExpr, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		if selExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
			if ident, ok := selExpr.X.(*ast.Ident); ok && ident.Name == "validate" && selExpr.Sel.Name == "Parameter" {
				if len(callExpr.Args) >= 4 {
					firstArg, secondArg := callExpr.Args[0], callExpr.Args[1]
					if err := checkValidationCall(firstArg, secondArg, fset, projectRoot); err != nil {
						t.Errorf("%v", err)
					}
				}
			}
		}

		return true
	})

	return nil
}

func checkValidationCall(firstArg, secondArg ast.Expr, fset *token.FileSet, projectRoot string) error {
	// Extract variable name from the first argument
	varName := extractVarName(firstArg)

	// Extract string literal from the second argument
	if stringLit, ok := secondArg.(*ast.BasicLit); ok && stringLit.Kind == token.STRING {
		// Remove quotes from the string literal
		expectedVarName := strings.Trim(stringLit.Value, "\"")

		// Check if the variable name part of the first argument matches the expected variable name (case-insensitive)
		if len(varName) > 0 && len(expectedVarName) > 0 {
			// Take the last part of the variable name after '.'
			splitVarName := strings.Split(varName, ".")
			actualVarNamePart := splitVarName[len(splitVarName)-1]

			// Skip known exceptions:
			// 1. Single-letter loop variables (like 'v' in a for range loop)
			if len(actualVarNamePart) == 1 {
				return nil
			}
			// 2. Variables that have conversion suffixes (like statusesAsInts -> statuses)
			if strings.HasSuffix(strings.ToLower(actualVarNamePart), "asints") ||
				strings.HasSuffix(strings.ToLower(actualVarNamePart), "asstrings") {
				return nil
			}

			// Compare case-insensitively to catch actual mismatches (not just case differences)
			if !strings.EqualFold(actualVarNamePart, expectedVarName) {
				return reportMismatch(varName, expectedVarName, stringLit, fset, actualVarNamePart, projectRoot)
			}
		}
	}
	return nil
}

func extractVarName(expr ast.Expr) string {
	switch v := expr.(type) {
	case *ast.Ident:
		return v.Name
	case *ast.SelectorExpr:
		if ident, ok := v.X.(*ast.Ident); ok {
			return ident.Name + "." + v.Sel.Name
		}
	case *ast.CallExpr:
		// This block handles type conversion expressions.
		// Check if the Fun part of the CallExpr is an identifier, which would indicate a type conversion.
		if _, ok := v.Fun.(*ast.Ident); ok {
			// The first argument of the CallExpr should be the variable being converted.
			if len(v.Args) > 0 {
				return extractVarName(v.Args[0]) // Recursively extract the name from the first argument.
			}
		}
	}
	return ""
}

func reportMismatch(varName, varNameAsString string, stringLit *ast.BasicLit, fset *token.FileSet, formattedVarName string, projectRoot string) error {
	position := fset.Position(stringLit.Pos())
	displayPath := trimPath(position.Filename, projectRoot)
	return fmt.Errorf("mismatch found in %s at line %d: parameter '%s' should have the string literal '%s' next to it, not '%s'",
		displayPath, position.Line, varName, formattedVarName, varNameAsString)
}

func trimPath(path string, projectRoot string) string {
	// Make path relative to project root
	relPath, err := filepath.Rel(projectRoot, path)
	if err != nil {
		return path
	}
	return relPath
}
