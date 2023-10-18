package flowchecker

import (
	"fmt"
	"strconv"
	"strings"
	"tg-service/idl"
)

func replaceWorkflowID(id string, workflowId int64) (string, error) {
	idSeg := strings.Split(id, "-")
	if len(idSeg) != 3 || idSeg[0] != "action" {
		return "", fmt.Errorf("action id " + id + "is illegal")
	}

	// 将 action ID 中的 Workflow ID 替换为当前的 Workflow ID
	idSeg[1] = strconv.Itoa(int(workflowId))
	return strings.Join(idSeg, "-"), nil
}

func replaceIDsForAction(action *idl.Action, workflowId int64) error {
	if action == nil {
		return fmt.Errorf("action is nil")
	}

	// replace action id
	newId, err := replaceWorkflowID(action.ActionId, workflowId)
	if err != nil {
		return fmt.Errorf("action id illegal: %v", err)
	}
	action.ActionId = newId

	if len(action.NextActionIds) == 0 {
		return nil
	}

	// replace next action id
	newNextIds := make([]string, 0)
	for _, id := range action.NextActionIds {
		newId, err = replaceWorkflowID(id, workflowId)
		if err != nil {
			return fmt.Errorf("action id illegal: %v", err)
		}
		newNextIds = append(newNextIds, newId)
	}
	action.NextActionIds = newNextIds
	return nil
}

func isDup(strs []string) bool {
	m := make(map[string]bool)
	for _, s := range strs {
		if m[s] {
			return true
		}
		m[s] = true
	}
	return false
}

func checkFlowChart(flowChart *idl.WorkflowChart) error {

	if flowChart == nil {
		return fmt.Errorf("flow chart is nil")
	}

	if len(flowChart.ActionMap) == 0 {
		return fmt.Errorf("not action in flow chart")
	}

	for id, action := range flowChart.ActionMap {
		if id != action.ActionId {
			return fmt.Errorf("action id is not match, id: %v", id)
		}

		if isDup(action.NextActionIds) {
			return fmt.Errorf("next action ids duplicated, action id %v", action.ActionId)
		}

		for _, id := range action.NextActionIds {
			if flowChart.ActionMap[id] == nil {
				return fmt.Errorf("next action id is not in flow chart, id: %v", id)
			}
		}

	}
	return nil
}

func reaplceActionIdsAndCreateNewFlow(flowChart *idl.WorkflowChart, workflowId int64) (*idl.WorkflowChart, error) {
	newActions := make([]*idl.Action, 0)
	for aid, action := range flowChart.ActionMap {
		err := replaceIDsForAction(action, workflowId)
		if err != nil {
			return nil, fmt.Errorf("replace id for action %v error: %v", aid, err)
		}
		newActions = append(newActions, action)
	}

	newFLowChart := &idl.WorkflowChart{
		ActionMap: make(map[string]*idl.Action),
	}

	for _, action := range newActions {
		newFLowChart.ActionMap[action.ActionId] = action
	}

	if len(newFLowChart.ActionMap) != len(newActions) {
		return nil, fmt.Errorf("there are action ids duplicated")
	}
	return newFLowChart, nil
}

// ReplaceActionIDsForFlow 替换流程中的 ActionID 以匹配新的 workflowId, action ID 格式必须为 action-{workflowId}-{uid}
func ReplaceActionIDsForFlow(flowChart *idl.WorkflowChart, workflowId int64) (*idl.WorkflowChart, error) {

	if err := checkFlowChart(flowChart); err != nil {
		return nil, err
	}

	newFlowChart, err := reaplceActionIdsAndCreateNewFlow(flowChart, workflowId)
	if err != nil {
		return nil, err
	}

	if err := checkFlowChart(newFlowChart); err != nil {
		return nil, err
	}

	return newFlowChart, nil
}
