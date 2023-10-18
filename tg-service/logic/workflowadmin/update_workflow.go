package workflowadmin

import (
	"context"
	"github.com/didi/tg-flow/tg-core/common/mysql"
	"github.com/didi/tg-flow/tg-core/common/tlog"
	"strconv"
	"strings"
	"tg-service/common/logs"
	"tg-service/common/template"
	"tg-service/idl"
	"tg-service/logic"
	"tg-service/logic/sceneadmin"
	"time"
)

//更新workflow信息
func UpdateDBWorkFlow(modifyData *idl.WorkFlowConfig, responseInfo *idl.ResponseInfo) {
	var flowType int64
	flowType = -1
	if modifyData.SceneId != 0 {
		sceneData := &idl.SceneConfig{
			Id: modifyData.SceneId,
		}
		sceneConfigList := sceneadmin.SelectSceneConf(sceneData, -1, 0, false)
		if len(sceneConfigList) != 0 {
			flowType = sceneConfigList[0].FlowType
		}
	}
	ctx := context.TODO()
	//1、先将请求记录中的流量占比取出来,第2步中会用数据库中的占比离散数字字符串，替换Proportion属性值
	proportion, _ := strconv.Atoi(modifyData.Proportion)
	//2、查询对应主流量、数据库中该流量的信息和不包括该非主流量的，其他非主流量的总占比 和总占比值
	workFlow, thismodifyWorkFlow, haveUseProportion, _, err := SelectDBMainAndThisWorkFlow(modifyData)
	if err != nil {
		responseInfo.Tag = false
		responseInfo.ErrMsg = "修改数据失败:" + err.Error()
		return
	}
	if workFlow == nil && thismodifyWorkFlow == nil {
		responseInfo.Tag = false
		responseInfo.ErrMsg = "查询数据库失败，导致更新操作失败，请重新更新!"
		return
	}

	if proportion > (100-haveUseProportion) && flowType == 0 {
		responseInfo.Tag = false
		responseInfo.ErrMsg = "修改的占比，超过了可分配占比,可分配占比 = " + strconv.Itoa(100-haveUseProportion)
		return
	}

	nowTime := time.Now().Format("2006-01-02 15:04:05")

	//3、表示要修改的是主流量
	if modifyData.OldWorkFlowId == workFlow.WorkFlowId {
		sql := "update workflow set id=?,modules =?,remark =?,update_time =?,operator=?,manual_slot_ids=?,group_name=?,flow_chart=? where id = ?"
		_, updateMainErr := mysql.Handler.Exec(sql, modifyData.WorkFlowId, modifyData.Modules, modifyData.Remark, nowTime, modifyData.Operator, modifyData.ManualSlotIds, modifyData.GroupName, modifyData.FlowChartJson, modifyData.OldWorkFlowId)
		if updateMainErr != nil {
			responseInfo.Tag = false
			responseInfo.ErrMsg = updateMainErr.Error()
		}
		tlog.Handler.Infof(ctx, logs.DLTagProcessLog, "etype=workflowadmin.UpdateDBWorkFlow||data=%v||err=主流量修改成功", modifyData)
		return
	}

	//4、该流量对应的数据库中的数据
	dbMainWorkFlowRangeArry := strings.Split(workFlow.Proportion, ",")
	dbThismodifyWorkFlowRangeArry := strings.Split(thismodifyWorkFlow.Proportion, ",")

	thisModifyWorkFlowRang1 := ""         //修改后的该条非主流量占比
	MainWorkFlowRang1 := ""               //修改后的主流量调整后的占比
	dbThismodifyWorkFlowRangeArryLen := 0 //该条非主流量，数据库原始占比离散数字个数

	if thismodifyWorkFlow.Proportion != "" {
		dbThismodifyWorkFlowRangeArryLen = len(dbThismodifyWorkFlowRangeArry)
	}

	//5、对流量占比,未做修改
	if dbThismodifyWorkFlowRangeArryLen == proportion {
		MainWorkFlowRang1 = workFlow.Proportion
		thisModifyWorkFlowRang1 = thismodifyWorkFlow.Proportion
	} else if dbThismodifyWorkFlowRangeArryLen < proportion {
		//6、新增了流量占比，新增部分，需要从主流量占比值中随机抽取
		count := proportion - dbThismodifyWorkFlowRangeArryLen
		addRangeStr, rangeStr := ModifyWorkFlowRange(dbMainWorkFlowRangeArry, count)
		MainWorkFlowRang1 = rangeStr

		//将新增部分占比的离散数字，和原始的占比离散数字，按从小到大排序组合在一起
		if len(addRangeStr) <= 0 {
			thisModifyWorkFlowRang1 = logic.SortStr(thismodifyWorkFlow.Proportion)
		} else {
			if thismodifyWorkFlow.Proportion == "" {
				thisModifyWorkFlowRang1 = logic.SortStr(addRangeStr)
			} else {
				thisModifyWorkFlowRang1 = logic.SortStr(addRangeStr + "," + thismodifyWorkFlow.Proportion)
			}
		}
	} else {
		//7、减少了流量占比，从该流量占比值中随机抽取减少部分，且还原给主流量
		count := dbThismodifyWorkFlowRangeArryLen - proportion
		addRangeStr, rangeStr := ModifyWorkFlowRange(dbThismodifyWorkFlowRangeArry, count)
		thisModifyWorkFlowRang1 = rangeStr
		if len(addRangeStr) <= 0 {
			MainWorkFlowRang1 = logic.SortStr(workFlow.Proportion)
		} else {
			MainWorkFlowRang1 = logic.SortStr(addRangeStr + "," + workFlow.Proportion)
		}
	}

	//8、校验修改操作后的流量总占比，是否为100
	err = logic.CheckModifyOP(haveUseProportion, MainWorkFlowRang1, thisModifyWorkFlowRang1)
	if err != nil && flowType == 0 {
		responseInfo.Tag = false
		responseInfo.ErrMsg = err.Error()
		return
	}

	//9、生成对应的区间表达式
	thisModifyWorkFlowRang2 := logic.CreateRangeStr(thisModifyWorkFlowRang1)
	MainWorkFlowRang2 := logic.CreateRangeStr(MainWorkFlowRang1)

	//10、判断是否需要将原主流量变成副流量
	mainDefult := "1"
	if thismodifyWorkFlow.Defult == "0" && modifyData.Defult == "主" {
		mainDefult = "0"
	}

	//11、更新需要修改的流量信息

	//强校验，最终入库的流量比例是否和修改传入的流量比例一致
	var thisModifyWorkFlowRang1Array []string
	if thisModifyWorkFlowRang1 != "" {
		thisModifyWorkFlowRang1Array = strings.Split(thisModifyWorkFlowRang1, ",")
	}
	if len(thisModifyWorkFlowRang1Array) != proportion && flowType == 0 {
		responseInfo.Tag = false
		responseInfo.ErrMsg = "最终入库的流量比例 和 修改传入的流量比例 不一致，入库失败，请重新更新！"
		return
	}

	updateModifySql := "update workflow set id=?,range1 = ?,range2 = ?,is_default = ?,modules =?,remark =?,update_time =?,operator=?,manual_slot_ids=?,group_name=?,flow_chart=? where id = ?"
	_, updateModifyErr := mysql.Handler.Exec(updateModifySql, modifyData.WorkFlowId, thisModifyWorkFlowRang1, thisModifyWorkFlowRang2, template.DefultNameAndIdMap[modifyData.Defult], modifyData.Modules, modifyData.Remark, nowTime, modifyData.Operator, modifyData.ManualSlotIds, modifyData.GroupName, modifyData.FlowChartJson, modifyData.OldWorkFlowId)
	if updateModifyErr != nil {
		responseInfo.Tag = false
		responseInfo.ErrMsg = updateModifyErr.Error()
		return
	}

	tlog.Handler.Infof(ctx, logs.DLTagProcessLog, "etype=workflowadmin.UpdateDBWorkFlow||data=%v||Rang1=%v||Rang2=%v||err=本条记录修改成功", modifyData, thisModifyWorkFlowRang1, thisModifyWorkFlowRang2)

	//12、更新主流量信息
	updateMainSql := "update workflow set range1 = ?,range2 = ?,is_default = ?,update_time =?,operator=?,manual_slot_ids=? where id = ?"
	_, updateMainErr := mysql.Handler.Exec(updateMainSql, MainWorkFlowRang1, MainWorkFlowRang2, mainDefult, nowTime, modifyData.Operator, workFlow.ManualSlotIds, workFlow.WorkFlowId)
	if updateMainErr != nil {
		responseInfo.Tag = false
		responseInfo.ErrMsg = updateMainErr.Error()
		return
	}
	tlog.Handler.Infof(ctx, logs.DLTagProcessLog, "etype=workflowadmin.UpdateDBWorkFlow||Rang1=%v||Rang2=%v||err=主流量修改成功", MainWorkFlowRang1, MainWorkFlowRang2)
}
