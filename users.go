package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/http"
)

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type PasswordUpdate struct {
	Email       string `json:"email"`
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

func GetAllUsers() ([]*User, error) {
	var users []*User

	client, ctx, cancel := getConnection()
	defer cancel()
	defer client.Disconnect(ctx)
	db := client.Database("humetis")
	collection := db.Collection("users")
	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	err = cursor.All(ctx, &users)
	if err != nil {
		log.Printf("Failed marshalling %v", err)
		return nil, err
	}
	return users, nil
}

// GetUserByUsername  get a user by its username from the db
func GetUserByUsername(email string) (*User, error) {
	var user *User

	client, ctx, cancel := getConnection()
	defer cancel()
	defer client.Disconnect(ctx)
	db := client.Database("humetis")
	collection := db.Collection("users")
	result := collection.FindOne(ctx, bson.M{"email": email})
	if result == nil {
		return nil, errors.New("Could not find a user")
	}
	err := result.Decode(&user)

	if err != nil {
		log.Printf("Failed marshalling %v", err)
		return nil, err
	}
	log.Printf("Users: %v", user)
	return user, nil
}

func HandlerUpdatePassword(c *gin.Context) {
	var u PasswordUpdate
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}
	var user, err = Update(&u)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}

//Update updating an existing task in a mongo
func Update(cp *PasswordUpdate) (*User, error) {
	var updatedUser *User
	client, ctx, cancel := getConnection()
	defer cancel()
	defer client.Disconnect(ctx)

	_, err := client.Database("humetis").Collection("users").UpdateOne(ctx, bson.D{{"email", cp.Email}}, bson.D{{"$set", bson.D{{"password", cp.NewPassword}}}})
	if err != nil {
		log.Printf("Could not save Task: %v", err)
		return nil, err
	}
	return updatedUser, nil
}
