package services

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/maxskaink/proyect01-api-go/models"
	"github.com/maxskaink/proyect01-api-go/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collectionUsers *mongo.Collection

// InitDataBase the initializate the conections with the mongo database.
// create the collections of the document for persistence
// it also, create some rules for the database
func InitDataBase() *mongo.Client {
	MONGO_URI := os.Getenv("MONGO_URI")

	if MONGO_URI == "" {
		log.Fatal("MONGO_URI is required or load .env file")
	}

	clientOptions := options.Client().ApplyURI(MONGO_URI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Ping(context.Background(), nil); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB")
	collectionUsers = client.Database("users_manager").Collection("clients")

	indexModel := mongo.IndexModel{
		Keys:    bson.M{"email": 1}, // Índice en el campo Email
		Options: options.Index().SetUnique(true),
	}

	_, err = collectionUsers.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		log.Fatal(err)
	}

	return client
}

// CreateUser create the newUser validating the information of the structure
// if it is validate, create the users on the data base
// return error if is any error in the process
func CreateUser(newUser *models.User) (models.User, error) {

	if err := newUser.ValidateToCreate(); err != nil {
		return *newUser, err
	}

	newUser.Password = utils.GetHash(newUser.Password)
	newUser.IsActive = true
	newUser.LastSession = primitive.NewDateTimeFromTime(time.Now())

	insertResult, err := collectionUsers.InsertOne(context.Background(), newUser)

	if err != nil {
		return *newUser, err
	}

	newUser.ID = insertResult.InsertedID.(primitive.ObjectID)
	newUser.Password = ""
	return *newUser, nil
}

// GetAllUsers obtains as maxUsers active users of the database
// and a part of the segmen users. iF its any error it will be returned
func GetAllUsers(page int, maxUsers int) ([]models.User, error) {
	if page <= 0 || maxUsers <= 0 {
		return []models.User{}, nil
	}

	var users []models.User

	offset := (page - 1) * maxUsers

	findOption := options.Find()
	findOption.SetLimit(int64(maxUsers))
	findOption.SetSkip(int64(offset))

	cursor, err := collectionUsers.Find(context.Background(), bson.M{
		"isActive": true,
	}, findOption)
	if err != nil {
		fmt.Println(err)
		return users, err
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return users, err
		}
		user.Password = ""
		user.IsActive = false
		users = append(users, user)
	}
	return users, nil
}

// GetUserByID obtains a users searching by ID
// if it doesnt exist return an error, but if it exist
// return de struct of the user
func GetUserByID(id string) (models.User, error) {
	user := new(models.User)
	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return *user, err
	}

	response := collectionUsers.FindOne(context.Background(), bson.M{
		"_id": objectID,
	})

	if response.Err() != nil {
		return *user, response.Err()
	}
	err = response.Decode(user)

	if err != nil {
		return *user, err
	}

	return *user, nil
}

// GetTotalUsers obtains the total of active user in the database
// just return a number, and if its any error it also return it
func GetTotalUsers() (int, error) {
	total, err := collectionUsers.CountDocuments(context.Background(), bson.M{
		"isActive": true,
	})
	if err != nil {
		return 0, err
	}
	return int(total), nil
}

// LoginUser verify if the email and password match, if it doesnt, it return an error
// but if its correct, it return a Json Web Token with name and email information
func LogInUser(email string, password string) (string, error) {

	if !utils.IsEmail(email) {
		return "", fmt.Errorf("EMAIL IS REQUIRED AN VALID")
	}

	response := collectionUsers.FindOne(context.Background(), bson.M{
		"email":    email,
		"password": utils.GetHash(password),
	})
	if response.Err() != nil {
		return "", response.Err()
	}

	userToLogin := new(models.User)

	err := response.Decode(&userToLogin)

	if err != nil {
		return "", err
	}

	jwt, err := userToLogin.CreateJWT()

	if err != nil {
		return "", err
	}
	return jwt, nil
}

// ReplaceUser update the user with the id, and return the old user
// if is any problem it return it.
func ReplaceUser(newUser *models.User, id string) (models.User, error) {
	if err := newUser.ValidateToUpdate(); err != nil {
		return *newUser, err
	}

	ObjectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return *newUser, err
	}

	newUser.ID = ObjectId
	newUser.Password = utils.GetHash(newUser.Password)
	newUser.IsActive = true

	response := collectionUsers.FindOneAndReplace(context.Background(), bson.M{
		"_id": ObjectId,
	}, newUser)

	if response.Err() != nil {
		return *newUser, response.Err()
	}

	oldUser := new(models.User)
	err = response.Decode(oldUser)

	if err != nil {
		return *newUser, err
	}

	return *oldUser, err
}

// UpdateUser update some values fo the info of some user
// it also verify the information and if its correct, it will updated
// reeturn the old user, if its any error it will be returned
func UpdateUser(newUser *models.User, id string) (*models.User, error) {

	ObjectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return newUser, err
	}

	updateData := make(map[string]interface{})

	if newUser.Name != "" {
		updateData["name"] = newUser.Name
	}
	if newUser.Email != "" {
		updateData["email"] = newUser.Email
	}
	if newUser.Password != "" {
		if len(newUser.Password) < 8 {
			return newUser, fmt.Errorf("PASSWORD MUST BE AT LEAST 8 CHARACTERS")
		}
		updateData["password"] = utils.GetHash(newUser.Password)
	}

	if len(updateData) == 0 {
		return newUser, fmt.Errorf("NO FIELD TO UPDATE")
	}

	update := bson.M{
		"$set": updateData,
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.Before)

	filter := bson.M{
		"_id": ObjectID,
	}

	var oldUser = new(models.User)

	err = collectionUsers.FindOneAndUpdate(context.Background(), filter, update, opts).Decode(oldUser)

	if err != nil {
		return newUser, err
	}

	return oldUser, nil
}
