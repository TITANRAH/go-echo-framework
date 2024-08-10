package rutas

import (
	"context"
	"echo-framework/database"
	"echo-framework/dto"
	"echo-framework/jwt"
	"echo-framework/middleware_custom"
	"echo-framework/validaciones"
	"encoding/json"
	"net/http"

	echo "github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// "password": "p2gHNiENUw"
func Seguridad_registro(c echo.Context) error {

	var body dto.UsuarioDto

	err := json.NewDecoder(c.Request().Body).Decode(&body)

	// error en la peticion
	if err != nil {
		respuesta := map[string]interface{}{
			"error":   "error",
			"mensaje": "Error en la peticion",
		}

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(respuesta)
	}

	// validador
	if len(body.Nombre) == 0 || len(body.Correo) == 0 || len(body.Password) == 0 {
		respuesta := map[string]interface{}{
			"error":   "error",
			"mensaje": "Faltan datos",
		}

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(respuesta)
	}

	if validaciones.Regex_correo.FindStringSubmatch(body.Correo) == nil {
		respuesta := map[string]interface{}{
			"error":   "error",
			"mensaje": "Correo invalido",
		}

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(respuesta)
	}

	if !validaciones.ValidarPassword(body.Password) {
		respuesta := map[string]interface{}{
			"error":   "error",
			"mensaje": "Password invalido",
		}

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(respuesta)
	}

	var user bson.M

	// SI ES IGUAL A NIL NO HAY ERROR QUIERE DECIR QUE ENCONTRO USUARIO Y ESO NO PERMITIRA GUARDAR USUARIO

	if err := database.UsuariosCollection.FindOne(context.TODO(), bson.D{{Key: "correo", Value: body.Correo}}).Decode(&user); err == nil {

		respuesta := map[string]interface{}{
			"error":   "error",
			"mensaje": "Usuario ya registrado",
		}

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(respuesta)
	}

	costo := 8
	bytes, _ := bcrypt.GenerateFromPassword([]byte(body.Password), costo)
	registro := bson.D{{Key: "nombre", Value: body.Nombre}, {Key: "telefono", Value: body.Telefono}, {Key: "correo", Value: body.Correo}, {Key: "password", Value: string(bytes)}}
	database.UsuariosCollection.InsertOne(context.TODO(), registro)

	// retorno respuesta
	respuesta := map[string]string{
		"estado":  "ok",
		"mensaje": "Usuario registrado correctamente",
	}

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c.Response().WriteHeader(http.StatusCreated)

	return json.NewEncoder(c.Response()).Encode(respuesta)

}


func Seguridad_protegida(c echo.Context) error {
	if middleware_custom.ValidarJWT(c) == 0 {
		respuesta := map[string]interface{}{
			"error":   "error",
			"mensaje": "No autorizado",
		}

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c.Response().WriteHeader(http.StatusUnauthorized)
		return json.NewEncoder(c.Response()).Encode(respuesta)
	}

	respuesta := map[string]string{
		"estado":  "ok",
		"mensaje": "Metodo protegido",

	}

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(c.Response()).Encode(respuesta)
}

func Seguridad_login(c echo.Context) error {

	// tomo los datos ingresados en un form login
	var body dto.LoginDto
	// los paso por json
	err := json.NewDecoder(c.Request().Body).Decode(&body)
	// defino si la peticion da error
	if err != nil {
		respuesta := map[string]interface{}{
			"error":   "error",
			"mensaje": "Error en la peticion",
		}

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(respuesta)
	}

	// validacion de correo
	if validaciones.Regex_correo.FindStringSubmatch(body.Correo) == nil {
		respuesta := map[string]interface{}{
			"error":   "error",
			"mensaje": "Correo invalido",
		}

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(respuesta)
	}

	// validacion de password
	if !validaciones.ValidarPassword(body.Password) {
		respuesta := map[string]interface{}{
			"error":   "error",
			"mensaje": "Password invalido",
		}

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(respuesta)
	}

	// busco el usuario
	var user bson.M
	if err := database.UsuariosCollection.FindOne(context.TODO(), bson.D{{Key: "correo", Value: body.Correo}}).Decode(&user); err != nil {
		respuesta := map[string]interface{}{
			"error":   "error",
			"mensaje": "Usuario no encontrado",
		}

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(respuesta)
	}

	// convierto el pass ingresado en el form login

	passwordBytes := []byte(body.Password)

	// tomo el pass guardado en bd
	passwordBD := []byte(user["password"].(string))

	// defino si hay error en la comparacion de pass
	errPass := bcrypt.CompareHashAndPassword(passwordBD, passwordBytes)

	if errPass != nil {
		respuesta := map[string]interface{}{
			"error":   "error",
			"mensaje": "Password incorrecto",
		}

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(respuesta)
	} else {

		// si no hay error genero el token
		stringObjectID := user["_id"].(primitive.ObjectID).Hex()

		// llamo a la funcion que gnera el token
		jwtKey, err := jwt.GenerarJWT(user["correo"].(string), user["nombre"].(string), stringObjectID)

		// si hay error respondo error
		if err != nil {
			respuesta := map[string]interface{}{
				"error":   "error",
				"mensaje": "Error al generar token",
			}

			c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			c.Response().WriteHeader(http.StatusBadRequest)
			return json.NewEncoder(c.Response()).Encode(respuesta)
		} else {

			// si no hay error retorno el token
			retorno := dto.LoginRespuestaDto{
				Token:  jwtKey,
				Nombre: user["nombre"].(string),
			}

			c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			c.Response().WriteHeader(http.StatusOK)

			return json.NewEncoder(c.Response()).Encode(retorno)
		}
	}

}
