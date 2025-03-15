package helpers

import (
	"context"

	"github.com/Befous/api.befous.com/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func MongoConnect(mconn models.DBInfo) (db *mongo.Database, err error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mconn.DBString))
	if err != nil {
		return
	}
	db = client.Database(mconn.DBName)
	return
}

func GetAllDoc[T any](db *mongo.Database, collection string, filter bson.M) (doc T, err error) {
	ctx := context.TODO()
	cur, err := db.Collection(collection).Find(ctx, filter)
	if err != nil {
		return
	}
	defer cur.Close(ctx)
	err = cur.All(ctx, &doc)
	if err != nil {
		return
	}
	return
}
