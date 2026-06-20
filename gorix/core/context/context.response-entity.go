package context

import (
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
)

type Template struct {
	TemplateFile string
	Data         any
}
type Blob struct {
	ContentType string
	Data        []byte
}

type File struct {
	FilePath string
}

type Download struct {
	FilePath string
	FileName string
}

type Stream struct {
	ContentType string
	Reader      io.Reader
}

type Redirect struct {
	URL string
}

func (c *Context) ResponseEntityJSON(callback func() (any, error)) (any, error) {
	if c == nil || c.W == nil {
		return nil, fmt.Errorf(
			"gorix context: response writer is unavailable",
		)
	}
	c.setResponseType(ResponseTypeJSON)
	if c.err != nil {
		return nil, c.err
	}
	data, err := callback()
	if isEmpty(data) {
		c.setStatus(StatusCode(http.StatusNoContent))
		return []string{}, nil
	}
	return data, err
}

func (c *Context) ResponseEntityXML(callback func() (any, error)) (any, error) {
	if c == nil || c.W == nil {
		return nil, fmt.Errorf(
			"gorix context: response writer is unavailable",
		)
	}
	c.setResponseType(ResponseTypeXML)
	if c.err != nil {
		return nil, c.err
	}
	data, err := callback()
	return err, c.xml(data)
}

func (c *Context) ResponseEntityText(data string, err error) (any, error) {
	if err != nil {
		return nil, err
	}
	return nil, c.text(data)
}

func (c *Context) ResponseEntityHTML(data string, err error) (any, error) {
	if err != nil {
		return nil, err
	}
	return nil, c.html(data)
}

func (c *Context) ResponseEntityTemplate(data Template, err error) (any, error) {
	if err != nil {
		return nil, err
	}
	return nil, c.template(data.TemplateFile, data.Data)
}

func (c *Context) ResponseEntityBlob(data Blob, err error) (any, error) {
	if err != nil {
		return nil, err
	}
	return nil, c.blob(data.ContentType, data.Data)
}

func (c *Context) ResponseEntityFile(data File, err error) (any, error) {
	if err != nil {
		return nil, err
	}
	return nil, c.file(data.FilePath)
}
func (c *Context) ResponseEntityDownload(data Download, err error) (any, error) {
	if err != nil {
		return nil, err
	}
	return nil, c.download(data.FilePath, data.FileName)
}

func (c *Context) ResponseEntityStream(data Stream, err error) (any, error) {
	if err != nil {
		return nil, err
	}
	return nil, c.stream(data.ContentType, data.Reader)
}

func (c *Context) ResponseEntityRedirect(url string, err error) (any, error) {
	if err != nil {
		return nil, err
	}
	return nil, c.redirect(url)
}

func isEmpty(value any) bool {
	if value == nil {
		return true
	}

	v := reflect.ValueOf(value)

	switch v.Kind() {
	case reflect.String:
		return strings.TrimSpace(v.String()) == ""

	case reflect.Slice, reflect.Array, reflect.Map:
		return v.Len() == 0

	case reflect.Pointer, reflect.Interface, reflect.Chan, reflect.Func:
		return v.IsNil()

	default:
		return v.IsZero()
	}
}
