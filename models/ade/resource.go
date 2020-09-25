package ade

import (
	"encoding/xml"
	"fmt"
)

type Resource struct {
	XMLName  xml.Name `xml:"resource"`
	Category string   `xml:"category,attr"`
	Name     string   `xml:"name,attr"`
}

func (r Resource) String() string {
	return fmt.Sprintf("\t  %s : %s\n", r.Category, r.Name)
}
