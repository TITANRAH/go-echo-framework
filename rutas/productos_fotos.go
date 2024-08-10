package rutas

import (
	"context"
	"echo-framework/database"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ProductoFotos_upload(c echo.Context) error {

	fmt.Println("entro a subir foto")

	file, err := c.FormFile("foto")
	if err != nil {
		return err
	}

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	// nombrado del archivo
	var extension = strings.Split(file.Filename, ".")[1]
	time := strings.Split(time.Now().String(), " ")

	foto := string(time[4][6:14]) + "." + extension
	var archivo string = "public/uploads/productos/" + foto

	// destino
	dst, err := os.Create(archivo)
	if err != nil {
		fmt.Println("entro al  error de creacion" + archivo)

		return err
	}

	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		fmt.Println("entro al segundo error")

		return err
	}

	// aca deberiamos gaurdar en bd
	var result bson.M

	// tomamos el id
	objId, _ := primitive.ObjectIDFromHex(c.Param("id"))

	fmt.Println("sigue")

	// preguntamos si el producto existearchivo
	if err := database.ProductosCollection.FindOne(context.TODO(), bson.D{{Key: "_id", Value: objId}}).Decode(&result); err != nil {

		fmt.Println("entro al error")
		response := map[string]string{
			"estado":  "error",
			"mensaje": "Producto no encontrado",
		}
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(response)
	}

	registro := bson.D{{Key: "nombre", Value: foto}, {Key: "producto_id", Value: objId}}
	database.ProductosFotosCollection.InsertOne(context.TODO(), registro)

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	response := map[string]string{
		"estado":  "ok",
		"mensaje": "Guardado correctamente",
	}

	c.Response().WriteHeader(http.StatusCreated)
	return json.NewEncoder(c.Response()).Encode(response)

}

func ProductoFotosGet_con_parametros(c echo.Context) error {

	fmt.Println("entro a subir foto")
	var producto bson.M
	id := c.Param("id")
	objId, _ := primitive.ObjectIDFromHex(id)

	if err := database.ProductosCollection.FindOne(context.TODO(), bson.D{{Key: "_id", Value: objId}}).Decode(&producto); err != nil {
		response := map[string]string{
			"estado":  "error",
			"mensaje": "Producto no encontrado",
		}
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(response)
	}

	cursor, err := database.ProductosFotosCollection.Find(context.TODO(), bson.D{{Key: "producto_id", Value: objId}})

	if err != nil {
		panic(err)
	}

	var results []bson.M

	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}
	fmt.Println("llego aca")
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(c.Response()).Encode(results)
}

func ProductoFotos_delete(c echo.Context) error {

	fmt.Println("entro a eliminar foto")
	var producto bson.M
	id := c.Param("id")
	objId, _ := primitive.ObjectIDFromHex(id)

	if err := database.ProductosFotosCollection.FindOne(context.TODO(), bson.D{{Key: "_id", Value: objId}}).Decode(&producto); err != nil {
		response := map[string]string{
			"estado":  "error",
			"mensaje": "Foto no encontrada",
		}
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(response)
	}

	borrar := "public/uploads/productos/" + producto["nombre"].(string)
	e := os.Remove(borrar)

	if e != nil {
		log.Fatal(e)
	}
	
	database.ProductosFotosCollection.DeleteOne(context.TODO(), bson.D{{Key: "_id", Value: objId}})

	results := map[string]string{
		"estado":  "ok",
		"mensaje": "Foto eliminada correctamente",
	}

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(c.Response()).Encode(results)
}
