package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var ClientMySQL = ConnectDatabase()

func ConnectDatabase() *sql.DB {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, assuming environment variables are set (e.g., in Docker)")
	}

	username := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_ROOT_PASSWORD")
	database := os.Getenv("MYSQL_DATABASE")
	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")

	datasourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", username, password, host, port, database)

	db, err := sql.Open("mysql", datasourceName)
	if err != nil {
		log.Fatal(err)
	}

	return db
}
