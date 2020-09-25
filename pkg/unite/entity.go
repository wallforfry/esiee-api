package unite

import "go.mongodb.org/mongo-driver/bson/primitive"

type Unite struct {
	ID    primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Code  string             `json:"code" bson:"code,omitempty"`
	Label string             `json:"label" bson:"label,omitempty"`
}
