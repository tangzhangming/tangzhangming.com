package crontab

import (
	"fmt"

	"github.com/robfig/cron/v3"
)

var pull_redis_cache_num int = 0

func Task() {
	fmt.Println("-------------------- CRONTAB --------------------")
	return
	//New一个秒级定时任务
	c := cron.New(cron.WithSeconds())

	//挂载你的定时任务到这里
	c.AddFunc("*/10 * * * * *", pull_redis_cache)

	//启动定时任务
	c.Start()
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
