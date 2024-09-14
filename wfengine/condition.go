/**
	Description : condition for branch
	Author		: dayunzhangyunfeng@didiglobal.com
	Date		: 2021-05-14
*/

package wfengine

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

type Condition struct {}

const (
	YES	= "Y"
	NO	= "N"
)

type CondExecutors struct {
	InnerExecutor	*reflect.Value
	OuterExecutors	*sync.Map
}

var condExecutors *CondExecutors
var once sync.Once
func GetCondExecutors() *CondExecutors {
	once.Do(func(){
		condExecutors = newCondExecutors()
	})
	return condExecutors
}

func newCondExecutors() *CondExecutors {
	var cdt interface{} = &Condition{}
	innerExecutor := reflect.ValueOf(cdt)

	condExecutors = &CondExecutors {
		InnerExecutor	:	&innerExecutor,
		OuterExecutors	:	&sync.Map{},
	}

	return condExecutors
}

func (c CondExecutors) RegisterCondExecutor(conditionName string, cdt interface{}){
	cde := reflect.ValueOf(cdt)
	c.OuterExecutors.Store(conditionName, &cde)
}

func (c CondExecutors) Execute(actionName string, paramValues []interface{}) (string, error) {
	var executor *reflect.Value
	var methodName string
	idx := strings.Index(actionName, ".")
	if idx < 0 {
		executor = c.InnerExecutor
		methodName = actionName
	}else{
		structName := actionName[:idx]
		executorItf, ok := c.OuterExecutors.Load(structName)
		if !ok {
			return "", fmt.Errorf("struct %v not Registered", structName)
		}
		executor, ok = executorItf.(*reflect.Value)
		if !ok || executor == nil {
			return "", fmt.Errorf("struct %v must be *reflect.Value, while value is:%v", structName, executorItf)
		}

		methodName	= actionName[idx+1:]
	}

	m := executor.MethodByName(methodName)
	p := make([]reflect.Value, len(paramValues))
	for idx, param := range paramValues {
		p[idx] = reflect.ValueOf(param)
	}

	rets := m.Call(p)
	var err error
	if rets[1].IsNil() {
		err = nil
	}else{
		err = fmt.Errorf("error:%v", rets[1].String())
	}

	return rets[0].String(), err
}

/**
	比较两个数的大小是否相等，相等返回1，否则返回0
 */
func (c *Condition) EQ(itra, itrb interface{}) (string, error) {
	if fmt.Sprintf("%v",itra) == fmt.Sprintf("%v", itrb) {
		return YES, nil
	}

	return NO, nil
}

/**
	比较两个数的大小是否不相等，不相等返回y，否则返回n
*/
func (c *Condition) NE(itra, itrb interface{}) (string, error) {
	if fmt.Sprintf("%v",itra) != fmt.Sprintf("%v", itrb) {
		return YES, nil
	}

	return NO, nil
}

func (c *Condition) LT(itra, itrb interface{}) (string, error) {
	stra := fmt.Sprintf("%v", itra)
	strb := fmt.Sprintf("%v", itrb)
	fa, erra := strconv.ParseFloat(stra, 64)
	fb, errb := strconv.ParseFloat(strb, 64)
	if erra != nil || errb != nil {
		return "", fmt.Errorf("params must be number, while values are:%v, %v", itra, itrb)
	}

	if fa < fb {
		return YES, nil
	}

	return NO, nil
}

func (c *Condition) GT(itra, itrb interface{}) (string, error) {
	stra := fmt.Sprintf("%v", itra)
	strb := fmt.Sprintf("%v", itrb)
	fa, erra := strconv.ParseFloat(stra, 64)
	fb, errb := strconv.ParseFloat(strb, 64)
	if erra != nil || errb != nil {
		return "", fmt.Errorf("params must be number, while values are:%v, %v", itra, itrb)
	}

	if fa > fb {
		return YES, nil
	}

	return NO, nil
}

func (c *Condition) LE(itra, itrb interface{}) (string, error) {
	stra := fmt.Sprintf("%v", itra)
	strb := fmt.Sprintf("%v", itrb)
	fa, erra := strconv.ParseFloat(stra, 64)
	fb, errb := strconv.ParseFloat(strb, 64)
	if erra != nil || errb != nil {
		return "", fmt.Errorf("params must be number, while values are:%v, %v", itra, itrb)
	}

	if fa <= fb {
		return YES, nil
	}

	return NO, nil
}

func (c *Condition) GE(itra, itrb interface{}) (string, error) {
	stra := fmt.Sprintf("%v", itra)
	strb := fmt.Sprintf("%v", itrb)
	fa, erra := strconv.ParseFloat(stra, 64)
	fb, errb := strconv.ParseFloat(strb, 64)
	if erra != nil || errb != nil {
		return "", fmt.Errorf("params must be number, while values are:%v, %v", itra, itrb)
	}

	if fa >= fb {
		return YES, nil
	}

	return NO, nil
}

func (c *Condition) SW(itra interface{}) (string, error) {
	return fmt.Sprintf("%v",itra), nil
}