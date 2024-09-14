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
	"github.com/didi/tg-flow/common/timeutils"
	"sync"
)

const (
	KEY_TIMEOUT_ACTION = "timeout_action"
)

type StrategyContext struct {
	AppId     int64
	AppName   string
	SceneId   int64
	FlowId    int64
	IsLimited bool
	IsDebug   bool
	//CtxTrace        *trace.DefaultTrace
	UserId          string
	Phone           string
	GroupName       string
	ContextMap      *sync.Map
	ActionResultMap map[string]map[string]string
	contextMutex    sync.Mutex
	TC              *timeutils.TimeCoster
	errMap          *sync.Map
	ErrNo           int32
	ErrMsg          string

	skipFlag  bool
	skipMutex sync.Mutex
}

func NewStrategyContext(ctx context.Context) *StrategyContext {
	strategyContext := new(StrategyContext)
	strategyContext.ContextMap = &sync.Map{}
	strategyContext.ActionResultMap = make(map[string]map[string]string)
	strategyContext.errMap = &sync.Map{}
	strategyContext.TC = timeutils.NewTimeCosterUnit(timeutils.TimeUnitMillSecond)
	//if ctxTrace, ok := trace.GetCtxTrace(ctx); ok {
	//	strategyContext.CtxTrace = ctxTrace
	//} else {
	//	strategyContext.CtxTrace = trace.NewDefaultTrace()
	//}
	return strategyContext
}

func (s *StrategyContext) Set(key string, value interface{}) {
	s.ContextMap.Store(key, value)
}

func (s *StrategyContext) SetDebug(actionName string, key string, value interface{}) {
	if !s.IsDebug {
		return
	}
	s.contextMutex.Lock()
	var debugStr string
	if resultJson, err := json.Marshal(value); err == nil {
		debugStr = string(resultJson)
	}
	if _, ok := s.ActionResultMap[actionName]; !ok {
		s.ActionResultMap[actionName] = make(map[string]string)
	}
	s.ActionResultMap[actionName][key] = debugStr
	s.contextMutex.Unlock()
}

func (s *StrategyContext) Get(key string) interface{} {
	value, _ := s.ContextMap.Load(key)
	return value
}

func (s *StrategyContext) Skip(errNo int32, errMsg string) {
	s.skipMutex.Lock()
	s.ErrNo = errNo
	s.ErrMsg = errMsg
	s.skipFlag = true
	s.skipMutex.Unlock()
}

func (s *StrategyContext) IsSkip() bool {
	return s.skipFlag
}

func (s *StrategyContext) AddToArray(arrayKey string, value interface{}) {
	var itArray []interface{}
	if itr, ok := s.ContextMap.Load(arrayKey); ok {
		itArray = itr.([]interface{})
		itArray = append(itArray, value)
	} else {
		itArray = make([]interface{}, 0, 10)
		itArray = append(itArray, value)
	}
	s.ContextMap.Store(arrayKey, itArray)
}

func (s *StrategyContext) GetArray(arrayKey string) ([]interface{}, error) {
	if itf, ok := s.ContextMap.Load(arrayKey); ok {
		return itf.([]interface{}), nil
	}
	return nil, errors.New("No such key:" + arrayKey)
}

func (s *StrategyContext) AddTimeoutAction(actionName string) {
	s.AddToArray(KEY_TIMEOUT_ACTION, actionName)
}

func (s *StrategyContext) GetTimeoutActions() []string {
	itf, err := s.GetArray(KEY_TIMEOUT_ACTION)
	if err != nil || len(itf) == 0 {
		return nil
	}

	actionNames := make([]string, len(itf))
	for i, act := range itf {
		actionNames[i] = act.(string)
	}
	return actionNames
}

func (s *StrategyContext) SetError(actionId string, err error) {
	s.errMap.Store(actionId, err)
}

func (s *StrategyContext) GetErrorMap() *sync.Map {
	return s.errMap
}
