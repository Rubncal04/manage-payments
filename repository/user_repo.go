// repository/user_repo.go
package repository

import (
	"fmt"
	"github/Rubncal04/youtube-premium/db"
	"github/Rubncal04/youtube-premium/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepository struct {
	Mongo *db.MongoRepo
}

func NewUserRepository(mongo *db.MongoRepo) *UserRepository {
	return &UserRepository{Mongo: mongo}
}

func (r *UserRepository) GetAllUsers() ([]models.User, error) {
	var users []models.User
	err := r.Mongo.FindAll("users", "", &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) CreateUser(user *models.User) (models.User, error) {
	coll, err := r.Mongo.Create("users", user)
	newUser := models.User{
		ID:        coll.InsertedID.(primitive.ObjectID),
		Name:      user.Name,
		CellPhone: user.CellPhone,
		Paid:      user.Paid,
	}

	if err != nil {
		return models.User{}, err
	}

	return newUser, nil
}

func (r *UserRepository) UpdateUser(userID string, updateData bson.M) error {
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %v", err)
	}

	filter := bson.M{"_id": objID}

	// Asegurarse de que updateData use $set
	update := bson.M{
		"$set": updateData,
	}

	return r.Mongo.UpdateOne("users", filter, update)
}
