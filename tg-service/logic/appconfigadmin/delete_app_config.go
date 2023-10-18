package appconfigadmin

import (
	"context"
	"github.com/didi/tg-flow/tg-core/common/mysql"
	"github.com/didi/tg-flow/tg-core/common/tlog"
	"tg-service/common/logs"
	"tg-service/idl"
)

//删除数据
func DeleteAppConfig(deleteData *idl.AppConfigInfo) *idl.ResponseInfo {
	responseInfo := &idl.ResponseInfo{
		Tag:    true,
		ErrMsg: "",
	}
	sql := "delete from app_config where id = ?"
	_, err := mysql.Handler.Exec(sql, deleteData.Id)
	if err != nil {
		responseInfo.Tag = false
		responseInfo.ErrMsg = err.Error()
		tlog.Handler.Errorf(context.TODO(), logs.DLTagLoadDBFail, "etype=appconfigadmin.DeleteSystemConfig||sql=%v||err=%v", sql, err)
	}

	return responseInfo
}
