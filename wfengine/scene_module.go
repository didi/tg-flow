/**
Description : config info of scene and module
Author		: dayunzhangyunfeng@didiglobal.com
Date		: 2021-05-14
*/

package wfengine

import (
	"fmt"
	"time"
)

const (
	groupWorkflowKeySpliter = "."	
)

type SceneModule struct {
	Id			int64					`json:"id"`
	Name		string					`json:"name"`
	AppId		int64					`json:"appid"`
	BucketType	int						`json:"bucket_type"`
	SlotMap 	map[int]int64 			`json:"slots"`
	UpdateTime	time.Time               `json:"update_time"`
	FlowType	int                     `json:"flow_type"`
	GroupWorkflowMap map[string]int64	`json:"group_workflows"`
	DefaultWorkflowId int64             `json:"default_workflow_id"`
	DispatchExperimentName string       `json:"dispatch_experiment_name"`
}

func (this *SceneModule) GetWorkflowId(groupName string) (int64, error) {
	key := fmt.Sprintf("%v%v%v", this.DispatchExperimentName, groupWorkflowKeySpliter, groupName)
	// TODO: 升级直接去掉带实验名的逻辑
	if workflowId, ok := this.GroupWorkflowMap[key]; ok {
		return workflowId, nil
	} else if workflowId, ok := this.GroupWorkflowMap[groupName]; ok {
		return workflowId, nil
	} else {
		return workflowId, fmt.Errorf("no workflow_id found with key of:%v", key)
	}
}