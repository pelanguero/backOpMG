package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
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
	router.PUT("/iniciodesesion", iniciosesion)
	router.GET("/test", autoriza(), prueba)
	router.GET("/user/:id", autoriza(), userr)
	router.GET("/")
	//router.Use(corsmiddle())
	router.Run(":" + os.Getenv("PUERTO"))
}

func prueba(c *gin.Context) {
	if c.Request.Response.Status == "200 OK" {
		c.JSON(http.StatusOK, gin.H{"exito": "funciona bien"})
	}
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

//cors
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control,token, X-Requested-With,access-control-allow-origin")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

//da la informacion del usuario a partir del id de la url
func userr(c *gin.Context) {
	if c.Request.Response.StatusCode == 202 {
		user := usuarioxid(c.Param("id"))
		if user.Nombre != "" {
			c.JSON(http.StatusOK, user)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "usuario no existe"})
		}
	}
}
