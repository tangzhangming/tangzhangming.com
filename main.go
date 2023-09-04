package main

import (
	"flag"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"tangzhangming.com/initialize"
	"tangzhangming.com/internal/controller"
	"tangzhangming.com/internal/crontab"

	// "tangzhangming.com/internal/crontab"
	"tangzhangming.com/internal/config"
	// "tangzhangming.com/internal/pkg/database"
)

func main() {
	configFile := flag.String("config", "./config.yaml", "配置文件")
	flag.Parse()

	initialize.ConfigInit(configFile)

	initialize.LoggerInit()

	initialize.RedisInit(true)

	// database.SetConn()

	crontab.Task()

	HttpServer()

}

func HttpServer() {
	fmt.Println("\n -------------------- HTTP --------------------")
	fmt.Printf("[%s] 系统初始检测完成，正在启动HTTP服务... \n", config.Conf.Name)

	srv := gin.Default()
	srv.Use(ValidatorMiddleware())
	srv.LoadHTMLGlob("./web/template/*")
	srv.Static("/static", "./web/static")
	srv.Static("/upload", "./web/upload")
	srv.StaticFile("/favicon.ico", "./web/static/favicon-32x32.png")

	//注册业务路由并且启动
	controller.Routes(srv)
	srv.Run(config.Conf.Host + ":" + strconv.Itoa(config.Conf.Port)) //启动HTTP服务
}

func ValidatorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		validate := validator.New()

		// 将验证器存储到 Gin 上下文中
		c.Set("vd", validate)

		c.Next()
	}
}
