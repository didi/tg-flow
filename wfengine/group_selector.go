package wfengine

import (
	"fmt"
	"github.com/didi/tg-flow/model"
)

type GroupSelector struct {
	FlowSelector
}

/*
*
apollo分流,
如果出现error，就取缺省分桶，同时返回error信息
*/
func (a *GroupSelector) SelectWorkflowId(sc *model.StrategyContext, sceneModule *SceneModule) (int64, string, error) {
	//var err error
	//if a.ApolloInfo == nil {
	//	return -1, "", fmt.Errorf("apollo info not initialized")
	//}
	//
	//// 优先采用 ApolloInfo 中设置的分流实验名称，没有的话采用场景中配置的分流实验名称
	//if a.ApolloInfo.GetDispatchExperimentName() == "" {
	//	a.ApolloInfo.SetDispatchExperimentName(sceneModule.DispatchExperimentName)
	//}
	//
	//groupName, err := a.ApolloInfo.GetDispatchGroupName()
	//if err != nil {
	//	return -1, groupName, fmt.Errorf("get dispatch groupName fail,groupName=%v, err=%v", groupName, err)
	//}

	workflowId, err := sceneModule.GetWorkflowId(sc.GroupName)
	if err == nil {
		return workflowId, sc.GroupName, nil
	}

	return -1, sc.GroupName, fmt.Errorf("select workflowId error, groupName=%v,workflowId=%v,err=%v", sc.GroupName, workflowId, err)
}
