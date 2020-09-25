package group

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
	collection := database.Collection("groups")
	return &repo{
		collection: collection,
	}
}

func (r *repo) Find(id primitive.ObjectID) (*Group, error) {
	var result Group
	err := r.collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&result)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *repo) FindByUnite(code string) (*Group, error) {
	var result Group
	err := r.collection.FindOne(context.TODO(), bson.M{"unite": code}).Decode(&result)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *repo) FindByUsername(username string) ([]Group, error) {
	var Groups []Group
	cursor, err := r.collection.Find(context.TODO(), bson.M{"username": username})

	if err != nil {
		return nil, err
	}

	err = cursor.All(context.TODO(), &Groups)

	if err != nil {
		return nil, err
	}

	return Groups, nil
}

func (r *repo) Update(group *Group) error {
	upsert := true
	_, err := r.collection.ReplaceOne(
		context.TODO(),
		bson.M{"username": group.Username, "unite": group.Unite},
		group,
		&options.ReplaceOptions{Upsert: &upsert},
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *repo) Store(group *Group) (*primitive.ObjectID, error) {
	insertResult, err := r.collection.InsertOne(context.TODO(), group)

	if err != nil {
		return nil, err
	}

	id := insertResult.InsertedID.(primitive.ObjectID)
	return &id, nil
}

func (r *repo) StoreMany(groups []Group) error {
	y := make([]interface{}, len(groups))
	for i, v := range groups {
		y[i] = v
	}
	_, err := r.collection.InsertMany(context.TODO(), y)
	return err
}

func (r *repo) Dump() error {
	_, err := r.collection.DeleteMany(context.TODO(), bson.D{})
	return err
}
