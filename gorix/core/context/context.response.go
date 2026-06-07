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

func (c *Context) getStatus() int {
	if c == nil || c.status == 0 {
		return StatusOK.Int()
	}
	return c.status.Int()
}

func (c *Context) Status(status StatusCode) *Context {
	c.status = status
	return c
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
	c.W.Header().Set("Content-Type", "application/json")
	c.W.WriteHeader(c.getStatus())
	return json.NewEncoder(c.W).Encode(data)
}

func (c *Context) XML(data any) error {
	if c == nil || c.W == nil {
		return fmt.Errorf(
			"gorix context: response writer is unavailable",
		)
	}
	c.W.Header().Set("Content-Type", "application/xml; charset=utf-8")
	c.W.WriteHeader(c.getStatus())
	return xml.NewEncoder(c.W).Encode(data)
}

func (c *Context) Text(data string) error {
	if c == nil || c.W == nil {
		return fmt.Errorf(
			"gorix context: response writer is unavailable",
		)
	}
	c.W.Header().Set("Content-Type", "text/plain; charset=utf-8")
	c.W.WriteHeader(c.getStatus())
	_, err := c.W.Write([]byte(data))
	return err
}

func (c *Context) HTML(data string) error {
	if c == nil || c.W == nil {
		return fmt.Errorf(
			"gorix context: response writer is unavailable",
		)
	}
	c.W.Header().Set("Content-Type", "text/html; charset=utf-8")
	c.W.WriteHeader(c.getStatus())
	_, err := c.W.Write([]byte(data))
	return err
}

func (c *Context) Template(tpl string, data any) error {
	if c == nil || c.W == nil {
		return fmt.Errorf(
			"gorix context: response writer is unavailable",
		)
	}
	c.W.Header().Set("Content-Type", "text/html; charset=utf-8")
	c.W.WriteHeader(c.getStatus())

	t, err := template.ParseFiles(tpl)
	if err != nil {
		return err
	}

	return t.Execute(c.W, data)
}

func (c *Context) Blob(contentType string, data []byte) error {
	if c == nil || c.W == nil {
		return fmt.Errorf(
			"gorix context: response writer is unavailable",
		)
	}
	c.W.Header().Set("Content-Type", contentType)
	c.W.WriteHeader(c.getStatus())
	_, err := c.W.Write(data)
	return err
}
func (c *Context) serveLocalFile(filepath string, filename string, download bool) error {
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
	return nil
}

// File : local server file
func (c *Context) File(filepath string) error {
	return c.serveLocalFile(filepath, "", false)
}

// Download : local server file with force-download header
func (c *Context) Download(filepath string, filename string) error {
	return c.serveLocalFile(filepath, filename, true)
}

// Stream : stream data from reader with specified content type
func (c *Context) Stream(contentType string, reader io.Reader) error {
	if c == nil || c.W == nil {
		return fmt.Errorf(
			"gorix context: response writer is unavailable",
		)
	}
	c.W.Header().Set("Content-Type", contentType)
	c.W.WriteHeader(c.getStatus())
	if _, err := io.Copy(c.W, reader); err != nil {
		return err
	}
	return nil
}

// Redirect :send browser to signed S3 URL
func (c *Context) Redirect(url string) error {
	if c == nil || c.W == nil || c.R == nil {
		return fmt.Errorf(
			"gorix context: response writer or request is unavailable",
		)
	}
	http.Redirect(c.W, c.R, url, c.getStatus())
	return nil
}

func (c *Context) NoContent() error {
	if c == nil || c.W == nil {
		return fmt.Errorf(
			"gorix context: response writer is unavailable",
		)
	}
	c.W.WriteHeader(c.getStatus())
	return nil
}
