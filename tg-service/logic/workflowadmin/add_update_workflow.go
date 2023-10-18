package workflowadmin

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"tg-service/common/template"
	"tg-service/idl"
	"tg-service/logic"
)

func AddOrUpdateWorkFlowConfig(addData *idl.WorkFlowConfig) *idl.ResponseInfo {
	responseInfo := &idl.ResponseInfo{
		Tag:    true,
		ErrMsg: "",
	}

	if addData.OldWorkFlowId == "" { //新增数据
		AddWorkFlow(addData, responseInfo)
	} else { //更新数据
		UpdateDBWorkFlow(addData, responseInfo)
	}
	return responseInfo
}

//查询指定场景id下的主流量、非主流量占比等信息
func SelectDBMainAndThisWorkFlow(data *idl.WorkFlowConfig) (*idl.WorkFlowConfig, *idl.WorkFlowConfig, int, string, error) {
	//根据场景名称，获取场景id
	var sceneId int
	sceneId, ok := template.SceneNameAndIdMap[data.SceneName]
	if !ok {
		return nil, nil, 0, "", fmt.Errorf("SceneNameAndIdMap not has sceneid=%v", data.SceneName)
	}
	selectSql := "select id,dimension_id,experiment_id, modules,is_default,range1,range2, remark , create_time ,update_time , operator, manual_slot_ids,group_name, flow_chart from workflow where status=1 and experiment_id = ? and dimension_id = ?"
	sceneMap, err := SelectDBWorkFlow(selectSql, sceneId, data.DimensionId)
	if err != nil {
		return nil, nil, 0, "", err
	}

	haveUseProportion := 0                        //不包括该非主流量的，其他非主流量的总占比百分数
	haveUseRange1 := ""                           //不包括该非主流量的，其他非主流量的总占比离散数字值
	var defaultWorkFlow = new(idl.WorkFlowConfig) //主流量
	var thisWorkFlow = new(idl.WorkFlowConfig)    //本条流量在数据库中的记录
	for _, workFlowMap := range sceneMap {
		for _, workFlow := range workFlowMap {
			if workFlow.IsDefault == 1 { //主流量
				defaultWorkFlow.WorkFlowId = strconv.FormatInt(workFlow.Id, 10)
				defaultWorkFlow.Modules = workFlow.Modules
				defaultWorkFlow.Proportion = workFlow.Range1
				defaultWorkFlow.Remark = workFlow.Remark
				defaultWorkFlow.ManualSlotIds = workFlow.ManualSlotIds
			} else {
				//未修改workflow id时，需要将数据库中的流量占比离散数字，赋值给请求的记录中，便于后续计算新的流量占比离散数字
				if strconv.FormatInt(workFlow.Id, 10) == data.WorkFlowId {
					data.Proportion = workFlow.Range1
				}
				//修改了workflow id时，需要将数据库中该条流量对应的原始数据返回
				if strconv.FormatInt(workFlow.Id, 10) == data.OldWorkFlowId {
					thisWorkFlow.WorkFlowId = strconv.FormatInt(workFlow.Id, 10)
					thisWorkFlow.Modules = workFlow.Modules
					thisWorkFlow.Proportion = workFlow.Range1
					thisWorkFlow.Remark = workFlow.Remark
					thisWorkFlow.Defult = strconv.Itoa(workFlow.IsDefault)
					thisWorkFlow.ManualSlotIds = workFlow.ManualSlotIds
				} else {
					//将其他非主流量的占比,离散数字加在一起
					if len(haveUseRange1) <= 0 {
						haveUseRange1 += workFlow.Range1
					} else {
						haveUseRange1 += "," + workFlow.Range1
					}
					//将其他非主流量的占比，百分数加在一起
					if workFlow.Range1 != "" {
						proportionStr := strings.Split(workFlow.Range1, ",")
						haveUseProportion += len(proportionStr)
					}
				}
			}
		}
	}
	return defaultWorkFlow, thisWorkFlow, haveUseProportion, haveUseRange1, nil
}

//调整主流量占比值和给新增流量分配占比离散数字(0-99)
func ModifyWorkFlowRange(proportionStr []string, proportion int) (string, string) {
	var addRangeStr string //新增的分桶流量占比
	var rangeStr string    //更新主流量占比
	index := 0
	//在主流量proportionStr中，随机挑选proportion个新增的流量占比值
	indexArry := logic.GenerateRandomNumber(0, len(proportionStr), proportion)
	//对indexArry数组从小到大排序
	sort.Ints(indexArry)
	for i := 0; i < len(proportionStr); i++ {
		if index < len(indexArry) && i == indexArry[index] {
			if len(addRangeStr) <= 0 {
				addRangeStr += proportionStr[i]
			} else {
				addRangeStr += "," + proportionStr[i]
			}
			index += 1
		} else {
			if len(rangeStr) <= 0 {
				rangeStr += proportionStr[i]
			} else {
				rangeStr += "," + proportionStr[i]
			}
		}
	}

	return addRangeStr, rangeStr
}
