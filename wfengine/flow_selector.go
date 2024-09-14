package wfengine

import "github.com/didi/tg-flow/model"

/*
*
there are 3 flow selector:

	1: random
	2: custom
	3: apollo (apollo platform in didi)
*/
type FlowSelector interface {
	SelectWorkflowId(sc *model.StrategyContext, sceneModule *SceneModule) (int64, string, error)
}
