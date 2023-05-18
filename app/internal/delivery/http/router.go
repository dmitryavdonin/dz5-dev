package delivery

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (d *Delivery) initRouter() *gin.Engine {

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowCredentials = true

	var router = gin.Default()

	router.Use(cors.New(corsConfig))
	router.GET("/", d.GetVersion)

	return router
}
