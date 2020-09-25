package group

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"regexp"
)

type Group struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username string             `json:"username" bson:"username,omitempty"`
	Email    string             `json:"email" bson:"email,omitempty"`
	Unite    string             `json:"unite" bson:"unite,omitempty"`
	Groups   []string           `json:"groups" bson:"groups,omitempty"`
}

func (g Group) String() string {
	return fmt.Sprintf("%s %s\n", g.Unite, g.Groups)
}

func CreateGroupFromAurionEntry(username string, email string, uniteCode string) Group {
	yearRe := regexp.MustCompile(`^\d\d_`)
	uniteCode = yearRe.ReplaceAllString(uniteCode, "")

	promoRe := regexp.MustCompile(`^\w{1,4}_`)
	uniteCode = promoRe.ReplaceAllString(uniteCode, "")

	uniteRe := regexp.MustCompile(`^([^_]*)_([^_]*)_(.*)`)
	result := uniteRe.FindStringSubmatch(uniteCode)
	if len(result) == 0 {
		return Group{Unite: uniteCode, Groups: []string{}, Username: username, Email: email}
	}

	unite := fmt.Sprintf("%s-%s", result[1], result[2])
	return Group{Unite: unite, Groups: result[3:], Username: username, Email: email}
}
