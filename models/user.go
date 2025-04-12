package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// En models/user.go
type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string             `bson:"name" json:"name"`
	Username  string             `bson:"username" json:"username"`
	Email     string             `bson:"email" json:"email"`
	Password  string             `bson:"password" json:"-"`
	Role      string             `bson:"role" json:"role"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

// HashPassword hashes the user's password
func (u *User) HashPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword checks if the provided password matches the user's hashed password
func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

// NewUser creates a new user with the provided information
func NewUser(username, email, password string) (*User, error) {
	user := &User{
		Username:  username,
		Email:     email,
		Role:      "user", // Default role
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := user.HashPassword(password); err != nil {
		return nil, err
	}

	return user, nil
}
