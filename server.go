package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("fallo al cargar .env: " + err.Error())
	}
	router := gin.Default()
	//userRoutes:=router.Group()
	router.PUT("/registro", crearUsuarioruta)
	router.GET("/")
	router.Run(":" + os.Getenv("PUERTO"))
}

//conexion a mongodb retorna un cliente de mongo,contexto
func mongoConnection() (*mongo.Client, context.Context, context.CancelFunc) {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error al cargar .env")
	}

	connectTimeout := 5

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(connectTimeout)*time.Second)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_CONNECTION")))
	if err != nil {
		log.Printf("Fallo al crear el cliente: %v", err)
	}

	if err != nil {
		log.Printf("Fallo al conectar con el cluster: %v", err)
	}

	// Force a connection to verify our connection string
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Printf("Failed to ping cluster: %v", err)
	}

	fmt.Println("Conectado a MongoDB")
	return client, ctx, cancel
}
