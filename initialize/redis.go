package initialize

import (
	"context"
	"fmt"

	"github.com/spf13/viper"
	"tangzhangming.com/internal/config"
	"tangzhangming.com/internal/pkg/redis"
)

func RedisInit(isFirst bool) {
	if isFirst == true {
		fmt.Println("\n -------------------- REDIS --------------------")
		fmt.Printf("[%s] read %d redis connection configure\n", config.Conf.Name, len(config.Conf.RedisConf))
	}

	//如果有redis配置的情况下 必须配置一个名为default的redis连接
	if len(config.Conf.RedisConf) > 0 {
		_, ok := config.Conf.RedisConf["default"]
		if ok == false {
			fmt.Println("未检测到名为default的redis连接配置，redis配置将不生效")
			return
		}
	}

	for name, conf := range config.Conf.RedisConf {
		//viper配置文件更新到map时，viper只处理增量，需要特殊处理一下
		_, ok := viper.GetStringMap("RedisConf")[name]
		if ok == false {
			fmt.Printf("close redis connection %s \n", name)
			redis.DeleteConn(name)
			delete(config.Conf.RedisConf, name)
			continue
		}

		//设置连接
		connection, _ := redis.SetConn(name, conf)

		//Ping 测试
		show_name := fmt.Sprintf("%s(%s)", name, conf.Addr)
		if connection.Ping(context.Background()).Err() != nil {
			fmt.Printf("[%s] redis connection %s ping is fail \n", config.Conf.Name, show_name)
			break
		}

		fmt.Printf("[%s] redis connection %s ping is success \n", config.Conf.Name, show_name)
	}

	redis.Debug()
}
