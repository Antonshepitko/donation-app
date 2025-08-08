package routes

import (
	"context"
	"donation-backend/models"
	"donation-backend/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func RegisterAuthRoutes(r *gin.Engine, db *mongo.Database) {
	auth := r.Group("/api")
	auth.POST("/register", func(c *gin.Context) {
		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
			return
		}
		// Хешируем пароль
		hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
		user.Password = string(hash)
		_, err := db.Collection("users").InsertOne(context.TODO(), user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username taken or error"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "User registered"})
	})

	auth.POST("/login", func(c *gin.Context) {
		var input models.User
		if err := c.BindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
			return
		}
		var dbUser models.User
		err := db.Collection("users").FindOne(context.TODO(), bson.M{"username": input.Username}).Decode(&dbUser)
		if err != nil || bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(input.Password)) != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}
		token, _ := utils.GenerateToken(input.Username)
		c.JSON(http.StatusOK, gin.H{"token": token})
	})
}
