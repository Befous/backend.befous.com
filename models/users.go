package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Users struct {
	ID       *primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Username string              `json:"username" bson:"username"`
	Password string              `json:"password,omitempty" bson:"password,omitempty"`
	Nama     string              `json:"nama,omitempty" bson:"nama,omitempty"`
	Phone    string              `json:"phone,omitempty" bson:"phone,omitempty"`
	Email    string              `json:"email,omitempty" bson:"email,omitempty"`
	Role     string              `json:"role" bson:"role"`
}

type Session struct {
	ID         *primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Username   string              `json:"username" bson:"username"`
	Token      string              `json:"token" bson:"token"`
	User_Agent string              `json:"user_agent" bson:"user_agent"`
	Expire_At  time.Time           `json:"expire_at" bson:"expire_at"`
	Issued_At  time.Time           `json:"issued_at" bson:"issued_at"`
}
