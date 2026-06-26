package context

import "encoding/xml"

type SOAPFaultCode interface {
	IsValid(status SOAPStatusCode) bool
}
type SOAP11FaultCode string

const (
	SOAP11FaultVersionMismatch SOAP11FaultCode = "soapenv:VersionMismatch"
	SOAP11FaultMustUnderstand  SOAP11FaultCode = "soapenv:MustUnderstand"
	SOAP11FaultClient          SOAP11FaultCode = "soapenv:Client"
	SOAP11FaultServer          SOAP11FaultCode = "soapenv:Server"
)

func (s SOAP11FaultCode) IsValid(status SOAPStatusCode) bool {
	switch s {
	case "":
		return status == SOAPStatusOK || status == SOAPStatusAccepted

	case SOAP11FaultClient:
		return status.Int() >= 400 && status.Int() < 500

	case SOAP11FaultServer:
		return status.Int() >= 500 && status.Int() < 600

	case SOAP11FaultVersionMismatch:
		return status == SOAPStatusBadRequest

	case SOAP11FaultMustUnderstand:
		return status == SOAPStatusBadRequest

	default:
		return false
	}
}

type SOAP11Fault struct {
	XMLName     xml.Name        `xml:"soapenv:Fault"`
	FaultCode   SOAP11FaultCode `xml:"faultcode"`
	FaultString string          `xml:"faultstring"`
	FaultActor  string          `xml:"faultactor,omitempty"`
	Detail      any             `xml:"detail,omitempty"`
}
type SOAP12FaultCode string

const (
	SOAP12FaultVersionMismatch     SOAP12FaultCode = "env:VersionMismatch"
	SOAP12FaultMustUnderstand      SOAP12FaultCode = "env:MustUnderstand"
	SOAP12FaultDataEncodingUnknown SOAP12FaultCode = "env:DataEncodingUnknown"
	SOAP12FaultSender              SOAP12FaultCode = "env:Sender"
	SOAP12FaultReceiver            SOAP12FaultCode = "env:Receiver"
)

func (s SOAP12FaultCode) IsValid(status SOAPStatusCode) bool {
	switch s {
	case "":
		return status == SOAPStatusOK || status == SOAPStatusAccepted
	case SOAP12FaultSender:
		return status.Int() >= 400 && status.Int() < 500
	case SOAP12FaultReceiver:
		return status.Int() >= 500 && status.Int() < 600
	case SOAP12FaultVersionMismatch:
		return status == SOAPStatusBadRequest
	case SOAP12FaultMustUnderstand:
		return status == SOAPStatusBadRequest
	case SOAP12FaultDataEncodingUnknown:
		return status == SOAPStatusUnsupportedMediaType ||
			status == SOAPStatusBadRequest
	default:
		return false
	}
}

type SOAP12Fault struct {
	XMLName xml.Name     `xml:"env:Fault"`
	Code    SOAP12Code   `xml:"env:Code"`
	Reason  SOAP12Reason `xml:"env:Reason"`
	Node    string       `xml:"env:Node,omitempty"`
	Role    string       `xml:"env:Role,omitempty"`
	Detail  any          `xml:"env:Detail,omitempty"`
}

type SOAP12Code struct {
	Value   SOAP12FaultCode `xml:"env:Value"`
	Subcode *SOAP12Subcode  `xml:"env:Subcode,omitempty"`
}

type SOAP12Subcode struct {
	Value string `xml:"env:Value"`
}

type SOAP12Reason struct {
	Text SOAP12Text `xml:"env:Text"`
}

type SOAP12Text struct {
	Lang  string `xml:"xml:lang,attr,omitempty"`
	Value string `xml:",chardata"`
}

func NewSOAP11FaultEnvelope(
	code SOAP11FaultCode,
	message string,
	detail any,
) SOAP11Fault {
	return SOAP11Fault{
		FaultCode:   code,
		FaultString: message,
		Detail:      detail,
	}
}

func NewSOAP12FaultEnvelope(
	code SOAP12FaultCode,
	message string,
	detail any,
) SOAP12Fault {
	return SOAP12Fault{
		Code: SOAP12Code{
			Value: code,
		},
		Reason: SOAP12Reason{
			Text: SOAP12Text{
				Lang:  "en",
				Value: message,
			},
		},
		Detail: detail,
	}
}
