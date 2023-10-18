package workflowadmin

import (
	"context"
	"github.com/didi/tg-flow/tg-core/common/mysql"
	"github.com/didi/tg-flow/tg-core/common/tlog"
	"tg-service/common/logs"
	"tg-service/idl"
	"tg-service/logic"
	"time"
)

//删除数据
func DeleteWorkflowConfig(deleteData *idl.WorkFlowConfig) *idl.ResponseInfo {
	responseInfo := &idl.ResponseInfo{
		Tag:    true,
		ErrMsg: "",
	}

	//1、查询对应主流量信息
	workFlow, _, _, _, err := SelectDBMainAndThisWorkFlow(deleteData)
	if err != nil {
		responseInfo.Tag = false
		responseInfo.ErrMsg = "数据删除失败:" + err.Error()
		return responseInfo
	}

	//2、更新主流量range1和range2区间字段值
	var tempRange string
	if len(deleteData.Proportion) <= 0 || deleteData.Proportion == "" {
		tempRange = workFlow.Proportion
	} else {
		//因为SelectDBMainAndThisWorkFlow方法中，将这条被删除记录的Proportion属性值(请求中是百分数)，用数据库中原始的占比离散数字字符串替换了
		//所以，此处可以直接将deleteData.Proportion中的占比离散数字，直接追加到主流量中
		tempRange = deleteData.Proportion + "," + workFlow.Proportion
	}

	//3、删除操作后，对流量占比的校验
	err = logic.CheckDeleteOP(tempRange)
	if err != nil {
		responseInfo.Tag = false
		responseInfo.ErrMsg = err.Error()
		return responseInfo
	}

	nowTime := time.Now().Format("2006-01-02 15:04:05")

	//4、需要对range1的值按从小到大顺序排列
	range1 := logic.SortStr(tempRange)

	//5、对range1生成对应的区间表达式
	range2 := logic.CreateRangeStr(range1)

	//6、先更新主流量,如果失败，则不需要执行后续的删除操作
	updateSql := "update workflow set range1 = ?,range2 = ?,update_time =?,operator=?,manual_slot_ids=? where id = ?"
	_, updateErr := mysql.Handler.Exec(updateSql, range1, range2, nowTime, deleteData.Operator, workFlow.ManualSlotIds, workFlow.WorkFlowId)
	if updateErr != nil {
		responseInfo.Tag = false
		responseInfo.ErrMsg = updateErr.Error()
		tlog.Handler.Errorf(context.TODO(), logs.DLTagLoadDBFail, "etype=workflowadmin.DeleteWorkflowConfig||sql=%v||err=%v", updateSql, err)
		return responseInfo
	}

	//7、删除该条副流量实验
	deleteSql := "delete from workflow where id = ?"
	_, deleteErr := mysql.Handler.Exec(deleteSql, deleteData.WorkFlowId)
	if deleteErr != nil {
		responseInfo.Tag = false
		responseInfo.ErrMsg = deleteErr.Error()
		tlog.Handler.Errorf(context.TODO(), logs.DLTagLoadDBFail, "etype=workflowadmin.DeleteWorkflowConfig||sql=%v||err=%v", deleteSql, err)
	}
	return responseInfo
}
