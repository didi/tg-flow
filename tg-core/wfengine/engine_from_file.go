/**
Description : loader of workflow config info from files
Author		: dayunzhangyunfeng@didiglobal.com
Date		: 2021-05-14
*/

package wfengine

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/didi/tg-flow/tg-core/common/tlog"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

const (
	defaultRange    = "0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25,26,27,28,29,30,31,32,33,34,35,36,37,38,39,40,41,42,43,44,45,46,47,48,49,50,51,52,53,54,55,56,57,58,59,60,61,62,63,64,65,66,67,68,69,70,71,72,73,74,75,76,77,78,79,80,81,82,83,84,85,86,87,88,89,90,91,92,93,94,95,96,97,98,99"
	versionFileName = "/version"
)

func GetLatestVersionFromFile(path string) (string, error) {
	jsonFile, err := os.Open(path + versionFileName)
	if err != nil {
		return "", fmt.Errorf("os.open error, file:%v, err:%v", path, err)
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return "", fmt.Errorf("ioutil.ReadAll error, byteValue:%v,err=%v", string(byteValue), err)
	}

	return string(byteValue), nil
}

func NewWorkflowEngineFromFile(moduleObj ModuleObjBase, configPath string) (*WorkflowEngine, error) {
	smMap, err := LoadSceneModuleMapFromFile(configPath + "/scene.json")
	if err != nil {
		return nil, err
	}

	//更新workflow
	wfMap, err := LoadWorkflowFromFile(configPath, smMap)
	if err != nil {
		return nil, err
	}

	version, err1 := GetLatestVersionFromFile(configPath)
	if err1 != nil {
		tlog.ErrorCount(context.TODO(), "GetLatestVersionFromRedis_err", fmt.Sprintf("configPath=%v, err=%v", configPath, err1))
	}

	resetWorkflows(wfMap)

	mbMap, err := createModelMap(moduleObj, wfMap)
	if err != nil {
		return nil, err
	}

	workfowEngine := newWorkflowEngine(smMap, wfMap, mbMap, version)
	return workfowEngine, nil
}

func LoadSceneModuleMapFromFile(path string) (map[int64]*SceneModule, error) {
	jsonFile, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("os.open error, file:%v, err:%v", path, err)
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, fmt.Errorf("ioutil.ReadAll error, byteValue:%v,err=%v", string(byteValue), err)
	}

	var sceneModuleMap map[int64]*SceneModule
	err = json.Unmarshal(byteValue, &sceneModuleMap)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal error, byteValue:%v,err=%v", string(byteValue), err)
	}

	return sceneModuleMap, nil
}

func LoadWorkflowFromFile(dirPath string, smMap map[int64]*SceneModule) (map[int64]*Workflow, error) {
	//1.遍历目录
	filePaths := make(map[string]string)
	err := getWorkflowFiles(dirPath, filePaths)
	if err != nil {
		return nil, err
	}

	//对每个文件逐个加到map
	wfMap := make(map[int64]*Workflow)
	for filePath, _ := range filePaths {
		wf, err := createWorkflowFromFile(filePath)
		if wf == nil || err != nil {
			tlog.ErrorCount(context.TODO(), "loadWorkflowFromFile_err", fmt.Sprintf("wf:%v,err:%v", wf, err))
			continue
		}
		wf.FlowCharts, err = NewWorkflowChart(wf.FlowChart)
		if err != nil {
			tlog.ErrorCount(context.TODO(), "NewWorkflowChart_err", fmt.Sprintf("wf:%v,err:%v", wf, err))
			continue
		}
		wfMap[wf.Id] = wf
	}

	return wfMap, nil
}

func getWorkflowFiles(path string, filePaths map[string]string) error {
	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}

	rd, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	for _, fi := range rd {
		if fi.IsDir() {
			getWorkflowFiles(path+fi.Name(), filePaths)
		} else if strings.HasPrefix(fi.Name(), "workflow-") {
			filePath := path + fi.Name()
			filePaths[filePath] = ""
		}
	}

	return nil
}

func createWorkflowFromFile(path string) (*Workflow, error) {
	jsonFile, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("os.open error, file:%v, err:%v", path, err)
	}
	defer jsonFile.Close()

	fileName := jsonFile.Name()
	idx := strings.Index(fileName, ".")
	fileName = fileName[:idx]
	strs := strings.Split(fileName, "-")
	if len(strs) < 3 {
		return nil, fmt.Errorf("invalid workflow file name:%v, it must start with :workflow-sceneId-workflowId", jsonFile.Name())
	}
	sceneId, err := strconv.ParseInt(strs[1], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid sceneId:%v in file name:%v,err=%v", strs[1], fileName, err)
	}
	workflowId, err := strconv.ParseInt(strs[2], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid workflowId:%v in file name:%v,err=%v", strs[2], fileName, err)
	}

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, fmt.Errorf("ioutil.ReadAll error, byteValue:%v,err=%v", string(byteValue), err)
	}
	flowChartStr := string(byteValue)

	workflow := &Workflow{
		Id:          workflowId,
		DimensionId: -1,
		SceneId:     sceneId,
		FlowChart:   flowChartStr,
		FlowCharts:  nil,
		FlowBranch:  nil,
		IsDefault:   1,
		Range1:      defaultRange,
		Range2:      "",
		Remark:      "",
		//UpdateTime:,
		GroupName: "",
	}

	return workflow, nil
}
