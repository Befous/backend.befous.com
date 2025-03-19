package utils

import (
	"github.com/Befous/backend.befous.com/helpers"
	"github.com/Befous/backend.befous.com/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func InsertUser(mongoenv *mongo.Database, datauser models.Users) (interface{}, error) {
	return helpers.InsertOneDoc(mongoenv, "users", datauser)
}

func GetAllUser(mongoenv *mongo.Database) ([]models.Users, error) {
	return helpers.GetAllDoc[models.Users](mongoenv, "users")
}

func GetAllUserWithPagination(mongoenv *mongo.Database, page, limit int) ([]models.Users, int64, error) {
	filter := bson.M{}
	return helpers.GetAllDocByFilterWithPagination[models.Users](mongoenv, "users", page, limit, filter)
}

func FindUser(mongoenv *mongo.Database, datauser models.Users) models.Users {
	filter := bson.M{"username": datauser.Username}
	return helpers.GetOneDoc[models.Users](mongoenv, "users", filter)
}

func IsPasswordValid(mongoenv *mongo.Database, datauser models.Users) bool {
	filter := bson.M{"username": datauser.Username}
	res := helpers.GetOneDoc[models.Users](mongoenv, "users", filter)
	return helpers.CheckPasswordHash(datauser.Password, res.Password)
}

func UsernameExists(mongoenv *mongo.Database, datauser models.Users) bool {
	filter := bson.M{"username": datauser.Username}
	return helpers.DocExists(mongoenv, "users", filter, datauser)
}

func UpdateUser(mongoenv *mongo.Database, datauser models.Users) interface{} {
	filter := bson.M{"username": datauser.Username}
	return helpers.ReplaceOneDoc(mongoenv, "users", filter, datauser)
}

func DeleteUser(mongoenv *mongo.Database, datauser models.Users) interface{} {
	filter := bson.M{"username": datauser.Username}
	return helpers.DeleteOneDoc(mongoenv, "users", filter)
}

func GetSession(mongoenv *mongo.Database, tokenString, username string) models.Session {
	filter := bson.M{"token": tokenString, "username": username}
	return helpers.GetOneLatestDoc[models.Session](mongoenv, "session", filter)
}

func InsertSession(mongoenv *mongo.Database, datsession models.Session) (interface{}, error) {
	return helpers.InsertOneDoc(mongoenv, "session", datsession)
}
