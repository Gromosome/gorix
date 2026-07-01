package context

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"reflect"

	"github.com/Gromosome/gorix/gorix/internal/access"
	"github.com/Gromosome/gorix/gorix/internal/global"
	"github.com/Gromosome/gorix/gorix/logger"
)

func (c *Context) responseBody(data any) error {
	if c.status >= 400 {
		dto, valid := asErrorDTO(data)
		if valid {
			switch c.ResponseType() {
			case ResponseTypeJSON:
				return c.JSONFault(data)
			case ResponseTypeXML:
				return c.XMLFault(data)
			case ResponseTypeSOAP11:
				return c.soapFault(SOAP11, SOAP11FaultServer, dto.Message, dto.Error)
			case ResponseTypeSOAP12:
				return c.soapFault(SOAP12, SOAP12FaultReceiver, dto.Message, dto.Error)
			default:
				return fmt.Errorf("invalid response type: %s", c.ResponseType())
			}
		}
	}
	switch c.ResponseType() {
	case ResponseTypeJSON:
		return c.json(c.GetStatusOrDefault(StatusOK), data)
	case ResponseTypeXML:
		return c.xml(c.GetStatusOrDefault(StatusOK), data)
	case ResponseTypeSOAP11:
		return c.soap(SOAP11, c.GetStatusOrDefault(SOAPStatusOK), data)
	case ResponseTypeSOAP12:
		return c.soap(SOAP12, c.GetStatusOrDefault(SOAPStatusOK), data)
	default:
		return fmt.Errorf("invalid response type: %s", c.ResponseType())
	}
}
func asErrorDTO(data any) (*global.ErrorDTO, bool) {
	switch v := data.(type) {
	case global.ErrorDTO:
		return &v, true

	case *global.ErrorDTO:
		if v == nil {
			return nil, false
		}
		return v, true

	default:
		return nil, false
	}
}
func (c *Context) ResponseBodyInternal(_ access.Key, data any) error {
	return c.responseBody(data)
}
func (c *Context) SOAPFault11(code SOAP11FaultCode, message string, detail any) error {
	return c.soapFault(SOAP11, code, message, detail)
}
func (c *Context) SOAPFault12(code SOAP12FaultCode, message string, detail any) error {
	return c.soapFault(SOAP12, code, message, detail)
}
func (c *Context) JSONFault(data any) error {
	return c.json(c.GetStatusOrDefault(StatusBadRequest), data)
}
func (c *Context) XMLFault(data any) error {
	return c.xml(c.GetStatusOrDefault(StatusBadRequest), data)
}
func validateBody(data any) error {
	if data == nil {
		return fmt.Errorf("response body must be a struct or pointer to struct, got <nil>")
	}
	v := reflect.ValueOf(data)
	t := v.Type()
	for t.Kind() == reflect.Pointer {
		if v.IsNil() {
			return fmt.Errorf("response body must be a non-nil struct pointer, got nil %s", t.String())
		}

		v = v.Elem()
		t = v.Type()
	}
	if t.Kind() != reflect.Struct {
		return fmt.Errorf("response body must be a struct or pointer to struct, got %s", t.String())
	}
	return nil
}

func (c *Context) json(status int, data any) error {
	if c == nil || c.W == nil {
		return fmt.Errorf(
			"gorix context: response writer is unavailable",
		)
	}
	if c.committed {
		return nil
	}
	c.W.Header().Set("Content-Type", "application/json")
	c.W.WriteHeader(status)
	if status == StatusNoContent.Int() {
		c.committed = true
		return nil
	}
	if err := json.NewEncoder(c.W).Encode(data); err != nil {
		return err
	}
	c.committed = true
	return nil
}

func (c *Context) xml(status int, data any) error {
	processedData := data
	if err := validateBody(data); err != nil {
		c.logger.TypeWarn().CallerLevel(logger.Wrap4).Log("Default XML from framework, Customize response by DTO into interceptor")
		processedData = global.ResponseXML{Data: data}
	}
	if c == nil || c.W == nil {
		return fmt.Errorf(
			"gorix context: response writer is unavailable",
		)
	}
	if c.committed {
		return nil
	}
	c.W.Header().Set("Content-Type", "application/xml; charset=utf-8")
	c.W.WriteHeader(status)
	if err := xml.NewEncoder(c.W).Encode(processedData); err != nil {
		return err
	}
	c.committed = true
	return nil
}

func (c *Context) soapFault(version SOAPVersion, code SOAPFaultCode, message string, detail any) error {
	if !code.IsValid(SOAPStatusCode(c.GetStatusOrDefault(SOAPStatusBadRequest))) {
		c.logger.TypeError().CallerLevel(logger.Wrap2).Log("both Statuscode( %d ) and Faultcode( %s ) are mismatched ", SOAPStatusCode(c.GetStatusOrDefault(SOAPStatusBadRequest)))
	}
	switch version {
	case SOAP11:
		faultcode, _ := code.(SOAP11FaultCode)
		data := NewSOAP11FaultEnvelope(faultcode, message, detail)
		return c.soap(version, c.GetStatusOrDefault(SOAPStatusBadRequest), data)
	case SOAP12:
		faultcode, _ := code.(SOAP12FaultCode)
		data := NewSOAP12FaultEnvelope(faultcode, message, detail)
		return c.soap(version, c.GetStatusOrDefault(SOAPStatusBadRequest), data)
	default:
		return fmt.Errorf("invalid version: %s", version)
	}
}

func (c *Context) soap(version SOAPVersion, status int, data any) error {
	var contentType SOAPContentType
	var payload any
	switch version {
	case SOAP11:
		contentType = ContentTypeSOAP11
		payload = NewSOAP11Envelope[any](data)
	case SOAP12:
		contentType = ContentTypeSOAP12
		payload = NewSOAP12Envelope[any](data)
	}
	if c == nil || c.W == nil {
		return fmt.Errorf(
			"gorix context: response writer is unavailable",
		)
	}
	if c.committed {
		return nil
	}

	c.W.Header().Set("Content-Type", contentType.str())
	c.W.WriteHeader(status)
	if err := xml.NewEncoder(c.W).Encode(payload); err != nil {
		return err
	}
	c.committed = true
	return nil
}

func (c *Context) text(data string) error {
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
	return nil
}

func (c *Context) html(data string) error {
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
	return nil
}
func (c *Context) template(tpl string, data any) error {
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
	return nil
}

func (c *Context) blob(contentType string, data []byte) error {
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
	return nil
}

// File : local server file
func (c *Context) file(filepath string) error {
	if c != nil && c.committed {
		return nil
	}
	if c == nil || c.W == nil {
		return fmt.Errorf(
			"gorix context: response writer is unavailable",
		)
	}
	return c.serveLocalFile(filepath, "", false)
}

// Download : local server file with force-download header
func (c *Context) download(filepath string, filename string) error {
	if c != nil && c.committed {
		return nil
	}
	if c == nil || c.W == nil {
		return fmt.Errorf(
			"gorix context: response writer is unavailable",
		)
	}
	return c.serveLocalFile(filepath, filename, true)
}

// Stream : stream data from reader with specified content type
func (c *Context) stream(contentType string, reader io.Reader) error {
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
	return nil
}

// Redirect :send browser to signed S3 URL
func (c *Context) redirect(url string) error {
	if c != nil && c.committed {
		return nil
	}
	if c == nil || c.W == nil || c.R == nil {
		return fmt.Errorf(
			"gorix context: response writer or request is unavailable",
		)
	}
	http.Redirect(c.W, c.R, url, c.GetStatusOrDefault(StatusFound))
	c.committed = true
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
