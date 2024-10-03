package services

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/maxskaink/proyect01-api-go/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collectionUsers *mongo.Collection

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
	return client
}

func CreateUser(newUser *models.User) (models.User, error) {

	if err := newUser.ValidateToCreate(); err != nil {
		return *newUser, err
	}

	insertResult, err := collectionUsers.InsertOne(context.Background(), newUser)

	if err != nil {
		return *newUser, err
	}

	newUser.ID = insertResult.InsertedID.(primitive.ObjectID)

	return *newUser, nil
}

func GetAllUsers() ([]models.User, error) {
	var users []models.User
	cursor, err := collectionUsers.Find(context.Background(), bson.M{})
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
		users = append(users, user)
	}
	return users, nil
}
