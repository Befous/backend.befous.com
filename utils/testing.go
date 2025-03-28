package utils

import (
	"github.com/Befous/backend.befous.com/helpers"
	"github.com/Befous/backend.befous.com/models"
	"go.mongodb.org/mongo-driver/mongo"
)

func InsertTesting(mongoenv *mongo.Database, datsession models.Session) (interface{}, error) {
	return helpers.InsertOneDoc(mongoenv, "testing", datsession)
}
