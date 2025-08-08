package main

import (
	"context"
	"donation-backend/routes"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func main() {
	r := gin.Default()

	// Подключение к MongoDB
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://admin:prefectdinorah@donation-mongo:27017"))
	if err != nil {
		log.Fatal(err)
	}
	db := client.Database("donationdb")

	// Подключаем роуты
	routes.RegisterAuthRoutes(r, db)
	routes.RegisterDonationRoutes(r, db)
	routes.RegisterWebSocketRoute(r)

	// Запускаем сервер
	r.Run(":5000")
}
