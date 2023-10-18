package model

import (
	"time"
	"sync"
	"errors"
	"encoding/json"
	"fmt"
)

type Workflow struct {
	Id            int64      `json:"id"`
	DimensionId   int64      `json:"dimensionId"`
	ExperimentId  int64      `json:"experimentId"`
	ModulesArray  [][]string `json:"modulesArray"`
	FlowChart	  *WorkflowChart  `json:"flow_chart"`
	IsDefault     int        `json:"isDefault"`
	Range1        string     `json:"range1"`
	Range2        string     `json:"range2"`
	Remark        string     `json:"remark"`
	UpdateTime    time.Time  `json:"updateTime"`
	ManualSlotIds []int64    `json:"manual_slot_ids"`
	GroupName	  string	 `json:"group_name"`
}

type WorkflowChart struct {
	FirstActionId	string				`json:"first_action_id"`
	ActionMap	map[string]*Action		`json:"actions"`
	ParamMap	map[string]string		`json:"params"`
}

type Action struct {
	ActionId		string				`json:"action_id"`
	ActionName		string				`json:"action_name"`
	Params			map[string]string	`json:"params"`
	NextActionIds	[]string			`json:"next_action_ids"`
	PrevActionIds	[]string			`json:"prev_action_ids"`
}

func NewWorkflowChart(flowChartStr string) (*WorkflowChart, error) {
	if len(flowChartStr) ==0 {
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
    
    //setFirstActionId first, then setPrevActionIds
    flowChart.setPrevActionIds(flowChart.FirstActionId)
    return flowChart, nil
}

func (this *WorkflowChart) setPrevActionIds(actionId string){
	action, ok := this.ActionMap[actionId]
	if !ok {
		return
	}
	
	//默认第一个action的id为1,没有前驱节点
	if actionId == this.FirstActionId {
		action.PrevActionIds = []string{}
	}
	
	nextActionIds := action.NextActionIds
	for _, nextActionId := range nextActionIds {
		nextAction, ok := this.ActionMap[nextActionId]
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
				}
			}
			if !isAlreadyAdd {
				nextAction.PrevActionIds = append(nextAction.PrevActionIds, actionId)
			}
		}
		
		this.setPrevActionIds(nextActionId)
	}
}

/*
	耗时搜索
*/
func (this *WorkflowChart) setFirstActionId() error {
	nextActionIds := make(map[string]bool)
	for _, action := range this.ActionMap {
		if len(action.NextActionIds) <= 0 {
			continue
		}
		
		for _, nextId := range action.NextActionIds {
			nextActionIds[nextId] = true
		}
	}
	
	for actionId,_ := range this.ActionMap {
		if _, ok := nextActionIds[actionId]; !ok {
			this.FirstActionId = actionId
			return nil
		}
	}
	
	return errors.New("first action not found")
}

func (this *WorkflowChart) CreateWaitMap() map[string]*sync.WaitGroup {
	wgMap:= make(map[string]*sync.WaitGroup)
	for actionId, action := range this.ActionMap {
		prevCount := len(action.PrevActionIds)
		if prevCount >1 {
			wg := &sync.WaitGroup{}
			wg.Add(prevCount)
			wgMap[actionId]= wg
		}
	}
	
	return wgMap
}
