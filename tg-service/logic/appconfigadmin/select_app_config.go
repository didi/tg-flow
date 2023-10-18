package appconfigadmin

import (
	"context"
	"database/sql"
	"github.com/didi/tg-flow/tg-core/common/mysql"
	"github.com/didi/tg-flow/tg-core/common/tlog"
	"tg-service/common/logs"
	"tg-service/common/template"
	"tg-service/idl"
	"tg-service/logic"
	"time"
)

//查询AppConf数据
func SelectAppConf(selectData *idl.AppConfigInfo, pageLimit int64, pageNum int64, useLimit bool) []*idl.AppConfigInfo {
	sqlStr := "select id, app_name, machine_room, node_name, git_url, operator, create_time, update_time from app_config"
	tempSql := " order by id desc"

	if selectData.Id > 0 { //按条件查询
		sqlStr += " where id = ?" + tempSql
	} else { //全量️查询
		sqlStr += tempSql
	}

	if useLimit {
		sqlStr = logic.AddSelectLimit(sqlStr, pageLimit, pageNum)
	}

	return SelectAppConfPre(sqlStr, selectData)
}

//查询数据库
func SelectAppConfPre(sqlStr string, selectData *idl.AppConfigInfo) []*idl.AppConfigInfo {
	stmt, err := GetAppConfStmt(sqlStr)
	if err != nil {
		return make([]*idl.AppConfigInfo, 0)
	}
	return QueryAppConf(stmt, selectData)
}

func GetAppConfStmt(sqlStr string) (*sql.Stmt, error) {
	stmt, err := mysql.Handler.Prepare(sqlStr)
	if err != nil {
		tlog.Handler.Errorf(context.TODO(), logs.DLTagLoadDBFail, "etype=appconfigadmin.GetAppConfStmt||sql=%v||err=%v", sqlStr, err)
	}
	return stmt, err
}

//数据库查询
func QueryAppConf(stmt *sql.Stmt, selectData *idl.AppConfigInfo) []*idl.AppConfigInfo {
	defer stmt.Close()
	var rows *sql.Rows
	var err error
	if selectData.Id > 0 { //按条件查询
		rows, err = stmt.Query(selectData.Id)
	} else {
		rows, err = stmt.Query()
	}

	if err != nil {
		tlog.Handler.Errorf(context.TODO(), logs.DLTagLoadDBFail, "etype=appconfigadmin.Query||stmt query is fail||err=%v", err)
		return make([]*idl.AppConfigInfo, 0)
	}

	err = rows.Err()
	if err != nil {
		tlog.Handler.Errorf(context.TODO(), logs.DLTagLoadDBFail, "etype=appconfigadmin.RowsErr||respones row is err||err=%v", err)
		return make([]*idl.AppConfigInfo, 0)
	}

	return parasRowAppConf(rows)
}

func parasRowAppConf(rows *sql.Rows) []*idl.AppConfigInfo {
	defer rows.Close()
	appConfigList := make([]*idl.AppConfigInfo, 0)
	for rows.Next() {
		var appConfig = new(idl.AppConfigInfo)
		var createTime time.Time
		var updateTime time.Time
		err := rows.Scan(&appConfig.Id, &appConfig.AppName, &appConfig.MachineRoom, &appConfig.NodeName, &appConfig.GitUrl, &appConfig.Operator, &createTime, &updateTime)
		if err != nil {
			tlog.Handler.Errorf(context.TODO(), logs.DLTagLoadDBFail, "etype=appconfigadmin.RowsScan||rows=%v||err=%v", rows, err)
			return appConfigList
		}
		//时间格式转化
		appConfig.CreateTime = createTime.Format(logic.TIME_FORMAT)
		appConfig.UpdateTime = updateTime.Format(logic.TIME_FORMAT)

		appConfigList = append(appConfigList, appConfig)
	}

	return appConfigList
}

//更新系统id和对应中文映射关系
func UpdateAppMenuCacheData(appConfigList []*idl.AppConfigInfo) {
	logic.Mutex.Lock()
	var appNameMap map[string]int = make(map[string]int)
	var appIdMap map[int]string = make(map[int]string)
	for _, appConfig := range appConfigList {
		appNameMap[appConfig.AppName] = appConfig.Id
		appIdMap[appConfig.Id] = appConfig.AppName
	}
	if len(appConfigList) > 0 {
		template.AppNameMap = appNameMap
		template.AppIdMap = appIdMap
	}
	logic.Mutex.Unlock()
}
