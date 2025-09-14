package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"donation-backend/middleware"
	"donation-backend/routes"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Database

func main() {
	// 1) Mongo URI: сначала берем из env (совместимо с нашим Deployment),
	// при отсутствии — можно задать дефолт (ниже пример на mongodb сервисе)
	mongoURI := os.Getenv("MONGO_URL")
	if mongoURI == "" {
		// безопаснее завершиться с ошибкой, но чтобы не падать — оставлю дефолт.
		// Подстрой под свои реальные креды/базу:
		mongoURI = "mongodb://root:password@mongodb:27017/donation?authSource=admin"
		log.Println("WARN: MONGO_URL is empty; using default:", mongoURI)
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("MongoDB ping failed:", err)
	}

	// Базу лучше получать из URI — но если неудобно, оставим конкретную
	db = client.Database("donationdb") // <- при желании можно распарсить из URI

	// Уникальный индекс на username
	users := db.Collection("users")
	_, err = users.Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys:    bson.D{{Key: "username", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		log.Println("Index creation error:", err)
	}

	r := gin.Default()
	r.Use(middleware.CorsMiddleware())

	// Аутентификация/вебсокеты/health
	r.POST("/api/register", routes.Register(db))
	r.POST("/api/login", routes.Login(db))
	r.GET("/api/ws", routes.WebSocketHandler)
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// 2) Регистрируем РОУТЫ ДОНАТОВ (раньше их не было, из-за этого был 404)
	routes.RegisterDonationRoutes(r, db)

	// Запуск
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	log.Println("Server running on port", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Server failed:", err)
	}
}
