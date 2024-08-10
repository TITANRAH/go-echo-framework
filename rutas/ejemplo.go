package rutas

import (
	"echo-framework/dto"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

func Ejemplo_get(c echo.Context) error {

	response := map[string]string{
		"estado":        "ok",
		"mensaje":       "Hola Mundo desde ooh",
		"Authorization": c.Request().Header.Get("Authorization"),
	}

	// PUEDE SER ASI CON ECHO
	// existe statuds bad request que devuelve 400 u statusOk que devuekve 200
	// return c.JSON(http.StatusBadRequest,response)
	// return c.JSON(http.StatusBadRequest,response)

	// O ESTA FORMA
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)

	// para pasar el status code
	c.Response().WriteHeader(http.StatusBadRequest)
	return json.NewEncoder(c.Response()).Encode(response)
}

//	func Ejemplo_get(c echo.Context) error {
//		return c.String(http.StatusOK, "Hola Mundooooooo")
//	}
func Ejemplo_get_con_parmetros(c echo.Context) error {

	// este id coincide con el argumento pasado en la ruta

	id := c.Param("id")

	response := map[string]string{
		"estado":  "ok",
		"mensaje": "Hola Mundo desde Get PARAMETROS ID: " + id,
	}

	// PUEDE SER ASI CON ECHO
	// return c.JSON(200,response)

	// O ESTA FORMA
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	return json.NewEncoder(c.Response()).Encode(response)
	// return c.String(http.StatusOK, "Get parametro | id: "+id)
}
func Ejemplo_post(c echo.Context) error {

	var body dto.CategoriaDto
	if err := json.NewDecoder(c.Request().Body).Decode(&body); err != nil {
		response := map[string]string{
			"estado":  "error",
			"mensaje": "Ocurrio un error en el body",
		}
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(response)

	}

	response := map[string]string{
		"estado":  "ok",
		"mensaje": "Hola Mundo desde POST",
		"nombre":  body.Nombre,
	}

	// PUEDE SER ASI CON ECHO
	// return c.JSON(200,response)

	// O ESTA FORMA
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	return json.NewEncoder(c.Response()).Encode(response)
	// return c.String(http.StatusOK, "Metodo POST")
}
func Ejemplo_put(c echo.Context) error {
	id := c.Param("id")

	response := map[string]string{
		"estado":  "ok",
		"mensaje": "Hola Mundo desde PUT PARAMETROS ID: " + id,
	}

	// PUEDE SER ASI CON ECHO
	// return c.JSON(200,response)

	// O ESTA FORMA
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	return json.NewEncoder(c.Response()).Encode(response)
}
func Ejemplo_delete(c echo.Context) error {
	id := c.Param("id")

	response := map[string]string{
		"estado":  "ok",
		"mensaje": "Hola Mundo desde DELETE PARAMETROS ID: " + id,
	}

	// PUEDE SER ASI CON ECHO
	// return c.JSON(200,response)

	// O ESTA FORMA
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	return json.NewEncoder(c.Response()).Encode(response)
	// return c.String(http.StatusOK, "Metodo DELETE" + id)
}
func Ejemplo_query_string(c echo.Context) error {
	id := c.QueryParam("id")
	slug := c.QueryParam("slug")

	response := map[string]string{
		"estado":  "ok",
		"mensaje": "Hola Mundo desde Metodo queryparam" + " id: " + id + " slug: " + slug,
	}

	// PUEDE SER ASI CON ECHO
	// return c.JSON(200,response)

	// O ESTA FORMA
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	return json.NewEncoder(c.Response()).Encode(response)
	// return c.String(http.StatusOK, "Metodo queryparam" + " id: " + id + " slug: " + slug)
}

func Ejemplo_upload(c echo.Context) error {

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
	var archivo string = "public/uploads/fotos/" + foto

	// destino
	dst, err := os.Create(archivo)
	if err != nil {
		return err
	}

	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	// aca deberiamos gaurdar en bd


	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	response := map[string]string{
		"estado":  "ok",
		"mensaje": "Guardado correctamente",
		"foto":    foto,
	}

	c.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(c.Response()).Encode(response)

}
