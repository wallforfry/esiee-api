package aurion

import (
	"encoding/xml"
)

type GroupTable struct {
	XMLName xml.Name   `xml:"table"`
	Rows    []GroupRow `xml:"row"`
}
