// Package validate provides validation layer of gorix framework
package engine

import (
	"fmt"
	"os"
	"strings"
)

func ValidateRoot(root string) error {
	entries, err := os.ReadDir(root)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		// Only validate Go files in root layer
		if strings.HasSuffix(name, ".go") && name != "main.go" {
			return fmt.Errorf(
				"gorix validation error: root layer allows only main.go as a Go file, found %s",
				name,
			)
		}
	}
	return nil
}
