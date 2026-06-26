package global

import (
	"encoding/xml"
)

type ErrorDTO struct {
	XMLName xml.Name `json:"-" xml:"error"`
	Success bool     `json:"success" xml:"success"`
	Error   any      `json:"error" xml:"error"`
	Message string   `json:"message" xml:"message"`
}
