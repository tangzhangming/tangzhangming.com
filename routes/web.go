package routes

import (
	"github.com/gin-gonic/gin"
	"tangzhangming.com/controller"
)

func Web(r *gin.Engine) {
	r.GET("/", controller.Index)
	r.GET("options", controller.Config)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
}
