package mongodb

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type mongodb struct {
	collection *mongo.Collection
}

