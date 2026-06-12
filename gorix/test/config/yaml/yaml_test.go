package yaml

import (
	"reflect"
	"testing"

	yaml2 "github.com/Gromosome/gorix/gorix/config/yaml"
)

func TestParseYAMLReadsNestedMapsListsAndInlineValues(t *testing.T) {
	parsed := yaml2.ParseYAML(`
name: gorix
enabled: true
count: 3
ratio: 1.5
tags: [api, worker, 7]
database:
  host: localhost
  ports:
    - 5432
    - 5433
services:
  - name: api
    port: 8080
  - name: worker
    port: 9090
`)

	if yaml2.GetString(parsed, "name", "") != "gorix" {
		t.Fatalf("unexpected name: %#v", parsed["name"])
	}
	if yaml2.GetInt(parsed, "count", 0) != 3 {
		t.Fatalf("unexpected count")
	}
	if got := yaml2.GetStringSlice(parsed, "tags"); !reflect.DeepEqual(got, []string{"api", "worker", "7"}) {
		t.Fatalf("unexpected tags: %#v", got)
	}

	services, ok := parsed["services"].([]yaml2.YAMLValue)
	if !ok || len(services) != 2 {
		t.Fatalf("unexpected services: %#v", parsed["services"])
	}
}

func TestGetHelpersReturnFallbacksForMissingOrWrongTypes(t *testing.T) {
	parsed := yaml2.ParseYAML(`name: gorix`)

	if yaml2.GetString(parsed, "missing", "fallback") != "fallback" {
		t.Fatal("GetString did not return fallback")
	}
	if yaml2.GetInt(parsed, "name", 42) != 42 {
		t.Fatal("GetInt did not return fallback for string")
	}
	if yaml2.GetBoolPtr(parsed, "missing") != nil {
		t.Fatal("GetBoolPtr should return nil for missing path")
	}
	if _, ok := yaml2.GetMap(parsed, "name"); ok {
		t.Fatal("GetMap should reject non-map value")
	}
}
