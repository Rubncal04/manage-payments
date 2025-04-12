package db

import (
	"context"
	"log"
	"time"

	"github/Rubncal04/youtube-premium/config"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepo struct {
	Client *mongo.Client
	Db     *mongo.Database
}

// NewMongoRepo inicializa la conexión a MongoDB
func NewMongoRepo(variables *config.EnvVariables) (*MongoRepo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(variables.MONGO_URI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
		return nil, err
	}

	// Verifica la conexión
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Could not ping MongoDB: %v", err)
		return nil, err
	}

	log.Println("Connected to MongoDB")

	return &MongoRepo{
		Client: client,
		Db:     client.Database(variables.MONGO_DB),
	}, nil
}

// Desconectar la conexión
func (m *MongoRepo) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := m.Client.Disconnect(ctx); err != nil {
		log.Fatalf("Error disconnecting from MongoDB: %v", err)
	} else {
		log.Println("Disconnected from MongoDB")
	}
}

func (m *MongoRepo) FindAll(collectionName string, filter any, result any) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if str, ok := filter.(string); ok && str == "" {
		filter = bson.M{}
	}

	collection := m.Db.Collection(collectionName)
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, result); err != nil {
		return err
	}

	return nil
}

func (m *MongoRepo) FindOne(collectionName string, filter any, result any) (any, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if str, ok := filter.(string); ok && str == "" {
		filter = bson.M{}
	}

	collection := m.Db.Collection(collectionName)
	err := collection.FindOne(ctx, filter).Decode(result)
	if err != nil {
		return nil, err
	}

	return collection, nil
}

func (m *MongoRepo) Create(collectionName string, coll any) (*mongo.InsertOneResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := m.Db.Collection(collectionName)
	collec, err := collection.InsertOne(ctx, coll)
	if err != nil {
		return nil, err
	}

	return collec, nil
}

func (m *MongoRepo) UpdateOne(collectionName string, filter any, update any) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := m.Db.Collection(collectionName)
	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

// UpdateMany updates multiple documents in a collection
func (m *MongoRepo) UpdateMany(collectionName string, filter any, update any) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := m.Db.Collection(collectionName)
	result, err := collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// DeleteOne deletes a single document from a collection
func (m *MongoRepo) DeleteOne(collectionName string, filter any) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := m.Db.Collection(collectionName)
	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}
