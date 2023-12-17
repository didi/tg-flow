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
	"github.com/didi/tg-flow/tg-core/common/redis"
	"github.com/didi/tg-flow/tg-core/common/tlog"
)

const (
	RedisKeySceneModule = "scene_module_app_"
	RedisKeyWorkflow    = "workflow_app_"
	RedisKeyVersion     = "version_app_"
)

func GetLatestVersionFromRedis(appId int64) (string, error) {
	return redis.Handler.Get(context.TODO(), fmt.Sprintf("%v%v", RedisKeyVersion, appId))
}

func NewWorkflowEngine(moduleObj ModuleObjBase, appId int64) (*WorkflowEngine, error) {
	//更新系统下全部场景的节点对象
	smMap, err := LoadSceneModuleMap(appId)
	if err != nil {
		return nil, err
	}

	//load workflow
	wfMap, err := LoadWorkflow(appId, smMap)
	if err != nil {
		return nil, err
	}

	//版本号不是必须，兼容吧
	version, err1 := GetLatestVersionFromRedis(appId)
	if err1 != nil {
		tlog.ErrorCount(context.TODO(), "GetLatestVersionFromRedis_err", fmt.Sprintf("appId=%v, err=%v", appId, err1))
	}

	//TODO ZYF err
	resetWorkflows(wfMap)

	mbMap, err := createModelMap(moduleObj, wfMap)
	if err != nil {
		return nil, err
	}

	workfowEngine := newWorkflowEngine(smMap, wfMap, mbMap, version)
	return workfowEngine, nil
}

func LoadSceneModuleMap(appId int64) (map[int64]*SceneModule, error) {
	//加载redis中的数据到内存
	sceneModuleMap, err := loadSceneModules(appId)
	if err != nil {
		return sceneModuleMap, err
	}

	//根据app_id,获取对应的场景id
	for sceneId, sceneModule := range sceneModuleMap {
		if sceneModule.AppId != appId {
			delete(sceneModuleMap, sceneId)
		}
	}

	return sceneModuleMap, nil
}

// 获取redis中的数据
func loadSceneModules(appId int64) (map[int64]*SceneModule, error) {
	sceneModuleMapString, err := redis.Handler.Get(context.TODO(), fmt.Sprintf("%v%v", RedisKeySceneModule, appId))
	if err != nil && err.Error() != redis.ErrNil {
		return nil, err
	}

	var sceneModuleMap map[int64]*SceneModule
	err = json.Unmarshal([]byte(sceneModuleMapString), &sceneModuleMap)
	if err != nil {
		return nil, fmt.Errorf("err:%v, sceneModule:%v", err, sceneModuleMapString)
	}

	return sceneModuleMap, nil
}

func LoadWorkflow(appId int64, smMap map[int64]*SceneModule) (map[int64]*Workflow, error) {
	workflowMapStr, err := redis.Handler.Get(context.TODO(), fmt.Sprintf("%v%v", RedisKeyWorkflow, appId))
	if err != nil {
		return nil, err
	}

	var workflowMap map[int64]*Workflow
	err = json.Unmarshal([]byte(workflowMapStr), &workflowMap)
	if err != nil {
		return nil, fmt.Errorf("err:%v, workflow:%v", err, workflowMapStr)
	}

	wfMap := make(map[int64]*Workflow)
	for workflowId, workflow := range workflowMap {
		if _, ok := smMap[workflow.SceneId]; ok {
			workflow.FlowCharts, err = NewWorkflowChart(workflow.FlowChart)
			if err != nil {
				tlog.ErrorCount(context.TODO(), "NewWorkflowChart_err", fmt.Sprintf("wf:%v,err:%v", workflow, err))
				continue
			}
			wfMap[workflowId] = workflow
		}
	}

	return wfMap, nil
}
