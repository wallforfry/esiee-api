package event

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type repo struct {
	collection *mongo.Collection
}

func NewMongoRepository(database *mongo.Database) Repository {
	collection := database.Collection("events")
	return &repo{
		collection: collection,
	}
}

func (r *repo) Find(id primitive.ObjectID) (*Event, error) {
	var result Event
	err := r.collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&result)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *repo) FindAll() ([]Event, error) {
	var events []Event
	cursor, err := r.collection.Find(context.TODO(), bson.M{})

	if err != nil {
		return nil, err
	}

	err = cursor.All(context.TODO(), &events)

	if err != nil {
		return nil, err
	}

	return events, nil
}

func (r *repo) Count() (int64, error) {
	return r.collection.CountDocuments(context.TODO(), bson.M{})
}

func (r *repo) Update(event *Event) error {
	upsert := true
	_, err := r.collection.ReplaceOne(
		context.TODO(),
		bson.M{"event_id": event.EventId},
		event,
		&options.ReplaceOptions{Upsert: &upsert},
	)

	return err
}

func (r *repo) Store(event *Event) (*primitive.ObjectID, error) {
	insertResult, err := r.collection.InsertOne(context.TODO(), event)

	if err != nil {
		return nil, err
	}

	id := insertResult.InsertedID.(primitive.ObjectID)
	return &id, nil
}

func (r *repo) StoreMany(events []Event) error {
	y := make([]interface{}, len(events))
	for i, v := range events {
		y[i] = v
	}
	_, err := r.collection.InsertMany(context.TODO(), y)
	return err
}

func (r *repo) Dump() error {
	_, err := r.collection.DeleteMany(context.TODO(), bson.D{})
	return err
}
