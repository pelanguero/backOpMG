package main

import (
	"log"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var jwtkey = []byte(os.Getenv("JWT_KEY"))

//verifica el token (jwt) returna 0 si el token esta bien, 1 si la firma es invalida, 2 si el token no es valido y -1 si no se hizo la peticion de manera correcta
func verificarjwt(jjwt string, clai *Claims) int {

	tkn, errorr := jwt.ParseWithClaims(jjwt, clai, func(token *jwt.Token) (interface{}, error) {
		return jwtkey, nil
	})

	if errorr != nil {
		if errorr == jwt.ErrSignatureInvalid {
			return 1
		}
		return -1
	}
	if !tkn.Valid {
		return 2
	}

	return 0
}

//verifica el pw con hash y el posible pw plano
func verificarpw(hashedpw string, plain string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedpw), []byte(plain))
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

//middleware que maneja los cors
func corsmiddle() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
		c.Header("Access-Control-Allow-Methods", "PUT, POST, DELETE, GET")
	}

}

//es para no guardar el pw plano
func hashpw(pwd string) string {
	bytestr := []byte(pwd)
	hashh, err := bcrypt.GenerateFromPassword(bytestr, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hashh)
}

//middleware auth
func autoriza() gin.HandlerFunc {
	return func(c *gin.Context) {
		claim := &Claims{}
		jwwt, errror := c.Cookie("t")
		if errror != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Al leer la cookie " + errror.Error()})
			return
		}
		statuss := verificarjwt(jwwt, claim)
		if 0 == statuss {
			c.JSON(http.StatusOK, gin.H{})
			return
		} else if statuss == 1 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No estas Autorizado"})
			return
		} else if statuss == -1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Peticion invalida"})
			return
		} else if statuss == 2 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No estas Autorizado, token no valido"})
			return
		}
	}

}
