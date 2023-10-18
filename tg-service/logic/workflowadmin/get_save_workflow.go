package workflowadmin

import (
	"container/list"
	"context"
	"encoding/json"
	"github.com/didi/tg-flow/tg-core/common/mysql"
	"github.com/didi/tg-flow/tg-core/common/tlog"
	"strings"
	"tg-service/common/logs"
	"tg-service/idl"
	"tg-service/logic/workflowadmin/flowchecker"
	"time"
)

//获取可以绘制的流程图信息
func GetWorkFlowChart(exportData *idl.WorkFlowConfig) *idl.ResponseInfo {
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

	workFlowChart := formatWorkflow(workFlow)

	workFlowJson, err := json.Marshal(workFlowChart)
	responseInfo.Content = string(workFlowJson)
	return responseInfo
}

func formatWorkflow(workflow *idl.WorkflowExport) idl.ChartG6 {
	var nodeList []*idl.Node
	var edgeList []*idl.Edge

	actions := workflow.FlowCharts.ActionMap

	typeMap := map[string]string{"task": "rect", "condition": "diamond", "flow": "circle", "timeout": "clock"}

	for _, v := range actions {
		node := &idl.Node{
			Id:             v.ActionId,
			Label:          v.ActionName,
			NodeType:       typeMap[v.ActionType],
			Params:         v.Params,
			Timeout:        v.Timeout,
			RefWorkflowId:  v.RefWorkflowId,
			TimeoutAsync:   v.TimeoutAsync,
			TimeoutDynamic: v.TimeoutDynamic,
		}
		nodeList = append(nodeList, node)
		if strings.EqualFold(v.ActionType, "condition") {
			for i := 0; i < len(v.NextActionIds); i++ {
				edge := &idl.Edge{
					Id:       v.ActionId + "_" + v.NextActionIds[i],
					EdgeType: "line",
					Source:   v.ActionId,
					Target:   v.NextActionIds[i],
					Label:    v.NextConditions[i],
				}
				edgeList = append(edgeList, edge)
			}
			node.Params = v.Params
		} else {
			for _, actionId := range v.NextActionIds {
				edge := &idl.Edge{
					Id:       v.ActionId + "_" + actionId,
					EdgeType: "line",
					Source:   v.ActionId,
					Target:   actionId,
				}
				edgeList = append(edgeList, edge)
			}
		}

	}
	workflowChart := idl.ChartG6{
		Nodes: nodeList,
		Edges: edgeList,
	}
	return workflowChart
}

// SaveImportWorkFlowChart 处理导入信息的保存
func SaveImportWorkFlowChart(saveData *idl.WorkflowChart, workflowId int64, operator string) *idl.ResponseInfo {
	responseInfo := &idl.ResponseInfo{
		Tag:    true,
		ErrMsg: "",
	}

	flowChart, err := flowchecker.ReplaceActionIDsForFlow(saveData, workflowId)
	if err != nil {
		responseInfo.Tag = false
		responseInfo.ErrMsg = err.Error()
		tlog.Handler.Errorf(context.TODO(), logs.DLTagLoadDBFail, "etype=workflowAdmin.saveImportWorkFlowErr||err=%v", err)
		return responseInfo
	}

	updateJson, err := json.Marshal(flowChart)
	if err != nil {
		responseInfo.Tag = false
		responseInfo.ErrMsg = err.Error()
		tlog.Handler.Errorf(context.TODO(), logs.DLTagLoadDBFail, "etype=workflowAdmin.saveImportWorkFlowErr||err=%v", err)
		return responseInfo
	}
	err = SaveUpdateJson(string(updateJson), workflowId)
	if err != nil {
		responseInfo.Tag = false
		responseInfo.ErrMsg = err.Error()
		tlog.Handler.Errorf(context.TODO(), logs.DLTagLoadDBFail, "etype=workflowAdmin.saveImportWorkFlowErr||err=%v", err)
		return responseInfo
	}
	return responseInfo

}

func SaveWorkFlowChart(saveData *idl.ChartG6, workflowId int64, operator string) *idl.ResponseInfo {
	responseInfo := &idl.ResponseInfo{
		Tag:    true,
		ErrMsg: "",
	}
	workflowChart := formatChartG6(saveData)

	updateJson, err := json.Marshal(workflowChart)
	if err != nil {
		responseInfo.Tag = false
		responseInfo.ErrMsg = err.Error()
		tlog.Handler.Errorf(context.TODO(), logs.DLTagLoadDBFail, "etype=workflowAdmin.saveWorkFlowErr||err=%v", err)
		return responseInfo
	}
	err = SaveUpdateJson(string(updateJson), workflowId)
	if err != nil {
		responseInfo.Tag = false
		responseInfo.ErrMsg = err.Error()
		tlog.Handler.Errorf(context.TODO(), logs.DLTagLoadDBFail, "etype=workflowAdmin.saveWorkFlowErr||err=%v", err)
		return responseInfo
	}
	return responseInfo
}
func formatChartG6(saveData *idl.ChartG6) *idl.WorkflowChart {
	nodes := saveData.Nodes
	edges := saveData.Edges
	workflowChart := new(idl.WorkflowChart)
	actionMap := make(map[string]*idl.Action)
	typeMap := map[string]string{"rect": "task", "diamond": "condition", "circle": "flow", "clock": "timeout"}
	for _, node := range nodes {
		action := idl.Action{
			ActionType:    typeMap[node.NodeType],
			ActionId:      node.Id,
			ActionName:    node.Label,
			Params:        node.Params,
			Timeout:       node.Timeout,
			RefWorkflowId: node.RefWorkflowId,
		}
		deleteIndex := list.New()
		for edgeIndex, edge := range edges {
			if edge.Source == node.Id {
				action.NextActionIds = append(action.NextActionIds, edge.Target)
				if action.ActionType == "condition" {
					action.NextConditions = append(action.NextConditions, edge.Label)
				}
				deleteIndex.PushBack(edgeIndex)
			}
		}
		for i := deleteIndex.Back(); i != nil; i = i.Prev() {
			index := i.Value.(int)
			edges = append(edges[:index], edges[index+1:]...)
		}
		actionMap[node.Id] = &action
	}
	workflowChart.ActionMap = actionMap
	return workflowChart

}

//更新workflow信息
func SaveUpdateJson(updateJson string, workflowId int64) error {
	ctx := context.TODO()
	updateTime := time.Now()
	updateModifySql := "update workflow set flow_chart = ?, update_time = ? where id = ?"
	_, updateModifyErr := mysql.Handler.Exec(updateModifySql, updateJson, updateTime, workflowId)

	tlog.Handler.Infof(ctx, logs.DLTagProcessLog, "etype=workflowadmin.UpdateDBWorkFlow||data=%v||err=本条记录修改成功", updateJson)
	return updateModifyErr
}
