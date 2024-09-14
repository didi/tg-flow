package model

import (
	"time"
	"fmt"
)

const (
	groupWorkflowKeySpliter = "."	
)

type SceneModule struct {
	Id         int64  `json:"id"`
	Name       string `json:"name"`
	AppId      int64  `json:"appid"`
	BucketType int    `json:"bucketType"`
	//第一层key：维度id；第二层key：0-99的数字；value：workflowid
	DimensionMap map[int64]map[int]int64 `json:"dimensionMap"`
	/**第一层key：moduleName，如：“bonus.BonusStrategyAwareStrengthen”
	  第二层key：参数组合，如：“[p1=2&p2=9]”
	  第三层key：参数名称，如：“p1”
	     value：每个参数对应的值，支持string类型
	*/
	ModuleInfoMap map[string]map[string]map[string]string `json:"moduleInfoMap"`
	UpdateTime    time.Time                               `json:"updateTime"`
	FlowType      int                                     `json:"flow_type"`
	//第一层key：维度id；第二层key：manualSlotId；value：workflowid
	ManualSlotIdsMap map[int64]map[int64]int64 `json:"manualSlotIdsMap"`
	
	GroupWorkflowMap map[string]int64			`json:"groupWorkflowMap"`
}

func (this *SceneModule) GetWorkflowId(dispatchExperimentName, groupName string) (int64, error) {
	key := fmt.Sprintf("%v%v%v", dispatchExperimentName, groupWorkflowKeySpliter, groupName)
	if workflowId, ok := this.GroupWorkflowMap[key]; ok {
		return workflowId, nil
	}else{
		return workflowId, fmt.Errorf("no workflow_id found with key of:%v", key)
	}
}