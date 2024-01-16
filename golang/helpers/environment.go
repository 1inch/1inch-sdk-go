package helpers

import "os"

// TODO refactor this and all uses of os.Env to use this function and handle a simple error when missing

func GetenvSafe(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic("environment variable " + key + " is missing. Aborting...")
	}
	return value
}
