/**
	Description:	请求上下文，存放workflow处理过程中间结果
	Author:			dayunzhangyunfeng@didiglobal.com
	Date:			2018-07-13
**/

package model

import (
	"context"
	"encoding/json"
	"errors"
	trace "git.xiaojukeji.com/lego/context-go"
	"sync"
)

const (
	KEY_TIMEOUT_ACTION = "timeout_action"
)

type ModuleResultInfo struct {
	Id           string
	StrategyName string
	CostTime     int64
	ResultInfo   interface{}
}

type StrategyContext struct {
	AppId           int64
	AppName         string
	SceneId         int64
	FlowId          int64
	IsLimited       bool
	IsDebug         bool
	CtxTrace        *trace.DefaultTrace
	UserId          string
	Phone           string
	ContextMap      *sync.Map
	ActionResultMap map[string]map[string]string
	contextMutex    sync.Mutex

	moduleResultMap *sync.Map
	errMap          *sync.Map
	ErrNo           int32
	ErrMsg          string

	skipFlag  bool
	skipMutex sync.Mutex

	ApolloInfo *ApolloConfig
}

func NewStrategyContext(ctx context.Context) *StrategyContext {
	strategyContext := new(StrategyContext)
	strategyContext.ContextMap = &sync.Map{}
	strategyContext.ActionResultMap = make(map[string]map[string]string)
	strategyContext.moduleResultMap = &sync.Map{}
	strategyContext.errMap = &sync.Map{}
	if ctxTrace, ok := trace.GetCtxTrace(ctx); ok {
		strategyContext.CtxTrace = ctxTrace
	} else {
		strategyContext.CtxTrace = trace.NewDefaultTrace()
	}
	return strategyContext
}

func (this *StrategyContext) Set(key string, value interface{}) {
	this.ContextMap.Store(key,value)
}

func (this *StrategyContext) SetDebug(actionName string, key string, value interface{}) {
	if !this.IsDebug {
		return
	}
	this.contextMutex.Lock()
	var debugStr string
	if resultJson, err := json.Marshal(value); err == nil {
		debugStr = string(resultJson)
	}
	if _, ok := this.ActionResultMap[actionName]; !ok {
		this.ActionResultMap[actionName] = make(map[string]string)
	}
	this.ActionResultMap[actionName][key] = debugStr
	this.contextMutex.Unlock()
}

func (this *StrategyContext) Get(key string) interface{} {
	value, _ := this.ContextMap.Load(key)
	return value
}

func (this *StrategyContext) Skip(errNo int32, errMsg string) {
	this.skipMutex.Lock()
	this.ErrNo = errNo
	this.ErrMsg = errMsg
	this.skipFlag = true
	this.skipMutex.Unlock()
}

func (this *StrategyContext) IsSkip() bool {
	return this.skipFlag
}

func (this *StrategyContext) AddToArray(arrayKey string, value interface{}) {
	var itArray []interface{}
	if itr, ok := this.ContextMap.Load(arrayKey); ok {
		itArray = itr.([]interface{})
		itArray = append(itArray, value)
	} else {
		itArray = make([]interface{}, 0, 10)
		itArray = append(itArray, value)
	}
	this.ContextMap.Store(arrayKey, itArray)
}

func (this *StrategyContext) GetArray(arrayKey string) ([]interface{}, error) {
	if itf, ok := this.ContextMap.Load(arrayKey); ok {
		return itf.([]interface{}), nil
	}
	return nil, errors.New("No such key:" + arrayKey)
}

func (this *StrategyContext) AddTimeoutAction(actionName string) {
	this.AddToArray(KEY_TIMEOUT_ACTION, actionName)
}

func (this *StrategyContext) GetTimeoutActions() []string {
	itf, err := this.GetArray(KEY_TIMEOUT_ACTION)
	if err != nil || len(itf)==0 {
		return nil
	}

	actionNames := make([]string, len(itf))
	for i, act := range itf {
		actionNames[i] = act.(string)
	}
	return actionNames
}

func (this *StrategyContext) SetModuleResult(actionId, actionName string, costTime int64, resultInfo interface{}) {
	mri := &ModuleResultInfo{
		Id:           actionId,
		StrategyName: actionName,
		CostTime:     costTime,
		ResultInfo:   resultInfo,
	}

	this.moduleResultMap.Store(actionId, mri)
}

func (this *StrategyContext) GetModuleResultMap() *sync.Map {
	return this.moduleResultMap
}

func (this *StrategyContext) SetError(actionId string, err error) {
	this.errMap.Store(actionId, err)
}

func (this *StrategyContext) GetErrorMap() *sync.Map {
	return this.errMap
}
