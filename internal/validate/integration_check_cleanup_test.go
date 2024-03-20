package validate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const targetBlock = "helpers.Sleep()"

func TestIntegrationSleepCleanup(t *testing.T) {

	err := validateWorkingDirectory()
	if err != nil {
		t.Fatalf("Directory error: %v", err)
	}

	err = filepath.Walk("./..", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasSuffix(path, "_integration_test.go") || strings.HasSuffix(path, "_e2e_test.go") {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			inTRunBlock := false
			openBracesCount := 0
			tRunContent := ""

			lines := strings.Split(string(content), "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)

				if strings.HasPrefix(line, "t.Run(") {
					inTRunBlock = true
				}

				if inTRunBlock {
					tRunContent += line + "\n"
					openBracesCount += strings.Count(line, "{") - strings.Count(line, "}")

					if openBracesCount == 0 {
						if !strings.Contains(tRunContent, targetBlock) {
							t.Errorf(`File: %s - Missing cleanup block in one of the tests. Each test must have the following code in it:
	t.Cleanup(func() {
		helpers.Sleep()
	})
`, vanityPath(path))
						}
						inTRunBlock = false
						tRunContent = ""
					}
				}
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("Error walking through files: %v", err)
	}
}
