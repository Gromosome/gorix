package context

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
)

func (c *Context) StreamBody(
	handler func(reader io.Reader) error,
) *Context {
	if c == nil {
		return c
	}

	if c.R == nil || c.R.Body == nil {
		c.setBindingError(NewValidationError([]FieldError{
			NewFieldError("body", "stream", "request body cannot be nil"),
		}))
		return c
	}

	if handler == nil {
		c.setBindingError(NewValidationError([]FieldError{
			NewFieldError("body", "stream", "stream handler cannot be nil"),
		}))
		return c
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(c.R.Body)

	if err := handler(c.R.Body); err != nil {
		c.setBindingError(NewValidationError([]FieldError{
			NewFieldError(
				"body",
				"stream",
				fmt.Sprintf("failed to stream request body: %v", err),
			),
		}))
		return c
	}

	return c
}

func (c *Context) IsOctetStream() bool {
	if c == nil || c.R == nil {
		return false
	}

	contentType := c.R.Header.Get("Content-Type")
	contentType = strings.ToLower(strings.TrimSpace(strings.Split(contentType, ";")[0]))

	return contentType == "application/octet-stream"
}

type FileStream struct {
	FieldName   string
	FileName    string
	ContentType string
	Header      map[string][]string
	Reader      io.Reader
	Part        *multipart.Part
}

func (c *Context) StreamFile(
	fieldName string,
	handler func(file FileStream) error,
) *Context {
	if c == nil {
		return c
	}

	if c.R == nil || c.R.Body == nil {
		c.setBindingError(NewValidationError([]FieldError{
			NewFieldError("file", "stream", "request body cannot be nil"),
		}))
		return c
	}

	if handler == nil {
		c.setBindingError(NewValidationError([]FieldError{
			NewFieldError("file", "stream", "file stream handler cannot be nil"),
		}))
		return c
	}

	contentType := c.R.Header.Get("Content-Type")
	contentType = strings.ToLower(strings.TrimSpace(strings.Split(contentType, ";")[0]))

	if contentType != "multipart/form-data" {
		c.setBindingError(NewValidationError([]FieldError{
			NewFieldError(
				"headers.Content-Type",
				"content_type",
				"Content-Type must be multipart/form-data",
			),
		}))
		return c
	}

	reader, err := c.R.MultipartReader()
	if err != nil {
		c.setBindingError(NewValidationError([]FieldError{
			NewFieldError(
				"file",
				"multipart",
				fmt.Sprintf("invalid multipart body: %v", err),
			),
		}))
		return c
	}

	found := false

	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}

		if err != nil {
			c.setBindingError(NewValidationError([]FieldError{
				NewFieldError(
					"file",
					"multipart",
					fmt.Sprintf("failed to read multipart part: %v", err),
				),
			}))
			return c
		}

		if part.FileName() == "" {
			_ = part.Close()
			continue
		}

		if fieldName != "" && part.FormName() != fieldName {
			_ = part.Close()
			continue
		}

		found = true

		file := FileStream{
			FieldName:   part.FormName(),
			FileName:    part.FileName(),
			ContentType: part.Header.Get("Content-Type"),
			Header:      part.Header,
			Reader:      part,
			Part:        part,
		}

		if err := handler(file); err != nil {
			_ = part.Close()

			c.setBindingError(NewValidationError([]FieldError{
				NewFieldError(
					"file",
					"stream",
					fmt.Sprintf("failed to handle uploaded file: %v", err),
				),
			}))
			return c
		}

		_ = part.Close()
	}

	if !found {
		c.setBindingError(NewValidationError([]FieldError{
			NewFieldError(
				"file",
				"required",
				fmt.Sprintf("file field %q is required", fieldName),
			),
		}))
		return c
	}

	return c
}

func (c *Context) LimitBody(
	maxBytes int64,
) *Context {
	if c == nil {
		return c
	}

	if c.R == nil || c.W == nil {
		c.setBindingError(NewValidationError([]FieldError{
			NewFieldError("body", "limit", "request or response writer cannot be nil"),
		}))
		return c
	}

	c.R.Body = http.MaxBytesReader(c.W, c.R.Body, maxBytes)
	return c
}
