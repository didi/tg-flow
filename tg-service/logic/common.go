/**
    @Description:
    @Author:zhouzichun
    @Date:2023/5/8
**/

package logic

import (
	"container/list"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/didi/tg-flow/tg-core/common/tlog"
	"math/rand"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"tg-service/idl"
	"time"
)

const TIME_FORMAT = "2006-01-02 15:04:05"

var Mutex sync.Mutex

//解析请求参数
func ParsRequestParam(r *http.Request) *idl.WorkFlowData {
	data := &idl.WorkFlowData{
		ExperimentId:  CancelEmptyStr(r.FormValue("experimentid")),
		OldWorkFlowId: CancelEmptyStr(r.FormValue("oldworkflowid")),
		WorkFlowId:    CancelEmptyStr(r.FormValue("workflowid")),
		Modules:       CancelEmptyStr(r.FormValue("modules")),
		Proportion:    CancelEmptyStr(r.FormValue("proportion")),
		Defult:        CancelEmptyStr(r.FormValue("defult")),
		Remark:        CancelEmptyStr(r.FormValue("remark")),
		UserCookie:    CancelEmptyStr(r.FormValue("usercookie")),
	}
	return data
}

//解析workflow请求参数
func ParsRequsetWorkFlowParam(r *http.Request) (*idl.WorkFlowConfig, error) {
	err := CheckParam(r.FormValue("oldworkflowid"))
	if err != nil {
		return &idl.WorkFlowConfig{}, fmt.Errorf("实验编号需要是数字")
	}
	err = CheckParam(r.FormValue("workflowid"))
	if err != nil {
		return &idl.WorkFlowConfig{}, fmt.Errorf("新增实验编号需要是数字")
	}

	//去掉流量占比数字后面的“%”号
	proportion := strings.Split(CancelEmptyStr(r.FormValue("proportion")), "%")[0]
	err = CheckParam(proportion)
	if err != nil {
		return &idl.WorkFlowConfig{}, fmt.Errorf("流量占比需要是数字")
	}

	dimensionId, _ := strconv.ParseInt(r.FormValue("dimension_id"), 10, 64)
	dataJson := CancelEmptyStr(r.FormValue("dataJson"))
	var updateJson string
	if dataJson != "" {
		flowChart := new(idl.ChartG6)
		err = json.Unmarshal([]byte(dataJson), &flowChart)
		if err != nil {
			return &idl.WorkFlowConfig{}, fmt.Errorf("流程图不符合规范")
		}
		workflowChart := formatChartG6(flowChart)

		updateJsonByte, err := json.Marshal(workflowChart)
		updateJson = string(updateJsonByte)
		if err != nil {
			return &idl.WorkFlowConfig{}, fmt.Errorf("流程图不符合规范")
		}
	}

	data := &idl.WorkFlowConfig{
		SceneName:     CancelEmptyStr(r.FormValue("scenename")),
		OldWorkFlowId: CancelEmptyStr(r.FormValue("oldworkflowid")),
		WorkFlowId:    CancelEmptyStr(r.FormValue("workflowid")),
		DimensionId:   dimensionId,
		Modules:       CancelEmptyStr(r.FormValue("modules")),
		Proportion:    proportion,
		Defult:        CancelEmptyStr(r.FormValue("defult")),
		Remark:        CancelEmptyStr(r.FormValue("remark")),
		CreateTime:    r.FormValue("createtime"),
		Operator:      CancelEmptyStr(r.FormValue("operator")),
		ManualSlotIds: CancelEmptyStr(r.FormValue("manual_slot_ids")),
		GroupName:     r.FormValue("groupname"),
		FlowChartJson: updateJson,
	}
	return data, nil
}

//解析请求中的页码参数
func ParsRequestPageParam(r *http.Request) (pageLimit int64, pageNum int64) {
	pageNum, _ = strconv.ParseInt(r.FormValue("page_num"), 10, 64)
	pageLimit, _ = strconv.ParseInt(r.FormValue("page_limit"), 10, 64)
	return pageLimit, pageNum
}

func formatChartG6(saveData *idl.ChartG6) *idl.WorkflowChart {
	nodes := saveData.Nodes
	edges := saveData.Edges
	workflowChart := new(idl.WorkflowChart)
	actionMap := make(map[string]*idl.Action)
	typeMap := map[string]string{"rect": "task", "diamond": "condition", "circle": "flow", "clock": "timeout"}
	for _, node := range nodes {
		action := idl.Action{
			ActionType:     typeMap[node.NodeType],
			ActionId:       node.Id,
			ActionName:     node.Label,
			Params:         node.Params,
			Timeout:        node.Timeout,
			RefWorkflowId:  node.RefWorkflowId,
			TimeoutAsync:   node.TimeoutAsync,
			TimeoutDynamic: node.TimeoutDynamic,
			Location:       node.Location,
		}
		deleteIndex := list.New()
		for edgeIndex, edge := range edges {
			if edge.Source == node.Id {
				action.NextActionIds = append(action.NextActionIds, edge.Target)
				if action.ActionType == "condition" {
					action.NextConditions = append(action.NextConditions, edge.Label)
				}
				deleteIndex.PushBack(edgeIndex)
			}
		}
		for i := deleteIndex.Back(); i != nil; i = i.Prev() {
			index := i.Value.(int)
			edges = append(edges[:index], edges[index+1:]...)
		}
		actionMap[node.Id] = &action
	}
	workflowChart.ActionMap = actionMap
	return workflowChart

}

//检验请求参数
func CheckParam(param string) error {
	if param == "" {
		param = "0"
	}
	_, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return err
	}
	return nil
}

//去除换行符和空字符串
func CancelEmptyStr(str string) string {
	str = strings.Replace(str, "\n", "", -1)
	str = strings.Replace(str, " ", "", -1)
	return str
}

// EchoJSON json格式输出
func EchoJSON(w http.ResponseWriter, r *http.Request, data interface{}) {
	if cType := w.Header().Get("Content-Type"); cType == "" {
		w.Header().Set("Content-Type", "application/json")
	}
	b, err := json.Marshal(data)
	if err != nil {
		Echo(w, r, []byte(`{"errno":1, "errmsg":"`+err.Error()+`"}`))
	} else {
		Echo(w, r, b)
	}
}

func EchoToJSON(w http.ResponseWriter, r *http.Request, data interface{}) {
	if cType := w.Header().Get("Content-Type"); cType == "" {
		w.Header().Set("Content-Type", "application/json")
	}
	b, err := json.Marshal(data)
	if err != nil {
		EchoTo(w, r, []byte(`{"errno":1, "errmsg":"`+err.Error()+`"}`))
	} else {
		EchoTo(w, r, b)
	}
}

func EchoTo(w http.ResponseWriter, req *http.Request, body []byte) {

	if cType := w.Header().Get("Content-Type"); cType == "" {
		w.Header().Set("Content-Type", "text/html")
	}
	if req.Method == "options" {
		return
	}
	if origin := req.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization,X-Token")
	}
	w.Write(body)
}

// Echo 原始输出
func Echo(w http.ResponseWriter, req *http.Request, body []byte) {
	if cType := w.Header().Get("Content-Type"); cType == "" {
		w.Header().Set("Content-Type", "text/plain")
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(body)
}

//为Sql语句加上页码限制
func AddSelectLimit(sqlStr string, pageLimit int64, pageNum int64) (newSqlStr string) {
	startIndex := (pageNum - 1) * pageLimit
	newSqlStr = sqlStr + " limit " + strconv.FormatInt(startIndex, 10) + ", " + strconv.FormatInt(pageLimit, 10) + ";"
	tlog.Handler.Info(newSqlStr)
	return newSqlStr
}

//生成count个[start,end)结束的不重复的随机数
func GenerateRandomNumber(start int, end int, count int) []int {
	//范围检查
	if end < start || (end-start) < count {
		return nil
	}
	//存放结果的slice
	nums := make([]int, 0)
	//随机数生成器，加入时间戳保证每次生成的随机数不一样
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(nums) < count {
		//生成随机数
		num := r.Intn((end - start)) + start
		//查重
		exist := false
		for _, v := range nums {
			if v == num {
				exist = true
				break
			}
		}
		if !exist {
			nums = append(nums, num)
		}
	}
	return nums
}

//校验添加操作后，流量占比是否为100
func CheckAfterAddOP(experimentId string, dimendionId int64, mainWorkFlowId string, mainWorkFlowRange string, addRange string) error {
	sql := "select id,dimension_id,experiment_id, modules,is_default,range1, remark,create_time,update_time,operator,manual_slot_ids from workflow where status=1 and experiment_id = ? and dimension_id = ?"
	experimentMap := SelectWorkFlow(sql, experimentId, dimendionId)
	sum := 0
	for _, workFlowMap := range experimentMap {
		for workFlowId, workFlow := range workFlowMap {
			if strconv.FormatInt(workFlowId, 10) == mainWorkFlowId {
				continue
			}
			if workFlow.Range1 != "" {
				sum += len(strings.Split(workFlow.Range1, ","))
			}
		}
	}
	if len(mainWorkFlowRange) > 0 {
		if mainWorkFlowRange != "" {
			sum += len(strings.Split(mainWorkFlowRange, ","))
		}
	}

	if len(addRange) > 0 {
		sum += len(strings.Split(addRange, ","))
	}

	if sum > 100 {
		err := errors.New("添加操作造成该场景下所有实验的流量占比超过了100，添加操作失败.")
		return err
	}
	return nil
}

//根据字符串，生成区间表达式,如:"1,2,5,6,7,8,9" -> "[1-2][5-9]"
func CreateRangeStr(str string) string {
	resultStr := ""
	if len(str) <= 0 {
		return resultStr
	}
	list := make([]string, 0)
	addStr := strings.Split(str, ",")
	tempAddStr := addStr[0]
	for i := 1; i < len(addStr); i++ {
		before, _ := strconv.Atoi(addStr[i-1])
		after, _ := strconv.Atoi(addStr[i])
		if after-1 == before {
			if len(tempAddStr) <= 0 {
				tempAddStr += addStr[i]
			} else {
				tempAddStr += "," + addStr[i]
			}
		} else {
			list = append(list, tempAddStr)
			tempAddStr = addStr[i]
		}
	}
	list = append(list, tempAddStr)

	for _, v := range list {
		tempArry := strings.Split(v, ",")
		if len(tempArry) == 1 {
			resultStr += "[" + v + "]"
		} else {
			resultStr += "[" + tempArry[0] + "-" + tempArry[len(tempArry)-1] + "]"
		}
	}
	return resultStr
}

func CheckDeleteOP(mainWorkFlowRange string) error {
	if len(strings.Split(mainWorkFlowRange, ",")) > 100 {
		err := errors.New("删除操作造成该场景下主流量占比超过了100，删除操作失败.")
		return err
	}
	return nil
}

func CheckModifyOP(haveUseProportion int, mainWorkFlowRange string, addRange string) error {
	mainLen := 0
	if mainWorkFlowRange != "" {
		mainLen = len(strings.Split(mainWorkFlowRange, ","))

	}
	modifyLen := 0
	if addRange != "" {
		modifyLen = len(strings.Split(addRange, ","))
	}
	totalSum := haveUseProportion + mainLen + modifyLen
	if totalSum > 100 {
		err := errors.New("修改操作,造成该场景下所有实验的流量占比之和,超过了100，修改操作失败.")
		return err
	}
	return nil
}

//对字符串中的数字排序,如:“4,3,1,5,8,6” -> "1,3,4,5,6,8"
func SortStr(str string) string {
	// 去除空格
	str = strings.Replace(str, " ", "", -1)
	// 去除换行符
	str = strings.Replace(str, "\n", "", -1)
	resultStr := ""
	strs := strings.Split(str, ",")
	array := make([]int, 0)
	for _, v := range strs {
		temp, err := strconv.Atoi(v)
		if err == nil {
			array = append(array, temp)
		}
	}
	sort.Ints(array)
	for _, v := range array {
		if len(resultStr) <= 0 {
			resultStr += strconv.Itoa(v)
		} else {
			resultStr += "," + strconv.Itoa(v)
		}
	}
	return resultStr
}
