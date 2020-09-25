package group

import "go.mongodb.org/mongo-driver/bson/primitive"

type Reader interface {
	Find(id primitive.ObjectID) (*Group, error)
	FindByUnite(unite string) (*Group, error)
	FindByUsername(username string) ([]Group, error)
}

type Writer interface {
	Update(group *Group) error
	Store(group *Group) (*primitive.ObjectID, error)
	StoreMany(groups []Group) error
	Dump() error
}

type Repository interface {
	Reader
	Writer
}
