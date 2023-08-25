package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"tangzhangming.com/pkg/config"
)

var connection *redis.Client

func Conn() *redis.Client {
	return connection
}

func SetConn() {
	fmt.Println("\n -------------------- REDIS --------------------")

	connection = redis.NewClient(&redis.Options{
		Addr:     config.Conf.RedisConf.Addr,
		Username: config.Conf.RedisConf.Username,
		Password: config.Conf.RedisConf.Password,
		DB:       config.Conf.RedisConf.DB,
	})

	fmt.Printf("[%s] Redis %s 连接成功 \n", config.Conf.Name, config.Conf.RedisConf.Addr)
	if connection.Ping(context.Background()).Err() != nil {
		fmt.Printf("[%s] Redis %s Ping 测试失败 \n", config.Conf.Name, config.Conf.RedisConf.Addr)
	} else {
		fmt.Printf("[%s] Redis %s Ping 测试成功 \n", config.Conf.Name, config.Conf.RedisConf.Addr)
	}

}
