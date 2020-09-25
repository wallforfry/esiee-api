package aurion

import (
	"encoding/xml"
	"fmt"
)

type GroupRow struct {
	XMLName xml.Name `xml:"row",csv:"-"`
	Login   string   `xml:"login.Individu",csv:"login.Individu"`
	Mail    string   `xml:"Coordonnée.Coordonnée",csv:"Coordonnée.Coordonnée"`
	Group   string   `xml:"Code.Groupe",csv:"Code.Groupe"`
}

func (a GroupRow) String() string {
	return fmt.Sprintf("\t Username : %s, Email : %s, Groups : %s \n", a.Login, a.Mail, a.Group)
}

func (a GroupRow) CSV() []string {
	return []string{a.Login, a.Mail, a.Group}
}
