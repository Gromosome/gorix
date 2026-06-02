package engine

import (
	"fmt"
	"go/ast"
)

func validateMiddlewareFile(path string) error {
	file, err := parseGoFile(path)
	if err != nil {
		return err
	}

	hasMiddlewareFactory := false

	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		if fn.Recv != nil {
			continue
		}

		if fn.Type.Results == nil || len(fn.Type.Results.List) != 1 {
			continue
		}

		returnType := exprString(fn.Type.Results.List[0].Type)

		if returnType == "gorix.Middleware" {
			hasMiddlewareFactory = true
		}
	}

	if !hasMiddlewareFactory {
		return fmt.Errorf(
			"gorix validation error: middleware file %s must contain at least one function returning gorix.Middleware",
			path,
		)
	}

	return nil
}

func validateInterceptorFile(path string) error {
	file, err := parseGoFile(path)
	if err != nil {
		return err
	}

	structs := collectStructs(file)
	if len(structs) != 1 {
		return fmt.Errorf(
			"gorix validation error: interceptor file %s must have only one struct, found %d",
			path,
			len(structs),
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
			return fmt.Errorf(
				"gorix validation error: interceptor method %s in %s must use receiver of %s",
				fn.Name.Name,
				path,
				structName,
			)
		}

		switch fn.Name.Name {
		case "Before":
			if !validateBeforeAfterSignature(fn) {
				return fmt.Errorf(
					"gorix validation error: Before in %s must be Before(ctx *gorix.ExecutionContext) error",
					path,
				)
			}
			hasBefore = true

		case "After":
			if !validateBeforeAfterSignature(fn) {
				return fmt.Errorf(
					"gorix validation error: After in %s must be After(ctx *gorix.ExecutionContext) error",
					path,
				)
			}
			hasAfter = true

		default:
			return fmt.Errorf(
				"gorix validation error: interceptor file %s allows only receiver methods Before and After, found %s",
				path,
				fn.Name.Name,
			)
		}
	}

	if !hasBefore {
		return fmt.Errorf(
			"gorix validation error: interceptor file %s must have Before(ctx *gorix.ExecutionContext) error",
			path,
		)
	}

	if !hasAfter {
		return fmt.Errorf(
			"gorix validation error: interceptor file %s must have After(ctx *gorix.ExecutionContext) error",
			path,
		)
	}

	return nil
}

func validateBeforeAfterSignature(fn *ast.FuncDecl) bool {
	if fn.Type.Params == nil || len(fn.Type.Params.List) != 1 {
		return false
	}

	paramType := exprString(fn.Type.Params.List[0].Type)
	if paramType != "*gorix.ExecutionContext" {
		return false
	}

	if fn.Type.Results == nil || len(fn.Type.Results.List) != 1 {
		return false
	}

	resultType := exprString(fn.Type.Results.List[0].Type)
	if resultType != "error" {
		return false
	}

	return true
}

func validateFilterFile(path string) error {
	file, err := parseGoFile(path)
	if err != nil {
		return err
	}

	structs := collectStructs(file)
	if len(structs) != 1 {
		return fmt.Errorf(
			"gorix validation error: filter file %s must have only one struct, found %d",
			path,
			len(structs),
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
			return fmt.Errorf(
				"gorix validation error: filter method %s in %s must use receiver of %s",
				fn.Name.Name,
				path,
				structName,
			)
		}

		if fn.Name.Name != "Catch" {
			return fmt.Errorf(
				"gorix validation error: filter file %s allows only receiver method Catch, found %s",
				path,
				fn.Name.Name,
			)
		}

		if !validateCatchSignature(fn) {
			return fmt.Errorf(
				"gorix validation error: Catch in %s must be Catch(ctx *gorix.ExceptionContext)",
				path,
			)
		}

		hasCatch = true
	}

	if !hasCatch {
		return fmt.Errorf(
			"gorix validation error: filter file %s must have Catch(ctx *gorix.ExceptionContext)",
			path,
		)
	}

	return nil
}

func validateCatchSignature(fn *ast.FuncDecl) bool {
	if fn.Type.Params == nil || len(fn.Type.Params.List) != 1 {
		return false
	}
	paramType := exprString(fn.Type.Params.List[0].Type)
	if paramType != "*gorix.ExceptionContext" {
		return false
	}
	if fn.Type.Results != nil && len(fn.Type.Results.List) > 0 {
		return false
	}
	return true
}
