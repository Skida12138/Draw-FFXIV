package routes

import (
	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func init() {
	router = gin.Default()
	registerGlobalMiddlewares(router)
	registerRoomsRoutes(router)
}

// Start : start routing
func Start() {
	router.Run(":8080")
}
