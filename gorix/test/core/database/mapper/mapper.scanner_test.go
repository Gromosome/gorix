package mapper

import (
	"fmt"
	"testing"

	mapper2 "github.com/Gromosome/gorix/gorix/core/database/mapper"
)

type fakeRowScanner struct {
	columns []string
	values  []any
}

func (r fakeRowScanner) Columns() ([]string, error) {
	return r.columns, nil
}

func (r fakeRowScanner) Scan(destinations ...any) error {
	for index, destination := range destinations {
		if index >= len(r.values) {
			continue
		}
		switch pointer := destination.(type) {
		case *int:
			*pointer = r.values[index].(int)
		case *string:
			*pointer = r.values[index].(string)
		case *any:
			*pointer = r.values[index]
		default:
			return fmt.Errorf("unsupported destination %T", destination)
		}
	}
	return nil
}

type EmbeddedScanTarget struct {
	Label string `db:"label"`
}

type scanTarget struct {
	ID      int `db:"id"`
	Name    string
	Ignored string `db:"-"`
	*EmbeddedScanTarget
}

func TestScanIntoMapsColumnsToStructFields(t *testing.T) {
	row := fakeRowScanner{
		columns: []string{"id", "name", "label", "extra"},
		values:  []any{10, "bob", "admin", "ignored"},
	}
	var target scanTarget

	if err := mapper2.ScanInto(row, &target); err != nil {
		t.Fatalf("ScanInto returned error: %v", err)
	}
	if target.ID != 10 || target.Name != "bob" || target.EmbeddedScanTarget == nil || target.EmbeddedScanTarget.Label != "admin" {
		t.Fatalf("unexpected scan target: %#v", target)
	}
}

func TestScanIntoRejectsInvalidTargets(t *testing.T) {
	row := fakeRowScanner{}
	if err := mapper2.ScanInto(row, nil); err == nil {
		t.Fatal("nil target should be rejected")
	}
	var value int
	if err := mapper2.ScanInto(row, &value); err == nil {
		t.Fatal("non-struct target should be rejected")
	}
}
