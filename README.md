# Gin快速开发骨架

    Gin + Gorm


## Redis
    ```
    import "tangzhangming.com/pkg/redis"

    func FancName(){
        rdb := redis.Conn()
        view_count, _ := rdb.Incr(c, "view_count").Result()
    }
    ```

## 定时任务
    在项目 crontab 目录
    ```
    import "github.com/robfig/cron/v3"

    func Task() {
        //New一个秒级定时任务
        c := cron.New(cron.WithSeconds())

        //挂载你的自定义定时任务到这里
        c.AddFunc("*/10 * * * * *", myTask1)
        c.AddFunc("*/5 * * * *", myTask2)

        //启动定时任务
        c.Start()
    }

    func myTask1() {
        fmt.PrintIn("定时任务1执行了 \n")
    }

    func myTask2() {
        fmt.PrintIn("定时任务2执行了 \n")
    }
    ```


