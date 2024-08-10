package rutas

import (
	"context"
	"echo-framework/database"
	"echo-framework/dto"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gosimple/slug"
	echo "github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Categoria_get(c echo.Context) error {

	findOptions := options.Find()

	// LLAMO A LA COLECCION  Y APLICO METODO FIND DE MONGO DRIVER Y COMO TERCER ARGUMENTOS PUEDO PONER CONDIDCIONES
	// EN LA BUSQUEDA SI QUISIERA
	cursor, err := database.CategoriaCollection.Find(context.TODO(), bson.D{}, findOptions.SetSort(
		bson.D{{Key: "_id", Value: -1}},
	))

	if err != nil {
		panic(err)
	}

	// asi decaro que recibire un arreglo desde mongo
	var results []bson.M

	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	c.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(c.Response()).Encode(results)
}

func Categoria_get_con_parametro(c echo.Context) error {

	// tengo que convertir el id que viene en los params para poder hacer la busqueda en mongo
	objId, _ := primitive.ObjectIDFromHex(c.Param("id"))

	// en este caso no recibe un arreglo si no que un valor unico

	var results bson.M

	if err := database.CategoriaCollection.FindOne(context.TODO(), bson.D{{Key: "_id", Value: objId}}).Decode(&results); err != nil {
		response := map[string]string{
			"estado":  "error",
			"mensaje": "No se encontro la categoria",
		}
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c.Response().WriteHeader(http.StatusBadRequest)

		return json.NewEncoder(c.Response()).Encode(response)
	}

	fmt.Println(results)

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	c.Response().WriteHeader(http.StatusBadRequest)
	return json.NewEncoder(c.Response()).Encode(results)
}

func Categoria_post(c echo.Context) error {
	var body dto.CategoriaDto
	if err := json.NewDecoder(c.Request().Body).Decode(&body); err != nil {
		response := map[string]string{
			"estado":  "error",
			"mensaje": "Ocurrio un error en el body",
		}
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(response)
	}

	if len(body.Nombre) == 0 {
		response := map[string]string{
			"estado":  "error",
			"mensaje": "El nombre es requerido",
		}
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(response)
	}

	// guardamos en la bd

	registro := bson.D{
		{Key: "nombre", Value: body.Nombre},
		{Key: "slug", Value: slug.Make(body.Nombre)},
	}

	database.CategoriaCollection.InsertOne(context.TODO(), registro)

	respuesta := map[string]string{
		"estado":  "ok",
		"mensaje": "Categoria creada con exito",
	}

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c.Response().WriteHeader(http.StatusCreated)

	return json.NewEncoder(c.Response()).Encode(respuesta)

}

func Categoria_put(c echo.Context) error {

	// declaro lo que recibire
	var body dto.CategoriaDto

	err := json.NewDecoder(c.Request().Body).Decode(&body)

	// detecto que traiga el body
	if err != nil {
		response := map[string]string{
			"estado":  "error",
			"mensaje": "Ocurrio un error en el body",
		}
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(response)
	}

	// detecto que traiga el nombre
	if len(body.Nombre) == 0 {
		response := map[string]string{
			"estado":  "error",
			"mensaje": "El nombre es requerido",
		}
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(response)
	}

	var result bson.M

	// convierto el id que viene en los params

	objId, _ := primitive.ObjectIDFromHex(c.Param("id"))

	if err := database.CategoriaCollection.FindOne(context.TODO(), bson.D{{Key: "_id", Value: objId}}).Decode(&result); err != nil {
		response := map[string]string{
			"estado":  "error",
			"mensaje": "No se encontro la categoria",
		}
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(response)
	}

	// editamos el registro
	// creo la variable
	registro := make(map[string]interface{})

	registro["nombre"] = body.Nombre
	registro["slug"] = slug.Make(body.Nombre)

	// actualizo el registro
	updateString := bson.M{
		"$set": registro,
	}
	database.CategoriaCollection.UpdateOne(context.TODO(), bson.D{{Key: "_id", Value: objId}}, updateString)

	// retorno respuesta
	respuesta := map[string]string{
		"estado":  "ok",
		"mensaje": "Categoria actualizada con exito",
	}

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c.Response().WriteHeader(http.StatusCreated)

	return json.NewEncoder(c.Response()).Encode(respuesta)

}

func Categoria_delete(c echo.Context) error {

	// convierto el id que viene en los params

	objId, _ := primitive.ObjectIDFromHex(c.Param("id"))

	// busco el registro
	var result bson.M

	if err := database.CategoriaCollection.FindOne(context.TODO(), bson.D{{Key: "_id", Value: objId}}).Decode(&result); err != nil {
		response := map[string]string{
			"estado":  "error",
			"mensaje": "No se encontro la categoria",
		}
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(response)
	}

	// elimino el registro
	database.CategoriaCollection.DeleteOne(context.TODO(), bson.D{{Key: "_id", Value: objId}})

	// retorno respuesta
	respuesta := map[string]string{
		"estado":  "ok",
		"mensaje": "Categoria eliminada con exito",
	}

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c.Response().WriteHeader(http.StatusCreated)

	return json.NewEncoder(c.Response()).Encode(respuesta)

}
