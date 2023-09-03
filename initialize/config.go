package initialize

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"tangzhangming.com/internal/config"
)

func ConfigInit(f *string) {
	fmt.Println("-------------------- OPTION --------------------")
	// confFileName := "./config.yaml"
	confFileName := string(*f)
	viper.SetConfigFile(confFileName)

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	} else {
		fmt.Printf("[%s] 配置文件 '%s' 读取成功 \n", viper.Get("name"), confFileName)
	}

	if err := viper.Unmarshal(config.Conf); err != nil {
		panic(fmt.Errorf("unmarshal conf failed, err:%s \n", err))
	} else {
		fmt.Printf("[%s] 配置系统加载完成 \n", config.Conf.Name)
	}

	viper.WatchConfig()
	viper.OnConfigChange(onConfigChange)
}

// 配置文件重载 注意这破玩意可能会执行多次
func onConfigChange(e fsnotify.Event) {
	fmt.Println("检测到配置文件更新:", e.Name)

	if err := viper.Unmarshal(config.Conf); err != nil {
		panic(fmt.Errorf("配置文件更新同步到系统发生失败, err:%s \n", err))
	}

	reload_redis()
}

// 重载redis
func reload_redis() {
	fmt.Println("----------------- 重载 REDIS 配置 -----------------")
	RedisInit(false)
}
