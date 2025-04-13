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

type PriceConfigurationRepository struct {
	Mongo *db.MongoRepo
	Cache *cache.RedisCache
}

func NewPriceConfigurationRepository(mongo *db.MongoRepo, cache *cache.RedisCache) *PriceConfigurationRepository {
	return &PriceConfigurationRepository{
		Mongo: mongo,
		Cache: cache,
	}
}

func (r *PriceConfigurationRepository) Create(config *models.PriceConfiguration) error {
	// Verify if a configuration already exists for this user
	existing, _ := r.GetByUserID(config.UserID)
	if existing != nil {
		return fmt.Errorf("price configuration already exists for this user")
	}

	result, err := r.Mongo.Create("price_configurations", config)
	if err != nil {
		return err
	}

	config.ID = result.InsertedID.(primitive.ObjectID)

	// Invalidate cache
	if r.Cache != nil {
		key := fmt.Sprintf("price_config:%s", config.UserID.Hex())
		r.Cache.Delete(context.Background(), key)
	}

	return nil
}

func (r *PriceConfigurationRepository) GetByUserID(userID primitive.ObjectID) (*models.PriceConfiguration, error) {
	// Try to get from cache
	if r.Cache != nil {
		var config models.PriceConfiguration
		key := fmt.Sprintf("price_config:%s", userID.Hex())
		err := r.Cache.Get(context.Background(), key, &config)
		if err == nil {
			return &config, nil
		}
	}

	// If not in cache, get from MongoDB
	var config models.PriceConfiguration
	filter := bson.M{"user_id": userID}
	_, err := r.Mongo.FindOne("price_configurations", filter, &config)
	if err != nil {
		return nil, err
	}

	// Save in cache
	if r.Cache != nil {
		key := fmt.Sprintf("price_config:%s", userID.Hex())
		r.Cache.Set(context.Background(), key, config, 1*time.Hour)
	}

	return &config, nil
}

func (r *PriceConfigurationRepository) Update(userID primitive.ObjectID, amount float64) error {
	filter := bson.M{"user_id": userID}
	update := bson.M{
		"$set": bson.M{
			"amount":     amount,
			"updated_at": time.Now(),
		},
	}

	err := r.Mongo.UpdateOne("price_configurations", filter, update)
	if err != nil {
		return err
	}

	// Invalidate cache
	if r.Cache != nil {
		key := fmt.Sprintf("price_config:%s", userID.Hex())
		r.Cache.Delete(context.Background(), key)
	}

	return nil
}

func (r *PriceConfigurationRepository) Delete(userID primitive.ObjectID) error {
	filter := bson.M{"user_id": userID}
	err := r.Mongo.DeleteOne("price_configurations", filter)
	if err != nil {
		return err
	}

	// Invalidate cache
	if r.Cache != nil {
		key := fmt.Sprintf("price_config:%s", userID.Hex())
		r.Cache.Delete(context.Background(), key)
	}

	return nil
}
