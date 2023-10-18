package main

import (
	"github.com/didi/tg-flow/tg-core/common/httpserv"
	"log"
	"tg-service/conf"
	"tg-service/cron"
)

func main() {
	//初始化配置
	conf.InitConfig()

	//启动定时任务
	cron.StartCronTask()

	//启动服务监听
	httpserv.Handler.Run()

	log.Println("Finished!")
}
