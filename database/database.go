package database

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ClienteMongo = ConectarDB()
var MongoBD = "Cluster0"
var CategoriaCollection = ClienteMongo.Database(MongoBD).Collection("categorias")
var ProductosCollection = ClienteMongo.Database(MongoBD).Collection("productos")
var ProductosFotosCollection = ClienteMongo.Database(MongoBD).Collection("productos_fotos")
var UsuariosCollection = ClienteMongo.Database(MongoBD).Collection("usuarios")
// EDITAR IP EN LA CREACION DE LA BD
var clientOpts = options.Client().ApplyURI("mongodb+srv://granrah:biorapfia1@cluster0.u372s.mongodb.net/?retryWrites=true&w=majority&appName=" + MongoBD)

func ConectarDB() *mongo.Client {
	client, err := mongo.Connect(context.TODO(), clientOpts)
	if err != nil {
		log.Fatal(err.Error())
		return client
	}
	// no neceistamos hacer nada mas que un ping por eso nil
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err.Error())
		return client
	}

	log.Println("Conexion exitosa a la base de datos")
	return client

}

func CheckConnection() int {
	err := ClienteMongo.Ping(context.TODO(), nil)
	if err != nil {
		return 0
	}
	return 1
}
