package main

import (
	"gin-mongo-api/configs"
	"gin-mongo-api/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	//inicializamos la configuracion por default para la declaracion de rutas
	router := gin.Default()

	//Corremos la conexion a la BD
	configs.ConnectDB()

	//Rutas declaradas
	routes.UserRoute(router)
	//Iniciamos el ruteo en el puerto declarado
	router.Run("localhost:6000")
}
