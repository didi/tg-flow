/**
Description : loader of workflow config info from apollo config
Author		: dayunzhangyunfeng@didiglobal.com
Date		: 2023-01-09 14:30
*/

package wfengine

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/didi/tg-flow/common/tlog"
	"strconv"
	"strings"
)

const (
	dftRange   = "0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25,26,27,28,29,30,31,32,33,34,35,36,37,38,39,40,41,42,43,44,45,46,47,48,49,50,51,52,53,54,55,56,57,58,59,60,61,62,63,64,65,66,67,68,69,70,71,72,73,74,75,76,77,78,79,80,81,82,83,84,85,86,87,88,89,90,91,92,93,94,95,96,97,98,99"
	sceneKey   = "scene"
	versionKey = "version"
)

//func GetLatestVersionFromApollo(namespace, configName string) (string, error) {
//	apolloConfig, err := apollo.NewApolloConfig(namespace, configName)
//	if err != nil || apolloConfig == nil {
//		return "", fmt.Errorf("NewApolloConfig error, namespace=%v, configName=%v, err=%v", namespace, configName, err)
//	}
//
//	configParams := apolloConfig.GetConfigs()
//	if configParams == nil || configParams[versionKey] == "" {
//		return "", fmt.Errorf("GetConfigs fail, configParam is nil or version empty, namespace=%v, configName=%v, configParams.lenth=%v", namespace, configName, len(configParams))
//	}
//
//	return configParams[versionKey], nil
//}

func NewWorkflowEngineFromConfig(moduleObj ModuleObjBase, configParams map[string]string) (*WorkflowEngine, error) {
	smMap, wfMap, version, err := loadSceneModuleWorkflowFromApollo(configParams)
	if err != nil {
		return nil, err
	}

	return NewWorkflowEngine(smMap, wfMap, version, moduleObj)
}

func loadSceneModuleWorkflowFromApollo(configParams map[string]string) (map[int64]*SceneModule, map[int64]*Workflow, string, error) {
	if configParams == nil {
		return nil, nil, "", fmt.Errorf("GetConfigs error, configParams is nil")
	}

	if configParams[sceneKey] == "" {
		return nil, nil, "", fmt.Errorf("configParams error, sceneKey[%v] is empty", sceneKey)
	}

	smMap, err := loadSceneModuleString(configParams[sceneKey])
	if err != nil {
		return nil, nil, "", err
	}

	//更新workflow
	wfMap, err := LoadWorkflowFromApollo(configParams)
	if err != nil {
		return nil, nil, "", err
	}

	return smMap, wfMap, configParams[versionKey], nil
}

func loadSceneModuleString(sceneModuleMapString string) (map[int64]*SceneModule, error) {
	if len(sceneModuleMapString) <= 0 {
		return nil, fmt.Errorf("the param sceneModuleMapString:[%v] must not be empty", sceneModuleMapString)
	}

	var sceneModuleMap map[int64]*SceneModule
	err := json.Unmarshal([]byte(sceneModuleMapString), &sceneModuleMap)
	if err != nil {
		return nil, fmt.Errorf("invalid param sceneModuleMapString:[%v], unmarshal err:%v", sceneModuleMapString, err)
	}

	return sceneModuleMap, nil
}

func LoadWorkflowFromApollo(configParams map[string]string) (map[int64]*Workflow, error) {
	if configParams == nil || len(configParams) <= 0 {
		return nil, fmt.Errorf("configParam is empty:%v", configParams)
	}

	//对每个文件逐个加到map
	wfMap := make(map[int64]*Workflow)
	for key, val := range configParams {
		if key == sceneKey || key == versionKey {
			continue
		}
		wf, err := createWorkflowFromKV(key, val)
		if wf == nil || err != nil {
			tlog.Handler.ErrorCount(context.TODO(), "loadWorkflowFromKV_err", fmt.Sprintf("wf:%v,err:%v", wf, err))
			continue
		}
		wf.FlowCharts, err = NewWorkflowChart(wf.FlowChart)
		if err != nil {
			tlog.Handler.ErrorCount(context.TODO(), "NewWorkflowChart_err", fmt.Sprintf("wf:%v,err:%v", wf, err))
			continue
		}
		wfMap[wf.Id] = wf
	}

	return wfMap, nil
}

func createWorkflowFromKV(key, val string) (*Workflow, error) {
	strs := strings.Split(key, "-")
	if len(strs) < 3 {
		return nil, fmt.Errorf("invalid key:%v, it must start with :workflow-sceneId-workflowId", key)
	}

	sceneId, err := strconv.ParseInt(strs[1], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid sceneId:%v in key:%v,err=%v", strs[1], key, err)
	}
	workflowId, err := strconv.ParseInt(strs[2], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid workflowId:%v in key:%v,err=%v", strs[2], key, err)
	}

	workflow := &Workflow{
		Id:          workflowId,
		DimensionId: -1,
		SceneId:     sceneId,
		FlowChart:   val,
		FlowCharts:  nil,
		FlowBranch:  nil,
		IsDefault:   1,
		Range1:      dftRange,
		Range2:      "",
		Remark:      "",
		//UpdateTime:,
		GroupName: "",
	}

	return workflow, nil
}
