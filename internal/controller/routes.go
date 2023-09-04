package controller

import (
	"github.com/gin-gonic/gin"
)

func Routes(r *gin.Engine) {

	//基础路由
	r.GET("/", Index)
	r.GET("options", Config)
	r.GET("/ping", ping)
	r.GET("/crontab", crontab_list)

	//推特视频
	r.GET("/video/:video_id", video_detail)
	r.GET("/video/ajax/:video_id", get_video_url)

	//文件上传
	r.GET("/upload", upload)
	r.POST("/upload", uploadSave)

	//Api路由
	v1 := r.Group("/api")
	{
		v1.GET("/user/create", User_create)
		v1.GET("/user/:id", User_view)

	}
}
