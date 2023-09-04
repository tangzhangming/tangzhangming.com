package crontab

import (
	"fmt"

	"github.com/robfig/cron/v3"
)

var Cron *cron.Cron

var pull_redis_cache_num int = 0

func Task() {
	fmt.Println("-------------------- CRONTAB --------------------")

	//New一个秒级定时任务
	Cron = cron.New(cron.WithSeconds())

	//挂载你的定时任务到这里
	Cron.AddFunc("*/30 * * * * *", pull_redis_cache)

	//启动定时任务
	Cron.Start()
}

/*
 * 每2分钟执行一次
 */
//@Scheduled(cron = "*/2 * * * *")
func pull_redis_cache() {
	pull_redis_cache_num++
	fmt.Printf("第 %d 次更新缓存\n", pull_redis_cache_num)
}

func PullNum() int {
	return pull_redis_cache_num
}
