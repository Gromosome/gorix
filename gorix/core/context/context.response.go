package context

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
)

func (c *Context) GetStatusOrDefault(status StatusCode) int {
	if c == nil || c.status == 0 {
		return status.Int()
	}
	return c.status.Int()
}

func (c *Context) Status(status StatusCode) *Context {
	c.status = status
	return c
}

func (c *Context) setStatus(status StatusCode) {
	c.status = status
}

func (c *Context) Header(key, value string) *Context {
	if c != nil && c.W != nil {
		c.W.Header().Set(key, value)
	}
	return c
}

func (c *Context) JSON(data any) error {
	if c == nil || c.W == nil {
		return fmt.Errorf(
			"gorix context: response writer is unavailable",
		)
	}
	if c.err != nil {
		return c.err
	}
	if c.committed {
		return nil
	}
	c.W.Header().Set("Content-Type", "application/json")
	c.W.WriteHeader(c.GetStatusOrDefault(StatusCode(200)))
	if err := json.NewEncoder(c.W).Encode(data); err != nil {
		return err
	}
	c.committed = true
	return nil
}

func (c *Context) xml(data any) error {
	if c != nil && c.err != nil {
		return c.err
	}
	if c != nil && c.committed {
		return nil
	}
	if c == nil || c.W == nil {
		return fmt.Errorf(
			"gorix context: response writer is unavailable",
		)
	}
	c.W.Header().Set("Content-Type", "application/xml; charset=utf-8")
	c.W.WriteHeader(c.GetStatusOrDefault(StatusCode(200)))
	if err := xml.NewEncoder(c.W).Encode(data); err != nil {
		return err
	}
	c.committed = true
	c.setResponseType(ResponseTypeXML)
	return nil
}

func (c *Context) text(data string) error {
	if c != nil && c.err != nil {
		return c.err
	}
	if c != nil && c.committed {
		return nil
	}
	if c == nil || c.W == nil {
		return fmt.Errorf(
			"gorix context: response writer is unavailable",
		)
	}
	c.W.Header().Set("Content-Type", "text/plain; charset=utf-8")
	c.W.WriteHeader(c.GetStatusOrDefault(StatusCode(200)))
	_, err := c.W.Write([]byte(data))
	if err != nil {
		return err
	}
	c.committed = true
	c.setResponseType(ResponseTypeText)
	return nil
}

func (c *Context) html(data string) error {
	if c != nil && c.err != nil {
		return c.err
	}
	if c != nil && c.committed {
		return nil
	}
	if c == nil || c.W == nil {
		return fmt.Errorf(
			"gorix context: response writer is unavailable",
		)
	}
	c.W.Header().Set("Content-Type", "text/html; charset=utf-8")
	c.W.WriteHeader(c.GetStatusOrDefault(StatusCode(200)))
	_, err := c.W.Write([]byte(data))
	if err != nil {
		return err
	}
	c.committed = true
	c.setResponseType(ResponseTypeHTML)
	return nil
}
func (c *Context) template(tpl string, data any) error {
	if c != nil && c.err != nil {
		return c.err
	}
	if c != nil && c.committed {
		return nil
	}
	if c == nil || c.W == nil {
		return fmt.Errorf(
			"gorix context: response writer is unavailable",
		)
	}
	t, err := template.ParseFiles(tpl)
	if err != nil {
		return err
	}

	c.W.Header().Set("Content-Type", "text/html; charset=utf-8")
	c.W.WriteHeader(c.GetStatusOrDefault(StatusCode(200)))
	if err := t.Execute(c.W, data); err != nil {
		return err
	}
	c.committed = true
	c.setResponseType(ResponseTypeHTML)
	return nil
}

func (c *Context) blob(contentType string, data []byte) error {
	if c != nil && c.err != nil {
		return c.err
	}
	if c != nil && c.committed {
		return nil
	}
	if c == nil || c.W == nil {
		return fmt.Errorf(
			"gorix context: response writer is unavailable",
		)
	}
	c.W.Header().Set("Content-Type", contentType)
	c.W.WriteHeader(c.GetStatusOrDefault(StatusCode(200)))
	_, err := c.W.Write(data)
	if err != nil {
		return err
	}
	c.committed = true
	c.setResponseType(ResponseTypeFile)
	return nil
}

// File : local server file
func (c *Context) file(filepath string) error {
	if c != nil && c.err != nil {
		return c.err
	}
	if c != nil && c.committed {
		return nil
	}
	if c == nil || c.W == nil {
		return fmt.Errorf(
			"gorix context: response writer is unavailable",
		)
	}
	c.setResponseType(ResponseTypeDownload)
	return c.serveLocalFile(filepath, "", false)
}

// Download : local server file with force-download header
func (c *Context) download(filepath string, filename string) error {
	if c != nil && c.err != nil {
		return c.err
	}
	if c != nil && c.committed {
		return nil
	}
	if c == nil || c.W == nil {
		return fmt.Errorf(
			"gorix context: response writer is unavailable",
		)
	}
	c.setResponseType(ResponseTypeDownload)
	return c.serveLocalFile(filepath, filename, true)
}

// Stream : stream data from reader with specified content type
func (c *Context) stream(contentType string, reader io.Reader) error {
	if c != nil && c.err != nil {
		return c.err
	}
	if c != nil && c.committed {
		return nil
	}
	if c == nil || c.W == nil {
		return fmt.Errorf(
			"gorix context: response writer is unavailable",
		)
	}
	c.W.Header().Set("Content-Type", contentType)
	c.W.WriteHeader(c.GetStatusOrDefault(StatusCode(200)))
	if _, err := io.Copy(c.W, reader); err != nil {
		return err
	}
	c.committed = true
	c.setResponseType(ResponseTypeStream)
	return nil
}

// Redirect :send browser to signed S3 URL
func (c *Context) redirect(url string) error {
	if c != nil && c.err != nil {
		return c.err
	}
	if c != nil && c.committed {
		return nil
	}
	if c == nil || c.W == nil || c.R == nil {
		return fmt.Errorf(
			"gorix context: response writer or request is unavailable",
		)
	}
	http.Redirect(c.W, c.R, url, c.GetStatusOrDefault(302))
	c.committed = true
	c.setResponseType(ResponseTypeRedirect)
	return nil
}

// Utilities
func (c *Context) serveLocalFile(filepath string, filename string, download bool) error {
	if c != nil && c.committed {
		return nil
	}
	if c == nil || c.W == nil {
		return fmt.Errorf(
			"gorix context: response writer is unavailable",
		)
	}
	file, err := os.Open(filepath)
	if err != nil {
		http.Error(c.W, "file not found", http.StatusNotFound)
		return err
	}

	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println("gorix: error closing file:", err)
		}
	}()

	info, err := file.Stat()
	if err != nil {
		http.Error(c.W, "file stat error", http.StatusInternalServerError)
		return err
	}

	if download {
		c.W.Header().Set(
			"Content-Disposition",
			fmt.Sprintf(`attachment; filename="%s"`, filename),
		)
	}
	http.ServeContent(c.W, c.R, filename, info.ModTime(), file)
	c.committed = true
	return nil
}
