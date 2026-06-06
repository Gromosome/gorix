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

type Context struct {
	W      http.ResponseWriter
	R      *http.Request
	status StatusCode
	params map[string]string
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		W:      w,
		R:      r,
		params: make(map[string]string),
	}
}

func (c *Context) getStatus() int {
	if c.status == 0 {
		return StatusOK.Int()
	}
	return c.status.Int()
}

func (c *Context) Status(status StatusCode) *Context {
	c.status = status
	return c
}

func (c *Context) Header(key, value string) *Context {
	c.W.Header().Set(key, value)
	return c
}

func (c *Context) JSON(data any) error {
	c.W.Header().Set("Content-Type", "application/json")
	c.W.WriteHeader(c.getStatus())
	return json.NewEncoder(c.W).Encode(data)
}

func (c *Context) XML(data any) error {
	c.W.Header().Set("Content-Type", "application/xml; charset=utf-8")
	c.W.WriteHeader(c.getStatus())
	return xml.NewEncoder(c.W).Encode(data)
}

func (c *Context) Text(data string) error {
	c.W.Header().Set("Content-Type", "text/plain; charset=utf-8")
	c.W.WriteHeader(c.getStatus())
	_, err := c.W.Write([]byte(data))
	return err
}

func (c *Context) HTML(data string) error {
	c.W.Header().Set("Content-Type", "text/html; charset=utf-8")
	c.W.WriteHeader(c.getStatus())
	_, err := c.W.Write([]byte(data))
	return err
}

func (c *Context) Template(tpl string, data any) error {
	c.W.Header().Set("Content-Type", "text/html; charset=utf-8")
	c.W.WriteHeader(c.getStatus())

	t, err := template.ParseFiles(tpl)
	if err != nil {
		return err
	}

	return t.Execute(c.W, data)
}

func (c *Context) Blob(contentType string, data []byte) error {
	c.W.Header().Set("Content-Type", contentType)
	c.W.WriteHeader(c.getStatus())
	_, err := c.W.Write(data)
	return err
}
func (c *Context) serveLocalFile(filepath string, filename string, download bool) error {
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
	c.W.Header().Set("Content-Type", contentType)
	c.W.WriteHeader(c.getStatus())
	if _, err := io.Copy(c.W, reader); err != nil {
		return err
	}
	return nil
}

// Redirect :send browser to signed S3 URL
func (c *Context) Redirect(url string) {
	http.Redirect(c.W, c.R, url, c.getStatus())
}

func (c *Context) NoContent() {
	c.W.WriteHeader(c.getStatus())
}
