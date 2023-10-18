package workflowadmin

import (
	"context"
	"database/sql"
	"github.com/didi/tg-flow/tg-core/common/mysql"
	"github.com/didi/tg-flow/tg-core/common/tlog"
	"sort"
	"strconv"
	"strings"
	"tg-service/common/logs"
	"tg-service/common/template"
	"tg-service/idl"
	"tg-service/logic"
	"time"
)

//查询workflow数据
func SelectWorkFlowConfig(selectData *idl.WorkFlowConfig, pageLimit int64, pageNum int64, useLimit bool) ([]*idl.WorkFlowConfig, error) {
	sqlStr := "select id, dimension_id, experiment_id, modules,is_default,range1,range2,remark,create_time,update_time,operator,manual_slot_ids,group_name, flow_chart from workflow where status=1"
	workFlowMap := make(map[int64]map[int64]*idl.WorkFlow)
	err := new(error)
	//查询数据库,根据前端传回的名称，获取到场景id，再查询数据库
	if strings.HasPrefix(selectData.SceneName, "全部") && len(strings.Split(selectData.SceneName, ",")) > 1 {
		sceneNameList := strings.Split(selectData.SceneName, ",")
		sqlStr += " and experiment_id in ("
		for _, sceneName := range sceneNameList[1:] {
			sqlStr += strconv.Itoa(template.SceneNameAndIdMap[sceneName]) + ","
		}
		sqlStr = sqlStr[:len(sqlStr)-1]
		sqlStr += ") and dimension_id = ?"

		if useLimit {
			sqlStr = logic.AddSelectLimit(sqlStr, pageLimit, pageNum)
		}

		workFlowMap, *err = SelectDBWorkFlow(sqlStr, -1, selectData.DimensionId)
	} else {
		sqlStr += " and experiment_id = ? and dimension_id = ?"

		if useLimit {
			sqlStr = logic.AddSelectLimit(sqlStr, pageLimit, pageNum)
		}

		workFlowMap, *err = SelectDBWorkFlow(sqlStr, template.SceneNameAndIdMap[selectData.SceneName], selectData.DimensionId)
	}

	//组装返回结果
	workFlowDataList := Process(workFlowMap)
	return workFlowDataList, *err
}

//数据库查询
func SelectDBWorkFlow(sqlStr string, sceneId int, dimensionId int64) (map[int64]map[int64]*idl.WorkFlow, error) {
	workFlowMap := make(map[int64]map[int64]*idl.WorkFlow)
	stmt, err := mysql.Handler.Prepare(sqlStr)
	if err != nil {
		tlog.Handler.Errorf(context.TODO(), logs.DLTagLoadExperimentWorkflowFail, "etype=workflowadmin.SelectWorkFlow||econtent=sql:%v||err=%v", sqlStr, err)
		return workFlowMap, err
	}
	defer stmt.Close()

	var rows *sql.Rows
	if sceneId == -1 {
		rows, err = stmt.Query(dimensionId)
	} else {
		rows, err = stmt.Query(sceneId, dimensionId)
	}

	if err != nil {
		tlog.Handler.Errorf(context.TODO(), logs.DLTagLoadExperimentWorkflowFail, "etype=workflowadmin.SelectWorkFlow||econtent=sql:%v||err=%v", sqlStr, err)
		return workFlowMap, err
	}
	defer rows.Close()

	for rows.Next() {
		var createTime time.Time
		var updateTime time.Time
		var workFlow = new(idl.WorkFlow)
		err := rows.Scan(&workFlow.Id, &workFlow.DimensionId, &workFlow.ExperimentId, &workFlow.Modules, &workFlow.IsDefault, &workFlow.Range1, &workFlow.Range2, &workFlow.Remark, &createTime, &updateTime, &workFlow.Operator, &workFlow.ManualSlotIds, &workFlow.GroupName, &workFlow.FlowCharts)
		if err != nil {
			tlog.Handler.Errorf(context.TODO(), logs.DLTagLoadExperimentWorkflowFail, "etype=workflowadmin.workflowrowsScan||econtent=rows:%v||err=%v", rows, err)
			return workFlowMap, err
		}
		tempMap := workFlowMap[workFlow.ExperimentId]
		if tempMap == nil {
			tempMap = make(map[int64]*idl.WorkFlow)
			workFlowMap[workFlow.ExperimentId] = tempMap
		}

		//时间格式转化
		workFlow.CreateTime = createTime.Format(logic.TIME_FORMAT)
		workFlow.UpdateTime = updateTime.Format(logic.TIME_FORMAT)

		tempMap[workFlow.Id] = workFlow
	}
	err = rows.Err()
	if err != nil {
		tlog.Handler.Errorf(context.TODO(), logs.DLTagLoadExperimentWorkflowFail, "etype=workflowadmin.workflowrowsErr||econtent=||err=%v", err)
		return workFlowMap, err
	}
	return workFlowMap, nil
}

//对查询的数据库数据整理
func Process(workFlowMap map[int64]map[int64]*idl.WorkFlow) []*idl.WorkFlowConfig {
	workFlowDataList := make([]*idl.WorkFlowConfig, 0)
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

			var data = new(idl.WorkFlowConfig)
			sceneId, _ := strconv.Atoi(strconv.FormatInt(value.ExperimentId, 10))
			data.SceneName = template.SceneIdAndNameMap[sceneId]
			data.SceneId = value.ExperimentId
			data.DimensionId = value.DimensionId
			data.WorkFlowId = strconv.FormatInt(value.Id, 10)
			data.Modules = value.Modules
			data.ShowModules = value.FlowCharts
			data.Defult = template.DefultIdAndNameMap[value.IsDefault]
			data.Remark = value.Remark
			data.Range1 = value.Range1
			data.Range2 = value.Range2
			data.CreateTime = value.CreateTime
			data.UpdateTime = value.UpdateTime
			data.Operator = value.Operator
			data.ManualSlotIds = value.ManualSlotIds
			data.GroupName = value.GroupName
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
