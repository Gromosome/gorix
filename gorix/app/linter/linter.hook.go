package linter

import (
	"fmt"
	"go/ast"
	"go/token"
)

func validateMiddlewareFile(path string) error {
	fs, file, err := parseGoFile(path)
	if err != nil {
		return err
	}

	hasMiddlewareFactory := false
	var factoryNode ast.Node = file

	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		if fn.Recv != nil {
			line, col := positionOf(fs, fn)

			return newValidationError(
				path,
				line,
				col,
				"middleware file must not contain receiver methods",
			)
		}

		if fn.Type.Results == nil || len(fn.Type.Results.List) != 1 {
			continue
		}

		returnType := exprString(fn.Type.Results.List[0].Type)
		if returnType == "gorix.Middleware" {
			hasMiddlewareFactory = true
			factoryNode = fn
		}
	}

	if !hasMiddlewareFactory {
		line, col := positionOf(fs, factoryNode)

		return newValidationError(
			path,
			line,
			col,
			"middleware file must contain at least one function returning gorix.Middleware",
		)
	}

	return nil
}

func validateInterceptorFile(path string) error {
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
			fmt.Sprintf("interceptor file must have only one struct, found %d", len(structs)),
		)
	}

	structName := structs[0]
	hasBefore := false
	hasAfter := false

	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		if fn.Recv == nil {
			continue
		}

		if !receiverIsSameStruct(fn, structName) {
			line, col := positionOf(fs, fn)

			return newValidationError(
				path,
				line,
				col,
				fmt.Sprintf(
					"interceptor method %s must use receiver of %s",
					fn.Name.Name,
					structName,
				),
			)
		}

		switch fn.Name.Name {
		case "Before":
			if err := validateBeforeAfterSignature(fs, path, fn); err != nil {
				return err
			}
			hasBefore = true

		case "After":
			if err := validateBeforeAfterSignature(fs, path, fn); err != nil {
				return err
			}
			hasAfter = true

		default:
			line, col := positionOf(fs, fn)

			return newValidationError(
				path,
				line,
				col,
				fmt.Sprintf(
					"interceptor file allows only receiver methods Before and After, found %s",
					fn.Name.Name,
				),
			)
		}
	}

	if !hasBefore {
		line, col := positionOf(fs, file)

		return newValidationError(
			path,
			line,
			col,
			"interceptor file must have Before(ctx *gorix.ExecutionContext) error",
		)
	}

	if !hasAfter {
		line, col := positionOf(fs, file)

		return newValidationError(
			path,
			line,
			col,
			"interceptor file must have After(ctx *gorix.ExecutionContext) error",
		)
	}

	return nil
}
func validateBeforeAfterSignature(fs *token.FileSet, path string, fn *ast.FuncDecl) error {
	if fn.Type.Params == nil || len(fn.Type.Params.List) != 1 {
		line, col := positionOf(fs, fn)

		return newValidationError(
			path,
			line,
			col,
			fmt.Sprintf("%s must have exactly one parameter: ctx *gorix.ExecutionContext", fn.Name.Name),
		)
	}

	paramType := exprString(fn.Type.Params.List[0].Type)
	if paramType != "*gorix.ExecutionContext" {
		line, col := positionOf(fs, fn.Type.Params.List[0].Type)

		return newValidationError(
			path,
			line,
			col,
			fmt.Sprintf("%s parameter must be *gorix.ExecutionContext", fn.Name.Name),
		)
	}

	if fn.Type.Results == nil || len(fn.Type.Results.List) != 1 {
		line, col := positionOf(fs, fn)

		return newValidationError(
			path,
			line,
			col,
			fmt.Sprintf("%s must return error", fn.Name.Name),
		)
	}

	resultType := exprString(fn.Type.Results.List[0].Type)
	if resultType != "error" {
		line, col := positionOf(fs, fn.Type.Results.List[0].Type)

		return newValidationError(
			path,
			line,
			col,
			fmt.Sprintf("%s must return error", fn.Name.Name),
		)
	}

	return nil
}

func validateFilterFile(path string) error {
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
			fmt.Sprintf("filter file must have only one struct, found %d", len(structs)),
		)
	}

	structName := structs[0]
	hasCatch := false

	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		if fn.Recv == nil {
			continue
		}

		if !receiverIsSameStruct(fn, structName) {
			line, col := positionOf(fs, fn)

			return newValidationError(
				path,
				line,
				col,
				fmt.Sprintf(
					"filter method %s must use receiver of %s",
					fn.Name.Name,
					structName,
				),
			)
		}

		if fn.Name.Name != "Catch" {
			line, col := positionOf(fs, fn)

			return newValidationError(
				path,
				line,
				col,
				fmt.Sprintf(
					"filter file allows only receiver method Catch, found %s",
					fn.Name.Name,
				),
			)
		}

		if err := validateCatchSignature(fs, path, fn); err != nil {
			return err
		}

		hasCatch = true
	}

	if !hasCatch {
		line, col := positionOf(fs, file)

		return newValidationError(
			path,
			line,
			col,
			"filter file must have Catch(ctx *gorix.ExceptionContext)",
		)
	}

	return nil
}

func validateCatchSignature(fs *token.FileSet, path string, fn *ast.FuncDecl) error {
	if fn.Type.Params == nil || len(fn.Type.Params.List) != 1 {
		line, col := positionOf(fs, fn)

		return newValidationError(
			path,
			line,
			col,
			"Catch must have exactly one parameter: ctx *gorix.ExceptionContext",
		)
	}

	paramType := exprString(fn.Type.Params.List[0].Type)
	if paramType != "*gorix.ExceptionContext" {
		line, col := positionOf(fs, fn.Type.Params.List[0].Type)

		return newValidationError(
			path,
			line,
			col,
			"Catch parameter must be *gorix.ExceptionContext",
		)
	}

	if fn.Type.Results != nil && len(fn.Type.Results.List) > 0 {
		line, col := positionOf(fs, fn.Type.Results.List[0].Type)

		return newValidationError(
			path,
			line,
			col,
			"Catch must not return any value",
		)
	}

	return nil
}
