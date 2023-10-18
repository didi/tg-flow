package logic

import (
	"context"
	"database/sql"
	"github.com/didi/tg-flow/tg-core/common/mysql"
	"github.com/didi/tg-flow/tg-core/common/tlog"
	"sort"
	"strconv"
	"strings"
	"tg-service/common/logs"
	"tg-service/idl"
)

//查询数据
func Select(selectData *idl.WorkFlowData) []*idl.WorkFlowData {
	var sqlStr string
	if selectData.ExperimentId == "" { //全量️查询
		sqlStr = "select id,experiment_id, modules,is_default,range1,remark,create_time,update_time,operator,manual_slot_ids from workflow where status=1"
	} else { //按场景id查询
		_, err := strconv.ParseInt(selectData.ExperimentId, 10, 64)
		if err != nil {
			return make([]*idl.WorkFlowData, 0)
		}
		sqlStr = "select id,experiment_id, modules,is_default,range1, remark,create_time,update_time,operator,manual_slot_ids from workflow where status=1 and experiment_id = ?"
	}
	return SelectWorkFlowPre(sqlStr, selectData.ExperimentId, 0)
}

func SelectWorkFlowPre(sqlStr string, experimentId string, dimendionId int64) []*idl.WorkFlowData {
	//查询数据库
	workFlowMap := SelectWorkFlow(sqlStr, experimentId, dimendionId)
	//组装返回结果
	workFlowDataList := Process(workFlowMap)
	return workFlowDataList
}

//数据库查询
func SelectWorkFlow(sqlStr string, id string, dimendionId int64) map[int64]map[int64]*idl.WorkFlow {
	workFlowMap := make(map[int64]map[int64]*idl.WorkFlow)
	stmt, err := mysql.Handler.Prepare(sqlStr)
	if err != nil {
		tlog.Handler.Errorf(context.TODO(), logs.DLTagLoadExperimentWorkflowFail, "etype=logic_SelectWorkFlow_Prepare||econtent=sql:%v||err=%v", sqlStr, err)
		return workFlowMap
	}
	defer stmt.Close()

	var rows *sql.Rows
	if id == "" {
		rows, err = stmt.Query()
	} else {
		rows, err = stmt.Query(id, dimendionId)
	}

	if err != nil {
		tlog.Handler.Errorf(context.TODO(), logs.DLTagLoadExperimentWorkflowFail, "etype=logic_SelectWorkFlow_Query||econtent=sql:%v||err=%v", sqlStr, err)
		return workFlowMap
	}
	defer rows.Close()

	for rows.Next() {
		var workFlow = new(idl.WorkFlow)
		err := rows.Scan(&workFlow.Id, &workFlow.DimensionId, &workFlow.ExperimentId, &workFlow.Modules, &workFlow.IsDefault, &workFlow.Range1, &workFlow.Remark, &workFlow.CreateTime, &workFlow.UpdateTime, &workFlow.Operator, &workFlow.ManualSlotIds)
		if err != nil {
			tlog.Handler.Errorf(context.TODO(), logs.DLTagLoadExperimentWorkflowFail, "etype=logic_workflowrowsScan||econtent=rows:%v||err=%v", rows, err)
			return workFlowMap
		}
		tempMap := workFlowMap[workFlow.ExperimentId]
		if tempMap == nil {
			tempMap = make(map[int64]*idl.WorkFlow)
			workFlowMap[workFlow.ExperimentId] = tempMap
		}
		tempMap[workFlow.Id] = workFlow
	}
	err = rows.Err()
	if err != nil {
		tlog.Handler.Errorf(context.TODO(), logs.DLTagLoadExperimentWorkflowFail, "etype=logic_workflowrowsErr||econtent=||err=%v", err)
		return workFlowMap
	}
	return workFlowMap
}

//对查询的数据库数据整理
func Process(workFlowMap map[int64]map[int64]*idl.WorkFlow) []*idl.WorkFlowData {
	workFlowDataList := make([]*idl.WorkFlowData, 0)
	if workFlowMap == nil || len(workFlowMap) <= 0 {
		return workFlowDataList
	}

	keys := make([]int, 0)
	for key := range workFlowMap {
		keys = append(keys, int(key))
	}
	// 给key排序，从小到大
	sort.Sort(sort.IntSlice(keys))

	for _, key := range keys {
		tempMap := workFlowMap[int64(key)]
		if len(tempMap) <= 0 {
			continue
		}
		tempKeys := make([]int, 0)
		for tempkey := range tempMap {
			tempKeys = append(tempKeys, int(tempkey))
		}
		// 给key排序，从小到大
		sort.Sort(sort.IntSlice(tempKeys))

		for _, k := range tempKeys {
			value := tempMap[int64(k)]
			var data = new(idl.WorkFlowData)
			data.ExperimentId = strconv.FormatInt(value.ExperimentId, 10)
			data.WorkFlowId = strconv.FormatInt(value.Id, 10)
			data.Modules = value.Modules
			data.Defult = strconv.Itoa(value.IsDefault)
			data.Remark = value.Remark
			data.CreateTime = value.CreateTime
			data.UpdateTime = value.UpdateTime
			data.Operator = value.Operator
			data.ManualSlotIds = value.ManualSlotIds
			proportionStr := strings.Split(value.Range1, ",")
			if len(proportionStr) == 1 && proportionStr[0] == "" {
				data.Proportion = "0%"
			} else {
				data.Proportion = strconv.Itoa(len(proportionStr)) + "%"
			}
			workFlowDataList = append(workFlowDataList, data)
		}
	}
	return workFlowDataList
}
