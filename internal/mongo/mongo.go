package mongo

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB(uri string, dbname string) (*mongo.Database, error) {
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environmental variable.")
	}
	
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	return client.Database(dbname), nil
}

func CloseMongoDB(db *mongo.Database) error {
	return db.Client().Disconnect(context.Background())
}