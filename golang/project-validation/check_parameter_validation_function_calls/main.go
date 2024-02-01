package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// This is a helper script to ensure all parameter validation functions in the SDK properly pair each parameter name with the correct string
// Can only be run when the working directory is set to the "golang" folder of the SDK

func main() {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Failed to get current directory:", err)
		return
	}

	dirName := filepath.Base(currentDir)
	if dirName != "golang" {
		fmt.Println(`This script must be run specifically from the "golang" folder of the SDK project.`)
		return
	}

	err = filepath.Walk(currentDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(path, "types.go") {
			fmt.Println("Checking file:", trimPath(path))
			checkFile(path)
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error walking through files:", err)
	}
}

func checkFile(filename string) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		fmt.Println("Error parsing file:", err)
		return
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
					checkValidationCall(firstArg, secondArg, fset)
				}
			}
		}

		return true
	})
}

func checkValidationCall(firstArg, secondArg ast.Expr, fset *token.FileSet) {
	// Extract variable name from the first argument
	varName := extractVarName(firstArg)

	// Extract string literal from the second argument
	if stringLit, ok := secondArg.(*ast.BasicLit); ok && stringLit.Kind == token.STRING {
		// Remove quotes from the string literal
		expectedVarName := strings.Trim(stringLit.Value, "\"")

		// Check if the variable name part of the first argument matches the expected variable name with the first letter lowercase
		if len(varName) > 0 && len(expectedVarName) > 0 {
			// Take the last part of the variable name after '.'
			splitVarName := strings.Split(varName, ".")
			actualVarNamePart := splitVarName[len(splitVarName)-1]

			// Make the first letter of actualVarNamePart lowercase
			formattedVarName := strings.ToLower(string(actualVarNamePart[0])) + actualVarNamePart[1:]

			if formattedVarName != expectedVarName {
				reportMismatch(varName, expectedVarName, stringLit, fset, formattedVarName)
			}
		}
	}
}

func extractVarName(expr ast.Expr) string {
	switch v := expr.(type) {
	case *ast.Ident:
		return v.Name
	case *ast.SelectorExpr:
		if ident, ok := v.X.(*ast.Ident); ok {
			return ident.Name + "." + v.Sel.Name
		}
	}
	return ""
}

func reportMismatch(varName, varNameAsString string, stringLit *ast.BasicLit, fset *token.FileSet, formattedVarName string) {
	position := fset.Position(stringLit.Pos())
	displayPath := trimPath(position.Filename)
	fmt.Printf("Mismatch found in %s at line %d: parameter '%s' should have the string literal '%s' next to it, not '%s'.\n",
		displayPath, position.Line, varName, formattedVarName, varNameAsString)
}

func trimPath(path string) string {
	// Split the file path on "1inch-sdk" and use the part after it
	parts := strings.Split(path, "1inch-sdk/")
	if len(parts) > 1 {
		return "1inch-sdk/" + parts[1]
	}
	return path
}
