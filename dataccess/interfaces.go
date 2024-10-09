package dataccess

import "github.com/maxskaink/proyect01-api-go/models"

type IUserRepository interface {
	CloseClient()
	CreateUser(newUser *models.User) (models.User, error)
	GetAllUsers(page int, maxUsers int) ([]models.User, error)
	GetUserByID(id string) (models.User, error)
	GetTotalUsers() (int, error)
	GetUserByEmailAndPass(email string, password string) (*models.User, error)
	ReplaceUser(newUser *models.User, id string) (models.User, error)
	UpdateUserById(newUser *models.User, id string) (*models.User, error)
	DeleteUserById(id string) (*models.User, error)
}
