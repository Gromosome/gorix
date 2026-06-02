// Package engine provides whole framework layer
package engine

import (
	"os"
	"path/filepath"
	"strings"
)

func ValidateProject(root string) error {
	if err := ValidateRoot(root); err != nil {
		return err
	}

	return filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			if path == root {
				return nil
			}

			if strings.Contains(path, "vendor") {
				return filepath.SkipDir
			}

			return ValidatePackageDirectory(path)
		}

		return nil
	})
}
