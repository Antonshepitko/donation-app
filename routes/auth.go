package routes

import (
    "context"
    "errors"
    "net/http"
    "time"

    "donation-backend/utils"

    "github.com/gin-gonic/gin"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "golang.org/x/crypto/bcrypt"
)

type User struct {
    Username string `json:"username" bson:"username"`
    Password string `json:"password" bson:"password"`
}

func Register(db *mongo.Database) gin.HandlerFunc {
    return func(c *gin.Context) {
        var user User
        if err := c.ShouldBindJSON(&user); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
            return
        }

        hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
            return
        }
        user.Password = string(hashedPassword)

        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()

        _, err = db.Collection("users").InsertOne(ctx, user)
        if err != nil {
            var we *mongo.WriteException
            if errors.As(err, &we) {
                for _, e := range we.WriteErrors {
                    if e.Code == 11000 {
                        c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
                        return
                    }
                }
            }
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
            return
        }

        c.JSON(http.StatusCreated, gin.H{"message": "User registered"})
    }
}

func Login(db *mongo.Database) gin.HandlerFunc {
    return func(c *gin.Context) {
        var user User
        if err := c.ShouldBindJSON(&user); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
            return
        }

        var found User
        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()

        err := db.Collection("users").FindOne(ctx, bson.M{"username": user.Username}).Decode(&found)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
            return
        }

        if bcrypt.CompareHashAndPassword([]byte(found.Password), []byte(user.Password)) != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
            return
        }

        token, err := utils.GenerateJWT(user.Username)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
            return
        }

        c.JSON(http.StatusOK, gin.H{"token": token})
    }
}
