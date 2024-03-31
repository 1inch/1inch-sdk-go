package validate

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParameterFunctionConsistency(t *testing.T) {
	validateFile := "validate.go"
	funcCount, err := countFunctions(validateFile)
	assert.NoError(t, err, "Failed to count functions in validate.go")
	funcsCheckedCount, err := checkParameterFunctionConsistency(validateFile)
	assert.NoError(t, err, "Parameter function consistency check failed")
	if funcCount != funcsCheckedCount {
		t.Errorf("Expected %d functions to be checked, but found %d", funcCount, funcsCheckedCount)
	}
}

func countFunctions(filePath string) (int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	funcRegex := regexp.MustCompile(`^func\s+(\w+)`)
	scanner := bufio.NewScanner(file)

	count := 0
	for scanner.Scan() {
		line := scanner.Text()
		if funcRegex.MatchString(line) {
			count++
		}
	}

	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("error scanning file: %v", err)
	}

	return count, nil
}

func checkParameterFunctionConsistency(filePath string) (int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return -1, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	funcRegex := regexp.MustCompile(`^func Check(\w+)(Required)?\(`)
	scanner := bufio.NewScanner(file)

	var totalFunctionsChecked int
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()

		matches := funcRegex.FindStringSubmatch(line)
		if len(matches) > 1 { // Ensure at least the function name is captured
			totalFunctionsChecked++

			// Initialize caseLabel with the first capture group
			caseLabel := matches[1]
			
			caseLabel = strings.TrimSuffix(caseLabel, "Required")

			// Move the scanner four lines down
			for i := 0; i < 3; i++ {
				if !scanner.Scan() {
					return -1, fmt.Errorf("expected error message four lines down from case '%s' at line %d, but reached EOF", caseLabel, lineNumber)
				}
				lineNumber++
			}

			errorLine := scanner.Text()
			if !strings.HasSuffix(errorLine, fmt.Sprintf(`"%s")`, caseLabel)) {
				finalStringLiteralRegex := regexp.MustCompile(`"([^"]+)"\s*\)\s*$`)
				stringLiteral := finalStringLiteralRegex.FindStringSubmatch(errorLine)

				if len(stringLiteral) < 2 {
					return -1, fmt.Errorf("expected a string literal at the end of the error message at line %d, but none was found", lineNumber)
				}

				return -1, fmt.Errorf("mismatch found at line %d: case '%s' should be used in the error message. Have '%s', want '%s'", lineNumber, caseLabel, stringLiteral[1], caseLabel)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return -1, fmt.Errorf("error scanning file: %v", err)
	}

	if totalFunctionsChecked == 0 {
		return 0, fmt.Errorf("regex did not find any functions to check")
	}

	return totalFunctionsChecked, nil
}
