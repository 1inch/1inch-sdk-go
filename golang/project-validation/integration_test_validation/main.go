package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// This is a helper script to ensure all table-driven integration tests call the helpers.sleep() function
// Can only be run when the working directory is set to the "golang" folder of the SDK

const targetBlock = "helpers.Sleep()"

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

	err = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasSuffix(path, "_integration_test.go") || strings.HasSuffix(path, "_e2e_test.go") {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			functionName := extractFunctionName(string(content))

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
							fmt.Printf("File: %s - Missing cleanup block in the test: %s\n", path, functionName)
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
		fmt.Println("Failed to walk the path:", err)
		return
	}
	fmt.Println("Done!")
}

func extractFunctionName(content string) string {
	r := regexp.MustCompile(`func (\w+)\(t \*testing.T\)`)
	matches := r.FindStringSubmatch(content)
	if len(matches) > 1 {
		return matches[1]
	}
	return "Unknown"
}
