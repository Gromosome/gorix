package context

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	context2 "github.com/Gromosome/gorix/gorix/core/context"
)

func TestJSONWritesStatusContentTypeAndCommits(t *testing.T) {
	recorder := httptest.NewRecorder()
	ctx := context2.NewContext(recorder, httptest.NewRequest("GET", "/", nil))

	if err := ctx.Status(context2.StatusCreated).JSON(map[string]string{"ok": "true"}); err != nil {
		t.Fatalf("JSON returned error: %v", err)
	}
	if recorder.Code != http.StatusCreated {
		t.Fatalf("unexpected status: %d", recorder.Code)
	}
	if recorder.Header().Get("Content-Type") != "application/json" {
		t.Fatalf("unexpected content type: %s", recorder.Header().Get("Content-Type"))
	}
	if !ctx.IsCommitted() {
		t.Fatal("context should be committed")
	}
}

func TestTextHTMLBlobStreamAndNoContent(t *testing.T) {
	tests := []struct {
		name        string
		write       func(*context2.Context) error
		contentType string
		body        string
	}{
		{name: "text", write: func(c *context2.Context) error { return c.Text("hello") }, contentType: "text/plain; charset=utf-8", body: "hello"},
		{name: "html", write: func(c *context2.Context) error { return c.HTML("<b>hello</b>") }, contentType: "text/html; charset=utf-8", body: "<b>hello</b>"},
		{name: "blob", write: func(c *context2.Context) error { return c.Blob("application/octet-stream", []byte("bin")) }, contentType: "application/octet-stream", body: "bin"},
		{name: "stream", write: func(c *context2.Context) error { return c.Stream("text/csv", strings.NewReader("a,b")) }, contentType: "text/csv", body: "a,b"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			ctx := context2.NewContext(recorder, httptest.NewRequest("GET", "/", nil))
			if err := tt.write(ctx); err != nil {
				t.Fatalf("write returned error: %v", err)
			}
			if recorder.Header().Get("Content-Type") != tt.contentType {
				t.Fatalf("unexpected content type: %s", recorder.Header().Get("Content-Type"))
			}
			if recorder.Body.String() != tt.body {
				t.Fatalf("unexpected body: %q", recorder.Body.String())
			}
		})
	}

	recorder := httptest.NewRecorder()
	ctx := context2.NewContext(recorder, httptest.NewRequest("GET", "/", nil))
	if err := ctx.Status(context2.StatusNoContent).NoContent(); err != nil {
		t.Fatalf("NoContent returned error: %v", err)
	}
	if recorder.Code != http.StatusNoContent {
		t.Fatalf("unexpected no content status: %d", recorder.Code)
	}
}

func TestTemplateFileDownloadAndRedirect(t *testing.T) {
	dir := t.TempDir()
	templatePath := filepath.Join(dir, "index.html")
	if err := os.WriteFile(templatePath, []byte("hello {{.Name}}"), 0o600); err != nil {
		t.Fatalf("failed to write template: %v", err)
	}

	recorder := httptest.NewRecorder()
	ctx := context2.NewContext(recorder, httptest.NewRequest("GET", "/", nil))
	if err := ctx.Template(templatePath, map[string]string{"Name": "gorix"}); err != nil {
		t.Fatalf("Template returned error: %v", err)
	}
	if recorder.Body.String() != "hello gorix" {
		t.Fatalf("unexpected template body: %q", recorder.Body.String())
	}

	filePath := filepath.Join(dir, "file.txt")
	if err := os.WriteFile(filePath, []byte("file body"), 0o600); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	recorder = httptest.NewRecorder()
	ctx = context2.NewContext(recorder, httptest.NewRequest("GET", "/", nil))
	if err := ctx.Download(filePath, "download.txt"); err != nil {
		t.Fatalf("Download returned error: %v", err)
	}
	if recorder.Header().Get("Content-Disposition") != `attachment; filename="download.txt"` {
		t.Fatalf("unexpected disposition: %s", recorder.Header().Get("Content-Disposition"))
	}

	recorder = httptest.NewRecorder()
	ctx = context2.NewContext(recorder, httptest.NewRequest("GET", "/", nil))
	if err := ctx.Status(context2.StatusFound).Redirect("/next"); err != nil {
		t.Fatalf("Redirect returned error: %v", err)
	}
	if recorder.Code != http.StatusFound {
		t.Fatalf("unexpected redirect status: %d", recorder.Code)
	}
}
