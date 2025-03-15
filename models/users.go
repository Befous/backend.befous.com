package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Users struct {
	ID       *primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Username string              `json:"username" bson:"username"`
	Password string              `json:"password,omitempty" bson:"password,omitempty"`
	Nama     string              `json:"nama" bson:"nama"`
	No_Telp  string              `json:"no_telp" bson:"no_telp"`
	Email    string              `json:"email" bson:"email"`
	Role     string              `json:"role" bson:"role"`
}
