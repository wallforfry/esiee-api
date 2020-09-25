package aurion

import (
	"encoding/xml"
	"fmt"
)

type UniteRow struct {
	XMLName xml.Name `xml:"row",csv:"-"`
	Code    string   `xml:"Code.Unité",csv:"Code.Unité"`
	Label   string   `xml:"Libellé.Unité",csv:"Libellé.Unité"`
}

func (a UniteRow) String() string {
	return fmt.Sprintf("Code : %s, Label : %s\n", a.Code, a.Label)
}

func (a UniteRow) CSV() []string {
	return []string{a.Code, a.Label}
}
