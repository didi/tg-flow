package idl

import (
	"errors"
	"fmt"
	"github.com/didi/tg-flow/tg-core/wfengine"
	"reflect"
)

type WorkflowExport struct {
	Id         int64          `json:"id"`
	SceneId    int64          `json:"scene_id"`
	FlowCharts *WorkflowChart `json:"flow_charts"`
	SceneName  string         `json:"scene_name"`
	GroupName  string         `json:"group_name"`
}

type Workflow struct {
	Id         int64                     `json:"id"`
	SceneId    int64                     `json:"scene_id"`
	FlowCharts map[string]*WorkflowChart `json:"flow_charts"`
	SceneName  string                    `json:"scene_name"`
}

type WorkflowChart struct {
	ActionMap map[string]*Action `json:"actions"`
}

type Action struct {
	ActionType     string            `json:"action_type"`
	ActionId       string            `json:"action_id"`
	ActionName     string            `json:"action_name"`
	Params         []*wfengine.Param `json:"params"`
	NextActionIds  []string          `json:"next_action_ids,omitempty"`
	NextConditions []string          `json:"next_conditions,omitempty"`
	Description    string            `json:"description"`
	Timeout        int               `json:"timeout"`
	RefWorkflowId  int               `json:"ref_workflow_id"`
	TimeoutAsync   bool              `json:"timeout_async"`
	TimeoutDynamic bool              `json:"timeout_dynamic"`
	Location       string            `json:"location"`
}

func SimpleCopyProperties(dst, src interface{}) (err error) {
	// 防止意外panic
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprintf("%v", e))
		}
	}()

	dstType, dstValue := reflect.TypeOf(dst), reflect.ValueOf(dst)
	srcType, srcValue := reflect.TypeOf(src), reflect.ValueOf(src)

	// dst必须结构体指针类型
	if dstType.Kind() != reflect.Ptr || dstType.Elem().Kind() != reflect.Struct {
		return errors.New("dst type should be a struct pointer")
	}

	// src必须为结构体或者结构体指针，.Elem()类似于*ptr的操作返回指针指向的地址反射类型
	if srcType.Kind() == reflect.Ptr {
		srcType, srcValue = srcType.Elem(), srcValue.Elem()
	}
	if srcType.Kind() != reflect.Struct {
		return errors.New("src type should be a struct or a struct pointer")
	}

	// 取具体内容
	dstType, dstValue = dstType.Elem(), dstValue.Elem()

	// 属性个数
	propertyNums := dstType.NumField()

	for i := 0; i < propertyNums; i++ {
		// 属性
		property := dstType.Field(i)
		// 待填充属性值
		propertyValue := srcValue.FieldByName(property.Name)

		// 无效，说明src没有这个属性 || 属性同名但类型不同
		if !propertyValue.IsValid() || property.Type != propertyValue.Type() {
			continue
		}

		if dstValue.Field(i).CanSet() {
			dstValue.Field(i).Set(propertyValue)
		}
	}

	return nil
}
