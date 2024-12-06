package services

import (
	"time"

	"github.com/maxskaink/proyect01-api-go/dataccess"
	custom_errors "github.com/maxskaink/proyect01-api-go/errors"
	"github.com/maxskaink/proyect01-api-go/models"
	"github.com/maxskaink/proyect01-api-go/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UsersService struct {
	Repository dataccess.IUserRepository
}

// NewUsersService return a instance of UsersService
// it needs de repository to use
func NewUsersService(repository dataccess.IUserRepository) *UsersService {
	return &UsersService{
		Repository: repository,
	}
}

// CreateUser create the newUser validating the information of the structure
// if it is validate, create the users on the repository
// return error if is any error in the process
func (us *UsersService) CreateUser(newUser *models.User) (models.User, error) {

	if err := newUser.ValidateToCreate(); err != nil {
		return *newUser, err
	}

	newUser.Password = utils.GetHash(newUser.Password)
	newUser.IsActive = true
	newUser.LastSession = primitive.NewDateTimeFromTime(time.Now())

	createdUser, err := us.Repository.CreateUser(newUser)

	createdUser.Password = ""
	return createdUser, err
}

// GetAllUsers obtains as maxUsers active users of the database
// and a part of the segmen users. iF its any error it will be returned
func (us *UsersService) GetAllUsers(page int, maxUsers int) ([]models.User, error) {
	users, err := us.Repository.GetAllUsers(page, maxUsers)

	return users, err
}

// GetUserByID obtains a users searching by ID
// if it doesnt exist return an error, but if it exist
// return de struct of the user
func (us *UsersService) GetUserByID(id string) (models.User, error) {
	user, err := us.Repository.GetUserByID(id)

	return user, err
}

// GetTotalUsers obtains the total of active user in the respository
// just return a number, and if its any error it also return it
func (us *UsersService) GetTotalUsers() (int, error) {
	total, err := us.Repository.GetTotalUsers()

	return int(total), err
}

// LoginUser verify if the email and password match, if it doesnt, it return an error
// but if its correct, it return a Json Web Token with name and email information
func (us *UsersService) LogInUser(email string, password string) (string, error) {

	if !utils.IsEmail(email) {
		return "", custom_errors.NewInvalidFormat(
			"EMAIL IS INVALID",
			"EMAIL",
		)
	}

	userToLogin, err := us.Repository.GetUserByEmailAndPass(email, utils.GetHash(password))

	if err != nil {
		return "", err
	}

	return userToLogin.CreateJWT()
}

// ReplaceUser update the user with the id, and return the old user
// if is any problem it return it.
func (us *UsersService) ReplaceUser(newUser *models.User, id string) (models.User, error) {
	if err := newUser.ValidateToUpdate(); err != nil {
		return *newUser, err
	}

	newUser.Password = utils.GetHash(newUser.Password)

	oldUser, err := us.Repository.ReplaceUser(newUser, id)

	return oldUser, err
}

// UpdateUserById update some values fo the info of some user
// it also verify the information and if its correct, it will updated
// reeturn the old user, if its any error it will be returned
func (us *UsersService) UpdateUserById(newUser *models.User, id string) (*models.User, error) {
	if newUser.Password != "" {
		if len(newUser.Password) < 8 {
			return newUser, custom_errors.NewInvalidFormat(
				"PASSWORD MUST BE AT LEAST 8 CHARACTERS",
				"PASSWORD",
			)
		}
		newUser.Password = utils.GetHash(newUser.Password)
	}

	oldUser, err := us.Repository.UpdateUserById(newUser, id)

	return oldUser, err
}

// DeleteUserById update de statate of the user isAvtive to false
// it woul return the user info and if its any error it will be returned
func (us *UsersService) DeleteUserById(id string) (*models.User, error) {
	deletedUser, err := us.Repository.DeleteUserById(id)

	return deletedUser, err
}
