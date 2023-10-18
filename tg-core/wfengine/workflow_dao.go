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
	RedisKeyWorkflow	= "workflow_app_"
)

func LoadWorkflow(appId int64, smMap map[int64]*SceneModule) (map[int64]*Workflow, error) {
	workflowMapStr, err := redis.Handler.Get(context.TODO(), fmt.Sprintf("%v%v",RedisKeyWorkflow, appId))
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
				tlog.ErrorCount(context.TODO(),"NewWorkflowChart_err", fmt.Sprintf("wf:%v,err:%v", workflow, err))
				continue
			}
			wfMap[workflowId] = workflow
		}
	}

	return wfMap, nil
}