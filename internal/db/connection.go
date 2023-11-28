package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var MongoClient = ConnectionMongodb()

// conexion a mongodb, recibe el contexto, URI de mongo, nombre DB, nombre de Colecion. Retona *Colletion.
func ConnectionMongodb() *mongo.Client {

	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	uri := fmt.Sprintf("mongodb+srv://%s:%s@%s/", os.Getenv("MONGODB_USERNAME"), os.Getenv("MONGODB_PASSWORD"), os.Getenv("MONGODB_CLUSTER_URI"))
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Conexión a MongoDB establecida")

	return client
}

func GetCollection(database, collectionName string, client *mongo.Client) *mongo.Collection {
	return client.Database(database).Collection(collectionName)
}

func PingMongoDB(ctx context.Context, client *mongo.Client) error {
	err := client.Ping(ctx, readpref.Primary())
	if err != nil {
		return fmt.Errorf("couldn't connect to MongoDB: %v", err)
	}
	return nil
}
