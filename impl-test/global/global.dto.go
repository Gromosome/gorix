package global

import "encoding/xml"

type ResponseDTO[T any] struct {
	XMLName xml.Name `json:"-" xml:"response"`

	Success bool `json:"success" xml:"success"`
	Data    T    `json:"data" xml:"data"`
}

type ErrorDTO struct {
	XMLName xml.Name `json:"-" xml:"error"`
	Success bool     `json:"success" xml:"success"`
	Error   any      `json:"error" xml:"error"`
	Message string   `json:"message" xml:"message"`
}
