package redis

import (
	"context"
	"fmt"
	"sync"

	"github.com/redis/go-redis/v9"
	"tangzhangming.com/internal/pkg/config"
)

var connection *redis.Client

var lock sync.RWMutex

func Conn() *redis.Client {
	lock.RLock()
	c := connection
	lock.RUnlock()

	return c
}

func SetConn() {
	defer lock.Unlock()
	lock.Lock()

	fmt.Println("\n -------------------- REDIS --------------------")
	connection = redis.NewClient(&redis.Options{
		Addr:     config.Conf.RedisConf.Addr,
		Username: config.Conf.RedisConf.Username,
		Password: config.Conf.RedisConf.Password,
		DB:       config.Conf.RedisConf.DB,
	})

	fmt.Printf("[%s] Redis %s 配置成功 \n", config.Conf.Name, config.Conf.RedisConf.Addr)
	if connection.Ping(context.Background()).Err() != nil {
		fmt.Printf("[%s] Redis %s Ping 测试失败 \n", config.Conf.Name, config.Conf.RedisConf.Addr)
	} else {
		fmt.Printf("[%s] Redis %s Ping 测试成功 \n", config.Conf.Name, config.Conf.RedisConf.Addr)
	}

}
