package controller

import (
	"github.com/gin-gonic/gin"
)

func Routes(r *gin.Engine) {

	//基础路由
	r.GET("/", Index)
	r.GET("options", Config)
	r.GET("/ping", ping)

	//Api路由
	v1 := r.Group("/api")
	{
		v1.GET("/user/create", User_create)
		v1.GET("/user/:id", User_view)

	}
}
