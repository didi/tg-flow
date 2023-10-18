package workflowadmin

import (
	"context"
	"encoding/json"
	"github.com/didi/tg-flow/tg-core/common/mysql"
	"github.com/didi/tg-flow/tg-core/common/tlog"
	"github.com/didi/tg-flow/tg-core/wfengine"
	"strconv"
	"tg-service/common/logs"
	"tg-service/common/template"
	"tg-service/idl"
	"time"
)

//导出数据
func ExportWorkFlow(exportData *idl.WorkFlowConfig) *idl.ResponseInfo {
	responseInfo := &idl.ResponseInfo{
		Tag:    true,
		ErrMsg: "",
	}
	workFlow, err := GetWorkFlowExport(exportData)
	if err != nil {
		responseInfo.Tag = false
		responseInfo.ErrMsg = err.Error()
		tlog.Handler.Errorf(context.TODO(), logs.DLTagLoadDBFail, "etype=workflowAdmin.getWorkFlowErr||err=%v", err)
		return responseInfo
	}
	response, err := json.Marshal(workFlow)
	if err != nil {
		responseInfo.Tag = false
		responseInfo.ErrMsg = err.Error()
		tlog.Handler.Errorf(context.TODO(), logs.DLTagLoadDBFail, "etype=workflowAdmin.getWorkFlowErr||err=%v", err)
		return responseInfo
	}
	responseInfo.Content = string(response)
	return responseInfo
}
func GetWorkFlowExport(exportData *idl.WorkFlowConfig) (*idl.WorkflowExport, error) {
	workflow, flowChartStr, err := getWorkFlow(exportData)
	if err != nil {
		return nil, err
	}
	var actionMap = new(idl.WorkflowChart)
	if len(flowChartStr) == 0 {
		actionMap.ActionMap = make(map[string]*idl.Action)
	} else {
		err = json.Unmarshal([]byte(flowChartStr), actionMap)
	}
	if err != nil {
		return nil, err
	}
	var exportWorkflow = idl.WorkflowExport{
		Id:         workflow.Id,
		SceneId:    workflow.SceneId,
		FlowCharts: actionMap,
	}
	sceneId, err := strconv.Atoi(strconv.FormatInt(workflow.SceneId, 10))
	if err != nil {
		return nil, err
	}
	exportWorkflow.SceneName = template.SceneIdAndNameMap[sceneId]
	return &exportWorkflow, nil

}
func getWorkFlowShow(exportData *idl.WorkFlowConfig) (*wfengine.Workflow, string, error) {
	workflow, flowChartStr, err := getWorkFlow(exportData)
	return workflow, flowChartStr, err
}
func getWorkFlow(workflowConfig *idl.WorkFlowConfig) (*wfengine.Workflow, string, error) {
	sql := "select id,dimension_id,experiment_id,flow_chart,is_default,range1,range2,remark,group_name from workflow where id = ?"
	rows, err := mysql.Handler.Query(sql, workflowConfig.WorkFlowId)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()
	err = rows.Err()
	if err != nil {
		return nil, "", err
	}
	var workflow = new(wfengine.Workflow)
	var flowChartStr string
	for rows.Next() {
		err := rows.Scan(&workflow.Id, &workflow.DimensionId, &workflow.SceneId, &flowChartStr, &workflow.IsDefault, &workflow.Range1, &workflow.Range2, &workflow.Remark, &workflow.GroupName)
		if err != nil {
			return nil, "", err
		}
		//if len(flowChartStr) > 0 {
		//	workflow.FlowCharts, workflow.FlowBranch, err = wfengine.NewWorkflowChart(flowChartStr)
		//	if err != nil {
		//		return nil, "", err
		//	}
		//}
		workflow.UpdateTime = time.Now()
	}
	return workflow, flowChartStr, nil
}
