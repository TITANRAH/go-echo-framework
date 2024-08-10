package rutas

import (
	"context"
	"echo-framework/database"
	"echo-framework/dto"
	"encoding/json"
	"net/http"

	"github.com/gosimple/slug"
	echo "github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	// "go.mongodb.org/mongo-driver/mongo/options"
)

func Producto_post(c echo.Context) error {
	var body dto.ProductoDto
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
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(response)
	}

	// guardamos en la bd
	CategoriaID, _ := primitive.ObjectIDFromHex(body.CategoriaID)
	registro := bson.D{
		{Key: "nombre", Value: body.Nombre},
		{Key: "slug", Value: slug.Make(body.Nombre)},
		{Key: "stock", Value: body.Stock},
		{Key: "precio", Value: body.Precio},
		{Key: "descripcion", Value: body.Descripcion},
		{Key: "categoria_id", Value: CategoriaID},
	}

	database.ProductosCollection.InsertOne(context.TODO(), registro)

	respuesta := map[string]string{
		"estado":  "ok",
		"mensaje": "Product creado con exito",
	}

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c.Response().WriteHeader(http.StatusCreated)

	return json.NewEncoder(c.Response()).Encode(respuesta)

}

func Producto_get(c echo.Context) error {

	// findOptions := options.Find()

	// LLAMO A LA COLECCION  Y APLICO METODO FIND DE MONGO DRIVER Y COMO TERCER ARGUMENTOS PUEDO PONER CONDIDCIONES
	// EN LA BUSQUEDA SI QUISIERA
	// cursor, err := database.ProductosCollection.Find(context.TODO(), bson.D{}, findOptions.SetSort(
	// 	bson.D{{Key: "_id", Value: -1}},
	// ))

	// if err != nil {
	// 	panic(err)
	// }

	// relacion entre producto y categria
	pipeline := []bson.M{
		{"$match": bson.M{}},
		{"$lookup": bson.M{"from": "categorias", "localField": "categoria_id", "foreignField": "_id", "as": "categoria"}},
		{"$sort": bson.M{"_id": -1}},
	}
	cursor, err := database.ProductosCollection.Aggregate(context.TODO(), pipeline)
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

func Producto_get_con_parametro(c echo.Context) error {

	objId, _ := primitive.ObjectIDFromHex(c.Param("id"))

	pipeline := []bson.M{
		{"$match": bson.M{"_id": objId}},
		{"$lookup": bson.M{"from": "categorias", "localField": "categoria_id", "foreignField": "_id", "as": "categoria"}},
		{"$sort": bson.M{"_id": -1}},
	}
	cursor, err := database.ProductosCollection.Aggregate(context.TODO(), pipeline)
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
	return json.NewEncoder(c.Response()).Encode(results[0])

}

func Producto_put(c echo.Context) error {

	// declaro lo que recibire
	var body dto.ProductoDto

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

	if err := database.ProductosCollection.FindOne(context.TODO(), bson.D{{Key: "_id", Value: objId}}).Decode(&result); err != nil {
		response := map[string]string{
			"estado":  "error",
			"mensaje": "No se encontro el producto",
		}
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(response)
	}

	categoria_id, _ := primitive.ObjectIDFromHex(body.CategoriaID)

	// editamos el registro
	// creo la variable
	registro := make(map[string]interface{})

	registro["nombre"] = body.Nombre
	registro["slug"] = slug.Make(body.Nombre)
	registro["stock"] = body.Stock
	registro["precio"] = body.Precio
	registro["descripcion"] = body.Descripcion
	registro["categoria_id"] = categoria_id

	// actualizo el registro
	updateString := bson.M{
		"$set": registro,
	}
	database.ProductosCollection.UpdateOne(context.TODO(), bson.D{{Key: "_id", Value: objId}}, updateString)

	// retorno respuesta
	respuesta := map[string]string{
		"estado":  "ok",
		"mensaje": "Producto actualizado con exito",
	}

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c.Response().WriteHeader(http.StatusCreated)

	return json.NewEncoder(c.Response()).Encode(respuesta)

}

func Producto_delete(c echo.Context) error {
	
	// convierto el id que viene en los params

	objId, _ := primitive.ObjectIDFromHex(c.Param("id"))

	// busco el registro
	var result bson.M

	if err := database.ProductosCollection.FindOne(context.TODO(), bson.D{{Key: "_id", Value: objId}}).Decode(&result); err != nil {
		response := map[string]string{
			"estado":  "error",
			"mensaje": "No se encontro el producto",
		}
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(response)
	}

	// elimino el registro
	database.ProductosCollection.DeleteOne(context.TODO(), bson.D{{Key: "_id", Value: objId}})

	// retorno respuesta
	respuesta := map[string]string{
		"estado":  "ok",
		"mensaje": "Producto eliminado con exito",
	}

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c.Response().WriteHeader(http.StatusCreated)

	return json.NewEncoder(c.Response()).Encode(respuesta)
}
