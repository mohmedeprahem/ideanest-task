package repository

import (
	"context"
	database "ideanest-task/pkg/database/mongodb"
	"ideanest-task/pkg/database/mongodb/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	userCollection *mongo.Collection
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		userCollection: database.GetCollection("users"),
	}
}

func (r *UserRepository) Create(user *models.User) error {
	_, err := r.userCollection.InsertOne(context.Background(), user)
	return err
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User

	err := r.userCollection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)

	return &user, err
}
