package repository

import (
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func Init(connectionString string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
	if err != nil {
		return err
	}

	// Ping MongoDB to ensure the connection is live.
	if err = client.Ping(ctx, nil); err != nil {
		return err
	}

	logger.Info("Connected to MongoDB!")

	return nil
}

func Close() error {
	if client == nil {
		return errors.New("mongo client not initialized")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return client.Disconnect(ctx)
}

func getCollection(database, collectionName string) (*mongo.Collection, error) {
	if client == nil {
		return nil, errors.New("mongo client not initialized")
	}
	return client.Database(database).Collection(collectionName), nil
}

func WriteMongo(database, collectionName string, document any) error {
	collection, err := getCollection(database, collectionName)
	if err != nil {
		return err
	}

	_, err = collection.InsertOne(context.Background(), document)
	if err != nil {
		return err
	}

	log.Println("Inserted document")
	return nil
}

func UpdateMongo(database, collectionName string, document any, filter any) error {
	collection, err := getCollection(database, collectionName)
	if err != nil {
		return err
	}

	_, err = collection.UpdateOne(context.Background(), filter, document)
	if err != nil {
		return err
	}

	log.Println("Updated document")
	return nil
}

func ReadMongo[T any](database, collectionName string, filter any) (T, error) {
	var result T

	collection, err := getCollection(database, collectionName)
	if err != nil {
		return result, err
	}

	err = collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func CountMongo(database, collectionName string, filter any) (int64, error) {
	collection, err := getCollection(database, collectionName)
	if err != nil {
		return 0, err
	}

	count, err := collection.CountDocuments(context.Background(), filter)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func DeleteMongo(database, collectionName string, filter any) (int64, error) {
	collection, err := getCollection(database, collectionName)
	if err != nil {
		return 0, err
	}

	result, err := collection.DeleteMany(context.Background(), filter)
	if err != nil {
		return 0, err
	}

	return result.DeletedCount, nil
}
