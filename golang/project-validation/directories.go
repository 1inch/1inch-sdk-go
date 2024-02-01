package project_validation

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const projectValidationError = `this test must be run with the working directory set to the "golang" folder of the SDK project`

func validateWorkingDirectory() error {
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %v", err)
	}

	dirName := filepath.Base(currentDir)
	if dirName != "project-validation" {
		return fmt.Errorf("%v. Current directory: %v", projectValidationError, currentDir)
	}
	return nil
}

func vanityPath(path string) string {
	return strings.Replace(path, "..", "1inch-sdk/golang", 1)
}
