package crontask
//
//import (
//	"context"
//	"github.com/didi/tg-flow/common/tlog"
//	"reflect"
//	"time"
//)
//
//const(
//	TagCronTask = " _cron_task"
//)
//// 定时任务接口
//type TaskInterface interface {
//	Run()
//}
//
///*
//********************参数含义***************
//
//	task:定时任务结构体
//	period:定时周期的单位:秒
//	flag:是否在系统启动前启动定时任务
//
//******************************************
//*/
//func StartTask(task TaskInterface, period time.Duration, flag bool) {
//	ctx := context.TODO()
//	taskName := reflect.TypeOf(task).String()
//
//	defer RecoverTaskRun(taskName, ctx)
//	if flag { //如果是需要在系统启动前启动定时任务，那么定时任务正式开始执行时间后延一个周期
//		time.Sleep(period * time.Second)
//	}
//
//	for {
//		tlog.Handler.Infof(ctx, TagCronTask, "taskName=%v is start.", taskName)
//		task.Run()
//		time.Sleep(period * time.Second)
//	}
//}
