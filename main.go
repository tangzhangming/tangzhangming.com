package main

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"tangzhangming.com/api"
	"tangzhangming.com/pkg/config"
	"tangzhangming.com/pkg/database"
	"tangzhangming.com/pkg/log"
	"tangzhangming.com/pkg/redis"
	"tangzhangming.com/routes"
)

func main() {

	config.Load()

	log.InitLogger()

	redis.SetConn()

	database.SetConn()

	HttpServer()
}

func HttpServer() {
	fmt.Println("\n -------------------- HTTP --------------------")
	fmt.Printf("[%s] 系统初始检测完成，正在启动HTTP服务... \n", config.Conf.Name)

	r := gin.Default()
	r.Use(ValidatorMiddleware())

	api.Route(r)
	routes.Web(r)

	//启动HTTP服务
	r.Run(config.Conf.Host + ":" + strconv.Itoa(config.Conf.Port))
}

func ValidatorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		validate := validator.New()

		// 将验证器存储到 Gin 上下文中
		c.Set("vd", validate)

		c.Next()
	}
}
