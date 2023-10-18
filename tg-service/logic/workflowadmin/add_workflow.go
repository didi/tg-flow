package workflowadmin

import (
	"github.com/didi/tg-flow/tg-core/common/mysql"
	"strconv"
	"strings"
	"tg-service/common/template"
	"tg-service/idl"
	"tg-service/logic"
	"tg-service/logic/sceneadmin"
	"time"
)

func AddWorkFlow(addData *idl.WorkFlowConfig, responseInfo *idl.ResponseInfo) {
	var flowType int64
	flowType = -1
	if addData.SceneId != 0 {
		sceneData := &idl.SceneConfig{
			Id: addData.SceneId,
		}
		sceneConfigList := sceneadmin.SelectSceneConf(sceneData, -1, 0, false)
		if len(sceneConfigList) != 0 {
			flowType = sceneConfigList[0].FlowType
		}
	}

	nowTime := time.Now().Format("2006-01-02 15:04:05")
	//新增流量占比
	proportion, _ := strconv.Atoi(addData.Proportion)

	//根据场景id，查询主流量占比,并更新
	var sceneId int
	sceneId, ok := template.SceneNameAndIdMap[addData.SceneName]
	if !ok {
		sceneId = 0
	}

	//查询该场景和维度下的信息
	sql := "select id,dimension_id,experiment_id, modules,is_default,range1,range2, remark ,create_time,update_time,operator,manual_slot_ids,group_name, flow_chart from workflow where status=1 and experiment_id = ? and dimension_id = ?"
	workFlowMap, err := SelectDBWorkFlow(sql, sceneId, addData.DimensionId)
	if err != nil {
		responseInfo.Tag = false
		responseInfo.ErrMsg = "添加数据失败:" + err.Error()
		return
	}

	//找出主流量
	defaultWorkFlow := new(idl.WorkFlow)
	for _, sceneMap := range workFlowMap { //同一个场景和维度下，主流量只会有一条记录
		for _, workFlow := range sceneMap {
			if workFlow.IsDefault == 1 {
				defaultWorkFlow = workFlow
				break
			}
		}
	}

	//新增的分桶流量占比离散数字(0-99)字符串
	var addRangeStr string
	//需要更新的主流量的分桶流量占比离散数字(0-99)字符串
	var defaultRangeStr string
	//主流量实验编号
	var defaultWorkFlowId int64

	defaultProportion := strings.Split(defaultWorkFlow.Range1, ",")

	defaultProportionLength := 0
	if len(defaultWorkFlow.Range1) > 0 { //排除主流量占比值为空的情况
		defaultProportionLength = len(defaultProportion)
	}

	if flowType == 0 {
		if defaultProportionLength < proportion {
			responseInfo.Tag = false
			responseInfo.ErrMsg = "新增的流量占比，超过了可分配的流量"
			return
		} else {
			defaultWorkFlowId = defaultWorkFlow.Id
			//生成新增流量和主流量占比的离散数字(0-99)字符串
			addRangeStr, defaultRangeStr = ModifyWorkFlowRange(defaultProportion, proportion)

			//校验添加操作后的流量总占比，是否为100
			err := logic.CheckAfterAddOP(strconv.FormatInt(defaultWorkFlow.ExperimentId, 10), addData.DimensionId, strconv.FormatInt(defaultWorkFlow.Id, 10), defaultRangeStr, addRangeStr)
			if err != nil {
				responseInfo.Tag = false
				responseInfo.ErrMsg = err.Error()
				return
			}
		}
	}

	/**************先入库新增流量成功了，才能去更新之前的主流量占比*********************/

	//if len(addRangeStr) <= 0 || addRangeStr == "" {
	//	responseInfo.Tag = false
	//	responseInfo.ErrMsg = "新增分桶流量分配占比，必须大于0"
	//	return
	//}

	var addSql string
	var addErr error
	addRange2 := logic.CreateRangeStr(addRangeStr)
	if addData.WorkFlowId == "" { //新增时，没有填写实验编号
		addSql = "insert into workflow(dimension_id,experiment_id, modules,is_default,status,range1,range2,remark,create_time,update_time,operator,manual_slot_ids,group_name,flow_chart) values(?,?,?,0,1,?,?,?,?,?,?,?,?,?)"
		_, addErr = mysql.Handler.Exec(addSql, addData.DimensionId, sceneId, addData.Modules, addRangeStr, addRange2, addData.Remark, nowTime, nowTime, addData.Operator, addData.ManualSlotIds, addData.GroupName, addData.FlowChartJson)
	} else {
		addSql = "insert into workflow(id,dimension_id,experiment_id, modules,is_default,status,range1,range2,remark,create_time,update_time,operator,manual_slot_ids,group_name,flow_chart) values(?,?,?,?,0,1,?,?,?,?,?,?,?,?,?)"
		_, addErr = mysql.Handler.Exec(addSql, addData.WorkFlowId, addData.DimensionId, sceneId, addData.Modules, addRangeStr, addRange2, addData.Remark, nowTime, nowTime, addData.Operator, addData.ManualSlotIds, addData.GroupName, addData.FlowChartJson)
	}

	if addErr != nil {
		responseInfo.Tag = false
		responseInfo.ErrMsg = addErr.Error()
		return
	}

	//更新主流量占比，修改数据库
	range2 := logic.CreateRangeStr(defaultRangeStr)
	updateSql := "update workflow set range1 = ? ,range2 = ? ,update_time = ?, operator =?, manual_slot_ids =? where id= ?"
	_, updateErr := mysql.Handler.Exec(updateSql, defaultRangeStr, range2, nowTime, addData.Operator, defaultWorkFlow.ManualSlotIds, defaultWorkFlowId)
	if updateErr != nil {
		responseInfo.Tag = false
		responseInfo.ErrMsg = updateErr.Error()
		return
	}
}
