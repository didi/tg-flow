/**
  description: 统计结构上报工具
  后续将删除或抽象为接口供外部实现
**/

package metric

import (
	"context"
	"fmt"
	"github.com/didi/tg-flow/tg-core/common/timeutils"
	"github.com/didi/tg-flow/tg-core/common/tlog"
	"github.com/didi/tg-flow/tg-core/consts"
	statsd "go.intra.xiaojukeji.com/foundation/didi-standard-lib/metric/golang-v2/statsdlib"
	"time"
)

const (
	dlTagTimeCost	= " time_cost_ms"
	eTag			= "time_cost_ms"
)

/**
   metricName：统计的名称
   source：数据来源
   metricType：统计的类型
**/
func SendMetric(startTime time.Time, metricName string, source string, metricType string) {
	latency := time.Now().Sub(startTime)
	statsd.RpcMetric(metricName, source, metricType, latency, "ok")
}

/**
   metricName：统计的名称
   caller:		调用方
   callee:		被调用方
   startTime:	开始计时时间内
   code: 		调用结果,  取值 "ok" "0" "200" "201" "203"为成功、其他均为失败
**/
func RpcMetric(metricName string, caller string, callee string, startTime time.Time, code interface{}) {
	latency := time.Now().Sub(startTime)
	statsd.RpcMetric(metricName, caller, callee, latency, code)
}

//打印各个策略节点耗时，并上报metrics
func PrintSectionTime(startTime time.Time, traceId string, tc *timeutils.TimeCoster, metricsName string, sceneId int64) {
	tlog.Handler.Infof(context.TODO(), consts.DLTagControll, "etype=%v||%v||traceid=%v", eTag, tc.ToCountString(), traceId)

	//上报系统,包含所有场景,作为一个整体耗时给metrcis，便于配置报警策略
	go SendMetric(startTime, fmt.Sprintf("%v_total", metricsName), "all", "rt")

	//上报系统,分每个场景各自的耗时给metrcis，便于配置报警策略
	go SendMetric(startTime, fmt.Sprintf("%v_%v_total", metricsName, sceneId), "all", "rt")

	//系统内部各节点分段耗时上报metrics统计服务
	for key, useTime := range tc.GetAllCounts() {
		go statsd.RpcMetric(fmt.Sprintf("%v_%v", metricsName, key), "all", "rt", time.Millisecond*time.Duration(useTime), "ok")
	}
}

///////////////////////////////////////////////////////////////////////////////////////////////////
//上面是base库中老的metric封装, metricName由自己命名，无product和hintcode,将废弃
//下面是按点架构新出的规范做的的metric封装，用于替换上线老的metric封装。
//点架构metric规范: http://wiki.intra.xiaojukeji.com/pages/viewpage.action?pageId=609594742
///////////////////////////////////////////////////////////////////////////////////////////////////
func PrintAllMetrics(caller, callee string, sceneId int64, errCode, productId, hintCode string, traceId string, tc *timeutils.TimeCoster) {
	latency := time.Now().Sub(tc.GetStartTime())
	tlog.Handler.Infof(context.TODO(), dlTagTimeCost, "etype=%v||%v||traceid=%v", eTag, tc.ToCountString(), traceId)

	//接口
	tags := getTagMap(productId, hintCode)
	go SendRpcMetric("rpc_outer", caller, callee, latency, errCode, tags)

	//场景
	go SendRpcMetric("rpc_outer_scene", caller, fmt.Sprintf("scene_%v", sceneId), latency, errCode, nil)

	//节点
	for key, useTime := range tc.GetAllCounts() {
		go SendRpcMetric("rpc_outer_section", caller, key, time.Millisecond*time.Duration(useTime), errCode, nil)
	}
}

/**
	标签，注意返回内容必须不为空
 */
func getTagMap(productId, hintCode string) map[string]string {
	tags := make(map[string]string)
	if productId == "" {
		productId = "unknown"
	}
	if hintCode == "" {
		hintCode = "unknown"
	}
	tags["product"] = productId
	tags["hintcode"] = hintCode
	return tags
}
/**
	上报服务被调用的请求量、延时、错误率指标，打印耗时日志
*/
func PrintOuterMetric(caller, callee string, errCode, productId, hintCode string, traceId string, tc *timeutils.TimeCoster) {
	latency := time.Now().Sub(tc.GetStartTime())
	tlog.Handler.Infof(context.TODO(), dlTagTimeCost, "etype=%v||%v||traceid=%v", eTag, tc.ToCountString(), traceId)
	tags := getTagMap(productId, hintCode)
	SendRpcMetric("rpc_outer", caller, callee, latency, errCode, tags)
}

func SendRpcMetric(metricsName, caller, callee string, latency time.Duration, errCode string, tags map[string]string) {
	if len(tags) > 0 {
		statsd.RpcMetric(metricsName, caller, callee, latency, errCode, tags)
	}else{
		statsd.RpcMetric(metricsName, caller, callee, latency, errCode)
	}

}

/**
上报服务调用下游的请求量、延时、错误率指标
*/
func PrintAccessMetric(fun, callee string, latency time.Duration, errCode, productId, hintCode string) {
	tags := getTagMap(productId, hintCode)
	SendRpcMetric("rpc_access", fun, callee, latency, errCode, tags)
}
