package context

import "encoding/xml"

type SOAPVersion string

const (
	SOAP11 SOAPVersion = "1.1"
	SOAP12 SOAPVersion = "1.2"
)

const (
	SOAP11Namespace = "http://schemas.xmlsoap.org/soap/envelope/"
	SOAP12Namespace = "http://www.w3.org/2003/05/soap-envelope"
)

type SOAPContentType string

const (
	ContentTypeSOAP11 SOAPContentType = `text/xml; charset="utf-8"`
	ContentTypeSOAP12 SOAPContentType = `application/soap+xml; charset="utf-8"`
)

func (t SOAPContentType) str() string {
	return string(t)
}

type SOAP11Envelope[T any] struct {
	XMLName xml.Name      `xml:"soapenv:Envelope"`
	SoapEnv string        `xml:"xmlns:soapenv,attr"`
	Header  *SOAP11Header `xml:"soapenv:Header,omitempty"`
	Body    SOAP11Body[T] `xml:"soapenv:Body"`
}

type SOAP11Header struct {
	Value any `xml:",omitempty"`
}

type SOAP11Body[T any] struct {
	Content T `xml:",any"`
}

type SOAP12Envelope[T any] struct {
	XMLName xml.Name      `xml:"env:Envelope"`
	Env     string        `xml:"xmlns:env,attr"`
	Header  *SOAP12Header `xml:"env:Header,omitempty"`
	Body    SOAP12Body[T] `xml:"env:Body"`
}

type SOAP12Header struct {
	Value any `xml:",omitempty"`
}

type SOAP12Body[T any] struct {
	Content T `xml:",any"`
}

func NewSOAP11Envelope[T any](data T) SOAP11Envelope[T] {
	return SOAP11Envelope[T]{
		SoapEnv: SOAP11Namespace,
		Body: SOAP11Body[T]{
			Content: data,
		},
	}
}

func NewSOAP12Envelope[T any](data T) SOAP12Envelope[T] {
	return SOAP12Envelope[T]{
		Env: SOAP12Namespace,
		Body: SOAP12Body[T]{
			Content: data,
		},
	}
}
