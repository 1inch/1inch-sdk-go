package validate

//
//import (
//	"fmt"
//	"go/ast"
//	"go/parser"
//	"go/token"
//	"os"
//	"path/filepath"
//	"strings"
//	"testing"
//)
//
//// This is a project validation test to ensure all parameter validation functions in the SDK properly pair each parameter name with the correct string
//
//func TestParameterConsistency(t *testing.T) {
//
//	err := validateWorkingDirectory()
//	if err != nil {
//		t.Fatalf("Directory error: %v", err)
//	}
//
//	err = filepath.Walk("./..", func(path string, info os.FileInfo, err error) error {
//		if err != nil {
//			return err
//		}
//		if strings.HasSuffix(path, "types.go") {
//			if err := checkFile(t, path); err != nil {
//				t.Errorf("Failed to check file %v: %v", trimPath(path), err)
//			}
//		}
//		return nil
//	})
//	if err != nil {
//		t.Fatalf("Error walking through files: %v", err)
//	}
//}
//
//func checkFile(t *testing.T, filePath string) error {
//	fset := token.NewFileSet()
//	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
//	if err != nil {
//		return fmt.Errorf("error parsing file: %w", err)
//	}
//
//	var lastErr error
//	ast.Inspect(node, func(n ast.Node) bool {
//		callExpr, ok := n.(*ast.CallExpr)
//		if !ok {
//			return true
//		}
//
//		if selExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
//			if ident, ok := selExpr.X.(*ast.Ident); ok && ident.Name == "validate" && selExpr.Sel.Name == "Parameter" {
//				if len(callExpr.Args) >= 4 {
//					firstArg, secondArg := callExpr.Args[0], callExpr.Args[1]
//					if err := checkValidationCall(firstArg, secondArg, fset); err != nil {
//						t.Errorf("%v", err)
//					}
//				}
//			}
//		}
//
//		return true
//	})
//
//	return lastErr
//}
//
//func checkValidationCall(firstArg, secondArg ast.Expr, fset *token.FileSet) error {
//	// Extract variable name from the first argument
//	varName := extractVarName(firstArg)
//
//	// Extract string literal from the second argument
//	if stringLit, ok := secondArg.(*ast.BasicLit); ok && stringLit.Kind == token.STRING {
//		// Remove quotes from the string literal
//		expectedVarName := strings.Trim(stringLit.Value, "\"")
//
//		// Check if the variable name part of the first argument matches the expected variable name with the first letter lowercase
//		if len(varName) > 0 && len(expectedVarName) > 0 {
//			// Take the last part of the variable name after '.'
//			splitVarName := strings.Split(varName, ".")
//			actualVarNamePart := splitVarName[len(splitVarName)-1]
//
//			// Make the first letter of actualVarNamePart lowercase
//			formattedVarName := strings.ToLower(string(actualVarNamePart[0])) + actualVarNamePart[1:]
//
//			if formattedVarName != expectedVarName {
//				return reportMismatch(varName, expectedVarName, stringLit, fset, formattedVarName)
//			}
//		}
//	}
//	return nil
//}
//
//func extractVarName(expr ast.Expr) string {
//	switch v := expr.(type) {
//	case *ast.Ident:
//		return v.Name
//	case *ast.SelectorExpr:
//		if ident, ok := v.X.(*ast.Ident); ok {
//			return ident.Name + "." + v.Sel.Name
//		}
//	case *ast.CallExpr:
//		// This block handles type conversion expressions.
//		// Check if the Fun part of the CallExpr is an identifier, which would indicate a type conversion.
//		if ident, ok := v.Fun.(*ast.Ident); ok {
//			// The first argument of the CallExpr should be the variable being converted.
//			if len(v.Args) > 0 {
//				return extractVarName(v.Args[0]) // Recursively extract the name from the first argument.
//			}
//			return ident.Name // If no arguments are found, return the type being converted to.
//		}
//	}
//	return ""
//}
//
//func reportMismatch(varName, varNameAsString string, stringLit *ast.BasicLit, fset *token.FileSet, formattedVarName string) error {
//	position := fset.Position(stringLit.Pos())
//	displayPath := trimPath(position.Filename)
//	return fmt.Errorf("Mismatch found in %s at line %d: parameter '%s' should have the string literal '%s' next to it, not '%s'.\n",
//		vanityPath(displayPath), position.Line, varName, formattedVarName, varNameAsString)
//}
//
//func trimPath(path string) string {
//	// Split the file path on "1inch-sdk" and use the part after it
//	parts := strings.Split(path, "1inch-sdk/")
//	if len(parts) > 1 {
//		return "1inch-sdk/" + parts[1]
//	}
//	return path
//}
