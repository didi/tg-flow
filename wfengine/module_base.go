/**
Description : ModelBase interface define
Author		: dayunzhangyunfeng@didiglobal.com
Date		: 2021-05-14
*/

package wfengine

import (
	"context"
	"fmt"
	"github.com/didi/tg-flow/common/tlog"
	"github.com/didi/tg-flow/model"
	"reflect"
	"strconv"
)

type IModelBase interface {
	DoAction(context.Context, *model.StrategyContext) interface{}
	OnTimeout(context.Context, *model.StrategyContext)
	SetName(string)
	GetName() string
}

type ModelBase struct {
	IModelBase
	Name string
}

func (m *ModelBase) DoAction(context.Context, *model.StrategyContext) interface{} {
	return nil
}

func (m *ModelBase) OnTimeout(context.Context, *model.StrategyContext) {
	//nothing to do here, you can override this function in inherited struct as you need
}

func (m *ModelBase) SetName(name string) {
	m.Name = name
}

func (m *ModelBase) GetName() string {
	return m.Name
}

type ModuleObjBase interface {
	NewObj(moduleName string) IModelBase
}

func reflectModuleField(obj interface{}, vMap map[string]string) error {
	if len(vMap) == 0 {
		return nil
	}

	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr && !v.Elem().CanSet() {
		err := fmt.Errorf("this obj is not match reflect")
		tlog.Handler.ErrorCount(context.TODO(), "ReflectModuleField_err", fmt.Sprintf("obj:%v, err:%v", obj, err))
		return err
	}

	reflectType := reflect.Indirect(v).Type()
	v = v.Elem()
	for i := 0; i < reflectType.NumField(); i++ {
		field := reflectType.Field(i)
		fieldValue := v.FieldByName(field.Name)
		if !fieldValue.IsValid() {
			err := fmt.Errorf("this obj(" + fmt.Sprintf("%v", obj) + ") field(" + field.Name + ")")
			tlog.Handler.ErrorCount(context.TODO(), "ReflectModuleField_err", fmt.Sprintf("%v", err))
			continue
		}

		if vMap[field.Name] == "" { //如果数据库该节点，没有设置参数，则不进行字段的反射赋值，使用初始化时默认值(未设置，则为该类型的自有默认值)
			continue
		}

		if fieldValue.Kind() == reflect.String {
			fieldValue.SetString(vMap[field.Name])
		} else if fieldValue.Kind() == reflect.Int64 || fieldValue.Kind() == reflect.Int32 || fieldValue.Kind() == reflect.Int {
			var tempInt int64
			tempInt, err := strconv.ParseInt(vMap[field.Name], 10, 64)
			if err != nil {
				err := fmt.Errorf("obj(" + fmt.Sprintf("%v", obj) + ") field(" + field.Name + ")'s value(" + vMap[field.Name] + ")")
				tlog.Handler.ErrorCount(context.TODO(), "ReflectModuleField_err", fmt.Sprintf("%v", err))
				continue
			}
			fieldValue.SetInt(tempInt)
		} else if fieldValue.Kind() == reflect.Float64 || fieldValue.Kind() == reflect.Float32 {
			var tempFolat float64
			tempFolat, err := strconv.ParseFloat(vMap[field.Name], 64)
			if err != nil {
				err := fmt.Errorf("obj(" + fmt.Sprintf("%v", obj) + ") field(" + field.Name + ")'s value(" + vMap[field.Name] + ")")
				tlog.Handler.ErrorCount(context.TODO(), "ReflectModuleField_err", fmt.Sprintf("%v", err))
				continue
			}
			fieldValue.SetFloat(tempFolat)
		}
	}
	return nil
}
