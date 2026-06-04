package engine

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
			if err := validateControllerFile(fullPath); err != nil {
				return err
			}
		}
		if strings.HasSuffix(name, ".module.go") {
			if err := validateModuleFile(fullPath, folderName); err != nil {
				return err
			}
		}
		if strings.HasSuffix(name, ".middleware.go") {
			if err := validateMiddlewareFile(fullPath); err != nil {
				return err
			}
		}
		if strings.HasSuffix(name, ".interceptor.go") {
			if err := validateInterceptorFile(fullPath); err != nil {
				return err
			}
		}
		if strings.HasSuffix(name, "filter.go") {
			if err := validateFilterFile(fullPath); err != nil {
				return err
			}
		}
		if strings.HasSuffix(name, ".service.go") {
			if err := validateServiceFile(fullPath); err != nil {
				return err
			}
		}
		if strings.HasSuffix(name, ".dto.go") {
			if err := validateDTOFile(fullPath); err != nil {
				return err
			}
		}
	}
	return nil
}
