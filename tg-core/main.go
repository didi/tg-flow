package main

import (
	"fmt"
)

/*func main(){
	path := "/Users/didi/gopath/src/github.com/didi/tg-flow/tg-core/wfengine/workflow.json"
	flowStr, err := readStrFromFile(path)
	if err != nil {
		fmt.Println("flowStr,err", flowStr, err)
	}
	//fmt.Println("flowStr======>", fmt.Sprintf("%+v",flowStr))
	fc, err := wfengine.NewWorkflowChart(flowStr)
	if err != nil {
		fmt.Println("create WorkflowChart error:", err)
	}

	fcm, err := json.Marshal(fc)
	fmt.Println("fcm===============>", fmt.Sprintf("%+v", string(fcm)))

	action := getAction()
	params :=make([]interface{}, 2)
	params[0] = 4.4
	params[1] = "4.4"
	ce := wfengine.GetCondExecutors()
	res, err := ce.Execute(action.ActionName, params)

	fmt.Println("res, err=====>", res, err)

	actionIn := getInAction()
	paramsIn := make([]interface{}, 2)
	paramsIn[0] = "testc"
	coll := make(map[string]interface{})
	coll["testa"] = "vala"
	coll["testb"] = "valb"
	paramsIn[1] = coll

	ceIn := wfengine.GetCondExecutors()
	var cde interface{} = &test.ConditionExternal{}
	ceIn.RegisterCondExecutor("ConditionExternal", cde)
	fmt.Println("actionName,paramsIn====>", actionIn.ActionName, paramsIn)
	resIn, errIn := ceIn.Execute(actionIn.ActionName, paramsIn)

	fmt.Println("resIn, errIn=====>", resIn, errIn)
}

func getAction() *wfengine.Action {
	param0 := &wfengine.Param{
		Name:  "pl",
		Value: "$version",
		Type:  "string",
	}
	param1 := &wfengine.Param{
		Name:  "pr",
		Value: "6.0",
		Type:  "string",
	}

	action := &wfengine.Action{
		ActionType: "condition",
		ActionId:   "Action_11",
		ActionName: "NE",
		Params:     []*wfengine.Param{param0, param1},
	}
	return action
}

func getInAction() *wfengine.Action {
	param1 := &wfengine.Param{
		Name:  "Param1",
		Value: "$elem",
		Type:  "string",
	}
	param2 := &wfengine.Param{
		Name:  "Param2",
		Value: "$coll",
		Type:  "map[string]interface{}",
	}

	action := &wfengine.Action{
		ActionType: "condition",
		ActionId:   "Action_11",
		ActionName: "ConditionExternal.In",
		Params:     []*wfengine.Param{param1, param2},
	}
	return action
}

func readStrFromFile(path string)  (string, error){
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(f), nil
}
*/

type Action struct {
	ActionType		string				`json:"action_type"`
	ActionId		string				`json:"action_id"`
	ActionName		string				`json:"action_name"`
	NextActionIds	[]string			`json:"next_action_ids"`
	NextConditions	[]string			`json:"next_conditions"`
	PrevActionIds	[]string			`json:"prev_action_ids"`
	Description		string				`json:"description"`
}

func main() {

	//创建和初始化数组
	//使用简写声明
	action1 := Action{
		ActionType:     "action-1",
		ActionId:       "",
		ActionName:     "",
		NextActionIds:  nil,
		NextConditions: nil,
		PrevActionIds:  []string{"aaa","bbb","ccc"},
		Description:    "",
	}

	action2 := Action{
		ActionType:     "action-2",
		ActionId:       "",
		ActionName:     "",
		NextActionIds:  nil,
		NextConditions: nil,
		PrevActionIds:  action1.PrevActionIds,
		Description:    "",
	}

	fmt.Println("action1: ", action1)
	fmt.Println("action2:", action2)

	action1.PrevActionIds = nil
	fmt.Println("action11: ", action1)
	fmt.Println("action22:", action2)
}