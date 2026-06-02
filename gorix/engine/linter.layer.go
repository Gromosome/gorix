package engine

import (
	"fmt"
	"go/ast"
	"path/filepath"
)

func validateControllerFile(path string) error {
	file, err := parseGoFile(path)
	if err != nil {
		return err
	}

	structs := collectStructs(file)
	if len(structs) != 1 {
		return fmt.Errorf(
			"gorix validation error: controller file %s must have only one struct, found %d",
			path,
			len(structs),
		)
	}
	structName := structs[0]
	constructorCount := 0
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		if fn.Recv == nil {
			if returnsType(fn, structName) {
				constructorCount++
			}
			continue
		}
		if !receiverIsSameStruct(fn, structName) {
			return fmt.Errorf(
				"gorix validation error: receiver method %s in %s must use receiver of %s",
				fn.Name.Name,
				path,
				structName,
			)
		}

		if err := validateControllerReceiverReturn(fn, path); err != nil {
			return err
		}
	}

	if constructorCount != 1 {
		return fmt.Errorf(
			"gorix validation error: controller file %s must have exactly one constructor returning %s",
			path,
			structName,
		)
	}

	return nil
}

func validateModuleFile(path string, folderName string) error {
	file, err := parseGoFile(path)
	if err != nil {
		return err
	}
	filename := filepath.Base(path)
	expectedPrefix := folderName + ".module.go"
	if filename != expectedPrefix {
		return fmt.Errorf(
			"gorix validation error: module file %s must be named %s",
			filename,
			expectedPrefix,
		)
	}
	if file.Name.Name != folderName {
		return fmt.Errorf(
			"gorix validation error: package name in %s must be same as folder name %s",
			path,
			folderName,
		)
	}
	structs := collectStructs(file)
	if len(structs) != 1 {
		return fmt.Errorf(
			"gorix validation error: module file %s must have only one struct, found %d",
			path,
			len(structs),
		)
	}
	structName := structs[0]
	constructorCount := 0
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		if fn.Recv == nil {
			if returnsPointerType(fn, structName) {
				constructorCount++
			}
			continue
		}
		if !receiverIsSameStruct(fn, structName) {
			return fmt.Errorf(
				"gorix validation error: receiver method %s in %s must use receiver of %s",
				fn.Name.Name,
				path,
				structName,
			)
		}
		if err := validateModuleReceiverReturn(fn, path); err != nil {
			return err
		}
	}
	if constructorCount != 1 {
		return fmt.Errorf(
			"gorix validation error: module file %s must have exactly one constructor returning *%s",
			path,
			structName,
		)
	}
	return nil
}

func validateControllerReceiverReturn(fn *ast.FuncDecl, path string) error {
	if fn.Type.Results == nil || len(fn.Type.Results.List) != 3 {
		return fmt.Errorf(
			"gorix validation error: controller receiver method %s in %s must return gorix.Method, gorix.Path, response",
			fn.Name.Name,
			path,
		)
	}
	first := exprString(fn.Type.Results.List[0].Type)
	second := exprString(fn.Type.Results.List[1].Type)
	if first != "gorix.Method" {
		return fmt.Errorf(
			"gorix validation error: controller method %s first return must be gorix.Method",
			fn.Name.Name,
		)
	}
	if second != "gorix.Path" {
		return fmt.Errorf(
			"gorix validation error: controller method %s second return must be gorix.Path",
			fn.Name.Name,
		)
	}
	return nil
}

func validateModuleReceiverReturn(fn *ast.FuncDecl, path string) error {
	if fn.Type.Results == nil || len(fn.Type.Results.List) != 2 {
		return fmt.Errorf(
			"gorix validation error: module receiver method %s in %s must return gorix.BasePath, Controller",
			fn.Name.Name,
			path,
		)
	}
	first := exprString(fn.Type.Results.List[0].Type)
	if first != "gorix.BasePath" {
		return fmt.Errorf(
			"gorix validation error: module method %s first return must be gorix.BasePath",
			fn.Name.Name,
		)
	}
	return nil
}
