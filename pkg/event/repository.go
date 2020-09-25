package event

import "go.mongodb.org/mongo-driver/bson/primitive"

type Reader interface {
	Find(id primitive.ObjectID) (*Event, error)
	FindAll() ([]Event, error)
	Count() (int64, error)
}

type Writer interface {
	Update(event *Event) error
	Store(event *Event) (*primitive.ObjectID, error)
	StoreMany(events []Event) error
	Dump() error
}

type Repository interface {
	Reader
	Writer
}
