package redis

import (
	"fmt"
	"sync"

	"github.com/redis/go-redis/v9"
	"tangzhangming.com/internal/config"
)

var connections map[string]*redis.Client

var lock sync.Mutex

func Conn() *redis.Client {
	return connections["default"]
}

func Connection(name string) *redis.Client {
	return connections[name]
}

func SetConn(name string, conf *config.RedisConf) (*redis.Client, error) {
	lock.Lock()
	defer lock.Unlock()

	if connections == nil {
		connections = make(map[string]*redis.Client)
	}

	//如果该连接名已经存在配置 关闭旧连接
	if _, ok := connections[name]; ok == true {
		connections[name].Close()
		connections[name] = nil
	}

	connections[name] = redis.NewClient(&redis.Options{
		Addr:     conf.Addr,
		Username: conf.Username,
		Password: conf.Password,
		DB:       conf.DB,
	})

	return connections[name], nil
}

func DeleteConn(name string) {
	conn, ok := connections[name]
	if ok {
		conn.Close()
		delete(connections, name)
	}
}

func Debug() {
	var connName string
	for name, _ := range connections {
		connName += name + " "
	}
	fmt.Printf("当前有连接: %s \n", connName)
}
