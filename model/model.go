package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID            primitive.ObjectID   `json:"_id,omitempty" bson:"_id,omitempty"`
	UserName      string               `json:"username"   validate:"required,min=5,max=100"`
	First_name    string               `json:"first_name" validate:"required,min=2,max=100"`
	Last_name     string               `json:"last_name" validate:"required,min=2,max=100"`
	Password      string               `json:"password,omitempty" validate:"required,min=8,max=100"`
	Email         string               `json:"email"    validate:"required"`
	Friends       []primitive.ObjectID `json:"friends"`
	profilePicUrl string               `json:"profilepicurl"`
	Bio           string               `json:"bio"`
	Created_at    time.Time            `json:"created_at"`
	Updated_at    time.Time            `json:"updated_at"`
	Token         string               `json:"token"`
	Refresh_token string               `json:"refresh_token"`
	User_id       string               `json:"user_id"`
}

type Post struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Description string             `json:"description"`
	File        string             `json:"file" validate:"required"`
	Comments    []Comment          `json:"comments,omitempty" bson:"comments, omitempty"`
	User        primitive.ObjectID `json:"user"`
	Likes       []Like             `json:"likes,omitempty" bson:"likes, omitempty"`
	Created_at  time.Time          `json:"created_at`
}

type Comment struct {
	User       primitive.ObjectID `json:"user" validate:"required"`
	Date       time.Time          `json:"date" validate:"required"`
	Comment    string             `json:"comment" validate:"required"`
	Created_at time.Time          `json:"created_at`
}

type Like struct {
	User primitive.ObjectID `json:"user"   validate:"required"`
	Date time.Time          `json:"date" validate:"required"`
}

//"62b07dc97c5bbfdd83d27dc6"
