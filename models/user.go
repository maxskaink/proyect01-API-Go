package models

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/maxskaink/proyect01-api-go/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represent the entity of Users, as json response or entity
// for the database
type User struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name,omitempty" bson:"name"`
	Email       string             `json:"email,omitempty" bson:"email"`
	Password    string             `json:"password,omitempty" bson:"password"`
	IsActive    bool               `json:"isActive,omitempty" bson:"isActive"`
	LastSession primitive.DateTime `json:"lastSession,omitempty" bson:"lastSession"`
}

// ValidateToCreate validate if the information of the struct
// have the enough and correct data for been stored
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
	if len(u.Password) < 8 {
		return errors.New("PASSWORD MUST BE AT LEAST 8 CHARACTERS")
	}

	return nil
}

// ValidateToUpdate validate if the information of the struct
// have the enough and correc data for been updated
func (u *User) ValidateToUpdate() error {
	if u.ID != primitive.NilObjectID {
		return errors.New("ID MUST BE EMPTY")
	}
	return u.ValidateToCreate()
}

// CreateJWT Generate a Json Web Token, with the name and email of the struct
func (u *User) CreateJWT() (string, error) {
	claims := jwt.MapClaims{
		"name":  u.Name,
		"email": u.Email,
		"exp":   time.Now().Add(time.Hour * 12).Unix(),
	}

	return utils.CreateJWT(claims)
}
