package linter

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

func parseGoFile(path string) (*token.FileSet, *ast.File, error) {
	fs := token.NewFileSet()

	file, err := parser.ParseFile(fs, path, nil, parser.ParseComments)
	if err != nil {
		return nil, nil, newValidationError(
			path,
			1,
			1,
			fmt.Sprintf("failed to parse file: %v", err),
		)
	}

	return fs, file, nil
}
func positionOf(fs *token.FileSet, node ast.Node) (int, int) {
	if node == nil {
		return 1, 1
	}

	pos := fs.Position(node.Pos())
	return pos.Line, pos.Column
}

func collectStructs(file *ast.File) []string {
	var structs []string

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

			_, ok = typeSpec.Type.(*ast.StructType)
			if ok {
				structs = append(structs, typeSpec.Name.Name)
			}
		}
	}

	return structs
}

func receiverIsSameStruct(fn *ast.FuncDecl, structName string) bool {
	if fn.Recv == nil || len(fn.Recv.List) == 0 {
		return false
	}

	recvType := exprString(fn.Recv.List[0].Type)

	return recvType == structName || recvType == "*"+structName
}

func returnsType(fn *ast.FuncDecl, typeName string) bool {
	if fn.Type.Results == nil {
		return false
	}

	for _, result := range fn.Type.Results.List {
		if exprString(result.Type) == typeName {
			return true
		}
	}

	return false
}

func returnsPointerType(fn *ast.FuncDecl, typeName string) bool {
	if fn.Type.Results == nil {
		return false
	}

	for _, result := range fn.Type.Results.List {
		if exprString(result.Type) == "*"+typeName {
			return true
		}
	}

	return false
}

func exprString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name

	case *ast.StarExpr:
		return "*" + exprString(t.X)

	case *ast.SelectorExpr:
		return exprString(t.X) + "." + t.Sel.Name

	case *ast.ArrayType:
		return "[]" + exprString(t.Elt)

	case *ast.MapType:
		return "map[" + exprString(t.Key) + "]" + exprString(t.Value)

	case *ast.InterfaceType:
		return "interface{}"

	case *ast.FuncType:
		return "func"

	default:
		return fmt.Sprintf("%T", expr)
	}
}
