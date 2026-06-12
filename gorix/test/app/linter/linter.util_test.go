package linter

import (
	"go/ast"
	"testing"

	linter2 "github.com/Gromosome/gorix/gorix/app/linter"
)

func TestExprStringFormatsSupportedExpressions(t *testing.T) {
	expr := &ast.SelectorExpr{
		X:   ast.NewIdent("gorix"),
		Sel: ast.NewIdent("Method"),
	}
	if got := linter2.ExprString(expr); got != "gorix.Method" {
		t.Fatalf("unexpected selector expression: %s", got)
	}

	array := &ast.ArrayType{Elt: ast.NewIdent("any")}
	if got := linter2.ExprString(array); got != "[]any" {
		t.Fatalf("unexpected array expression: %s", got)
	}
}
