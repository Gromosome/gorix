package context

import (
	"io"
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
	c.setResponseType(ResponseTypeJSON)
	if c.bindingErr != nil {
		return nil, c.bindingErr
	}
	data, err := callback()
	if err != nil {
		return nil, err
	}
	if isEmpty(data) {
		c.setStatus(StatusNoContent)
		return []any{}, nil
	}
	return data, nil
}

func (c *Context) ResponseEntityXML(callback func() (any, error)) (any, error) {
	c.setResponseType(ResponseTypeXML)
	if c != nil && c.bindingErr != nil {
		return nil, c.bindingErr
	}
	data, err := callback()
	if err != nil {
		return nil, err
	}
	return data, nil
}
func (c *Context) ResponseEntitySoap11(callback func() (any, error)) (any, error) {
	c.setResponseType(ResponseTypeSOAP11)
	if c != nil && c.bindingErr != nil {
		return nil, c.bindingErr
	}
	data, err := callback()
	if err != nil {
		return nil, err
	}
	return nil, c.soap(SOAP11, c.GetStatusOrDefault(SOAPStatusOK), data)
}
func (c *Context) ResponseEntitySoap12(callback func() (any, error)) (any, error) {
	c.setResponseType(ResponseTypeSOAP12)
	if c != nil && c.bindingErr != nil {
		return nil, c.bindingErr
	}
	data, err := callback()
	if err != nil {
		return nil, err
	}
	return nil, c.soap(SOAP12, c.GetStatusOrDefault(SOAPStatusOK), data)
}

func (c *Context) ResponseEntityText(callback func() (string, error)) (any, error) {
	c.setResponseType(ResponseTypeJSON)
	if c != nil && c.bindingErr != nil {
		return nil, c.bindingErr
	}
	data, err := callback()
	if err != nil {
		return nil, err
	}
	return nil, c.text(data)
}

func (c *Context) ResponseEntityHTML(callback func() (string, error)) (any, error) {
	c.setResponseType(ResponseTypeJSON)
	if c != nil && c.bindingErr != nil {
		return nil, c.bindingErr
	}
	data, err := callback()
	if err != nil {
		return nil, err
	}
	return nil, c.html(data)
}

func (c *Context) ResponseEntityTemplate(callback func() (Template, error)) (any, error) {
	c.setResponseType(ResponseTypeJSON)
	if c != nil && c.bindingErr != nil {
		return nil, c.bindingErr
	}
	data, err := callback()
	if err != nil {
		return nil, err
	}
	return nil, c.template(data.TemplateFile, data.Data)
}

func (c *Context) ResponseEntityBlob(callback func() (Blob, error)) (any, error) {
	c.setResponseType(ResponseTypeJSON)
	if c != nil && c.bindingErr != nil {
		return nil, c.bindingErr
	}
	data, err := callback()
	if err != nil {
		return nil, err
	}
	return nil, c.blob(data.ContentType, data.Data)
}

func (c *Context) ResponseEntityFile(callback func() (File, error)) (any, error) {
	c.setResponseType(ResponseTypeJSON)
	if c != nil && c.bindingErr != nil {
		return nil, c.bindingErr
	}
	data, err := callback()
	if err != nil {
		return nil, err
	}
	return nil, c.file(data.FilePath)
}
func (c *Context) ResponseEntityDownload(callback func() (Download, error)) (any, error) {
	c.setResponseType(ResponseTypeJSON)
	if c != nil && c.bindingErr != nil {
		return nil, c.bindingErr
	}
	data, err := callback()
	if err != nil {
		return nil, err
	}
	return nil, c.download(data.FilePath, data.FileName)
}

func (c *Context) ResponseEntityStream(callback func() (Stream, error)) (any, error) {
	c.setResponseType(ResponseTypeJSON)
	if c != nil && c.bindingErr != nil {
		return nil, c.bindingErr
	}
	data, err := callback()
	if err != nil {
		return nil, err
	}
	return nil, c.stream(data.ContentType, data.Reader)
}

func (c *Context) ResponseEntityRedirect(callback func() (string, error)) (any, error) {
	c.setResponseType(ResponseTypeJSON)
	if c != nil && c.bindingErr != nil {
		return nil, c.bindingErr
	}
	data, err := callback()
	if err != nil {
		return nil, err
	}
	return nil, c.redirect(data)
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
