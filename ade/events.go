package ade

import (
	"encoding/xml"
)

type Events struct {
	XMLName xml.Name `xml:"events"`
	Events  []Event  `xml:"event"`
}
