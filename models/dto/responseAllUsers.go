package dto

import (
	"github.com/maxskaink/proyect01-api-go/models"
)

type ResponseAllUsers struct {
	Page           int           `json:"page"`
	TotalPages     int           `json:"totalPages"`
	TotalUsersPage int           `json:"totalUsersPage"`
	Data           []models.User `json:"users"`
}
