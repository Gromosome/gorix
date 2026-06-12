package yaml

import (
	"os"
	"reflect"
	"testing"

	yaml2 "github.com/Gromosome/gorix/gorix/config/yaml"
)

func TestParseScalarOrInline(t *testing.T) {
	t.Setenv("GORIX_SCALAR_HOST", "localhost")

	tests := []struct {
		input string
		want  any
	}{
		{input: "null", want: nil},
		{input: "true", want: true},
		{input: "42", want: 42},
		{input: "1.25", want: 1.25},
		{input: `"quoted"`, want: "quoted"},
		{input: "${GORIX_SCALAR_HOST:-fallback}", want: "localhost"},
		{input: "${GORIX_SCALAR_MISSING:-fallback}", want: "fallback"},
	}

	for _, tt := range tests {
		if got := yaml2.ParseScalarOrInline(tt.input); !reflect.DeepEqual(got, tt.want) {
			t.Fatalf("parseScalarOrInline(%q) = %#v, want %#v", tt.input, got, tt.want)
		}
	}

	_ = os.Unsetenv("GORIX_SCALAR_MISSING")
}

func TestSplitCommaAwareIgnoresQuotedAndNestedCommas(t *testing.T) {
	input := `one, "two, still two", [three, four], {five: "six, seven"}`
	want := []string{"one", `"two, still two"`, "[three, four]", `{five: "six, seven"}`}

	if got := yaml2.SplitCommaAware(input); !reflect.DeepEqual(got, want) {
		t.Fatalf("splitCommaAware returned %#v, want %#v", got, want)
	}
}
