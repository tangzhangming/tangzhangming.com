package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-module/carbon"
	"tangzhangming.com/internal/config"
	"tangzhangming.com/internal/pkg/log"
	"tangzhangming.com/internal/pkg/redis"
)

func Index(c *gin.Context) {

	id, _ := redis.Connection("cache").Incr(c, "welcomenjsq").Result()

	c.HTML(http.StatusOK, "index.html", gin.H{
		"appName":  config.Conf.Name,
		"viewnum":  id,
		"datetime": carbon.Now().String(),
	})
}

func Config(c *gin.Context) {

	log.Info("failed to fetch URL",
		log.String("url", "http://www.qq.com"),
		log.Int("attempt", 3),
	)

	log.Error("一条错误信息")

	c.JSON(200, config.Conf)
}

func ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func crontab_list(c *gin.Context) {
	// var list []map[string]string

	// for id, xx := range crontab.Cron.Entries() {
	// 	list[id] = map[string]string{
	// 		"ID": strconv.Itoa(int(xx.ID)),
	// 	}
	// }

	// storage.Put("ddd", "sswsw")

	c.JSON(200, gin.H{
		"message": "crontab",
		// "data":    list,
	})
}
