/**
Description : flow splitter
Author		: dayunzhangyunfeng@didiglobal.com
Date		: 2021-05-14
*/
package wfengine

import (
	"fmt"
	"github.com/didi/tg-flow/tg-core/model"
	"hash/crc32"
	"time"
)

//在线随机分流
func FlowByOnlineRandom(sc *model.StrategyContext, sceneModule *SceneModule) (int64, int) {
	//1、根据用户id，算出一个0-99的数字
	slotId := getSlotId(sc.UserId, sceneModule.BucketType)

	//2、在对应的维度id内，根据slotId，选择对应的workflow
	return sceneModule.SlotMap[slotId], slotId
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

// FlowByApollo apollo分流
func FlowByApollo(sc *model.StrategyContext, sceneModule *SceneModule) (int64, string, error) {
	var err error
	if sc.ApolloInfo == nil {
		return 0, "", fmt.Errorf("apollo info not initialized")
	}

	// 优先采用 ApolloInfo 中设置的分流实验名称，没有的话采用场景中配置的分流实验名称
	if sc.ApolloInfo.GetDispatchExperimentName() == "" {
		sc.ApolloInfo.SetDispatchExperimentName(sceneModule.DispatchExperimentName)
	}

	groupName, err := sc.ApolloInfo.GetDispatchGroupName()
	if err != nil {
		return 0, "", err
	}

	workflowId, err := sceneModule.GetWorkflowId(groupName)
	if err == nil {
		return workflowId, groupName, nil
	} else {
		return workflowId, groupName, err
	}
}
