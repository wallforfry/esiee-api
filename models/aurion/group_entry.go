package aurion

import "fmt"

type GroupEntry struct {
	Unite  string   `json:"unite"`
	Groups []string `json:"groups"`
}

func (g GroupEntry) String() string {
	return fmt.Sprintf("%s %s\n", g.Unite, g.Groups)
}
