package engine

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"
)

func validateDTOFile(path string) error {
	fs, file, err := parseGoFile(path)
	if err != nil {
		return err
	}

	structs := collectStructs(file)
	if len(structs) == 0 {
		line, col := positionOf(fs, file)

		return newValidationError(
			path,
			line,
			col,
			"dto file must contain at least one DTO struct",
		)
	}

	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}

		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}

			if err := validateDTOStruct(fs, path, typeSpec.Name.Name, structType); err != nil {
				return err
			}
		}
	}

	return nil
}

func validateDTOStruct(fs *token.FileSet, path string, structName string, structType *ast.StructType) error {
	if !hasSuffix(structName, "Dto") && !hasSuffix(structName, "DTO") {
		line, col := positionOf(fs, structType)

		return newValidationError(
			path,
			line,
			col,
			fmt.Sprintf("DTO struct %s should end with Dto or DTO", structName),
		)
	}

	for _, field := range structType.Fields.List {
		if len(field.Names) == 0 {
			continue
		}

		fieldName := field.Names[0].Name

		if !field.Names[0].IsExported() {
			line, col := positionOf(fs, field)

			return newValidationError(
				path,
				line,
				col,
				fmt.Sprintf("DTO field %s must be exported", fieldName),
			)
		}

		if field.Tag == nil {
			line, col := positionOf(fs, field)

			return newValidationError(
				path,
				line,
				col,
				fmt.Sprintf("DTO field %s must have json, query, or param tag", fieldName),
			)
		}

		tag := field.Tag.Value

		if !contains(tag, "json:") && !contains(tag, "query:") && !contains(tag, "param:") {
			line, col := positionOf(fs, field.Tag)

			return newValidationError(
				path,
				line,
				col,
				fmt.Sprintf("DTO field %s must have json, query, or param tag", fieldName),
			)
		}
	}

	return nil
}
func hasSuffix(value string, suffix string) bool {
	if len(value) < len(suffix) {
		return false
	}

	return value[len(value)-len(suffix):] == suffix
}

func contains(value string, part string) bool {
	return strings.Contains(value, part)
}
