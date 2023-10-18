package sceneadmin

import (
	"context"
	"fmt"
	"github.com/didi/tg-flow/tg-core/common/mysql"
	"github.com/didi/tg-flow/tg-core/common/tlog"
	"strconv"
	"tg-service/common/logs"
	"tg-service/common/template"
	"tg-service/idl"
	"time"
	"unicode"
)

//新增 数据
func AddAddOrUpdateSceneConfig(addData *idl.SceneConfig) *idl.ResponseInfo {
	responseInfo := &idl.ResponseInfo{
		Tag:    true,
		ErrMsg: "",
	}

	if len(addData.Name) <= 0 {
		responseInfo.Tag = false
		responseInfo.ErrMsg = "场景名称 不能为空或者空格，无法添加本条数据！"
		return responseInfo
	}

	for _, v := range addData.Name {
		if unicode.Is(unicode.Han, v) {
			responseInfo.Tag = false
			responseInfo.ErrMsg = "场景名称 不支持中文，请输入英文！"
			return responseInfo
		}
	}
	if len(addData.NameZh) > 0 {
		hasZh := false
		for _, zhv := range addData.NameZh {
			if unicode.Is(unicode.Han, zhv) {
				hasZh = true
				break
			}
		}
		if !hasZh {
			responseInfo.Tag = false
			responseInfo.ErrMsg = "中文场景名称不为空时则至少输入一个中文字符！"
			return responseInfo
		}
	}
	nowTime := time.Now().Format("2006-01-02 15:04:05")
	var sql string
	var err error
	//判断数据库中是否已经存在该场景id
	selectData := &idl.SceneConfig{
		Id: addData.Id,
	}
	sceneConfigIdList := SelectSceneConf(selectData, -1, 0, false)
	//判断数据库中是否已经存在该场景名称
	selectData = &idl.SceneConfig{
		Name: addData.Name,
	}
	sceneConfigNameList := SelectSceneConf(selectData, -1, 0, false)

	if len(sceneConfigNameList) > 0 && sceneConfigNameList[0].Id != addData.Id {
		responseInfo.Tag = false
		responseInfo.ErrMsg = "数据库中已存在 其他不同场景编号，但相同场景名称 的记录，无法添加本条数据！"
		return responseInfo
	}

	appId := template.AppNameMap[addData.AppName]

	if addData.OldId <= 0 { //添加数据
		if len(sceneConfigIdList) > 0 {
			responseInfo.Tag = false
			responseInfo.ErrMsg = "数据库中已存在 相同场景编号 的记录，无法添加本条数据！"
			return responseInfo
		}
		sql = "insert into scene_config(id,name,app_id,bucket_type,operator,create_time,update_time,flow_type,name_zh,exp_name) values(?,?,?,?,?,?,?,?,?,?)"
		_, err = mysql.Handler.Exec(sql, addData.Id, addData.Name, appId, addData.BucketType, addData.Operator, nowTime, nowTime, addData.FlowType, addData.NameZh, addData.ExpName)
	} else { //更新数据
		sql = "UPDATE scene_config SET name=?,app_id= ?, bucket_type=?, operator=?, update_time= ?, flow_type=?, name_zh=?, exp_name=? WHERE id=?"
		_, err = mysql.Handler.Exec(sql, addData.Name, appId, addData.BucketType, addData.Operator, nowTime, addData.FlowType, addData.NameZh, addData.ExpName, addData.Id)
	}

	if err != nil {
		responseInfo.Tag = false
		responseInfo.ErrMsg = err.Error()
		tlog.Handler.Errorf(context.TODO(), logs.DLTagLoadDBFail, "etype=sceneadmin.AddAddOrUpdateSceneConfig||sql=%v||err=%v", sql, err)
		return responseInfo
	}

	//如果是新增场景,需要向workflow表中插入一条该场景的主流量，且占比100%
	if addData.OldId <= 0 {
		var range1 string //流量占比
		for i := 0; i < 100; i++ {
			value := strconv.Itoa(i)
			if len(range1) <= 0 {
				range1 = value
			} else {
				range1 += "," + value
			}
		}
		addSql := "insert into workflow(dimension_id, experiment_id, modules,is_default,status,range1,range2,remark,create_time,update_time,operator,flow_chart) values(-1,?,'',1,1,?,'[0-99]','',?,?,?,'')"
		_, addErr := mysql.Handler.Exec(addSql, addData.Id, range1, nowTime, nowTime, addData.Operator)
		if addErr != nil {
			responseInfo.Tag = false
			responseInfo.ErrMsg = fmt.Sprintf("插入数据失败:%v", addErr.Error())
			tlog.Handler.Errorf(context.TODO(), logs.DLTagLoadDBFail, "etype=sceneadmin.AddOrUpdateSceneConfig||sql=%v||err=%v", addSql, err)
		}
	}

	//对于推荐引擎，即appid=1的场景，需要向rec_scene_config表中插入一条该场景和场景组的对应关系
	/*	if appId == 1 {
		recSceneAddData := &idl.RecSceneConfig{
			SceneId:           addData.Id,
			SceneName:         addData.Name,
			Operator:          addData.Operator,
		}
		recSceneConfigList := recalladmin.SelectRecSceneConf(recSceneAddData)
		if len(recSceneConfigList) > 0 {
			recSceneAddData.Id = recSceneConfigList[0].Id
		}
		responseInfo = recalladmin.AddOrUpdateRecSceneConfig(recSceneAddData)
	}*/
	return responseInfo
}
