package store

import (
	"context"
	"encoding/json"
	"github.com/didi/tg-flow/tg-core/common/redis"
	"github.com/didi/tg-flow/tg-core/common/tlog"
	"github.com/didi/tg-flow/tg-core/consts"
)

var MachineIpMap map[string][]string = make(map[string][]string)

func LoadMachineIp(appId int) error {
	machineIpStr := redis.GetRedis(context.TODO(), nil, redis.Handler, consts.RedisKeyMachineIp, "task.LoadMachineIp")

	if machineIpStr == "" {
		return nil
	}

	var machineIpMap map[int]map[string][]string
	err := json.Unmarshal([]byte(machineIpStr), &machineIpMap)
	if err != nil {
		return err
	}

	//只获取本系统下的机器ip
	tempMachineIpMap := make(map[string][]string)
	//appId := *(*int)(unsafe.Pointer(&module.AppId))

	if tempMap := machineIpMap[appId]; tempMap != nil {
		for machineRoom, ipList := range tempMap {
			tempMachineIpMap[machineRoom] = ipList
		}
	}

	if len(tempMachineIpMap) <= 0 {
		return nil
	}

	MachineIpMap = tempMachineIpMap
	tlog.Handler.Infof(context.TODO(), consts.DLTagCronTask, "etype=task.LoadMachineIp||machineIpMap=%v", MachineIpMap)
	return nil
}
