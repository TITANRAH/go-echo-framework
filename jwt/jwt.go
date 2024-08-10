package jwt

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func GenerarJWT(correo string, nombre string, id string) (string, error) {
	errorVariables := godotenv.Load()

	if errorVariables != nil {
		panic("Error al cargar las variables de entorno")
	}

	miClave := []byte(os.Getenv("SECRET_JWT"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"correo":       correo,
		"nombre":       nombre,
		"generado_por": "Sergio",
		"id":           id,
		"iat":          time.Now().Unix(),
		"exp":          time.Now().Add(time.Hour * 24).Unix(), // 24hrs
	})

	tokenString, err := token.SignedString(miClave)

	return tokenString, err
}
