/**
	Description : workflow v3.0 with branch
	Author		: dayunzhangyunfeng@didiglobal.com
	Date		: 2021-05-14
 */
package wfengine

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	ActionTypeTask	= "task"
	ActionTypeCond	= "condition"
	ActionTypeFlow	= "flow"
	ActionTypeTimeout= "timeout"
	BranchKeyDefault = "default"
	BranchKeyJoiner = "_"
	defaultBranch	= "*"
)

type Workflow struct {
	Id          int64                     `json:"id"`
	DimensionId int64                     `json:"dimension_id"`
	SceneId     int64                     `json:"scene_id"`
	FlowChart  	string					  `json:"flow_chart"`
	FlowCharts  *WorkflowChart			  `json:"flow_charts"`
	FlowBranch  *WorkflowBranch           `json:"flow_branch"`
	IsDefault   int                       `json:"is_default"`
	Range1      string                    `json:"range1"`
	Range2      string                    `json:"range2"`
	Remark      string                    `json:"remark"`
	UpdateTime  time.Time                 `json:"update_time"`
	GroupName   string                    `json:"group_name"`
}

type WorkflowChart struct {
	FirstActionId	string				`json:"first_action_id"`
	LastActionId	string				`json:"last_action_id"`
	HashCondition	bool				`json:"has_condition"`
	ActionMap		map[string]*Action	`json:"actions"`
}

type Action struct {
	ActionType		string				`json:"action_type"`
	ActionId		string				`json:"action_id"`
	ActionName		string				`json:"action_name"`
	Params			[]*Param			`json:"params"`
	NextActionIds	[]string			`json:"next_action_ids"`
	NextConditions	[]string			`json:"next_conditions"`
	PrevActionIds	[]string			`json:"prev_action_ids"`
	Timeout			int64				`json:"timeout"`
	TimeoutAsync	bool				`json:"timeout_async"`
	TimeoutDynamic	bool				`json:"timeout_dynamic"`
	RefWorkflowId	int64				`json:"ref_workflow_id"`
	Description		string				`json:"description"`
}

type Param struct {
	Name			string				`json:"name"`
	Value			string				`json:"value"`
	Type			string				`json:"type"`
}

type WorkflowBranch struct {
	SortedBranch  []string           `json:"sorted_branch"`
	CurrentBranch map[string]int     `json:"current_branch"`
	ActionMap     map[string]*Action `json:"actions"`
}

func (w *Workflow) GetWorkflowChart() *WorkflowChart {
	if w.FlowCharts.HashCondition {
		return w.FlowCharts.clone()
	}
	return w.FlowCharts
}

func NewWorkflowChart(flowChartStr string) (*WorkflowChart, error) {
	if len(flowChartStr) == 0 {
		return nil, errors.New("create WorkflowChart fail, flowChartStr is empty")
	}

	flowChart := &WorkflowChart{}
    //读取的数据为json格式，需要进行解码
    err := json.Unmarshal([]byte(flowChartStr), flowChart)
    if err != nil {
        return nil, fmt.Errorf("create WorkflowChart fail, invalid json:%v, err:%v", flowChartStr, err)
    }

    err = flowChart.setFirstActionId()
    if err != nil {
    	return nil, fmt.Errorf("create WorkflowChart fail, err:%v", err)
    }

	flowChart.setPrevActionIds(flowChart.FirstActionId)

    return flowChart, nil
}

func (w *WorkflowChart) setPrevActionIds(actionId string){
	action, ok := w.ActionMap[actionId]
	if !ok {
		return
	}

	if actionId == w.FirstActionId {
		action.PrevActionIds = []string{}
	}

	nextActionIds := action.NextActionIds
	for _, nextActionId := range nextActionIds {
		nextAction, ok := w.ActionMap[nextActionId]
		if !ok {
			//报个error？
			continue
		}

		if nextAction.PrevActionIds == nil || len(nextAction.PrevActionIds) == 0 {
			nextAction.PrevActionIds = []string{actionId}
		}else{
			isAlreadyAdd := false
			var prevActionId string
			for _, prevActionId = range nextAction.PrevActionIds {
				if prevActionId == actionId {
					isAlreadyAdd = true
					break
				}
			}
			if !isAlreadyAdd {
				nextAction.PrevActionIds = append(nextAction.PrevActionIds, actionId)
			}
		}

		w.setPrevActionIds(nextActionId)
	}
}

/*
	耗时搜索
*/
func (w *WorkflowChart) setFirstActionId() error {
	nextActionIds := make(map[string]bool)
	for _, action := range w.ActionMap {
		if len(action.NextActionIds) <= 0 {
			w.LastActionId = action.ActionId
			continue
		}

		for _, nextId := range action.NextActionIds {
			nextActionIds[nextId] = true
		}
	}

	for actionId,_ := range w.ActionMap {
		if _, ok := nextActionIds[actionId]; !ok {
			w.FirstActionId = actionId
			return nil
		}
	}

	return errors.New("first action not found")
}

func (w *WorkflowChart) CreateWaitMap() (map[string]*sync.WaitGroup, map[string]*TimeWaiter) {
	wgMap:= make(map[string]*sync.WaitGroup)
	tsMap:= make(map[string]*TimeWaiter)
	for actionId, action := range w.ActionMap {
		prevCount := len(action.PrevActionIds)
		if prevCount >1 {
			wg := &sync.WaitGroup{}
			wg.Add(prevCount)
			wgMap[actionId]= wg
		}

		if action.Timeout > 0 {
			tsMap[actionId] = NewTimeWaiter(action.Timeout)
		}
	}

	return wgMap,tsMap
}

func (p *Param) clone() *Param {
	return &Param{
		Name:	p.Name,
		Value:	p.Value,
		Type:	p.Type,
	}
}

func (a *Action) createParamSlice(paramValues *sync.Map) ([]interface{}, error) {
	pCount := len(a.Params)
	if pCount == 0 {
		return make([]interface{}, 0), nil
	}

	p := make([]interface{}, pCount)
	for idx, param := range a.Params {
		//获取实际值
		str := param.Value
		if strings.HasPrefix(param.Value,"$") {
			paramValue, _ := paramValues.Load(str[1:])
			str = fmt.Sprintf("%v", paramValue)
		}

		//获取参数
		if param.Type == "string" {
			p[idx] = str
		}else if param.Type == "int" {
			val, err := strconv.ParseInt(str, 10, 64)
			if err != nil {
				return nil, err
			}
			p[idx] = val
		}else if param.Type == "float" {
			val, err := strconv.ParseFloat(str, 64)
			if err != nil {
				return nil, err
			}
			p[idx] = val
		}else if param.Type == "bool" {
			var val bool
			strl := strings.ToLower(str)
			if strl == "true" {
				val = true
			}else if strl == "false" {
				val = false
			}else {
				return nil, fmt.Errorf("invalid bool value, it must be: true or false")
			}
			p[idx] = val
		}else if param.Type == "interface" {
			var val interface{} = str
			p[idx] = val
		}else{
			//TODO ZYF 先暂时支持4种最常见的类型
			return nil, fmt.Errorf("unknown param type:%v", param.Type)
		}
	}

	return p, nil
}

func (a *Action) Detach(prevAction *Action) {
	if prevAction == nil || len(a.PrevActionIds)==0 {
		//todo error
		return
	}

	prevId := -1
	for pId, prevActionId := range a.PrevActionIds {
		if prevActionId == prevAction.ActionId {
			prevId = pId
			break
		}
	}
	if prevId > -1{
		a.PrevActionIds = append(a.PrevActionIds[:prevId], a.PrevActionIds[prevId+1:]...)
	}

	nextId := -1
	for nId, nextActionId := range prevAction.NextActionIds {
		if nextActionId == a.ActionId {
			nextId = nId
			break
		}
	}
	//fmt.Println("nextId============>", nextId)
	//fmt.Println("before set prevAction.NextActionIds===>", strings.Join(prevAction.NextActionIds,","))
	if nextId>-1 {
		prevAction.NextActionIds = append(prevAction.NextActionIds[:nextId], prevAction.NextActionIds[nextId+1:]...)
		//fmt.Println("after set prevAction.NextActionIds===>", strings.Join(prevAction.NextActionIds,","))
		if len(prevAction.NextConditions)> nextId {
			prevAction.NextConditions = append(prevAction.NextConditions[:nextId], prevAction.NextConditions[nextId+1:]...)
		}
		//fmt.Println("after set prevAction.NextConditions===>", strings.Join(prevAction.NextConditions,","))
	}
	//fmt.Println("\nfinish detach, prevActionId, actionId:", prevAction.ActionId, a.ActionId)
}

func (a *Action) toString() string {
	return fmt.Sprintf("actionId:%+v, nextActionIds:%+v, nextConditions:%+v, prevActionIds:%+v",a.ActionId, strings.Join(a.NextActionIds,","),strings.Join(a.NextConditions,","))
}

func (a *Action) clone() *Action {
	var params []*Param
	if a.Params != nil {
		params = make([]*Param, len(a.Params))
		for i, param := range a.Params {
			params[i] = param.clone()
		}
	}

	act:= &Action{
		ActionType:		a.ActionType,
		ActionId:		a.ActionId,
		ActionName:		a.ActionName,
		Params: 		params,
		NextActionIds:	copyStringArray(a.NextActionIds),
		NextConditions:	copyStringArray(a.NextConditions),
		PrevActionIds:	copyStringArray(a.PrevActionIds),
		Timeout: 		a.Timeout,
		TimeoutAsync:   a.TimeoutAsync,
		TimeoutDynamic: a.TimeoutDynamic,
		RefWorkflowId:  a.RefWorkflowId,
		Description:	a.Description,
	}
	return act
}

func (w *WorkflowChart) clone() *WorkflowChart {
	chart := &WorkflowChart{}
	chart.FirstActionId = w.FirstActionId
	chart.LastActionId  = w.LastActionId
	chart.HashCondition = w.HashCondition

	actionMap := make(map[string]*Action)
	for actionId, action := range w.ActionMap {
		actionMap[actionId] = action.clone()
	}
	chart.ActionMap = actionMap

	return chart
}

func copyStringArray(sources []string) []string {
	if sources == nil {
		return nil
	}

	dests:=make([]string, len(sources))
	for i, source := range sources {
		dests[i] = source
	}
	return dests
}

func (a *Action) executeCond(paramValues *sync.Map) (retActionId string, err error) {
	//fmt.Println("\n\n\n开始执行executeCond:",fmt.Sprintf("action:%+v", a))
	//若报错，有缺省用缺省,无缺省用最后一个。条件必含后继，配置保存时校验
	defaultIndex := len(a.NextActionIds) - 1
	for idx, cdt := range a.NextConditions {
		if cdt == defaultBranch {
			defaultIndex = idx
		}
	}
	//fmt.Println("defaultIndex=======>", defaultIndex)

	retActionId = a.NextActionIds[defaultIndex]
	//fmt.Println("retActionId=>", retActionId)
	err = nil
	defer func() {
		if err0 := recover(); err0 != nil {
			err = fmt.Errorf("executeCond error, default value: %v used, a.Params:%v, paramValues:%v, err0:%v", a.NextActionIds[defaultIndex], a.Params, paramValues, err0)
		}
	}()

	params, err := a.createParamSlice(paramValues)
	if err != nil {
		err = fmt.Errorf("createParamSlice error, default value: %v used, a.Params:%v, paramValues:%v, err:%v", a.NextActionIds[defaultIndex], a.Params, paramValues, err)
		return
	}
	//fmt.Println("prepare to exe:", a.ActionName, params)
	val, err := GetCondExecutors().Execute(a.ActionName, params)
	//fmt.Println("条件执行结果 result: val====>", val, "err====>", err)
	if err != nil {
		err = fmt.Errorf("GetCondExecutors().Execute error, default value: %v used, actionName:%v, params:%v, err:%v", a.NextActionIds[defaultIndex], a.ActionName, params, err)
		return
	}

	for idx, cdt := range a.NextConditions {
		if cdt == val {
			retActionId = a.NextActionIds[idx]
			err = nil
			return
		}
	}

	err = fmt.Errorf("no matched value error, default value:%v used, a.ActionId:%v, execute result:%v, nextActionIds:%v", a.NextActionIds[defaultIndex], a.ActionId, val, a.NextConditions)
	return
}

func (w *WorkflowChart) newWorkflowBranch() *WorkflowBranch {
	sortedArr := []string{}
	aMap := make(map[string]*Action)
	currentMap := make(map[string]int)
	for _, act := range w.ActionMap {
		if act.ActionType == ActionTypeCond {
			sortedArr = append(sortedArr, act.ActionId)
			currentMap[act.ActionId] = 0
			aMap[act.ActionId] = act.clone()
		}
	}

	sort.Strings(sortedArr)
	return &WorkflowBranch{
		SortedBranch:  sortedArr,
		CurrentBranch: currentMap,
		ActionMap:     aMap,
	}
}

func (w *WorkflowBranch) getCurrentBranchKey() string {
	return w.getBranchKey(w.CurrentBranch)
}

func (w *WorkflowBranch) getBranchKey(branchMap map[string]int) string {
	if len(w.SortedBranch) ==0 || len(branchMap) == 0 {
		return BranchKeyDefault
	}

	var buffer bytes.Buffer
	for _, sb := range w.SortedBranch {
		buffer.WriteString(sb)
		buffer.WriteString(BranchKeyJoiner)
		buffer.WriteString(strconv.Itoa(branchMap[sb]))
		buffer.WriteString(BranchKeyJoiner)
	}
	return buffer.String()
}

func (w *WorkflowBranch) currentIndex(actionId string) int {
	//实际currentBranch 不会为nil, 防一下
	if w.CurrentBranch == nil {
		return 0
	}

	return w.CurrentBranch[actionId]
}

func (w *WorkflowBranch) hasBranch() bool {
	if w.SortedBranch == nil || len(w.SortedBranch) ==0 {
		return false
	}

	return true
}

func (w *WorkflowBranch) nextStep() bool {
	for idx, branchId := range w.SortedBranch {
		if w.CurrentBranch[branchId] < len(w.ActionMap[branchId].NextActionIds)-1 {
			w.CurrentBranch[branchId] += 1
			return true
		}

		for i:=0;i<=idx;i++ {
			branchId := w.SortedBranch[i]
			w.CurrentBranch[branchId] = 0
		}
	}

	return false
}
