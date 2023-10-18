/**
Description	:	dispatcher of workflow engine
Author:			dayunzhangyunfeng@didiglobal.com
Date:			2021-07-20
*/
package dispatcher

import (
	"context"
	"fmt"
	trace "git.xiaojukeji.com/lego/context-go"
	"github.com/didi/tg-flow/tg-core/common/timeutils"
	"github.com/didi/tg-flow/tg-core/common/tlog"
	"github.com/didi/tg-flow/tg-core/common/utils"
	"github.com/didi/tg-flow/tg-core/model"
	"github.com/didi/tg-flow/tg-core/wfengine"
	"sync"
)

type Dispatcher interface {
	BuildRequest(ctx context.Context, requestParam interface{}) *model.StrategyContext
	BuildResponse(sc *model.StrategyContext) interface{}
	WriteLog(ctx context.Context, sc *model.StrategyContext) map[string]interface{}
	GetWorkflowEngine() *wfengine.WorkflowEngine
	GetPublicKey() string
	GetInterfaceName() string
}

func getInterfaceErrorTag(interfaceName string) string {
	return interfaceName + "_err"
}

/**最新逻辑
  moduleName:业务模块名称
  systemName:系统名称
  tags:public日志类型标记
*/
func DoStrategy(ctx context.Context, requestParam interface{}, d Dispatcher, tc *timeutils.TimeCoster) (interface{}, *model.StrategyContext) {
	errTag := getInterfaceErrorTag(d.GetInterfaceName())
	defer utils.RecoverPanic(ctx, errTag)

	//1、请求参数解析
	tc.StartSectionCount("BuildRequest")
	sc := d.BuildRequest(ctx, requestParam)
	tc.StopSectionCount("BuildRequest")

	//2. 执行业务逻辑
	d.GetWorkflowEngine().Run(ctx, sc)
	errMap := sc.GetErrorMap()
	errMap.Range(func(key, val interface{}) bool {
		tlog.ErrorCount(ctx, "WorkflowEngine.Run_err", fmt.Sprintf("workflowengine run error, key=%v, val=%v", key, val))
		return true
	})

	//3. 耗时
	resultMap := sc.GetModuleResultMap()
	resultMap.Range(func(key, val interface{}) bool {
		moduleResult, ok := val.(*model.ModuleResultInfo)
		if ok && moduleResult != nil {
			tc.SetSectionCount(moduleResult.StrategyName, moduleResult.CostTime)
			return true
		}
		return false
	})

	//3、组装返回结果
	tc.StartSectionCount("BuildResponse")
	responseInfo := d.BuildResponse(sc)
	tc.StopSectionCount("BuildResponse")

	//4、异步日志
	go writeLog(ctx, sc, d, errTag)

	return responseInfo, sc
}

// DoStrategyBatch 批量执行 Workflow 的接口，根据 requestParam 批量执行多次 Workflow，并将结果返回，返回结果的数量总是和请求的数量相同。
func DoStrategyBatch(ctx context.Context, requestParams []interface{}, d Dispatcher) ([]interface{}, []*model.StrategyContext, []*timeutils.TimeCoster) {
	errTag := getInterfaceErrorTag(d.GetInterfaceName())
	defer utils.RecoverPanic(ctx, errTag)

	responseInfos := make([]interface{}, len(requestParams))
	scs := make([]*model.StrategyContext, len(requestParams))
	tcs := make([]*timeutils.TimeCoster, len(requestParams))

	wg := sync.WaitGroup{}
	wg.Add(len(requestParams))

	for i, req := range requestParams {

		tcs[i] = timeutils.NewTimeCoster()

		go func(curIdx int, curReq interface{}) {

			tc := tcs[curIdx]
			tc.StartCount()

			defer func() {

				if err := recover(); err != nil {
					tlog.ErrorCount(ctx, errTag, fmt.Sprintf("Recover system panic : %v", err))
					// 如果某个出现 panic，停止计时
					for section := range tc.GetAllCounts() {
						tc.StopSectionCount(section)
					}
				} else {
					tc.StopCount()
				}

				wg.Done()

			}()

			//1、请求参数解析
			tc.StartSectionCount("BuildRequest")
			sc := d.BuildRequest(ctx, curReq)
			tc.StopSectionCount("BuildRequest")

			//2. 执行业务逻辑
			d.GetWorkflowEngine().Run(ctx, sc)
			errMap := sc.GetErrorMap()
			errMap.Range(func(key, val interface{}) bool {
				tlog.ErrorCount(ctx, "WorkflowEngine.BatchRun_err", fmt.Sprintf("workflowengine batch run error, key=%v, val=%v, reqNum=%v", key, val, curIdx))
				return true
			})

			//3. 耗时
			resultMap := sc.GetModuleResultMap()
			resultMap.Range(func(key, val interface{}) bool {
				moduleResult, ok := val.(*model.ModuleResultInfo)
				if ok && moduleResult != nil {
					tc.SetSectionCount(moduleResult.StrategyName, moduleResult.CostTime)
					return true
				}
				return false
			})

			//3、组装返回结果
			tc.StartSectionCount("BuildResponse")
			responseInfo := d.BuildResponse(sc)
			tc.StopSectionCount("BuildResponse")

			//4、异步日志，每个请求单独打日志
			go writeLog(ctx, sc, d, errTag)

			responseInfos[curIdx] = responseInfo
			scs[curIdx] = sc

		}(i, req)
	}

	wg.Wait()

	return responseInfos, scs, tcs
}

func GetCtxTrace(ctx context.Context) *trace.DefaultTrace {
	if ctxTrace, ok := trace.GetCtxTrace(ctx); ok {
		return ctxTrace
	} else {
		return trace.NewDefaultTrace()
	}
}

/****************************************************记录public日志时，调用该方法
         ctx: 上下文环境
          sc: base框架中存放的临时数据，可以传nil
   publicKey: 数据采集平台的唯一表名（ 建议保持格式: g_系统名_服务接口名 ）
           d: base框架中的处理流程的结构体，不同系统都实现base中这个结构体内的方法
     tagName: 数据采集平台的唯一标识（ 建议保持格式: 系统名_服务接口名, 注:本标识 全部大写 ）
*************************************************************************************/
func writeLog(ctx context.Context, sc *model.StrategyContext, d Dispatcher, tagName string) {
	defer utils.RecoverPanic(ctx, tagName)
	params := d.WriteLog(ctx, sc)
	if params == nil {
		return
	}

	//公有日志信息
	params["uid"] = sc.UserId
	params["scene_id"] = sc.SceneId
	params["is_rateLimit"] = sc.IsLimited
	params["workflow_id"] = sc.FlowId

	mergeLog(params)
	tlog.Handler.Public(ctx, d.GetPublicKey(), params, false)
}

//将要记录到日志中的内容，再做一次非nil的过滤
func mergeLog(pairs map[string]interface{}) {
	for k, v := range pairs {
		if v == nil {
			pairs[k] = "NULL"
		}
	}
}
