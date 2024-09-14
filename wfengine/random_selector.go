package wfengine

import (
	"fmt"
	"github.com/didi/tg-flow/model"
	"hash/crc32"
	"strconv"
	"time"
)

type RandomSelector struct {
	FlowSelector
}

// 在线随机分流
func (r *RandomSelector) SelectWorkflowId(sc *model.StrategyContext, sceneModule *SceneModule) (int64, string, error) {
	//1、根据用户id，算出一个0-99的数字
	slotId := getSlotId(sc.UserId, sceneModule.BucketType)
	//2、在对应的维度id内，根据slotId，选择对应的workflow
	workflowId, ok := sceneModule.SlotMap[slotId]
	if !ok {
		return -1, "", fmt.Errorf("no workflowId found, UserId=%v, BucketType=%v", sc.UserId, sceneModule.BucketType)
	}

	return workflowId, strconv.Itoa(slotId), nil
}

func getSlotId(str string, bucketType int) int {
	if bucketType == 0 {
		return time.Now().Nanosecond() % 100
	}

	v := int(crc32.ChecksumIEEE([]byte(str)))
	if v < 0 {
		v = -v
	}
	return v % 100
}
