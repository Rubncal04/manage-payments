package repository

import (
	"context"
	"fmt"
	"time"

	"github/Rubncal04/youtube-premium/cache"
	"github/Rubncal04/youtube-premium/db"
	"github/Rubncal04/youtube-premium/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ClientRepository struct {
	Mongo *db.MongoRepo
	Cache *cache.RedisCache
}

func NewClientRepository(mongo *db.MongoRepo, cache *cache.RedisCache) *ClientRepository {
	return &ClientRepository{
		Mongo: mongo,
		Cache: cache,
	}
}

func (r *ClientRepository) Create(client models.Client) (models.Client, error) {
	result, err := r.Mongo.Create("clients", client)

	if err != nil {
		return models.Client{}, err
	}

	// Invalidar caché de lista de clientes
	if r.Cache != nil {
		cacheKey := cache.GenerateKey("clients", client.UserID.Hex())
		r.Cache.Delete(context.Background(), cacheKey)
	}

	client.ID = result.InsertedID.(primitive.ObjectID)
	return client, nil
}

func (r *ClientRepository) GetByID(id string) (*models.Client, error) {
	ctx := context.Background()
	cacheKey := cache.GenerateKey("client", id)

	// Intentar obtener de caché
	if r.Cache != nil {
		var client models.Client
		err := r.Cache.Get(ctx, cacheKey, &client)
		if err == nil {
			return &client, nil
		}
	}

	// Si no está en caché, obtener de MongoDB
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid client ID: %v", err)
	}

	filter := bson.M{"_id": objID}
	var client models.Client
	_, err = r.Mongo.FindOne("clients", filter, &client)
	if err != nil {
		return nil, err
	}

	// Guardar en caché
	if r.Cache != nil {
		err = r.Cache.Set(ctx, cacheKey, client, 1*time.Hour)
		if err != nil {
			// Log error pero no fallar la operación
			fmt.Printf("Error caching client: %v\n", err)
		}
	}

	return &client, nil
}

func (r *ClientRepository) Update(id string, updateData bson.M) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid client ID: %v", err)
	}

	filter := bson.M{"_id": objID}
	update := bson.M{
		"$set": updateData,
	}

	err = r.Mongo.UpdateOne("clients", filter, update)
	if err != nil {
		return err
	}

	// Invalidar caché
	if r.Cache != nil {
		cacheKey := cache.GenerateKey("client", id)
		r.Cache.Delete(context.Background(), cacheKey)
	}

	return nil
}

func (r *ClientRepository) GetAll(userID primitive.ObjectID) ([]models.Client, error) {
	ctx := context.Background()
	cacheKey := cache.GenerateKey("clients", userID.Hex())

	// Intentar obtener de caché
	if r.Cache != nil {
		var clients []models.Client
		err := r.Cache.Get(ctx, cacheKey, &clients)
		if err == nil {
			return clients, nil
		}
	}

	// Si no está en caché, obtener de MongoDB
	filter := bson.M{"user_id": userID}
	var clients []models.Client
	err := r.Mongo.FindAll("clients", filter, &clients)
	if err != nil {
		return nil, err
	}

	// Guardar en caché
	if r.Cache != nil {
		err = r.Cache.Set(ctx, cacheKey, clients, 1*time.Hour)
		if err != nil {
			fmt.Printf("Error caching clients: %v\n", err)
		}
	}

	return clients, nil
}

func (r *ClientRepository) UpdateLastPaymentDate(clientID primitive.ObjectID, lastPaymentDate primitive.DateTime) error {
	filter := bson.M{"_id": clientID}
	update := bson.M{
		"$set": bson.M{
			"last_payment_date": lastPaymentDate,
			"updated_at":        time.Now(),
		},
	}

	err := r.Mongo.UpdateOne("clients", filter, update)
	if err != nil {
		return err
	}

	// Invalidar caché
	if r.Cache != nil {
		cacheKey := cache.GenerateKey("client", clientID.Hex())
		r.Cache.Delete(context.Background(), cacheKey)
	}

	return nil
}

func (r *ClientRepository) Delete(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid client ID: %v", err)
	}

	filter := bson.M{"_id": objID}
	return r.Mongo.DeleteOne("clients", filter)
}
