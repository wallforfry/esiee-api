package unite

import "go.mongodb.org/mongo-driver/bson/primitive"

type Reader interface {
	Find(id primitive.ObjectID) (*Unite, error)
	FindByCode(code string) (*Unite, error)
	FindByLabel(label string) (*Unite, error)
	FindAll() ([]Unite, error)
}

type Writer interface {
	Update(unite *Unite) error
	Store(unite *Unite) (*primitive.ObjectID, error)
}

type Repository interface {
	Reader
	Writer
}
