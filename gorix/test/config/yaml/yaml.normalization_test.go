package yaml

import (
	"testing"

	yaml2 "github.com/Gromosome/gorix/gorix/config/yaml"
)

func TestRemoveYAMLCommentPreservesQuotedHash(t *testing.T) {
	line := `value: "a # b" # comment`
	if got := yaml2.RemoveYAMLComment(line); got != `value: "a # b"` {
		t.Fatalf("removeYAMLComment returned %q", got)
	}
}

func TestNormalizeYAMLLinesHandlesCRLFAndComments(t *testing.T) {
	lines := yaml2.NormalizeYAMLLines("name: gorix # comment\r\n\r\nhost: localhost")
	if len(lines) != 3 {
		t.Fatalf("unexpected line count: %d", len(lines))
	}
	if lines[0] != "name: gorix" {
		t.Fatalf("unexpected first line: %q", lines[0])
	}
}
