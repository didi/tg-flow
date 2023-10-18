package crontask

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/didi/tg-flow/tg-core/common/tlog"
	"github.com/didi/tg-flow/tg-core/consts"
	"github.com/robfig/cron"
	"math"
	"math/big"
	"reflect"
	"strconv"
	"strings"
//	"time"
)

//定时任务管理实例
var CronTask *cron.Cron = cron.New()

//定时任务管道
var ChanMap map[string]chan interface{} = make(map[string]chan interface{})

/**注册定时任务
   taskTime:定时任务的周期
        job:具体的定时任务对象
         ch:job任务中的定时管道
       flag:是否随系统启动时，执行一次该任务
**/
func RegisterTask(taskTime string, job cron.Job, flag bool) {
	//定时任务管道初始化
	name := reflect.TypeOf(job)
	ChanMap[name.Name()] = make(chan interface{}, 1)

	//是否随系统启动执行
	if flag {
		job.Run()
	}

	//如果是秒以上的定时周期，则给一个随机的秒，以避免同台机器并发执行任务，给相应机器造成并发压力
	newTaskTime := ""
	tempTimes := strings.Split(taskTime, " ")
	if tempTimes[0] == "0" || tempTimes[0] == "*" {
		for i := 1; i < len(tempTimes); i++ {
			newTaskTime += " " + tempTimes[i]
		}
		newTaskTime = strconv.FormatInt(RangeRand(0, 59), 10) + newTaskTime
//		time.Sleep(1 * time.Second)
//		newTaskTime = strconv.Itoa(time.Now().Second()) + newTaskTime
	} else {
		newTaskTime = taskTime
	}
	
	//注册任务
	CronTask.AddJob(newTaskTime, job)
}

/**
   启动定时任务
**/
func TaskStart() {
	defer recoverTaskStart()
	CronTask.Start()
}

//定时任务运行异常,捕获异常后，需要释放该任务的任务管道
func RecoverTaskRun(taskName string, ctx context.Context) {
	if err := recover(); err != nil {
		content := fmt.Sprintf("taskName=%v", taskName)
		tlog.LogError(ctx, nil, consts.DLTagCronTask, "crontask.recoverTaskRun", content, fmt.Errorf("%v", err))
	}
}

//所有定时任务启动失败
func recoverTaskStart() {
	if err := recover(); err != nil {
		tlog.LogError(context.TODO(), nil, consts.DLTagCronTask, "crontask.recoverTaskStart", "all task start fail", fmt.Errorf("%v", err))
	}
}

// 生成区间[-m, n]的安全随机数
func RangeRand(min, max int64) int64 {
	if min > max {
		panic("the min is greater than max!")
	}

	if min < 0 {
		f64Min := math.Abs(float64(min))
		i64Min := int64(f64Min)
		result, _ := rand.Int(rand.Reader, big.NewInt(max+1+i64Min))

		return result.Int64() - i64Min
	} else {
		result, _ := rand.Int(rand.Reader, big.NewInt(max-min+1))
		return min + result.Int64()
	}
}
