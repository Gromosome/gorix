package yaml

import (
	"os"
	"path/filepath"
	"testing"

	yaml2 "github.com/Gromosome/gorix/gorix/config/yaml"
)

func TestIsValidEnvironmentName(t *testing.T) {
	valid := []string{"dev", "prod_1", "staging-us"}
	for _, name := range valid {
		if !yaml2.IsValidEnvironmentName(name) {
			t.Fatalf("expected valid environment name: %s", name)
		}
	}

	invalid := []string{"", "../prod", "prod.env", "prod/dev"}
	for _, name := range invalid {
		if yaml2.IsValidEnvironmentName(name) {
			t.Fatalf("expected invalid environment name: %s", name)
		}
	}
}

func TestLoadEnvFileParsesValuesAndPreservesExistingEnvironment(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	content := []byte(`
GORIX_ENV_EXISTING=from-file
GORIX_ENV_PLAIN=value # comment
GORIX_ENV_QUOTED="hello\nworld"
GORIX_ENV_SINGLE='literal value'
GORIX_ENV_EXPANDED=${GORIX_ENV_PLAIN}/suffix
`)
	if err := os.WriteFile(path, content, 0o600); err != nil {
		t.Fatalf("failed to write env file: %v", err)
	}

	t.Setenv("GORIX_ENV_EXISTING", "from-os")
	t.Cleanup(func() {
		_ = os.Unsetenv("GORIX_ENV_PLAIN")
		_ = os.Unsetenv("GORIX_ENV_QUOTED")
		_ = os.Unsetenv("GORIX_ENV_SINGLE")
		_ = os.Unsetenv("GORIX_ENV_EXPANDED")
	})

	if err := yaml2.LoadEnvFile(path, true); err != nil {
		t.Fatalf("LoadEnvFile returned error: %v", err)
	}

	if os.Getenv("GORIX_ENV_EXISTING") != "from-os" {
		t.Fatal("LoadEnvFile overwrote existing environment variable")
	}
	if os.Getenv("GORIX_ENV_PLAIN") != "value" {
		t.Fatalf("unexpected plain env: %q", os.Getenv("GORIX_ENV_PLAIN"))
	}
	if os.Getenv("GORIX_ENV_QUOTED") != "hello\nworld" {
		t.Fatalf("unexpected quoted env: %q", os.Getenv("GORIX_ENV_QUOTED"))
	}
	if os.Getenv("GORIX_ENV_EXPANDED") != "value/suffix" {
		t.Fatalf("unexpected expanded env: %q", os.Getenv("GORIX_ENV_EXPANDED"))
	}
}

func TestLoadEnvFileRequiredMissingReturnsError(t *testing.T) {
	if err := yaml2.LoadEnvFile(filepath.Join(t.TempDir(), "missing.env"), true); err == nil {
		t.Fatal("expected missing required env error")
	}
}
