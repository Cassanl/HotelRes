package db

import (
	"hoteRes/types"

	"go.mongodb.org/mongo-driver/mongo"
)

type UserStore interface {
	GetUserById(string) (*types.User, error)
}

type MongoUserStore struct {
	client *mongo.Client
}

func NewMongoUserStore(c *mongo.Client) *MongoUserStore {
	return &MongoUserStore{
		client: c,
	}
}

func (m *MongoUserStore) GetUserById(id string) (*types.User, error) {
	return nil, nil
}
