package linter

import (
	"os"
	"path/filepath"
	"strings"
)

func ValidatePackageDirectory(dir string) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	folderName := filepath.Base(dir)
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		name := file.Name()
		if !strings.HasSuffix(name, ".go") {
			continue
		}
		fullPath := filepath.Join(dir, name)
		if strings.HasSuffix(name, ".controller.go") {
			if err := ValidateControllerFile(fullPath); err != nil {
				return err
			}
		}
		if strings.HasSuffix(name, ".module.go") {
			if err := ValidateModuleFile(fullPath, folderName); err != nil {
				return err
			}
		}
		if strings.HasSuffix(name, ".middleware.go") {
			if err := ValidateMiddlewareFile(fullPath); err != nil {
				return err
			}
		}
		if strings.HasSuffix(name, ".interceptor.go") {
			if err := ValidateInterceptorFile(fullPath); err != nil {
				return err
			}
		}
		if strings.HasSuffix(name, "filter.go") {
			if err := ValidateFilterFile(fullPath); err != nil {
				return err
			}
		}
		if strings.HasSuffix(name, ".service.go") {
			if err := ValidateServiceFile(fullPath); err != nil {
				return err
			}
		}
		if strings.HasSuffix(name, ".global.go") {
			if err := ValidateDTOFile(fullPath); err != nil {
				return err
			}
		}
	}
	return nil
}
