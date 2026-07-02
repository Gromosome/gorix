package linter

import (
	"fmt"
	"go/ast"
	"go/token"
	"path/filepath"
)

func ValidateControllerFile(path string) error {
	fs, file, err := parseGoFile(path)
	if err != nil {
		return err
	}

	structs := collectStructs(file)
	if len(structs) != 1 {
		line, col := positionOf(fs, file)

		return newValidationError(
			path,
			line,
			col,
			fmt.Sprintf("controller file must have only one struct, found %d", len(structs)),
		)
	}

	structName := structs[0]
	constructorCount := 0
	var constructorNode ast.Node = file

	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		if fn.Recv == nil {
			if returnsType(fn, structName) || returnsPointerType(fn, structName) {
				constructorCount++
				constructorNode = fn
			}
			continue
		}

		if !receiverIsSameStruct(fn, structName) {
			line, col := positionOf(fs, fn)

			return newValidationError(
				path,
				line,
				col,
				fmt.Sprintf(
					"controller receiver method %s must use receiver of %s",
					fn.Name.Name,
					structName,
				),
			)
		}

		if err := validateControllerReceiverReturn(fs, path, fn); err != nil {
			return err
		}
	}

	if constructorCount != 1 {
		if fn := findConstructorLikeFunction(file, structName); fn != nil {
			constructorNode = fn
		}

		line, col := positionOf(fs, constructorNode)

		return newValidationError(
			path,
			line,
			col,
			fmt.Sprintf(
				"controller file must have exactly one constructor returning %s or *%s",
				structName,
				structName,
			),
		)
	}

	return nil
}

func ValidateModuleFile(path string, folderName string) error {
	fs, file, err := parseGoFile(path)
	if err != nil {
		return err
	}

	filename := filepath.Base(path)
	expectedName := folderName + ".module.go"

	if filename != expectedName {
		line, col := positionOf(fs, file)

		return newValidationError(
			path,
			line,
			col,
			fmt.Sprintf("module file %s must be named %s", filename, expectedName),
		)
	}

	if file.Name.Name != folderName {
		line, col := positionOf(fs, file.Name)

		return newValidationError(
			path,
			line,
			col,
			fmt.Sprintf("package name must be same as folder name %s", folderName),
		)
	}

	structs := collectStructs(file)
	if len(structs) != 1 {
		line, col := positionOf(fs, file)

		return newValidationError(
			path,
			line,
			col,
			fmt.Sprintf("module file must have only one struct, found %d", len(structs)),
		)
	}

	structName := structs[0]
	constructorCount := 0
	var constructorNode ast.Node = file

	hasBasePath := false
	hasControllers := false

	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		if fn.Recv == nil {
			if returnsPointerType(fn, structName) {
				constructorCount++
				constructorNode = fn
			}
			continue
		}

		if !receiverIsSameStruct(fn, structName) {
			line, col := positionOf(fs, fn)

			return newValidationError(
				path,
				line,
				col,
				fmt.Sprintf(
					"module receiver method %s must use receiver of %s",
					fn.Name.Name,
					structName,
				),
			)
		}

		switch fn.Name.Name {
		case "BasePath":
			if err := validateModuleBasePathSignature(fs, path, fn); err != nil {
				return err
			}
			hasBasePath = true

		case "APIVersion":
			if err := validateModuleAPIVersionSignature(fs, path, fn); err != nil {
				return err
			}

		case "Controllers":
			if err := validateModuleControllersSignature(fs, path, fn); err != nil {
				return err
			}
			hasControllers = true

		case "Providers":
			if err := validateModuleProvidersSignature(fs, path, fn); err != nil {
				return err
			}

		default:
			line, col := positionOf(fs, fn)

			return newValidationError(
				path,
				line,
				col,
				fmt.Sprintf(
					"module file allows only receiver methods BasePath, Providers, Controllers,APIVersion found %s",
					fn.Name.Name,
				),
			)
		}
	}

	if constructorCount != 1 {
		if fn := findConstructorLikeFunction(file, structName); fn != nil {
			constructorNode = fn
		}

		line, col := positionOf(fs, constructorNode)

		return newValidationError(
			path,
			line,
			col,
			fmt.Sprintf("module file must have exactly one constructor returning *%s", structName),
		)
	}

	if !hasBasePath {
		line, col := positionOf(fs, file)

		return newValidationError(
			path,
			line,
			col,
			"module file must have BasePath() gorix.BasePath",
		)
	}

	if !hasControllers {
		line, col := positionOf(fs, file)

		return newValidationError(
			path,
			line,
			col,
			"module file must have Controllers() []any",
		)
	}
	return nil
}
func validateModuleBasePathSignature(fs *token.FileSet, path string, fn *ast.FuncDecl) error {
	if fn.Type.Params != nil && len(fn.Type.Params.List) > 0 {
		line, col := positionOf(fs, fn)

		return newValidationError(
			path,
			line,
			col,
			"BasePath must not have parameters",
		)
	}

	if fn.Type.Results == nil || len(fn.Type.Results.List) != 1 {
		line, col := positionOf(fs, fn)

		return newValidationError(
			path,
			line,
			col,
			"BasePath must return gorix.BasePath",
		)
	}

	returnType := ExprString(fn.Type.Results.List[0].Type)
	if returnType != "gorix.BasePath" {
		line, col := positionOf(fs, fn.Type.Results.List[0].Type)

		return newValidationError(
			path,
			line,
			col,
			"BasePath must return gorix.BasePath",
		)
	}

	return nil
}
func validateModuleAPIVersionSignature(fs *token.FileSet, path string, fn *ast.FuncDecl) error {
	if fn.Type.Params != nil && len(fn.Type.Params.List) > 0 {
		line, col := positionOf(fs, fn)

		return newValidationError(
			path,
			line,
			col,
			"APIVersion must not have parameters",
		)
	}

	if fn.Type.Results == nil || len(fn.Type.Results.List) != 1 {
		line, col := positionOf(fs, fn)

		return newValidationError(
			path,
			line,
			col,
			"APIVersion must return gorix.APIVersion",
		)
	}

	returnType := ExprString(fn.Type.Results.List[0].Type)
	if returnType != "gorix.APIVersion" {
		line, col := positionOf(fs, fn.Type.Results.List[0].Type)

		return newValidationError(
			path,
			line,
			col,
			"APIVersion must return gorix.APIVersion",
		)
	}

	return nil
}

func validateModuleControllersSignature(fs *token.FileSet, path string, fn *ast.FuncDecl) error {
	if fn.Type.Params != nil && len(fn.Type.Params.List) > 0 {
		line, col := positionOf(fs, fn)

		return newValidationError(
			path,
			line,
			col,
			"Controllers must not have parameters",
		)
	}

	if fn.Type.Results == nil || len(fn.Type.Results.List) != 1 {
		line, col := positionOf(fs, fn)

		return newValidationError(
			path,
			line,
			col,
			"Controllers must return []any",
		)
	}

	returnType := ExprString(fn.Type.Results.List[0].Type)
	if returnType != "[]any" {
		line, col := positionOf(fs, fn.Type.Results.List[0].Type)

		return newValidationError(
			path,
			line,
			col,
			"Controllers must return []any",
		)
	}

	return nil
}

func validateModuleProvidersSignature(fs *token.FileSet, path string, fn *ast.FuncDecl) error {
	if fn.Type.Params != nil && len(fn.Type.Params.List) > 0 {
		line, col := positionOf(fs, fn)

		return newValidationError(
			path,
			line,
			col,
			"Providers must not have parameters",
		)
	}

	if fn.Type.Results == nil || len(fn.Type.Results.List) != 1 {
		line, col := positionOf(fs, fn)

		return newValidationError(
			path,
			line,
			col,
			"Providers must return []any",
		)
	}

	returnType := ExprString(fn.Type.Results.List[0].Type)
	if returnType != "[]any" {
		line, col := positionOf(fs, fn.Type.Results.List[0].Type)

		return newValidationError(
			path,
			line,
			col,
			"Providers must return []any",
		)
	}

	return nil
}

func ValidateServiceFile(path string) error {
	fs, file, err := parseGoFile(path)
	if err != nil {
		return err
	}

	structs := collectStructs(file)
	if len(structs) != 1 {
		line, col := positionOf(fs, file)

		return newValidationError(
			path,
			line,
			col,
			fmt.Sprintf("service file must have only one struct, found %d", len(structs)),
		)
	}

	structName := structs[0]
	constructorCount := 0
	var constructorNode ast.Node = file

	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		// Constructor function
		if fn.Recv == nil {
			if isServiceConstructor(fn, structName) {
				constructorCount++
				constructorNode = fn
			}
			continue
		}

		// Receiver methods must belong to same service struct
		if !receiverIsSameStruct(fn, structName) {
			line, col := positionOf(fs, fn)

			return newValidationError(
				path,
				line,
				col,
				fmt.Sprintf(
					"service receiver method %s must use receiver of %s",
					fn.Name.Name,
					structName,
				),
			)
		}
	}

	if constructorCount != 1 {
		if fn := findConstructorLikeFunction(file, structName); fn != nil {
			constructorNode = fn
		}

		line, col := positionOf(fs, constructorNode)

		return newValidationError(
			path,
			line,
			col,
			fmt.Sprintf(
				"service file must have exactly one constructor named New%s returning *%s",
				structName,
				structName,
			),
		)
	}

	return nil
}

func validateControllerReceiverReturn(fs *token.FileSet, path string, fn *ast.FuncDecl) error {
	if fn.Type.Params != nil && len(fn.Type.Params.List) > 0 {
		line, col := positionOf(fs, fn)

		return newValidationError(
			path,
			line,
			col,
			fmt.Sprintf("controller receiver method %s must not have parameters", fn.Name.Name),
		)
	}

	if fn.Type.Results == nil || len(fn.Type.Results.List) != 3 {
		line, col := positionOf(fs, fn)

		return newValidationError(
			path,
			line,
			col,
			fmt.Sprintf(
				"controller receiver method %s must return gorix.Method, gorix.Path, response",
				fn.Name.Name,
			),
		)
	}

	first := ExprString(fn.Type.Results.List[0].Type)
	second := ExprString(fn.Type.Results.List[1].Type)
	third := ExprString(fn.Type.Results.List[2].Type)

	if first != "gorix.Method" {
		line, col := positionOf(fs, fn.Type.Results.List[0].Type)

		return newValidationError(
			path,
			line,
			col,
			fmt.Sprintf("controller method %s first return must be gorix.Method", fn.Name.Name),
		)
	}

	if second != "gorix.Path" {
		line, col := positionOf(fs, fn.Type.Results.List[1].Type)

		return newValidationError(
			path,
			line,
			col,
			fmt.Sprintf("controller method %s second return must be gorix.Path", fn.Name.Name),
		)
	}

	if third != "gorix.RouteHandler" {
		line, col := positionOf(fs, fn.Type.Results.List[2].Type)

		return newValidationError(
			path,
			line,
			col,
			fmt.Sprintf("controller method %s third return must be gorix.RouteHandler", fn.Name.Name),
		)
	}

	return nil
}

func findConstructorLikeFunction(file *ast.File, structName string) *ast.FuncDecl {
	expectedName := "New" + structName

	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		if fn.Recv != nil {
			continue
		}

		if fn.Name.Name == expectedName {
			return fn
		}
	}

	return nil
}
func isServiceConstructor(fn *ast.FuncDecl, structName string) bool {
	if fn.Name.Name != "New"+structName {
		return false
	}

	if fn.Type.Results == nil || len(fn.Type.Results.List) != 1 {
		return false
	}

	return ExprString(fn.Type.Results.List[0].Type) == "*"+structName
}
