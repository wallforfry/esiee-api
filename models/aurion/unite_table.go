package aurion

import (
	"encoding/xml"
)

type UniteTable struct {
	XMLName xml.Name   `xml:"table"`
	Rows    []UniteRow `xml:"row"`
}
