/**
  author:dingyi
  time:2019-08-28
**/

package cron

import (
	"github.com/didi/tg-flow/tg-core/common/crontask"
)

const (
	//定时任务周期
	BASE_TASK_TIME = "0 0/2 * * * ?"
	HOUR_TASK_TIME = "0 0 0/4 * * ?"
)

func StartCronTask() {

	//启动所有定时任务
	crontask.TaskStart()
}
