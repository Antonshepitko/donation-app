package routes

import (
	"context"
	"donation-backend/middleware"
	"donation-backend/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var WebSocketBroadcast func(donation models.Donation)

func RegisterDonationRoutes(r *gin.Engine, db *mongo.Database) {
	donations := r.Group("/api")

	donations.Use(middleware.CorsMiddleware()) // CORS защита (если нужно)

	// Получить все донаты для текущего стримера (JWT)
	donations.GET("/donations", middleware.AuthMiddleware(), func(c *gin.Context) {
		username := c.GetString("username")
		cursor, err := db.Collection("donations").Find(
			context.TODO(),
			bson.M{"streamer": username},
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
			return
		}
		var results []models.Donation
		if err := cursor.All(context.TODO(), &results); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing data"})
			return
		}
		c.JSON(http.StatusOK, results)
	})

	// Публичный роут: создать донат
	donations.POST("/donate", func(c *gin.Context) {
		var d models.Donation
		if err := c.BindJSON(&d); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
			return
		}
		d.Timestamp = time.Now()

		// Проверим, что стример существует
		count, _ := db.Collection("users").CountDocuments(context.TODO(), bson.M{"username": d.Streamer})
		if count == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Streamer not found"})
			return
		}

		_, err := db.Collection("donations").InsertOne(context.TODO(), d)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
			return
		}

		// Уведомим по WebSocket (если подключен)
		if WebSocketBroadcast != nil {
			WebSocketBroadcast(d)
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Donation added"})
	})
}
