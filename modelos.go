package main

import (
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type rol int8

const (
	ADMINISTRADOR = 1
	DEFAULT       = 0
)

type Claims struct {
	Correo string `json:"correo"`
	jwt.StandardClaims
}

//modelo de usuario
type Usuario struct {
	ID primitive.ObjectID
	//verificar el tamaño de los strings
	Nombre string
	Correo string
	//verificar el minimo tamaño de el pw
	Clave string
	About string
	//por defecto --Default--
	Role      rol
	Historial []string
}
