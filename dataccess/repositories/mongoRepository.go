package repositories

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/maxskaink/proyect01-api-go/dataccess"
	"github.com/maxskaink/proyect01-api-go/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoUserRepository struct {
	collectionUsers *mongo.Collection
	client          *mongo.Client
}

// NewMongoRepository the initializate the conections with the mongo database.
// create the collections of the document for persistence
// it also, create some rules for the database
// and return a instance of MongoRepository
func NewUserMongoRepository() *MongoUserRepository {
	var collectionsUsers *mongo.Collection

	MONGOT_URI := os.Getenv("MONGO_URI")

	if MONGOT_URI == "" {
		log.Fatal("MONGO_URI  is required, load .env file")
	}

	clientOptions := options.Client().ApplyURI(MONGOT_URI)
	client, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	if err := client.Ping(context.Background(), nil); err != nil {
		log.Fatal(err)
	}

	collectionsUsers = client.Database("users_manager").Collection("clients")

	indexModel := mongo.IndexModel{
		Keys:    bson.M{"email": 1},
		Options: options.Index().SetUnique(true),
	}

	_, err = collectionsUsers.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		log.Fatal(err)
	}

	return &MongoUserRepository{
		collectionUsers: collectionsUsers,
		client:          client,
	}
}

// CloseClient Close the conections with de database
func (mr *MongoUserRepository) CloseClient() {
	mr.client.Disconnect(context.Background())
}

// CreateUser recive the newUser to be persisntace in de database
// return the newUser and error if it exist
func (mr *MongoUserRepository) CreateUser(newUser *models.User) (models.User, error) {

	insertResult, err := mr.collectionUsers.InsertOne(context.Background(), newUser)

	if err != nil {
		return *newUser, err
	}

	newUser.ID = insertResult.InsertedID.(primitive.ObjectID)
	newUser.Password = ""
	return *newUser, nil
}

// GetAllUsers return all the active users in the database
// recieve the page and max users per page, it also return
// a error if it exist
func (mr *MongoUserRepository) GetAllUsers(page int, maxUsers int) ([]models.User, error) {
	if page <= 0 || maxUsers <= 0 {
		return []models.User{}, nil
	}

	var users []models.User

	offset := (page - 1) * maxUsers

	findOption := options.Find()
	findOption.SetLimit(int64(maxUsers))
	findOption.SetSkip(int64(offset))

	cursor, err := mr.collectionUsers.Find(context.Background(), bson.M{
		"isActive": true,
	}, findOption)
	if err != nil {
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

// GetUserById recieve de id to search in the database
// if it exist it will be return and an error if it exist
func (mr *MongoUserRepository) GetUserByID(id string) (models.User, error) {
	user := new(models.User)
	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return *user, err
	}

	response := mr.collectionUsers.FindOne(context.Background(), bson.M{
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

// GetTotalUsers it wiil return the number of total active users
// it also return a error if it exist
func (mr *MongoUserRepository) GetTotalUsers() (int, error) {
	total, err := mr.collectionUsers.CountDocuments(context.Background(), bson.M{
		"isActive": true,
	})
	if err != nil {
		return 0, err
	}
	return int(total), nil
}

// GetUserByEmialAndPass will return a user in the database with this information
// the password must be encrypt if is needed. it return de user and error if exist
func (mr *MongoUserRepository) GetUserByEmailAndPass(email string, password string) (*models.User, error) {

	response := mr.collectionUsers.FindOne(context.Background(), bson.M{
		"email":    email,
		"password": password,
		"isActive": true,
	})
	if response.Err() != nil {
		return &models.User{}, response.Err()
	}

	userToLogin := new(models.User)

	err := response.Decode(&userToLogin)

	if err != nil {
		return &models.User{}, err
	}

	return userToLogin, nil
}

// ReplaceUser replace de user with id of the parameter, and the information of the
// model of newUser, it return the old User and an error if it exist
func (mr *MongoUserRepository) ReplaceUser(newUser *models.User, id string) (models.User, error) {

	ObjectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return *newUser, err
	}

	newUser.ID = ObjectId
	newUser.IsActive = true

	response := mr.collectionUsers.FindOneAndReplace(context.Background(), bson.M{
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

// UpdateUserById would update the user with the parameter id,
// all the information must be prepare, and it will return the old user
// and an error if it exist
func (mr *MongoUserRepository) UpdateUserById(newUser *models.User, id string) (*models.User, error) {
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
		updateData["password"] = newUser.Password
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

	err = mr.collectionUsers.FindOneAndUpdate(context.Background(), filter, update, opts).Decode(oldUser)

	if err != nil {
		return newUser, err
	}

	return oldUser, nil
}

// DeleteUserById wll delete the user with de paramaeter id
// and return the deleted user of the database. it also return an error
// if it exist
func (mr *MongoUserRepository) DeleteUserById(id string) (*models.User, error) {
	ObjectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	update := bson.M{
		"$set": bson.M{
			"isActive": false,
		},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	filter := bson.M{
		"_id": ObjectID,
	}

	var deletedUser = new(models.User)

	err = mr.collectionUsers.FindOneAndUpdate(context.Background(), filter, update, opts).Decode(deletedUser)

	return deletedUser, err
}

var _ dataccess.IUserRepository = (*MongoUserRepository)(nil)
