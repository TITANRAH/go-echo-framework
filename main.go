package main

import (
	"echo-framework/database"
	"echo-framework/rutas"
	"os"

	"github.com/joho/godotenv"
	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var prefijo string = "/api/v1"

func main() {
	e := echo.New()

	// archivos estaticos
	// este codigo permite comunicar archivos estaticos con el exterior
	e.Static("/public", "public")


	// CONEXION A LA BD
	database.CheckConnection()

	e.GET(prefijo+"/saludar", rutas.Ejemplo_get)
	e.POST(prefijo+"/saludar", rutas.Ejemplo_post)
	e.PUT(prefijo+"/saludar:id", rutas.Ejemplo_put)
	e.DELETE(prefijo+"/saludar/:id", rutas.Ejemplo_delete)
	e.GET(prefijo+"/saludar/:id", rutas.Ejemplo_get_con_parmetros)
	e.GET(prefijo+"/query-string", rutas.Ejemplo_query_string)
	e.POST(prefijo+"/upload", rutas.Ejemplo_upload)

	// Categorias
	e.POST(prefijo+"/categorias", rutas.Categoria_post)
	e.GET(prefijo+"/categorias", rutas.Categoria_get)
	e.GET(prefijo+"/categorias/:id", rutas.Categoria_get_con_parametro)
	e.PUT(prefijo+"/categorias/:id", rutas.Categoria_put)
	e.DELETE(prefijo+"/categorias/:id", rutas.Categoria_delete)

	// Productos
	e.POST(prefijo+"/productos", rutas.Producto_post)
	e.PUT(prefijo+"/productos/:id", rutas.Producto_put)
	e.GET(prefijo+"/productos", rutas.Producto_get)
	e.GET(prefijo+"/productos/:id", rutas.Producto_get_con_parametro)
	e.DELETE(prefijo+"/productos/:id", rutas.Producto_delete)

	//Productos_fotos

	e.POST(prefijo+"/productos_fotos/:id", rutas.ProductoFotos_upload)
	e.GET(prefijo+"/productos_fotos/:id", rutas.ProductoFotosGet_con_parametros)
	e.DELETE(prefijo+"/productos_fotos/:id", rutas.ProductoFotos_delete)

	// Seguridad

	e.POST(prefijo+"/seguridad/registro", rutas.Seguridad_registro)
	e.POST(prefijo+"/seguridad/login", rutas.Seguridad_login)
	e.GET(prefijo+"/seguridad/protegida", rutas.Seguridad_protegida)



	errorVariables := godotenv.Load()
	if errorVariables != nil {
		panic(errorVariables)
	}

	//CORS

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
	  AllowOrigins: []string{"http//localhost", "https://labstack.net"},
	  AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}
