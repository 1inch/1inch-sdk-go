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
	err := checkParameterFunctionConsistency("validate.go")
	assert.NoError(t, err, "Parameter function consistency check failed")
}

func checkParameterFunctionConsistency(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	funcRegex := regexp.MustCompile(`^func Check(\w+)\(`)
	scanner := bufio.NewScanner(file)

	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()

		matches := funcRegex.FindStringSubmatch(line)
		if len(matches) == 2 {
			caseLabel := matches[1]

			// Move the scanner four lines down
			for i := 0; i < 3; i++ {
				if !scanner.Scan() {
					return fmt.Errorf("expected error message four lines down from case '%s' at line %d, but reached EOF", caseLabel, lineNumber)
				}
				lineNumber++
			}

			errorLine := scanner.Text()
			if !strings.HasSuffix(errorLine, fmt.Sprintf(`"%s")`, caseLabel)) {

				finalStringLiteralRegex := regexp.MustCompile(`"([^"]+)"\s*\)\s*$`)
				stringLiteral := finalStringLiteralRegex.FindStringSubmatch(errorLine)

				return fmt.Errorf("mismatch found at line %d: case '%s' should be used in the error message. Have '%s', want '%s'", lineNumber, caseLabel, stringLiteral[1], caseLabel)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error scanning file: %v", err)
	}

	return nil
}
