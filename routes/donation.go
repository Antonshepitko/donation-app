package routes

import (
	"context"
	"donation-backend/middleware"
	"donation-backend/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive" // <- добавили
	"go.mongodb.org/mongo-driver/mongo"
)

var WebSocketBroadcast func(donation models.Donation)

func RegisterDonationRoutes(r *gin.Engine, db *mongo.Database) {
	donations := r.Group("/api")

	donations.Use(middleware.CorsMiddleware()) // CORS

	// Список донатов ТЕКУЩЕГО стримера (по JWT)
	donations.GET("/donations", middleware.AuthMiddleware(), func(c *gin.Context) {
		username := c.GetString("username")

		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		cur, err := db.Collection("donations").Find(ctx, bson.M{"streamer": username})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
			return
		}
		defer cur.Close(ctx)

		var results []models.Donation
		if err := cur.All(ctx, &results); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing data"})
			return
		}
		c.JSON(http.StatusOK, results)
	})

	// Получить донат по id (только свой)
	donations.GET("/donations/:id", middleware.AuthMiddleware(), func(c *gin.Context) {
		username := c.GetString("username")
		oid, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		var d models.Donation
		err = db.Collection("donations").FindOne(ctx, bson.M{
			"_id":      oid,
			"streamer": username,
		}).Decode(&d)
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
			return
		}
		c.JSON(http.StatusOK, d)
	})

	// Публичный роут: создать донат (оригинальный путь)
	createHandler := func(c *gin.Context) {
		var d models.Donation
		if err := c.ShouldBindJSON(&d); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
			return
		}
		d.Timestamp = time.Now().UTC()

		// Проверим, что стример существует
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		count, err := db.Collection("users").CountDocuments(ctx, bson.M{"username": d.Streamer})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
			return
		}
		if count == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Streamer not found"})
			return
		}

		res, err := db.Collection("donations").InsertOne(ctx, d)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
			return
		}
		if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
			d.ID = oid // если в модели есть поле bson:"_id,omitempty"
		}

		// Уведомим по WebSocket
		if WebSocketBroadcast != nil {
			WebSocketBroadcast(d)
		}

		c.JSON(http.StatusCreated, d) // отдаём созданный донат
	}

	donations.POST("/donate", createHandler)     // оставляем старый путь
	donations.POST("/donations", createHandler)  // и добавляем REST-алиас

	// Удалить донат по id (только свой)
	donations.DELETE("/donations/:id", middleware.AuthMiddleware(), func(c *gin.Context) {
		username := c.GetString("username")
		oid, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		res, err := db.Collection("donations").DeleteOne(ctx, bson.M{
			"_id":      oid,
			"streamer": username,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
			return
		}
		if res.DeletedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.Status(http.StatusNoContent)
	})
}
