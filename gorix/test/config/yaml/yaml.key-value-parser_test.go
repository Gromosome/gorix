package yaml

import (
	"testing"

	yaml2 "github.com/Gromosome/gorix/gorix/config/yaml"
)

func TestSplitYAMLKeyValueIgnoresColonInsideQuotes(t *testing.T) {
	key, value, ok := yaml2.SplitYAMLKeyValue(`url: "postgres://localhost:5432/app"`)
	if !ok {
		t.Fatal("expected key-value pair")
	}
	if key != "url" || value != `"postgres://localhost:5432/app"` {
		t.Fatalf("unexpected key/value: %q %q", key, value)
	}
}

func TestSplitYAMLKeyValueDetectsBlockKey(t *testing.T) {
	key, value, ok := yaml2.SplitYAMLKeyValue("gorix:")
	if ok {
		t.Fatal("block key should not report inline value")
	}
	if key != "gorix" || value != "" {
		t.Fatalf("unexpected block key parse: %q %q", key, value)
	}
}
