package store

import (
	"context"
	"encoding/json"
	"github.com/didi/tg-flow/tg-core/common/redis"
	"github.com/didi/tg-flow/tg-core/consts"
	"github.com/didi/tg-flow/tg-core/model"
)

var WorkflowMap map[int64]model.Workflow = make(map[int64]model.Workflow)

func LoadWorkflow() error {
	workflowMapByte, err := redis.Handler.Get(context.TODO(), consts.StrategyWorkflowMap)
	if err != nil && err.Error() != redis.ErrNil {
		return err
	}

	var workflowMap map[int64]model.Workflow
	err = json.Unmarshal([]byte(workflowMapByte), &workflowMap)
	if err != nil {
		return err
	}

	WorkflowMap = workflowMap

	return nil
}
