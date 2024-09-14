/**
Description : loader of workflow config info from redis
Author		: dayunzhangyunfeng@didiglobal.com
Date		: 2021-07-12
*/

package wfengine

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/didi/tg-flow/common/tlog"
)

const (
	RedisKeySceneModule = "scene_module_app_"
	RedisKeyWorkflow    = "workflow_app_"
	RedisKeyVersion     = "version_app_"
)

//func GetLatestVersionFromRedis(appId int64) (string, error) {
//	return redis.Handler.Get(context.TODO(), fmt.Sprintf("%v%v", RedisKeyVersion, appId))
//}

func NewWorkflowEngineFromKV(moduleObj ModuleObjBase, sceneModuleMapString, workflowMapStr, version string) (*WorkflowEngine, error) {
	var sceneModuleMap map[int64]*SceneModule
	err := json.Unmarshal([]byte(sceneModuleMapString), &sceneModuleMap)
	if err != nil {
		return nil, fmt.Errorf("err:%v, sceneModule:%v", err, sceneModuleMapString)
	}
	/*for sceneId, sceneModule := range sceneModuleMap {
		if sceneModule.AppId != appId {
			delete(sceneModuleMap, sceneId)
		}
	}*/

	var workflowMap map[int64]*Workflow
	err = json.Unmarshal([]byte(workflowMapStr), &workflowMap)
	if err != nil {
		return nil, fmt.Errorf("err:%v, workflow:%v", err, workflowMapStr)
	}

	wfMap := make(map[int64]*Workflow)
	for workflowId, workflow := range workflowMap {
		if _, ok := sceneModuleMap[workflow.SceneId]; ok {
			workflow.FlowCharts, err = NewWorkflowChart(workflow.FlowChart)
			if err != nil {
				tlog.Handler.ErrorCount(context.TODO(), "NewWorkflowChart_err", fmt.Sprintf("wf:%v,err:%v", workflow, err))
				continue
			}
			wfMap[workflowId] = workflow
		}
	}

	return NewWorkflowEngine(sceneModuleMap, wfMap, version, moduleObj)
}
