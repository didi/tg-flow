package appconfigadmin

import (
	"container/list"
	"context"
	"encoding/json"
	"github.com/didi/tg-flow/tg-core/common/mysql"
	"github.com/didi/tg-flow/tg-core/common/tlog"
	"github.com/didi/tg-flow/tg-core/wfengine"
	"tg-service/common/logs"
	"tg-service/idl"

	"strconv"
)

//导出数据
func ExportAppConfig(exportData *idl.AppConfigInfo) *idl.ResponseInfo {
	responseInfo := &idl.ResponseInfo{
		Tag:    true,
		ErrMsg: "",
	}
	sceneConfigList, err := getSceneListByAppId(exportData)
	if err != nil {
		responseInfo.Tag = false
		responseInfo.ErrMsg = err.Error()
		tlog.Handler.Errorf(context.TODO(), logs.DLTagLoadDBFail, "etype=appconfigadmin.getSceneIdErr||err=%v", err)
		return responseInfo
	}
	sceneMap := make(map[string]*wfengine.SceneModule)
	workflowList := list.New()
	for scene := sceneConfigList.Front(); scene != nil; scene = scene.Next() {
		sceneModule, ok := scene.Value.(*wfengine.SceneModule)
		if ok {
			sceneMap[strconv.FormatInt(sceneModule.Id, 10)] = sceneModule
			workflowListByScene, err := getSceneWorkflowList(sceneModule)
			if err != nil {
				responseInfo.Tag = false
				responseInfo.ErrMsg = err.Error()
				tlog.Handler.Errorf(context.TODO(), logs.DLTagLoadDBFail, "etype=appconfigadmin.getSceneWorkflowErr||err=%v", err)
				return responseInfo
			}
			workflowList.PushBackList(workflowListByScene)
		}
	}
	content := make(map[string]interface{})
	workflowStr := "["
	for workflow := workflowList.Front(); workflow != nil; workflow = workflow.Next() {
		json, err := json.Marshal(workflow.Value)
		if err != nil {
			responseInfo.Tag = false
			responseInfo.ErrMsg = err.Error()
			tlog.Handler.Errorf(context.TODO(), logs.DLTagLoadDBFail, "etype=appconfigadmin.getSceneWorkflowErr||err=%v", err)
			return responseInfo
		}
		workflowStr = workflowStr + "'" + string(json) + "',"
	}
	workflowStr = workflowStr + "]"
	sceneJson, err := json.Marshal(sceneMap)
	content["scene"] = string(sceneJson)
	content["workflowList"] = workflowStr
	response, err := json.Marshal(content)
	responseInfo.Content = string(response)
	return responseInfo
}

func ExportAppConfigForApollo(exportData *idl.AppConfigInfo) map[string]interface{} {
	appWorkflowMap := make(map[string]interface{})
	sceneConfigList, err := getSceneListByAppId(exportData)
	if err != nil {
		tlog.Handler.Errorf(context.TODO(), logs.DLTagLoadDBFail, "etype=asyncWorkflowToApollo.getSceneIdErr||err=%v", err)
		return appWorkflowMap
	}
	sceneMap := make(map[string]*wfengine.SceneModule)
	workflowList := list.New()
	for scene := sceneConfigList.Front(); scene != nil; scene = scene.Next() {
		sceneModule, ok := scene.Value.(*wfengine.SceneModule)
		if ok {
			sceneMap[strconv.FormatInt(sceneModule.Id, 10)] = sceneModule
			workflowListByScene, err := getSceneWorkflowList(sceneModule)
			if err != nil {
				tlog.Handler.Errorf(context.TODO(), logs.DLTagLoadDBFail, "etype=asyncWorkflowToApollo.getSceneWorkflowErr||err=%v", err)
				return appWorkflowMap
			}
			workflowList.PushBackList(workflowListByScene)
		}
	}
	appWorkflowMap["scene"] = sceneMap
	for workflow := workflowList.Front(); workflow != nil; workflow = workflow.Next() {
		workflowValue, ok := workflow.Value.(idl.WorkflowExport)
		if ok {
			workflowKey := "workflow-" + strconv.FormatInt(workflowValue.SceneId, 10) + "-" + strconv.FormatInt(workflowValue.Id, 10) + "-" + workflowValue.SceneName
			appWorkflowMap[workflowKey] = workflowValue.FlowCharts
		}
	}
	return appWorkflowMap
}

func getSceneListByAppId(exportData *idl.AppConfigInfo) (*list.List, error) {
	sql := "select id,`name`,app_id,bucket_type,update_time,flow_type,exp_name from scene_config where app_id = ?"
	rows, err := mysql.Handler.Query(sql, exportData.Id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	var sceneConfigList = list.New()
	for rows.Next() {
		var sceneConfig wfengine.SceneModule
		err := rows.Scan(&sceneConfig.Id, &sceneConfig.Name, &sceneConfig.AppId, &sceneConfig.BucketType, &sceneConfig.UpdateTime, &sceneConfig.FlowType, &sceneConfig.DispatchExperimentName)
		if err != nil {
			return nil, err
		}
		sceneConfigList.PushBack(&sceneConfig)
	}
	return sceneConfigList, nil
}

func getSceneWorkflowList(sceneModule *wfengine.SceneModule) (*list.List, error) {
	sql := "select id,experiment_id,flow_chart,group_name,is_default from workflow where status = 1 and experiment_id = ?"
	sceneWorkflowList := list.New()

	var sceneId = sceneModule.Id
	groupWorkflowMap := make(map[string]int64)
	rows, err := mysql.Handler.Query(sql, sceneId)
	if err != nil {
		return list.New(), err
	}
	defer rows.Close()
	err = rows.Err()
	if err != nil {
		return list.New(), err
	}
	var defaultWorkflowId int64
	for rows.Next() {
		var workflow = new(wfengine.Workflow)
		var flowChartStr string
		err := rows.Scan(&workflow.Id, &workflow.SceneId, &flowChartStr, &workflow.GroupName, &workflow.IsDefault)
		if err != nil {
			return list.New(), err
		}
		if defaultWorkflowId == 0 && workflow.IsDefault == 1 {
			defaultWorkflowId = workflow.Id
		}
		var actionMap = new(idl.WorkflowChart)
		if len(flowChartStr) == 0 {
			actionMap.ActionMap = make(map[string]*idl.Action)
		} else {
			err = json.Unmarshal([]byte(flowChartStr), actionMap)
		}
		if err != nil {
			return list.New(), err
		}
		var exportWorkflow = idl.WorkflowExport{
			Id:         workflow.Id,
			SceneId:    workflow.SceneId,
			FlowCharts: actionMap,
			SceneName:  sceneModule.Name,
			GroupName:  workflow.GroupName,
		}
		sceneWorkflowList.PushBack(exportWorkflow)
		groupWorkflowMap[workflow.GroupName] = workflow.Id
	}
	sceneModule.GroupWorkflowMap = groupWorkflowMap
	sceneModule.DefaultWorkflowId = defaultWorkflowId
	return sceneWorkflowList, nil
}
