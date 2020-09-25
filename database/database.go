package database

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"wallforfry/esiee-api/utils"
)

var logger = utils.InitLogger("database-logger")

var Database = &mongo.Database{}

func CreateMongoDatabase(host string, port int, database string, username string, password string) {
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017").
		SetAuth(options.Credential{
			AuthSource: "esiee-api",
			Username:   username,
			Password:   password,
		})

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		logger.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		logger.Fatal(err)
	}
	Database = client.Database(database)
}
