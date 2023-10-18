package sceneadmin

import (
	"context"
	"database/sql"
	"github.com/didi/tg-flow/tg-core/common/mysql"
	"github.com/didi/tg-flow/tg-core/common/tlog"
	"tg-service/common/logs"
	"tg-service/common/template"
	"tg-service/idl"
	"tg-service/logic"
	"tg-service/logic/appconfigadmin"

	"time"
)

//查询sceneconfig数据
func SelectSceneConf(selectData *idl.SceneConfig, pageLimit int64, pageNum int64, useLimit bool) []*idl.SceneConfig {
	sqlStr := "select id,name,app_id,bucket_type,operator,create_time,update_time,flow_type,name_zh,exp_name from scene_config"
	tempSql := " order by id desc"
	if selectData.Id > 0 { //按场景编号查询
		sqlStr += " where id = ?" + tempSql
	} else if selectData.Name != "" && selectData.Name != "-1" && selectData.Name != "全部" { //按场景名称查询
		sqlStr += " where name = ?" + tempSql
	} else if selectData.AppName != "" && selectData.AppName != "-1" { //按系统名称查询
		//传过来的是系统名称，但是系统相关信息已经拆分到app_config表，因此，此处需要将名称转化成系统id查询
		sqlStr += " where app_id = ?" + tempSql
	} else if selectData.BucketType > 0 { //按分桶类型查询
		sqlStr += " where bucket_type = ?" + tempSql
	} else { //全量️查询
		sqlStr += tempSql
	}

	if useLimit {
		sqlStr = logic.AddSelectLimit(sqlStr, pageLimit, pageNum)
	}

	return SelectSceneConfPre(sqlStr, selectData)
}

func SelectSceneConfPre(sqlStr string, selectData *idl.SceneConfig) []*idl.SceneConfig {
	//查询数据库
	sceneConfigList := SelectServerDB(sqlStr, selectData)
	return sceneConfigList
}

//数据库查询
func SelectServerDB(sqlStr string, selectData *idl.SceneConfig) []*idl.SceneConfig {
	sceneConfigList := make([]*idl.SceneConfig, 0)
	stmt, err := mysql.Handler.Prepare(sqlStr)
	if err != nil {
		tlog.Handler.Errorf(context.TODO(), logs.DLTagLoadDBFail, "etype=sceneadmin.SelectServerDB||sql=%v||err=%v", sqlStr, err)
		return sceneConfigList
	}
	defer stmt.Close()

	var rows *sql.Rows
	if selectData.Id > 0 { //按场景编号查询
		rows, err = stmt.Query(selectData.Id)
	} else if selectData.Name != "" && selectData.Name != "-1" && selectData.Name != "全部" { //按场景名称查询
		rows, err = stmt.Query(selectData.Name)
	} else if selectData.AppName != "" && selectData.AppName != "-1" { //按系统名称查询
		appId := template.AppNameMap[selectData.AppName]
		rows, err = stmt.Query(appId)
	} else if selectData.BucketType > 0 { //按分桶类型查询
		rows, err = stmt.Query(selectData.BucketType)
	} else {
		rows, err = stmt.Query()
	}

	if err != nil {
		tlog.Handler.Errorf(context.TODO(), logs.DLTagLoadDBFail, "etype=sceneadmin.SelectServerDB||sql=%v||err=%v", sqlStr, err)
		return sceneConfigList
	}
	defer rows.Close()

	for rows.Next() {
		var createTime time.Time
		var updateTime time.Time
		var sceneConfig = new(idl.SceneConfig)
		err := rows.Scan(&sceneConfig.Id, &sceneConfig.Name, &sceneConfig.AppId, &sceneConfig.BucketType, &sceneConfig.Operator, &createTime, &updateTime, &sceneConfig.FlowType, &sceneConfig.NameZh, &sceneConfig.ExpName)
		if err != nil {
			tlog.Handler.Errorf(context.TODO(), logs.DLTagLoadDBFail, "etype=sceneadmin.SceneRowsScan||rows=%v||err=%v", rows, err)
			return sceneConfigList
		}

		//根据appid，查询app_config表，获取appname
		appData := &idl.AppConfigInfo{
			Id: sceneConfig.AppId,
		}
		appInfo := appconfigadmin.SelectAppConf(appData, -1, 0, false)

		if len(appInfo) > 0 {
			sceneConfig.AppName = appInfo[0].AppName
		} else {
			sceneConfig.AppName = "数据异常"
		}

		//时间格式转化
		sceneConfig.CreateTime = createTime.Format(logic.TIME_FORMAT)
		sceneConfig.UpdateTime = updateTime.Format(logic.TIME_FORMAT)

		sceneConfigList = append(sceneConfigList, sceneConfig)
	}
	err = rows.Err()
	if err != nil {
		tlog.Handler.Errorf(context.TODO(), logs.DLTagLoadDBFail, "etype=sceneadmin.SceneRowsErr||respones row is err||err=%v", err)
		return sceneConfigList
	}
	return sceneConfigList
}

//查询数据库中最大的appid
func SelectMaxAppId() int {
	sqlStr := "select MAX(app_id) from scene_config"
	stmt, err := mysql.Handler.Prepare(sqlStr)
	if err != nil {
		tlog.Handler.Errorf(context.TODO(), logs.DLTagLoadDBFail, "etype=sceneadmin.SelectMaxAppId||sql=%v||err=%v", sqlStr, err)
		return -1
	}
	defer stmt.Close()

	var rows *sql.Rows
	rows, err = stmt.Query()

	if err != nil {
		tlog.Handler.Errorf(context.TODO(), logs.DLTagLoadDBFail, "etype=sceneadmin.SelectMaxAppId||sql=%v||err=%v", sqlStr, err)
		return -1
	}
	defer rows.Close()

	var maxAppId int
	for rows.Next() {
		err := rows.Scan(&maxAppId)
		if err != nil {
			tlog.Handler.Errorf(context.TODO(), logs.DLTagLoadDBFail, "etype=sceneadmin.SelectMaxAppIdRowsScan||rows=%v||err=%v", rows, err)
			return -1
		}
	}
	err = rows.Err()
	if err != nil {
		tlog.Handler.Errorf(context.TODO(), logs.DLTagLoadDBFail, "etype=sceneadmin.SelectMaxAppIdRowsErr||respones row is err||err=%v", err)
		return -1
	}
	return maxAppId
}
