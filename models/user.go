package models

import (
	"errors"

	"github.com/maxskaink/proyect01-api-go/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name,omitempty" bson:"name,omitempty"`
	Email       string             `json:"email,omitempty" bson:"email,omitempty"`
	Password    string             `json:"password,omitempty" bson:"password,omitempty"`
	IsActive    bool               `json:"isActive,omitempty" bson:"isActive,omitempty"`
	LastSession primitive.DateTime `json:"lastSession,omitempty" bson:"lastSession,omitempty"`
}

func (u *User) ValidateToCreate() error {

	if u.Name == "" {
		return errors.New("NAME IS REQUIRED")
	}
	if !utils.IsEmail(u.Email) {
		return errors.New("EMAIL IS REQUIRED AN DVALID")
	}
	if u.Password == "" {
		return errors.New("PASSWORD IS REQUIRED")
	}

	return nil
}
