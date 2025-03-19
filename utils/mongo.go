package utils

import (
	"log"
	"os"
	"sync"

	"github.com/Befous/backend.befous.com/helpers"
	"github.com/Befous/backend.befous.com/models"
	"go.mongodb.org/mongo-driver/mongo"
)

var dbInstance *mongo.Database
var once sync.Once

func SetConnection() *mongo.Database {
	once.Do(func() {
		dbInfo := models.DBInfo{
			DBString: os.Getenv("mongoenv"),
			DBName:   "befous",
		}
		var err error
		dbInstance, err = helpers.MongoConnect(dbInfo)
		if err != nil {
			log.Fatalf("Failed to connect to MongoDB: %v", err)
		}
	})
	return dbInstance
}
