package api

import "github.com/gin-gonic/gin"

func Route(r *gin.Engine) {
	v1 := r.Group("/api")
	{
		v1.GET("/user/create", user_create)
		v1.GET("/user/:id", user_view)

	}
}
