package appconfigadmin

import (
	"context"
	"github.com/didi/tg-flow/tg-core/common/mysql"
	"github.com/didi/tg-flow/tg-core/common/tlog"
	"tg-service/common/logs"
	"tg-service/idl"
	"time"
)

//新增数据
func AddOrUpdateAppConfig(addData *idl.AppConfigInfo) *idl.ResponseInfo {
	responseInfo := &idl.ResponseInfo{
		Tag:    true,
		ErrMsg: "",
	}

	nowTime := time.Now().Format("2006-01-02 15:04:05")
	var sql string
	var err error

	if addData.AppName == "" {
		responseInfo.Tag = false
		responseInfo.ErrMsg = "系统名称 不能为空，或空串"
		return responseInfo
	}
	if addData.MachineRoom == "" {
		responseInfo.Tag = false
		responseInfo.ErrMsg = "部署机房 不能为空，或空串"
		return responseInfo
	}
	if addData.NodeName == "" {
		responseInfo.Tag = false
		responseInfo.ErrMsg = "节点名称 不能为空，或空串"
		return responseInfo
	}

	//先判断数据库中是否有相同记录
	appConfigList := SelectAppConf(addData, -1, 0, false)

	if addData.OldId <= 0 { //添加数据
		if len(appConfigList) > 0 {
			responseInfo.Tag = false
			responseInfo.ErrMsg = "数据库中已存在 相同app id 的记录，无法添加本条数据！"
			return responseInfo
		}
		sql = "insert into app_config(id, app_name, machine_room, node_name, git_url, operator, create_time, update_time) values(?,?,?,?,?,?,?)"
		_, err = mysql.Handler.Exec(sql, addData.Id, addData.AppName, addData.MachineRoom, addData.NodeName, addData.GitUrl, addData.Operator, nowTime, nowTime)
	} else { //更新数据
		sql = "update app_config set app_name=?,machine_room=?,node_name=?, git_url=?, operator=?,update_time=? where id=?"
		_, err = mysql.Handler.Exec(sql, addData.AppName, addData.MachineRoom, addData.NodeName, addData.GitUrl, addData.Operator, nowTime, addData.Id)
	}
	if err != nil {
		responseInfo.Tag = false
		responseInfo.ErrMsg = err.Error()
		tlog.Handler.Errorf(context.TODO(), logs.DLTagLoadDBFail, "etype=appconfigadmin.AddOrUpdateSystemConfig||sql=%v||err=%v", sql, err)
	}
	return responseInfo
}
