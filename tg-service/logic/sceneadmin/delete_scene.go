package sceneadmin

import (
	"context"
	"github.com/didi/tg-flow/tg-core/common/mysql"
	"github.com/didi/tg-flow/tg-core/common/tlog"
	"tg-service/common/logs"
	"tg-service/idl"
)

//删除数据
func DeleteSceneConfig(deleteData *idl.SceneConfig) *idl.ResponseInfo {
	responseInfo := &idl.ResponseInfo{
		Tag:    true,
		ErrMsg: "",
	}

	//删除workflow表配置
	sql := "delete from workflow where experiment_id=?"
	_, err := mysql.Handler.Exec(sql, deleteData.Id)
	if err != nil {
		responseInfo.Tag = false
		responseInfo.ErrMsg = "删除workflow表 该场景的配置 失败！" + err.Error()
		return responseInfo
	}

	//删除rec_scene_config表配置
	//sql = "delete from rec_scene_config where scene_id = ?"
	//_, err = mysql.Handler.Exec(sql, deleteData.Id)
	//if err != nil {
	//	responseInfo.Tag = false
	//		responseInfo.ErrMsg = "删除rec_scene_config表 该场景的配置 失败！" + err.Error()
	//		return responseInfo
	//	}

	//删除dimension表配置
	sql = "delete from dimension where scene_id = ?"
	_, err = mysql.Handler.Exec(sql, deleteData.Id)
	if err != nil {
		responseInfo.Tag = false
		responseInfo.ErrMsg = "删除 dimension表 该场景的配置 失败！" + err.Error()
		return responseInfo
	}

	//删除场景配置
	sql = "delete from scene_config where id=?"
	_, err = mysql.Handler.Exec(sql, deleteData.Id)
	if err != nil {
		responseInfo.Tag = false
		responseInfo.ErrMsg = "删除该场景的配置 失败！" + err.Error()
		tlog.Handler.Errorf(context.TODO(), logs.DLTagLoadDBFail, "etype=sceneadmin.DeleteSceneConfig||sql=%v||err=%v", sql, err)
	}
	return responseInfo
}
