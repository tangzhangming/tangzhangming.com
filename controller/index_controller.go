package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"tangzhangming.com/pkg/config"
	"tangzhangming.com/pkg/log"
	"tangzhangming.com/pkg/redis"
)

func Index(c *gin.Context) {

	rdb := redis.Conn()
	id, _ := rdb.Incr(c, "welcomen").Result()

	c.String(200, "app name is : "+config.Conf.Name+", N=:"+strconv.Itoa(int(id)))
}

func Config(c *gin.Context) {

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	log.Info("failed to fetch URL",
		// Structured context as strongly typed Field values.
		zap.String("url", "dejijdwie"),
		zap.Int("attempt", 3),
		// zap.Duration("backoff", time.Second),
	)

	log.Info("一条错误信息")

	c.JSON(200, config.Conf)
}
