package middleware_custom

import (
	"echo-framework/database"
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	jwt "github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/joho/godotenv"
	echo "github.com/labstack/echo/v4"
)

func ValidarJWT(c echo.Context) int {
	errorVariables := godotenv.Load()
	if errorVariables != nil {
		http.Error(c.Response(), "Error con variables de entorno", http.StatusUnauthorized)
		return 0
	}
	miClave := []byte(os.Getenv("SECRET_JWT"))
	header := c.Request().Header.Get("Authorization")
	//fmt.Println(len(token))
	if len(header) == 0 {

		//http.Error(c.Response(), "No autorizado", http.StatusUnauthorized)
		return 0
	}
	splitBearer := strings.Split(header, " ")
	if len(splitBearer) != 2 {
		//http.Error(c.Response(), "No autorizado", http.StatusUnauthorized)
		return 0
	}
	splitToken := strings.Split(splitBearer[1], ".")
	if len(splitToken) != 3 {
		//http.Error(c.Response(), "No autorizado", http.StatusUnauthorized)
		return 0
	}
	tk := strings.TrimSpace(splitBearer[1])
	token, err := jwt.Parse(tk, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: ")

		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return miClave, nil
	})
	if err != nil {
		//http.Error(c.Response(), "No autorizado", http.StatusUnauthorized)
		return 0
	}
	// if err != nil {
	// 	//http.Error(c.Response(), "No autorizado", http.StatusUnauthorized)
	// 	return 0
	// }
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		var user bson.M
		if err := database.UsuariosCollection.FindOne(context.TODO(), bson.M{
			"correo": claims["correo"],
		}).Decode(&user); err != nil {
			//http.Error(c.Response(), "No autorizado", http.StatusUnauthorized)
			return 0
		}
		return 1
		//fmt.Println(user["nombre"])
	} else {
		//http.Error(c.Response(), "No autorizado", http.StatusUnauthorized)
		return 0
	}
	return 1
}
