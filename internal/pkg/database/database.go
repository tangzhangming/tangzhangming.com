package database

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"tangzhangming.com/internal/pkg/config"
)

var DB *gorm.DB

func SetConn() {
	fmt.Println("\n -------------------- MYSQL --------------------")

	dsn := "root:root@tcp(localhost:3306)/tang?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err == nil {
		fmt.Printf("[%s] Mysql 连接成功 \n", config.Conf.Name)
	} else {
		fmt.Printf("[%s] Mysql 连接失败: %s \n", config.Conf.Name, err)
	}

	DB = db
}
