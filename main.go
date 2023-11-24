package main

import (
	"clubhub-hotel-management/cmd/server/routes"
	"clubhub-hotel-management/utils/db"

	"github.com/gin-gonic/gin"
)

func main() {

	engine := gin.Default()

	mysql := db.ClientMySQL
	defer mysql.Close()

	router := routes.NewRouter(engine, mysql)
	router.MapRoutes()
	engine.Run(":8080")
}
