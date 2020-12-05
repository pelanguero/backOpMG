package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//rol del usuario

//crea un usuario
func crearUsuarioruta(c *gin.Context) {
	var userr Usuario
	err := c.ShouldBindJSON(&userr)
	if err != nil {
		log.Print(err)
		log.Print(userr.Nombre)
		c.JSON(http.StatusBadRequest, gin.H{"msg": err})
		return
	}
	userr.Clave = hashpw(userr.Clave)
	//"no es coneccion es conexion"
	id, err, insertod := Create(&userr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err})
		return
	}
	if insertod {
		c.JSON(http.StatusOK, gin.H{"error": "El Usuario ya existe"})
	} else {
		c.JSON(http.StatusOK, gin.H{"id": id})
	}
}

//crea el usuario en la base de datos
func Create(user *Usuario) (primitive.ObjectID, error, bool) {
	//if len(user.Clave)<10
	var oid primitive.ObjectID
	var existe bool
	client, ctx, cancel := mongoConnection()
	filtro := bson.M{"correo": user.Correo}
	var testt Usuario
	defer cancel()
	defer client.Disconnect(ctx)
	user.ID = primitive.NewObjectID()
	user.Role = DEFAULT
	errorr := client.Database("omgtest").Collection("usuarios").FindOne(context.TODO(), filtro).Decode(&testt)
	if errorr != nil {
		log.Println(errorr)
	}
	if testt.Correo != user.Correo {
		result, err := client.Database("omgtest").Collection("usuarios").InsertOne(ctx, user)
		if err != nil {
			log.Printf("No se pudo agregar el usuario: %v", err)
			return primitive.NilObjectID, err, false
		}
		oid = result.InsertedID.(primitive.ObjectID)
		existe = false
	} else {
		log.Println("El Usuario ya existe")
		existe = true
	}

	return oid, nil, existe
}

//inicio de sesion retorna un jwt, crea una cookie
func iniciosesion(c *gin.Context) {

	var creds Credenciales
	var testt Usuario
	err := c.ShouldBindJSON(&creds)
	if err != nil {
		log.Print(err)
		log.Print(creds.Correo)
		c.JSON(http.StatusLocked, gin.H{"msg": err})
		return
	}
	filtro := bson.M{"correo": creds.Correo}
	client, ctx, cancel := mongoConnection()
	defer cancel()
	defer client.Disconnect(ctx)

	errorr := client.Database("omgtest").Collection("usuarios").FindOne(context.TODO(), filtro).Decode(&testt)
	if errorr != nil {
		log.Println(errorr)
	}
	if testt.Correo == creds.Correo && verificarpw(testt.Clave, creds.Clave) {
		expirationTime := time.Now().Add(24 * time.Hour)
		claims := &Claims{
			Correo: creds.Correo,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expirationTime.Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, errr := token.SignedString(jwtkey)
		if errr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": errr})
		}
		c.SetCookie("t", tokenString, expirationTime.Second()-time.Now().Second(), "/iniciodesesion", "localhost", false, true)
		c.JSON(http.StatusAccepted, gin.H{"Name": "token", "token": tokenString, "Expira": expirationTime, "usuario": creds.Correo})
	} else {

		c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciales no validas"})
	}

}

//busca un usuario con el id provisto
func usuarioxid(id string) Usuario {
	var retorno Usuario
	pro, errror := primitive.ObjectIDFromHex(id)
	if errror != nil {
		fmt.Println(errror.Error())
	}
	filtro := bson.M{"id": pro}
	client, ctx, cancel := mongoConnection()
	defer cancel()
	defer client.Disconnect(ctx)
	erorr := client.Database("omgtest").Collection("usuarios").FindOne(ctx, filtro).Decode(&retorno)
	if erorr != nil {
		fmt.Println(erorr)
	}

	return retorno
}
