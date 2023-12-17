/*
*
Description : workflow engine v3.1
Author		: dayunzhangyunfeng@didiglobal.com
Date		: 2021-05-14
*/
package wfengine

import (
	"context"
	"errors"
	"fmt"
	"github.com/didi/tg-flow/tg-core/common/tlog"
	"github.com/didi/tg-flow/tg-core/consts"
	"github.com/didi/tg-flow/tg-core/model"
	"strings"
	"sync"
	"time"
)

const (
	action0      = "action0"
	ErrNoUnknown = 1000
)

type WorkflowEngine struct {
	sceneModuleMap map[int64]*SceneModule `json:"scene_module_map"`
	workflowMap    map[int64]*Workflow    `json:"workflows"`
	modelBaseMap   map[string]IModelBase  `json:"modules"`
	updateTime     string                 `json:"update_time"`
	condExecutors  *CondExecutors         `json:"cond_executors"`
}

func resetWorkflows(wfMap map[int64]*Workflow) {
	//fmt.Println("wfMap.len=", len(wfMap))
	if len(wfMap) == 0 {
		return
	}

	for _, workflow := range wfMap {
		err := resetWorkflow(wfMap, workflow)
		if err != nil {
			tlog.ErrorCount(context.TODO(), "resetWorkflow err", fmt.Sprintf("workflow=%+v", workflow))
			continue
		}
	}
}

func resetWorkflow(wfMap map[int64]*Workflow, workflow *Workflow) error {
	if workflow == nil {
		return errors.New("workflow is nil:%v")
	}
	//fmt.Println("开始for循环处理flowAction")
	for {
		flowActionMap := make(map[string]*Action)
		//fmt.Println("workflow.FlowCharts.ActionMap.len===>", len(workflow.FlowCharts.ActionMap), workflow.FlowCharts.HashCondition)
		for _, action := range workflow.FlowCharts.ActionMap {
			if action.ActionType == ActionTypeFlow {
				if action.RefWorkflowId <= 0 {
					tlog.ErrorCount(context.TODO(), "resetWorkflow_err", fmt.Sprintf("ref_workflow_id must >0, flowAction.Action=%+v", action))
					continue
				}
				flowActionMap[action.ActionId] = action
			}
		}
		if len(flowActionMap) < 1 {
			break
		}
		//fmt.Println()
		//对所有flowAction
		for _, flowAction := range flowActionMap {
			//get workflowid
			refWf := wfMap[flowAction.RefWorkflowId]
			if refWf == nil {
				tlog.ErrorCount(context.TODO(), "resetWorkflow_err", fmt.Sprintf("reference workflow not exist, flowAction.Action=%+v", flowAction))
				continue
			}

			//refwf actions add to workflow
			for refActionId, refAction := range refWf.FlowCharts.ActionMap {
				rAction := refAction.clone()
				workflow.FlowCharts.ActionMap[refActionId] = rAction
			}

			//join refwf head
			firstActionId := refWf.FlowCharts.FirstActionId
			if flowAction.ActionId == workflow.FlowCharts.FirstActionId {
				workflow.FlowCharts.FirstActionId = firstActionId
			} else {
				for _, prevId := range flowAction.PrevActionIds {
					strNextActionIds := strings.Join(workflow.FlowCharts.ActionMap[prevId].NextActionIds, ",")
					strNewNextActionIds := ""
					if flowAction.Timeout > 0 {
						strNewNextActionIds = strNextActionIds + "," + firstActionId
					} else {
						strNewNextActionIds = strings.ReplaceAll(strNextActionIds, flowAction.ActionId, firstActionId)
					}
					workflow.FlowCharts.ActionMap[prevId].NextActionIds = strings.Split(strNewNextActionIds, ",")
				}

				//reset PrevActionIds for first action
				workflow.FlowCharts.ActionMap[firstActionId].PrevActionIds = flowAction.PrevActionIds
			}

			//join refwf tail
			lastActionId := refWf.FlowCharts.LastActionId
			if flowAction.ActionId == workflow.FlowCharts.LastActionId {
				workflow.FlowCharts.LastActionId = lastActionId
			} else {
				for _, nextId := range flowAction.NextActionIds {
					strPrevActionIds := strings.Join(workflow.FlowCharts.ActionMap[nextId].PrevActionIds, ",")
					strNewPrevActionIds := ""
					if flowAction.Timeout > 0 {
						strNewPrevActionIds = strPrevActionIds + "," + lastActionId
					} else {
						strNewPrevActionIds = strings.ReplaceAll(strPrevActionIds, flowAction.ActionId, lastActionId)
					}
					workflow.FlowCharts.ActionMap[nextId].PrevActionIds = strings.Split(strNewPrevActionIds, ",")
				}
				//reset NextActionIds for last action
				workflow.FlowCharts.ActionMap[lastActionId].NextActionIds = flowAction.NextActionIds
			}

			//删除
			if flowAction.Timeout > 0 {
				flowAction.ActionType = ActionTypeTimeout
			} else {
				delete(workflow.FlowCharts.ActionMap, flowAction.ActionId)
			}
		}
	}

	//fmt.Println("before workflow.FlowCharts.HashCondition=>", workflow.FlowCharts.HashCondition)
	for _, action := range workflow.FlowCharts.ActionMap {
		//fmt.Println("actionId=", action.ActionId, "action.ActionType=", action.ActionType)
		if action.ActionType == ActionTypeCond {
			workflow.FlowCharts.HashCondition = true
			break
		}

	}
	//fmt.Println("after  workflow.FlowCharts.HashCondition, workflowId=",workflow.Id, "hasCondition=",workflow.FlowCharts.HashCondition)
	return nil
}

func newWorkflowEngine(sceneModuleMap map[int64]*SceneModule, workflowMap map[int64]*Workflow, modelBaseMap map[string]IModelBase, version string) *WorkflowEngine {
	//ut := fmt.Sprintf("%v", time.Now().Format("2006-01-02 15:04:05"))
	cExecutors := GetCondExecutors()
	wfe := &WorkflowEngine{
		sceneModuleMap: sceneModuleMap,
		workflowMap:    workflowMap,
		modelBaseMap:   modelBaseMap,
		updateTime:     version,
		condExecutors:  cExecutors,
	}

	return wfe
}

func createModelMap(moduleObj ModuleObjBase, wfMap map[int64]*Workflow) (map[string]IModelBase, error) {
	if len(wfMap) == 0 {
		return nil, fmt.Errorf("workflow map cannot be nil")
	}

	ctx := context.TODO()
	modelBaseMap := make(map[string]IModelBase)
	for _, wf := range wfMap {
		if wf.FlowCharts == nil {
			tlog.ErrorCount(ctx, "create_modelbase_err", fmt.Sprintf("flow_charts is nil, workflow=%v", wf))
			continue
		}

		for actionId, action := range wf.FlowCharts.ActionMap {
			if action == nil || (action.ActionType != ActionTypeTask && action.ActionType != ActionTypeTimeout) || modelBaseMap[actionId] != nil {
				continue
			}

			mb, err := createModelBase(moduleObj, action)
			if mb == nil || err != nil {
				tlog.ErrorCount(ctx, "create_modelbase_err", fmt.Sprintf("workflow=%v, action=%v, mb=%v, err=%v", wf, action, mb, err))
				continue
			}
			modelBaseMap[actionId] = mb
		}

	}

	return modelBaseMap, nil
}

func createModelBase(moduleObj ModuleObjBase, action *Action) (IModelBase, error) {
	if action == nil || moduleObj == nil {
		return nil, fmt.Errorf("action or moduleObj empty, action:%v,moduleObj:%v", action, moduleObj)
	}

	vMap := make(map[string]string)
	if len(action.Params) > 0 {
		for _, param := range action.Params {
			vMap[param.Name] = param.Value
		}
	}
	//vMap["Name"] = action.ActionName
	mb := moduleObj.NewObj(action.ActionName, vMap)
	if mb == nil {
		return mb, fmt.Errorf("create ModelBase instance error, action:%v, vMap:%v", action, vMap)
	}

	mb.SetName(action.ActionName)
	err := reflectModuleField(mb, vMap)
	if err != nil {
		tlog.ErrorCount(context.TODO(), "createModelBase_err", fmt.Sprintf("set module field fail, actionName:%v, vMap:%v, error:%v", action.ActionName, vMap, err))
	}

	return mb, nil
}

func (w *WorkflowEngine) GetVersion() string {
	return w.updateTime
}

func (w *WorkflowEngine) RegisterCondExecutor(conditionName string, executor interface{}) {
	w.condExecutors.RegisterCondExecutor(conditionName, executor)
}

func (w *WorkflowEngine) Run(ctx context.Context, sc *model.StrategyContext) {
	defer func() {
		if err := recover(); err != nil {
			err := fmt.Errorf("WorkflowEngine run panic:%v", err)
			sc.SetError(action0, err)
		}
	}()
	//选择一个workflow策略
	flow, err := w.selectWorkFlow(ctx, sc)
	if err != nil {
		sc.SetError(action0, fmt.Errorf("flow:%v, err:%v", flow, err))
		return
	}

	//fmt.Println(fmt.Sprintf("准备Run,先获取图,flow=%+v", flow.FlowCharts))
	flowChart := flow.GetWorkflowChart()
	//fmt.Println(fmt.Sprintf("准备Run,先获取图,flowChart=%+v", flowChart))
	if flowChart == nil {
		sc.SetError(action0, errors.New("flowChart is nil"))
		return
	}

	wgMap, tsMap := flowChart.CreateWaitMap()

	waitedMap := &sync.Map{}
	wgn := &sync.WaitGroup{}
	wgn.Add(1)
	skipedActionIdPairs := &sync.Map{}

	w.doExecuteModule(ctx, sc, flowChart, skipedActionIdPairs, wgMap, tsMap, waitedMap, flowChart.FirstActionId, wgn)
	wgn.Wait()
}

func (w *WorkflowEngine) selectWorkFlow(ctx context.Context, sc *model.StrategyContext) (*Workflow, error) {
	sceneModule, okE := w.sceneModuleMap[sc.SceneId]
	if !okE {
		return nil, fmt.Errorf("no sceneModule found for scene_id:%v", sc.SceneId)
	}

	var flowId int64
	var slotId int
	var groupName string
	var err error
	if sc.FlowId > 0 { //测试
		flowId = sc.FlowId
	} else {
		//根据该场景设置的分流方式，获取workflowId
		if sceneModule.FlowType == consts.FLOW_BY_ONLINE_RANDOM {
			flowId, slotId = FlowByOnlineRandom(sc, sceneModule)
		} else if sceneModule.FlowType == consts.FLOW_BY_APOLLO {
			flowId, groupName, err = FlowByApollo(sc, sceneModule)
		}

		// 如果没有找到 flowId，采用默认的 Workflow
		if err != nil || flowId == 0 {
			tlog.ErrorCount(ctx, "select_workflow_err", fmt.Sprintf("select workflow by apollo error, err=%v, use default workflow", err))
			flowId = sceneModule.DefaultWorkflowId
			err = nil
		}

		sc.FlowId = flowId
	}
	//根据分流选出的workflow_id，取得对应的实验策略配置
	flow, okW := w.workflowMap[flowId]
	if !okW {
		err := fmt.Errorf("no workflow found, slotId=%v||flowId=%v||groupId=%v||err=%v", slotId, flowId, groupName, err)
		return nil, err
	}
	sc.Set("flow", flow)
	sc.Set("groupId", slotId)
	sc.Set("scene", sceneModule)
	return flow, err
}

/*
*

	新版本（v1.2.11）后因为有动态条件节点，原流程执行模式有较大变化，不再包含预处理，需实时判断节点走向和执行、跳过节点等操作。

*
*/
func (w *WorkflowEngine) doExecuteModule(ctx context.Context, sc *model.StrategyContext, flowChart *WorkflowChart, skipedActionIdPairs *sync.Map, wgMap map[string]*sync.WaitGroup, tsMap map[string]*TimeWaiter, waitedMap *sync.Map, actionId string, wgn *sync.WaitGroup) {
	action, ok := flowChart.ActionMap[actionId]
	//fmt.Println("开始执行actionId===>", fmt.Sprintf("%+v",action), ok)
	if !ok {
		return
	}

	//fmt.Println("action.PrevActionIds===>", strings.Join(action.PrevActionIds, ","))
	//is merge
	if len(action.PrevActionIds) > 1 {
		_, ok := waitedMap.LoadOrStore(action.ActionId, true)
		if ok {
			return
		}

		timeoutActionId := ""
		for _, prevActionId := range action.PrevActionIds {
			if flowChart.ActionMap[prevActionId].ActionType == ActionTypeTimeout {
				timeoutActionId = prevActionId
				break
			}
		}

		if timeoutActionId == "" {
			wgMap[action.ActionId].Wait()
		} else {
			go func() {
				wgMap[action.ActionId].Wait()
				if tsMap[timeoutActionId] != nil {
					tsMap[timeoutActionId].Done()
				}
			}()

			if !tsMap[timeoutActionId].Wait() {
				if mb, _ := w.modelBaseMap[timeoutActionId]; mb != nil {
					if flowChart.ActionMap[timeoutActionId].TimeoutAsync {
						go mb.OnTimeout(ctx, sc)
					} else {
						mb.OnTimeout(ctx, sc)
					}
				}
			}
		}

	}

	toExecuteActionId := w.executeModule(ctx, sc, flowChart, action, skipedActionIdPairs, wgMap, tsMap, wgn)
	if len(action.NextActionIds) == 0 {
		return
	}

	if len(action.NextActionIds) == 1 {
		w.doExecuteModule(ctx, sc, flowChart, skipedActionIdPairs, wgMap, tsMap, waitedMap, action.NextActionIds[0], wgn)
		return
	}

	if _, ok := flowChart.ActionMap[toExecuteActionId]; ok {
		w.doExecuteModule(ctx, sc, flowChart, skipedActionIdPairs, wgMap, tsMap, waitedMap, toExecuteActionId, wgn)
		return
	}

	for _, nextActionId := range action.NextActionIds {
		go w.doExecuteModule(ctx, sc, flowChart, skipedActionIdPairs, wgMap, tsMap, waitedMap, nextActionId, wgn)
	}

}

func (w *WorkflowEngine) skipBranch(flowChart *WorkflowChart, wgMap map[string]*sync.WaitGroup, skipedActionIdPairs *sync.Map, toExcludeActionId string, prevAction *Action, action *Action) {
	if prevAction == nil || action == nil {
		return
	}

	if len(action.PrevActionIds) > 1 {
		if _, ok := skipedActionIdPairs.LoadOrStore(prevAction.ActionId+"_"+action.ActionId, ""); !ok {
			wgMap[action.ActionId].Done()
		}
	}

	action.Detach(prevAction)

	if len(action.PrevActionIds) > 0 {
		return
	}

	nextCount := len(action.NextActionIds)
	for i := nextCount - 1; i >= 0; i-- {
		nextActionId := action.NextActionIds[i]
		nextAction := flowChart.ActionMap[nextActionId]
		w.skipBranch(flowChart, wgMap, skipedActionIdPairs, toExcludeActionId, action, nextAction)
	}
}

func (w *WorkflowEngine) executeModule(ctx context.Context, sc *model.StrategyContext, flowChart *WorkflowChart, action *Action, skipedActionIdPairs *sync.Map, wgMap map[string]*sync.WaitGroup, tsMap map[string]*TimeWaiter, wgn *sync.WaitGroup) string {
	toExeActionId := ""

	defer func() {
		if err := recover(); err != nil {
			tlog.ErrorCount(ctx, "executeModule_err", fmt.Sprintf("actionId=%+v,err=%+v", action.ActionId, err))
			sc.SetError(action.ActionId, fmt.Errorf("%v", err))
			sc.Skip(ErrNoUnknown, fmt.Sprintf("executeModule_err, actionId=%+v, err=%+v", action.ActionId, err))
		}

		nextCount := len(action.NextActionIds)
		for i := nextCount - 1; i >= 0; i-- {
			nextActionId := action.NextActionIds[i]
			nextAction, ok := flowChart.ActionMap[nextActionId]
			if !ok {
				continue
			}

			if action.ActionType == ActionTypeCond {
				if nextActionId != toExeActionId {
					w.skipBranch(flowChart, wgMap, skipedActionIdPairs, toExeActionId, action, nextAction)
				}
			}

			if len(nextAction.PrevActionIds) > 1 {
				if _, ok := skipedActionIdPairs.LoadOrStore(action.ActionId+"_"+nextActionId, ""); !ok {
					wgMap[nextActionId].Done()
				}
			}

		}

		if action.ActionId == flowChart.LastActionId {
			wgn.Done()
		}

		//fmt.Println(fmt.Sprintf("====>Finish time=%v ,actionId=%v, actionName=%v, action=%+v", time.Now().Format("2006-01-02 15:04:05.000"), action.ActionId, action.ActionName,action))
	}()

	if action == nil {
		return ""
	}
	//条件节点执行
	if action.ActionType == ActionTypeCond {
		var err error
		toExeActionId, err = action.executeCond(sc.ContextMap)
		if err != nil {
			tlog.ErrorCount(ctx, "executeCond_err", fmt.Sprintf("execute condition error, actionId=%+v, toExeActionId=%+v, err=%+v", action.ActionId, toExeActionId, err))
			sc.SetError(action.ActionId, fmt.Errorf("action.executeCond error, ActionName:%v, err:%v", action.ActionName, err))
		}
		return toExeActionId
	}
	//任务节点执行
	if moduleBase, ok := w.modelBaseMap[action.ActionId]; ok {
		if !sc.IsSkip() {
			startTime := time.Now().UnixNano() / 1e6
			if action.Timeout > 0 && action.ActionType != ActionTypeTimeout {
				go func() {
					moduleBase.DoAction(ctx, sc)
					if tsMap[action.ActionId] != nil {
						tsMap[action.ActionId].Done()
					}
				}()
				if !tsMap[action.ActionId].Wait() {
					sc.AddTimeoutAction(action.ActionName)
					if action.TimeoutAsync {
						go moduleBase.OnTimeout(ctx, sc)
					} else {
						moduleBase.OnTimeout(ctx, sc)
					}
				} else {
					if action.TimeoutDynamic {
						for _, nextActionId := range action.NextActionIds {
							if ts, ok := tsMap[nextActionId]; ok {
								leftTime := action.Timeout - (time.Now().UnixNano()/1e6 - startTime)
								ts.AddTimeout(leftTime)
							}
						}
					}
				}
			} else {
				moduleBase.DoAction(ctx, sc)
			}

			endTime := time.Now().UnixNano() / 1e6
			sc.SetModuleResult(action.ActionId, action.ActionName, endTime-startTime)

		}
	} else {
		sc.SetError(action.ActionId, fmt.Errorf("module not found in map, moduleMap:%v, moduleName:%v", w.modelBaseMap, action.ActionName))
		sc.Skip(ErrNoUnknown, fmt.Sprintf("moduleMap=%v||moduleName=%v", w.modelBaseMap, action.ActionName))
	}

	return ""
}
