package unite

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
	collection := database.Collection("unites")
	return &repo{
		collection: collection,
	}
}

func (r *repo) Find(id primitive.ObjectID) (*Unite, error) {
	var result Unite
	err := r.collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&result)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *repo) FindByCode(code string) (*Unite, error) {
	var result Unite
	err := r.collection.FindOne(context.TODO(), bson.M{"code": code}).Decode(&result)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *repo) FindByLabel(label string) (*Unite, error) {
	var result Unite
	err := r.collection.FindOne(context.TODO(), bson.M{"label": label}).Decode(&result)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *repo) FindAll() ([]Unite, error) {
	var unites []Unite
	cursor, err := r.collection.Find(context.TODO(), bson.M{})

	if err != nil {
		return nil, err
	}

	err = cursor.All(context.TODO(), &unites)

	if err != nil {
		return nil, err
	}

	return unites, nil
}

func (r *repo) Count() (int64, error) {
	return r.collection.CountDocuments(context.TODO(), bson.M{})
}

func (r *repo) Update(unite *Unite) error {
	upsert := true
	_, err := r.collection.ReplaceOne(
		context.TODO(),
		bson.M{"code": unite.Code},
		unite,
		&options.ReplaceOptions{Upsert: &upsert},
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *repo) Store(unite *Unite) (*primitive.ObjectID, error) {
	insertResult, err := r.collection.InsertOne(context.TODO(), unite)

	if err != nil {
		return nil, err
	}

	id := insertResult.InsertedID.(primitive.ObjectID)
	return &id, nil
}
