package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"tangzhangming.com/api"
	"tangzhangming.com/crontab"
	"tangzhangming.com/pkg/config"
	"tangzhangming.com/pkg/database"
	"tangzhangming.com/pkg/log"
	"tangzhangming.com/pkg/redis"
	"tangzhangming.com/routes"
)

func main() {

	name := flag.String("action", "start", "命令")
	flag.Parse()

	if *name == "stop" {
		//停机
		StopHttpServer()
		return

	} else if *name == "restart" {
		//重启
		fmt.Println("重启停机")
		return
	} else if *name == "d" {
		shjc_srv()
		// return
	} else if *name == "start" {
		start_srv()
		// return
	} else {
		fmt.Println("未知的启动方式")

	}

	// fmt.Println(*name)
	// return

}

func start_srv() {
	config.Load()

	log.InitLogger()

	redis.SetConn()

	database.SetConn()

	crontab.Task()

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

func StopHttpServer() {
	fmt.Println("停机")

	pfname := "./daemon.pid"
	data, err := ioutil.ReadFile(pfname)
	if err != nil {
		fmt.Printf("守护进程启动失败, 错误信息：%s", err)
		return
	}

	fmt.Println(string(data))
}

// 守护进程方式启动
func shjc_srv() {
	fmt.Println("守护进程方式启动")

	pfname := "./logs/daemon.pid"

	//判断是否已经有进程
	// data, err := ioutil.ReadFile(pfname)
	// if err != nil {
	// 	fmt.Printf("守护进程启动失败, 错误信息：%s", err)
	// 	return
	// }

	pf, _ := os.Create(pfname)
	pf.Write([]byte(strconv.Itoa(os.Getpid()))) //把当前进程pid写入文件

	//守护进程方式启动
	// start_srv()
	// os.Exit(0)
}
