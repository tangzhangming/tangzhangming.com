package config

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func Load() {
	fmt.Println("-------------------- OPTION --------------------")

	viper.SetConfigFile("./conf/conf.yaml")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	} else {
		fmt.Printf("[%s] 配置文件 conf/conf.yaml 读取成功 \n", Conf.Name)
	}

	if err := viper.Unmarshal(Conf); err != nil {
		panic(fmt.Errorf("unmarshal conf failed, err:%s \n", err))
	} else {
		fmt.Printf("[%s] 配置系统加载完成 \n", Conf.Name)
	}

	//监听配置文件的更新并且写入系统
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("检测到配置文件更新:", e.Name)

		if err := viper.Unmarshal(Conf); err != nil {
			panic(fmt.Errorf("配置文件更新同步到系统发生失败, err:%s \n", err))
		}
	})

}
