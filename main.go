package main

import (
	"clubhub-hotel-management/cmd/server/routes"
	"clubhub-hotel-management/internal/db"
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	engine := gin.Default()

	mongodb := db.MongoClient

	router := routes.NewRouter(engine, mongodb)
	router.MapRoutes()
	defer mongodb.Disconnect(context.Background())
	engine.Run(":8080")
}
